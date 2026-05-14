// Example 11_transport_parity shows one app running through long polling or webhook.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/middleware"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
	"github.com/xDilettante/ai-gram/transport/webhook"
)

const (
	logComponent     = "transport_parity"
	transportPolling = "polling"
	transportWebhook = "webhook"
)

var parityAllowedUpdates = []string{"message"}

type appConfig struct {
	Transport   string
	WebhookURL  string
	SecretToken string
	ListenAddr  string
	WebhookPath string
}

type app struct {
	bot              *aigram.Bot
	dispatcher       *dispatch.Dispatcher
	accessController *exampleutil.AccessController
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
	config, err := appConfigFromEnv()
	if err != nil {
		return err
	}
	application, err := newApp(b)
	if err != nil {
		return err
	}

	switch config.Transport {
	case transportPolling:
		return runPolling(ctx, b, application.dispatcher)
	case transportWebhook:
		return runWebhook(ctx, b, application.dispatcher, config)
	default:
		return fmt.Errorf("unsupported transport %q", config.Transport)
	}
}

func appConfigFromEnv() (appConfig, error) {
	transport, err := parseTransport(exampleutil.OptionalEnv("AIGRAM_TRANSPORT", transportPolling))
	if err != nil {
		return appConfig{}, err
	}

	config := appConfig{
		Transport:   transport,
		WebhookURL:  strings.TrimSpace(exampleutil.OptionalEnv("AIGRAM_WEBHOOK_URL", "")),
		SecretToken: strings.TrimSpace(exampleutil.OptionalEnv("AIGRAM_WEBHOOK_SECRET", "")),
		ListenAddr:  exampleutil.OptionalEnv("AIGRAM_LISTEN_ADDR", ":8080"),
		WebhookPath: normalizePath(exampleutil.OptionalEnv("AIGRAM_WEBHOOK_PATH", "/webhook")),
	}
	if config.Transport == transportWebhook && config.WebhookURL == "" {
		return appConfig{}, errors.New("AIGRAM_WEBHOOK_URL is required when AIGRAM_TRANSPORT=webhook")
	}
	return config, nil
}

func parseTransport(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", transportPolling:
		return transportPolling, nil
	case transportWebhook:
		return transportWebhook, nil
	default:
		return "", errors.New("AIGRAM_TRANSPORT must be polling or webhook")
	}
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/webhook"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func newApp(b *aigram.Bot) (*app, error) {
	accessConfig, err := exampleutil.AccessConfigFromEnv()
	if err != nil {
		return nil, err
	}
	accessController := exampleutil.NewAccessController(accessConfig)

	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("%s event=handler_error update_id=%d err=%v", logComponent, update.UpdateID, err)
	}))
	application := &app{bot: b, dispatcher: dp, accessController: accessController}
	dp.Use(middleware.AccessWithPolicy(accessController, application.accessDenyHandler()))
	if err := application.registerRoutes(); err != nil {
		return nil, err
	}

	return application, nil
}

func (a *app) registerRoutes() error {
	if err := a.dispatcher.OnCommandFunc("start", a.helpHandler()); err != nil {
		return err
	}
	if err := a.dispatcher.OnCommandFunc("help", a.helpHandler()); err != nil {
		return err
	}
	if err := a.dispatcher.OnCommandFunc("status", a.statusHandler()); err != nil {
		return err
	}
	return a.dispatcher.OnMessageFunc(a.echoHandler())
}

func (a *app) helpHandler() dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "help")
		return a.reply(ctx, update, strings.Join([]string{
			"Transport parity example.",
			"",
			"The same dispatcher, middleware, and handlers run through polling or webhook.",
			"Set AIGRAM_TRANSPORT=polling or AIGRAM_TRANSPORT=webhook.",
			"",
			"Commands:",
			"/status - show shared handler status",
			"Any non-command text receives a shared-handler reply.",
		}, "\n"))
	}
}

func (a *app) statusHandler() dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "status")
		return a.reply(ctx, update, statusText(update, a.accessController.Mode()))
	}
}

func (a *app) echoHandler() dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || strings.TrimSpace(message.Text) == "" || message.Command() != "" {
			return nil
		}
		logSafeUpdate(update, "echo")
		return a.reply(ctx, update, "Shared handler received your message.\n"+statusText(update, a.accessController.Mode()))
	}
}

func (a *app) accessDenyHandler() func(context.Context, telegram.Update) error {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "access_denied")
		return a.reply(ctx, update, "Access denied. Configure AIGRAM_ACCESS_MODE, AIGRAM_ADMIN_USER_IDS, or AIGRAM_ALLOWED_CHAT_IDS.")
	}
}

func (a *app) reply(ctx context.Context, update telegram.Update, text string) error {
	message := update.EffectiveMessage()
	if message == nil {
		return nil
	}
	_, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:          aigram.ChatIDInt(message.Chat.ID),
		MessageThreadID: message.MessageThreadID,
		Text:            text,
		ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
	})
	return err
}

func runPolling(ctx context.Context, b *aigram.Bot, handler dispatch.Handler) error {
	if _, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{DropPendingUpdates: true}); err != nil {
		return fmt.Errorf("delete webhook before polling: %w", err)
	}
	runner, err := longpoll.New(b, handler, longpoll.Config{
		Timeout:        30,
		AllowedUpdates: parityAllowedUpdates,
		OnError: func(ctx context.Context, err error) {
			log.Printf("%s event=polling_error err=%v", logComponent, err)
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%s event=transport_started transport=%s", logComponent, transportPolling)
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Printf("%s event=transport_stopped transport=%s", logComponent, transportPolling)
	return nil
}

func runWebhook(ctx context.Context, b *aigram.Bot, handler dispatch.Handler, config appConfig) error {
	receiver, err := webhook.New(handler, webhook.Config{
		SecretToken: config.SecretToken,
		OnError: func(ctx context.Context, update *telegram.Update, err error) {
			updateID := int64(0)
			if update != nil {
				updateID = update.UpdateID
			}
			log.Printf("%s event=webhook_handler_error update_id=%d err=%v", logComponent, updateID, err)
		},
	})
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    config.ListenAddr,
		Handler: mux(receiver, config.WebhookPath),
	}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("%s event=http_listen addr=%s webhook_path=%s", logComponent, config.ListenAddr, config.WebhookPath)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	ok, err := b.SetWebhook(ctx, aigram.SetWebhookParams{
		URL:            config.WebhookURL,
		SecretToken:    config.SecretToken,
		AllowedUpdates: parityAllowedUpdates,
	})
	if err != nil {
		shutdownHTTP(server)
		<-serverErr
		return fmt.Errorf("set webhook: %w", err)
	}
	log.Printf("%s event=transport_started transport=%s webhook_registered=%t", logComponent, transportWebhook, ok)

	select {
	case <-ctx.Done():
		shutdownHTTP(server)
		<-serverErr
		log.Printf("%s event=transport_stopped transport=%s", logComponent, transportWebhook)
		return nil
	case err := <-serverErr:
		return err
	}
}

func mux(receiver http.Handler, webhookPath string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(webhookPath, receiver)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

func shutdownHTTP(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("%s event=http_shutdown_error err=%v", logComponent, err)
	}
}

func statusText(update telegram.Update, accessMode middleware.AccessMode) string {
	chatLabel := "none"
	actorLabel := "unknown"
	if chat := update.EffectiveChat(); chat != nil {
		chatLabel = exampleutil.MaskInt64(chat.ID) + ":" + safeValue(chat.Type)
	}
	if actor := update.Actor(); !actor.IsZero() {
		actorLabel = actorText(actor)
	}
	return strings.Join([]string{
		"Shared app status:",
		"access_mode=" + string(accessMode),
		"chat=" + chatLabel,
		"actor=" + actorLabel,
	}, "\n")
}

func actorText(actor telegram.Actor) string {
	if actor.User != nil {
		return "user:" + exampleutil.MaskInt64(actor.User.ID)
	}
	if actor.Chat != nil {
		return "chat:" + exampleutil.MaskInt64(actor.Chat.ID)
	}
	return "unknown"
}

func logSafeUpdate(update telegram.Update, action string) {
	chatID := "0"
	actorUserID := "0"
	actorChatID := "0"
	command := ""
	messageID := int64(0)
	if chat := update.EffectiveChat(); chat != nil {
		chatID = exampleutil.MaskInt64(chat.ID)
	}
	actor := update.Actor()
	if actor.User != nil {
		actorUserID = exampleutil.MaskInt64(actor.User.ID)
	}
	if actor.Chat != nil {
		actorChatID = exampleutil.MaskInt64(actor.Chat.ID)
	}
	if message := update.EffectiveMessage(); message != nil {
		messageID = message.MessageID
		command = message.Command()
	}

	log.Printf("%s event=update action=%s update_id=%d message_id=%d chat_id=%s actor_user_id=%s actor_chat_id=%s command=%s",
		logComponent, action, update.UpdateID, messageID, chatID, actorUserID, actorChatID, command)
}

func safeValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "<none>"
	}
	return strings.NewReplacer("\n", " ", "\r", " ", "\t", " ").Replace(value)
}
