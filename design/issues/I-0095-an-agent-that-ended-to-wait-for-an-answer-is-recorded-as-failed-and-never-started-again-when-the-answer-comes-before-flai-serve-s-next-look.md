---
id: I-0095
title: An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look
class: defect
status: closed
count: 1
cost: 1m
first_reported: 2026-10-06T21:01:06Z
last_reported: 2026-10-06T21:01:06Z
updated: 2026-10-08T00:17:23Z
---

# I-0095 An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look

## Description

An agent that ends with a question open on a thread is recorded as `asked`, and flai serve starts it again when the thread is answered. flai serve decides this when it notices the agent has ended, at a look, up to a minute later. If the operator answers in between, the look finds an ended agent and no open question. It records the run as `failed`, and nothing starts the agent again: the story stays in progress with no agent, holding every story that waits for it.

Two things made the window wide on 2026-10-06. flai serve had been restarted since it started the agent, so it was not the agent's parent, could not see it exit, and learned of the end only at its next look. And the operator was at the dashboard and answered within thirty seconds of the agent ending.

## Instances

### 2026-10-06T21:01:06Z
Story: S-0223.
S-0223's agent ended at 20:53:33Z to wait for the operator on TH-0196. The operator answered at 20:54:03Z. flai serve, restarted twice since it started that agent (20:36:50Z and 20:48:06Z) and so not its parent, looked at 20:54:03Z, found the run ended with no question open, and recorded it failed: 'ended (an exit code nobody saw) with S-0223 in in-progress'. A failed run is not started again on an answer. Claude, watching the board, ran flai serve agent restart at 20:54:38Z; the cost is those 35 seconds, and has no limit when nobody is watching.

## Remediation

Directions to weigh: when a run ends, read whether its story has a thread whose last answer came after the agent's last entry on it, and treat that as an answer to act on rather than as no question; or record at the moment an agent asks, through `thread_open` or `thread_reply` on its own story, that it is waiting, so that the end is `asked` whatever the thread says by the time of the look. Either way, a story in progress whose agent has ended, with no thread awaiting the operator and no block, could be restarted once by flai serve, which is the second remediation of I-0084.

Story S-0317 remediates this issue, created from it at 2026-10-07T18:59:47Z.
Closed 2026-10-08T00:17:23Z: S-0317: flai serve's judgeRun now records a run asked when its agent wrote on a thread of its story after the run started and that question was answered before the end was judged (askedAnswered in flai/internal/serve/agents.go), so the next look starts the agent again at once in its session instead of recording it failed. TestAnAnswerBeforeTheEndIsJudgedStartsTheAgentAgain reproduces I-0095 and fails without the change.
