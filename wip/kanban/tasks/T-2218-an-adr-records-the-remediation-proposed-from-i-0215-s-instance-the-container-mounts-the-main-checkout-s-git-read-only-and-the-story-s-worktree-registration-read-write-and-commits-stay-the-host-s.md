---
id: T-2218
type: task
nature: remediation
title: "An ADR records the remediation proposed from I-0215's instance: the container mounts the main checkout's .git read-only and the story's worktree registration read-write, and commits stay the host's"
status: backlog
parent: S-0517
owner: alex
created: 2026-10-28T13:19:58Z
updated: 2026-10-28T13:19:58Z
transitions: []
stream: S-0517
tags: [cli, docker, adr]
touches: [design/adrs/README.md]
---
# T-2218 An ADR records the remediation proposed from I-0215's instance: the container mounts the main checkout's .git read-only and the story's worktree registration read-write, and commits stay the host's

## Work

Read I-0215's instance, S-0505's results (§ Git inside the container, § Recommendation), and the code: `flai/internal/storygit/open.go`, where `git worktree add` writes the worktree's `.git` file naming `<root>/.git/worktrees/S-nnnn`, whose `commondir` names the main checkout's `.git`; `scripts/agent-container.sh`, which mounts the whole `.git` read-write under `FLAI_AGENT_GIT_MOUNT=1`; and ADR-0203 § Decision (mounts, flai's MCP server) and ADR-0027 (the dashboard's read-only git paths). Propose the remediation from the instance and record it with `flai adr new`, then `flai adr accept`, refining ADR-0203: a container that runs a story's agent mounts, beside the worktree, the main checkout's `.git` read-only and the story's own `.git/worktrees/S-nnnn` read-write, both at their host paths, so that git in the worktree reads and refreshes its index while the objects, refs, `packed-refs`, `config`, and hooks stay the host's. Say why read-only suffices: every commit is the host's, through `task_done` and the MCP `verify` the close-out uses, and S-0521 makes `flai verify`, `stream sync`, and `task done` refuse inside the container; a `git commit` or `git add` in the container fails on the read-only objects with git's own error. Say what is given up (the agent cannot stash, commit, or change what it has checked out beyond the worktree's `HEAD`) and what stays read-only visible (other worktrees' registrations and every branch's objects, which the agent sees on the host today). Record the alternatives: the whole `.git` read-write, as S-0505 did; a clone of the story branch fetched back, which ADR-0203 keeps for Kubernetes and S-0505's mount timing found not needed from the start; and the objects and refs mounted one by one, which git's layout makes brittle. Add the ADR to `design/adrs/README.md`'s index as `flai adr new` does, and give it topics `harness` and `docker` with `flai adr topics`. If the operator answered the plan's thread with another remediation, record that one instead and adjust the tasks that follow.

## Done when

- An accepted ADR under `design/adrs/` refining ADR-0203 records the two mounts, their modes, why read-only suffices, the alternatives, and its consequences for the container adapter, the hand-start script, and the docs
- `design/adrs/README.md` lists it, its topics are `harness` and `docker`, and `flai check --strict` is clean on the ADR

## Notes

The ADR's file is written under a number and slug no task can name in advance; the story keeps the folder `design/adrs` for it, and `flai task done` widens the touches with the committed path (ADR-0144). Layer 1: the tasks that build the remediation wait for this one, since it fixes the paths and the modes.
