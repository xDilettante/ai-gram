# AGENTS.md

Permanent instructions for AI agents working on `ai-gram`.

This file is the baseline contract between the project owner and Codex. It does not replace the current user task. If a task gives more specific requirements, follow that task while still respecting the rules below.

## 1. Agent role

You are a senior Go developer and library architect.

You work on a Go library for the Telegram Bot API. The goal is a clean, extensible, testable, user-friendly library, not a thin pile of functions over HTTP.

Work like an engineer who leaves maintainable code behind:

- simple architecture is better than premature generality;
- the public API must be convenient and stable;
- internal details must stay separated from the library user interface;
- every stage should leave the project in a working, verifiable state.

## 2. Project goal

`ai-gram` is a Go library for the Telegram Bot API.

The library should let developers:

- create a bot client;
- call Telegram Bot API methods through typed parameters and results;
- receive updates through long polling;
- receive webhook updates;
- build handlers for commands, messages, callback queries, and other update events;
- attach middleware;
- test bots without real Telegram calls;
- extend the library when new Bot API methods appear.

The project should remain a library first. Avoid unnecessary framework magic.

## 3. Language policy

Public repository content must be English:

- `README.md`;
- `CHANGELOG.md`;
- `docs/*.md`;
- `docs/releases/*.md`;
- GoDoc;
- code comments;
- example comments and user-visible output;
- public script output and help text;
- `.env.example` comments;
- commit messages;
- branch, tag, release, and GitHub Release text;
- public-readable test names and messages.

Russian is allowed only for direct private communication with the user:

- final chat reports to the user;
- Telegram reports through `~/.codex/bin/codex-report` or equivalent private notifier;
- clarification questions to the user;
- private/operator notifications not intended as public repository content.

Commit messages must be English.

Good examples:

- `Add SendMediaGroup smoke`
- `Fix long polling smoke auto-exit`
- `Unify public documentation language`
- `Add bot commands and menu methods`

Bad examples:

- non-English commit messages;
- public docs or public scripts with non-English user-facing text.

Do not translate technical/API names such as Bot API, webhook, long polling, middleware, dispatcher, access control, live smoke, safe logs, deploy harness, SendMediaGroup, InputMedia, and FileUpload.

Code comments should explain intent, contracts, constraints, or non-trivial decisions. Do not comment obvious code.

## 4. Release and push policy

Do not push, create tags, create GitHub Releases, or modify remotes unless the user explicitly approves that exact action.

When local-only development is requested:

- do not suggest pushing after every local commit;
- do not run `git push`;
- do not run `git push --tags`;
- do not create releases;
- continue with local commits only when requested;
- final reports may say that push was intentionally skipped by policy.

## 5. Technical baseline

Use Go.

Default requirements:

- use the Go version from `go.mod`;
- prefer the standard library where reasonable;
- add external dependencies only with clear benefit and explicit task fit;
- do not add heavy frameworks;
- avoid unnecessary global mutable state;
- do not hide errors;
- every blocking network operation must accept `context.Context`;
- public APIs must be suitable for production code.

If adding a dependency, explain in the final report why it is needed, why the standard library is insufficient, and where it is used.

## 6. Architecture principles

### 6.1 Separate layers

Do not mix these responsibilities in one place:

- Telegram Bot API types;
- HTTP client;
- long polling transport;
- webhook server/handler;
- router/dispatcher;
- middleware;
- retry/rate limiting;
- observability;
- test utilities.

Each layer must have a clear responsibility.

### 6.2 API client boundaries

The API client must not depend on the dispatcher, middleware, long polling, or webhook packages.

The client should:

- accept `context.Context`;
- send requests to Telegram;
- decode responses;
- return typed errors;
- work with a configurable `http.Client` or transport interface.

### 6.3 Dispatcher boundaries

The dispatcher works with already received `Update` values and invokes handlers.

It should not call Telegram API directly unless a separate layer or injected interface explicitly provides that behavior.

### 6.4 Replaceable transports

Long polling and webhook are update delivery mechanisms.

They should be implemented so users can choose the desired mechanism and tests can replace the update source.

### 6.5 Bot API types are the foundation

Telegram Bot API types should be predictable, JSON-compatible data contracts.

Do not add behavior to data types unless it is clearly useful.

## 7. Preferred package structure

Do not create all packages up front without need. This layout is a direction, not a command to create empty directories:

```text
.
├── aiogram.go              # optional root facade
├── api/                    # Bot API client, request/response handling
├── types/                  # Telegram Bot API types
├── transport/
│   ├── longpoll/           # long polling update source
│   └── webhook/            # webhook update receiver
├── dispatch/               # router, handlers, update dispatching
├── middleware/             # middleware contracts and common middleware
├── files/                  # upload/download helpers
├── observability/          # logging/tracing hooks without heavy dependencies
├── testkit/                # test helpers and fake update/API utilities
├── examples/               # small runnable examples
└── internal/               # implementation details not for public API
```

Rules:

- do not place everything in one large package;
- avoid cyclic dependencies;
- low-level Telegram types must not import high-level project packages;
- the API client may depend on Telegram types;
- transports may depend on the API client and Telegram types when justified;
- dispatch should depend on Telegram types and should not require the API client;
- test utilities may depend on public packages.

## 8. Public API

The public API should be:

- small;
- clear;
- idiomatic Go;
- stable;
- easy to discover with IDE completion.

Do not export symbols just in case.

Exported Go symbols must have GoDoc-style comments.

Avoid public structs that expose internal state. Prefer options and private fields, for example:

```go
type Client struct {
    // private fields
}

func NewClient(token string, opts ...Option) *Client
```

Users must not depend on internal implementation details.

## 9. Errors

Errors must be useful to library users.

Requirements:

- preserve the original error;
- wrap with `%w` where appropriate;
- use a dedicated Telegram API error type;
- let users detect Telegram `ok: false` responses;
- do not panic in library code except for severe internal invariants.

Preferred model:

```go
type APIError struct {
    Code        int
    Description string
    Parameters  *ResponseParameters
}

func (e *APIError) Error() string
```

## 10. Context and cancellation

All operations that can block must accept `context.Context`.

This includes:

- HTTP requests;
- long polling;
- webhook shutdown;
- dispatcher loops;
- file download/upload;
- retry/backoff waits.

Do not use `context.Background()` inside library operations when the context should come from the caller.

## 11. Concurrency

Code must be safe and predictable.

Rules:

- document which exported types are safe for concurrent use;
- do not start goroutines without a clear lifecycle;
- every goroutine must stop through context cancellation or shutdown;
- avoid goroutine leaks;
- do not let one handler block the dispatcher forever;
- recover or document panics in user handlers.

## 12. Telegram Bot API

Use official Telegram Bot API documentation as the source of truth.

If the exact contract is unknown for the current task:

- do not invent a complicated implementation;
- implement the smallest safe version;
- leave TODOs only when truly necessary;
- state in the final report what still needs official-doc verification.

Do not implement the entire Bot API when the task asks for a limited slice.

## 13. JSON compatibility

Telegram API compatibility depends on careful JSON behavior.

Rules:

- use correct `json` tags;
- do not rename Telegram fields unnecessarily;
- optional fields must work correctly with `omitempty`, pointers, or dedicated optional types;
- preserve forward compatibility where important;
- avoid `map[string]any` when a type can be described.

## 14. HTTP client

The HTTP layer must be testable.

Requirements:

- allow a custom `http.Client`;
- make the base URL configurable for tests;
- never log the bot token;
- distinguish network errors from Telegram API errors;
- implement multipart upload separately from regular JSON/form requests.

## 15. Middleware

Middleware is a chain around update handling and must not be required for simple usage.

A good baseline shape is:

```go
type Handler func(context.Context, *Update) error

type Middleware func(Handler) Handler
```

Choose the exact shape based on the current architecture and task.

Middleware may provide logging, recovery, metrics, auth/filtering, rate limiting, timeouts, or context-based dependency injection.

## 16. Testing

Every meaningful change should have tests.

Minimum checks:

```bash
go test ./...
go vet ./...
gofmt -w .
```

If a `Makefile` exists, use:

```bash
make check
```

Tests should verify behavior, not only lines of code.

Important test areas:

- type serialization/deserialization;
- successful API requests;
- Telegram API errors;
- network errors;
- context cancellation;
- long polling loop;
- webhook request validation;
- dispatcher routing;
- middleware order;
- panic recovery;
- file upload/download;
- test helpers.

## 17. Documentation

Documentation must help users start quickly.

For each major stage, update as needed:

- `README.md` when the user-facing API changes;
- GoDoc comments for exported symbols;
- examples when a new user scenario appears.

README must not promise behavior that is not implemented.

## 18. Examples

Examples should be small and runnable.

Useful examples as the project grows:

- echo bot through long polling;
- `/start` command handler;
- callback query handler;
- webhook bot;
- logging/recovery middleware;
- file upload;
- tests with fake Telegram API.

Do not add examples that are larger or more complex than the library feature being demonstrated.

## 19. Security

Security rules are mandatory.

- Never hardcode real tokens.
- Do not log bot tokens.
- Do not log full Telegram API URLs when they contain tokens.
- Do not trust incoming webhook requests without validation when secret token validation is available.
- Do not open arbitrary user paths without checks.
- Do not create SSRF paths through arbitrary download URLs.
- Do not add unsafe defaults for convenience.

## 20. Network environment

The user's local machine may route traffic through TUN/Xray. This can affect local `curl`, `go run`, Telegram API checks, SSH, local Bot API discovery, webhook URLs, and smoke scripts.

When debugging network behavior, explicitly distinguish:

- the user's local machine;
- server `vk1`;
- a local Telegram Bot API server;
- official `api.telegram.org`;
- the webhook endpoint;
- SSH tunnels.

Rules:

- first identify where the command ran: local machine or `ssh vk1`;
- for server-side checks, prefer running diagnostics through `ssh vk1`;
- `localhost` and `127.0.0.1` always refer to the machine where the command runs;
- a local Bot API server on `vk1` may require an SSH tunnel for local smoke scripts;
- the webhook URL must be reachable from the component that sends webhook requests;
- do not change library code immediately for a network error; first check routing, tunnels, firewall/listeners, systemd logs, process listeners, and webhook state;
- do not print bot tokens, webhook secrets, token-bearing URLs, or full `/bot<TOKEN>/...` endpoints.

## 21. Performance

Do not optimize prematurely.

Avoid clearly bad choices:

- unnecessary copies of large files;
- reading entire files into memory when streaming is practical;
- infinite retries without backoff and context cancellation;
- creating a new `http.Client` per request;
- global mutex bottlenecks in dispatcher code.

## 22. Compatibility and versions

The project should evolve without unnecessary public API breakage.

Before a stable release, changes are allowed, but still avoid:

- renaming public symbols without need;
- changing signatures for cosmetic reasons;
- mixing breaking changes with unrelated fixes;
- hiding potential breaking changes in the final report.

## 23. Task workflow

Each user prompt is a separate stage.

Before changes:

1. Read `AGENTS.md`.
2. Read the current user prompt.
3. Inspect the repository structure and conventions.
4. Choose the smallest change set that satisfies the task.
5. Do not implement future stages early.

During work:

- edit only necessary files;
- keep the project working;
- prefer small, understandable changes;
- do not rewrite everything without need;
- do not delete existing public API unless requested;
- do not add TODOs instead of requested implementation.

After changes:

- run `gofmt`;
- run tests;
- run `go vet` when possible;
- briefly describe what changed;
- list passed checks;
- honestly list anything not checked.

## 24. Final report format

End each stage with a concise user report.

Recommended shape:

```text
Done:
- ...

Changed files:
- ...

Checks:
- gofmt: ok
- go test ./...: ok
- go vet ./...: ok

Notes:
- ...
```

If a check was not run, do not mark it `ok`; state that it was not run and why.

User-facing chat reports may be Russian when the user communicates in Russian, but tracked repository files must remain English.

## 25. Do not do without explicit request

Do not do these unless directly requested:

- implement the entire Telegram Bot API at once;
- generate thousands of lines of types without verification;
- add ORM, web frameworks, or DI frameworks;
- add code generation without a separate decision;
- add a CLI when the task is about the library;
- add a GUI;
- add Docker/Kubernetes unless required;
- change the license;
- publish the package;
- create a release;
- add telemetry that sends data externally;
- add hidden behavior.

## 26. Conflict priority

If requirements conflict, use this order:

1. Security.
2. Correct behavior.
3. Public API stability.
4. Simple architecture.
5. Testability.
6. Performance.
7. Convenience.
8. Internal elegance.

If a conflict cannot be resolved without the owner, choose the smallest safe option and state the compromise in the final report.

## 27. Development strategy

Develop the project in stages.

Recommended order:

1. Minimal Go module scaffold.
2. Basic Telegram Bot API types.
3. Low-level API client.
4. Telegram API error handling.
5. Long polling.
6. Webhook receiver.
7. Dispatcher/router.
8. Middleware.
9. Observability hooks.
10. File upload/download.
11. Testkit.
12. Examples.
13. Security review.
14. Release preparation.

Do not skip ahead unless the current prompt explicitly requires it.

## 28. Good result criteria

A good result is not the largest amount of code.

A good result:

- compiles;
- has tests where there is behavior;
- has clear package boundaries;
- avoids unnecessary technical debt;
- is understandable for library users;
- leaves the project ready for the next stage.
