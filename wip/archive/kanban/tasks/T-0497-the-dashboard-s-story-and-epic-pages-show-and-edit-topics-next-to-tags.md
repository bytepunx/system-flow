---
id: T-0497
type: task
nature: feature
title: The dashboard's story and epic pages show and edit topics next to tags
status: done
parent: S-0135
owner: alex
created: 2026-09-28T22:41:23Z
updated: 2026-09-28T22:50:25Z
transitions:
  - to: ready
    at: 2026-09-28T22:41:40Z
    by: agent-S-0135
  - to: in-progress
    at: 2026-09-28T22:48:21Z
    by: agent-S-0135
  - to: done
    at: 2026-09-28T22:50:25Z
    by: agent-S-0135
stream: S-0135
tags: []
touches: [flaiover/src]
---

# T-0497 The dashboard's story and epic pages show and edit topics next to tags

## Work

The item editor on story and epic pages shows and edits topics next to tags, sending `topics` through the dashboard's edit route to the host's `item.edit`; the new-item form takes topics too if it takes tags.

## Done when

Vitest tests cover showing, editing, and sending topics; flaiover lint and tests pass.

## Notes
