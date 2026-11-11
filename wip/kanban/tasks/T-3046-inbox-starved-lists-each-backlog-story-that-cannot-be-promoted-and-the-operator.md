---
id: T-3046
type: task
nature: improvement
title: inbox.Starved lists each backlog story that cannot be promoted and the operator's decision it waits on
status: in-progress
parent: S-0691
owner: alex
created: 2026-11-10T10:34:50Z
updated: 2026-11-11T05:11:22Z
transitions:
  - to: ready
    at: 2026-11-11T05:11:22Z
    by: agent-S-0691
  - to: in-progress
    at: 2026-11-11T05:11:22Z
    by: agent-S-0691
stream: S-0691
tags: [flai, go]
touches: [flai/internal/inbox/starved.go, flai/internal/inbox/starved_test.go]
after: [T-3043]
---
# T-3046 inbox.Starved lists each backlog story that cannot be promoted and the operator's decision it waits on

## Work

Add `flai/internal/inbox/starved.go`. It takes the items, the holds, and the open threads, and gives one entry per backlog story that `workitem.PromotionCandidates` lists among `others`, in the pull order. Each entry carries the story, its reasons as the candidates give them, and the decision it waits on:

- A draft: finalize it, with the ID and title of each open thread on it or its tasks that `inbox.Summarize` marks `awaiting: "operator"`.
- A story held by `after`: the story at the root of its chain, followed through each named story not done, with that root's own entry.
- A blocked story: its block reason, and the thread it names, when it names one.
- Any other reason, such as `no forecast duration`: the reason as it is, which a planner run can settle.

The package is `inbox` because `threads` imports `workitem`, and `inbox` imports both. Give the entries a stable text form, so that a caller can tell when the list changed.

It waits on T-3043, whose ADR fixes what the list holds.

## Done when

- `starved_test.go` covers the case from I-0288. Two drafts each have a thread awaiting the operator, and five stories are held by `after` behind one of them. The list names both drafts with their threads, and gives each held story the draft at the root of its chain.
- It also covers a blocked story, a draft with no thread, and a resolved thread, which is left out.
- `flai test flai/internal/inbox/starved.go flai/internal/inbox/starved_test.go` passes.

## Notes
