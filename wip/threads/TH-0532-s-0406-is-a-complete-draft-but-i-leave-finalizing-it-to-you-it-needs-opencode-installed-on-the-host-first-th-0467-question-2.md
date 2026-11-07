---
id: TH-0532
title: "S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2)"
anchor:
  path: wip/kanban/stories/S-0406-a-story-s-agent-may-run-on-opencode-headless-with-flai-s-mcp-server-no-guard.md
  item: S-0406
status: resolved
participants: [orchestrator, alex]
created: 2026-10-20T21:14:16Z
updated: 2026-11-07T09:29:48Z
---

# TH-0532 S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2)

On wip/kanban/stories/S-0406-a-story-s-agent-may-run-on-opencode-headless-with-flai-s-mcp-server-no-guard.md.

## Entries

### 2026-10-20T21:14:16Z orchestrator
**Recommended:** install OpenCode on this host, then finalize S-0406 yourself. Alternatively, answer TH-0467 question 2 with "fixtures only" and tell me here, and I will finalize it.

## Why I have not finalized it

`flai promote --drafts` lists S-0406 as complete: criteria, touches, a forecast of 1h15m, and a value from E-0023's input. But its first criterion checks the four open OpenCode findings against a real `opencode` binary, and flai never installs one (E-0023's Notes).

- planner-E-0023 asked on TH-0467, question 2, whether you install OpenCode before this story starts, or whether it runs on recorded fixtures. That question is still unanswered.
- The planner noted that questions 2 and 3 "can wait until S-0406 and S-0412 are finalized". Finalizing S-0406 now would let an agent start it with neither.

## What I did

I finalized S-0407 to S-0411. Each is complete, consistent, and held behind S-0406 by `after`, so none can start before you decide here.

### 2026-11-07T09:27:45Z alex
fixtures only for now - I have no intention of configuring OpenCode at this time

### 2026-11-07T09:29:48Z orchestrator
Resolved: alex answered "fixtures only" (D-0506), the planner's alternative on TH-0467 question 2. That makes S-0406's title, goal, criterion 6, and T-1682, which settle the four findings against a real opencode binary, unmeetable as written. The planner noted on TH-0467 that a different answer means re-planning. I started planner-E-0023 at 2026-11-07T09:29:43Z to revisit S-0406 and S-0412 (TH-0533) on fixtures. S-0406 stays a draft. I finalize it under finalize_drafts once flai promote --drafts lists the re-planned story as complete and consistent.
