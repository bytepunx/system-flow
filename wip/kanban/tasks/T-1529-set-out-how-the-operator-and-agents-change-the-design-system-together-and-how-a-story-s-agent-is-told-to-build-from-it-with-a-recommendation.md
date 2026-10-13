---
id: T-1529
type: task
nature: research
title: Set out how the operator and agents change the design system together, and how a story's agent is told to build from it, with a recommendation
status: done
parent: S-0380
owner: alex
created: 2026-10-13T17:46:07Z
updated: 2026-10-13T18:16:16Z
transitions:
  - to: ready
    at: 2026-10-13T18:07:46Z
    by: agent-S-0380
  - to: in-progress
    at: 2026-10-13T18:07:46Z
    by: agent-S-0380
  - to: done
    at: 2026-10-13T18:16:16Z
    by: agent-S-0380
stream: S-0380
tags: [dashboard, cli]
touches: [design/system/design-system-tool.md]
after: [T-1528]
usage:
  source: log
  seconds: 510
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 96
      output: 33209
      cache_read: 5052639
      cache_write: 163167
      cost: 2.6664
---
# T-1529 Set out how the operator and agents change the design system together, and how a story's agent is told to build from it, with a recommendation

## Work

Add to the finding the ways the operator and agents can work on the design system: the operator editing tokens on the Design page; an agent's proposal as a story whose branch changes the document, previewed beside the main branch's before acceptance; a thread on a heading of the document; or a store of draft variants. Weigh each for how a decision is reviewed, recorded, and undone. Then how agents use it: the prime pack briefing the document for stories on the UI it governs, a baseline convention, and a stylesheet flai writes from the tokens and checks. End with a recommendation and the whole design's recommendation. It waits for T-1528, whose recommendation on format it builds on.

## Done when

- The finding has a section on working together and a section on agents' use, each with its options and a recommendation, and a Recommendation section for the whole.

## Notes

- Layer 3 of 5.
