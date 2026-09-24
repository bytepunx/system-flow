---
id: T-0400
type: task
nature: feature
title: Generate the command reference from the cobra tree
status: done
parent: S-0017
owner: alex
created: 2026-09-24T08:12:55Z
updated: 2026-09-24T08:16:02Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T08:13:04Z
    by: system-flow
  - to: done
    at: 2026-09-24T08:16:02Z
    by: system-flow
stream: S-0017
tags: [cli]
touches: [flai/cmd, scripts, Makefile]
---
# T-0400 Generate the command reference from the cobra tree

## Work
Add a hidden `flai reference` command that walks the cobra command tree and prints markdown: every visible command with its usage, description, examples, aliases, and flags (local and inherited), in tree order. Standard library only, no cobra/doc dependency. Add `scripts/flai-reference.sh` and a `make flai-reference` target that write `docs/users/flai-reference.md`, and a test that fails when the committed file is stale.

## Done when
`make flai-reference` writes the file, the renderer has a behavior test, the staleness test runs in the integration tier, and lint is clean.

## Notes
