---
id: ADR-0274
title: "The close-out of a story whose agent runs in a container runs its verify in that container, and its commit and sync run on the host through flai's own code with the main checkout's hooks"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0229, ADR-0246, ADR-0253]
---

# ADR-0274 The close-out of a story whose agent runs in a container runs its verify in that container, and its commit and sync run on the host through flai's own code with the main checkout's hooks

## Context

A story's agent on the `container` harness runs in a container that mounts the story's worktree read-write ([ADR-0246](0246-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-git-read-only.md), [ADR-0253](0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md)). Everything in that worktree is the agent's to write: `scripts/`, `system-flow.yaml`'s copy on the branch, `.claude/settings.json`, the `.git` file, and the story's registration under the main checkout's `.git/worktrees/<story>`, its `config.worktree` among them.

Two things then run that code on the host, as the operator, with nothing reviewed.

| What | Where | What it runs on the host |
|------|-------|--------------------------|
| The close-out job ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)) | `flai/internal/serve/verifyjob.go`, `command` and `runCloseOut` | `sh scripts/close-out.sh <story>` from the worktree, and through it `scripts/flai.sh`, the test tiers' commands, and git with the hooks and configuration the worktree's registration names |
| The operator's checks (S-0082) | `flai/internal/serve/checks.go`, `RunChecksWith` and `runOneCheck` | each named check's command, in the worktree |

The job and the checks run them whatever the story's harness. I-0295 records it. ADR-0253 says the close-out runs inside the container, against the main checkout's state. The guides do not say that the host runs the worktree's scripts, and `docs/users/flai.md` says the close-out cannot run in the container (S-0520). So the container harness reads as a boundary it is not.

The close-out has two halves. `flai verify <story> --record-issues` runs the check and the tiers: it reads the worktree, the main checkout's `.git`, `wip/`, and `system-flow.yaml`, and writes only its report under `.flai-cache/verify/`. The container already has exactly those mounts. The commit, `flai verify <story> --sync-only`, and the sync rounds ([ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md)) write the branch, which only the host can: the main checkout's `.git` is read-only inside (ADR-0246).

## Decision

The close-out of a story whose agent runs in a container runs its verify in that container, and its commit and sync run on the host through flai's own code with the main checkout's hooks.

The story's agent is the story's own over the project's default (`Agent.With`, [ADR-0251](0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md)). When its harness is `container`:

1. **The verify runs in the container.** The job runs the image's flai, `flai verify <story> --record-issues`, with `all` and `fresh` as asked, in a container from the story's image: the project's layer, else the host's image for the inner harness, as `containerImageHost` in `flai/internal/serve/container.go` decides for the agent. It has the story's mounts and no others:

   | Mount | Mode |
   |-------|------|
   | the worktree, the working directory | read-write |
   | the main checkout's `.git` | read-only |
   | the story's registration, `.git/worktrees/<story>` | read-write |
   | the main checkout's `wip/` and `system-flow.yaml` | read-only |
   | the main checkout's `.flai-cache/verify` | read-write |

   It runs as the operator's uid, with `safe.directory` for the main checkout, on the host's network. It has no login mount, no MCP variables, no provider key, and no `.flai-cache/agents/<story>` folder. So the worktree's scripts, the tiers they select, and the hooks the registration names run inside.
2. **The commit and the sync run on the host, through flai.** The job commits what is outstanding with git and the request's message, runs this flai's own `flai verify <story> --sync-only`, and, when main moved, syncs and verifies again in the container, at most three rounds in all, as `scripts/close-out.sh` does (ADR-0156). No script from the worktree runs on the host.
3. **The host's git does not trust the worktree's git files.** The container can write the worktree's `.git` file and the registration, `config.worktree` included. So each git command on the host names its git dir and work tree from the main checkout's own record of the worktree, and sets `core.hooksPath` to the main checkout's `.git/hooks`. A `core.hooksPath` or `core.fsmonitor` the container wrote runs nothing on the host.
4. **The operator's checks run in the same kind of container.** Each check's command runs through the same argv, in the same image, with the run's environment passed by name. Cancelling a check stops its container by name, as a stop does for an agent ([ADR-0058](0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md)).
5. **No fallback to the host.** When docker is not on `PATH` or the image is not on this host, the job's verdict is `failed` with the refusal `containerAllowed` gives, which names the fix. It never runs the close-out on the host instead.
6. **Any other harness keeps today's job:** `sh scripts/close-out.sh <story>` in the worktree, or this flai's `flai verify` where the worktree has no script.

What still runs on the host, by design:

| What | Why |
|------|-----|
| The merge queue's batch tiers | They run in the main checkout, on stories an acceptance decision put there |
| An acceptance's `fix_command` | It runs in the main checkout, after the operator's or the orchestrator's decision to accept |
| `flai test` and `flai verify` run by hand | The operator chose to run them |
| A story's agent on a host harness | Its code runs on the host anyway |

## Consequences

- A container story's code reaches the host only through review: an acceptance merges it into the main checkout, where the batch and `fix_command` run it.
- The close-out job for a container story no longer runs `scripts/close-out.sh`. Its commit step and its three sync rounds are flai's code, and a change to the script's steps changes them in Go too.
- The close-out needs the story's image on the host, as the agent's start does. A host whose image is gone fails the close-out, naming the build.
- The tiers inside the container use the layer's tools (S-0510). A tier that needs something only the host has fails in the container, as it would for the agent.
- The host's commit runs the main checkout's hooks, not the branch's. A story that changes a hook does not run it in its own close-out.
- `docs/users/flai.md`, `docs/operators`, and `design/system/flai-cli.md` say what a container story runs inside and what still runs on the host.

## Alternatives considered

- **Document the boundary alone.** It costs nothing in code, but every close-out of a container story would still run the agent's unreviewed scripts on the host as the operator. It leaves I-0295's cause in place and ADR-0253 unmet.
- **Run the whole close-out script inside the container.** The script commits, and a commit inside is refused: the main checkout's `.git` is read-only there (ADR-0246). Making it writable would let the container rewrite any branch and the hooks.
- **Run the close-out on the host with the worktree's scripts checked against main's.** A script can call anything the branch changed, such as `scripts/env.sh` or a tier's own script, so a check of a few files does not bound what runs.
