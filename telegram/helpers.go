package telegram

import (
	"strings"
	"unicode"
)

// IsText reports whether m contains non-empty text.
func (m *Message) IsText() bool {
	return m != nil && m.Text != ""
}

// IsCommand reports whether m starts with the given bot command without a leading slash.
func (m *Message) IsCommand(command string) bool {
	return validCommandName(command) && m.Command() == command
}

// Command returns the command name without a leading slash or bot username.
func (m *Message) Command() string {
	if m == nil || !strings.HasPrefix(m.Text, "/") {
		return ""
	}

	commandToken, _ := splitCommand(m.Text)
	commandToken = strings.TrimPrefix(commandToken, "/")
	if commandToken == "" {
		return ""
	}
	if at := strings.IndexByte(commandToken, '@'); at >= 0 {
		commandToken = commandToken[:at]
	}

	return commandToken
}

// CommandArguments returns the text after the leading bot command.
func (m *Message) CommandArguments() string {
	if m == nil || !strings.HasPrefix(m.Text, "/") {
		return ""
	}

	commandToken, rest := splitCommand(m.Text)
	if commandToken == "/" || strings.TrimPrefix(commandToken, "/") == "" {
		return ""
	}

	return strings.TrimLeftFunc(rest, unicode.IsSpace)
}

// HasPhoto reports whether m contains at least one photo size.
func (m *Message) HasPhoto() bool {
	return m != nil && len(m.Photo) > 0
}

// HasDocument reports whether m contains a document.
func (m *Message) HasDocument() bool {
	return m != nil && m.Document != nil
}

// HasMedia reports whether m contains a supported media payload.
func (m *Message) HasMedia() bool {
	return m != nil && (len(m.Photo) > 0 || m.Document != nil || m.Animation != nil || m.Audio != nil || m.Video != nil || m.Voice != nil || m.Sticker != nil)
}

// EffectiveMessage returns the most relevant message contained in u.
func (u *Update) EffectiveMessage() *Message {
	if u == nil {
		return nil
	}
	if u.Message != nil {
		return u.Message
	}
	if u.EditedMessage != nil {
		return u.EditedMessage
	}
	if u.ChannelPost != nil {
		return u.ChannelPost
	}
	if u.EditedChannelPost != nil {
		return u.EditedChannelPost
	}
	if u.BusinessMessage != nil {
		return u.BusinessMessage
	}
	if u.EditedBusinessMessage != nil {
		return u.EditedBusinessMessage
	}
	if u.GuestMessage != nil {
		return u.GuestMessage
	}
	if u.CallbackQuery != nil && u.CallbackQuery.Message != nil {
		return u.CallbackQuery.Message.AccessibleMessage()
	}

	return nil
}

// IsAccessible reports whether m contains a full accessible Message.
func (m *MaybeInaccessibleMessage) IsAccessible() bool {
	return m != nil && m.Message != nil
}

// AccessibleMessage returns the full message when it is accessible to the bot.
func (m *MaybeInaccessibleMessage) AccessibleMessage() *Message {
	if m == nil {
		return nil
	}
	return m.Message
}

// EffectiveChat returns the chat most directly associated with u.
func (u *Update) EffectiveChat() *Chat {
	message := u.EffectiveMessage()
	if message != nil {
		return &message.Chat
	}
	if u != nil && u.ChatJoinRequest != nil {
		return &u.ChatJoinRequest.Chat
	}
	if u != nil && u.MessageReaction != nil {
		return &u.MessageReaction.Chat
	}
	if u != nil && u.MessageReactionCount != nil {
		return &u.MessageReactionCount.Chat
	}
	if u != nil && u.PollAnswer != nil && u.PollAnswer.VoterChat != nil {
		return u.PollAnswer.VoterChat
	}
	if u != nil && u.DeletedBusinessMessages != nil {
		return &u.DeletedBusinessMessages.Chat
	}
	if u != nil && u.MyChatMember != nil {
		return &u.MyChatMember.Chat
	}
	if u != nil && u.ChatMember != nil {
		return &u.ChatMember.Chat
	}
	if u != nil && u.ChatBoost != nil {
		return &u.ChatBoost.Chat
	}
	if u != nil && u.RemovedChatBoost != nil {
		return &u.RemovedChatBoost.Chat
	}

	return nil
}

// EffectiveUser returns the user most directly associated with u.
func (u *Update) EffectiveUser() *User {
	if u == nil {
		return nil
	}
	if u.Message != nil && u.Message.From != nil {
		return u.Message.From
	}
	if u.EditedMessage != nil && u.EditedMessage.From != nil {
		return u.EditedMessage.From
	}
	if u.BusinessConnection != nil {
		return &u.BusinessConnection.User
	}
	if u.BusinessMessage != nil && u.BusinessMessage.From != nil {
		return u.BusinessMessage.From
	}
	if u.EditedBusinessMessage != nil && u.EditedBusinessMessage.From != nil {
		return u.EditedBusinessMessage.From
	}
	if u.GuestMessage != nil {
		if u.GuestMessage.From != nil {
			return u.GuestMessage.From
		}
		if u.GuestMessage.GuestBotCallerUser != nil {
			return u.GuestMessage.GuestBotCallerUser
		}
	}
	if u.CallbackQuery != nil {
		return &u.CallbackQuery.From
	}
	if u.InlineQuery != nil {
		return &u.InlineQuery.From
	}
	if u.ChosenInlineResult != nil {
		return &u.ChosenInlineResult.From
	}
	if u.ShippingQuery != nil {
		return &u.ShippingQuery.From
	}
	if u.PreCheckoutQuery != nil {
		return &u.PreCheckoutQuery.From
	}
	if u.PurchasedPaidMedia != nil {
		return &u.PurchasedPaidMedia.From
	}
	if u.ManagedBot != nil {
		return &u.ManagedBot.User
	}
	if u.PollAnswer != nil && u.PollAnswer.User != nil {
		return u.PollAnswer.User
	}
	if u.ChatJoinRequest != nil {
		return &u.ChatJoinRequest.From
	}
	if u.MyChatMember != nil {
		return &u.MyChatMember.From
	}
	if u.ChatMember != nil {
		return &u.ChatMember.From
	}
	if u.ChatBoost != nil {
		return effectiveChatBoostUser(u.ChatBoost.Boost.Source)
	}
	if u.RemovedChatBoost != nil {
		return effectiveChatBoostUser(u.RemovedChatBoost.Source)
	}
	if u.MessageReaction != nil && u.MessageReaction.User != nil {
		return u.MessageReaction.User
	}

	return nil
}

func effectiveChatBoostUser(source ChatBoostSource) *User {
	switch value := source.(type) {
	case ChatBoostSourcePremium:
		return &value.User
	case *ChatBoostSourcePremium:
		if value != nil {
			return &value.User
		}
	case ChatBoostSourceGiftCode:
		return &value.User
	case *ChatBoostSourceGiftCode:
		if value != nil {
			return &value.User
		}
	case ChatBoostSourceGiveaway:
		return value.User
	case *ChatBoostSourceGiveaway:
		if value != nil {
			return value.User
		}
	}
	return nil
}

func validCommandName(command string) bool {
	return command != "" && !strings.HasPrefix(command, "/") && !strings.ContainsFunc(command, unicode.IsSpace)
}

func splitCommand(text string) (command string, rest string) {
	for index, r := range text {
		if unicode.IsSpace(r) {
			return text[:index], text[index:]
		}
	}

	return text, ""
}
