---
id: T-103
type: task
nature: feature
title: OpenTelemetry traces via OTLP when configured, W3C context propagation, one span per request
status: done
parent: S-032
owner: alex
created: 2026-09-17T05:52:20Z
updated: 2026-09-17T06:13:52Z
transitions:
  - to: ready
    at: 2026-09-17T06:13:52Z
    by: alex
  - to: in-progress
    at: 2026-09-17T06:13:52Z
    by: alex
  - to: done
    at: 2026-09-17T06:13:52Z
    by: alex
stream: S-032
tags: [dashboard, telemetry]
---

# T-103 OpenTelemetry traces via OTLP when configured, W3C context propagation, one span per request

## Work
src/lib/server/otel.ts starts the NodeSDK with the OTLP/HTTP trace exporter only when OTEL_EXPORTER_OTLP_ENDPOINT is set; one SERVER span per request named method and path with W3C context extracted from the headers, http.route and status attributes, error status on 5xx; trace_id surfaced to the log line.

## Done when
Unit test for context extraction and trace id; span attributes checked with an in-memory exporter.

## Notes
