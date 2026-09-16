---
title: Logging
updated: 2026-09-16
audience: agent
order: 110
status: active
---

# Logging

What is logged, at which level, in what shape, and what never appears in a log.

## Rules

- Logs are structured events: one event per line, machine-parseable (JSON in deployed environments, key-value text locally), never free prose.
- Every event carries `ts` (UTC, RFC 3339 with milliseconds), `level`, `component` (the package or service), and `msg`. Events on a request or job path also carry the correlation ID (`trace_id` when tracing is on, otherwise `request_id` or `job_id`).
- `msg` is a static, lower-case phrase that describes what failed so that the operator knows what is affected and where to look. Variable data goes in named fields, never interpolated into the message, so events can be grouped and counted.
- Five levels, with fixed meanings: `fatal` means the process cannot continue and exits; `error` means a human or an automatic recovery must act and the operation failed; `warn` means something unexpected was handled and the operation continued; `info` marks lifecycle and state changes (start, stop, configuration loaded, job finished, external call failed and retried); `debug` is diagnostic detail, off by default.
- Errors log once, at the boundary where they are handled, with an `err` field carrying the full error chain. Do not log an error and also return it to a caller that will log it again.
- Never log secrets, credentials, tokens, session identifiers, personal data, full request or response bodies, or anything the safety convention forbids in the repository. Redact at the source, not in the collector.
- No per-item logging inside loops at `info` or above; log the aggregate (count, duration, failures) when the loop ends.
- Command-line tools log to stderr and keep stdout for their output. Services log to stdout for the platform's collector. Nothing writes log files of its own.
- The level is set by `LOG_LEVEL` and, for command-line tools, a `--verbose` flag that selects `debug`; the format by `LOG_FORMAT` (`json` or `text`), defaulting to `json` when not attached to a terminal.
- Log output is part of the interface: a test that depends on a log line pins the field names, not the message text.
- If an error is recoverable (i.e. a request fails), it should not result in an exit.
- Fatal level log entries should be issued at any point where a detected issue or caught error cannot be recovered from and should exit with non-zero.

## When in doubt

- If you would page someone for it, `error`. If you would want to know it happened next week, `info`. If only the person debugging cares, `debug`.
- If a field might contain user data, leave it out.

<!-- system-flow:end-of-baseline -->

## Project additions
- flai logs with the standard library `log/slog` to stderr: text handler on a terminal, JSON otherwise; `--verbose` selects debug; `LOG_LEVEL` and `LOG_FORMAT` are honoured. Component is the Go package. `fatal` is a custom slog level above error, emitted once by the command boundary before a non-zero exit; a failed sub-step that the command recovers from is `warn`. Implemented in S-031.
- flaiover logs with `pino` to stdout: one line per request with `trace_id`, method, route, status, and duration, plus lifecycle events. Implemented in S-032.
- Narratives and work items never receive log output; summaries only (safety.md).
