---
id: T-2008
type: task
nature: remediation
title: An item edit refused by the check puts back only its own change to an item file, keeping what was written meanwhile
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:13:00Z
updated: 2026-10-21T12:13:31Z
transitions: []
stream: S-0479
tags: []
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go]
after: [T-2007]
---
# T-2008 An item edit refused by the check puts back only its own change to an item file, keeping what was written meanwhile

## Work

`itemedit.Apply` keeps each file it changes and, when the check after the write refuses the edit or a later step fails, writes every file back as it was, which undoes a move or a usage write made to an item while the check ran. Waits for T-2007, whose merge it uses. For an item file it wrote, put back only its own change: through the item lock, a file still as the edit wrote it is written back whole, and one changed since is merged as `Repo.Save` merges, with the edit's write as the base and the file as it was before as the copy, so the keys the other writer changed stay. A file the edit created, renamed, or that is not an item is put back as before.

## Done when

- A test in `flai/internal/itemedit/lock_test.go` notes an item, saves the edit's change, saves another key to the item from a second copy, restores, and finds the edit undone and the other key kept; it fails on main.
- `flai test` on the changed paths passes.

## Notes
