// Package dispatch defines update routing, handlers, and middleware.
package dispatch

import (
	"context"
	stderrors "errors"
	"reflect"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// Handler handles one Telegram update.
type Handler interface {
	HandleUpdate(context.Context, telegram.Update) error
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(context.Context, telegram.Update) error

// HandleUpdate calls f(ctx, update).
func (f HandlerFunc) HandleUpdate(ctx context.Context, update telegram.Update) error {
	return f(ctx, update)
}

// Middleware wraps a Handler with additional behavior.
type Middleware func(Handler) Handler

// MiddlewareFunc adapts a function to Middleware.
type MiddlewareFunc = Middleware

// Predicate decides whether a route should handle an update.
type Predicate func(telegram.Update) bool

// ErrorHandler handles errors returned by route or fallback handlers.
type ErrorHandler func(context.Context, telegram.Update, error)

// Option configures a Dispatcher.
type Option func(*Dispatcher)

// Dispatcher routes Telegram updates to registered handlers.
type Dispatcher struct {
	routes       []route
	middleware   []Middleware
	fallback     Handler
	errorHandler ErrorHandler
}

type route struct {
	predicate Predicate
	handler   Handler
}

// New creates a Dispatcher.
func New(options ...Option) *Dispatcher {
	dispatcher := &Dispatcher{}
	for _, option := range options {
		if option != nil {
			option(dispatcher)
		}
	}

	return dispatcher
}

// WithFallback configures a fallback handler for unmatched updates.
func WithFallback(handler Handler) Option {
	return func(dispatcher *Dispatcher) {
		if dispatcher != nil && !isNilHandler(handler) {
			dispatcher.fallback = handler
		}
	}
}

// WithErrorHandler configures an error handler for route and fallback errors.
func WithErrorHandler(handler ErrorHandler) Option {
	return func(dispatcher *Dispatcher) {
		if dispatcher != nil {
			dispatcher.errorHandler = handler
		}
	}
}

// Use appends middleware to the dispatcher. Nil middleware is ignored.
func (d *Dispatcher) Use(middleware ...Middleware) {
	if d == nil {
		return
	}
	for _, item := range middleware {
		if item != nil {
			d.middleware = append(d.middleware, item)
		}
	}
}

// Handle registers a route.
func (d *Dispatcher) Handle(predicate Predicate, handler Handler) error {
	if d == nil {
		return stderrors.New("dispatcher is required")
	}
	if predicate == nil {
		return stderrors.New("predicate is required")
	}
	if isNilHandler(handler) {
		return stderrors.New("handler is required")
	}

	d.routes = append(d.routes, route{predicate: predicate, handler: handler})
	return nil
}

// HandleFunc registers a function route.
func (d *Dispatcher) HandleFunc(predicate Predicate, handler HandlerFunc) error {
	return d.Handle(predicate, handler)
}

// HandleUpdate routes update to the first matching handler.
func (d *Dispatcher) HandleUpdate(ctx context.Context, update telegram.Update) error {
	if ctx == nil {
		return stderrors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if d == nil {
		return stderrors.New("dispatcher is required")
	}

	for _, route := range d.routes {
		if route.predicate(update) {
			return d.dispatch(ctx, update, route.handler)
		}
	}
	if !isNilHandler(d.fallback) {
		return d.dispatch(ctx, update, d.fallback)
	}

	return nil
}

// OnMessage registers a handler for updates with message.
func (d *Dispatcher) OnMessage(handler Handler) error {
	return d.Handle(Message(), handler)
}

// OnMessageFunc registers a function handler for updates with message.
func (d *Dispatcher) OnMessageFunc(handler HandlerFunc) error {
	return d.OnMessage(handler)
}

// OnGuestMessage registers a handler for guest message updates.
func (d *Dispatcher) OnGuestMessage(handler Handler) error {
	return d.Handle(GuestMessage(), handler)
}

// OnGuestMessageFunc registers a function handler for guest message updates.
func (d *Dispatcher) OnGuestMessageFunc(handler HandlerFunc) error {
	return d.OnGuestMessage(handler)
}

// OnChannelPost registers a handler for channel post updates.
func (d *Dispatcher) OnChannelPost(handler Handler) error {
	return d.Handle(ChannelPost(), handler)
}

// OnChannelPostFunc registers a function handler for channel post updates.
func (d *Dispatcher) OnChannelPostFunc(handler HandlerFunc) error {
	return d.OnChannelPost(handler)
}

// OnEditedChannelPost registers a handler for edited channel post updates.
func (d *Dispatcher) OnEditedChannelPost(handler Handler) error {
	return d.Handle(EditedChannelPost(), handler)
}

// OnEditedChannelPostFunc registers a function handler for edited channel post updates.
func (d *Dispatcher) OnEditedChannelPostFunc(handler HandlerFunc) error {
	return d.OnEditedChannelPost(handler)
}

// OnCommand registers a handler for a slash command.
func (d *Dispatcher) OnCommand(command string, handler Handler) error {
	if !validCommand(command) {
		return stderrors.New("command must be non-empty, without slash or spaces")
	}

	return d.Handle(Command(command), handler)
}

// OnCommandFunc registers a function handler for a slash command.
func (d *Dispatcher) OnCommandFunc(command string, handler HandlerFunc) error {
	return d.OnCommand(command, handler)
}

// OnCallbackQuery registers a handler for callback query updates.
func (d *Dispatcher) OnCallbackQuery(handler Handler) error {
	return d.Handle(CallbackQuery(), handler)
}

// OnCallbackQueryFunc registers a function handler for callback query updates.
func (d *Dispatcher) OnCallbackQueryFunc(handler HandlerFunc) error {
	return d.OnCallbackQuery(handler)
}

// OnCallbackData registers a handler for callback query updates with matching data.
func (d *Dispatcher) OnCallbackData(data string, handler Handler) error {
	if data == "" {
		return stderrors.New("callback data is required")
	}

	return d.Handle(CallbackData(data), handler)
}

// OnCallbackDataFunc registers a function handler for callback query updates with matching data.
func (d *Dispatcher) OnCallbackDataFunc(data string, handler HandlerFunc) error {
	return d.OnCallbackData(data, handler)
}

// OnInlineQuery registers a handler for inline query updates.
func (d *Dispatcher) OnInlineQuery(handler Handler) error {
	return d.Handle(InlineQuery(), handler)
}

// OnInlineQueryFunc registers a function handler for inline query updates.
func (d *Dispatcher) OnInlineQueryFunc(handler HandlerFunc) error {
	return d.OnInlineQuery(handler)
}

// OnChosenInlineResult registers a handler for chosen inline result updates.
func (d *Dispatcher) OnChosenInlineResult(handler Handler) error {
	return d.Handle(ChosenInlineResult(), handler)
}

// OnChosenInlineResultFunc registers a function handler for chosen inline result updates.
func (d *Dispatcher) OnChosenInlineResultFunc(handler HandlerFunc) error {
	return d.OnChosenInlineResult(handler)
}

// OnShippingQuery registers a handler for shipping query updates.
func (d *Dispatcher) OnShippingQuery(handler Handler) error {
	return d.Handle(ShippingQuery(), handler)
}

// OnShippingQueryFunc registers a function handler for shipping query updates.
func (d *Dispatcher) OnShippingQueryFunc(handler HandlerFunc) error {
	return d.OnShippingQuery(handler)
}

// OnPreCheckoutQuery registers a handler for pre-checkout query updates.
func (d *Dispatcher) OnPreCheckoutQuery(handler Handler) error {
	return d.Handle(PreCheckoutQuery(), handler)
}

// OnPreCheckoutQueryFunc registers a function handler for pre-checkout query updates.
func (d *Dispatcher) OnPreCheckoutQueryFunc(handler HandlerFunc) error {
	return d.OnPreCheckoutQuery(handler)
}

// OnPaidMediaPurchased registers a handler for paid media purchase updates.
func (d *Dispatcher) OnPaidMediaPurchased(handler Handler) error {
	return d.Handle(PaidMediaPurchased(), handler)
}

// OnPaidMediaPurchasedFunc registers a function handler for paid media purchase updates.
func (d *Dispatcher) OnPaidMediaPurchasedFunc(handler HandlerFunc) error {
	return d.OnPaidMediaPurchased(handler)
}

// OnManagedBot registers a handler for managed bot updates.
func (d *Dispatcher) OnManagedBot(handler Handler) error {
	return d.Handle(ManagedBot(), handler)
}

// OnManagedBotFunc registers a function handler for managed bot updates.
func (d *Dispatcher) OnManagedBotFunc(handler HandlerFunc) error {
	return d.OnManagedBot(handler)
}

// OnBusinessConnection registers a handler for business connection updates.
func (d *Dispatcher) OnBusinessConnection(handler Handler) error {
	return d.Handle(BusinessConnection(), handler)
}

// OnBusinessConnectionFunc registers a function handler for business connection updates.
func (d *Dispatcher) OnBusinessConnectionFunc(handler HandlerFunc) error {
	return d.OnBusinessConnection(handler)
}

// OnBusinessMessage registers a handler for business message updates.
func (d *Dispatcher) OnBusinessMessage(handler Handler) error {
	return d.Handle(BusinessMessage(), handler)
}

// OnBusinessMessageFunc registers a function handler for business message updates.
func (d *Dispatcher) OnBusinessMessageFunc(handler HandlerFunc) error {
	return d.OnBusinessMessage(handler)
}

// OnEditedBusinessMessage registers a handler for edited business message updates.
func (d *Dispatcher) OnEditedBusinessMessage(handler Handler) error {
	return d.Handle(EditedBusinessMessage(), handler)
}

// OnEditedBusinessMessageFunc registers a function handler for edited business message updates.
func (d *Dispatcher) OnEditedBusinessMessageFunc(handler HandlerFunc) error {
	return d.OnEditedBusinessMessage(handler)
}

// OnDeletedBusinessMessages registers a handler for deleted business messages updates.
func (d *Dispatcher) OnDeletedBusinessMessages(handler Handler) error {
	return d.Handle(DeletedBusinessMessages(), handler)
}

// OnDeletedBusinessMessagesFunc registers a function handler for deleted business messages updates.
func (d *Dispatcher) OnDeletedBusinessMessagesFunc(handler HandlerFunc) error {
	return d.OnDeletedBusinessMessages(handler)
}

// OnPollAnswer registers a handler for poll answer updates.
func (d *Dispatcher) OnPollAnswer(handler Handler) error {
	return d.Handle(PollAnswer(), handler)
}

// OnPollAnswerFunc registers a function handler for poll answer updates.
func (d *Dispatcher) OnPollAnswerFunc(handler HandlerFunc) error {
	return d.OnPollAnswer(handler)
}

// OnPoll registers a handler for standalone poll updates.
func (d *Dispatcher) OnPoll(handler Handler) error {
	return d.Handle(Poll(), handler)
}

// OnPollFunc registers a function handler for standalone poll updates.
func (d *Dispatcher) OnPollFunc(handler HandlerFunc) error {
	return d.OnPoll(handler)
}

// OnChatJoinRequest registers a handler for chat join request updates.
func (d *Dispatcher) OnChatJoinRequest(handler Handler) error {
	return d.Handle(ChatJoinRequest(), handler)
}

// OnChatJoinRequestFunc registers a function handler for chat join request updates.
func (d *Dispatcher) OnChatJoinRequestFunc(handler HandlerFunc) error {
	return d.OnChatJoinRequest(handler)
}

// OnMessageReaction registers a handler for message reaction updates.
func (d *Dispatcher) OnMessageReaction(handler Handler) error {
	return d.Handle(MessageReaction(), handler)
}

// OnMessageReactionFunc registers a function handler for message reaction updates.
func (d *Dispatcher) OnMessageReactionFunc(handler HandlerFunc) error {
	return d.OnMessageReaction(handler)
}

// OnMessageReactionCount registers a handler for anonymous message reaction count updates.
func (d *Dispatcher) OnMessageReactionCount(handler Handler) error {
	return d.Handle(MessageReactionCount(), handler)
}

// OnMessageReactionCountFunc registers a function handler for anonymous message reaction count updates.
func (d *Dispatcher) OnMessageReactionCountFunc(handler HandlerFunc) error {
	return d.OnMessageReactionCount(handler)
}

// OnMyChatMember registers a handler for bot chat member status updates.
func (d *Dispatcher) OnMyChatMember(handler Handler) error {
	return d.Handle(MyChatMember(), handler)
}

// OnMyChatMemberFunc registers a function handler for bot chat member status updates.
func (d *Dispatcher) OnMyChatMemberFunc(handler HandlerFunc) error {
	return d.OnMyChatMember(handler)
}

// OnChatMember registers a handler for chat member status updates.
func (d *Dispatcher) OnChatMember(handler Handler) error {
	return d.Handle(ChatMember(), handler)
}

// OnChatMemberFunc registers a function handler for chat member status updates.
func (d *Dispatcher) OnChatMemberFunc(handler HandlerFunc) error {
	return d.OnChatMember(handler)
}

// OnChatBoost registers a handler for chat boost updates.
func (d *Dispatcher) OnChatBoost(handler Handler) error {
	return d.Handle(ChatBoost(), handler)
}

// OnChatBoostFunc registers a function handler for chat boost updates.
func (d *Dispatcher) OnChatBoostFunc(handler HandlerFunc) error {
	return d.OnChatBoost(handler)
}

// OnRemovedChatBoost registers a handler for removed chat boost updates.
func (d *Dispatcher) OnRemovedChatBoost(handler Handler) error {
	return d.Handle(RemovedChatBoost(), handler)
}

// OnRemovedChatBoostFunc registers a function handler for removed chat boost updates.
func (d *Dispatcher) OnRemovedChatBoostFunc(handler HandlerFunc) error {
	return d.OnRemovedChatBoost(handler)
}

// Any matches every update.
func Any() Predicate {
	return func(telegram.Update) bool { return true }
}

// Message matches updates with a message.
func Message() Predicate {
	return func(update telegram.Update) bool { return update.Message != nil }
}

// GuestMessage matches updates with a guest message.
func GuestMessage() Predicate {
	return func(update telegram.Update) bool { return update.GuestMessage != nil }
}

// ChannelPost matches updates with a channel post.
func ChannelPost() Predicate {
	return func(update telegram.Update) bool { return update.ChannelPost != nil }
}

// EditedChannelPost matches updates with an edited channel post.
func EditedChannelPost() Predicate {
	return func(update telegram.Update) bool { return update.EditedChannelPost != nil }
}

// Command matches updates with the given slash command.
func Command(command string) Predicate {
	return func(update telegram.Update) bool {
		return validCommand(command) && update.Message != nil && update.Message.IsCommand(command)
	}
}

// CallbackQuery matches updates with a callback query.
func CallbackQuery() Predicate {
	return func(update telegram.Update) bool { return update.CallbackQuery != nil }
}

// CallbackData matches callback query updates with exact data.
func CallbackData(data string) Predicate {
	return func(update telegram.Update) bool {
		return data != "" && update.CallbackQuery != nil && update.CallbackQuery.Data == data
	}
}

// InlineQuery matches updates with an inline query.
func InlineQuery() Predicate {
	return func(update telegram.Update) bool { return update.InlineQuery != nil }
}

// ChosenInlineResult matches updates with a chosen inline result.
func ChosenInlineResult() Predicate {
	return func(update telegram.Update) bool { return update.ChosenInlineResult != nil }
}

// ShippingQuery matches updates with a shipping query.
func ShippingQuery() Predicate {
	return func(update telegram.Update) bool { return update.ShippingQuery != nil }
}

// PreCheckoutQuery matches updates with a pre-checkout query.
func PreCheckoutQuery() Predicate {
	return func(update telegram.Update) bool { return update.PreCheckoutQuery != nil }
}

// PaidMediaPurchased matches updates with a paid media purchase.
func PaidMediaPurchased() Predicate {
	return func(update telegram.Update) bool { return update.PurchasedPaidMedia != nil }
}

// ManagedBot matches updates about managed bot creation or changes.
func ManagedBot() Predicate {
	return func(update telegram.Update) bool { return update.ManagedBot != nil }
}

// BusinessConnection matches updates about business account connections.
func BusinessConnection() Predicate {
	return func(update telegram.Update) bool { return update.BusinessConnection != nil }
}

// BusinessMessage matches new business message updates.
func BusinessMessage() Predicate {
	return func(update telegram.Update) bool { return update.BusinessMessage != nil }
}

// EditedBusinessMessage matches edited business message updates.
func EditedBusinessMessage() Predicate {
	return func(update telegram.Update) bool { return update.EditedBusinessMessage != nil }
}

// DeletedBusinessMessages matches deleted business messages updates.
func DeletedBusinessMessages() Predicate {
	return func(update telegram.Update) bool { return update.DeletedBusinessMessages != nil }
}

// PollAnswer matches updates with a poll answer.
func PollAnswer() Predicate {
	return func(update telegram.Update) bool { return update.PollAnswer != nil }
}

// Poll matches standalone poll state updates.
func Poll() Predicate {
	return func(update telegram.Update) bool { return update.Poll != nil }
}

// ChatJoinRequest matches updates with a chat join request.
func ChatJoinRequest() Predicate {
	return func(update telegram.Update) bool { return update.ChatJoinRequest != nil }
}

// MessageReaction matches updates with a message reaction change.
func MessageReaction() Predicate {
	return func(update telegram.Update) bool { return update.MessageReaction != nil }
}

// MessageReactionCount matches updates with anonymous message reaction count changes.
func MessageReactionCount() Predicate {
	return func(update telegram.Update) bool { return update.MessageReactionCount != nil }
}

// MyChatMember matches updates with the bot's chat member status changes.
func MyChatMember() Predicate {
	return func(update telegram.Update) bool { return update.MyChatMember != nil }
}

// ChatMember matches updates with chat member status changes.
func ChatMember() Predicate {
	return func(update telegram.Update) bool { return update.ChatMember != nil }
}

// ChatBoost matches updates with added or changed chat boosts.
func ChatBoost() Predicate {
	return func(update telegram.Update) bool { return update.ChatBoost != nil }
}

// RemovedChatBoost matches updates with removed chat boosts.
func RemovedChatBoost() Predicate {
	return func(update telegram.Update) bool { return update.RemovedChatBoost != nil }
}

// Chain wraps handler with middleware in the order it is provided.
func Chain(handler Handler, middleware ...Middleware) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		if middleware[i] != nil {
			handler = middleware[i](handler)
		}
	}

	return handler
}

func (d *Dispatcher) dispatch(ctx context.Context, update telegram.Update, handler Handler) error {
	wrapped, err := d.applyMiddleware(handler)
	if err == nil {
		err = wrapped.HandleUpdate(ctx, update)
	}
	if err == nil {
		return nil
	}
	if d.errorHandler != nil {
		d.errorHandler(ctx, update, err)
		return nil
	}

	return err
}

func (d *Dispatcher) applyMiddleware(handler Handler) (Handler, error) {
	wrapped := handler
	for i := len(d.middleware) - 1; i >= 0; i-- {
		middleware := d.middleware[i]
		if middleware == nil {
			continue
		}
		wrapped = middleware(wrapped)
		if isNilHandler(wrapped) {
			return nil, stderrors.New("middleware returned nil handler")
		}
	}

	return wrapped, nil
}

func isNilHandler(handler Handler) bool {
	if handler == nil {
		return true
	}

	value := reflect.ValueOf(handler)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func validCommand(command string) bool {
	return command != "" && !strings.HasPrefix(command, "/") && !strings.ContainsAny(command, " \t\n\r")
}
