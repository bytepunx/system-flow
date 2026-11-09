---
id: T-2756
type: task
nature: improvement
title: telemetry.md and code-quality.md say what flai and flaiover expose and where flai's integration tests live
status: in-progress
parent: S-0626
owner: alex
created: 2026-11-09T17:13:09Z
updated: 2026-11-09T17:27:15Z
transitions:
  - to: ready
    at: 2026-11-09T17:27:14Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:27:15Z
    by: system-flow
stream: S-0626
tags: [conventions, docs]
touches: [design/conventions/telemetry.md, design/conventions/code-quality.md]
---
# T-2756 telemetry.md and code-quality.md say what flai and flaiover expose and where flai's integration tests live

## Work

Both sentences are under `## Project additions`, below the marker, so they are this project's to change. The baseline and the template's copies stay as they are.

- `design/conventions/telemetry.md` line 37: flai has no OpenTelemetry module (`flai/go.mod`) and reads no `OTEL_*` variable, and S-0031 built logging, not spans. Say that flai emits no traces, and that `flai serve`, which runs for days, serves pprof behind `FLAI_PPROF_ADDR` (`flai/cmd/serve_pprof.go`, `docs/operators/settings.md`). Drop "flai exposes nothing: it is a short-lived command".
- `design/conventions/telemetry.md` line 38: flaiover's `/_ready` checks `host_flai` and `project` through flai on the host (`flaiover/src/routes/_ready/+server.ts:17-31`, `ready.test.ts`). Nothing is mounted since ADR-0031. Say so, as `design/system/flaiover-dashboard.md` line 225 already does. Keep the true clause on `@opentelemetry/sdk-node` traces.
- `design/conventions/code-quality.md` line 85: `flai/tests/integration/` exists, with `mcp_http_test.go`, `golangci_lock_test.go`, `forge_github_test.go`, and `opencode_test.go`. There is no root `tests/smoke/`. Say what is there and what each holds.

Bump each document's `updated`.

It waits for no task: no other task changes these files.

## Done when

- Neither project additions section names an OTEL span in flai, a mounted repository, or a `flai/tests/integration/` still to be made.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
