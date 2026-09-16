---
id: S-032
type: story
nature: feature
title: "flaiover observability: request logs, health, metrics, traces"
status: backlog
parent: E-003
owner: alex
created: 2026-09-16T23:24:18Z
updated: 2026-09-16T23:24:18Z
transitions: []
tags: [dashboard, telemetry, logging]
---

# S-032 flaiover observability: request logs, health, metrics, traces

## Goal
flaiover follows the logging and telemetry conventions: pino request and lifecycle logs on stdout, `/healthz`, `/readyz`, `/metrics` with `flaiover_` golden signals and build info, OpenTelemetry traces to a collector, and a local collector in the compose stack.

## Acceptance criteria
- [ ] Request log line per request with trace_id, method, route, status, duration; lifecycle events for start, index built, watcher events
- [ ] `/healthz` always 200 while the process runs; `/readyz` checks the mount and index with timeouts and names the failing dependency
- [ ] `/metrics` exposes request counter with status label, latency histogram, in-flight gauge, `flaiover_build_info`
- [ ] Traces exported via OTLP when `OTEL_EXPORTER_OTLP_ENDPOINT` is set; W3C context propagated
- [ ] docker compose for local development includes a collector and a viewer; docs/operators documents the endpoints and variables
- [ ] design/tech records pino and the OpenTelemetry packages with versions

## Tasks

## Notes
- Created by S-030. Depends on S-011 (dashboard scaffold). Nature feature; bump flaiover minor on acceptance.
