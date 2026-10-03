---
id: T-0752
type: task
nature: improvement
title: An ADR and the design and user guide describe how an epic follows its stories
status: done
parent: S-0200
owner: alex
created: 2026-10-03T07:11:28Z
updated: 2026-10-03T07:37:05Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T07:32:13Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T07:37:05Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [design/adrs, design/system/workflow.md, design/system/work-hierarchy.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-0748, T-0749, T-0750]
usage:
  source: log
  seconds: 292
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 65
      output: 22253
      cache_read: 4915232
      cache_write: 110424
      cost: 2.1695
---
# T-0752 An ADR and the design and user guide describe how an epic follows its stories

## Work

Record the rule with `flai adr new`, refining ADR-0004. Describe it in `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/system/flai-cli.md`, and the user guide in `docs/users/`, and in the template's copies where it has them.

It waits for T-0748, T-0749, and T-0750, so that it describes what they built.

## Done when

- [x] The ADR is accepted and refines ADR-0004
- [x] The design, the user guide, and the template say what the code does

## Notes
