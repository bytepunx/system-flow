---
id: D-0407
title: flai in an agent container that mounts a story's worktree without the main checkout refuses verify, check, stream sync, and task done, naming the host's MCP tools, rather than taking the worktree for the project
kind: cause
sources: [I-0219, S-0521]
topics: [harness, docker, flai, cli, go, code]
paths: [flai/internal/workitem/repo.go, flai/cmd/verify.go, flai/cmd/check.go, flai/cmd/stream_sync.go, flai/internal/taskdone/taskdone.go, scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md]
date: 2026-10-28
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0407 flai in an agent container that mounts a story's worktree without the main checkout refuses verify, check, stream sync, and task done, naming the host's MCP tools, rather than taking the worktree for the project

I-0219 was closed on 2026-10-28. The cause: inside the agent container, which mounts the story's worktree but not the main checkout, `flai verify` took the worktree for the project root, looked for `<worktree>/.flai-cache/worktrees/S-nnnn`, and wrote a config under the container's home.

The fix, by S-0521, from the issue's close:

```text
T-2207: workitem.Open in a story worktree whose main checkout is not here records it in Repo.MainAbsent, and WorktreePath of the worktree's own story is the worktree, never <worktree>/.flai-cache/worktrees/S-nnnn. T-2211: flai verify, stream sync, task done, and check refuse there in one line naming the main checkout and the host's MCP tools (verify, task_done, stream_sync), in place of taking the worktree for the project [...] T-2209: scripts/agent-container.sh passes FLAI_CONFIG and FLAI_CACHE_DIR under .flai-cache/agents/S-nnnn/flai, so the container's flai writes nothing under /home/agent.
```

Why a refusal and not making these commands work in the container, from S-0521's Decisions (`wip/archive/agents/S-0521.md`):

```text
The fix is a refusal, not running these commands in the container: `wip/` and the main checkout stay on the host, reached through flai mcp (agent-adapters.md, ADR-0019); S-0510 builds the container's close-out on it.
```

For a later container story (S-0510, S-0520): an agent in the container verifies, syncs, and closes tasks through the host's MCP tools; do not mount the main checkout or teach these commands to run without it.

Sources: I-0219, S-0521.
