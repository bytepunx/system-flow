---
id: T-0301
type: task
nature: feature
title: The agent action and its command live in the host's configuration, off and empty by default
status: done
parent: S-0079
owner: alex
created: 2026-09-20T15:47:55Z
updated: 2026-09-20T15:53:48Z
transitions:
  - to: ready
    at: 2026-09-20T15:47:56Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:47:57Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:53:48Z
    by: system-flow
stream: S-0079
tags: []
---
# T-0301 The agent action and its command live in the host's configuration, off and empty by default

## Work
A host action named agent, enabled per project like push. The command is an argument list in the host's configuration with no default; flai serve agent set, show, and clear manage it. The tokens {story} and {root} inside an argument are replaced; nothing else is interpreted and no shell is involved.

## Done when
- Command tests: set, show, clear, enable; an empty command cannot be enabled into doing anything
- The Go tests and lint pass

## Notes
