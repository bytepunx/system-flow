---
id: D-0641
title: A conversation between two stories closes once either is in an accepted or cancelled state, a story talks from an active or waiting state, and a share's held story is in a state a pull transition leaves and its holder in an active state
kind: decision
sources: [S-0749, T-3233]
topics: [flai, go, code]
paths: [flai/internal/messages/messages.go, flai/internal/messages/share.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0641 A conversation between two stories closes once either is in an accepted or cancelled state, a story talks from an active or waiting state, and a share's held story is in a state a pull transition leaves and its holder in an active state

S-0749 moved messages between stories (ADR-0120) off state names onto the workflow's kinds (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0749.md`):

> Messages (T-3233): a conversation is closed once either story is in an accepted or cancelled state; a story talks from an active or waiting state; a held story is in a state a pull transition leaves, and a holder is in an active state. The state words in refusals come from the workflow.

On the baseline that is `done` or `cancelled`, `in-progress` or `review`, `ready`, and `in-progress`, as before. A refusal now names the states of the board it runs on. `flai/internal/workitem/hold.go` still names `InProgress` for a share's holder (I-0379), so the hold itself keeps the baseline until a story moves it.

Source: S-0749's narrative, `wip/archive/agents/S-0749.md`, § Decisions; I-0379.
