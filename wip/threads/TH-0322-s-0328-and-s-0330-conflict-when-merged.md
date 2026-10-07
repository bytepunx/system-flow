---
id: TH-0322
title: S-0328 and S-0330 conflict when merged
anchor:
  path: wip/kanban/stories/S-0330-flai-message-sends-a-message-from-one-open-story-s-agent-to-another-s-kept-apart-from-the-operator-s-threads.md
  item: S-0330
status: answered
participants: [flai, agent-S-0330]
created: 2026-10-07T20:26:04Z
updated: 2026-10-07T20:26:10Z
---

# TH-0322 S-0328 and S-0330 conflict when merged

On wip/kanban/stories/S-0330-flai-message-sends-a-message-from-one-open-story-s-agent-to-another-s-kept-apart-from-the-operator-s-threads.md.

## Entries

### 2026-10-07T20:26:04Z flai
A trial merge of story/S-0328 with story/S-0330 at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0328 and S-0330 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T20:26:10Z agent-S-0330
S-0330 takes the second place. The only conflict is `design/adrs/README.md`: S-0328 adds the ADR-0119 row and S-0330 adds the ADR-0120 row below it. Neither story narrows anything. If S-0328 is accepted first, S-0330's next `flai stream sync` keeps both rows.
