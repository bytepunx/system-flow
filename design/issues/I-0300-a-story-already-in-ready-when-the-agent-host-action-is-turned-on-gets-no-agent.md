---
id: I-0300
title: A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-07T08:31:47Z
last_reported: 2026-11-07T08:31:47Z
updated: 2026-11-07T10:23:33Z
---

# I-0300 A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Description
A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Instances

### 2026-11-07T08:31:47Z
Story: S-0564.
S-0564's import run: the board's Import made S-0001 in the new project flai-scratch-import, in ready at 08:25:48Z. The operator turned on agent for the project at 08:26:06Z, as the import formula needs. flai serve started no agent for S-0001: still ready, no journal entry, at 08:32Z, with the in-progress limit free in that project. flai serve starts an agent when a story enters ready, so a story already in ready when agent comes on waits for a manual Start agent. The board's Import cannot have agent on before the project exists, so every import formula run hits this. Enabling agent, or flai serve's periodic look, could start the ready stories that have none.

## Remediation

Story S-0589 remediates this issue, created from it at 2026-11-07T09:44:04Z.
Closed 2026-11-07T10:23:33Z: S-0589: the cause this issue guessed does not occur. flai serve's launcher reads the agent action at every look, on each change under kanban and threads and every minute, and starts a story already in ready once the action is on (ADR-0041); TestAStoryAlreadyReadyGetsItsAgentOnceTheActionIsOnAndSomethingStartsIt in flai/internal/serve/agents_test.go (T-2573) pins it. S-0001 in flai-scratch-import got no agent because its project's default agent named no harness and no command was set on the host: serve.log logged "S-0001 names no harness, and no command is set on the host" at 08:26:49Z, 43 seconds after the action was enabled, and S-0001 started at 08:42:14Z once the default agent was given a harness. S-0587 makes a default harness work out of the box. The operators guide now says a story already in ready needs no Start agent (T-2575).
