---
id: T-0610
type: task
nature: feature
title: An open question in a narrative leads to its story's page, shown and answerable there
status: done
parent: S-0173
owner: alex
created: 2026-09-30T01:26:59Z
updated: 2026-10-01T07:40:57Z
transitions:
  - to: ready
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:40:57Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 449
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 20639
      cache_read: 4965994
      cache_write: 71123
      cost: 1.9753
---

# T-0610 An open question in a narrative leads to its story's page, shown and answerable there

## Work

The inbox links a question entry to /items/<story>?question=<key>, not to the narrative document. The item page lists its story's open questions from the shared inbox, each with the answer form the inbox page has, and scrolls to and marks the one the link names.

## Done when

hrefFor never returns /docs for a question; the story page shows and answers its open questions; tests cover both.

## Notes
