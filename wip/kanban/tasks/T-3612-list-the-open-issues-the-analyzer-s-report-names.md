---
id: T-3612
type: task
nature: improvement
title: List the open issues the analyzer's report names
status: backlog
parent: S-0850
owner: alex
created: 2026-11-23T06:00:22Z
updated: 2026-11-23T06:00:22Z
transitions: []
stream: S-0850
tags: []
after: [T-3611]
---
# T-3612 List the open issues the analyzer's report names

## Work

Run `scripts/flai.sh issue list --report latest --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh issue list --report latest --json` exits 0 in the main checkout.

## Notes
