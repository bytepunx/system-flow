---
title: Logging and telemetry libraries
updated: 2026-09-16
status: active
---

# Logging and telemetry libraries

What the logging and telemetry conventions imply for each sub-project. Decided per sub-project when its story lands; recorded here so the choice is visible before the code exists.

| Sub-project | Logging | Telemetry | Story |
|-------------|---------|-----------|-------|
| flai | Go standard library `log/slog` (Go 1.21+): text handler on a terminal, JSON handler otherwise, to stderr. No third-party logger. | None by default; an optional root span via the OpenTelemetry Go SDK only when `OTEL_*` is set, decided in the story. | S-031 |
| flaiover | `pino` (JSON lines to stdout, child loggers per component, request logging via `pino-http` or an equivalent hook in SvelteKit). | `@opentelemetry/sdk-node` with the OTLP exporter for traces; `prom-client` for `/metrics`; health endpoints are plain routes. Versions pinned when S-032 lands. | S-032 |
| template | Ships the conventions only; no library. | Compose stack for local testing gains an OpenTelemetry collector and a viewer when the first service needs one. | S-032 |

## Why

- `log/slog` is standard, structured, and has handlers for both formats the convention requires; a third-party logger would add a dependency for no capability.
- `pino` is the fastest structured logger in Node, JSON by default, and the usual choice next to OpenTelemetry.
- OpenTelemetry is the vendor-neutral standard the telemetry convention is written against, so exporters can change without touching code.

## Considered

- `zerolog` or `zap` for Go: faster than slog in benchmarks, unnecessary for a CLI.
- `winston` for Node: more configurable, slower, and pushes toward prose messages.
- Custom `/metrics` formatting: `prom-client` handles histograms and the exposition format correctly; hand-rolling it is a known source of subtle bugs.
