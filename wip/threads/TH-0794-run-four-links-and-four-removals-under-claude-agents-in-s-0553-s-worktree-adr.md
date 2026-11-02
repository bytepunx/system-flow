---
id: TH-0794
title: Run four links and four removals under .claude/agents in S-0553's worktree (ADR-0236)?
anchor:
  path: wip/kanban/tasks/T-2352-the-planner-s-orchestrator-s-analyzer-s-and-curator-s-claude-agents-files-here.md
  item: T-2352
status: resolved
participants: [agent-S-0553, orchestrator, alex]
created: 2026-11-02T00:06:32Z
updated: 2026-11-02T02:27:19Z
---

# TH-0794 Run four links and four removals under .claude/agents in S-0553's worktree (ADR-0236)?

On wip/kanban/tasks/T-2352-the-planner-s-orchestrator-s-analyzer-s-and-curator-s-claude-agents-files-here.md.

## Entries

### 2026-11-02T00:06:32Z agent-S-0553
ADR-0236, committed on story/S-0553, settles where the strategic agents' definitions live. The trial with Claude Code 2.1.311 found that it follows a symbolic link under `.claude/agents`. So each strategic agent's `.claude/agents/<name>.md` becomes a committed link to `../../.flai-cache/agents/<name>.md`, which `flai role definition --write` generates from the role file. The template ships no copy of the four.

`.claude/` is protected, and a symbolic link or a removal cannot go through `permission_prompt`. So I ask you to run these in the worktree. I commit them with T-2352.

```sh
cd /home/alex/git/bytepunx/system-flow/.flai-cache/worktrees/S-0553
for n in planner orchestrator analyzer curator; do
  ln -sf "../../.flai-cache/agents/$n.md" ".claude/agents/$n.md"
  git rm -q "template/root/.claude/agents/$n.md"
done
git status --short .claude template/root/.claude
```

The explorer and the verifier are left alone.

**Recommended:** run them as written. Reply here once done, or say what to change.

One step at acceptance, so that the installed flai keeps working: the links dangle in the main checkout until the generated files exist there. The installed release still reads `.claude/agents/<name>.md` to start the planner, the orchestrator, the analyzer, and the curator. So after accepting S-0553, before the next strategic run, run this once in the main checkout:

```sh
scripts/flai.sh role definition --write
```

I will repeat this step in the review summary.

### 2026-11-02T00:07:40Z orchestrator (recommendation)
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

### 2026-11-02T02:25:48Z alex
Confirmed the recommendation of 2026-11-02T00:07:40Z orchestrator.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision

### 2026-11-02T02:27:19Z agent-S-0553
Resolved: The operator made the four links and four removals; T-2352 committed them as 8ae657a on story/S-0553, and flai role definition --write finds every link in place.
