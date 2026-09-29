---
id: T-0505
type: task
nature: feature
title: The context package loads design, tech, and ADRs and selects them by topics, links, one step of links, and supersession
status: done
parent: S-0137
owner: alex
created: 2026-09-29T00:36:07Z
updated: 2026-09-29T00:39:12Z
transitions:
  - to: ready
    at: 2026-09-29T00:36:13Z
    by: agent-S-0137
  - to: in-progress
    at: 2026-09-29T00:36:14Z
    by: agent-S-0137
  - to: done
    at: 2026-09-29T00:39:12Z
    by: agent-S-0137
stream: S-0137
tags: [cli]
touches: [flai/internal/context]
---
# T-0505 The context package loads design, tech, and ADRs and selects them by topics, links, one step of links, and supersession

## Work

- Load `design/system`, `design/tech`, and `design/adrs` (not README, not the ADR template) into documents cut at headings with `topics.Parse`, with ADR ID, refines, and superseded_by.
- A selection that owns each section once: the first reason wins, later reasons are listed on it.
- Steps as functions: by topics; linked from the story, epic, and tasks (markdown links, repository paths, ADR IDs); one step (ADRs a selected section links, ADRs a selected ADR refines); a superseded ADR replaced by what supersedes it.
- Fixture repository under `flai/internal/context/testdata/`.

## Done when

- Behavior tests cover each step and "no document printed twice" and pass with `go test -race -short ./internal/context/`.

## Notes
