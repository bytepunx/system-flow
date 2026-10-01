---
id: TH-0042
title: "S-0175's measurement criterion: which two delegating runs count?"
anchor:
  path: wip/kanban/stories/S-0175-agents-the-claude-code-adapter-starts-hand-noisy-work-and-verification-to-sub-agents-that-cannot-act-for-the-story.md
  item: S-0175
status: resolved
participants: [agent-S-0175, alex]
created: 2026-10-01T07:53:18Z
updated: 2026-10-01T08:27:52Z
---

# TH-0042 S-0175's measurement criterion: which two delegating runs count?

On wip/kanban/stories/S-0175-agents-the-claude-code-adapter-starts-hand-noisy-work-and-verification-to-sub-agents-that-cannot-act-for-the-story.md.

## Entries

### 2026-10-01T07:53:18Z agent-S-0175
The criterion asks for cost, cache reads, and turns on at least two stories that delegated, against comparable runs that did not. Stories only get the new prompt once a flai with it is released and installed on the host, which is after S-0175 is accepted. flai serve keeps 8 logs today; only S-0118's shows a sub-agent (one Explore).

Recommended: measure (1) this run of S-0175, which delegates search, test runs, and pre-review verification to sub-agents by hand as the new prompt says, and (2) S-0118, against comparable non-delegating runs from the kept logs and the usage in archived items' front matter; record that table in agent-context.md, and file a follow-up story to measure two stories that run with the released prompt.

Alternative: leave the criterion unchecked, move S-0175 to review with it explained in the notes, and measure after release in the same story.

Which do you want?

### 2026-10-01T08:05:17Z alex
take the recommended approach

### 2026-10-01T08:06:59Z agent-S-0175
Recorded in the narrative's Decisions. T-0624 measures this run and S-0118 into agent-context.md, and I will file the follow-up story to measure two runs with the released prompt before S-0175 goes to review.

### 2026-10-01T08:20:23Z alex
Resolved.

### 2026-10-01T08:27:52Z agent-S-0175
Resolved: This run and S-0118 are measured in design/system/agent-context.md § Sub-agents › Measured; S-0188 measures two stories run with the released prompt.
