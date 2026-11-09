---
id: T-2824
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0646
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-09T22:22:55Z
transitions: []
stream: S-0646
tags: []
after: [T-2823]
---
# T-2824 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
