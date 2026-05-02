# Security Policy

## Reporting vulnerabilities

Please report security issues privately through GitHub Security Advisories when available for this repository. If advisories are not available, open a minimal issue that states the affected area without including secrets, tokens, private chat IDs, webhook URLs with secrets, payment payloads, Passport data, or exploit details.

Do not disclose real bot tokens, webhook secrets, managed bot tokens, private keys, payment payloads, Passport data, or token-bearing URLs in public issues, pull requests, logs, screenshots, or examples.

## Supported versions

`ai-gram` is currently pre-1.0. The supported development target is the current `main` branch and the latest release once releases are created.

## Sensitive areas

Extra care is required around:

- payments, Stars, paid media, and gifts;
- Telegram Passport encrypted payloads;
- managed bot token methods;
- business APIs and business message payloads;
- webhook secret validation and certificate upload;
- admin or destructive chat moderation methods;
- file upload and download paths.
