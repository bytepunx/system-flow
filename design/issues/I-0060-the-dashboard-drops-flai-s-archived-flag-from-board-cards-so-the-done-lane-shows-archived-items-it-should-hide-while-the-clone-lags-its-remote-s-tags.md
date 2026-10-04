---
id: I-0060
title: The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags
class: defect
status: closed
count: 1
cost: 45m
first_reported: 2026-10-03T06:24:48Z
last_reported: 2026-10-03T06:24:48Z
updated: 2026-10-04T21:49:35Z
---

# I-0060 The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags

## Description
The dashboard drops flai's archived flag from board cards, so the done lane shows archived items it should hide while the clone lags its remote's tags

## Instances

### 2026-10-03T06:24:48Z
Operator saw done, archived items still on the dashboard's done lane and asked why. flai keeps a done, archived item on the done column only while its release looks pending (S-0087, workitem.NewBoardView), and marks it archived on the card; while the clone lags its remote's release tags, pending cannot be told from published, so the board page's doneLane() is meant to hide those cards (S-0174). It never can: flaiover/src/lib/server/board.ts maps flai's board.get cards into its own Card type without the archived field, so every card reaches the browser with archived undefined and the filter in flaiover/src/lib/publish.ts keeps them all. board.svelte.test.ts passes because it mocks /api/board with cards that carry archived, bypassing board.ts. Only stories and epics can be pending, so archived tasks cannot reach the done lane this way; at the time of the examination no served project had an archived card in done, and system-flow's six done cards were the done tasks of the in-progress S-0199, which is by design. Fix: carry archived through board.ts's Card and mapping, and add a test through board() rather than a mocked /api/board.

## Remediation
Closed 2026-10-04T21:49:35Z: S-0247: flaiover's board.ts now carries flai's archived flag onto each board card, so doneLane() leaves out the done, archived cards while the clone lags its remote's tags; repo-channel.test.ts checks it through board() with a fake board.get rather than a mocked /api/board.
