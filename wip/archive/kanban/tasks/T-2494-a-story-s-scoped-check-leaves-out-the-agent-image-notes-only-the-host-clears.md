---
id: T-2494
type: task
nature: improvement
title: A story's scoped check leaves out the agent.image notes only the host clears, unless the story changes the manifest or the layer
status: done
parent: S-0582
owner: alex
created: 2026-11-07T00:30:46Z
updated: 2026-11-07T00:40:53Z
transitions:
  - to: ready
    at: 2026-11-07T00:35:12Z
    by: agent-S-0582
  - to: in-progress
    at: 2026-11-07T00:35:13Z
    by: agent-S-0582
  - to: done
    at: 2026-11-07T00:40:53Z
    by: agent-S-0582
stream: S-0582
tags: [flai, check]
touches: [flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/check/check.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-2493]
usage:
  source: log
  seconds: 340
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 7280
      cache_read: 3445485
      cache_write: 68857
      cost: 1.3584
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 13
          cache_read: 82505
          cache_write: 1711
          cost: 0.0203
---
# T-2494 A story's scoped check leaves out the agent.image notes only the host clears, unless the story changes the manifest or the layer

## Work

Build what T-2493's ADR decides. It waits for T-2493 because that ADR settles which findings are left out and when they stay inside.

- In `flai/internal/check/image.go`, mark the docker-not-on-`PATH`, docker-could-not-say, and not-built findings of `AgentImage` so that the scoped check can tell them apart, as `unchecked` marks the `components.examples` note (`check.go`). They sit on `system-flow.yaml` but are judged by two paths: the manifest and the layer's Dockerfile.
- In `flai/internal/check/scope.go`, have `ScopeToStories` drop them unless a changed path is `system-flow.yaml` or the layer, or a changed folder holds either. Keep the not-a-file warning as it is. Update the doc comments.
- In `flai/cmd/check.go`, name the left-out findings in the `--story` help and in `recordOutside`'s comment. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Tests: in `scope_test.go`, a story that changes neither path loses each of the three findings; one that changes the manifest or the layer keeps them, and the not-a-file warning stays outside. In `cmd/check_test.go`, `--story --record-issues` records no `agent.image` issue when docker is not on `PATH`, which reproduces I-0264's instance. Adjust `image_test.go` for the mark.

## Done when

- The tests above pass with `flai test` on the changed paths.
- `flai check --story S-nnnn --record-issues` on a probe with no docker writes no issue for a story that changes neither path.
- `flai check` without `--story` still reports all four `agent.image` findings as before.

## Notes
