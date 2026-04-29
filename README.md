# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is currently in early development. At this stage it only contains a minimal Go module and a root package skeleton. Telegram Bot API methods, transports, dispatching, middleware, and test utilities are not implemented yet.

## Current status

- Minimal Go module: present.
- Root Go package: present.
- Telegram Bot API implementation: not implemented.
- Public API stability: not guaranteed before the first stable release.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
