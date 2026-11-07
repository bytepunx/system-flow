---
id: I-0245
title: flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-11-02T00:53:54Z
last_reported: 2026-11-02T00:53:54Z
updated: 2026-11-07T08:37:07Z
---

# I-0245 flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Description
flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Instances

### 2026-11-02T00:53:54Z
Story: S-0553.
S-0553's task sub-agent for T-2351 could not run `scripts/flai.sh role definition planner` to check the command it wrote: flai guard refused it as a flai command that changes work items (ADR-0060). The guard's `cliReads` table in flai/internal/guard/guard.go lists no `role` subcommand, so `flai role list`, `flai role show`, and `flai role definition` without `--write`, which only read, are refused to every sub-agent; the story's agent ran it instead.

## Remediation

Story S-0568 remediates this issue, created from it at 2026-11-02T20:31:49Z.
Closed 2026-11-07T08:37:07Z: Fixed by S-0568: flai guard's flagReads now has a role entry, --write in without and a form (readsRoles) that allows role alone, list, show, and definition, so every sub-agent and every strategic agent's own call reads the roles, while role run and role definition --write stay refused. TestEveryAgentReadsTheRolesButNoneRunsOneOrWritesItsDefinition in flai/internal/guard/roles_test.go and the a7 and a8 TestGuard cases in flai/cmd/guard_test.go reproduce the refusal; the help, the reference, the user guide, and the design say so.
