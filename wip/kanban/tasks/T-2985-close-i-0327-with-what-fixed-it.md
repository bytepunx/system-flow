---
id: T-2985
type: task
nature: remediation
title: Close I-0327 with what fixed it
status: done
parent: S-0693
owner: alex
created: 2026-11-10T10:29:34Z
updated: 2026-11-11T05:11:09Z
transitions:
  - to: ready
    at: 2026-11-11T05:11:02Z
    by: agent-S-0693
  - to: in-progress
    at: 2026-11-11T05:11:02Z
    by: agent-S-0693
  - to: done
    at: 2026-11-11T05:11:09Z
    by: agent-S-0693
stream: S-0693
tags: [issues]
touches: [design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/summary.md]
after: [T-2982]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 222
      cache_read: 74569
      cache_write: 4538
      cost: 0.0557
---
# T-2985 Close I-0327 with what fixed it

## Work

Run `flai issue close I-0327 --reason "<reason>"` in the story's worktree. The reason names S-0693 and says that `TestRepositoryManifestLintsChangedFlaioverFiles` no longer skips under `-short`, so the go-test tier catches a tier it must list. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

Waits for T-2982: the reason states what that task changed, so the fix must be in and passing first.

## Done when

- `flai issue show I-0327` reports it closed, with the reason naming S-0693 and the dropped gate.
- `design/issues/summary.md` no longer lists I-0327 as open.
- `flai check --strict` is clean for the story.

## Notes
