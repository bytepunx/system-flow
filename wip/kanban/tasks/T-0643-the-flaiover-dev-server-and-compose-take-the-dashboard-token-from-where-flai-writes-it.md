---
id: T-0643
type: task
nature: remediation
title: The flaiover dev server and compose take the dashboard token from where flai writes it
status: ready
parent: S-0183
owner: arobson
created: 2026-10-01T08:34:18Z
updated: 2026-10-01T08:34:25Z
transitions:
  - to: ready
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
stream: S-0183
tags: []
touches: [scripts/flaiover-dev.sh, flaiover/compose.yaml]
---
# T-0643 The flaiover dev server and compose take the dashboard token from where flai writes it

## Work
scripts/flaiover-dev.sh takes the token file from flai dashboard token --json. flaiover/compose.yaml mounts that file as the secret and drops the repository mount, PROJECT_DIR, and their comments.

## Done when
The dev script points FLAIOVER_TOKEN_FILE at the file flai reports; docker compose config renders with the token file and no project mount.

## Notes
