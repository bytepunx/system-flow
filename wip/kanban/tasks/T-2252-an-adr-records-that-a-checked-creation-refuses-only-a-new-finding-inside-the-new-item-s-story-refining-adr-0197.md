---
id: T-2252
type: task
nature: remediation
title: An ADR records that a checked creation refuses only a new finding inside the new item's story, refining ADR-0197
status: backlog
parent: S-0525
owner: alex
created: 2026-10-28T14:52:37Z
updated: 2026-10-28T14:52:37Z
transitions: []
stream: S-0525
tags: [flai, adr]
touches: [design/adrs/README.md]
---
# T-2252 An ADR records that a checked creation refuses only a new finding inside the new item's story, refining ADR-0197

## Work

Write a new ADR with `flai adr new` that refines ADR-0197 (`--refines ADR-0197`), whose point 6 keeps the unscoped diff for `flai story new`, `flai task new`, and the other checked writers. The decision: a checked creation of a work item (`itemnew.Create`, reached by `story new`, `task new`, and `epic new` with `--body-stdin`, `--after`, or `--autocommit`, by the MCP tool `item_new` with a body or `after`, and by `flai issue story`) scopes its second check run to the new item's story as a checked edit does, refuses only a new finding inside that scope or an error on the item's own file, returns a new finding outside it as a warning, and keeps `wip.overlap`, `touches.folder`, and `wip.after-hub` advisory. A document save and `flai adr new` keep their unscoped diff.

The context names I-0191: every work item write takes the items lock under `.flai-cache/` (ADR-0184), so the first checked creation in a project with no `.gitignore` made the folder and refused itself on `layout.gitignore`, a finding on the repository root that the creation did not bring. The alternatives considered: the lock outside the project, which supersedes ADR-0184 and needs a per-project key under the container's per-story `FLAI_CACHE_DIR`; flai writing the `.gitignore` entry, which edits a tracked file the operator owns and does nothing for a project without git. Give the ADR the topic `cli`, as ADR-0197 has, and accept it with `flai adr accept`. `flai adr new` numbers it and adds its row to `design/adrs/README.md`, which this task's touches name; the ADR's own file widens the touches when the task closes.

This task runs with the fix, not after it: both rest on the same decision, and the ADR's words do not depend on the code.

## Done when

- A new accepted ADR under `design/adrs/` refines ADR-0197, states the decision above, names I-0191 in its context, and lists the two alternatives with why they were not taken.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean on the ADR and the index.

## Notes

The story's touches keep the folder `design/adrs` for the ADR's file, since its number and slug are `flai adr new`'s to give.
