---
id: T-0461
type: task
nature: research
title: Survey how an agent's context can be primed with only the documentation its story needs
status: done
parent: S-0125
owner: alex
created: 2026-09-26T08:04:21Z
updated: 2026-09-26T08:08:03Z
transitions:
  - to: ready
    at: 2026-09-26T08:04:26Z
    by: agent-S-0125
  - to: in-progress
    at: 2026-09-26T08:04:27Z
    by: agent-S-0125
  - to: done
    at: 2026-09-26T08:08:03Z
    by: agent-S-0125
stream: S-0125
tags: []
touches: [design/system/agent-context.md, design/system/README.md]
---
# T-0461 Survey how an agent's context can be primed with only the documentation its story needs

## Work

Measure what an agent loads today (`flai prime --cat`, CLAUDE.md, what it reads after) and the size of each documentation type. Survey how agent tools and research select context: always-on versus path-scoped rules, descriptions loaded with bodies on demand, repository maps, lexical and embedding retrieval, link graphs. Map the signals a story already carries (tags, touches, parent epic, linked ADRs and design documents, manifest projects, `design/tech` "Where" rows). Write the findings and ranked designs to `design/system/agent-context.md` and index it.

## Done when

`design/system/agent-context.md` states what is loaded today with numbers, what the options are, and a recommendation; `design/system/README.md` lists it; the markdown lint passes.

## Notes
