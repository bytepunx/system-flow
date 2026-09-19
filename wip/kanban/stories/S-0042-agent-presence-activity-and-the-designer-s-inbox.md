---
id: S-0042
type: story
nature: feature
title: Agent presence, activity, and the designer's inbox
status: ready
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-19T01:56:17Z
transitions:
  - to: ready
    at: 2026-09-19T01:56:17Z
    by: alex
tags: [dashboard, cli]
---

# S-0042 Agent presence, activity, and the designer's inbox

## Goal
The designer sees which agents are active and what they are doing, and has one inbox of things that need a human: open questions, threads awaiting an answer, stories in review, blocked items, and overlapping touches.

## Acceptance criteria
- [ ] An activity view built from `wip/agents/index.md` and the narratives: active streams, agent and session, last log entry and its age, current task, blocked flags
- [ ] An inbox view listing open threads and questions, stories in review, blocked items, and `wip.overlap` warnings, each linking to its page; a badge in the navigation with the count
- [ ] Optional desktop notifications and a webhook (`dashboard.notify_url`) for new inbox entries, off by default
- [ ] Tests on fixtures; docs/users updated

## Tasks

## Notes
Presence is derived from files, so nothing new is written; an agent that stops logging simply ages out of the view.
