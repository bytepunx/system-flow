---
id: T-0719
type: task
nature: improvement
title: flai release --pending and Publish refuse a clone behind its remote, naming the fetch and merge to run
status: in-progress
parent: S-0195
owner: arobson
created: 2026-10-02T23:31:03Z
updated: 2026-10-02T23:46:36Z
transitions:
  - to: ready
    at: 2026-10-02T23:31:26Z
    by: agent-S-0195
  - to: in-progress
    at: 2026-10-02T23:46:36Z
    by: agent-S-0195
stream: S-0195
tags: []
touches: [flai/internal/release/remote.go, flai/internal/release/remote_test.go, flai/cmd/release.go, flai/cmd/release_test.go, flai/internal/preview]
after: [T-0717]
---
# T-0719 flai release --pending and Publish refuse a clone behind its remote, naming the fetch and merge to run

## Work

`release.CheckRemote` also asks the remote for the head of the branch this clone tracks; when that commit is not in this clone's history (not fetched, or fetched and not merged), the clone is behind, and `flai release --pending` refuses with exit 3 before planning, naming what to run (`git fetch <remote>`, then merge or rebase onto `<remote>/<branch>`, as ADR-0067 allows), as S-0174 does for tags. `--dry-run` and the dashboard's publish preview (`preview.Publish`) say the same, so Publish offers nothing until the clone is in step. flai still never fetches by itself. Waits for T-0717: it builds what ADR-0067 decides.

## Done when

- A publish from a clone whose remote branch has commits it lacks is refused before anything is tagged, with the commands to run
- Tests cover the refusal, a publish after fetching and merging, and the dry run's message

## Notes
