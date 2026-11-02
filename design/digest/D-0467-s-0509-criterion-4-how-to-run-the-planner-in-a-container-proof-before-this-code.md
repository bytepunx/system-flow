---
id: D-0467
title: "S-0509 criterion 4: how to run the planner-in-a-container proof before this code is released: Recommended: option 2, written as you answered for S-0508 on TH-0812 (D-0466): move the live planner run…"
kind: answer
sources: [TH-0815, S-0509]
topics: [harness, docker, planning, flai, cli, go, code, all]
paths: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/activity.go, flai/internal/serve/claudecheck.go, design/adrs/README.md, design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/summary.md, design/issues, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/orchestrate_hold.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0467 S-0509 criterion 4: how to run the planner-in-a-container proof before this code is released: Recommended: option 2, written as you answered for S-0508 on TH-0812 (D-0466): move the live planner run…

The operator, alex, answered TH-0815 at 2026-11-02T18:56:22Z.

What was asked:

```text
Criterion 4 asks for a planner run in a container that plans a story and logs its activity in `wip/agents/planner.md`. flai serve cannot do that run during this story:
```

The answer:

```text
**Recommended:** option 2, written as you answered for S-0508 on TH-0812 (D-0466): move the live planner run to a follow-up draft rather than leave criterion 4 unticked.

- S-0509 keeps what it can prove: the four start paths on a fake adapter, the argv and mounts per role, and the serve test that feeds a recorded container log to the planner's run end and checks the activity entry's duration and cost.
- S-0509's agent narrows criterion 4 to that, and writes a draft story in E-0029 with `after: S-0509`, holding the live part: `flai serve` starts a planner on `harness: container`, it plans one backlog story, and `wip/agents/planner.md` gains its entry. Its Notes say it waits for the release and for your host settings (`flai serve agent harness container --login <path>`).
- S-0509's Notes record the move.

**Why:** the same case as D-0466 and D-0461. The role start paths exist only on `story/S-0509`, and this host's `flai serve` is the installed 1.80.2. An unticked criterion would keep S-0509 from the definition of done.

**Why not option 1 now:** it mounts your Claude login into a container and spends a planner session, for a proof the follow-up gives on the real path. If you want it anyway, say so and name the story to plan. **Option 3** breaks your rule against a second `flai serve`; I don't recommend it.

**Why I recommend rather than answer:** it changes a finalized story's criteria and adds a story. That is your call, and so is finalizing the draft.

Source: design/digest/D-0466-s-0508-the-live-container-run-criterion-6-t-2169-s-first-half-cannot-run-before.md
```

Source: TH-0815.

What S-0509 did with it, as its narrative records under Decisions: criterion 4 now covers a planner's container run with docker faked, logged with its duration and cost, and the live `flai serve` run moved to draft S-0567, with `after: S-0509`. S-0509's Notes record the change. S-0509 was accepted on 2026-11-02.

Source: S-0509 (wip/archive/agents/S-0509.md).
