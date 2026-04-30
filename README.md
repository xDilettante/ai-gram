# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The latest public release is `v0.2.0`. Local development is now focused on reaching full Telegram Bot API 9.6 coverage before the next push, tag, or GitHub Release. The library provides practical incoming update types, a typed HTTP Bot API core, selected public Bot API methods, media sending by file_id, URL, or multipart upload, file download support, webhook management methods, a managed long polling runner, an inbound webhook HTTP handler, a small update dispatcher/router, helper middleware, examples, and manual smoke tooling. It does not yet implement FSM, scenes, storage, full thumbnail coverage for every media method, or full Bot API 9.6 coverage.

## Status

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: include practical incoming message fields for text entities, captions, media, contacts, locations, venues, and callback queries.
- Bot client package: present with token validation, private token storage, token-safe string output, an internal HTTP call core, JSON/multipart calls, and typed method wrappers.
- Typed Telegram API errors: scaffolded.
- Dispatcher/router: supports predicates, message/command/callback routes, middleware, fallback, and error handling.
- Middleware helpers: recover, timeout, hook-based observability, and reusable access control are available.
- Long polling transport: managed runner is available. Webhook transport: inbound HTTP handler is available.
- Telegram Bot API method coverage: `GetMe`, `SendMessage`, `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`, `SendContact`, `SendLocation`, `SendVenue`, `SendPoll`, `StopPoll`, `SendDice`, `SendSticker`, `SendAnimation`, `SendVideoNote`, `SendMediaGroup`, `GetStickerSet`, `GetCustomEmojiStickers`, `UploadStickerFile`, `CreateNewStickerSet`, `AddStickerToSet`, `SetStickerPositionInSet`, `DeleteStickerFromSet`, `ReplaceStickerInSet`, `SetStickerEmojiList`, `SetStickerKeywords`, `SetStickerMaskPosition`, `SetStickerSetTitle`, `SetStickerSetThumbnail`, `SetCustomEmojiStickerSetThumbnail`, `DeleteStickerSet`, `SetMyCommands`, `DeleteMyCommands`, `GetMyCommands`, `SetChatMenuButton`, `GetChatMenuButton`, `SetMyDefaultAdministratorRights`, `SetMyName`, `GetMyName`, `SetMyDescription`, `GetMyDescription`, `SetMyShortDescription`, `GetMyShortDescription`, `GetMyDefaultAdministratorRights`, `SetMyProfilePhoto`, `RemoveMyProfilePhoto`, `AnswerCallbackQuery`, `EditMessageText`, `EditMessageCaption`, `EditMessageReplyMarkup`, `DeleteMessage`, `DeleteMessages`, `ForwardMessage`, `ForwardMessages`, `CopyMessage`, `CopyMessages`, `SendChatAction`, `PinChatMessage`, `UnpinChatMessage`, `UnpinAllChatMessages`, `GetChat`, `GetChatMember`, `GetChatAdministrators`, `GetChatMemberCount`, `SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, `LeaveChat`, `SetChatStickerSet`, `DeleteChatStickerSet`, `CreateForumTopic`, `EditForumTopic`, `CloseForumTopic`, `ReopenForumTopic`, `DeleteForumTopic`, `UnpinAllForumTopicMessages`, `EditGeneralForumTopic`, `CloseGeneralForumTopic`, `ReopenGeneralForumTopic`, `HideGeneralForumTopic`, `UnhideGeneralForumTopic`, `UnpinAllGeneralForumTopicMessages`, `SetMessageReaction`, `ExportChatInviteLink`, `CreateChatInviteLink`, `EditChatInviteLink`, `RevokeChatInviteLink`, `ApproveChatJoinRequest`, `DeclineChatJoinRequest`, `PromoteChatMember`, `SetChatAdministratorCustomTitle`, `SetChatPermissions`, `BanChatMember`, `UnbanChatMember`, `RestrictChatMember`, reply markup for supported send and edit methods, the manual `GetUpdates` API call, `GetFile`, `DownloadFile`, multipart upload for media send methods, and JSON-only webhook management methods (`SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`) are implemented. The rest of the Bot API is not implemented yet.
- Public API stability: not guaranteed before the first stable release.
- Latest public release: `v0.2.0`.
- Current local strategy: full Telegram Bot API 9.6 coverage before the next push/tag/release.

## Planned architecture

The library is split into small packages with clear responsibilities:

- `telegram` contains basic Telegram Bot API data contracts such as `Update`, `Message`, `User`, `Chat`, and `CallbackQuery`.
- `bot` contains the primary Bot API client type and configuration. It owns token handling, keeps the raw token out of public accessors, and provides typed Telegram method wrappers. Use `GetMe` for bot identity and redacted string output for diagnostics.
- `internal/httpclient` contains low-level HTTP sending helpers, response body handling, and HTTP status checks. It is internal and must not leak into the public API.
- `errors` contains typed errors returned by Telegram Bot API responses.
- `dispatch` defines update routing, middleware, fallback handling, and error handling without depending on HTTP details.
- `middleware` provides reusable dispatch middleware helpers for panic recovery, per-update timeout contexts, hook-based observability, and admin/public/off access control.
- `transport/longpoll` provides a managed runner that repeatedly calls `GetUpdates` and passes updates to a handler.
- `transport/webhook` provides an inbound `net/http` handler for Telegram webhook updates.
- `aigram` is a lightweight root facade that re-exports the most important public types.

The intended dependency direction is data types first, then the Bot API client and transports, then dispatching and middleware. Transports deliver updates; dispatchers process already received updates; the API client does not know about dispatching.

## Installation

For public GitHub module usage, install from the canonical module path:

```bash
go get github.com/xDilettante/ai-gram@latest
```

`v0.1.0` was the first milestone tag. Use `v0.2.0` for the latest public release:

```bash
go get github.com/xDilettante/ai-gram@v0.2.0
```

## Usage examples

Create a bot and call `getMe`:

```go
ctx := context.Background()

b, err := aigram.New(aigram.BotConfig{Token: token})
if err != nil {
    return err
}

me, err := b.GetMe(ctx)
if err != nil {
    return err
}
fmt.Println(me.Username)
```

Send a text message:

```go
message, err := b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID: aigram.ChatIDInt(123456789),
    Text:   "Hello from ai-gram",
})
if err != nil {
    return err
}
fmt.Println(message.MessageID)
```

Set bot commands:

```go
ok, err := b.SetMyCommands(ctx, aigram.SetMyCommandsParams{
    Commands: []telegram.BotCommand{
        {Command: "start", Description: "Start the bot"},
        {Command: "help", Description: "Show help"},
    },
})
if err != nil {
    return err
}
fmt.Println(ok)
```

Read and update bot profile metadata:

```go
name, err := b.GetMyName(ctx, aigram.GetMyNameParams{})
if err != nil {
    return err
}
fmt.Println(name.Name)

ok, err = b.SetMyDescription(ctx, aigram.SetMyDescriptionParams{
    Description:  "Support bot for ai-gram examples",
    LanguageCode: "en",
})
if err != nil {
    return err
}
fmt.Println(ok)

rights, err := b.GetMyDefaultAdministratorRights(ctx, aigram.GetMyDefaultAdministratorRightsParams{
    ForChannels: true,
})
if err != nil {
    return err
}
fmt.Println(rights.CanManageChat)
```

Profile and metadata set methods change real bot state. Run live checks only with a dedicated test bot and explicit confirmation.

Send a contact or location:

```go
contactMessage, err := b.SendContact(ctx, aigram.SendContactParams{
    ChatID:      aigram.ChatIDInt(123456789),
    PhoneNumber: "+15551234567",
    FirstName:   "Ada",
})
if err != nil {
    return err
}
fmt.Println(contactMessage.MessageID)

locationMessage, err := b.SendLocation(ctx, aigram.SendLocationParams{
    ChatID:    aigram.ChatIDInt(123456789),
    Latitude:  51.5074,
    Longitude: -0.1278,
})
if err != nil {
    return err
}
fmt.Println(locationMessage.MessageID)
```

Send a poll or dice message:

```go
pollMessage, err := b.SendPoll(ctx, aigram.SendPollParams{
    ChatID:   aigram.ChatIDInt(123456789),
    Question: "Pick one",
    Options:  []string{"A", "B"},
})
if err != nil {
    return err
}
fmt.Println(pollMessage.MessageID)

diceMessage, err := b.SendDice(ctx, aigram.SendDiceParams{
    ChatID: aigram.ChatIDInt(123456789),
    Emoji:  "🎲",
})
if err != nil {
    return err
}
fmt.Println(diceMessage.MessageID)
```

Reply to an incoming message and keep a forum topic/thread when Telegram provides one:

```go
msg := update.EffectiveMessage()
if msg == nil {
    return nil
}

_, err := b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:          aigram.ChatIDInt(msg.Chat.ID),
    MessageThreadID: msg.MessageThreadID,
    Text:            "Reply from ai-gram",
    ReplyParameters: &aigram.ReplyParameters{MessageID: msg.MessageID},
})
if err != nil {
    return err
}
```

Attach reply markup:

```go
inlineKeyboard := aigram.NewInlineKeyboard(
    []aigram.InlineKeyboardButton{
        aigram.InlineButtonCallback("Confirm", "confirm"),
        aigram.InlineButtonURL("Docs", "https://core.telegram.org/bots/api"),
    },
)

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Choose an action",
    ReplyMarkup: inlineKeyboard,
})
if err != nil {
    return err
}

d := dispatch.New()
if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    ok, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Done",
    })
    if err != nil {
        return err
    }
    fmt.Println("callback answered:", ok)
    return nil
}); err != nil {
    return err
}

if err := d.OnCallbackDataFunc("danger", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Confirmation required",
        ShowAlert:       true,
    })
    return err
}); err != nil {
    return err
}
```

Use a regular reply keyboard and remove it later:

```go
replyKeyboard := aigram.NewReplyKeyboard(
    []aigram.KeyboardButton{
        aigram.KeyboardButtonText("Help"),
        aigram.KeyboardButtonContact("Share phone"),
    },
)
replyKeyboard.ResizeKeyboard = true

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Pick an option",
    ReplyMarkup: replyKeyboard,
})
if err != nil {
    return err
}

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Keyboard removed",
    ReplyMarkup: aigram.RemoveKeyboard(false),
})
if err != nil {
    return err
}
```

Edit a message from a callback query and remove the inline keyboard:

```go
d := dispatch.New()

if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Done",
    }); err != nil {
        return err
    }

    msg := update.CallbackQuery.Message
    if msg == nil {
        return nil
    }

    if _, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
        Target: aigram.EditTargetChat(aigram.ChatIDInt(msg.Chat.ID), msg.MessageID),
        Text:   "Confirmed",
    }); err != nil {
        return err
    }

    _, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
        Target: aigram.EditTargetChat(aigram.ChatIDInt(msg.Chat.ID), msg.MessageID),
        // nil ReplyMarkup removes the inline keyboard.
    })
    return err
}); err != nil {
    return err
}
```

Edit a media caption and delete a message:

```go
captionResult, err := b.EditMessageCaption(ctx, aigram.EditMessageCaptionParams{
    Target:  aigram.EditTargetChat(aigram.ChatIDInt(123456789), 42),
    Caption: "Updated caption",
})
if err != nil {
    return err
}
fmt.Println("caption edit ok:", captionResult.IsOK())

deleted, err := b.DeleteMessage(ctx, aigram.DeleteMessageParams{
    ChatID:    aigram.ChatIDInt(123456789),
    MessageID: 42,
})
if err != nil {
    return err
}
fmt.Println("deleted:", deleted)
```

Forward or copy a message:

```go
forwarded, err := b.ForwardMessage(ctx, aigram.ForwardMessageParams{
    ChatID:     aigram.ChatIDInt(123456789),
    FromChatID: aigram.ChatIDInt(987654321),
    MessageID:  42,
})
if err != nil {
    return err
}
fmt.Println("forwarded:", forwarded.MessageID)

copied, err := b.CopyMessage(ctx, aigram.CopyMessageParams{
    ChatID:     aigram.ChatIDInt(123456789),
    FromChatID: aigram.ChatIDInt(987654321),
    MessageID:  42,
    Caption:    "Copied with a new caption",
})
if err != nil {
    return err
}
fmt.Println("copied:", copied.MessageID)
```

Batch message methods support up to 100 message IDs per call. `DeleteMessages` is destructive and should be tested only on disposable test messages:

```go
forwardedIDs, err := b.ForwardMessages(ctx, aigram.ForwardMessagesParams{
    ChatID:     aigram.ChatIDInt(123456789),
    FromChatID: aigram.ChatIDInt(987654321),
    MessageIDs: []int64{42, 43},
})
if err != nil {
    return err
}
fmt.Println("forwarded batch:", len(forwardedIDs))

deletedBatch, err := b.DeleteMessages(ctx, aigram.DeleteMessagesParams{
    ChatID:     aigram.ChatIDInt(123456789),
    MessageIDs: []int64{44, 45},
})
if err != nil {
    return err
}
fmt.Println("deleted batch:", deletedBatch)
```

Send a chat action and pin or unpin a message:

```go
if _, err := b.SendChatAction(ctx, aigram.SendChatActionParams{
    ChatID: aigram.ChatIDInt(123456789),
    Action: aigram.ChatActionTyping,
}); err != nil {
    return err
}

pinned, err := b.PinChatMessage(ctx, aigram.PinChatMessageParams{
    ChatID:              aigram.ChatIDInt(123456789),
    MessageID:           42,
    DisableNotification: true,
})
if err != nil {
    return err
}
fmt.Println("pinned:", pinned)

unpinned, err := b.UnpinChatMessage(ctx, aigram.UnpinChatMessageParams{
    ChatID:    aigram.ChatIDInt(123456789),
    MessageID: 42, // omit to unpin the most recent pinned message
})
if err != nil {
    return err
}
fmt.Println("unpinned:", unpinned)
```


Get chat and member information:

```go
chat, err := b.GetChat(ctx, aigram.GetChatParams{
    ChatID: aigram.ChatIDInt(123456789),
})
if err != nil {
    return err
}
fmt.Println("chat type:", chat.Type)

member, err := b.GetChatMember(ctx, aigram.GetChatMemberParams{
    ChatID: aigram.ChatIDInt(123456789),
    UserID: 987654321,
})
if err != nil {
    return err
}
fmt.Println("member status:", member.Status)

admins, err := b.GetChatAdministrators(ctx, aigram.GetChatAdministratorsParams{
    ChatID: aigram.ChatIDInt(123456789),
})
if err != nil {
    return err
}
fmt.Println("admins:", len(admins))
```

Moderation methods require suitable bot admin rights in a group or supergroup. Use them only with explicit operator intent and test subjects:

```go
restricted, err := b.RestrictChatMember(ctx, aigram.RestrictChatMemberParams{
    ChatID:      aigram.ChatIDInt(-1001234567890),
    UserID:      987654321,
    Permissions: aigram.ChatPermissions{}, // zero permissions restrict all supported actions.
})
if err != nil {
    return err
}
fmt.Println("restricted:", restricted)

banned, err := b.BanChatMember(ctx, aigram.BanChatMemberParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    UserID: 987654321,
})
if err != nil {
    return err
}
fmt.Println("banned:", banned)

unbanned, err := b.UnbanChatMember(ctx, aigram.UnbanChatMemberParams{
    ChatID:       aigram.ChatIDInt(-1001234567890),
    UserID:       987654321,
    OnlyIfBanned: true,
})
if err != nil {
    return err
}
fmt.Println("unbanned:", unbanned)
```

Chat management and admin management methods require suitable bot admin rights where applicable and change real chat state. Keep live checks limited to dedicated test groups and revert changes after testing:

```go
renamed, err := b.SetChatTitle(ctx, aigram.SetChatTitleParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    Title:  "ai-gram test group",
})
if err != nil {
    return err
}
fmt.Println("renamed:", renamed)

described, err := b.SetChatDescription(ctx, aigram.SetChatDescriptionParams{
    ChatID:      aigram.ChatIDInt(-1001234567890),
    Description: "temporary ai-gram test description",
})
if err != nil {
    return err
}
fmt.Println("description set:", described)

left, err := b.LeaveChat(ctx, aigram.LeaveChatParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
})
if err != nil {
    return err
}
fmt.Println("left:", left)

promoted, err := b.PromoteChatMember(ctx, aigram.PromoteChatMemberParams{
    ChatID:            aigram.ChatIDInt(-1001234567890),
    UserID:            987654321,
    CanManageChat:     true,
    CanInviteUsers:    true,
    CanDeleteMessages: true,
})
if err != nil {
    return err
}
fmt.Println("promoted:", promoted)

permissionsSet, err := b.SetChatPermissions(ctx, aigram.SetChatPermissionsParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    Permissions: aigram.ChatPermissions{
        CanSendMessages: true,
        CanSendPhotos:   true,
    },
})
if err != nil {
    return err
}
fmt.Println("permissions set:", permissionsSet)
```

Forum topic methods require bot admin rights in forum supergroups and change real forum topic state. Keep live checks limited to dedicated test forum supergroups, create only test topics, and restore the General topic state after testing:

```go
topic, err := b.CreateForumTopic(ctx, aigram.CreateForumTopicParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    Name:   "ai-gram test topic",
})
if err != nil {
    return err
}
fmt.Println("created topic:", topic.MessageThreadID)

closed, err := b.CloseForumTopic(ctx, aigram.CloseForumTopicParams{
    ChatID:          aigram.ChatIDInt(-1001234567890),
    MessageThreadID: topic.MessageThreadID,
})
if err != nil {
    return err
}
fmt.Println("closed topic:", closed)
```

Reaction methods change real message reaction state. Keep live checks limited to dedicated test chats and messages, and run them only after explicit confirmation:

```go
reacted, err := b.SetMessageReaction(ctx, aigram.SetMessageReactionParams{
    ChatID:    aigram.ChatIDInt(-1001234567890),
    MessageID: 123,
    Reaction: []telegram.ReactionType{
        telegram.NewReactionTypeEmoji("👍"),
    },
    IsBig: true,
})
if err != nil {
    return err
}
fmt.Println("reaction set:", reacted)
```

Invite link methods require suitable bot admin rights in the chat. They create, edit, export, or revoke real invite links, so keep live checks limited to dedicated test groups:

```go
inviteLink, err := b.CreateChatInviteLink(ctx, aigram.CreateChatInviteLinkParams{
    ChatID:      aigram.ChatIDInt(-1001234567890),
    Name:        "smoke-test",
    MemberLimit: 5,
})
if err != nil {
    return err
}
fmt.Println("created invite link")

revoked, err := b.RevokeChatInviteLink(ctx, aigram.RevokeChatInviteLinkParams{
    ChatID:     aigram.ChatIDInt(-1001234567890),
    InviteLink: inviteLink.InviteLink,
})
if err != nil {
    return err
}
fmt.Println("revoked:", revoked.IsRevoked)
```

Chat join request methods require bot admin rights with invite-user permission. They process real pending join requests, usually created by invite links with `CreatesJoinRequest: true`, so use them only in dedicated test groups:

```go
approved, err := b.ApproveChatJoinRequest(ctx, aigram.ApproveChatJoinRequestParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    UserID: 987654321,
})
if err != nil {
    return err
}
fmt.Println("approved:", approved)

declined, err := b.DeclineChatJoinRequest(ctx, aigram.DeclineChatJoinRequestParams{
    ChatID: aigram.ChatIDInt(-1001234567890),
    UserID: 987654322,
})
if err != nil {
    return err
}
fmt.Println("declined:", declined)
```

Inline mode basics support `inline_query`/`chosen_inline_result` updates and `AnswerInlineQuery` with article, location, venue, contact, and game results. Text, location, venue, contact, and invoice input message content variants are available. Enable inline mode in BotFather before live testing; inline live smoke is manual-only:

```go
if err := d.OnInlineQueryFunc(func(ctx context.Context, update telegram.Update) error {
    query := update.InlineQuery
    if query == nil {
        return nil
    }

    _, err := b.AnswerInlineQuery(ctx, aigram.AnswerInlineQueryParams{
        InlineQueryID: query.ID,
        Results: []aigram.InlineQueryResult{
            aigram.InlineArticle(
                "echo",
                "Echo",
                aigram.InputText("Inline response"),
            ),
        },
        CacheTime:  0,
        IsPersonal: true,
    })
    return err
}); err != nil {
    return err
}
```

Reply markup currently supports inline keyboards, reply keyboards, keyboard removal, and force reply for send methods. Edit methods intentionally accept only inline keyboard markup. `AnswerCallbackQuery` can acknowledge callback taps with a toast or alert. `editMessageMedia`, WebApp/LoginUrl buttons, payments, and a keyboard builder DSL will be added separately later.

Protect handlers with access control middleware:

```go
d := dispatch.New()
d.Use(middleware.Access(middleware.AccessConfig{
    Mode:         middleware.AccessModeAdmin,
    AdminUserIDs: []int64{123456789},
}))
```

The examples default to admin-only mode through `AIGRAM_ACCESS_MODE=admin`. Use `AIGRAM_ADMIN_USER_IDS` for admins, `AIGRAM_ALLOWED_USER_IDS` / `AIGRAM_ALLOWED_CHAT_IDS` for temporary allow lists, or the runtime `/access_open` and `/access_close` commands in the examples.

Webhook examples also support Telegram deep-link smoke panels: open `https://t.me/<bot_username>?start=smoke` for the main smoke keyboard or `https://t.me/<bot_username>?start=access_panel` for the access-control panel.

Send media by `file_id`, URL, or multipart upload:

```go
photoMessage, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID:  aigram.ChatIDInt(123456789),
    Photo:   aigram.FileID("existing-photo-file-id"),
    Caption: "Photo from file_id",
})
if err != nil {
    return err
}
fmt.Println(photoMessage.MessageID)

photoByURL, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Photo:  aigram.FileURL("https://example.com/photo.jpg"),
})
if err != nil {
    return err
}
fmt.Println(photoByURL.MessageID)

documentMessage, err := b.SendDocument(ctx, aigram.SendDocumentParams{
    ChatID:   aigram.ChatIDInt(123456789),
    Document: aigram.FileID("existing-document-file-id"),
    Caption:  "Document from file_id",
})
if err != nil {
    return err
}
fmt.Println(documentMessage.MessageID)
```

Upload a photo from `os.File`:

```go
photoFile, err := os.Open("photo.jpg")
if err != nil {
    return err
}
defer photoFile.Close()

uploadedPhoto, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Photo: aigram.FileUpload(aigram.UploadFile{
        Name:        "photo.jpg",
        Reader:      photoFile,
        ContentType: "image/jpeg",
    }),
    Caption: "Uploaded photo",
})
if err != nil {
    return err
}
fmt.Println(uploadedPhoto.MessageID)
```

Upload a document from `bytes.Reader`:

```go
report := []byte("report contents")

uploadedDocument, err := b.SendDocument(ctx, aigram.SendDocumentParams{
    ChatID: aigram.ChatIDInt(123456789),
    Document: aigram.FileUpload(aigram.UploadFile{
        Name:        "report.txt",
        Reader:      bytes.NewReader(report),
        ContentType: "text/plain",
    }),
})
if err != nil {
    return err
}
fmt.Println(uploadedDocument.MessageID)
```

Send video, audio, and voice messages:

```go
videoMessage, err := b.SendVideo(ctx, aigram.SendVideoParams{
    ChatID:            aigram.ChatIDInt(123456789),
    Video:             aigram.FileID("existing-video-file-id"),
    Caption:           "Video from file_id",
    SupportsStreaming: true,
})
if err != nil {
    return err
}
fmt.Println(videoMessage.MessageID)

audioMessage, err := b.SendAudio(ctx, aigram.SendAudioParams{
    ChatID:    aigram.ChatIDInt(123456789),
    Audio:     aigram.FileURL("https://example.com/audio.mp3"),
    Performer: "Example Artist",
    Title:     "Example Track",
})
if err != nil {
    return err
}
fmt.Println(audioMessage.MessageID)

voiceMessage, err := b.SendVoice(ctx, aigram.SendVoiceParams{
    ChatID: aigram.ChatIDInt(123456789),
    Voice:  aigram.FileID("existing-voice-file-id"),
})
if err != nil {
    return err
}
fmt.Println(voiceMessage.MessageID)
```

Send sticker, animation, or video note messages:

```go
stickerMessage, err := b.SendSticker(ctx, aigram.SendStickerParams{
    ChatID:  aigram.ChatIDInt(123456789),
    Sticker: aigram.FileID("existing-sticker-file-id"),
})
if err != nil {
    return err
}
fmt.Println(stickerMessage.MessageID)

animationMessage, err := b.SendAnimation(ctx, aigram.SendAnimationParams{
    ChatID:    aigram.ChatIDInt(123456789),
    Animation: aigram.FileURL("https://example.com/animation.gif"),
    Caption:   "Animation from URL",
})
if err != nil {
    return err
}
fmt.Println(animationMessage.MessageID)

videoNoteMessage, err := b.SendVideoNote(ctx, aigram.SendVideoNoteParams{
    ChatID:    aigram.ChatIDInt(123456789),
    VideoNote: aigram.FileID("existing-video-note-file-id"),
})
if err != nil {
    return err
}
fmt.Println(videoNoteMessage.MessageID)
```

Sticker set management methods mutate real sticker sets and are manual-only for live testing. Use a dedicated test user/bot and disposable sticker set names:

```go
set, err := b.GetStickerSet(ctx, aigram.GetStickerSetParams{Name: "animals_by_bot"})
if err != nil {
    return err
}
fmt.Println(set.Title)

created, err := b.CreateNewStickerSet(ctx, aigram.CreateNewStickerSetParams{
    UserID: 123456789,
    Name:   "animals_by_bot",
    Title:  "Animals",
    Stickers: []aigram.InputSticker{
        aigram.NewInputSticker(aigram.FileID("existing-sticker-file-id"), "static", "🐱"),
    },
})
if err != nil {
    return err
}
fmt.Println("created:", created)

removed, err := b.DeleteStickerSet(ctx, aigram.DeleteStickerSetParams{Name: "animals_by_bot"})
if err != nil {
    return err
}
fmt.Println("deleted:", removed)
```

Send a media group:

```go
album, err := b.SendMediaGroup(ctx, aigram.SendMediaGroupParams{
    ChatID: aigram.ChatIDInt(123456789),
    Media: []aigram.InputMedia{
        aigram.MediaPhoto(aigram.FileID("existing-photo-file-id")),
        aigram.MediaDocument(aigram.FileURL("https://example.com/file.pdf")),
    },
})
if err != nil {
    return err
}
fmt.Println(len(album))
```

Upload a video from `os.File`:

```go
videoFile, err := os.Open("video.mp4")
if err != nil {
    return err
}
defer videoFile.Close()

uploadedVideo, err := b.SendVideo(ctx, aigram.SendVideoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Video: aigram.FileUpload(aigram.UploadFile{
        Name:        "video.mp4",
        Reader:      videoFile,
        ContentType: "video/mp4",
    }),
    SupportsStreaming: true,
})
if err != nil {
    return err
}
fmt.Println(uploadedVideo.MessageID)
```

`FileID` and `FileURL` are sent as JSON requests. `FileUpload` uses multipart/form-data and ai-gram generates the internal `attach://` value for the file field. The library consumes `UploadFile.Reader` but does not close it; the caller owns reader lifecycle. Thumbnail upload is supported for `SendAnimation`, `SendVideoNote`, and `SendMediaGroup` input media that accept thumbnails; other media thumbnail parameters are still deferred. `SendVideoNote` accepts file IDs and uploads, not HTTP URLs.

Fetch updates manually with one `getUpdates` API call:

```go
updates, err := b.GetUpdates(ctx, aigram.GetUpdatesParams{
    Limit:   10,
    Timeout: 0,
})
if err != nil {
    return err
}
for _, update := range updates {
    if update.Message != nil {
        fmt.Println(update.Message.Text)
    }
}
```

Create a small dispatcher:

```go
d := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
    fmt.Println("handler error:", err)
}))

if err := d.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
    fmt.Println("start command")
    return nil
}); err != nil {
    return err
}

if err := d.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
    fmt.Println(update.Message.Text)
    return nil
}); err != nil {
    return err
}
```


Handle common incoming message shapes:

```go
if err := d.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
    message := update.EffectiveMessage()
    if message == nil {
        return nil
    }

    switch {
    case message.IsCommand("start"):
        fmt.Println("command args:", message.CommandArguments())
    case message.HasPhoto():
        largest := message.Photo[len(message.Photo)-1]
        fmt.Println("photo:", largest.FileID)
    case message.HasDocument():
        fmt.Println("document:", message.Document.FileName)
    default:
        fmt.Println("text:", message.Text)
    }

    return nil
}); err != nil {
    return err
}

if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    fmt.Println("callback data:", update.CallbackQuery.Data)
    return nil
}); err != nil {
    return err
}
```

Telegram types currently support decoding incoming media and helper methods for handling them. Sending is currently available for text, photo, document, video, audio, voice, contact, location, venue, poll, dice, sticker, animation, video note, and media group methods; remaining specialized Bot API areas will be added separately later.

Add helper middleware:

```go
type Observer struct{}

func (Observer) OnUpdateStart(ctx context.Context, update telegram.Update) {}
func (Observer) OnUpdateFinish(ctx context.Context, update telegram.Update, err error, duration time.Duration) {
    fmt.Println("handled update in", duration, "error:", err)
}

d.Use(
    middleware.Recover(nil),
    middleware.Timeout(5*time.Second),
    middleware.Observe(Observer{}),
)
```

Observability is hook-based only for now; the library does not include Prometheus, OpenTelemetry, or a logger.

Run managed long polling with the dispatcher:

```go
runner, err := longpoll.New(b, d, longpoll.Config{
    Timeout: 30,
})
if err != nil {
    return err
}

if err := runner.Run(ctx); err != nil {
    return err
}
```

The long polling runner fetches updates and calls a handler; `dispatch.Dispatcher` is one compatible handler implementation.


Download a file by `file_id` from an incoming document:

```go
if message.Document == nil {
    return nil
}

file, err := b.GetFile(ctx, aigram.GetFileParams{
    FileID: message.Document.FileID,
})
if err != nil {
    return err
}

var buf bytes.Buffer
if err := b.DownloadFile(ctx, file.FilePath, &buf); err != nil {
    return err
}
fmt.Println("downloaded bytes:", buf.Len())
```

For large files pass an `*os.File` or another streaming `io.Writer` instead of `bytes.Buffer`. Telegram download URLs contain the bot token; ai-gram builds them internally and does not expose them as a public API. The regular cloud Bot API has Telegram-side file download limits. Upload is currently implemented for `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`, `SendSticker`, `SendAnimation`, and `SendVideoNote`; download helpers never expose a full token-bearing download URL.

Serve inbound webhook updates with `net/http`:

```go
webhookHandler, err := webhook.New(d, webhook.Config{
    SecretToken: "your-secret-token",
})
if err != nil {
    return err
}

http.Handle("/telegram/webhook", webhookHandler)
if err := http.ListenAndServe(":8080", nil); err != nil {
    return err
}
```

Manage webhook registration through outbound Bot API methods:

```go
secret := "my_secret_123"

ok, err := b.SetWebhook(ctx, aigram.SetWebhookParams{
    URL:         "https://example.com/telegram/webhook",
    SecretToken: secret,
})
if err != nil {
    return err
}
if !ok {
    return fmt.Errorf("set webhook returned false")
}

webhookHandler, err := webhook.New(d, webhook.Config{
    SecretToken: secret,
})
if err != nil {
    return err
}
http.Handle("/telegram/webhook", webhookHandler)
```

The `SecretToken` passed to `SetWebhook` must match `transport/webhook.Config.SecretToken` so inbound requests can be verified.

Read current webhook status:

```go
info, err := b.GetWebhookInfo(ctx)
if err != nil {
    return err
}
fmt.Println(info.URL, info.PendingUpdateCount)
```

Delete a webhook:

```go
ok, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{
    DropPendingUpdates: true,
})
if err != nil {
    return err
}
fmt.Println(ok)
```

Webhook management is JSON-only for now. Webhook certificate upload, full thumbnail coverage, editMessageMedia, remaining inline result variants, WebApp/LoginUrl buttons, payments, FSM, scenes, storage, dependency injection, and full Bot API coverage are not implemented yet.


## Examples

Runnable examples are available under `examples/`:

- `examples/echo_longpoll` — basic long polling echo bot.
- `examples/inline_longpoll` — admin-only inline keyboard callbacks with `AnswerCallbackQuery` and runtime access commands.
- `examples/webhook_server` — admin-only inbound webhook server with `SetWebhook`, safe action logs, callback edit/delete flows, and caption edit smoke.
- `examples/media_upload` — document upload and file download smoke checks.
- `examples/v02_send_methods` — safe v0.2 send-method live smoke helper for contact, location, venue, poll, dice, and optional media sends.
- `examples/media_group_smoke` — targeted `SendMediaGroup` live smoke helper with generated upload fallback.
- `examples/local_api_server` — connectivity check for a local Telegram Bot API server.

## Manual testing

Manual smoke testing instructions are in [`docs/MANUAL_TESTING.md`](docs/MANUAL_TESTING.md). The examples require real environment variables at runtime, but they are written so `go test ./...` can compile them without a token.

Coverage and planning documents:

- [`CHANGELOG.md`](CHANGELOG.md) — release history and v0.2/v0.1.x summaries.
- [`docs/releases/v0.2.0.md`](docs/releases/v0.2.0.md) — release notes for the v0.2 coverage milestone.
- [`docs/releases/v0.1.1.md`](docs/releases/v0.1.1.md) — GitHub-release-ready notes for the canonical module path patch release.
- [`docs/releases/v0.1.0.md`](docs/releases/v0.1.0.md) — historical first milestone notes.
- [`docs/API_COVERAGE.md`](docs/API_COVERAGE.md) — implemented methods, missing Bot API areas, risk classification, and v0.1 recommendation.
- [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](docs/BOT_API_9_6_COVERAGE_PLAN.md) — local-only full Telegram Bot API 9.6 coverage plan and freeze policy.
- [`docs/V0_2_CHECKPOINT.md`](docs/V0_2_CHECKPOINT.md) — v0.2 coverage checkpoint and release recommendation.
- [`docs/V0_3_PLAN.md`](docs/V0_3_PLAN.md) — superseded v0.3 planning notes; full Bot API 9.6 coverage is tracked in the coverage plan.
- [`docs/ROADMAP.md`](docs/ROADMAP.md) — stabilization and expansion roadmap.
- [`docs/MANUAL_TESTING.md`](docs/MANUAL_TESTING.md) — local/manual smoke checklist.
- [`docs/DEPLOY_TESTING.md`](docs/DEPLOY_TESTING.md) — deploy/manual integration harness.
- [`docs/LIVE_SMOKE_MATRIX.md`](docs/LIVE_SMOKE_MATRIX.md) — safe and dangerous live smoke flows.
- [`docs/RELEASE_CHECKLIST.md`](docs/RELEASE_CHECKLIST.md) — pre-tag v0.1 release checklist.

Deployment-oriented manual integration checks are described in [`docs/DEPLOY_TESTING.md`](docs/DEPLOY_TESTING.md). The deploy harness can start from a minimal `.env.local` with bot token, chat ID, and SSH alias, then write discovered values to ignored `.deploy/generated.env`. Smoke scripts can open a temporary SSH tunnel when a discovered local Bot API server listens only on a remote loopback; the Bot API server may live on a separate SSH target from the webhook deploy target. Remote logs are redacted before printing.

The integration harness supports role-specific test bot tokens (`MAIN`, `LOCAL`, `WEBHOOK`, `NOTIFY`, and others) while preserving the legacy single-token `AIGRAM_BOT_TOKEN` mode. Set `AIGRAM_BOTAPI_SSH_TARGET` when the local Telegram Bot API server runs on a different SSH host than the webhook example.

Manual smoke scripts can also send Telegram notifications with the target `@username`, `t.me` deep link, exact panel/button guidance, and what Codex will verify in safe logs. Webhook deploy notifications use `AIGRAM_SMOKE_MODE=targeted` by default, so deploys do not ask for a full checklist unless `AIGRAM_SMOKE_MODE=full` is explicitly set. See [`docs/DEPLOY_TESTING.md`](docs/DEPLOY_TESTING.md#telegram-notifications-during-smoke-checks).

Use `./scripts/smoke_v02_send_methods.sh` for the targeted v0.2 send-method live smoke. It sends contact/location/venue/poll/dice test messages to the configured smoke chat; sticker, animation, and video note checks run only when the matching optional media environment variables are set.

Use `./scripts/smoke_media_group.sh` for the targeted `SendMediaGroup` live smoke. By default it sends two generated small text documents as a media group, and optional file ID/path env variables can switch it to FileID or upload mode.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
