---
id: T-2196
type: task
nature: improvement
title: docs/operators/runbooks/agent-container.md describes building the agent images and running a story's agent in a container by hand, and the runbooks README indexes it
status: done
parent: S-0513
owner: alex
created: 2026-10-28T12:28:31Z
updated: 2026-10-28T12:29:19Z
transitions:
  - to: ready
    at: 2026-10-28T12:28:37Z
    by: agent-S-0513
  - to: in-progress
    at: 2026-10-28T12:28:38Z
    by: agent-S-0513
  - to: done
    at: 2026-10-28T12:29:19Z
    by: agent-S-0513
stream: S-0513
tags: []
touches: [docs/operators/runbooks/agent-container.md, docs/operators/runbooks/README.md]
usage:
  source: log
  seconds: 0
  models: []
---
# T-2196 docs/operators/runbooks/agent-container.md describes building the agent images and running a story's agent in a container by hand, and the runbooks README indexes it

## Work

Write `docs/operators/runbooks/agent-container.md` from S-0505's `scripts/agent-image.sh`, `scripts/agent-container.sh`, the `agent-images` Make target, and its results document: what the two images are (local, unsigned, E-0029), building them, what the run mounts and never mounts, what it needs on the host, the settings it reads, the limitations the results record (linked, not copied), and cleaning up. Add its row to `docs/operators/runbooks/README.md`. One task: the two files are one change and wait for nothing.

## Done when

- The runbook covers criteria 1 and 2 of S-0513 and the README indexes it
- `flai check --strict` scoped to S-0513 and the markdown lint pass

## Notes

The scripts and the results document are on `story/S-0505` and reach main when S-0505 is accepted.
