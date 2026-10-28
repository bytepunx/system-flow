---
id: D-0096
title: "Should wip/ warnings outside a story stop its close-out: Go with the recommendation, notes are preferable when its a finding outside the story"
kind: answer
sources: [TH-0056]
topics: []
paths: []
date: 2026-10-01
status: absorbed
absorbed_by: "design/adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md#decision"
asked_again: []
cited: []
---

# D-0096 Should wip/ warnings outside a story stop its close-out: Go with the recommendation, notes are preferable when its a finding outside the story

The operator, alex, answered TH-0056 at 2026-10-01T10:57:26Z.

What was asked:

```text
From S-0187. `scripts/close-out.sh` runs `flai check --strict`, which reads the main checkout's `wip/` and fails on warnings the story's agent may not fix: `threads.archived` (TH-0032 today, which holds an open request of yours), `item.archive` for done epics, another story's narrative lint. The new rule in `work-management.md` says to move to review only when the close-out ends clean, so the agent either waits on you or finishes the steps by hand. S-0187 went to review with only TH-0032's warning left, and its narrative says why.
```

The answer:

```text
Go with the recommendation, notes are preferable when its a finding outside the story. Have flai record the finding as an issue (if one doesn't exist) or increment the count so operators can see where these are occurring.
```

Source: TH-0056.
