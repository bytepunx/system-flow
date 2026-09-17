---
id: S-032
type: story
nature: feature
title: "flaiover observability: request logs, health, metrics, traces"
status: review
parent: E-003
owner: alex
created: 2026-09-16T23:24:18Z
updated: 2026-09-17T06:13:53Z
transitions:
  - to: ready
    at: 2026-09-17T05:52:20Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:52:20Z
    by: agent
  - to: review
    at: 2026-09-17T06:13:53Z
    by: alex
tags: [dashboard, telemetry, logging]
---

# S-032 flaiover observability: request logs, health, metrics, traces

## Goal
flaiover follows the logging and telemetry conventions: pino request and lifecycle logs on stdout, `/_health`, `/_ready`, `/metrics` with `flaiover_` golden signals and build info, OpenTelemetry traces to a collector, and a local collector in the compose stack.

## Acceptance criteria
- [x] Request log line per request with trace_id, method, route, status, duration; lifecycle events for start, index built, watcher events
- [x] `/_health` always 200 while the process runs; `/_ready` checks the mount and index with timeouts and names the failing dependency
- [x] `/metrics` exposes request counter with status label, latency histogram, in-flight gauge, `flaiover_build_info`
- [x] Traces exported via OTLP when `OTEL_EXPORTER_OTLP_ENDPOINT` is set; W3C context propagated
- [x] docker compose for local development includes a collector and a viewer; docs/operators documents the endpoints and variables
- [x] design/tech records pino and the OpenTelemetry packages with versions

## Tasks
- T-101 pino logging: request line with trace_id, lifecycle events, LOG_LEVEL and LOG_FORMAT
- T-102 /_health, /_ready with dependency checks, /metrics with prom-client golden signals and build_info
- T-103 OpenTelemetry traces via OTLP when configured, W3C context propagation, one span per request
- T-104 compose stack with a collector and viewer; operator docs; tech pins; tests

## Notes
- Created by S-030. Depends on S-011 (dashboard scaffold). Nature feature; bump flaiover minor on acceptance.
