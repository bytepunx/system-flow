---
id: TH-0132
title: "S-0218 plan: six tasks in three layers, and the assumptions behind them"
anchor:
  path: wip/kanban/stories/S-0218-the-orchestrator-is-a-long-running-agent-per-project-behind-the-orchestrate-host-action-with-permissions-the-operator-sets-and-a-decision-log.md
  item: S-0218
status: resolved
participants: [planner-S-0218, alex]
created: 2026-10-05T04:48:40Z
updated: 2026-10-05T04:57:36Z
---

# TH-0132 S-0218 plan: six tasks in three layers, and the assumptions behind them

On wip/kanban/stories/S-0218-the-orchestrator-is-a-long-running-agent-per-project-behind-the-orchestrate-host-action-with-permissions-the-operator-sets-and-a-decision-log.md.

## Entries

### 2026-10-05T04:48:40Z planner-S-0218
The plan for S-0218, drafted as six tasks in the backlog.

| Layer | Task | Waits for | What it does |
|-------|------|-----------|--------------|
| 1 | T-0881 | nothing | `orchestration.permissions` and `orchestration.agent` in the manifest, `OrchestrationAgent()`, and `manifest.orchestration` findings |
| 1 | T-0882 | nothing | `orchestratePrompt`, the claude-code adapter's `--agent orchestrator`, `.claude/agents/orchestrator.md` and the template's copy, the settings' edit hook for `orchestrate`, and a template release |
| 2 | T-0883 | T-0881, T-0882 | the `orchestrate` host action; `flai serve` starts one run per project, restarts it when it ends, and stops it when the action goes off; the run in `agents.json` and `agent.status`, its end in `orchestrator.md` |
| 2 | T-0887 | T-0881 | `flai guard` with `FLAI_ROLE=orchestrate`: each permission's calls, a refusal naming the permission, and the refusal logged in `orchestrator.md`; the MCP tool `plan` open to the orchestrator for backlog epics |
| 3 | T-0889 | T-0883 | the orchestrator's run on the Activity page and `orchestrate` on the Settings page |
| 3 | T-0892 | T-0881, T-0882, T-0883, T-0887 | a new ADR, `strategic-agents.md`, `flai-cli.md`, `project-manifest.md`, the operator guide, and the user guide and reference |

No two tasks of a layer share a path. I propose no split, merge, or drop.

The story is now forecast at 1h40m, delivered about 10:52Z after S-0217, at 102.04 USD a week of cost of delay. The reasons are under `### Planning` in its Notes.

Assumptions for you to confirm or correct:

1. `orchestration.policy`: the second criterion names `throughput` or `cost_of_delay`, but S-0217's T-0809, which lands first, defines it as `cod`, `wsjf`, `throughput`, or `fifo`, the names of `flai order --by`. I planned on S-0217's definition and added no policy value here. My recommendation: keep S-0217's names, and read the criterion as satisfied by them.
2. `answer_threads` takes `off` as well as `recommend` and `autonomous`, so that it can be off by default like the others.
3. The guard maps each permission to calls as follows. `plan_backlog_epics` allows `plan` on a backlog epic. `finalize_drafts` allows `flai edit --no-draft`. `promote_to_ready` allows moves to ready. `order_ready` allows `flai order` writes. `answer_threads` allows `thread_reply` and `flai reply`. `accept_reviews` allows `flai accept`. `publish` allows `flai release --pending` and `flai push`. Reads, `inbox`, `board`, `activity_log`, `wait_for_events`, `thread_open`, and S-0217's evaluations always pass. Telling `recommend` apart from `autonomous` in what a reply says is S-0220's.
4. A refusal is logged as an entry in `wip/agents/orchestrator.md`, with no seconds or cost, written by `flai guard` itself.
5. "Visible in the dashboard" is met here by the run on the Activity page, the toggle on Settings, and `activity.document` for kind `orchestrator`, which S-0206 already answers. The Orchestrator page, with its decisions and Stop, stays S-0228's.
6. A run that fails is restarted no sooner than a minute later, so one that fails at once does not spin.
7. The orchestrator's agent is named `orchestrator`, one per project, and its log is `<key>-orchestrator-<start>.log`.
8. T-0882 writes under `.claude/`, which a story's agent asks you for on a thread (S-0257).

### 2026-10-05T04:57:36Z alex
Resolved.
