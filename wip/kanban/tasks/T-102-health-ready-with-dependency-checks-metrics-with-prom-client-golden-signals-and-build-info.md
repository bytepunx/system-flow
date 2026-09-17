---
id: T-102
type: task
nature: feature
title: "/_health, /_ready with dependency checks, /metrics with prom-client golden signals and build_info"
status: done
parent: S-032
owner: alex
created: 2026-09-17T05:52:20Z
updated: 2026-09-17T06:13:51Z
transitions:
  - to: ready
    at: 2026-09-17T06:13:51Z
    by: alex
  - to: in-progress
    at: 2026-09-17T06:13:51Z
    by: alex
  - to: done
    at: 2026-09-17T06:13:51Z
    by: alex
stream: S-032
tags: [dashboard, telemetry]
---

# T-102 /_health, /_ready with dependency checks, /metrics with prom-client golden signals and build_info

## Work
/_health always 200; /_ready checks the manifest and item listing with a two second timeout each and reports flai as writable or read-only, 503 naming the failing check; /metrics exposes flaiover_http_requests_total{method,route,status}, flaiover_http_request_duration_seconds, flaiover_http_requests_in_flight, flaiover_build_info{version,commit}, and default process metrics with the flaiover_ prefix, using @prometheus-io/client (prom-client's successor). Route labels are SvelteKit route ids so IDs and paths never become labels.

## Done when
Unit tests for the metrics registry and labels; endpoints verified on the served build.

## Notes
