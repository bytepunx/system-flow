---
id: T-0226
type: task
nature: feature
title: flai push --pending on the host
status: done
parent: S-0063
owner: alex
created: 2026-09-19T08:34:59Z
updated: 2026-09-19T08:37:34Z
transitions:
  - to: ready
    at: 2026-09-19T08:37:33Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:37:33Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:37:34Z
    by: system-flow
stream: S-0063
tags: []
---

# T-0226 flai push --pending on the host

## Work
`flai push --pending [--dry-run] [--remote]`: on the host, with the host's own credentials, push the branch and the tags on unpushed commits when the commits ahead include an acceptance; say "nothing pending" otherwise; refuse when the branch has diverged from the remote-tracking branch and say to fetch and merge; never force. `--dry-run` prints what it would push. Tests against a scratch bare remote.

## Done when
- The command pushes an unpushed acceptance with its tags, and the tests cover nothing pending, no acceptance among the commits, diverged, and dry run

## Notes
