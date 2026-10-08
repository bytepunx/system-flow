---
id: T-1309
type: task
nature: remediation
title: flai release raises flai.minimum no higher than the flai that publishes
status: done
parent: S-0321
owner: alex
created: 2026-10-08T00:23:14Z
updated: 2026-10-08T08:06:47Z
transitions:
  - to: ready
    at: 2026-10-08T07:58:42Z
    by: agent-S-0321
  - to: in-progress
    at: 2026-10-08T07:58:42Z
    by: agent-S-0321
  - to: done
    at: 2026-10-08T08:06:47Z
    by: agent-S-0321
stream: S-0321
tags: [flai, release]
touches: [flai/internal/release/release.go, flai/internal/release/release_test.go, flai/cmd/release.go]
usage:
  source: log
  seconds: 485
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11206
      cache_read: 2247616
      cache_write: 52548
      cost: 1.016
---
# T-1309 flai release raises flai.minimum no higher than the flai that publishes

## Work

Remove the cause of I-0107: `release.RaiseMinimum` (`flai/internal/release/release.go`), called from `computeApplyAndTagPending` in `flai/cmd/release.go` for `flai release --pending` and `flai push --pending`, writes the release being published as `flai.minimum` in the publish commit, before its binaries exist, so the flai serving the project drops it.

Raise the minimum only to a flai release that can already be installed. Walk the `flai/v*` tags: a tag whose `flai/internal/workitem/front-matter-fields.txt` differs from its previous tag's is a candidate, and the minimum becomes the newest candidate no newer than the flai running the publish (`buildinfo`). The release being published is a candidate only for a flai at or above it. So a raise the publish cannot make yet is made by the first publish from the upgraded flai, and the minimum never rises above the flai that writes it. A build that is not a release (`scripts/flai.sh`) raises nothing. Never lower a minimum already set. Keep a first release and a fields file the previous tag lacked raising nothing, and a git that cannot say an error.

Say in the publish's output what was raised, or that a raise waits for a publish from flai X.Y.Z or newer, in place of the warning `flai.minimum raised: upgrade the host's flai once the release's binaries are built`.

Waits for nothing: it shares no path with the serve task.

## Done when

- A test in `flai/internal/release/release_test.go` reproduces I-0107: flai X publishing X+1, whose fields file changed, leaves `flai.minimum` at or below X. A later publish from flai X+1 raises it to X+1.
- A build that is not a release raises nothing, and the output says why.
- The cases of `TestRaiseMinimum` still hold where the rule above keeps them, and are rewritten where it changes them.
- `flai test` passes on the changed paths.

## Notes
