---
id: T-0104
type: task
nature: feature
title: "compose stack with a collector and viewer; operator docs; tech pins; tests"
status: done
parent: S-0032
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
stream: S-0032
tags: [dashboard, docs]
---

# T-0104 compose stack with a collector and viewer; operator docs; tech pins; tests

## Work
flaiover/compose.yaml runs the dashboard next to grafana/otel-lgtm (collector, Tempo, Prometheus, Grafana) with the endpoint wired; docs/operators telemetry reference (endpoints, metric names, variables, compose usage); design/tech/observability.md pins; I-0010 for the Node 25 host mismatch.

## Done when
Docs match; compose config validates; tiers green.

## Notes
