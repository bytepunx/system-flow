---
id: S-0779
type: story
nature: remediation
title: flai's MCP address tests no longer count on a port they do not hold staying free
status: done
owner: alex
created: 2026-11-13T20:29:36Z
updated: 2026-11-13T22:46:11Z
transitions:
  - to: ready
    at: 2026-11-13T20:37:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T20:37:57Z
    by: agent-S-0779
  - to: review
    at: 2026-11-13T22:18:02Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T22:46:11Z
    by: alex
tags: [flai]
touches: [flai/cmd/serve_mcp_test.go, flai/cmd/mcp_http_test.go, flai/cmd/serve_mcp.go, flai/cmd/host.go, design/issues/I-0372-flai-s-mcp-address-tests-fail-now-and-then-when-another-test-takes-a-port-they.md, design/issues/summary.md, flai/cmd/host_integration_test.go, flai/cmd/serve_pprof_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6027
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 4
      test_runs: 6
      hand_edits: 3
      work: 85
  models:
    - model: claude-opus-5-5
      input: 206
      output: 61982
      cache_read: 26829667
      cache_write: 329353
      cost: 9.2412
  strategic:
    - kind: planner
      seconds: 380
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 130
          output: 35365
          cache_read: 8512374
          cache_write: 252360
          cost: 4.4292
    - kind: orchestrator
      seconds: 236
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 56
          output: 509
          cache_read: 2150981
          cache_write: 294991
          cost: 0.7195
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: orchestrator
    at: 2026-11-13T20:35:26Z
  value: 300
  by: planner-S-0779
  at: 2026-11-13T20:36:21Z
forecast:
  duration: 45m
  delivery: 2026-11-14T03:37:00Z
  basis: "flai forecast's 21m (137 s per unit of size over 65 done remediation stories on claude-opus-5-5, size 9) raised to 45m for five tasks in three layers, eight tests moved onto a new seam, and a fifty-count race run beside the whole package; delivery moved by the same 24 minutes"
  by: planner-S-0779
  at: 2026-11-13T20:34:24Z
finalized:
  by: orchestrator
  at: 2026-11-13T20:37:03Z
---
# S-0779 flai's MCP address tests no longer count on a port they do not hold staying free

## Goal

The `test (ubuntu-latest, -race)` job of the `flai` workflow fails now and then on a test that loses a race for a port, so a red run on `main` says nothing about the commit and has to be read, then run again. Two such failures in the last eight runs on `main`:

- Run 43495183748 (commit `9b63bacf3`, 2026-11-13), "Integration tests": `TestTheHostMovesAProjectOffAPortAnotherProcessHolds` failed at `serve_mcp_test.go:516` with `the first start, on the remembered port: "127.0.0.1:4243" <nil>`. The same test had passed minutes before in the job's "Behavior tests" step, on the same commit.
- Run 43455794271 (flai 1.93.1): `TestMCPListenFailureClosesWhatItOpened` failed at `mcp_http_test.go:805` with `listen tcp 127.0.0.2:42545: bind: address already in use`.

Both count on a port the test does not hold being free:

- `unusedPort` (`flai/cmd/serve_mcp_test.go` ~483) asks the kernel for a port, closes it, and returns it. The test writes it as the project's remembered address and expects the first start to give it back. `mcpAddrsFor` (`flai/cmd/serve_mcp.go` ~45) gives the remembered port only while `freeOnAll` finds it free, and a port from the default range otherwise, which is the `4243`. Between the close and that probe something else took the port. The test is parallel, in a package of some 640 parallel tests that open listeners and connections on ports the kernel assigns. Six tests in the file use `unusedPort` this way.
- `TestMCPListenFailureClosesWhatItOpened` holds `127.0.0.1:<port>` and counts on `127.0.0.2:<port>` being free, before and after, which nothing ensures.

flai did what it should each time: a remembered port another process holds is moved off (I-0271). Only the tests' assumption broke.

## Acceptance criteria
- [x] No test in `flai/cmd` counts on a port it closed, or never held, being free when flai probes or listens on it: the six tests that use `unusedPort` and `TestMCPListenFailureClosesWhatItOpened` decide for themselves which addresses are free or held
- [x] `TestTheHostMovesAProjectOffAPortAnotherProcessHolds` still shows what I-0271 fixed: a remembered port another process holds is moved off, the new one kept, and the log says so
- [x] At least one test still listens on a real port, so that what `freePort` and `listenMCP` do with the operating system stays covered
- [x] `go test -race -count=50 -run 'MCP|TheHost' ./cmd/` passes from `flai/`, run beside the whole package's tests or under like load, and the narrative records the run
- [x] An issue records the flake with both occurrences, and is closed with what fixed it

## Tasks
- T-3329 Record the MCP port flake as an issue with both occurrences
- T-3330 Give the MCP address choice its free-port check as a value
- T-3331 The address tests say which ports are free and held instead of racing for real ones
- T-3332 The listen-failure test no longer counts on 127.0.0.2 at a port it does not hold staying free
- T-3333 Run the MCP and host tests fifty times under race and load, record the run, and close the issue
- T-3337 The every-address test checks its listeners closed without counting on their ports staying free
- T-3338 The host integration test tries new ports when another process takes the host's or the project's
- T-3339 The profiles test lets flai serve pick the port and reads it from the log

## Notes

- Not confirmed: what took the port. The log does not show it; another parallel test's listener or outgoing connection is the likely holder.
- One way: give `mcpAddrsFor` and the launcher the free-port check as a value (`freePort` in `flai/cmd/host.go` ~796 by default), so a test says which ports are free and held, and no real port is raced for. Another, for a test that needs a real port: hold the listener and hand it over, or retry with a new port when the first probe finds it taken. The story's agent chooses, and records the choice in the narrative's `## Decisions`.
- `unusedPort`'s comment says its port is "outside the default range that the live host and the parallel tests probe", which holds for the default range and not for the kernel's own. Correct or remove it with the helper.
- The other tests that fail this job, `TestRepositoryLintsClean`, are I-0363 and not this story's.
- The failing job's log: `gh api repos/bytepunx/system-flow/actions/jobs/132126294302/logs`.
- No change to flai's behaviour, so no design or user documentation changes.
- Remediates I-0372, which records both occurrences.

### Planning

Touches, all four declared and kept:

- `flai/cmd/serve_mcp.go`, `flai/cmd/host.go` (declared): `mcpAddrsFor`, `freeOnAll`, `hostLauncher.mcpAddrs`, and `freePort` get the free-port check as a value (T-3330).
- `flai/cmd/serve_mcp_test.go` (declared): the tests that use `unusedPort` (T-3330's call sites, T-3331).
- `flai/cmd/mcp_http_test.go` (declared): `TestMCPListenFailureClosesWhatItOpened` (T-3332).
- Not added: `flai/cmd/mcp_http.go` (co-change, 10 of 30 commits). No task changes it unless T-3332's agent gives `listenMCP` its listen as a value; then `flai task done` widens the touches. The other co-change paths are under 25% and nothing in the goal reaches them.
- The tasks' `design/issues/summary.md` (T-3329, T-3333) is left off the story: `design/issues` is in `claims.shared`, and the touches drift leaves it out (ADR-0284). The issue file's name is flai's to give.

No folder touch. No `after`: the story reads nothing another open story writes.

Layers: T-3329, T-3330, and T-3332 first, sharing no path; T-3331 after T-3330, whose value it passes and whose test file it shares; T-3333 last, after T-3329, T-3331, and T-3332, whose tests it runs and whose issue it closes.

Forecast: 45m, raised from `flai forecast`'s 21m. History sizes the story by its 5 criteria and 4 touches. It does not count the eight tests that move onto a new seam, nor a fifty-count race run beside the whole package, which alone takes minutes. Delivery is moved by the same 24 minutes, to 2026-11-14T03:37:00Z.

Cost of delay: 300 USD a week, as `flai cod` works it out and kept. The input, `time_lost_per_cycle` 2h, was set on TH-1182: about ten red runs a week on `main` at some 15 minutes each to read and run again. 2h of time lost per 168h cycle at `planning.hour_rate` 150 is 300 USD a week.
