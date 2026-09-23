---
id: S-0094
type: story
nature: remediation
title: Tagging should occur before push
status: done
parent: E-0003
owner: alex
created: 2026-09-23T00:32:19Z
updated: 2026-09-23T01:31:06Z
transitions:
  - to: ready
    at: 2026-09-23T00:32:28Z
    by: alex
  - to: in-progress
    at: 2026-09-23T00:57:57Z
    by: system-flow
  - to: review
    at: 2026-09-23T00:58:03Z
    by: system-flow
  - to: done
    at: 2026-09-23T01:31:06Z
    by: alex
tags: [dashboard, cli]
touches: [flai/cmd, flai/internal/release]
---
# S-0094 Tagging should occur before push

## Goal

After making changes to how the workflow should happen so that stories can batch in done before a push occur, it appears that no tagging is taking place, meaning the release workflow in GitHub does not run.

This should be corrected so that tagging occurs when the user says to push/publish changes to the remote. It's very important that tagging occurs before the push to the remote.

## Acceptance criteria
- [x] No tagging is calculated when a story is moved to done
- [x] Tagging does occur when the user clicks push
- [x] Tagging must occur before the push to the remote
- [x] Tags must be included in the push to the remote

## Tasks
- T-0336 flai push --pending tags everything pending before it pushes

## Notes

Once the tests pass, build the flai command locally and use bin/flai to tag this so the push results in a release so that the working environment can get an upgrade and begin functioning correctly again.
