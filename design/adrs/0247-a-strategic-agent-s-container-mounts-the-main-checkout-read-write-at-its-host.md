---
id: ADR-0247
title: "A strategic agent's container mounts the main checkout read-write at its host path with its git hooks read-only, and the role's cache folder, and its writes are the ones its guard allows on the host"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0203, ADR-0212]
---

# ADR-0247 A strategic agent's container mounts the main checkout read-write at its host path with its git hooks read-only, and the role's cache folder, and its writes are the ones its guard allows on the host

## Context

[ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) runs a story's agent in a container that mounts the story's worktree, the story's cache folder, and the operator's login read-only, and nothing else: no home, no `~/.ssh`, no `~/.flai`, no main checkout, no `wip/`. S-0508 built that adapter ([ADR-0212](0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md)), and it refuses a start that names no story.

The planner, the orchestrator, the analyzer, and the curator work in the main checkout, not in a worktree. `flai serve` starts each there (`plan.go`, `orchestrate.go`, `analyze.go`, `curate.go`), and each writes in its own way:

| Role | What it writes, and how |
|------|-------------------------|
| Planner | Work items and threads, through the MCP tools and the `flai` CLI, which write the main checkout's `wip/` |
| Orchestrator | Work items, threads, and the board through MCP; `flai accept --by orchestrator`, which merges, archives, and commits on main; `flai subproject ... --autocommit` and `flai formula run --autocommit`, which commit on main |
| Analyzer | Its report and the index under `design/analysis/` with the file tools; issues with `flai issue new` and `bump` |
| Curator | Digest entries under `design/digest/` with `flai digest new`, `edit`, `absorb`, `cited`, and `compact`, left uncommitted |

`flai guard` bounds each role's calls on the host, from its role file ([ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md)). It sees the file tools, the MCP tools, and the `flai` and git commands; other shell commands pass, as a sub-agent's do ([ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)).

S-0509 must say what a strategic agent's container mounts. Two choices were weighed:

1. **The main checkout read-write at its host path.** The role writes as it does on the host, and the guard bounds it as on the host.
2. **The main checkout read-only, with every write over the host's HTTP MCP server.** The container could change nothing of the project itself. But the analyzer's report, the curator's `flai digest` commands, and the orchestrator's `flai accept`, `flai subproject`, and `flai formula run` have no MCP tool today. Each would need one before its role could run in a container.

## Decision

A strategic agent's container mounts the main checkout read-write at its host path with its git hooks read-only, and the role's cache folder, and its writes are the ones its guard allows on the host.

1. **Mounts.** When the start request names a role and no story, the container adapter mounts these, and nothing else:
   - the project's main checkout, read-write at its host path, as the working directory, its `.git` inside it;
   - the checkout's `.git/hooks` read-only over it, which `flai serve` makes when it is missing, since a hook written there would run on the host at its next git command;
   - the role's cache folder, `.flai-cache/agents/<role>` with `<role>` the run's role (`plan`, `orchestrate`, `analyze`, or `curate`), read-write at its host path, for the inner harness's configuration and sessions and the container's `flai`'s configuration and cache, as a story's cache folder is;
   - the operator's login read-only at `/run/claude-login`, for an inner harness that needs it.

   It still never mounts the home, `~/.ssh`, `~/.flai`, or the Docker socket, so no credential of the operator's but the login reaches the container.
2. **The container.** Its name is `flai-agent-<role>-<session>`, by which `flai serve` stops it, as it stops a story's. The variables passed by name are a story's but `FLAI_STORY`, with `FLAI_ROLE` added. The orchestrator's container lives for its run, which ends once it has acted on its reasons ([ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md)); a stop of the `orchestrate` host action stops it by name.
3. **pilot's start request** carries the run's `Role`, `Item`, `Focus`, `Since`, and `AskedOn`, as the host's start does, so the inner harness starts with `--agent <role>` and `FLAI_ROLE`. pilot reads the role file a built-in role runs from out of the mounted checkout, the project's or flai's baseline, since the request does not carry it. A role a project defines in a role file of its own is refused on the container harness until a story carries its file to pilot.
4. **The guard** runs inside the container as the inner harness's hook from the mounted checkout's `.claude/settings.json`. It reads `FLAI_ROLE` from the environment pilot gives the inner harness and the role file from the mounted checkout, exactly as on the host. What it refuses each role is unchanged:

   | Role | The guard refuses inside the container, as on the host |
   |------|--------------------------------------------------------|
   | Planner | Every file edit; any write to an item it was not started for, or outside planning; git writes |
   | Orchestrator | Each act whose `orchestration.permissions` entry is off; `flai release` but `--evaluate`, `flai push`, `git push`, and `git tag`; every file edit |
   | Analyzer | Every file edit outside `design/analysis/`; any write to a work item; git writes |
   | Curator | Every file edit; `flai digest seed`; any write to a work item; git writes |

5. **The claude check.** `flai serve` reads Claude Code's version from the image before a strategic container starts, as it does for a story's (S-0508), and does not probe the host's `claude` for it.

## Consequences

- A strategic agent on the container harness sees the project and one cache folder, not the machine: no home, no `~/.ssh`, no `~/.flai`, no Docker socket. That is the boundary ADR-0203 draws for a story's agent, widened to the main checkout for these four roles alone.
- Nothing a role does today needs a new MCP tool. The analyzer, the curator, and the orchestrator write as they do on the host.
- The activity documents and the journal are unchanged. `flai serve` reads the container's log through the inner harness's reader, which pilot relays unchanged, and logs each run's activity as it does a host run's.
- `.git/hooks` read-only closes the plainest way back to the host. `.git/config` stays writable, since `flai accept` deletes a merged branch's section from it.

Risks:

- **The checkout is the role's to write.** A shell command the guard does not see can write any file of the checkout, `.git` and every story's worktree under `.flai-cache` among them, as it can on the host today. The container narrows the rest of the machine away, not the checkout.
- **`.git/config` can name commands** git runs on the host, such as `core.fsmonitor` or `core.sshCommand`. The guard refuses git writes to every role but the orchestrator's acceptance; a write through the shell is not seen.
- **The container's `flai` has no host configuration.** Its `FLAI_CONFIG` is in the cache folder, so a command that reads the host's actions, such as `flai accept` while the merge queue batches, or `flai verify`, which needs the test tiers' tools, sees none of them there. The orchestrator's acceptances are the role most exposed: each limitation met is recorded as an issue.
- **This repository's guard hook runs `scripts/flai.sh`,** which builds flai with `go`, absent from the image (S-0520). A project made from the template runs the image's `flai guard`.

## Alternatives considered

- **The main checkout read-only, with writes over MCP.** Rejected for now: four of the roles' writes have no MCP route, and the guard already bounds the writes they make. Taken up again once those routes exist, it would leave the container nothing of the project's to write.
- **The main checkout read-only, with the role's own folders read-write.** `design/analysis/` for the analyzer and `design/digest/` for the curator would do for those two. The planner and the orchestrator write `wip/`, `system-flow.yaml`, and `.git` through the CLI, so their mounts would be the checkout anyway, and two mount sets would differ by role for little gain.
- **`.git` read-only as a whole.** It would refuse the orchestrator's `flai accept`, which commits and merges on main.
- **A clone of main per run.** It is the route to Kubernetes ADR-0203 keeps open, but every write would then need pushing back, and the curator's entries are left uncommitted by design.
