---
id: I-0300
title: A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T08:31:47Z
last_reported: 2026-11-07T08:31:47Z
updated: 2026-11-07T08:31:47Z
---

# I-0300 A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Description
A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Instances

### 2026-11-07T08:31:47Z
Story: S-0564.
S-0564's import run: the board's Import made S-0001 in the new project flai-scratch-import, in ready at 08:25:48Z. The operator turned on agent for the project at 08:26:06Z, as the import formula needs. flai serve started no agent for S-0001: still ready, no journal entry, at 08:32Z, with the in-progress limit free in that project. flai serve starts an agent when a story enters ready, so a story already in ready when agent comes on waits for a manual Start agent. The board's Import cannot have agent on before the project exists, so every import formula run hits this. Enabling agent, or flai serve's periodic look, could start the ready stories that have none.

## Remediation
