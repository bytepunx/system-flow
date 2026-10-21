---
id: TH-0533
title: "S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3)"
anchor:
  path: wip/kanban/stories/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md
  item: S-0412
status: open
participants: [orchestrator]
created: 2026-10-20T21:14:19Z
updated: 2026-10-20T21:14:19Z
---

# TH-0533 S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3)

On wip/kanban/stories/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md.

## Entries

### 2026-10-20T21:14:19Z orchestrator
**Recommended:** as planner-E-0023 proposed on TH-0467, question 3:
- set S-0412's `agent.model` to the non-Claude model you want measured;
- put an OpenRouter key with a 5 USD limit in the host's `providers` map;
- then finalize S-0412 yourself.

## Why I have not finalized it

`flai promote --drafts` lists S-0412 as complete: criteria, touches, a forecast of 30m, and a value from E-0023's input. But the plan leaves its model "chosen by you when you finalize that story", and its run needs a funded key. Both are spend and your choice, so I leave them to you.

There is no hurry: S-0412 waits behind S-0406 to S-0411 by `after`.
