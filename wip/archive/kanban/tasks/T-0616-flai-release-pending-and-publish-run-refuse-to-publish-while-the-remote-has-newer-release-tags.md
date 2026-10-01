---
id: T-0616
type: task
nature: remediation
title: flai release --pending and publish.run refuse to publish while the remote has newer release tags
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T07:51:03Z
updated: 2026-10-01T08:00:02Z
transitions:
  - to: ready
    at: 2026-10-01T08:00:02Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T08:00:02Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:00:02Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [flai/cmd/release.go, flai/internal/hostapi]
usage:
  source: log
  seconds: 0
  models: []
---
# T-0616 flai release --pending and publish.run refuse to publish while the remote has newer release tags

## Work

- Before anything is applied, committed, tagged, or pushed, `flai release --pending` (and so `publish.run`, and `flai push --pending` with auto-publish) asks the remote; a component behind, or a remote that cannot be asked, refuses with exit 3 (a conflict to the dashboard) and says what to do.

## Done when

- Tests show the refusal changes nothing locally (no commit, no tag) and that a clone in step still publishes; `make test`, `make integration`, and lint pass.

## Notes
