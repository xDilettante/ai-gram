# v0.3 Plan

## Goal

Expand `ai-gram` beyond v0.2 by covering group/chat administration and advanced interaction surfaces. v0.3 should stay incremental: each slice should keep the public API typed, testable, and documented without turning the library into a framework.

## Scope candidates

### Slice 1: Chat management

Methods:

- `SetChatTitle`
- `SetChatDescription`
- `SetChatPhoto`
- `DeleteChatPhoto`
- `LeaveChat`

Why first:

- Natural continuation of v0.2 admin/group coverage.
- Smaller than forum topics or inline mode.
- Mostly straightforward request/response methods.
- Reuses existing `ChatID`, `FileUpload`, JSON/multipart, validation, and error patterns.
- Live checks are admin/state-changing, so unit/httptest coverage should come first and live verification should stay manual-only.

Expected verification shape:

- Unit/httptest success and failure coverage for each method.
- Multipart coverage for `SetChatPhoto`.
- Token-leak regression checks.
- Manual-only live checklist for dedicated test chats.

### Slice 2: Forum topics

Methods:

- `CreateForumTopic`
- `EditForumTopic`
- `CloseForumTopic`
- `ReopenForumTopic`
- `DeleteForumTopic`
- `UnpinAllForumTopicMessages`

Why second:

- Builds on chat administration but adds topic-specific state.
- Useful for supergroups with forum mode enabled.
- Requires careful docs because live tests need a dedicated forum-enabled supergroup.

### Slice 3: Reactions

Methods:

- `SetMessageReaction`

Types:

- `ReactionType` if needed
- reaction-related update/type support if needed

Why third:

- Adds modern interaction coverage without the breadth of inline mode.
- Needs careful type design because Telegram reaction types are polymorphic.

### Slice 4: Inline mode basics

Methods:

- `AnswerInlineQuery`

Types:

- `InlineQuery`
- minimal `InlineQueryResult` variants
- chosen inline result support if needed

Why fourth:

- High-value but broader than the previous slices.
- Requires polymorphic result serialization, update decoding, dispatcher predicates, and examples.
- Should be planned after simpler v0.3 slices validate the pattern for new update surfaces.

## Recommended order

1. Chat management
2. Forum topics
3. Reactions
4. Inline mode basics

## First implementation slice

Start with **Slice 1: Chat management**. It has the best balance of value, low design risk, and reuse of existing v0.2 patterns. Recommended first stage:

1. Add params and methods for `SetChatTitle`, `SetChatDescription`, `DeleteChatPhoto`, and `LeaveChat`.
2. Add multipart `SetChatPhoto` using existing upload helpers.
3. Add unit/httptest coverage for success, validation, API errors, invalid JSON, HTTP errors, cancelled context, and token leakage.
4. Update README, API coverage, changelog, and manual testing docs.
5. Do not run live smoke automatically.

## Live smoke policy

- Safe/read-only flows may be live-smoked when they do not mutate bot/chat state.
- Admin/state-changing flows require a dedicated test chat and explicit user confirmation.
- No destructive/admin automatic live smoke.
- For chat management, live checks must use a dedicated test group/channel and must include a rollback plan for title, description, and photo changes.
- `LeaveChat` should never be part of automatic smoke because it removes the bot from the chat.

## Release policy

`v0.3.0` should be tagged only after:

- static checks pass;
- API coverage is updated;
- README and manual testing docs are updated;
- changelog entries are grouped and accurate;
- safe live smoke, if any, is green;
- manual-only risks are documented;
- no token/secret leaks are found;
- public `go get` is verified after the tag is published.
