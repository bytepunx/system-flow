---
id: S-0516
type: story
nature: remediation
title: flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere
status: backlog
owner: alex
created: 2026-10-28T12:52:40Z
updated: 2026-10-29T18:03:28Z
transitions: []
tags: [cli, serve]
topics: [harness]
touches: [design/adrs, design/adrs/README.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/internal/serve/commit.go, flai/internal/serve/commit_test.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0214-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 200
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 713
          cache_read: 13032920
          cache_write: 18772
          cost: 2.8652
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-10-28T12:52:40Z
  value: 40
  by: planner-S-0516
  at: 2026-10-28T13:12:42Z
forecast:
  duration: 1h
  delivery: 2026-10-30T02:20:00Z
  basis: "Its own forecast of 1h; 17th in the pull order with an in-progress limit of 5, behind S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511 and S-0512."
  by: flai
  at: 2026-10-29T18:03:28Z
---
# S-0516 flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere

## Goal

This story remediates [I-0214](../../../design/issues/I-0214-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md), "flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0214 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0214 is closed with `flai issue close I-0214 --reason` saying what fixed it

## Tasks
- T-2213 An ADR records the remediation proposed from I-0214's instances: harness none on a story's agent says its agent runs elsewhere, and flai serve starts nothing for it
- T-2214 harness.None is a reserved name: Elsewhere reports an agent that runs elsewhere, and harness.For refuses it saying flai serve starts nothing for it
- T-2215 flai serve starts nothing for a story whose agent runs elsewhere: look names the reason in waiting, and the automatic restart, the resume, and the send backs leave it alone, with tests that reproduce I-0214
- T-2216 flai serve agent start, restart, and commit refuse a story whose agent runs elsewhere, saying flai serve starts nothing for it
- T-2217 The users guide, the operators guide, and agent-adapters.md describe harness none, and I-0214 is closed saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0214. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-28T12:27:00Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0516 on 2026-10-28. The remediation the plan assumes, proposed from I-0214's two instances: a reserved harness name, `none`, on a story's agent says its agent runs elsewhere, by hand or in a container nobody has made a harness for yet. `flai serve` starts nothing for such a story, neither when it enters ready, nor as an automatic restart (ADR-0108), nor on a resume or a send back, and names the reason in `waiting`; `flai serve agent start`, `restart`, and `commit` refuse it saying so; `flai story start` and `story_start` by hand work as they do. Once S-0508 makes the container a harness, `none` stays as the way to say "by hand".

Touches. The story declared none, so `flai touches suggest S-0516` ran on four seeds (112 of 2478 commits): `flai/internal/serve/agents.go`, `flai/internal/harness/harness.go`, `flai/internal/serve/restart.go`, and `flai/internal/serve/start.go`. Of what it listed, `harness_test.go` (57%) and `serve/agents_test.go` (28%) are kept (co-change). The rest are left out: `cmd/serve_actions.go` and its test, `docs/users/flai-reference.md`, and `template/CHANGELOG.md` change only with a new flag or a template file, and this story adds neither; the conventions, `template/template.yaml`, `hostapi/writes.go`, `mcpserver/server.go`, and `design/system/flai-cli.md` change with the serve commands, not with what a harness name means. The rest are the tasks' (layout):

- `design/adrs/README.md`, and the folder `design/adrs` (layout, folder kept): `flai adr new` in the first task writes the ADR under a number and slug no task can name yet; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).
- `flai/internal/harness/harness.go`, `harness_test.go`: the reserved name and the helper beside `For`.
- `flai/internal/serve/agents.go`, `agents_test.go`: the six start paths that check "no harness and no command" today (`look`, `resume`, `startSentBack`, `startQueueSentBack`, `startAcceptSentBack`, `restartEnded`) and the tests that reproduce the two instances.
- `flai/internal/serve/start.go`, `restart.go`, `commit.go` and their tests: the operator's starts, which refuse with the reason.
- `docs/users/flai.md`, `docs/operators/index.md`, `design/system/agent-adapters.md` (design: the harness list, § Who works a story and § Starting an agent when a story becomes ready).
- `design/issues/I-0214-…md`, `design/issues/summary.md`: `flai issue close` (criterion 2).

Overlaps on 2026-10-28: S-0502, in progress, touches `harness.go` and `harness_test.go` (T-2145) and `docs/users/flai.md`; S-0414, in progress, touches `docs/users/flai.md` and `docs/operators/index.md`; S-0521, in progress, touches `docs/users/flai.md` and `design/issues/summary.md`. This story is held in ready until they move to review (ADR-0096); none is named in `after`, since it needs none of their changes on main. S-0508, in the backlog, touches `harness.go`, `serve/agents.go`, `start.go`, `restart.go`, and the same docs: whichever starts second syncs over the other.

Forecast. `flai forecast S-0516` gives 31m (median 95 s per unit of size over 38 done remediation stories on claude-opus-5-5 in the large band; size 19: 2 criteria, 17 touches) and delivery 2026-10-29T01:59Z, 29th in the pull order. Raised to 1h: the story writes an ADR, threads one reserved name through nine start paths in `flai serve` (six automatic, three the operator's) with tests that reproduce both instances in `agents_test.go`, flai's largest test file, and changes three documents. The delivery written is flai's play-out start, 2026-10-29T01:59Z less 31m, plus 1h: 2026-10-29T02:28Z.

Cost of delay. `flai cod S-0516` gives 40.00 USD a week: 16m of time lost per 168h cycle (the inputs flai set from I-0214) at 150 USD an hour, one cycle a week. Written as flai gives it.
