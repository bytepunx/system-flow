---
id: T-3171
type: task
nature: feature
title: flai image ensure pulls or builds each missing or stale harness image, and flai image build builds without the pull
status: done
parent: S-0740
owner: alex
created: 2026-11-12T15:59:50Z
updated: 2026-11-13T12:42:01Z
transitions:
  - to: ready
    at: 2026-11-13T12:26:13Z
    by: agent-S-0740
  - to: in-progress
    at: 2026-11-13T12:26:13Z
    by: agent-S-0740
  - to: done
    at: 2026-11-13T12:42:01Z
    by: agent-S-0740
stream: S-0740
tags: [cli]
touches: [flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/cmd/image.go, flai/cmd/image_test.go]
after: [T-3169, T-3170]
usage:
  source: log
  seconds: 948
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 13176
      cache_read: 3284654
      cache_write: 97453
      cost: 1.4956
---
# T-3171 flai image ensure pulls or builds each missing or stale harness image, and flai image build builds without the pull

## Work

`flai/internal/agentimage/ensure.go`:

- takes S-0735's status;
- for each harness image missing or stale and not the operator's own, pulls it with T-3169 and falls back to T-3170's build as S-0733's ADR allows;
- answers, per image, which way it took, the log file, and the error.

It never touches an image the host's `--image` names. Add the subcommands `ensure` and `build` to `flai/cmd/image.go`, with `--json`, streaming the output. `build` takes `--inner` to build one harness. Both exit non-zero on a failure, naming the log file.

Tests over a fake docker and release client cover:

- nothing to do;
- a pull;
- a fallback build;
- an own image left alone;
- a failure that names its log.

It waits for T-3169 and T-3170, whose pull and build it chooses between.

## Done when

- `go test ./internal/agentimage/ ./cmd/ -run 'Ensure|Image'` passes.
- After an `ensure` over the fake docker, `flai image status` reports every harness image current.

## Notes
