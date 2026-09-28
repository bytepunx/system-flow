---
id: T-0496
type: task
nature: feature
title: A story's topics with their sources, printed by flai show and --json
status: done
parent: S-0135
owner: alex
created: 2026-09-28T22:41:23Z
updated: 2026-09-28T22:48:21Z
transitions:
  - to: ready
    at: 2026-09-28T22:41:39Z
    by: agent-S-0135
  - to: in-progress
    at: 2026-09-28T22:45:40Z
    by: agent-S-0135
  - to: done
    at: 2026-09-28T22:48:21Z
    by: agent-S-0135
stream: S-0135
tags: []
touches: [flai/internal/topics, flai/cmd]
---

# T-0496 A story's topics with their sources, printed by flai show and --json

## Work

A function in `flai/internal/topics` returns a story's topics, each with where it came from: own, epic, sub-project via tag, sub-project via claim (touches and open tasks', as ADR-0046 defines the claim), code, all. `flai show` prints them and `flai show --json` returns them.

## Done when

Behavior tests cover each source, a story with no tags or touches, and a touch outside every sub-project; `flai show` text and JSON carry the topics and their sources.

## Notes
