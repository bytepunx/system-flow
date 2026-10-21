---
id: TH-0545
title: Merge queue batch 20261021-014447 fails with S-0454 at smoke
anchor:
  path: wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md
  item: S-0454
status: resolved
participants: [flai serve]
created: 2026-10-21T01:52:07Z
updated: 2026-10-21T02:06:08Z
---

# TH-0545 Merge queue batch 20261021-014447 fails with S-0454 at smoke

On wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md.

## Entries

### 2026-10-21T01:52:07Z flai serve
Merge queue batch `20261021-014447` fails with S-0454 at `smoke`. The queue verified S-0454 alone on main at `259e69ae8d0d`, at `6a4ae3b52fd6`, the copy of its verified head `6a4ae3b52fd6` rebased in the batch, and `smoke` failed:

```text
smoke failed: smoke: smoke: template render and check
```

flai moved S-0454 back to in-progress and took it out of the merge queue. Fix it on `story/S-0454`, run flai verify S-0454, move it to review, and enqueue it again with flai queue add S-0454.

### 2026-10-21T02:06:08Z flai serve
Resolved: S-0454 is back in review at 49c5260ffc16, a new commit; it is batched again once it is enqueued
