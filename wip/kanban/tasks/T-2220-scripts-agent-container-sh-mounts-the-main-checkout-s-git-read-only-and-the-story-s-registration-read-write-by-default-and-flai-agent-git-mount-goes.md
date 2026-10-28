---
id: T-2220
type: task
nature: remediation
title: scripts/agent-container.sh mounts the main checkout's .git read-only and the story's registration read-write by default, and FLAI_AGENT_GIT_MOUNT goes
status: backlog
parent: S-0517
owner: alex
created: 2026-10-28T13:20:34Z
updated: 2026-10-28T13:20:34Z
transitions: []
stream: S-0517
tags: [cli, docker, scripts]
touches: [scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, scripts/README.md]
after: [T-2218]
---
# T-2220 scripts/agent-container.sh mounts the main checkout's .git read-only and the story's registration read-write by default, and FLAI_AGENT_GIT_MOUNT goes

## Work

In `scripts/agent-container.sh`, replace the `FLAI_AGENT_GIT_MOUNT=1` branch, which mounts the whole `.git` read-write, with the two mounts the ADR of T-2218 records, always: `-v "$CACHE_ROOT/.git:$CACHE_ROOT/.git:ro"` and `-v "$CACHE_ROOT/.git/worktrees/$STORY:$CACHE_ROOT/.git/worktrees/$STORY"`, read-write, in that order, since Docker nests the longer destination over the shorter. Refuse the start, naming the path, when `$CACHE_ROOT/.git/worktrees/$STORY` is not a folder, as the script refuses a missing worktree. Drop `FLAI_AGENT_GIT_MOUNT` and `GIT_ARGS`, and say in the header comment what is mounted and that commits are the host's through `task_done`. Keep the script POSIX `sh` under `set -eu`, quoted, and ShellCheck-clean (`flai test scripts/agent-container.sh`). Check it by hand with the image built (`make agent-images`) on a story worktree: `scripts/agent-container.sh S-nnnn -- git status` prints the branch, and `scripts/agent-container.sh S-nnnn -- git commit --allow-empty -m probe` is refused by git on the read-only objects; record both outcomes in the narrative.

Update `docs/operators/runbooks/agent-container.md`: the mounts table gains the two rows (mode and what for), loses the `FLAI_AGENT_GIT_MOUNT` row, the variables table loses the variable, and the sentence "Without `FLAI_AGENT_GIT_MOUNT=1`, git in the worktree fails" becomes what holds now: git reads in the worktree, and a commit is the host's through `task_done`. Update the script's row in `scripts/README.md` the same way.

## Done when

- `scripts/agent-container.sh S-nnnn -- git status` works with no variable set, `git commit` inside is refused by git, nothing under `.git` but `.git/worktrees/S-nnnn` is writable from the container, and the script refuses a story whose registration is absent, naming it
- The runbook's mounts and variables tables and `scripts/README.md` say what is mounted and that commits are the host's, `FLAI_AGENT_GIT_MOUNT` appears in neither, and ShellCheck and the markdown lint pass

## Notes

Layer 2, with T-2219: both wait for the ADR and touch no path in common. S-0521 (in progress, T-2209 done) changed the script and the runbook to pass `FLAI_CONFIG` and `FLAI_CACHE_DIR`; this story waits on S-0508 and syncs over that change. S-0507's T-2174 (backlog) touches the script too: whichever starts second syncs over the other.
