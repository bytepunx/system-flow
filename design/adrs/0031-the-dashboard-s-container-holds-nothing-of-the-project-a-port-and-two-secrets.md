---
id: ADR-0031
title: "The dashboard's container holds nothing of the project: a port and two secrets"
status: accepted
date: 2026-09-20
supersedes: [ADR-0022, ADR-0026, ADR-0027]
superseded_by: []
refines: [ADR-0018, ADR-0029]
---

# ADR-0031 The dashboard's container holds nothing of the project: a port and two secrets

## Context

The dashboard's container was given the project three ways, each with a decision behind it. [ADR-0022](0022-repository-mounted-at-its-host-path.md) mounted the repository read-write at its host path, so that git's links to story worktrees resolved inside the container. [ADR-0026](0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md) let the operator hand it an SSH key to push acceptances with. [ADR-0027](0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md) mounted `.git/hooks`, `.git/config`, and `.git/info` read-only over that, after I-0022 showed that a compromised container could leave a hook that ran on the host as the operator, and it said plainly what stayed open: tracked files, refs, and ignored files the host executes.

[ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md) decided that the dashboard reaches a project only by asking flai on the host. Reads moved in S-0073 and S-0074, writes in S-0075, which also took flai, git, and ssh out of the image, and MCP left the dashboard in S-0076 ([ADR-0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)). Since S-0075 the mount has been passed and not used: with it taken away by hand, every page and every write worked.

## Decision

**The container holds nothing of the project.** `flai dashboard` gives it a published port and two secrets, each a single file mounted read-only under `/run/secrets`: the login token and the credential `flai serve` proves itself with. It passes no project volume and no `PROJECT_DIR`; no read-only git paths; no git identity; no global excludes file; no push key, known hosts, or generated passwd. The image names no project directory. `flai dashboard` no longer audits the clone for hooks and settings, because nothing it starts can write them.

**The push key is retired.** `--push-key`, `--push-known-hosts`, `dashboard.push_key`, and `dashboard.push_known_hosts` are still accepted, so an old configuration loads and can be cleared; when one is set, `flai dashboard` names it, says it is ignored, and says what to do instead. An acceptance made from the board is committed and tagged by `flai serve` on the host as the operator and pushed from the host with the operator's own credentials: `flai push --pending` today, and by `flai serve` on the operator's word when S-0078 lands. The standing "accepted, not pushed" state on the board, the item page, `flai board`, and the agents' `inbox` is kept from ADR-0026 unchanged.

**The container still runs as the user who started it.** The two secrets are files with mode 0600, and a container running as anyone else could not read them. Handing them over as environment variables instead would put them in `docker inspect` and in the process's environment, which ADR-0018 rules out. Running as that user costs nothing now: there is no file of the host within the container's reach to write as them.

**`flai dashboard status` says when a container is stale.** A container started by an older flai keeps its mounts until it is restarted. `status` reads the running container's mounts and, when any lies outside `/run/secrets`, says so and how to restart. It no longer reports a push key or read-only git paths.

**What of [ADR-0018](0018-dashboard-token.md) still holds: all of it that concerns the dashboard.** The token is generated on first run, kept at `.flai-cache/dashboard.token` with mode 0600, mounted read-only, named by `FLAIOVER_TOKEN_FILE`, never an environment variable, an argument, a log field, or a committed file; the login link, the cookie, the bearer check, rotation, and plain HTTP behind the operator's TLS are unchanged. One sentence changed with ADR-0030: agents no longer present this token, because MCP is flai's and has a token of its own.

**`worktrees.relative_paths` stays** as the opt-in ADR-0022 made it. Nothing about the dashboard needs it any more; it remains for operators who move or share clones.

## Consequences

- I-0022 and what ADR-0027 left open are closed by removal, not by another guard. A process in the container cannot write a hook, a git setting, a tracked file, a ref, or an ignored file the host executes, because it can reach no file of the project at all. What it can still do is what the channel offers: ask flai for the named, typed methods of ADR-0029, each of which flai validates and performs itself.
- What a compromised container can still reach: the dashboard token and the agent credential (both readable by it, as they must be), the channel's methods with whatever they allow (moves, saves under the three folders, an acceptance of a story in review), and the network. Host actions beyond those are off unless the operator enables them (ADR-0029).
- The container cannot pin a stale `.git/config`, and the operators' advice to restart the dashboard after changing git settings is gone: commits are made on the host with the configuration of the moment.
- A Windows or otherwise unmirrorable host path no longer costs anything: acceptance of a story with a branch works wherever `flai serve` runs, because git runs there.
- An operator who had a push key configured loses unattended pushes from the board until S-0078, and is told so at every start until they clear the setting. The key itself and its deploy-key registration are theirs to remove.
- `.flai-cache/dashboard.known_hosts` and `.flai-cache/dashboard.passwd`, written by earlier releases, are no longer read or written; they can be deleted.
- A dashboard started before this release keeps the old mounts until restarted; `flai dashboard status` says so.

## Alternatives considered

- **Keep the mount read-only, as a fallback for reads.** Nothing reads it; it would keep every project file visible to a compromised container for no function.
- **Run the container as the image's own unprivileged user and pass the secrets as environment variables.** Removes the last tie to the host user, but puts both secrets where ADR-0018 says they must not be.
- **Make the secret files world-readable so any user can read them.** Trades a harmless `--user` for secrets any local user can read.
- **Refuse to start when a push key is configured.** Safer against silent loss of unattended pushes, but it would break every such operator's dashboard on upgrade for a setting that can simply be ignored and explained.
- **Remove the retired flags and config keys outright.** An unknown flag or key is an error; the operator would meet a failure instead of an explanation.
