---
id: S-0016
type: story
nature: feature
title: Conventions guide
status: done
parent: E-0004
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-24T08:08:41Z
transitions:
  - to: ready
    at: 2026-09-24T07:58:22Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:01:03Z
    by: agent-S-0016
  - to: review
    at: 2026-09-24T08:05:24Z
    by: agent-S-0016
  - to: done
    at: 2026-09-24T08:08:41Z
    by: alex
tags: []
touches: [docs/users/conventions.md, docs/users/index.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0016 Conventions guide

## Goal
docs/users explains the layout, hierarchy, workflow, and narrative in user terms.

## Acceptance criteria
- [x] A reader can set up a conforming repo by hand from the guide

## Tasks
- T-0397 Write the conventions guide
- T-0398 Verify the by-hand walkthrough
- T-0399 Link the guide from the users index

## Notes
Verified 2026-09-24 by following docs/users/conventions.md, section "Set up a conforming repository by hand", literally in a scratch directory with no flai new or flai import: the code blocks were extracted from the guide and applied in order. The result passed flai check --strict with no warnings, as did a fresh clone of it; a story was then taken through in-progress, a task, review, and flai accept with check clean after each step, and the check stayed clean with the template's twelve baseline conventions copied in. The first run found that the example timestamps, being later than the real time, made the first flai move fail item.chronology; the guide now tells the reader to write the real UTC time.
