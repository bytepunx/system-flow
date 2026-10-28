---
id: I-0219
title: flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-28T12:31:04Z
last_reported: 2026-10-28T12:31:04Z
updated: 2026-10-28T13:42:07Z
---

# I-0219 flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home

## Description
flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home

## Instances

### 2026-10-28T12:31:04Z
Story: S-0505.
S-0505 T-2166: agent-S-0513 ran the image's flai (1.71.1 from the tree) as `flai verify S-0513` in its worktree inside flai-agent-claude-code:local, with the worktree and the main checkout's .git mounted. It failed with `verify S-0513: it has no worktree at <root>/.flai-cache/worktrees/S-0513/.flai-cache/worktrees/S-0513`: with the main checkout's files absent, flai took the worktree for the project and looked for the story's worktree under it. It also created /home/agent/.flai/config.json in the container's home. The MCP verify tool, on the host, passed. Remediation for S-0508: flai in the container is told the project and its story's worktree (as FLAI_STORY names the story), or the container's commands that need the main checkout go to the host over MCP, as verify did.

## Remediation

Story S-0521 remediates this issue, created from it at 2026-10-28T12:53:27Z.
Closed 2026-10-28T13:42:07Z: Fixed by S-0521. T-2207: workitem.Open in a story worktree whose main checkout is not here records it in Repo.MainAbsent, and WorktreePath of the worktree's own story is the worktree, never <worktree>/.flai-cache/worktrees/S-nnnn. T-2211: flai verify, stream sync, task done, and check refuse there in one line naming the main checkout and the host's MCP tools (verify, task_done, stream_sync), in place of taking the worktree for the project; TestVerifyRefusesAWorktreeWhoseMainCheckoutIsAbsent reproduces the instance, and each other command has a refusal test. T-2209: scripts/agent-container.sh passes FLAI_CONFIG and FLAI_CACHE_DIR under .flai-cache/agents/S-nnnn/flai, so the container's flai writes nothing under /home/agent. Checked in an image built from the branch with the worktree and the main .git mounted alone: verify, check, and stream sync refused naming /home/alex/git/bytepunx/system-flow, and the home stayed empty.
