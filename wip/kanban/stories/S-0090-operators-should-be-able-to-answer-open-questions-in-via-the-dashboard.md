---
id: S-0090
type: story
nature: improvement
title: Operators should be able to answer open questions in via the dashboard
status: review
parent: E-0003
owner: alex
created: 2026-09-21T04:03:18Z
updated: 2026-09-23T00:28:48Z
transitions:
  - to: ready
    at: 2026-09-22T22:39:01Z
    by: alex
  - to: in-progress
    at: 2026-09-23T00:28:45Z
    by: system-flow
  - to: review
    at: 2026-09-23T00:28:48Z
    by: system-flow
tags: [dashboard]
touches: [flai/cmd, flai/internal/hostapi, flai/internal/workitem, flaiover/src]
---
# S-0090 Operators should be able to answer open questions in via the dashboard

## Goal

When an agent has a question, it should place that question in the inbox where the operator can answer it, have it recorded, and then the open question should close out.

## Acceptance criteria
- [x] Operators need to be able to answer open questions in the dashboard
- [x] Open questions should block a story from being moved to review since the agent can't complete the story without the answer
- [x] Answered open questions should not remain in the inbox

## Tasks
- T-0334 flai stream answer moves a hand-written open question to Decisions; the dashboard's inbox offers it inline

## Notes
Criterion 2 duplicates S-0089's own second criterion exactly (the two stories were filed four minutes apart); the rule itself (`internal/workitem/rules.go`'s `Move`) came with S-0089. Since S-0089's acceptance this branch is rebased onto it, and the two halves are joined here: the refusal now says how to answer (the dashboard's inbox, or `flai stream answer`), `TestAnsweringAnOpenQuestionUnblocksReview` holds that an unanswered question refuses review and answering it lets the story through, and `design/system/agent-narrative.md` records the rule. This story stands on its own for what it actually adds: the answer capability itself (criterion 1), which criterion 3 follows from directly, since flai's `inbox.designer` re-reads the narrative fresh on every call, so a bullet `stream answer` removes is gone from the very next inbox read — no caching to invalidate, verified by a real narrative-file diff before/after in both `internal/workitem`'s and flaiover's own tests.
