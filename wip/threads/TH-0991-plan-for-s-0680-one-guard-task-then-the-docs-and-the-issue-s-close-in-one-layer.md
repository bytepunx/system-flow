---
id: TH-0991
title: "Plan for S-0680: one guard task, then the docs and the issue's close in one layer"
anchor:
  path: wip/kanban/stories/S-0680-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md
  item: S-0680
status: resolved
participants: [planner-S-0680, orchestrator, agent-S-0680, alex]
created: 2026-11-10T08:47:45Z
updated: 2026-11-10T09:08:45Z
---

# TH-0991 Plan for S-0680: one guard task, then the docs and the issue's close in one layer

On wip/kanban/stories/S-0680-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md.

## Entries

### 2026-11-10T08:47:45Z planner-S-0680
S-0680 is planned as three tasks in two layers.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2901 The guard counts the read forms of flai serve among the reads | `flai/internal/guard/guard_tables.go`, `guard.go`, `hostconfig.go`, `guard_test.go`, `hostconfig_test.go`, `orchestrate_test.go` |
| 2 | T-2902 The design and the user guide say the read forms of flai serve are open to every agent (after T-2901) | `design/system/strategic-agents.md`, `docs/users/flai.md` |
| 2 | T-2903 Close I-0318 with what fixed it (after T-2901) | `design/issues/I-0318-…-reads.md`, `design/issues/summary.md` |

T-2902 and T-2903 share no path, so they run together.

Proposed fix: add `serve` to `flagReads` with a form function that passes only the read forms ADR-0220 lists. Those are `status`, `actions`, `journal`, `agent` bare or `show`, `agent harness` and `agent provider` bare or with a name alone, `checks show`, `components show`, `forge` bare or with a forge alone, `notify show`, `mcp show`, and `import list`. Every form that starts, stops, or writes the host stays refused.

Figures: forecast 45m, delivery 2026-11-10T12:19:00Z. flai gave 22m; I raised it for the exhaustive table test and the integration and smoke tiers. Cost of delay value 5 USD a week, from the operator's 2m-per-cycle input.

Assumptions:
- These read forms open to every agent the guard holds to its reads: the orchestrator, the planner, the analyzer, the curator, defined roles, and every sub-agent. I-0318 calls that scope the operator's choice. I took the story's creation from the issue as that choice, because `askForHost` already tells all of those agents the forms are open. Alternative: open `flai serve status` to the orchestrator alone, through `design/workflows/roles/orchestrator.md`. That would replace T-2901's table entry with a role-file entry, and the sub-agent tests would go.
- No new ADR. Earlier read additions (S-0428 digest, S-0491 subproject, S-0612 hostapi) took none. The story's agent may add one if it judges ADR-0220's note on sub-agents needs superseding. `design/adrs` is a shared path, so that would hold nothing.
- No `after` on another story, and no folder touch.

### 2026-11-10T08:58:04Z orchestrator (recommendation)
**Recommended:** confirm the planner's scope. The read forms of `flai serve` that ADR-0220 lists open to every agent the guard holds to its reads: the orchestrator, the planner, the analyzer, the curator, defined roles, and read-only sub-agents. Every form that starts, stops, or writes the host stays refused. The plan stands as written: T-2901, then T-2902 and T-2903 together, 45m, 5 USD a week.

**Why:**

- ADR-0220 already names `flai serve status`, `actions`, `journal`, and the show and list forms as reads, and refuses a sub-agent one only because ADR-0060's read table does not list it. Adding `serve` to that table is how reads have been opened before: `flai digest` (S-0428), `flai subproject` (S-0491), and `flai components list` and `serve` for sub-agents (D-0380). None of those took a new ADR.
- `askForHost` already tells every one of these agents the forms are open. Opening them to the orchestrator alone would leave the guard and that message still disagreeing for the rest, which is the inconsistency I-0318 records.
- These forms only read the host's state. None of them changes what `flai serve` may do, so ADR-0220's reason for its refusals does not reach them.

**Why I recommend rather than approve:** I-0318 says that which agents may run these reads is your choice. The plan is otherwise approved: its three tasks, ten file touches, and figures fit the story. T-2901 tells the story's agent to go on unless you answer otherwise here.

**Alternative:** open `flai serve status` to the orchestrator alone, through its role file. T-2901's table entry then gives way to a role-file entry, and its sub-agent tests go.

Source: design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md

Source: design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md

### 2026-11-10T09:02:03Z agent-S-0680
Went ahead with the plan as written. T-2901 is done (27cc458). The read forms ADR-0220 lists now pass for the orchestrator, the planner, the analyzer, the curator, and sub-agents. Every form that runs, starts, stops, commits, or writes the host is still refused. Along the way I found that the container harness's `--image`, `--network`, and `--login` were not counted as host writes, so I fixed that in the same task. If you want these reads open to the orchestrator alone, say so here and I will narrow it before review.

### 2026-11-10T09:07:39Z orchestrator (recommendation)
**Recommended:** keep S-0680's scope as built.

- The read forms of `flai serve` that ADR-0220 lists stay open to every agent the guard holds to its reads. That covers the orchestrator, the planner, the analyzer, the curator, defined roles, and read-only sub-agents.
- Keep the container harness's `--image`, `--network`, and `--login` among the refused writes. The agent added them in T-2901 (27cc458).

**Why:**

- Your answer on the scope is still owed. agent-S-0680's entry cleared my earlier recommendation, so I am posting it again. That lets it be confirmed, or answered otherwise.
- The scope is unchanged from my 08:58 recommendation. ADR-0220 § Decision names these forms as reads. `askForHost` already tells every one of these agents that the forms are open. Opening them to the orchestrator alone would leave the guard and that message disagreeing for everyone else.
- The harness flags fall under ADR-0220's rule. It refuses every command that writes the host configuration, in any form. `flai serve agent harness` with a name and a setting is in its table of writes. `--image`, `--network`, and `--login` are such settings, added after the ADR was written. Counting them as writes closes a gap under the same decision, and the change sits inside T-2901's touches. It refuses more than before and opens nothing.
- All three tasks are done.

**Alternative:** open the reads to the orchestrator alone, through `design/workflows/roles/orchestrator.md`. agent-S-0680 narrows it before review if you answer that here.

Source: design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md § Decision

### 2026-11-10T09:07:50Z alex
Confirmed the recommendation of 2026-11-10T09:07:39Z orchestrator.

Source: design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md § Decision

### 2026-11-10T09:08:45Z orchestrator
Resolved: The operator confirmed the 09:07:39Z recommendation, so S-0680 keeps its scope as built. The flai serve read forms that ADR-0220 lists are open to every guarded agent. The container harness's --image, --network, and --login count as host writes. All three tasks are done.
