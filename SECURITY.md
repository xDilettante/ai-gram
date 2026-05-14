# Security Policy

## Reporting vulnerabilities

Please report security issues privately through GitHub Security Advisories when available for this repository. If advisories are not available, open a minimal issue that states the affected area without including secrets, tokens, private chat IDs, webhook URLs with secrets, payment payloads, Passport data, or exploit details.

Do not disclose real bot tokens, webhook secrets, managed bot tokens, private keys, payment payloads, Passport data, or token-bearing URLs in public issues, pull requests, logs, screenshots, or examples.

## Supported versions

`ai-gram` is currently pre-1.0. The supported development targets are the current `main` branch and the latest regular public release.

| Version | Support status |
| --- | --- |
| `main` | Supported for current development and fixes |
| Latest regular release | Supported for user reports and security fixes |
| Older pre-v1 releases | Best effort only |

Breaking changes may still happen before v1.0, but token redaction, safe diagnostics, and private disclosure expectations should remain stable.

## What to report privately

Report privately when an issue could expose, mishandle, or amplify sensitive data or bot control. Examples include:

- bot token or webhook secret disclosure;
- token-bearing URL leakage in errors, logs, examples, or diagnostics;
- unsafe handling of Passport, payment, Stars, gift, business, or managed-token payloads;
- request signing, webhook secret validation, or certificate handling mistakes;
- file path or file payload disclosure in upload/download helpers;
- bugs that could make destructive/admin methods easier to call unintentionally.

For normal bugs without a security impact, use the public issue templates and follow [`SUPPORT.md`](SUPPORT.md).

## Sensitive areas

Extra care is required around:

- payments, Stars, paid media, and gifts;
- Telegram Passport encrypted payloads;
- managed bot token methods;
- business APIs and business message payloads;
- webhook secret validation and certificate upload;
- admin or destructive chat moderation methods;
- file upload and download paths.
