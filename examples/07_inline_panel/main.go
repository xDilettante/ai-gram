// Example 07_inline_panel shows a production-style inline panel with typed callbacks.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/callback"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
)

const (
	panelNamespace = "ops_panel"
	actionPage     = "page"
	actionSelect   = "select"
	pageSize       = 3
)

type panelItem struct {
	ID      string
	Title   string
	Details string
}

var panelItems = []panelItem{
	{ID: "queue", Title: "Queue status", Details: "Inspect pending background jobs."},
	{ID: "alerts", Title: "Alert routing", Details: "Review alert destinations and silence windows."},
	{ID: "reports", Title: "Reports", Details: "Generate the next reporting bundle."},
	{ID: "limits", Title: "Rate limits", Details: "Inspect current sending limits."},
	{ID: "deploy", Title: "Deploy window", Details: "Prepare a deploy window checklist."},
	{ID: "health", Title: "Health checks", Details: "Review external service health."},
	{ID: "audit", Title: "Audit log", Details: "Open the latest operational audit summary."},
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := exampleutil.SignalContext()
	defer stop()

	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}
	if _, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{DropPendingUpdates: true}); err != nil {
		return fmt.Errorf("delete webhook before long polling: %w", err)
	}

	dp, err := newDispatcher(b)
	if err != nil {
		return err
	}

	runner, err := longpoll.New(b, dp, longpoll.Config{
		Timeout:        30,
		AllowedUpdates: []string{"message", "callback_query"},
		OnError: func(ctx context.Context, err error) {
			log.Printf("longpoll error err=%v", err)
		},
	})
	if err != nil {
		return err
	}

	log.Println("inline panel bot started; press Ctrl+C to stop")
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Println("inline panel bot stopped")
	return nil
}

func newDispatcher(b *aigram.Bot) (*dispatch.Dispatcher, error) {
	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("handler error update_id=%d err=%v", update.UpdateID, err)
	}))

	if err := dp.OnCommandFunc("start", showPanel(b)); err != nil {
		return nil, err
	}
	if err := dp.OnCommandFunc("panel", showPanel(b)); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackActionFunc(panelNamespace, actionPage, handlePage(b)); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackActionFunc(panelNamespace, actionSelect, handleSelect(b)); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackActionFunc(panelNamespace, callback.ActionConfirm, handleConfirm(b)); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackActionFunc(panelNamespace, callback.ActionCancel, handleCancel(b)); err != nil {
		return nil, err
	}
	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || strings.HasPrefix(strings.TrimSpace(message.Text), "/") {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(message.Chat.ID),
			Text:   "Use /panel to open the inline panel.",
		})
		return err
	}); err != nil {
		return nil, err
	}

	return dp, nil
}

func showPanel(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		keyboard := panelKeyboard(0)
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:      aigram.ChatIDInt(message.Chat.ID),
			Text:        panelText(0),
			ReplyMarkup: &keyboard,
		})
		return err
	}
}

func handlePage(b *aigram.Bot) dispatch.CallbackDataHandlerFunc {
	return func(ctx context.Context, update telegram.Update, data callback.Data) error {
		query := update.CallbackQuery
		if query == nil {
			return nil
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: query.ID}); err != nil {
			return err
		}
		if query.Message == nil {
			return nil
		}

		page := pageFromData(data)
		keyboard := panelKeyboard(page)
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      editTarget(query),
			Text:        panelText(page),
			ReplyMarkup: &keyboard,
		})
		return err
	}
}

func handleSelect(b *aigram.Bot) dispatch.CallbackDataHandlerFunc {
	return func(ctx context.Context, update telegram.Update, data callback.Data) error {
		query := update.CallbackQuery
		if query == nil {
			return nil
		}
		item, ok := findItem(data.ID)
		if !ok {
			_, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Unknown panel item",
				ShowAlert:       true,
			})
			return err
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            "Confirm or cancel the selected action.",
		}); err != nil {
			return err
		}
		if query.Message == nil {
			return nil
		}

		keyboard := confirmKeyboard(item.ID, pageFromData(data))
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      editTarget(query),
			Text:        confirmText(item),
			ReplyMarkup: &keyboard,
		})
		return err
	}
}

func handleConfirm(b *aigram.Bot) dispatch.CallbackDataHandlerFunc {
	return func(ctx context.Context, update telegram.Update, data callback.Data) error {
		query := update.CallbackQuery
		if query == nil {
			return nil
		}
		item, ok := findItem(data.ID)
		if !ok {
			_, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Unknown panel item",
				ShowAlert:       true,
			})
			return err
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            "Confirmed",
		}); err != nil {
			return err
		}
		if query.Message == nil {
			return nil
		}

		page := pageFromData(data)
		keyboard := panelKeyboard(page)
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      editTarget(query),
			Text:        fmt.Sprintf("Confirmed: %s\n\n%s\n\nThis example only edits the message; put real side effects behind your own authorization and audit checks.", item.Title, panelText(page)),
			ReplyMarkup: &keyboard,
		})
		return err
	}
}

func handleCancel(b *aigram.Bot) dispatch.CallbackDataHandlerFunc {
	return func(ctx context.Context, update telegram.Update, data callback.Data) error {
		query := update.CallbackQuery
		if query == nil {
			return nil
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            "Cancelled",
		}); err != nil {
			return err
		}
		if query.Message == nil {
			return nil
		}

		page := pageFromData(data)
		keyboard := panelKeyboard(page)
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      editTarget(query),
			Text:        "Cancelled.\n\n" + panelText(page),
			ReplyMarkup: &keyboard,
		})
		return err
	}
}

func panelText(page int) string {
	page = clampPage(page)
	start, end := pageBounds(page)
	var builder strings.Builder
	fmt.Fprintf(&builder, "Operations panel\nPage %d/%d\n\n", page+1, totalPages())
	for i, item := range panelItems[start:end] {
		fmt.Fprintf(&builder, "%d. %s\n%s\n", start+i+1, item.Title, item.Details)
		if i < end-start-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func panelKeyboard(page int) telegram.InlineKeyboardMarkup {
	page = clampPage(page)
	start, end := pageBounds(page)
	rows := make([][]telegram.InlineKeyboardButton, 0, pageSize+1)
	for _, item := range panelItems[start:end] {
		rows = append(rows, []telegram.InlineKeyboardButton{
			callback.MustButton(item.Title, callback.New(panelNamespace, actionSelect).WithID(item.ID).WithPage(page)),
		})
	}

	var nav []telegram.InlineKeyboardButton
	if page > 0 {
		nav = append(nav, callback.MustButton("Previous", callback.ForPage(panelNamespace, actionPage, page-1)))
	}
	if page < totalPages()-1 {
		nav = append(nav, callback.MustButton("Next", callback.ForPage(panelNamespace, actionPage, page+1)))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	return telegram.NewInlineKeyboard(rows...)
}

func confirmKeyboard(itemID string, page int) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboard(
		[]telegram.InlineKeyboardButton{
			callback.MustButton("Confirm", callback.Confirm(panelNamespace, itemID).WithPage(page)),
			callback.MustButton("Cancel", callback.Cancel(panelNamespace, itemID).WithPage(page)),
		},
	)
}

func confirmText(item panelItem) string {
	return fmt.Sprintf("Confirm action\n\n%s\n%s", item.Title, item.Details)
}

func editTarget(query *telegram.CallbackQuery) aigram.EditMessageTarget {
	return aigram.EditTargetChat(aigram.ChatIDInt(query.Message.Chat.ID), query.Message.MessageID)
}

func findItem(id string) (panelItem, bool) {
	for _, item := range panelItems {
		if item.ID == id {
			return item, true
		}
	}
	return panelItem{}, false
}

func pageFromData(data callback.Data) int {
	if !data.HasPage {
		return 0
	}
	return clampPage(data.Page)
}

func pageBounds(page int) (int, int) {
	page = clampPage(page)
	start := page * pageSize
	end := start + pageSize
	if end > len(panelItems) {
		end = len(panelItems)
	}
	return start, end
}

func clampPage(page int) int {
	if page < 0 {
		return 0
	}
	last := totalPages() - 1
	if page > last {
		return last
	}
	return page
}

func totalPages() int {
	return int(math.Ceil(float64(len(panelItems)) / float64(pageSize)))
}
