---
id: T-1715
type: task
nature: feature
title: The design and operator docs say where the status of a story on OpenCode comes from
status: backlog
parent: S-0411
owner: alex
created: 2026-10-16T05:19:51Z
updated: 2026-10-16T05:19:51Z
transitions: []
stream: S-0411
tags: [cli]
touches: [design/system/flai-cli.md, design/system/agent-adapters.md, docs/operators/index.md]
after: [T-1714]
---
# T-1715 The design and operator docs say where the status of a story on OpenCode comes from

## Work

- In `design/system/flai-cli.md`, where `flai serve`'s stream is described, say that it reads through the run's adapter and name OpenCode's source, linking the ADR.
- Record the choice in `agent-adapters.md` § OpenCode, and update the table in `docs/operators/index.md` § Harnesses.
- Waits for T-1714: the docs describe the tested behaviour.

## Done when

The docs say what the story built, and the markdown lint passes.

## Notes
