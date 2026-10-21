---
id: TH-0560
title: Merge queue batch 20261021-080642 fails with S-0456 at smoke
anchor:
  path: wip/kanban/stories/S-0456-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md
  item: S-0456
status: answered
participants: [flai serve, agent-S-0456]
created: 2026-10-21T08:12:42Z
updated: 2026-10-21T08:13:12Z
---

# TH-0560 Merge queue batch 20261021-080642 fails with S-0456 at smoke

On wip/kanban/stories/S-0456-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md.

## Entries

### 2026-10-21T08:12:42Z flai serve
Merge queue batch `20261021-080642` fails with S-0456 at `smoke`. The queue verified S-0456 alone on main at `53cdfda3fda1`, at `9db1029f06e4`, the copy of its verified head `097b48d89db6` rebased in the batch, and `smoke` failed:

```text
smoke: smoke: template render and check
    building bin/flai
    Created /tmp/tmp.eX8gjwIqx1/sample from template /home/alex/git/bytepunx/system-flow/.flai-cache/worktrees/merge-queue/template (1.0.106)
      66 files written
    
    Next:
      cd /tmp/tmp.eX8gjwIqx1/sample
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

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261021-080642`.

flai moved S-0456 back to in-progress and took it out of the merge queue. Fix it on `story/S-0456`, run flai verify S-0456, move it to review, and enqueue it again with flai queue add S-0456.

### 2026-10-21T08:13:12Z agent-S-0456
The batch failed on nothing S-0456 changed. Smoke's repository check is an unscoped `flai check --strict`, and it stopped on a `wip.overlap` between S-0458 and S-0465, two other stories in progress. That is I-0174, which S-0464 (in review) fixes by scoping the batch's check to its stories. The batch's `smoke.log` ends at that check, so the gateway check S-0456 changes never ran.

I bumped I-0174 with this occurrence on story/S-0456, which gives the branch a new commit. I am closing out again and moving S-0456 back to review. A batch that runs before S-0464 lands, while S-0458 and S-0465 still overlap, will fail the same way. Landing S-0464 first, or batching S-0456 with it, avoids that.
