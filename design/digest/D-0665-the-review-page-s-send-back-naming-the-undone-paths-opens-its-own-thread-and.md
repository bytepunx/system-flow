---
id: D-0665
title: "The review page's Send back, naming the undone paths opens its own thread and then moves the story, so since S-0823 the story gets a second, by-hand thread with the reason too, the one flai serve starts the agent on; left as is, since the agent reads both"
kind: decision
sources: [S-0823, T-3545]
topics: [flai, flaiover, code]
paths: [flai/cmd/move.go, flai/internal/sendback/sendback.go, docs/users/flaiover.md]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0665 The review page's Send back, naming the undone paths opens its own thread and then moves the story, so since S-0823 the story gets a second, by-hand thread with the reason too, the one flai serve starts the agent on; left as is, since the agent reads both

S-0823 made `flai move` along the send-back transition, by anyone but `flai` and the story's own agent, open a thread on the story with the move's reason (I-0397). The dashboard's review page already had a button, **Send back, naming the undone paths** (`docs/users/flaiover.md`), which opens a thread naming the paths and then moves the story through the host write `item.move`. That move now opens the by-hand thread as well. From S-0823's Decisions (`wip/archive/agents/S-0823.md`):

> T-3545: the review page's Send back, naming the undone paths opens its own thread and then moves the story through item.move, so flai move now opens a second, by-hand thread with the reason too; that is the one flai serve starts the agent on. Left as is: the agent reads both.

So two threads on a story sent back from that button are expected, not a defect: S-0823 chose to leave them. A story that wants one thread there changes the button or `flai move`'s by-hand send back in `flai/cmd/move.go` and `flai/internal/sendback/sendback.go`.
