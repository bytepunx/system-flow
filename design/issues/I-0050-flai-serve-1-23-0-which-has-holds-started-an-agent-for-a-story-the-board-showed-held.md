---
id: I-0050
title: flai serve 1.23.0, which has holds, started an agent for a story the board showed held
class: defect
status: closed
count: 2
cost: 8m
first_reported: 2026-09-29T02:06:24Z
last_reported: 2026-09-29T02:13:28Z
updated: 2026-10-01T09:53:41Z
---

# I-0050 flai serve 1.23.0, which has holds, started an agent for a story the board showed held

## Description
flai serve 1.23.0, which has holds, started an agent for a story the board showed held

## Instances

### 2026-09-29T02:06:24Z
2026-09-29: S-0144 entered ready at 01:10:17Z, claiming flai/cmd, while S-0138 (in progress since 01:07:54Z) claimed flai/cmd/prime.go. flai board and inbox showed S-0144 held (overlap), yet the host's flai serve 1.23.0 started agent-S-0144 for it with the instruction to work it to review. Unlike I-0049 the host flai postdates holds, so the agent start path appears not to consult them; not investigated. The agent narrowed the claim to what it changes and worked it; the remaining overlap (docs/users, design/system/flai-cli.md) fails the smoke tier's strict repository check until S-0138 is accepted.

### 2026-09-29T02:13:28Z
2026-09-29: S-0145 entered ready at 01:57:54Z claiming flai/cmd, while S-0138 (in progress) claimed flai/cmd/prime.go. flai board and inbox showed S-0145 held (overlap), yet the host's flai serve 1.23.0 (pid 2077672, running since 01:06:55Z) started agent-S-0145 for it at 02:10:59Z with the instruction to work it to review. Second instance in one evening, same serve. The story is research whose output is an ADR draft and touches no code, so the agent narrowed the claim to design/adrs and design/issues, which cleared the hold, and pulled it.

## Remediation

Cause (S-0182, TH-0049): the operator's Start agent, not the launcher. The other host's journal was not reachable, so this comes from flai 1.23.0's code and the stories' timing. The launcher (`look`) checked `holds.Of` before every start, against the same `Repo.Holds` the board and `inbox` read, and a file and its folder overlap both ways. The story page's Start agent (`agent.start`, `flai serve agent start`, S-0115) starts a held story on purpose (S-0129), but it gave the agent the launcher's prompt word for word ("started by flai serve ... because it entered ready") and journalled a `serve.agent` entry by `flai serve`. Its warning, `agent started for a held story`, went only to the log of the `flai serve agent start` process. Both agents therefore reported a start that flai serve had not made. S-0145's start, at 02:10:59Z, came 3.5 minutes after S-0144 was accepted while S-0138 still held it; the launcher looks at every kanban change and every minute, so it would have started S-0145 at 02:07:25Z or not at all. The designer resolved TH-0049 without the other host's journal.

Fixed in S-0182:

- An agent the operator starts past a hold or a full in-progress limit is told so in its prompt, with the hold's reason, and asked to keep its touches to what it changes rather than narrow them only to clear the hold. The journal's `serve.agent` entry says `on the operator's word, past …`.
- `resume`, the one launcher path that did go past holds, restarts an answered agent only while its story is in progress or in review. A story sent back to ready waits for its hold and the limit, and one in backlog is not started.
- Regression tests: `TestTheLauncherHoldsAFileAndItsFolderBothWays`, `TestAnAnsweredAgentWhoseStoryIsBackInReadyWaitsForItsHold`, `TestAnAnsweredAgentWhoseStoryIsInBacklogIsNotStarted` (`flai/internal/serve`), and `TestAnAgentTheOperatorStartedIsToldWhatItWentPast` (`flai/internal/harness`).
Closed 2026-10-01T09:53:41Z: S-0182: the operator's Start agent started both held stories with the launcher's prompt; the prompt and the journal now say what an operator's start went past, and resume no longer passes a hold
