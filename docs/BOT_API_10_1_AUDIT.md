# Bot API 10.1 Audit

Date: 2026-06-12

Telegram published Bot API 10.1 on June 11, 2026. This audit tracks the update from the Bot API 10.0-complete implementation to Bot API 10.1-complete coverage.

Source: <https://core.telegram.org/bots/api#recent-changes>

## Upstream Delta

### Rich Messages

- New rich text classes:
  - `RichTextBold`
  - `RichTextItalic`
  - `RichTextUnderline`
  - `RichTextStrikethrough`
  - `RichTextSpoiler`
  - `RichTextDateTime`
  - `RichTextTextMention`
  - `RichTextSubscript`
  - `RichTextSuperscript`
  - `RichTextMarked`
  - `RichTextCode`
  - `RichTextCustomEmoji`
  - `RichTextMathematicalExpression`
  - `RichTextUrl`
  - `RichTextEmailAddress`
  - `RichTextPhoneNumber`
  - `RichTextBankCardNumber`
  - `RichTextMention`
  - `RichTextHashtag`
  - `RichTextCashtag`
  - `RichTextBotCommand`
  - `RichTextAnchor`
  - `RichTextAnchorLink`
  - `RichTextReference`
  - `RichTextReferenceLink`
- New rich block helper classes:
  - `RichBlockCaption`
  - `RichBlockTableCell`
  - `RichBlockListItem`
- New rich block classes:
  - `RichBlockParagraph`
  - `RichBlockSectionHeading`
  - `RichBlockPreformatted`
  - `RichBlockFooter`
  - `RichBlockDivider`
  - `RichBlockMathematicalExpression`
  - `RichBlockAnchor`
  - `RichBlockList`
  - `RichBlockBlockQuotation`
  - `RichBlockPullQuotation`
  - `RichBlockCollage`
  - `RichBlockSlideshow`
  - `RichBlockTable`
  - `RichBlockDetails`
  - `RichBlockMap`
  - `RichBlockAnimation`
  - `RichBlockAudio`
  - `RichBlockPhoto`
  - `RichBlockVideo`
  - `RichBlockVoiceNote`
  - `RichBlockThinking`
- New polymorphic classes: `RichText`, `RichBlock`, `RichMessage`.
- New incoming field: `Message.rich_message`.
- New input classes: `InputRichMessage`, `InputRichMessageContent`.
- New methods:
  - `sendRichMessage`;
  - `sendRichMessageDraft`.
- Changed method:
  - `editMessageText.rich_message`.
- `InputRichMessageContent` is allowed in inline, guest, and Web App query results.

### Join Request Queries

- New `User.supports_join_request_queries`.
- New `ChatFullInfo.guard_bot`.
- New `ChatJoinRequest.query_id`.
- New methods:
  - `answerChatJoinRequestQuery`;
  - `sendChatJoinRequestWebApp`.

### Polls

- New `Link` object.
- New `PollMedia.link`.
- New `InputMediaLink`, allowed as `InputPollOptionMedia`.

## Implemented

- `telegram.User.SupportsJoinRequestQueries`.
- `telegram.ChatFullInfo.GuardBot`.
- `telegram.ChatJoinRequest.QueryID`.
- `telegram.Link` and `telegram.PollMedia.Link`.
- `bot.InputMediaLink` and `bot.MediaLink`.
- `bot.AnswerChatJoinRequestQuery`.
- `bot.SendChatJoinRequestWebApp`.
- `telegram.RichMessage`.
- `telegram.RichText`, including plain strings, rich text arrays, and all official rich text variants.
- `telegram.RichBlock`, including all official rich block helper classes and block variants.
- `telegram.Message.RichMessage`;
- `bot.InputRichMessage`;
- `bot.InputRichMessageContent`;
- `bot.SendRichMessage`;
- `bot.SendRichMessageDraft`;
- `bot.EditMessageTextParams.RichMessage`;
- validation and request encoding tests for join request queries, poll link media, rich-message methods, rich edit payloads, and inline/guest/Web App rich input content.

## Pending

- No known Bot API 10.1 code coverage gaps remain.
- Live Rich Message, join request query, and advanced poll checks remain manual-only because they are user-visible or state-changing.

## Safety Classification

- `answerChatJoinRequestQuery` is admin/state-changing because it can approve or decline a real user, or leave the request queued.
- `sendChatJoinRequestWebApp` opens a Mini App during a join request query and is manual-only.
- `InputMediaLink` in poll options is poll-state-changing/noisy when live-smoked, so unit tests remain the default evidence.
- Rich message sends, edits, and drafts are user-visible message mutations and should be manual-only until dedicated test examples exist.

## Verification

Focused checks used during implementation:

```bash
go test ./bot ./telegram
```

Full completion check:

```bash
scripts/check.sh
```
