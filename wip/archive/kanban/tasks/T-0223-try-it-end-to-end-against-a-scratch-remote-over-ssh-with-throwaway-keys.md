---
id: T-0223
type: task
nature: feature
title: Try it end to end against a scratch remote over SSH with throwaway keys
status: done
parent: S-0062
owner: alex
created: 2026-09-19T08:21:01Z
updated: 2026-09-19T08:29:07Z
transitions:
  - to: ready
    at: 2026-09-19T08:29:07Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:29:07Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:29:07Z
    by: system-flow
stream: S-0062
tags: []
---

# T-0223 Try it end to end against a scratch remote over SSH with throwaway keys

## Work
In the scratchpad, never this repository or the operator's key: a scratch remote served over SSH from a container, throwaway keys, a scratch project whose `origin` is that remote, and the branch's image started by the branch's `flai dashboard --push-key`. Accept a story from the board's API and see the commit and tag on the remote. Then: the default (no key: accepted locally, not pushed), a user ID the image does not know, a key with a passphrase, a host with no known_hosts entry. Tear everything down.

## Done when
- Each case is recorded in the narrative with what was run and seen
- Nothing of the lab is left behind

## Notes
