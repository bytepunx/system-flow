---
id: T-2060
type: task
nature: improvement
title: I-0152 is closed with what fixed it and the story's criteria are verified and ticked
status: done
parent: S-0468
owner: alex
created: 2026-10-22T21:06:52Z
updated: 2026-10-23T04:59:05Z
transitions:
  - to: ready
    at: 2026-10-23T04:58:40Z
    by: agent-S-0468
  - to: in-progress
    at: 2026-10-23T04:58:41Z
    by: agent-S-0468
  - to: done
    at: 2026-10-23T04:59:05Z
    by: agent-S-0468
stream: S-0468
tags: [cli]
touches: [design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md, design/issues/summary.md]
after: [T-2059]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 3149
      cache_read: 616123
      cache_write: 20094
      cost: 0.3031
---
# T-2060 I-0152 is closed with what fixed it and the story's criteria are verified and ticked

## Work

Waits for everything else: the close names what was built and documented.

- In the story's worktree (issue commands write under the checkout they run in), run `flai issue close I-0152 --reason "…"` with a reason that names the ADR, the refusal of a task's folder touch, the planner's refusal without a Planning note, the `touches.folder` and `wip.after-hub` findings, and the `after` rule in the convention and prompt. flai regenerates `design/issues/summary.md`; both files are inside `claims.shared`, so the change holds no story.
- Verify the story's criteria and tick each through flai (`flai criteria tick S-0468 <n>`, or `criteria_tick`): criterion 1 on the tests the rule tasks wrote, which reproduce the cause (a task with a folder touch, a planner's story folder touch, a story in progress whose claim keeps a folder, an `after` hub) and show it refused or reported; criterion 2 on the closed issue.

## Done when

- I-0152's status is `closed`, its Remediation says what fixed it, and `summary.md` no longer lists it as open.
- Both criteria are ticked through flai, and `flai check --strict` passes scoped to the story.

## Notes

Layer 5, alone.
