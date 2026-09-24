---
id: S-0112
type: story
nature: remediation
title: "A restarted flai serve never starts an agent for a story that was already ready: the stale rule predates flai host"
status: ready
owner: alex
created: 2026-09-24T06:17:48Z
updated: 2026-09-24T06:18:31Z
transitions:
  - to: ready
    at: 2026-09-24T06:18:31Z
    by: alex
tags: [cli]
touches: [flai/internal/serve]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0112 A restarted flai serve never starts an agent for a story that was already ready: the stale rule predates flai host

## Goal

A story in ready gets its agent whether it entered ready before or after the serving flai started, so that a serve restart (an upgrade, a host restart, a crash) never leaves ready stories waiting for nothing.

## Acceptance criteria

- [ ] a story in ready that no agent has been started for since it entered ready is started by the next look after `flai serve` starts, the in-progress limit and the attended rule permitting, whether or not it was ready before serve started
- [ ] a story an agent was already started for since it entered ready is not started again after a restart (`serve/agents.json` already records each story's runs)
- [ ] a story that is skipped for any reason is said in `agent.status`'s `waiting` and the log, so that a ready story never sits without a stated reason
- [ ] the launcher tests that assert "a restart starts nothing" are replaced by ones that assert "a restart starts nothing twice"

## Tasks

## Notes

**What was seen (2026-09-24).** S-0109 entered ready at 03:52Z. `flai host` started a new `flai serve` at 04:44Z (1.15.1) and again at 05:37Z (1.15.2). Neither serve has logged anything about S-0109: no `agent started`, no `agent not started`. The last agent activity in `~/.flai/serve/serve.log` is 02:08Z.

**Why.** `internal/serve/agents.go` (`launcher.look`): the launcher's first look records every story then in ready as *stale*, and a stale story is never a candidate until it leaves ready. The comment says why: "a restart of flai serve must not start a session nobody asked for" (S-0079, kept in S-0104). That rule was written when the operator started `flai serve` by hand and a restart was rare. Since S-0106, `flai host` restarts serve on every upgrade, every host restart, and every crash, so each restart silently retires every ready story until someone moves it out of ready and back. The skip is silent: a stale story produces no `waiting` reason and no log line.

**Why the rule is no longer needed.** Its purpose, not starting a second session for a story that already has one, is served since S-0104 by `serve/agents.json`: each story's newest run, with its start time, survives a restart, and `look` already skips a story tried since it entered ready (`startedBefore`) and settles runs that outlived an earlier serve (`settleOrphans`). With that, the first look can behave like any other.

**What the fix leaves alone.** The attended rule (an agent connected over MCP or writing a narrative recently means someone will pull the story) and the in-progress limit still hold; both say why in `waiting`. A story that names no harness while no command is set still waits, and says so.

**Also seen, not the cause.** S-0109 itself carries no `agent` (it was created before the project's default was set), and the host has no command, so once the stale rule is gone it will wait with "names no harness, and no command is set" until it is given one with `flai edit S-0109 --harness ...` or from its page.
