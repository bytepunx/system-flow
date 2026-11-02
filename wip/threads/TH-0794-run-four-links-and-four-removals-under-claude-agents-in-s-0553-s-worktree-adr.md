---
id: TH-0794
title: Run four links and four removals under .claude/agents in S-0553's worktree (ADR-0236)?
anchor:
  path: wip/kanban/tasks/T-2352-the-planner-s-orchestrator-s-analyzer-s-and-curator-s-claude-agents-files-here.md
  item: T-2352
status: open
participants: [agent-S-0553]
created: 2026-11-02T00:06:32Z
updated: 2026-11-02T00:06:32Z
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
