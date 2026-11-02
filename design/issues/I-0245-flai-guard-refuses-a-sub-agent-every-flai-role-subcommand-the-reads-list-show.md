---
id: I-0245
title: flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T00:53:54Z
last_reported: 2026-11-02T00:53:54Z
updated: 2026-11-02T00:53:54Z
---

# I-0245 flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Description
flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Instances

### 2026-11-02T00:53:54Z
Story: S-0553.
S-0553's task sub-agent for T-2351 could not run `scripts/flai.sh role definition planner` to check the command it wrote: flai guard refused it as a flai command that changes work items (ADR-0060). The guard's `cliReads` table in flai/internal/guard/guard.go lists no `role` subcommand, so `flai role list`, `flai role show`, and `flai role definition` without `--write`, which only read, are refused to every sub-agent; the story's agent ran it instead.

## Remediation
