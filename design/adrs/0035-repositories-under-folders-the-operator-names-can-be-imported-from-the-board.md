---
id: ADR-0035
title: Repositories under folders the operator names can be imported from the board
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: []
refines: [ADR-0029]
---

# ADR-0035 Repositories under folders the operator names can be imported from the board

## Context

The operator asked (S-0098) to import a repository into system-flow from the board: pick a repository that is not a project yet, be asked whether to import it, and have the import run its tests and commit.

The board can only show what `flai serve` connects for ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). `flai serve` connects only for registered projects, and registering needs a manifest key. The dashboard's container holds no project files ([ADR-0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)), so it cannot look for repositories itself. And the channel's methods never take a path from the dashboard: each acts on the project its connection names.

## Decision

**The operator names folders on the host, and `flai serve` offers the board the repositories in them.** `flai serve import add <folder>` keeps the folder in the host's configuration (`import_roots`). `flai serve` looks through each one, three levels deep, every 30 seconds for git repositories with no `system-flow.yaml`. It opens a connection for each to the dashboards its projects already reach.

- **A new kind of connection: a candidate.** Its hello says `kind: candidate`. It names the repository by `import-` and a key made from the folder's name, unique among what is served, and it offers only `import.preview` (`flai import --dry-run`) and `import.run` (`flai import --yes --commit`). The dashboard lists it as not imported, never counts its missing methods, and never treats it as the one project.
- **Naming the folder is the operator's consent.** Importing writes to the repository, runs its tests (its code, as the operator, on the host), and commits. No host action gates `import.run`: the operator said yes by naming the folder, in a shell on the host, and nothing a dashboard can send names a folder or a path. Every import is journalled (`action: import`).
- **`flai import --commit` is the import.** It refuses a repository that is not git, or has uncommitted changes, so the commit holds the import alone. It runs the checks the host names (`checks.commands`) or else the tests the repository has (its own Makefile's `test` target, `go test`, the package manager's test script, `cargo test`, pytest). It commits exactly the paths the import wrote or moved when every test passes or none are found. When a test fails it leaves the import uncommitted, says which failed and what it printed, and exits 5 (`NotCommitted` over the channel).
- **Once applied, it is a project.** Committed or not, `flai serve` registers the repository with its new manifest's key and serves it from then on, and the candidate goes. A declined import is not remembered: picking the repository again asks again.

## Consequences

- A repository can join the board without anyone running `flai import` or `flai dashboard` in it by hand, once its folder is named.
- `flai serve` scans the named folders and holds a connection per candidate per dashboard: a few directory reads every 30 seconds, and an idle WebSocket per repository offered.
- Candidates reach only dashboards that some served project already reaches; with no project served, none is offered.
- Whoever can open the board can import any repository under a named folder, running its tests on the host as the operator. The operator guide says so where it says how to name a folder.
- The channel's hello has an optional field (`kind`); an older dashboard ignores it and would show a candidate as a project it cannot load, so a flai that offers candidates wants a dashboard from this release on.

## Alternatives considered

- **Register each repository by hand** (`flai serve add <dir>`). Nothing is scanned, but it is a host step per repository, and the story asks for the board to find them.
- **Type a path in the board, checked against an allowlist.** One fewer scan, but it would be the first method taking a path from the dashboard, which ADR-0029 kept out.
- **Gate `import.run` on a host action as well.** Consent twice for the same thing; naming the folder already is the operator's say, in the same shell a host action would be enabled from.
- **Commit even when tests fail.** The operator chose not to: a failed import is left uncommitted for them to look at.
