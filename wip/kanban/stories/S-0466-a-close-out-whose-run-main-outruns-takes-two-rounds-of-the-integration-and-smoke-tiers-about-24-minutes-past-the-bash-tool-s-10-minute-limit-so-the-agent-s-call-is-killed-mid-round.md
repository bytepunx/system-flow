---
id: S-0466
type: story
nature: improvement
title: A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round
status: in-progress
owner: alex
created: 2026-10-21T11:57:48Z
updated: 2026-10-22T21:04:13Z
transitions:
  - to: ready
    at: 2026-10-22T20:50:02Z
    by: alex
  - to: in-progress
    at: 2026-10-22T21:00:00Z
    by: agent-S-0466
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, design/adrs/README.md, flai/cmd/closeout.go, flai/cmd/closeout_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/closeout/detach.go, flai/internal/closeout/detach_test.go, scripts/close-out.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/README.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/workflow.md, design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 268
  priced_by: harness
  turns:
    - day: 2026-10-22
      ceremony: 1
      work: 39
  models:
    - model: claude-opus-5-5
      input: 92
      output: 24815
      cache_read: 4156036
      cache_write: 127733
      cost: 2.3497
  strategic:
    - kind: orchestrator
      seconds: 293
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 158
          output: 2756
          cache_read: 33835404
          cache_write: 43577
          cost: 7.4808
cost_of_delay:
  inputs:
    time_lost_per_cycle: 25m
    by: flai
    at: 2026-10-21T11:57:48Z
  value: 62.5
  by: planner-S-0466
  at: 2026-10-22T20:59:56Z
forecast:
  duration: 3h
  delivery: 2026-10-23T00:00:00Z
  basis: "flai's 42m at size 28 (2 criteria, 26 touches), 89 s a unit over 87 done improvement stories on claude-opus-5-5, raised to 3h for a Go reimplementation of the close-out against git fixtures, process-group detachment with signal tests, an ADR awaiting the operator, and two .claude/ writes behind permission prompts; delivery from the start at 2026-10-22T21:00Z."
  by: planner-S-0466
  at: 2026-10-22T21:04:13Z
finalized:
  by: alex
  at: 2026-10-21T20:51:54Z
---
# S-0466 A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round

## Goal

This story remediates [I-0148](../../../design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md), "A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0148 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0148 is closed with `flai issue close I-0148 --reason` saying what fixed it

## Tasks
- T-2030 An ADR settles how a close-out outlives the agent's call: flai close-out runs the rounds detached and is followed by calls that return within the tool's limit
- T-2031 flai close-out runs the close-out's rounds in Go, verify, commit, sync check, stream sync, at most three rounds, with the script's last line
- T-2033 The close-out's run is detached from the call that starts it, kept under .flai-cache/close-out, followed by calls that return within --wait, and survives the caller's death
- T-2046 scripts/close-out.sh, here and in the template, calls flai close-out and passes the git commit options through
- T-2050 The conventions, the agents' prompts, the user docs, and the design describe the detached close-out and the calls that follow it, in place of one 600000 ms call
- T-2054 I-0148 is closed with the reason naming flai close-out and its detached run

## Notes

Cost of delay inputs set by flai from I-0148. time_lost_per_cycle 25m: 25m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-15T00:07:49Z, 6.5 days before this story; under one cycle counts as one).

### Planning

- **Plan thread.** TH-0596 names the proposed fix, the tasks, their layers, and the assumptions: the close-out moves into flai as `flai close-out`, detached from the call that starts it and followed by calls that return within `--wait`, with `scripts/close-out.sh` a wrapper on it; T-2030's ADR confirms or amends it.
- **Layers.** T-2030 (the ADR) first; T-2031 (the rounds in Go) after it; T-2033 (the detached run) and T-2046 (the scripts) after T-2031, sharing no path; T-2050 (conventions, prompts, docs, design) after both; T-2054 (closing I-0148) last.
- **Touches.** 26, one folder among them, every one on a task.
  - **Declared:** none; the story was made from I-0148 with no touches, and `flai touches suggest S-0466` refused for want of a path to start from, so there is no co-change list.
  - **From the design:** `design/adrs/README.md` for the ADR's index row; `design/conventions/delegation.md` and `design/conventions/work-management.md` with their template copies, which give the 600000 ms call and describe the close-out; `design/system/flai-cli.md` (the Commands table) and `design/system/workflow.md` § Branches and collisions, which describe the close-out's rounds; `design/issues/I-0148-….md` and `design/issues/summary.md`, which `flai issue close` rewrites.
  - **From the code layout:** `flai/cmd/closeout.go` and `flai/internal/closeout/{closeout,detach}.go` with their tests, a new command and package beside `flai/cmd/verify.go`, `flai/cmd/checks.go`, and `flai/internal/verify`, whose `Verify`, `SyncOnly`, and detach patterns it reuses; `scripts/close-out.sh`, `scripts/README.md`, and their template copies, which the template tier keeps in step; `flai/internal/harness/harness.go` and `harness_test.go`, whose delegation prompt repeats the 600000 ms call; `.claude/agents/verifier.md` and its template copy, step 4 of which names it; `docs/users/flai.md` and the generated `docs/users/flai-reference.md`.
  - **Folder touch kept:** `design/adrs/`, for the ADR T-2030 adds, whose number `flai adr new` allocates when it runs; `flai task done` widens the touches to the file. While the story is in progress it holds a ready story that touches an ADR.
- **Forecast.** flai's 42m at size 28 (2 criteria, 26 touches), 89 s a unit over 87 done improvement stories on claude-opus-5-5 in the large band, raised to 3h: a Go reimplementation of the close-out against git fixtures, process-group detachment with tests that kill the caller, an ADR that waits for the operator's agreement, and two `.claude/` writes that each ask through `permission_prompt` are work the touch count does not weigh. Delivery 2026-10-23T00:00Z, 3h from the start at 2026-10-22T21:00Z, the story being in progress when flai forecast it.
- **Cost of delay.** 62.50 USD a week, flai's figure from the operator's input: 25m lost per 168h cycle at 150 USD an hour, 1.00 cycles a week. Unchanged: the input rests on one instance, and the second round it measures would recur with every close-out main outruns.
- **Tags and topics.** Tags `flai, template` added: `flai check` refused a task while the story touched both components with no component tag. Topics `cli, conventions, template`, which the tasks reach.
