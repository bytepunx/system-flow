---
id: TH-0533
title: "S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3)"
anchor:
  path: wip/kanban/stories/S-0412-a-story-is-worked-to-review-end-to-end-on-the-opencode-adapter-by-a-scripted.md
  item: S-0412
status: resolved
participants: [orchestrator, alex]
created: 2026-10-20T21:14:19Z
updated: 2026-11-07T09:30:49Z
---

# TH-0533 S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3)

On wip/kanban/stories/S-0412-a-story-is-worked-to-review-end-to-end-on-the-opencode-adapter-by-a-scripted.md.

## Entries

### 2026-10-20T21:14:19Z orchestrator
**Recommended:** as planner-E-0023 proposed on TH-0467, question 3:
- set S-0412's `agent.model` to the non-Claude model you want measured;
- put an OpenRouter key with a 5 USD limit in the host's `providers` map;
- then finalize S-0412 yourself.

## Why I have not finalized it

`flai promote --drafts` lists S-0412 as complete: criteria, touches, a forecast of 30m, and a value from E-0023's input. But the plan leaves its model "chosen by you when you finalize that story", and its run needs a funded key. Both are spend and your choice, so I leave them to you.

There is no hurry: S-0412 waits behind S-0406 to S-0411 by `after`.

### 2026-11-07T09:28:49Z alex
We will only be using OpenRouter fixtures, no direct integration will be part of this repository for now

### 2026-11-07T09:30:49Z orchestrator
Resolved: alex answered "OpenRouter fixtures only, no direct integration in this repository for now" (D-0507). S-0412 as written needs a real run of a non-Claude model through a funded OpenRouter key, so its goal and criteria 1 to 3 cannot be met. planner-E-0023 has been running since 2026-11-07T09:29:43Z (pid 955118) to revisit S-0406 and S-0412 on fixtures (TH-0532, D-0506). S-0412 stays a draft. I finalize it under finalize_drafts once flai promote --drafts lists the re-planned story as complete and consistent.
