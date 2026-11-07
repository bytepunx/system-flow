---
id: T-2594
type: task
nature: improvement
title: TestSettingsIndexIsComplete checks the settings index's keys under -short
status: done
parent: S-0592
owner: alex
created: 2026-11-07T22:02:31Z
updated: 2026-11-07T22:19:19Z
transitions:
  - to: ready
    at: 2026-11-07T22:06:23Z
    by: agent-S-0592
  - to: in-progress
    at: 2026-11-07T22:06:27Z
    by: agent-S-0592
  - to: done
    at: 2026-11-07T22:19:19Z
    by: agent-S-0592
stream: S-0592
tags: [flai, tests]
touches: [flai/cmd/settings_doc_test.go]
usage:
  source: log
  seconds: 772
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 4343
      cache_read: 2618728
      cache_write: 53138
      cost: 1.0183
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 19330
          cache_write: 573
          cost: 0.005
---
# T-2594 TestSettingsIndexIsComplete checks the settings index's keys under -short

## Work

`TestSettingsIndexIsComplete` in `flai/cmd/settings_doc_test.go` skips whole under `-short`, the flag the `go-test` tier passes, so `flai test` on a change to the configuration's or the manifest's shape never checks the settings index (I-0304's second instance). Its key checks, over `config.Config` and `manifest.Manifest` by reflection, read one page and cost nothing; only the walks for environment variables read the sources.

Split it: the key checks of `## flai configuration` and `## Project manifest` run under `-short`, as a test of their own or with the skip moved below them, and the environment-variable walks keep the skip. Time the test without `-short` first. If the whole test runs in under a second once the flaiover walk leaves out `node_modules`, drop the skip instead, and say which you chose in the narrative's `## Decisions`.

This task waits for no other: it changes one file no other task touches.

## Done when

- `go test -short -run '^TestSettingsIndex' ./cmd` in `flai` runs the key checks and fails when a configuration or manifest key has no row in `docs/operators/settings.md`, checked by removing a row by hand and restoring it.
- `flai test flai/cmd/settings_doc_test.go` passes.

## Notes
