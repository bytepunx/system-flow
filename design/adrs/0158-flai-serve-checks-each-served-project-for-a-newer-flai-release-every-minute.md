---
id: ADR-0158
title: "flai serve checks each served project for a newer flai release every minute, upgrades the host to it behind the shell-only auto-upgrade host action, and otherwise asks the operator on one thread per project and release"
status: proposed
date: 2026-10-15
supersedes: []
superseded_by: []
refines: [ADR-0040, ADR-0142]
---

# ADR-0158 flai serve checks each served project for a newer flai release every minute, upgrades the host to it behind the shell-only auto-upgrade host action, and otherwise asks the operator on one thread per project and release

## Context

[I-0141](../issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md) records four occurrences of one cause: the host's installed flai stays behind main after a flai release is tagged. The MCP tool `test` of flai 1.40.1 refused the manifest's `{importers}` placeholder, which S-0364 added, and told the agent to fix the manifest. Three times the MCP tool `task_done` of flai 1.48.0 refused under [ADR-0142](0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md), since main's history held 1.49.x. Each agent worked around it with the tree's `scripts/flai.sh`, which a project made from the template does not have.

Nothing upgrades the host. `flai serve` logs `flai is older than the project it serves` once, when it starts serving a project. A host that keeps running through a release does not say so again. The inbox carries `flai_outdated`, but only agents read it.

`flai host upgrade` already does the upgrade: it runs `flai self-upgrade`, which refuses a release that is not published or whose signature does not verify ([ADR-0070](0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md)). Then the host restarts on the new binary and starts its children again, `flai serve` and each project's MCP server ([ADR-0040](0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)).

## Decision

`flai serve` checks every minute whether its flai is behind each project it serves. With the `auto-upgrade` host action on for that project, it upgrades the host to the newer release. With it off, it asks the operator on one thread.

1. **The check.** On its minute loop, and when it starts serving a project, `flai serve` calls `release.FlaiOutdated` for each served project. It compares its own version with the newest `flai/vX.Y.Z` tag reachable from the project's HEAD. A development build is never behind, as `FlaiOutdated` already holds.
2. **The action.** `auto-upgrade` is a host action like `auto-approve`. The operator turns it on only in a shell, with `flai serve enable auto-upgrade [root]`; the dashboard cannot ([ADR-0039](0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md)). It is off by default, because it replaces the binary on the operator's machine.
3. **The upgrade.** With `auto-upgrade` on for a project that is behind, `flai serve` asks the host to upgrade to that project's newest release, by version, as `flai host upgrade --version X.Y.Z` does. Where several such projects are behind, it asks for the highest version. The host verifies and installs the release and restarts, as for the operator's upgrade. `flai serve` journals the upgrade as a host action and logs `flai upgrade asked` (info) with the running and newest versions.
4. **What the restart keeps.** `flai serve` starts story agents, planners, the orchestrator, and the analyzer detached, in sessions of their own. A host restart does not signal them. The new `flai serve` adopts the live ones, as it adopts any agent that outlived an earlier serve. Their next `flai mcp` runs the new flai. The dashboard's connection drops and dials again.
5. **What it waits for.** A dashboard write that `flai serve` runs in its own process, such as an acceptance, is cut by a restart. `flai serve` asks for no upgrade while such a write is running, and asks again on the next minute. It does not wait for agents, which survive the restart.
6. **A refused upgrade.** A release is tagged before its binaries are published, so the first ask can be refused as not published. `flai serve` logs the refusal (warn) and asks again for the same release no sooner than 15 minutes later. When the release is still not installed an hour after the first ask, it opens the thread of point 7, naming the last refusal, and keeps asking every 15 minutes.
7. **The thread.** With `auto-upgrade` off, `flai serve` opens one thread on the project's `system-flow.yaml`, as author `flai serve`. It names the running and newest versions and asks the operator to run `flai host upgrade`, or to turn on `auto-upgrade`. It opens one thread per project and newest release. A record kept under the serve directory holds the release asked about, so neither the minute loop nor a restart of `flai serve` opens it again. A newer release opens a new thread.
8. **The tier refusal names the version.** When the manifest's test tiers do not validate and `FlaiOutdated` answers, `flai test`, `flai verify`, and the MCP tools `test` and `verify` say so. The refusal says the running flai is older than the newest release in the project's history, so the manifest may use what it does not know, and names the upgrade command. It no longer tells the agent only to fix the manifest. When `FlaiOutdated` answers nothing, the refusal stays as it is.

[ADR-0142](0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md)'s refusal stays. This decision shortens the time it holds.

## Consequences

- With `auto-upgrade` on, a host catches up within about a minute of a release being published, often within 15 minutes of its tag. Then ADR-0142 stops no one for long.
- With it off, the operator hears of each release that leaves the host behind, once per project, on a thread the board shows.
- An upgrade restarts the MCP servers too. An agent mid-call on one gets an error for that call and calls again.
- An unpublished or unverifiable release is asked for at most four times an hour, never in a tight loop.

## Alternatives considered

- **Raise `flai.minimum` when the manifest's test vocabulary changes.** A host below the minimum stops serving the project at all, which is worse than either kind of refusal.
- **Upgrade to the newest published release, not the project's tag.** An unpublished tag would make the host answer "already the latest" every minute, and the host would upgrade past what any project asked for.
- **Upgrade without waiting for writes.** An acceptance cut halfway leaves a merge unfinished in the main checkout.
- **Only make the refusal clearer.** The agent would know why, but the host would still lag until the operator noticed.
