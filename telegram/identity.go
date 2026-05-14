package telegram

// Actor describes the user or chat that caused a Telegram event.
//
// Actor is a convenience helper for application code. It does not decide
// whether the actor is allowed to perform an action.
type Actor struct {
	// User is set when Telegram exposes a user actor.
	User *User
	// Chat is set when Telegram exposes a chat actor, such as message sender_chat.
	Chat *Chat
	// AnonymousAdmin is true for messages sent by an anonymous group administrator.
	AnonymousAdmin bool
}

// IsZero reports whether no actor identity is available.
func (a Actor) IsZero() bool {
	return a.User == nil && a.Chat == nil
}

// Actor returns the user or chat actor for m, if Telegram exposed one.
func (m *Message) Actor() Actor {
	if m == nil {
		return Actor{}
	}

	actor := Actor{
		User:           m.From,
		Chat:           m.SenderChat,
		AnonymousAdmin: m.IsAnonymousAdmin(),
	}
	if actor.User == nil && m.GuestBotCallerUser != nil {
		actor.User = m.GuestBotCallerUser
	}
	if actor.Chat == nil && m.GuestBotCallerChat != nil {
		actor.Chat = m.GuestBotCallerChat
	}
	return actor
}

// Actor returns the user that pressed the inline keyboard button.
func (q *CallbackQuery) Actor() Actor {
	if q == nil {
		return Actor{}
	}
	return Actor{User: &q.From}
}

// Actor returns the user or chat actor for u, if Telegram exposed one.
func (u *Update) Actor() Actor {
	if u == nil {
		return Actor{}
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Actor()
	}
	if message := u.EffectiveMessage(); message != nil {
		return message.Actor()
	}
	if u.MessageReaction != nil {
		if u.MessageReaction.User != nil {
			return Actor{User: u.MessageReaction.User}
		}
		if u.MessageReaction.ActorChat != nil {
			return Actor{Chat: u.MessageReaction.ActorChat}
		}
	}
	if u.PollAnswer != nil {
		if u.PollAnswer.User != nil {
			return Actor{User: u.PollAnswer.User}
		}
		if u.PollAnswer.VoterChat != nil {
			return Actor{Chat: u.PollAnswer.VoterChat}
		}
	}
	if user := u.EffectiveUser(); user != nil {
		return Actor{User: user}
	}
	return Actor{}
}

// IsAnonymousAdmin reports whether m was sent by an anonymous group administrator.
func (m *Message) IsAnonymousAdmin() bool {
	if m == nil || m.SenderChat == nil || m.Chat.ID == 0 {
		return false
	}
	if m.SenderChat.ID != m.Chat.ID {
		return false
	}
	return m.Chat.Type == "group" || m.Chat.Type == "supergroup"
}

// ReplyTarget returns the actor of the message that m replies to, if available.
func (m *Message) ReplyTarget() Actor {
	if m == nil || m.ReplyToMessage == nil {
		return Actor{}
	}
	return m.ReplyToMessage.Actor()
}

// ReplyTargetUser returns the user actor of the message that m replies to, if available.
func (m *Message) ReplyTargetUser() *User {
	return m.ReplyTarget().User
}

// ReplyTarget returns the actor of the message that u replies to, if available.
func (u *Update) ReplyTarget() Actor {
	if u == nil {
		return Actor{}
	}
	message := u.EffectiveMessage()
	if message == nil {
		return Actor{}
	}
	return message.ReplyTarget()
}

// ReplyTargetUser returns the user actor of the message that u replies to, if available.
func (u *Update) ReplyTargetUser() *User {
	return u.ReplyTarget().User
}
