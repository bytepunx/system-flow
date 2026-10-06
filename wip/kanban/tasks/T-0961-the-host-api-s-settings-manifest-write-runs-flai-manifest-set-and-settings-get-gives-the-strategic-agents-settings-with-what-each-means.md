---
id: T-0961
type: task
nature: feature
title: The host API's settings.manifest write runs flai manifest set, and settings.get gives the strategic agents' settings with what each means
status: done
parent: S-0229
owner: alex
created: 2026-10-05T05:46:42Z
updated: 2026-10-06T21:48:35Z
transitions:
  - to: ready
    at: 2026-10-06T21:41:10Z
    by: agent-S-0229
  - to: in-progress
    at: 2026-10-06T21:41:10Z
    by: agent-S-0229
  - to: done
    at: 2026-10-06T21:48:35Z
    by: agent-S-0229
stream: S-0229
tags: [flai]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/contract_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flaiover/src/lib/server/agent.ts]
after: [T-0952]
usage:
  source: log
  seconds: 445
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 82
      output: 414
      cache_read: 4060443
      cache_write: 156298
      cost: 1.8452
---
# T-0961 The host API's settings.manifest write runs flai manifest set, and settings.get gives the strategic agents' settings with what each means

## Work

The dashboard reaches the manifest only through flai on the host (ADR-0029), so the panels need a write and a read.

- `flai/internal/hostapi/settings.go`: `settings.manifest`, a project write gated by the `settings` host action like `settings.default_agent`. It takes `{set: {key: value}, unset: [key]}`, checks each key is in T-0952's catalog and each value is a scalar or, for an agent key, an agent `Agent.Validate` passes, and builds `manifest set ... --autocommit --trailer=<Trailer> --json`. A refusal flai gives comes back as a `channel.Error` whose data carries each `{field, reason}`, so the route can answer 422 with them. With `settings` off it answers the `Disabled` error, naming `flai serve enable settings`.
- `flai/cmd/serve_actions.go`: `settings.get` gains a `strategic` block with, for each catalog key, its value as the manifest has it (or that it is unset), its default, its kind and values, its sentence, and for a permission its risk; and whether the project's `settings` action is on, with the command that turns it on. Keep `planning` (`planningTriggers`) as it is for the Settings page.
- `flaiover/src/lib/server/agent.ts`: add `settings.manifest` to `REQUIRED_METHODS`, since `contract_test.go` fails when that list and flai's table differ.

This task waits for T-0952, for the command it runs and the catalog it reads.

## Done when

- `writes_test.go` builds the command line for a save, refuses a key outside the catalog, and answers `Disabled` with `settings` off
- a refusal from the command reaches the caller with its field and reason, in a test
- `serve_actions_test.go` reads the `strategic` block of `settings.get` for a fixture manifest, values and defaults
- `contract_test.go` passes, and `go test ./internal/hostapi/ ./cmd/` passes

## Notes
