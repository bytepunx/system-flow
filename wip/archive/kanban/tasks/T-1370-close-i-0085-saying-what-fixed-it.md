---
id: T-1370
type: task
nature: remediation
title: Close I-0085 saying what fixed it
status: done
parent: S-0288
owner: alex
created: 2026-10-08T08:42:39Z
updated: 2026-10-08T09:06:56Z
transitions:
  - to: ready
    at: 2026-10-08T09:06:46Z
    by: agent-S-0288
  - to: in-progress
    at: 2026-10-08T09:06:46Z
    by: agent-S-0288
  - to: done
    at: 2026-10-08T09:06:56Z
    by: agent-S-0288
stream: S-0288
tags: [issues]
touches: [design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md, design/issues/summary.md]
after: [T-1369]
usage:
  source: log
  seconds: 10
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1773
      cache_read: 424168
      cache_write: 27138
      cost: 0.324
---
# T-1370 Close I-0085 saying what fixed it

## Work

Second layer: it waits for T-1369, because the reason must name the fix T-1369 made and the runs that showed it.

- In the story's worktree, run `flai issue close I-0085 --reason "<reason>"`. The reason names the cause (a background manifest read started by an earlier change, cached by the `Repo`, while `setManifest` truncated the manifest and never reported the change), the fix in `flaiover/src/lib/server/notify.test.ts`, and the reproducing test.
- Check that `design/issues/summary.md` no longer lists I-0085.

## Done when

- I-0085's front matter says `status: closed` with the reason.
- `design/issues/summary.md` no longer lists I-0085.
- `flai check --strict` is clean for the story.

## Notes

`design/issues` is in the manifest's `claims.shared`, so these touches hold no other story.
