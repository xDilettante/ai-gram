package telegram

import "testing"

func TestMessageActorUser(t *testing.T) {
	user := &User{ID: 7, FirstName: "Alice"}
	message := &Message{
		MessageID: 1,
		From:      user,
		Chat:      Chat{ID: -1001, Type: "supergroup", Title: "Group"},
		Text:      "hello",
	}

	actor := message.Actor()
	if actor.User != user {
		t.Fatalf("unexpected actor user: %+v", actor.User)
	}
	if actor.Chat != nil {
		t.Fatalf("unexpected actor chat: %+v", actor.Chat)
	}
	if actor.AnonymousAdmin {
		t.Fatal("regular user message should not be anonymous admin")
	}
	if actor.IsZero() {
		t.Fatal("actor should not be zero")
	}
}

func TestMessageActorAnonymousAdmin(t *testing.T) {
	chat := Chat{ID: -1001, Type: "supergroup", Title: "Group"}
	message := &Message{
		MessageID:       1,
		Chat:            chat,
		SenderChat:      &chat,
		AuthorSignature: "admin",
		Text:            "moderation note",
	}

	actor := message.Actor()
	if actor.User != nil {
		t.Fatalf("anonymous admin should not invent user actor: %+v", actor.User)
	}
	if actor.Chat == nil || actor.Chat.ID != chat.ID {
		t.Fatalf("unexpected actor chat: %+v", actor.Chat)
	}
	if !actor.AnonymousAdmin {
		t.Fatal("expected anonymous admin actor")
	}
	if !message.IsAnonymousAdmin() {
		t.Fatal("expected IsAnonymousAdmin to report true")
	}
}

func TestMessageActorSenderChatIsNotAlwaysAnonymousAdmin(t *testing.T) {
	message := &Message{
		MessageID:  1,
		Chat:       Chat{ID: -1001, Type: "supergroup", Title: "Group"},
		SenderChat: &Chat{ID: -2002, Type: "channel", Title: "Channel"},
		Text:       "channel post in group",
	}

	actor := message.Actor()
	if actor.Chat == nil || actor.Chat.ID != -2002 {
		t.Fatalf("unexpected sender chat actor: %+v", actor.Chat)
	}
	if actor.AnonymousAdmin {
		t.Fatal("different sender_chat should not be anonymous admin")
	}
	if message.IsAnonymousAdmin() {
		t.Fatal("expected IsAnonymousAdmin to report false")
	}
}

func TestUpdateActorPrefersCallbackUser(t *testing.T) {
	messageUser := &User{ID: 7, FirstName: "Message user"}
	callbackUser := User{ID: 8, FirstName: "Callback user"}
	message := &Message{
		MessageID: 1,
		From:      messageUser,
		Chat:      Chat{ID: -1001, Type: "supergroup", Title: "Group"},
	}
	update := &Update{
		CallbackQuery: &CallbackQuery{
			ID:      "callback-id",
			From:    callbackUser,
			Message: &MaybeInaccessibleMessage{Message: message, MessageID: message.MessageID, Chat: message.Chat},
		},
	}

	actor := update.Actor()
	if actor.User == nil || actor.User.ID != callbackUser.ID {
		t.Fatalf("unexpected callback actor: %+v", actor.User)
	}
}

func TestUpdateActorFromJoinRequestAndAnonymousPollAnswer(t *testing.T) {
	join := &Update{ChatJoinRequest: &ChatJoinRequest{
		Chat: Chat{ID: -1001, Type: "supergroup"},
		From: User{ID: 11, FirstName: "Joiner"},
	}}
	if actor := join.Actor(); actor.User == nil || actor.User.ID != 11 {
		t.Fatalf("unexpected join request actor: %+v", actor)
	}

	poll := &Update{PollAnswer: &PollAnswer{
		PollID:    "poll-id",
		VoterChat: &Chat{ID: -1002, Type: "supergroup", Title: "Anonymous voters"},
		OptionIDs: []int{1},
	}}
	if actor := poll.Actor(); actor.Chat == nil || actor.Chat.ID != -1002 {
		t.Fatalf("unexpected anonymous poll actor: %+v", actor)
	}
}

func TestUpdateActorFromReactionActorChat(t *testing.T) {
	update := &Update{MessageReaction: &MessageReactionUpdated{
		Chat:      Chat{ID: -1001, Type: "supergroup"},
		MessageID: 10,
		ActorChat: &Chat{
			ID:    -1002,
			Type:  "channel",
			Title: "Channel actor",
		},
	}}

	actor := update.Actor()
	if actor.Chat == nil || actor.Chat.ID != -1002 {
		t.Fatalf("unexpected reaction actor chat: %+v", actor.Chat)
	}
}

func TestReplyTargetHelpers(t *testing.T) {
	replyUser := &User{ID: 12, FirstName: "Original"}
	reply := &Message{
		MessageID: 1,
		From:      replyUser,
		Chat:      Chat{ID: -1001, Type: "supergroup"},
		Text:      "original",
	}
	message := &Message{
		MessageID:      2,
		From:           &User{ID: 13, FirstName: "Responder"},
		Chat:           Chat{ID: -1001, Type: "supergroup"},
		ReplyToMessage: reply,
		Text:           "reply",
	}

	if target := message.ReplyTarget(); target.User != replyUser {
		t.Fatalf("unexpected reply target: %+v", target)
	}
	if user := message.ReplyTargetUser(); user != replyUser {
		t.Fatalf("unexpected reply target user: %+v", user)
	}

	update := &Update{Message: message}
	if user := update.ReplyTargetUser(); user != replyUser {
		t.Fatalf("unexpected update reply target user: %+v", user)
	}
}

func TestNilIdentityHelpers(t *testing.T) {
	var message *Message
	if !message.Actor().IsZero() {
		t.Fatal("nil message actor should be zero")
	}
	if !message.ReplyTarget().IsZero() {
		t.Fatal("nil message reply target should be zero")
	}
	if message.ReplyTargetUser() != nil {
		t.Fatal("nil message reply target user should be nil")
	}

	var update *Update
	if !update.Actor().IsZero() {
		t.Fatal("nil update actor should be zero")
	}
	if !update.ReplyTarget().IsZero() {
		t.Fatal("nil update reply target should be zero")
	}
	if update.ReplyTargetUser() != nil {
		t.Fatal("nil update reply target user should be nil")
	}

	var query *CallbackQuery
	if !query.Actor().IsZero() {
		t.Fatal("nil callback query actor should be zero")
	}
}
