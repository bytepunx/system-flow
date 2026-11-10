---
id: ADR-0281
title: "A strategic agent's container mounts the main checkout's .git/config read-only over the read-write checkout, and flai serve refuses a strategic container start whose checkout has no .git/config file"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0247]
---

# ADR-0281 A strategic agent's container mounts the main checkout's .git/config read-only over the read-write checkout, and flai serve refuses a strategic container start whose checkout has no .git/config file

## Context

[ADR-0247](0247-a-strategic-agent-s-container-mounts-the-main-checkout-read-write-at-its-host.md) runs the planner, the orchestrator, the analyzer, and the curator in a container that mounts the project's main checkout read-write at its host path, with `.git/hooks` read-only over it. It kept `.git/config` writable for one reason: the orchestrator's `flai accept` deletes a merged branch's section from it.

`.git/config` can name commands git runs, such as `core.fsmonitor`, `core.sshCommand`, `core.hooksPath`, a `filter.<driver>.clean`, or an `include.path` to a file that names one. git on the host runs them at its next command in the checkout. `flai guard` refuses each role `git config`, but it does not see a write through the shell. [I-0257](../issues/I-0257-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md) records this as a way back out of the container, found twice: by S-0509 and by the analysis of 2026-11-07.

The reason ADR-0247 gave is gone. [ADR-0278](0278-the-orchestrator-accepts-through-an-mcp-tool-that-runs-flai-accept-on-the-host.md) runs the orchestrator's acceptance on the host, through the MCP tool `accept`, and the shell's `flai accept` refuses in a role's container (`orchestratorInRoleContainer` in `flai/cmd/verify.go`). Nothing else a role runs there writes `.git/config`. `flai subproject --autocommit` and `flai formula run --autocommit` only add and commit. The other branch deletions, in `flai/cmd/branch.go`, `flai/cmd/accept_remote.go`, and `flai/internal/mergequeue/batch.go`, run on the host. No flai code runs `git config` to write.

I-0257 names two remedies: mount `.git/config` read-only, or have `flai serve` check `.git/config` for command-naming keys after each strategic run.

## Decision

A strategic agent's container mounts the main checkout's `.git/config` read-only over the read-write checkout, and flai serve refuses a strategic container start whose checkout has no `.git/config` file.

1. **The mount.** `roleStart` adds `-v <root>/.git/config:<root>/.git/config:ro` after the `.git/hooks` mount, at its host path, so the strategic container's mounts are: the main checkout read-write, its `.git/hooks` read-only, its `.git/config` read-only, the role's cache folder, and the login read-only. A write to `.git/config` in the container fails, whether git writes it in place or renames `config.lock` over it, since a bind mount cannot be replaced.
2. **The file must be there.** Docker makes a missing bind-mount source a folder, owned by root, which would break every git command on the host. `flai serve` checks before it starts the container that `<root>/.git/config` is a regular file, beside its check of `.git` and its making of `.git/hooks`, and refuses the start otherwise, naming the file. A git checkout always has one.
3. **The run's record.** `flai serve` records a container run's mounts from its arguments (S-0665), so the new mount is recorded with the rest.

ADR-0247's other mounts, its guard, and its risks but `.git/config` stand.

## Consequences

- A role in its container can no longer name a command git runs on the host through `.git/config`. I-0257 closes.
- A git command in the container that writes the config fails: `git config` with a write, `git branch --set-upstream-to`, `git remote add`, or `git worktree add --relative-paths`, which sets an extension. No role needs one, and the guard already refuses them.
- The checkout is still the role's to write through the shell, `.git` among it but for its hooks and its config, as ADR-0247 says. A write to `.git/info/attributes` names no command without a filter in the config. A story's worktree under `.flai-cache` stays writable from the container; that is the story half of the boundary, outside this decision.

## Alternatives considered

- **Check `.git/config` after each strategic run.** `flai serve` would compare the file with what it held before the run, or look for command-naming keys, and raise what changed. It finds a write only after it is made, and git on the host may already have run the command between the write and the run's end, as `flai serve` runs git in the checkout throughout. Not built; the read-only mount makes the write impossible instead.
- **Mount `.git` read-only as a whole.** The orchestrator's `flai subproject --autocommit` and `flai formula run --autocommit` commit on main from the container, which needs the index, the objects, and the refs writable. Rejected, as ADR-0247 rejected it.
- **Keep it, as no worse than a host run.** On the host, a role's shell can already write `.git/config`. But the container exists to narrow what a role reaches, and this is the plainest way back to the host left after `.git/hooks`.
