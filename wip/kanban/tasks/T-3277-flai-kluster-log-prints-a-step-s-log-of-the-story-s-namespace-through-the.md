---
id: T-3277
type: task
nature: feature
title: flai kluster log prints a step's log of the story's namespace through the reader the host method kluster.log uses
status: in-progress
parent: S-0763
owner: alex
created: 2026-11-13T09:03:01Z
updated: 2026-11-13T09:08:16Z
transitions:
  - to: ready
    at: 2026-11-13T09:08:16Z
    by: agent-S-0763
  - to: in-progress
    at: 2026-11-13T09:08:16Z
    by: agent-S-0763
stream: S-0763
tags: [cli]
touches: [flai/internal/kluster/log.go, flai/internal/kluster/log_test.go, flai/internal/hostapi/kluster.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
usage:
  source: log
  seconds: 2335
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 38698
      cache_read: 3867035
      cache_write: 438602
      cost: 4.0417
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 15
          cache_read: 81709
          cache_write: 2328
          cost: 0.0242
---
# T-3277 flai kluster log prints a step's log of the story's namespace through the reader the host method kluster.log uses

## Work

Criterion 4. The host method `kluster.log` reads a step's log in `flai/internal/hostapi/kluster.go` (`klusterLogIn` checks the namespace and the step, `klusterLogDo` reads at most 64 KB from an offset and the step's outcome), unexported and tied to the channel. So that `flai kluster log` prints a log as `kluster.log` reads it, move the reading into the kluster package and have both call it.

- A new `flai/internal/kluster/log.go`: the step name check (`build`, `build-<name>`, `load`, `deploy`, `teardown`), and `ReadLog(root, namespace, step, offset)` answering what `kluster.log` answers today: the path, `exists`, `offset`, `next`, `size`, `text` ending on a whole character, `reset`, and the step's outcome from the record. Tests in `log_test.go`.
- `hostapi/kluster.go`: `klusterLogIn` and `klusterLogDo` call it, with every field and refusal as before. Its tests in `hostapi/kluster_test.go` pass unchanged.
- `flai/cmd/kluster.go`: `flai kluster log S-nnnn --step <step> [--json]`. It reads the story's namespace's log from 0 to its end in `ReadLog`'s chunks and prints it; `--json` prints the answers' fields with the whole text. It refuses a step `kluster.log` refuses, and says when the log is absent, exiting 1. It needs no `kluster:` block and no flai serve, as `kluster.log` reads only files. A session with a story of its own is refused another story, as `--job` is (`sessionStory`). Tests in `kluster_test.go`.
- Run `scripts/flai-reference.sh`, which writes `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` from the command's help, so `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` pass.

Waits for nothing: no other task of the first layer changes these files.

## Done when

- `flai kluster log S-nnnn --step deploy` prints the namespace's deploy log as `kluster.log` reads it, and refuses another story for a session with its own, and a step `kluster.log` refuses.
- `kluster.log` answers as before, through the shared reader.
- `go test ./internal/kluster/ ./internal/hostapi/ ./cmd/` passes, `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` among them.

## Notes
