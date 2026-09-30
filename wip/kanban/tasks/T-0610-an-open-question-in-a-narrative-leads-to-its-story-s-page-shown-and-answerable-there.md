---
id: T-0610
type: task
nature: feature
title: An open question in a narrative leads to its story's page, shown and answerable there
status: in-progress
parent: S-0173
owner: alex
created: 2026-09-30T01:26:59Z
updated: 2026-09-30T01:27:10Z
transitions:
  - to: ready
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-09-30T01:27:10Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [flaiover/src]
---

# T-0610 An open question in a narrative leads to its story's page, shown and answerable there

## Work

The inbox links a question entry to /items/<story>?question=<key>, not to the narrative document. The item page lists its story's open questions from the shared inbox, each with the answer form the inbox page has, and scrolls to and marks the one the link names.

## Done when

hrefFor never returns /docs for a question; the story page shows and answers its open questions; tests cover both.

## Notes
