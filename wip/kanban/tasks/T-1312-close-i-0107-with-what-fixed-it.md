---
id: T-1312
type: task
nature: remediation
title: Close I-0107 with what fixed it
status: done
parent: S-0321
owner: alex
created: 2026-10-08T00:23:40Z
updated: 2026-10-08T08:11:52Z
transitions:
  - to: ready
    at: 2026-10-08T08:11:30Z
    by: agent-S-0321
  - to: in-progress
    at: 2026-10-08T08:11:31Z
    by: agent-S-0321
  - to: done
    at: 2026-10-08T08:11:52Z
    by: agent-S-0321
stream: S-0321
tags: [issues]
touches: [design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md, design/issues/summary.md]
after: [T-1311]
usage:
  source: log
  seconds: 21
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 115
      cache_read: 1387692
      cache_write: 3447
      cost: 0.6236
---
# T-1312 Close I-0107 with what fixed it

## Work

Run `flai issue close I-0107 --reason "<reason>"` in the story's worktree. The reason names S-0321 and what fixed it: a publish raises `flai.minimum` no higher than the flai that runs it (T-1309), and a served project that stops loading lets its requests in flight finish (T-1310). The command regenerates `design/issues/summary.md`.

Leave the duplicated `## 1.0.67` heading in `template/CHANGELOG.md`, which the issue names as a separate defect, out of the reason. If no issue records it, say so in the story's notes for the operator.

Waits for T-1311, so that the issue closes only once the fix is built and described.

## Done when

- I-0107 is closed, with a reason that names what fixed it.
- `design/issues/summary.md` no longer lists I-0107 as open.
- `flai check --strict` is clean for the story.

## Notes
