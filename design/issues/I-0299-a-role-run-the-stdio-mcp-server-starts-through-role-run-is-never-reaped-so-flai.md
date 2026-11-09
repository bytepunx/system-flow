---
id: I-0299
title: A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it
class: defect
status: closed
count: 1
cost: 20m
first_reported: 2026-11-07T08:21:58Z
last_reported: 2026-11-07T08:21:58Z
updated: 2026-11-09T22:13:18Z
---

# I-0299 A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it

## Description
A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it

## Instances

### 2026-11-07T08:21:58Z
Story: S-0584.
S-0584's formula step T-2545 started the analyzer with role_run from the story agent's stdio `flai mcp --agent agent-S-0584` (pid 3871368) at 2026-11-07T08:04:52Z, pid 3881796. The analyzer ended its work at about 08:18Z (log /home/alex/.flai/serve/agents/sf-analyzer-20261107T080452Z.log ends with its result: report design/analysis/2026-11-07-all.md, I-0284 to I-0298 filed). The process then stayed a zombie (state Z) of the flai mcp process. flai serve's run record in agents.json was not settled, and no agent event reached the instance's agent. Cause, read from the code: the hand-over launcher's await (flai/internal/serve/agents.go, handOver branch) closes the output and calls cmd.Process.Release without cmd.Wait, so a long-lived parent never reaps the child. Also, Alive/Owns (flai/internal/serve/proc_unix.go, proc_linux.go) read the session and start time from /proc/<pid>/stat but not the state, so settleOrphans (agents.go) counts a zombie as running. flai analyze exits after starting the run, so its child is reaped. Only a hand-over from a long-lived parent, flai mcp over stdio, leaves the zombie. Possible fixes: reap in a goroutine when the hand-over parent stays alive, or have Owns treat state Z as not running.

## Remediation

Story S-0632 remediates this issue, created from it at 2026-11-07T23:40:59Z.
Closed 2026-11-09T22:13:18Z: S-0632: a launcher that hands its agent over now reaps it in a goroutine instead of releasing it (flai/internal/serve/agents.go, await), so flai mcp over stdio keeps no zombie of a role run; and Alive and Owns count a process in state Z or X as gone on Linux (proc_linux.go, proc_unix.go), so flai serve settles a zombie of any parent. TestAHandedOverAgentIsReaped and TestAZombieIsNeitherAliveNorOwned reproduce both halves.
