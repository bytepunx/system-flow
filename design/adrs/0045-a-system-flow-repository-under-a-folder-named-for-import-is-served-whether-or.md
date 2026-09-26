---
id: ADR-0045
title: A system-flow repository under a folder named for import is served whether or not it is registered
status: accepted
date: 2026-09-26
supersedes: []
superseded_by: []
refines: [ADR-0035]
---

# ADR-0045 A system-flow repository under a folder named for import is served whether or not it is registered

## Context

[ADR-0035](0035-repositories-under-folders-the-operator-names-can-be-imported-from-the-board.md) has `flai serve` offer the board the git repositories with no `system-flow.yaml` below the folders the operator names (`import_roots`), and register one once it is imported from the board. A repository imported with `flai import` on the command line got a manifest and was never registered (I-0047, S-0120). `FindCandidates` skips any repository with a manifest, so such a repository was neither offered nor served, and nothing short of `flai dashboard` in it brought it to the board. The story asks that such a repository be served, or offered to be served, and not dropped from both.

## Decision

**`flai serve` serves every git repository with a `system-flow.yaml` below a folder named for import that is not registered.** It does this as it serves the projects below the folder it was started in (S-0102, [ADR-0036](0036-a-folder-that-is-not-a-project-is-served-whole-by-flai-mcp-and-flai-dashboard.md)): found again every 30 seconds, served for the first dashboard, by address, that the registered projects or a recorded one reach, and never written to the registry. A registered project is left to the registry. One whose manifest does not load or has no key, whose key another project has, or for which no dashboard is known is not served, and `flai serve status` lists it with why. Only git repositories count there, as only they are candidates, and nothing inside a repository is looked at.

`flai import` on the command line also registers the project itself when a `flai host` runs for the config, with the entry `flai dashboard` writes, so that it is served before the next scan and stays served if the folder is no longer named.

## Consequences

- A repository imported on the command line reaches the board without another step when a host runs, and, below a named folder, even when none ran at the time.
- Naming a folder for import now also lets anyone who can use the board act on the system-flow projects already in it, as they can on registered ones: read them, and change them through the host actions enabled for them. The operator guide says so where it says how to name a folder.
- Removing a folder from `import_roots` stops serving the projects found there that are not registered.
- No dashboard change: a served project is one the switcher already knows how to show.

## Alternatives considered

- **Offer it in the switcher, to be served on request.** Needs a second kind of candidate in the dashboard: a candidate's connection answers only the import methods, and the switcher would offer to import a repository that is already imported. More work for a step the operator already consented to by naming the folder.
- **Register it in `projects.json` when found.** It would stay served after the folder is no longer named, with nothing saying why.
