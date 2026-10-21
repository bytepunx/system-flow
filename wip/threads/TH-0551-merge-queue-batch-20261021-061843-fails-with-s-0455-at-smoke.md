---
id: TH-0551
title: Merge queue batch 20261021-061843 fails with S-0455 at smoke
anchor:
  path: wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md
  item: S-0455
status: resolved
participants: [flai serve]
created: 2026-10-21T06:24:57Z
updated: 2026-10-21T06:40:11Z
---

# TH-0551 Merge queue batch 20261021-061843 fails with S-0455 at smoke

On wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md.

## Entries

### 2026-10-21T06:24:57Z flai serve
Merge queue batch `20261021-061843` fails with S-0455 at `smoke`. The queue verified S-0455 alone on main at `f7d46658b74d`, at `b1d7174e3475`, the copy of its verified head `b1d7174e3475` rebased in the batch, and `smoke` failed:

```text
smoke failed: smoke: building bin/flai
```

flai moved S-0455 back to in-progress and took it out of the merge queue. Fix it on `story/S-0455`, run flai verify S-0455, move it to review, and enqueue it again with flai queue add S-0455.

### 2026-10-21T06:40:11Z flai serve
Resolved: S-0455 is back in review at 42a178cdea36, a new commit; it is batched again once it is enqueued
