---
id: T-1990
type: task
nature: remediation
title: Close I-0163 with a reason naming S-0448's fix, its commits, and the releases that carry it
status: done
parent: S-0475
owner: alex
created: 2026-10-21T12:07:08Z
updated: 2026-10-21T12:12:07Z
transitions:
  - to: ready
    at: 2026-10-21T12:11:34Z
    by: agent-S-0475
  - to: in-progress
    at: 2026-10-21T12:11:35Z
    by: agent-S-0475
  - to: done
    at: 2026-10-21T12:12:07Z
    by: agent-S-0475
stream: S-0475
tags: [issues]
touches: [design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md, design/issues/summary.md]
after: [T-1987]
usage:
  source: log
  seconds: 32
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 976
      cache_read: 360180
      cache_write: 19284
      cost: 0.2458
---
# T-1990 Close I-0163 with a reason naming S-0448's fix, its commits, and the releases that carry it

## Work

Close I-0163 in the story's worktree, as the second criterion asks. It waits for T-1987: the reason may only say the cause is fixed once T-1987 has confirmed it.

Run `flai issue close I-0163 --reason "<reason>"`. Make the reason say:

- S-0448 fixed it: the components methods take their sub-project as `sub_project` (382fab56f), and the dashboard's agent puts the channel's project key after a call's params (d94c114d4).
- flai v1.63.1 and flaiover v0.52.5 carry the fix.
- The regression tests are `TestComponentsAnswerThroughTheChannelWithTheProjectKeyInTheirParams` and `TestNoMethodTakesAParameterNamedProject` in `flai/internal/hostapi`, and the S-0448 cases in flaiover's `agent.test.ts` and `components.test.ts`.
- S-0475 confirmed the fix on main.

If T-1987 had to fix something, the reason names that fix too.

Closing the issue rewrites the issue file and regenerates `design/issues/summary.md`.

## Done when

- I-0163's status is closed, and its reason names S-0448, both commits, and both releases.
- `design/issues/summary.md` no longer lists I-0163.
- `flai check --strict` is clean.

## Notes
