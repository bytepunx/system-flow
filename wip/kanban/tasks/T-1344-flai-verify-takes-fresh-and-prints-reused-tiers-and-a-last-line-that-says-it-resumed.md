---
id: T-1344
type: task
nature: improvement
title: flai verify takes --fresh and prints reused tiers and a last line that says it resumed
status: done
parent: S-0341
owner: alex
created: 2026-10-08T08:05:15Z
updated: 2026-10-08T08:49:51Z
transitions:
  - to: ready
    at: 2026-10-08T08:43:09Z
    by: agent-S-0341
  - to: in-progress
    at: 2026-10-08T08:43:09Z
    by: agent-S-0341
  - to: done
    at: 2026-10-08T08:49:51Z
    by: agent-S-0341
stream: S-0341
tags: [cli]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1343]
usage:
  source: log
  seconds: 402
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 161
      cache_read: 1594267
      cache_write: 86887
      cost: 0.7588
---
# T-1344 flai verify takes --fresh and prints reused tiers and a last line that says it resumed

## Work

The command's side. It waits for T-1343, which adds the `Reused` state, `reused_from`, and the `Fresh` option it passes on.

- Add `--fresh` to `flai verify`, passed to the story run as `Fresh`; its help says it forces a full run. It has no effect with `--last`, which runs nothing.
- In `verifyText`, print a reused tier as `reused tier <name> (from <time>)`, with the time of the run it comes from, in UTC as flai prints timestamps.
- Make the last line say when tiers were reused, for a pass and for a stop: `verify: S-nnnn passed every step (reused N tiers from <time>)`, and the same after `stopped at <step> (...)`. The close-out reads it (T-1347).
- `--json` prints the report as it is, so it carries the reused states and `reused_from` without more code; check it does.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- Tests in `flai/cmd/verify_test.go` cover `--fresh` reaching the run, a resumed run's text lines and last line, and its `--json` carrying `reused` and `reused_from`.
- `docs/users/flai-reference.md` lists `--fresh`, and the reference test passes.
- `flai test flai/cmd/verify.go flai/cmd/verify_test.go docs/users/flai-reference.md` passes.

## Notes

`verifyText` and the last line are in `flai/cmd/verify.go`, not in `flai/internal/verify/text.go`, which prints `flai test`'s result.
