---
title: Telemetry
updated: 2026-09-16
audience: agent
order: 120
status: draft
---

# Telemetry

Metrics, traces, and health signals every service emits, and how they are named. Proposed by the agent in S-030; the operator confirms or edits before this leaves draft.

## Rules

- Every service exposes `/healthz` (process is alive), `/readyz` (dependencies reachable, safe to receive traffic), and `/metrics` (Prometheus exposition format). Command-line tools expose none of these.
- Metrics and traces follow OpenTelemetry: the OTLP exporter is configured by the standard `OTEL_*` environment variables, and W3C trace context is propagated on every outbound call so the `trace_id` in logs joins spans across services.
- The minimum metric set per service is the four golden signals: request rate and error rate (a counter with a `status` label), latency (a histogram in seconds), and saturation (in-flight requests or queue depth), plus a `build_info` gauge with `version` and `commit` labels.
- Names are `snake_case`, prefixed with the service name, and end with the unit: `_seconds`, `_bytes`, `_total` for counters, `_ratio` for 0 to 1 values. Histogram buckets are chosen per metric and stated in the code.
- Labels are a fixed, low-cardinality set declared with the metric. Never label with identifiers, user data, free text, or anything unbounded.
- Spans are named `<component>.<operation>`, carry attributes from the OpenTelemetry semantic conventions where one exists, and set span status to error when the operation fails, with the error recorded once.
- Sampling is configured, not hard-coded: head sampling via `OTEL_TRACES_SAMPLER` and `OTEL_TRACES_SAMPLER_ARG`, with 100 percent locally and a stated default in deployed environments.
- A change that adds a metric, label, or span name updates the service's telemetry reference in `docs/operators`, and alert rules live in the repository with the service and are reviewed with the code that emits the signals.
- Health endpoints never require authentication and never touch data stores in `/healthz`; `/readyz` checks each dependency with a short timeout and reports which one failed.
- Local runs see their own telemetry: the compose stack for local testing includes a collector, and a service started locally exports to it by default.

## When in doubt

- If a dashboard could not be built from it, it is a log, not a metric.
- If a label could take more than a few dozen values, it is an attribute on a span or a field in a log, not a label.

<!-- system-flow:end-of-baseline -->

## Project additions
- flai exposes nothing: it is a short-lived command. It emits no traces unless `OTEL_*` variables are set, and then only one root span per invocation (S-031, optional).
- flaiover exposes `/healthz`, `/readyz` (checks the mounted repository is readable and the index is built), and `/metrics` with `flaiover_` prefixed golden signals and `flaiover_build_info`. Traces via `@opentelemetry/sdk-node` to the collector named by `OTEL_EXPORTER_OTLP_ENDPOINT`. Implemented in S-032.
- Local telemetry stack (collector plus a viewer) is added to the flaiover compose file in S-032; there is no cluster here.
