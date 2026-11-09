---
id: I-0318
title: flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-11-09T09:06:05Z
last_reported: 2026-11-09T09:06:05Z
updated: 2026-11-09T09:06:05Z
---

# I-0318 flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Description
flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Instances

### 2026-11-09T09:06:05Z
Story: S-0613.
S-0613 found that 11 of I-0230's orchestrator refusals were `flai serve status` with a redirection, but the redirection was not their cause: the guard refuses the orchestrator `flai serve status` without one too, because `cliReads` in flai/internal/guard/guard.go has no `serve` entry, while the host-configuration check (hostconfig_test.go) already counts `flai serve status` among the reads. S-0613 fixed the redirection and left this out of scope: allowing it changes what every strategic agent and read-only sub-agent may run, which is the operator's choice. TestARedirectionNeitherAllowsNorRefusesTheOrchestratorsCommand pins that it is refused alike with and without a redirection.

## Remediation
