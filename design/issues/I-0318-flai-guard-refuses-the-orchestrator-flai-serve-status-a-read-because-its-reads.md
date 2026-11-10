---
id: I-0318
title: flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status
class: efficiency
status: closed
count: 1
cost: 2m
first_reported: 2026-11-09T09:06:05Z
last_reported: 2026-11-09T09:06:05Z
updated: 2026-11-10T09:02:58Z
---

# I-0318 flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Description
flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Instances

### 2026-11-09T09:06:05Z
Story: S-0613.
S-0613 found that 11 of I-0230's orchestrator refusals were `flai serve status` with a redirection, but the redirection was not their cause: the guard refuses the orchestrator `flai serve status` without one too, because `cliReads` in flai/internal/guard/guard.go has no `serve` entry, while the host-configuration check (hostconfig_test.go) already counts `flai serve status` among the reads. S-0613 fixed the redirection and left this out of scope: allowing it changes what every strategic agent and read-only sub-agent may run, which is the operator's choice. TestARedirectionNeitherAllowsNorRefusesTheOrchestratorsCommand pins that it is refused alike with and without a redirection.

## Remediation

Story S-0680 remediates this issue, created from it at 2026-11-10T08:34:57Z.
Closed 2026-11-10T09:02:58Z: S-0680 fixed it: flagReads in flai/internal/guard/guard_tables.go has a serve entry whose readsServe passes the read forms ADR-0220 leaves open (status, actions, journal, agent bare or show, agent harness and provider bare or with a name alone, checks, components, and notify show, forge with no setting, mcp bare or show, import list) for every agent held to its reads, and never a form writesHost finds. TestEveryAgentReadsTheReadFormsOfFlaiServe pins it for the orchestrator, the planner, the analyzer, the curator, and sub-agents, and TestARedirectionNeitherAllowsNorRefusesTheOrchestratorsCommand now expects flai serve status to pass with and without a redirection.
