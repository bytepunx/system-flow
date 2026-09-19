---
id: T-0222
type: task
nature: feature
title: The image carries an SSH client
status: done
parent: S-0062
owner: alex
created: 2026-09-19T08:21:01Z
updated: 2026-09-19T08:29:07Z
transitions:
  - to: ready
    at: 2026-09-19T08:29:06Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:29:06Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:29:07Z
    by: system-flow
stream: S-0062
tags: []
---

# T-0222 The image carries an SSH client

## Work
`flaiover/Dockerfile`: add the OpenSSH client to the final stage (about 0.7 MB on Alpine, measured in S-0052) and nothing else. `design/tech` records it.

## Done when
- The image builds and `ssh -V` runs in it as an arbitrary user ID with a passwd entry mounted

## Notes
