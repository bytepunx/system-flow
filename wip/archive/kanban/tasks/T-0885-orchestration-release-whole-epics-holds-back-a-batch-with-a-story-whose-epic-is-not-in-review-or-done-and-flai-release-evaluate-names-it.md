---
id: T-0885
type: task
nature: feature
title: orchestration.release.whole_epics holds back a batch with a story whose epic is not in review or done, and flai release --evaluate names it
status: done
parent: S-0222
owner: alex
created: 2026-10-05T04:46:20Z
updated: 2026-10-06T11:58:34Z
transitions:
  - to: ready
    at: 2026-10-06T11:53:20Z
    by: agent-S-0222
  - to: in-progress
    at: 2026-10-06T11:53:20Z
    by: agent-S-0222
  - to: done
    at: 2026-10-06T11:58:34Z
    by: agent-S-0222
stream: S-0222
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 314
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 323
      cache_read: 1983987
      cache_write: 96448
      cost: 1.3173
---
# T-0885 orchestration.release.whole_epics holds back a batch with a story whose epic is not in review or done, and flai release --evaluate names it

## Work

Add `whole_epics`, a boolean that is off when unset, to the `orchestration.release` block that S-0217's T-0809 adds to `manifest.Manifest` in `flai/internal/manifest/manifest.go`.

Extend the evaluation that S-0217's T-0812 adds in `flai/internal/release/evaluate.go`. When `whole_epics` is set, take each story accepted and not yet released whose epic is in neither `review` nor `done`. Return them as the evaluation's `held_by_epic`, each with its epic and the epic's state. While any is held, the verdict is not met under every policy, and the sentence says which epic holds the batch back. Under `judgement` the list tells the orchestrator what it must not publish. A story with no epic is never held.

Print the held stories in `flai release --evaluate` and include them in its `--json` (`flai/cmd/release.go`).

Describe `whole_epics`, its default, and what it holds back in `design/system/project-manifest.md` and `docs/operators/settings.md`, beside the `orchestration.release` keys.

This task waits for nothing in this story. It needs S-0217's `orchestration.release` block and its evaluation, which the story's `after` on S-0218 brings in. It runs with the guard task, whose paths it does not share.

## Done when

- the manifest parses `whole_epics`, and a test pins that unset means off
- fixture tests in `evaluate_test.go` pin, for `threshold`, `theme`, and `judgement`, a batch held by a story whose epic is in progress, a batch whose epics are all in review or done, a story with no epic, and `whole_epics` off
- `flai release --evaluate` prints the held stories and their epics, and changes no file, tag, or remote
- `project-manifest.md` and `settings.md` describe the key
- `go test ./internal/manifest/ ./internal/release/ ./cmd/` passes

## Notes
