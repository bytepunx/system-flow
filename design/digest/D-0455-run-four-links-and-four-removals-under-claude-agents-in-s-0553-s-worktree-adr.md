---
id: D-0455
title: "Run four links and four removals under .claude/agents in S-0553's worktree (ADR-0236): Recommended: alex, run agent-S-0553's commands as written in…"
kind: answer
sources: [TH-0794]
topics: [cli, orchestration, planning, template, conventions, flai, go, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/internal/roles/definition.go, flai/internal/roles/definition_test.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/guard/writes.go, flai/internal/guard/writes_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/cmd/role.go, flai/cmd/role_test.go, ".claude/agents/planner.md", ".claude/agents/orchestrator.md", ".claude/agents/analyzer.md", ".claude/agents/curator.md", template/root/.claude/agents/planner.md, template/root/.claude/agents/orchestrator.md, template/root/.claude/agents/analyzer.md, template/root/.claude/agents/curator.md, flai/internal/harness/harness.go, flai/internal/serve/orchestrate_test.go, design/issues/I-0244-testaheldwaitfailsonatruncatedstory-times-out-instead-of-failing-on-the.md, design/issues/summary.md, docs/operators/settings.md, flai/internal/pilot/run_test.go, flai/internal/serve/plan_test.go, design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0455 Run four links and four removals under .claude/agents in S-0553's worktree (ADR-0236): Recommended: alex, run agent-S-0553's commands as written in…

The operator, alex, answered TH-0794 at 2026-11-02T02:25:48Z.

What was asked:

```text
ADR-0236, committed on story/S-0553, settles where the strategic agents' definitions live. The trial with Claude Code 2.1.311 found that it follows a symbolic link under `.claude/agents`. So each strategic agent's `.claude/agents/<name>.md` becomes a committed link to `../../.flai-cache/agents/<name>.md`, which `flai role definition --write` generates from the role file. The template ships no copy of the four.
```

The answer:

```text
**Recommended:** alex, run agent-S-0553's commands as written in `/home/alex/git/bytepunx/system-flow/.flai-cache/worktrees/S-0553`, then reply `done`. You did the same kind of step for S-0224 (D-0188) and S-0270 (D-0212).

**Why:**

- The commands match ADR-0236 § Decision, committed on `story/S-0553` (c246b8af1). Each `.claude/agents/<name>.md` becomes a relative link to `../../.flai-cache/agents/<name>.md`, and the template ships no copy of the four. The explorer's and verifier's files stay as they are.
- I checked the worktree. Both `.claude/agents/` and `template/root/.claude/agents/` hold the four files the commands replace and remove, plus `explorer.md` and `verifier.md`, which the commands leave alone.
- `.claude/` is a protected path. `permission_prompt` cannot make a symbolic link or run `git rm`, so only you can run these commands.
- Under ADR-0106 § 1, a branch that changes a protected path is accepted by the operator only. So S-0553 will come to you for acceptance whatever happens here.

**One thing to watch at acceptance:** after the merge, the four links in the main checkout dangle until you run `scripts/flai.sh role definition --write` there. The installed flai reads `.claude/agents/<name>.md` when it starts the planner, the orchestrator, the analyzer, and the curator. Until you run that command, those runs will fail. Run it straight after accepting S-0553.

I marked this operator-only because only you can carry it out. If `flai serve` confirmed it automatically after `confirm_after`, the thread would be marked answered with nothing run.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision
```

Source: TH-0794.
