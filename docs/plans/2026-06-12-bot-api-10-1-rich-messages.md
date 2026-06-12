# Bot API 10.1 Rich Messages Implementation Plan

Goal: finish Telegram Bot API 10.1 coverage by adding Rich Message contracts, send/edit methods, inline content support, docs, tests, and release metadata.

Source: <https://core.telegram.org/bots/api#recent-changes>

## Already Done

- Bot API 10.1 audit is recorded in `docs/BOT_API_10_1_AUDIT.md`.
- First 10.1 slice is implemented and pushed in `e7a7bdd`:
  - join request query helpers;
  - `User.supports_join_request_queries`;
  - `ChatFullInfo.guard_bot`;
  - `ChatJoinRequest.query_id`;
  - `Link`, `PollMedia.link`, and `InputMediaLink`.
- Local `scripts/check.sh` and GitHub Actions CI passed for the first slice.

## Files

- Create `telegram/rich_messages.go`: received Rich Message contracts and polymorphic JSON decoding.
- Create `telegram/rich_messages_test.go`: Rich Message decode coverage.
- Modify `telegram/types.go`: add `Message.RichMessage`.
- Modify `telegram/message_metadata.go`: decode `Message.rich_message`.
- Create `bot/rich_message_methods.go`: `InputRichMessage`, `sendRichMessage`, and `sendRichMessageDraft`.
- Create `bot/rich_message_methods_test.go`: request encoding, validation, API error, and context tests.
- Modify `bot/edit_methods.go` and `bot/edit_methods_test.go`: support `EditMessageTextParams.RichMessage`.
- Modify `bot/inline_methods.go` and related tests: allow `InputRichMessageContent`.
- Update docs: README, CHANGELOG, API coverage, 10.1 audit, roadmap, smoke matrix.

## Tasks

- [x] Task 1: Add failing decode tests for `RichText`, `RichBlock`, `RichMessage`, and `Message.rich_message`.
  - Command: `go test ./telegram`
  - Expected before implementation: compile failure for missing Rich Message types.

- [x] Task 2: Implement `telegram/rich_messages.go` and `Message.rich_message` decoding.
  - Use interfaces `RichText` and `RichBlock`.
  - Decode plain rich text strings, arrays, and typed objects.
  - Decode all official rich text and rich block variants from Bot API 10.1.
  - Command: `go test ./telegram`
  - Expected after implementation: pass.

- [x] Task 3: Add failing bot tests for `sendRichMessage`, `sendRichMessageDraft`, and rich edit payloads.
  - Command: `go test ./bot`
  - Expected before implementation: compile failure for missing bot APIs.

- [x] Task 4: Implement `bot.InputRichMessage`, send/draft methods, edit integration, and validation.
  - Require exactly one of `HTML` or `Markdown`.
  - Validate chat IDs, non-negative thread IDs, non-zero draft IDs, reply parameters, reply markup, and edit target rules.
  - Allow `EditMessageTextParams` to send either text or rich message, but not both.
  - Command: `go test ./bot`
  - Expected after implementation: pass.

- [x] Task 5: Add `InputRichMessageContent` to inline/guest/Web App query flows.
  - Extend `InputMessageContent`.
  - Reuse `InputRichMessage` validation.
  - Command: `go test ./bot`
  - Expected after implementation: pass.

- [x] Task 6: Update docs and public status.
  - Mark Bot API 10.1 as complete only after Rich Messages pass tests.
  - Update `docs/BOT_API_10_1_AUDIT.md` pending/completed sections.
  - Command: `scripts/check.sh`
  - Expected: all checks passed.

- [ ] Task 7: Commit, push, and wait for CI.
  - Commit message: `Complete Bot API 10.1 rich messages`
  - Commands:
    - `git status --short`
    - `git diff --cached --check`
    - `git commit -m "Complete Bot API 10.1 rich messages"`
    - `git push origin main`
    - `gh run watch <run-id> --exit-status`

## Notes

- Do not run live Telegram checks for Rich Messages without explicit approval.
- Keep media upload support out of this slice unless the official input object requires it; `InputRichMessage` sends formatted HTML/Markdown strings, not multipart files.
- Keep `telegram` package free of `bot` package dependencies.
