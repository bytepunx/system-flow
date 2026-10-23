---
id: I-0152
title: Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five
class: efficiency
status: closed
count: 1
cost: 44h
first_reported: 2026-10-15T22:48:43Z
last_reported: 2026-10-15T22:48:43Z
updated: 2026-10-23T04:58:49Z
usage:
  strategic:
    - kind: analyzer
      seconds: 113
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 97
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
---

# I-0152 Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five

## Description
Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five

## Instances

### 2026-10-15T22:48:43Z
Report: design/analysis/2026-10-15-bottlenecks.md.
E-0016's planned stories carried folder-wide touches and after chains on S-0218 and S-0223, so one story in progress or one hub awaiting acceptance held the whole ready column; I-0087 (closed) recorded the earlier form through docs/users/flai.md. Shares under ADR-0134 cut it in W42.

## Impact
flai stats claims.weeks held_seconds: W40 overlap 653,549 + after 582,160; W41 overlap 690,335 + after 887,656 (438 h); W42 to date overlap 25,902 + after 133,264. claims.days held 9 on 2026-10-03, 9 on 10-05, 10 on 10-06 with in_progress 0, 1, 2 of limit 5. items held_seconds: S-0138 180,695, S-0137 172,010, S-0136 171,537, S-0135 170,672, S-0227 141,332, S-0229 141,310, S-0224 139,613. Archived S-0207..S-0211 and S-0220..S-0229 declare flai/cmd, flai/internal/*, flaiover/src/routes, docs/users; S-0218 and S-0223 are after hubs. exact_share 0.49 to 0.60. W42 shares.started 11, shares.held_seconds 179,046 avoided.

- time_lost_per_cycle: 44h

## Remediation

Found by the analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md).

Story S-0468 remediates this issue, created from it at 2026-10-21T11:58:08Z.
Closed 2026-10-23T04:58:49Z: Fixed by S-0468 (ADR-0188). flai refuses a task a folder touch through flai task new, flai edit --touches, flai touches, item_new, and item_edit, so ADR-0096's narrowing always has files to narrow to. It refuses the planner's folder touch on a story unless the story's ### Planning note names the folder. flai check gives the advisory warnings touches.folder, on a story in progress whose claim keeps a folder, and wip.after-hub, on a story three or more open stories name in after. The planner's convention, planPrompt, and planner.md say that after names only a story whose accepted change is needed on main, never an order of preference.
