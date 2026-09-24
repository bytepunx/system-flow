---
id: S-0106
type: story
nature: feature
title: The flai host command should become a process-managing shell
status: done
owner: alex
created: 2026-09-24T01:23:37Z
updated: 2026-09-24T02:57:39Z
transitions:
  - to: ready
    at: 2026-09-24T01:23:43Z
    by: alex
  - to: in-progress
    at: 2026-09-24T01:27:43Z
    by: system-flow
  - to: review
    at: 2026-09-24T01:49:17Z
    by: system-flow
  - to: done
    at: 2026-09-24T02:57:39Z
    by: alex
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0106 The flai host command should become a process-managing shell

## Goal

A new command, `flai host`, is a bootstrap/loader that manages the other processes. For this to work, there needs to be a line of communication from `flai serve` process back to the the host.

The goal is to enable the addition of extensions to the dashboard's host panel that show the status of the other managed processes (flai serve and flai MCP) and enable management commands that can update the flai version and restart those sub-processes.

## Acceptance criteria
- [x] flai host starts and manages child processes for flai serve and flai MCP
- [x] flai serve is able to communicate commands back to the host in order to get status, version, issue upgrade commands, issue restart commands
- [x] flai serve is no longer responsible for managing MCP or dashboard processes directly but instead sends the signal up to the host process which then manages the other sub-processes
- [x] the flai host process stops all flai sub-processes when it crashes or stops (no zombie processes)
- [x] the flai dashboard command no longer starts flai serve or flai MCP processes but instead boots the host process which handles the other two
- [x] one flai host process can run per machine

## Tasks
- T-0374 The host package supervises child processes behind a control API on one machine-wide address
- T-0375 flai host runs, starts, stops, and reports the host, and flai serve runs under it
- T-0376 flai dashboard boots the host instead of starting flai serve
- T-0377 The dashboard asks the host through flai serve for status, version, upgrade, and restart
- T-0378 ADR, design, and documentation for flai host

## Notes

Delivered on `story/S-0106` in five commits, one per task. The design is in [ADR-0040](../../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md).

How each criterion was verified:

- **Children, and none left behind.** `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` (`flai/cmd`, integration tier) builds flai from the tree and runs `flai host start` on a spare port with a temp config. It sees serve and the project's MCP server start, and serve restart without the MCP server moving. After a SIGKILL of the host, both children are gone. After `flai host stop`, everything is gone. `internal/host`'s tests cover restart after a crash, the grace kill, stop, start, and restart per process, and external processes.
- **Serve talks to the host.** `flai serve` tells the host its projects over the control API (`internal/serve` tests with a fake, the integration test for real). The dashboard's `host.status`, `host.check`, `host.process`, and `host.upgrade` run `flai host … --json` from `flai serve` (`internal/hostapi` tests).
- **One per machine.** `TestOneHostHoldsTheAddress`: a second host on the same address is refused. `FLAI_HOST_ADDR` can put a second host on another address, on purpose, for tests and for a port another program already holds.
- **Upgrade: not tried live.** The upgrade was tested with a fake installer: the host answers, stops its children, and returns `ErrRestart`. The real `flai self-upgrade` followed by the host exec'ing the new binary has not been run end to end, because the newest release, 1.14.0, has no `flai host` to come back as. Try `flai host upgrade` once a release with this story is out.
- **The checks.** `flai check --strict` and the smoke repository check report 0 errors and 4 warnings. The warnings are `wip.overlap` with S-0107 and T-0381, which another agent is working at the same time. They clear when either story leaves in-progress. Everything else in smoke passed on its own: template render, markdown lint (worktree and main checkout), and the installer. Behaviour and integration tiers and golangci-lint pass. flaiover's three test files that use `REQUIRED_METHODS` pass.
