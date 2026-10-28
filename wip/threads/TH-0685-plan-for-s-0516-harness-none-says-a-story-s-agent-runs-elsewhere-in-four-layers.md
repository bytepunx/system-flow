---
id: TH-0685
title: "Plan for S-0516: harness none says a story's agent runs elsewhere, in four layers"
anchor:
  path: wip/kanban/stories/S-0516-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md
  item: S-0516
status: open
participants: [planner-S-0516, orchestrator]
created: 2026-10-28T13:14:54Z
updated: 2026-10-28T13:15:10Z
---

# TH-0685 Plan for S-0516: harness none says a story's agent runs elsewhere, in four layers

On wip/kanban/stories/S-0516-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md.

## Entries

### 2026-10-28T13:14:54Z planner-S-0516
Planned by planner-S-0516 on 2026-10-28. The story remediates I-0214 and says to propose a solution from the instances before building it; the plan assumes this one, and the first task records it as an ADR, so say here if you want another.

The remediation assumed. A reserved harness name, `none`, on a story's agent says its agent runs elsewhere, by hand or in a container nobody has made a harness for yet. `flai serve` starts nothing for such a story: not when it enters ready, not as an automatic restart (ADR-0108), not on a resume or a send back; it names the reason in `waiting`, as it names a story with no harness and no command. `flai serve agent start`, `restart`, and `commit` refuse it saying so. `flai story start` and `story_start` work as they do. `harness.For` refuses the name too, so nothing above the adapters starts it by mistake. Today a story with no harness is started with the host's command, and any name but `claude-code` and `command` fails the start as a failed run, which is why a story had no way to say it: both instances of I-0214 came from that.

Tasks, in layers:

1. T-2213: the ADR recording the remediation, from I-0214's two instances, with the alternatives (a new field on the story, a host setting, waiting for S-0508).
2. T-2214, after T-2213: `harness.None`, `harness.Elsewhere`, and `For`'s refusal, in `flai/internal/harness`.
3. T-2215 and T-2216, after T-2214, together, since they touch no path in common: T-2215 makes `look`, `restartEnded`, `resume`, and the three send backs in `serve/agents.go` skip such a story, with the tests that reproduce both instances; T-2216 makes `flai serve agent start`, `restart`, and `commit` refuse it.
4. T-2217, after T-2215 and T-2216: `docs/users/flai.md`, `docs/operators/index.md`, `design/system/agent-adapters.md`, and `flai issue close I-0214`.

Figures. Forecast 1h, raised from flai's 31m for the ADR, nine start paths, and the tests in `agents_test.go`; delivery 2026-10-29T02:28Z, 29th in the pull order. Cost of delay 40.00 USD a week, as `flai cod` gives it from the inputs flai set (16m per cycle). Touches: 17 files and the folder `design/adrs`, kept for the ADR's file alone; the Planning note says where each came from.

Assumptions:

- A reserved harness name rather than a new field or a host setting: a story's agent is the one place that says where it runs, and only `flai serve`'s adapters read it. `none` stays once S-0508 makes the container a harness, as the way to say "by hand".
- The strategic agents are not changed: one with `harness: none` fails its start at `harness.For` with the refusal, which is enough for a value no role is expected to carry.
- No dashboard change: the reason reaches the story's page through `waiting` and the refusals through the `agent.start`, `agent.restart`, and `agent.commit` writes, as the no-harness case does today.
- No `after`: the story needs nothing of S-0414, S-0502, or S-0521 on main, though it is held in ready while they are in progress (they touch `harness.go` and the two guides). S-0508, in the backlog, touches the same serve files; whichever starts second syncs over the other.
- `docs/users/flai-reference.md` and `docs/operators/settings.md` are left out: no flag and no setting changes.

Nothing to split, merge, or drop.

### 2026-10-28T13:15:10Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan as written, with the reserved harness name `none` and T-2213's ADR recording it, then finalize S-0516.

## The plan fits

- **Tasks:** five tasks in four layers. T-2213, the ADR, comes first. T-2214 builds `harness.None` and `For`'s refusal on it. T-2215 and T-2216 share no path and run together. T-2217 documents them and closes I-0214.
- **Touches:** 17 files, plus `design/adrs` for the ADR's file alone, named under `### Planning`.
- **Figures:** forecast 1h, kept over flai's 31m with the reason stated. Value 40 USD a week, from `flai cod` on the 16m input the issue step set.
- **No `after`:** right. The holds on S-0414, S-0502 and S-0521 clear when each moves to review.

## Why I recommend `none`

- A story's agent is already the one place that says how its agent starts (ADR-0037), and only `flai serve`'s adapters read it. A reserved name adds no field and no host setting.
- It stays useful once S-0508 makes the container a harness, as the way to say "worked by hand".

## Why it is yours

A reserved value that `flai serve`, `flai serve agent`, and `harness.For` all honour changes the contract of a story's `agent.harness`, a tool-enforced rule. Under `decisions.md` § Rules that needs an ADR, and none records it yet.

**The alternatives** the ADR weighs: a new field on the story, a host setting, or waiting for S-0508's container harness.

The thread stays open for your answer. S-0516 is still a draft, so finalizing it is your call too.

Source: design/conventions/decisions.md § Rules
