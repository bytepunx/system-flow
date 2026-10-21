---
id: TH-0532
title: "S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2)"
anchor:
  path: wip/kanban/stories/S-0406-a-story-s-agent-may-run-on-opencode-headless-with-flai-s-mcp-server-no-guard-and-no-roles-and-the-findings-the-opencode-section-leaves-to-check-are-settled-against-a-real-opencode-binary.md
  item: S-0406
status: open
participants: [orchestrator]
created: 2026-10-20T21:14:16Z
updated: 2026-10-20T21:14:16Z
---

# TH-0532 S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2)

On wip/kanban/stories/S-0406-a-story-s-agent-may-run-on-opencode-headless-with-flai-s-mcp-server-no-guard-and-no-roles-and-the-findings-the-opencode-section-leaves-to-check-are-settled-against-a-real-opencode-binary.md.

## Entries

### 2026-10-20T21:14:16Z orchestrator
**Recommended:** install OpenCode on this host, then finalize S-0406 yourself. Alternatively, answer TH-0467 question 2 with "fixtures only" and tell me here, and I will finalize it.

## Why I have not finalized it

`flai promote --drafts` lists S-0406 as complete: criteria, touches, a forecast of 1h15m, and a value from E-0023's input. But its first criterion checks the four open OpenCode findings against a real `opencode` binary, and flai never installs one (E-0023's Notes).

- planner-E-0023 asked on TH-0467, question 2, whether you install OpenCode before this story starts, or whether it runs on recorded fixtures. That question is still unanswered.
- The planner noted that questions 2 and 3 "can wait until S-0406 and S-0412 are finalized". Finalizing S-0406 now would let an agent start it with neither.

## What I did

I finalized S-0407 to S-0411. Each is complete, consistent, and held behind S-0406 by `after`, so none can start before you decide here.
