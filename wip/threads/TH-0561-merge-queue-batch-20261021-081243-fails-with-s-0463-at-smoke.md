---
id: TH-0561
title: Merge queue batch 20261021-081243 fails with S-0463 at smoke
anchor:
  path: wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md
  item: S-0463
status: answered
participants: [flai serve, agent-S-0463]
created: 2026-10-21T08:19:09Z
updated: 2026-10-21T08:19:36Z
---

# TH-0561 Merge queue batch 20261021-081243 fails with S-0463 at smoke

On wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md.

## Entries

### 2026-10-21T08:19:09Z flai serve
Merge queue batch `20261021-081243` fails with S-0463 at `smoke`. The queue verified S-0463 alone on main at `53cdfda3fda1`, at `e69059a341aa`, the copy of its verified head `630c8d48f0a6` rebased in the batch, and `smoke` failed:

```text
smoke: smoke: template render and check
    building bin/flai
    Created /tmp/tmp.RZ7dZjkgc8/sample from template /home/alex/git/bytepunx/system-flow/.flai-cache/worktrees/merge-queue/template (1.0.106)
      66 files written
    
    Next:
      cd /tmp/tmp.RZ7dZjkgc8/sample
      flai epic new "First deliverable"
      flai dashboard
    0 items checked, 0 errors, 0 warnings
    template renders and passes check
    smoke: repository check
    ../../../wip/kanban/stories/S-0458-the-dashboard-s-pages-are-reached-from-a-left-nav-rail-of-five-collapsible-groups-and-the-top-bar-holds-only-the-name-the-project-picker-the-theme-and-the-status-pills.md:20: warning: wip.overlap: S-0458 touches design/system/flaiover-dashboard.md and docs/users/flaiover.md, which S-0465 (in progress) also touches as design/system/flaiover-dashboard.md and docs/users/flaiover.md
    design/components/flaiover:1: warning: components.examples: the props of the examples under design/components/flaiover are not checked against their components': svelte is not installed in flaiover, and flai reads components' props with its own svelte/compiler: run npm install (or its package manager's install) in flaiover, then run this again
    2474 items checked, 0 errors, 2 warnings (1 that --strict passes over: review over its limit and an epic behind its stories, which only the operator clears, or a narrative not yet written, which flai stream state writes, or component examples left unchecked for want of node, a Svelte sub-project, or a sub-project's svelte)
    exit status 1
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261021-081243`.

flai moved S-0463 back to in-progress and took it out of the merge queue. Fix it on `story/S-0463`, run flai verify S-0463, move it to review, and enqueue it again with flai queue add S-0463.

### 2026-10-21T08:19:36Z agent-S-0463
S-0463's change did not cause this failure. The batch's smoke ran `scripts/check.sh` with no CLOSE_OUT_STORY. Its unscoped `flai check --strict` failed on a `wip.overlap` between S-0458 and S-0465, two other stories in progress. That is I-0174, which I bumped on S-0463's branch (count 3). S-0464, now in review, scopes the batch's checks to its stories and fixes it. S-0463's own close-out, scoped to it, passed smoke. I am closing out again at the new commit and moving S-0463 to review. Until S-0464 lands, any batch's smoke fails the same way whenever two open stories overlap. So S-0463 is best batched after S-0464 or together with it.
