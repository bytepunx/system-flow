---
id: T-3524
type: task
nature: remediation
title: Close I-0327 with what fixed it
status: done
parent: S-0827
owner: alex
created: 2026-11-17T08:28:28Z
updated: 2026-11-17T08:45:13Z
transitions:
  - to: ready
    at: 2026-11-17T08:45:04Z
    by: agent-S-0827
  - to: in-progress
    at: 2026-11-17T08:45:04Z
    by: agent-S-0827
  - to: done
    at: 2026-11-17T08:45:13Z
    by: agent-S-0827
stream: S-0827
tags: [issues]
touches: [design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/summary.md]
after: [T-3523]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 360
      cache_read: 224300
      cache_write: 7294
      cost: 0.1104
  strategic:
    - kind: orchestrator
      seconds: 9
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 50
          cache_read: 217510
          cache_write: 2536
          cost: 0.0682
---
# T-3524 Close I-0327 with what fixed it

## Work

Run `flai issue close I-0327 --reason "<reason>"` in the story's worktree. The reason names S-0827 and says that `TestRepositoryManifestLintsChangedFlaioverFiles` no longer skips under `-short`, so the go-test tier catches a tier it must list. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

Waits for T-3523: the reason states what that task changed, so the fix must be in and passing first.

## Done when

- `flai issue show I-0327` reports it closed, with the reason naming S-0827 and the dropped gate.
- `design/issues/summary.md` no longer lists I-0327 as open.
- `flai check --strict` is clean for the story.

## Notes
