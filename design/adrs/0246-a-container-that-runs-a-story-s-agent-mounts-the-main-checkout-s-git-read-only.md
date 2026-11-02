---
id: ADR-0246
title: "A container that runs a story's agent mounts the main checkout's .git read-only and the story's worktree registration read-write over it, so git works in the worktree and every commit stays the host's"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0203]
topics: [harness, docker]
---

# ADR-0246 A container that runs a story's agent mounts the main checkout's .git read-only and the story's worktree registration read-write over it, so git works in the worktree and every commit stays the host's

## Context

`flai stream open` makes a story's worktree with `git worktree add` (`flai/internal/storygit/open.go`). The worktree's `.git` is a file, `gitdir: <root>/.git/worktrees/S-nnnn`. That folder, the story's registration, holds the worktree's `index`, `HEAD`, logs, and per-worktree refs, and a `commondir` of `../..` that names the main checkout's `.git`. Under `worktrees.relative_paths` the `gitdir` line is relative to the worktree instead ([ADR-0022](0022-repository-mounted-at-its-host-path.md)).

[ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) § Mounts gives a story's container the worktree, the story's cache folder, and the login, and "no main checkout". § flai's MCP server says the container holds "no `.git` of the main checkout". S-0505 found that a worktree mounted alone is not a git repository: `git status` fails with `fatal: not a git repository` ([I-0215](../issues/I-0215-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md)). The spike got round it with `FLAI_AGENT_GIT_MOUNT=1` in `scripts/agent-container.sh`, which mounts the whole `.git` read-write. Then the container can write every branch, every other worktree's index, and the hooks and config of the whole repository ([S-0505's results](../experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md), § Git inside the container).

What the agent's git needs in the container is to read: `git status`, `diff`, `log`, and `show` of its own work. The one write those make is the index refresh `git status` does, with its lock, in the registration. Every commit is the host's: `task_done` and the close-out's MCP `verify` run on the host over HTTP, in the worktree the container mounts. S-0521 makes `flai verify`, `flai stream sync`, and `flai task done` refuse inside the container.

## Decision

A container that runs a story's agent mounts, beside the worktree, the main checkout's `.git` read-only and the story's own registration, `.git/worktrees/S-nnnn`, read-write over it, both at their host paths.

Git in the worktree then reads the objects, refs, `packed-refs`, and `config` from the common dir, and reads and writes its own `index`, `HEAD`, logs, and locks in the registration. The objects, refs, `packed-refs`, `config`, and hooks stay read-only, as [ADR-0027](0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md) keeps the dashboard's hooks, config, and info.

The paths come from the worktree itself: its `.git` file names the registration, and the registration's `commondir` names the common dir, each resolved against the folder that holds it. A relative `gitdir` under `worktrees.relative_paths` resolves to the same host paths, so the mounts are the same. A worktree whose `.git` is a folder carries its repository, and gets no extra mount. A worktree whose `.git` file or registration is missing, or whose registration has no `commondir`, fails the start with the path, rather than starting a container without git.

This narrows ADR-0203. Where it said "no `.git` of the main checkout", it now says no writable `.git` but the story's registration. Nothing else of the main checkout is mounted.

Read-only suffices, because:

- Every commit is the host's, through `task_done` and the MCP `verify` of the close-out.
- S-0521 makes `flai verify`, `stream sync`, and `task done` refuse inside the container, so nothing in it tries to commit for the agent.
- A `git commit`, a `git add` of new content, a `git stash`, or a branch made in the container fails on the read-only objects or refs with git's own write error. The docs say the host commits.

## Consequences

- The container adapter: `Start` in `flai/internal/harness/container.go` adds the two mounts to the `docker run` argument list, the common dir with `:ro` and the registration after it. A helper beside `open.go`, `storygit.GitMounts`, reads the worktree's `.git` file and `commondir` and gives the two paths and their modes. A normal checkout gives none, and an absent registration fails the start. The adapter's check of what a container never mounts allows these two paths under the main checkout and nothing else. The argument-list test asserts that the common dir is read-only and that nothing under `.git` but the story's registration is writable.
- The hand-start script: `scripts/agent-container.sh` mounts both by default, and `FLAI_AGENT_GIT_MOUNT` goes.
- The docs: the runbook `docs/operators/runbooks/agent-container.md`, `docs/users/flai.md`, `design/system/agent-adapters.md`, and `design/tech/docker.md` name the two git mounts and their modes, and say that commits are the host's.
- What the agent gives up in the container: it cannot stash, commit, add new content, or make or move a branch. It can still change what its worktree has checked out, since a checkout of an existing branch or commit writes only the worktree's files, `index`, and `HEAD`, all its own; an agent on the host can do the same today, and the host's commands meet it as they do there.
- What the agent sees read-only: the other worktrees' registrations, every branch's objects and refs, and the repository's `config` and hooks. An agent on the host sees all of these today.
- The registration is the agent's to write, and its `config.worktree` is part of it. Git reads that file only when the common `config` sets `extensions.worktreeConfig`, which the container cannot set. Hooks written into the registration do not run: git takes hooks from the common dir, as ADR-0027 found.
- A read-only bind mount of the common dir shows the host's later writes, such as a commit `task_done` makes, since the mount is of the folder, not of a copy.
- [I-0215](../issues/I-0215-a-story-s-worktree-mounted-alone-into-a-container-is-not-a-git-repository-since-its-git-file-names-the-main-checkout-s-git-worktrees-folder.md) is closed by this once the adapter and the script mount the two paths.

## Alternatives considered

- The whole `.git` read-write, as S-0505 did with `FLAI_AGENT_GIT_MOUNT=1`. Git works, and so does a commit in the container, but the container can write every branch, every other worktree's index, and the hooks and config the host's git later runs. That is the route ADR-0027 closed for the dashboard, and no work of the agent needs it.
- A clone of the story branch in the container, fetched back by flai. ADR-0203 keeps it for Kubernetes, where there is no host path. On this host, S-0505 measured `git status` at the same 27 ms through the bind mount as on the host, so the clone is not needed from the start. It also changes how `task_done` reaches the worktree.
- The objects, refs, `packed-refs`, and `config` mounted one by one, read-only, beside the registration. It mounts the same as the common dir read-only, but must follow git's layout: a file that git adds to the common dir, a new ref store, or an alternates file would be missing from the container, with failures far from the cause.
