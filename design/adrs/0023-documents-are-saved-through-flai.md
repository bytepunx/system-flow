---
id: ADR-0023
title: Documents edited in the dashboard are saved through flai
status: accepted
date: 2026-09-19
supersedes: []
superseded_by: []
---

# ADR-0023 Documents edited in the dashboard are saved through flai

## Context

The designer reads every document in flaiover and, until now, changed them in an editor on the host. E-0006 makes the dashboard the place where the designer edits design, docs, and the bodies of work items (S-0040). [ADR-0016](0016-dashboard-delegates-to-flai.md) already says the dashboard does not write the repository itself: it runs the bundled `flai`, so a change made in the browser is the change the command line would have made, and the rules live in one place.

Editing raises four things a move or a thread reply did not. Some files are partly or wholly owned by flai: the front matter of a work item is its state, and `wip/agents/index.md` is generated. Two people can have the same document open, and agents write `wip/` all the time. A saved document must still pass `flai check`, or the next agent starts from a broken repository. And the designer's edits are on `main`, where `flai stream sync` brings them to every story branch ([ADR-0019](0019-story-branches-and-touches.md)), so an edit that is saved and not committed is an edit agents do not see.

## Decision

A document edited in the dashboard is saved by `flai doc save`, and the editor learns what it may edit from `flai doc show`. The dashboard's server passes the content through and holds none of the rules.

- **What may be edited.** A path must be a markdown file under the design, docs, or wip folders, the same guard as the MCP `doc_get` tool. Each file has a mode. `full`: design and docs files, body and front matter, the front matter validated by the check. `body`: work items, narratives, and `board.md`, whose front matter flai owns; a save that changes it is refused. `none`: generated files (`wip/agents/index.md`, `design/issues/summary.md`), threads, issue files, and everything under `wip/archive`, with the reason given.
- **Concurrent edits.** `flai doc show` returns the SHA-256 of the file's content. `flai doc save` takes the hash the editor loaded and, when the file no longer has it, fails with a conflict that carries the current content, its hash, and a unified diff between the current file and what was submitted. Saving over the other change is the designer's explicit second save with the current hash.
- **Validation.** `flai check` runs before and after the new content is written. Any finding the edit introduces, at any level and on any path, restores the old content and refuses the save with the findings: the repository's gate is the strict check, and every document rule is a warning. An error already on the file refuses too, until the edit fixes it. Warnings that were already on the file are reported and do not refuse.
- **Dates.** For a design or docs file whose front matter has `updated`, flai sets it to today when the designer did not change it themselves.
- **Commit.** Unless the project sets `dashboard.autocommit: false` in `system-flow.yaml`, the save commits that one path on the main checkout's current branch: `docs:` for design and docs, `chore:` for wip, with the item's ID in brackets when the file is a work item or a narrative. The author is the git identity in the environment, which `flai dashboard` passes in from the host, and a trailer names the dashboard. Nothing is pushed; that question belongs to S-0052.

Creating, renaming, and deleting documents are not part of this decision.

## Consequences

- One implementation of the ownership rules, reachable from the command line and from MCP later, and tested in Go with real git.
- A save costs a process start and the check, which on this repository is well under a second; an editor that saved on every keystroke would not be acceptable, and the editor saves on demand.
- The main checkout gains small commits authored by the designer between story merges. They are path-limited, so uncommitted `wip/` changes made by agents are never swept into them, and `flai accept` rebases story branches onto them as it already does for any commit on `main`.
- A save is refused for a file that already had an error finding before the edit, unless the edit fixes it. That is deliberate: the editor is a way to fix such a file. A file with an old warning can still be edited.
- A save costs two runs of the check.
- With autocommit off, the designer's edits reach story branches only when someone commits them, and an acceptance from the board will list them as uncommitted changes (S-0051).
- Accepted ADRs are immutable by convention. This decision does not enforce that in the editor; it is raised as an open question in S-0040.

## Alternatives considered

- The dashboard's server writes the file and runs `flai check` afterwards: two implementations of what may be edited, and the first write the dashboard would make without flai.
- Lock a document while it is open: locks outlive browser tabs, and agents writing files do not take them.
- Last write wins: silently loses an agent's or a colleague's change to the same document.
- Validate a temporary copy instead of writing and restoring: the check reads the repository as a whole (parents, indexes, overlaps), so a copy outside it does not get the same answer.
- Save to a branch and open a pull request: the designer's intent would not reach story branches until someone merged it, which is the opposite of what ADR-0019 set up.
