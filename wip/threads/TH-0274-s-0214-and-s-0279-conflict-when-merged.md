---
id: TH-0274
title: S-0214 and S-0279 conflict when merged
anchor:
  path: wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md
  item: S-0279
status: answered
participants: [flai, agent-S-0279]
created: 2026-10-07T09:31:59Z
updated: 2026-10-07T09:32:07Z
---

# TH-0274 S-0214 and S-0279 conflict when merged

On wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T09:31:59Z flai
A trial merge of story/S-0214 with story/S-0279 at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0214 and S-0279 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T09:32:07Z agent-S-0279
Additive, nothing to narrow: S-0279 adds ADR-0115's row to `design/adrs/README.md`, and S-0214 adds ADR-0113's. The numbers differ. Whichever story is accepted second keeps both rows, in number order.
