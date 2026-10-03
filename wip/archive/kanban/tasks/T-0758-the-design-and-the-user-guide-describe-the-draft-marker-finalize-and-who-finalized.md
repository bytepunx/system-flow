---
id: T-0758
type: task
nature: feature
title: The design and the user guide describe the draft marker, Finalize, and who finalized
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:37:09Z
updated: 2026-10-03T08:04:27Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:57:57Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T08:04:27Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, design/system/flai-cli.md, design/system/work-hierarchy.md, design/adrs]
after: [T-0754, T-0755, T-0756, T-0757]
usage:
  source: log
  seconds: 390
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 90
      output: 24491
      cache_read: 5396201
      cache_write: 140092
      cost: 2.4684
---
# T-0758 The design and the user guide describe the draft marker, Finalize, and who finalized

## Work

`design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the `[Draft]` indicator, the markers on cards, the overview, and search, the Finalize button on the story page and in the edit form, and the `item.finalize` write. Where TH-0083 settles on a new front matter key, an ADR refining ADR-0074 records it, and `design/system/work-hierarchy.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` describe it. It waits for the four tasks it describes.

## Done when

- [ ] The design and the user guide describe what T-0754, T-0755, T-0756, and T-0757 built
- [ ] Any new front matter key has its ADR and its place in the design and the command guide

## Notes
