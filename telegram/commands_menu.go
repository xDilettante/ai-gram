package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	botCommandScopeDefaultType               = "default"
	botCommandScopeAllPrivateChatsType       = "all_private_chats"
	botCommandScopeAllGroupChatsType         = "all_group_chats"
	botCommandScopeAllChatAdministratorsType = "all_chat_administrators"
	botCommandScopeChatType                  = "chat"
	botCommandScopeChatAdministratorsType    = "chat_administrators"
	botCommandScopeChatMemberType            = "chat_member"

	menuButtonCommandsType = "commands"
	menuButtonWebAppType   = "web_app"
	menuButtonDefaultType  = "default"
)

// BotCommand represents a bot command shown in Telegram clients.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// BotCommandScope marks Telegram bot command scope objects.
type BotCommandScope interface {
	botCommandScope()
}

// BotCommandScopeDefault represents commands used when no narrower scope matches.
type BotCommandScopeDefault struct {
	Type string `json:"type"`
}

// BotCommandScopeAllPrivateChats represents commands for all private chats.
type BotCommandScopeAllPrivateChats struct {
	Type string `json:"type"`
}

// BotCommandScopeAllGroupChats represents commands for all group and supergroup chats.
type BotCommandScopeAllGroupChats struct {
	Type string `json:"type"`
}

// BotCommandScopeAllChatAdministrators represents commands for all group and supergroup administrators.
type BotCommandScopeAllChatAdministrators struct {
	Type string `json:"type"`
}

// BotCommandScopeChat represents commands for a specific chat.
type BotCommandScopeChat struct {
	Type   string `json:"type"`
	ChatID any    `json:"chat_id"`
}

// BotCommandScopeChatAdministrators represents commands for administrators of a specific chat.
type BotCommandScopeChatAdministrators struct {
	Type   string `json:"type"`
	ChatID any    `json:"chat_id"`
}

// BotCommandScopeChatMember represents commands for a specific member of a group or supergroup.
type BotCommandScopeChatMember struct {
	Type   string `json:"type"`
	ChatID any    `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

func (BotCommandScopeDefault) botCommandScope()               {}
func (BotCommandScopeAllPrivateChats) botCommandScope()       {}
func (BotCommandScopeAllGroupChats) botCommandScope()         {}
func (BotCommandScopeAllChatAdministrators) botCommandScope() {}
func (BotCommandScopeChat) botCommandScope()                  {}
func (BotCommandScopeChatAdministrators) botCommandScope()    {}
func (BotCommandScopeChatMember) botCommandScope()            {}

// NewBotCommandScopeDefault creates the default bot command scope.
func NewBotCommandScopeDefault() BotCommandScopeDefault {
	return BotCommandScopeDefault{Type: botCommandScopeDefaultType}
}

// NewBotCommandScopeAllPrivateChats creates a scope for all private chats.
func NewBotCommandScopeAllPrivateChats() BotCommandScopeAllPrivateChats {
	return BotCommandScopeAllPrivateChats{Type: botCommandScopeAllPrivateChatsType}
}

// NewBotCommandScopeAllGroupChats creates a scope for all group and supergroup chats.
func NewBotCommandScopeAllGroupChats() BotCommandScopeAllGroupChats {
	return BotCommandScopeAllGroupChats{Type: botCommandScopeAllGroupChatsType}
}

// NewBotCommandScopeAllChatAdministrators creates a scope for all group and supergroup administrators.
func NewBotCommandScopeAllChatAdministrators() BotCommandScopeAllChatAdministrators {
	return BotCommandScopeAllChatAdministrators{Type: botCommandScopeAllChatAdministratorsType}
}

// NewBotCommandScopeChat creates a scope for a specific chat.
func NewBotCommandScopeChat(chatID any) BotCommandScopeChat {
	return BotCommandScopeChat{Type: botCommandScopeChatType, ChatID: chatID}
}

// NewBotCommandScopeChatAdministrators creates a scope for administrators of a specific chat.
func NewBotCommandScopeChatAdministrators(chatID any) BotCommandScopeChatAdministrators {
	return BotCommandScopeChatAdministrators{Type: botCommandScopeChatAdministratorsType, ChatID: chatID}
}

// NewBotCommandScopeChatMember creates a scope for a specific chat member.
func NewBotCommandScopeChatMember(chatID any, userID int64) BotCommandScopeChatMember {
	return BotCommandScopeChatMember{Type: botCommandScopeChatMemberType, ChatID: chatID, UserID: userID}
}

// MarshalJSON encodes BotCommandScopeDefault with the required Telegram type field.
func (s BotCommandScopeDefault) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeDefaultType
	type scope BotCommandScopeDefault
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeAllPrivateChats with the required Telegram type field.
func (s BotCommandScopeAllPrivateChats) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeAllPrivateChatsType
	type scope BotCommandScopeAllPrivateChats
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeAllGroupChats with the required Telegram type field.
func (s BotCommandScopeAllGroupChats) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeAllGroupChatsType
	type scope BotCommandScopeAllGroupChats
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeAllChatAdministrators with the required Telegram type field.
func (s BotCommandScopeAllChatAdministrators) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeAllChatAdministratorsType
	type scope BotCommandScopeAllChatAdministrators
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeChat with the required Telegram type field.
func (s BotCommandScopeChat) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeChatType
	type scope BotCommandScopeChat
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeChatAdministrators with the required Telegram type field.
func (s BotCommandScopeChatAdministrators) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeChatAdministratorsType
	type scope BotCommandScopeChatAdministrators
	return json.Marshal(scope(s))
}

// MarshalJSON encodes BotCommandScopeChatMember with the required Telegram type field.
func (s BotCommandScopeChatMember) MarshalJSON() ([]byte, error) {
	s.Type = botCommandScopeChatMemberType
	type scope BotCommandScopeChatMember
	return json.Marshal(scope(s))
}

// ValidateBotCommandScope checks whether scope can be encoded for Telegram.
func ValidateBotCommandScope(scope BotCommandScope) error {
	if scope == nil {
		return nil
	}
	if isNilInterfaceValue(scope) {
		return stderrors.New("bot command scope must not be nil")
	}

	data, err := json.Marshal(scope)
	if err != nil {
		return fmt.Errorf("bot command scope is invalid: %w", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("bot command scope is invalid: %w", err)
	}

	scopeType, err := stringField(payload, "type")
	if err != nil {
		return err
	}

	switch scopeType {
	case botCommandScopeDefaultType,
		botCommandScopeAllPrivateChatsType,
		botCommandScopeAllGroupChatsType,
		botCommandScopeAllChatAdministratorsType:
		return nil
	case botCommandScopeChatType, botCommandScopeChatAdministratorsType:
		return validateScopeChatID(payload)
	case botCommandScopeChatMemberType:
		if err := validateScopeChatID(payload); err != nil {
			return err
		}
		userID, err := int64Field(payload, "user_id")
		if err != nil {
			return err
		}
		if userID <= 0 {
			return stderrors.New("bot command scope user_id must be greater than zero")
		}
		return nil
	default:
		return stderrors.New("bot command scope type is unsupported")
	}
}

// MenuButton marks Telegram menu button objects.
type MenuButton interface {
	menuButton()
}

// MenuButtonCommands opens the bot's command list.
type MenuButtonCommands struct {
	Type string `json:"type"`
}

// WebAppInfo describes a Telegram Web App URL.
type WebAppInfo struct {
	URL string `json:"url"`
}

// MenuButtonWebApp launches a Telegram Web App.
type MenuButtonWebApp struct {
	Type   string     `json:"type"`
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}

// MenuButtonDefault restores the default menu button behavior.
type MenuButtonDefault struct {
	Type string `json:"type"`
}

func (MenuButtonCommands) menuButton() {}
func (MenuButtonWebApp) menuButton()   {}
func (MenuButtonDefault) menuButton()  {}

// NewMenuButtonCommands creates a menu button that opens the bot's commands.
func NewMenuButtonCommands() MenuButtonCommands {
	return MenuButtonCommands{Type: menuButtonCommandsType}
}

// NewMenuButtonDefault creates a menu button that restores Telegram's default behavior.
func NewMenuButtonDefault() MenuButtonDefault {
	return MenuButtonDefault{Type: menuButtonDefaultType}
}

// NewMenuButtonWebApp creates a menu button that opens a Web App.
func NewMenuButtonWebApp(text string, url string) MenuButtonWebApp {
	return MenuButtonWebApp{Type: menuButtonWebAppType, Text: text, WebApp: WebAppInfo{URL: url}}
}

// MarshalJSON encodes MenuButtonCommands with the required Telegram type field.
func (m MenuButtonCommands) MarshalJSON() ([]byte, error) {
	m.Type = menuButtonCommandsType
	type button MenuButtonCommands
	return json.Marshal(button(m))
}

// MarshalJSON encodes MenuButtonWebApp with the required Telegram type field.
func (m MenuButtonWebApp) MarshalJSON() ([]byte, error) {
	m.Type = menuButtonWebAppType
	type button MenuButtonWebApp
	return json.Marshal(button(m))
}

// MarshalJSON encodes MenuButtonDefault with the required Telegram type field.
func (m MenuButtonDefault) MarshalJSON() ([]byte, error) {
	m.Type = menuButtonDefaultType
	type button MenuButtonDefault
	return json.Marshal(button(m))
}

// ValidateMenuButton checks whether button can be sent to Telegram.
func ValidateMenuButton(button MenuButton) error {
	if button == nil {
		return nil
	}
	if isNilInterfaceValue(button) {
		return stderrors.New("menu_button must not be nil")
	}

	switch value := button.(type) {
	case MenuButtonCommands, *MenuButtonCommands, MenuButtonDefault, *MenuButtonDefault:
		return nil
	case MenuButtonWebApp:
		return validateMenuButtonWebApp(value)
	case *MenuButtonWebApp:
		return validateMenuButtonWebApp(*value)
	default:
		return stderrors.New("unsupported menu_button")
	}
}

// UnmarshalMenuButton decodes a polymorphic Telegram MenuButton object.
func UnmarshalMenuButton(data []byte) (MenuButton, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case menuButtonCommandsType:
		var button MenuButtonCommands
		if err := json.Unmarshal(data, &button); err != nil {
			return nil, err
		}
		button.Type = menuButtonCommandsType
		return button, nil
	case menuButtonDefaultType:
		var button MenuButtonDefault
		if err := json.Unmarshal(data, &button); err != nil {
			return nil, err
		}
		button.Type = menuButtonDefaultType
		return button, nil
	case menuButtonWebAppType:
		var button MenuButtonWebApp
		if err := json.Unmarshal(data, &button); err != nil {
			return nil, err
		}
		button.Type = menuButtonWebAppType
		return button, nil
	default:
		return nil, stderrors.New("unsupported menu button type")
	}
}

func validateMenuButtonWebApp(button MenuButtonWebApp) error {
	if strings.TrimSpace(button.Text) == "" {
		return stderrors.New("menu_button.web_app text is required")
	}
	if err := validateWebAppInfo(button.WebApp, "menu_button.web_app"); err != nil {
		return err
	}
	return nil
}

// ChatAdministratorRights describes administrator rights requested by a bot.
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous,omitempty"`
	CanManageChat       bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool `json:"can_change_info,omitempty"`
	CanInviteUsers      bool `json:"can_invite_users,omitempty"`
	CanPostStories      bool `json:"can_post_stories,omitempty"`
	CanEditStories      bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories    bool `json:"can_delete_stories,omitempty"`
	CanPostMessages     bool `json:"can_post_messages,omitempty"`
	CanEditMessages     bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanManageTopics     bool `json:"can_manage_topics,omitempty"`
}

func validateScopeChatID(payload map[string]json.RawMessage) error {
	raw, ok := payload["chat_id"]
	if !ok {
		return stderrors.New("bot command scope chat_id is required")
	}
	if bytes.Equal(raw, []byte("null")) {
		return stderrors.New("bot command scope chat_id is required")
	}

	var chatID any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&chatID); err != nil {
		return err
	}
	switch value := chatID.(type) {
	case string:
		if strings.TrimSpace(value) == "" {
			return stderrors.New("bot command scope chat_id is required")
		}
	case json.Number:
		if value.String() == "0" {
			return stderrors.New("bot command scope chat_id is required")
		}
	default:
		return stderrors.New("bot command scope chat_id must be an integer or string")
	}

	return nil
}

func stringField(payload map[string]json.RawMessage, name string) (string, error) {
	raw, ok := payload[name]
	if !ok {
		return "", fmt.Errorf("bot command scope %s is required", name)
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", err
	}
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("bot command scope %s is required", name)
	}
	return value, nil
}

func int64Field(payload map[string]json.RawMessage, name string) (int64, error) {
	raw, ok := payload[name]
	if !ok {
		return 0, fmt.Errorf("bot command scope %s is required", name)
	}
	var value int64
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func isNilInterfaceValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
