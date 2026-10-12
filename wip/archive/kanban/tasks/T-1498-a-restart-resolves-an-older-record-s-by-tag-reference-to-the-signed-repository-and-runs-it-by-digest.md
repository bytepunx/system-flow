---
id: T-1498
type: task
nature: remediation
title: A restart resolves an older record's by-tag reference to the signed repository and runs it by digest
status: done
parent: S-0236
owner: alex
created: 2026-10-12T19:52:35Z
updated: 2026-10-12T19:52:40Z
transitions:
  - to: ready
    at: 2026-10-12T19:52:38Z
    by: agent-S-0236
  - to: in-progress
    at: 2026-10-12T19:52:39Z
    by: agent-S-0236
  - to: done
    at: 2026-10-12T19:52:40Z
    by: agent-S-0236
stream: S-0236
tags: []
touches: [flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go, flai/internal/dashboard/resolve.go, flai/internal/dashboard/resolve_test.go, docs/users/flai.md, design/system/flai-cli.md]
after: [T-1244, T-1246]
usage:
  source: log
  seconds: 1
  estimated: true
  models: []
---
# T-1498 A restart resolves an older record's by-tag reference to the signed repository and runs it by digest

## Work

The verifier's review of the close-out found that a restart record written by an older flai names `ghcr.io/bytepunx/flaiover` by tag. `flai dashboard restart` and `flai host`'s watch started it as written, so a restart could pull a floating tag with no verified list. `pinnedRecordRef` resolves such a reference through the resolver in both restart paths. The release then runs by digest and is recorded, and a refused tag starts nothing. The `not-a-version` refusal no longer suggests `--image` for the signed repository. The user guide and `flai-cli.md` say so.

Waits for T-1244 and T-1246, which built the restart paths and their docs.

## Done when

- The watch and `flai dashboard restart` run an older by-tag record of the signed repository by digest, and record that digest.
- A refused tag in such a record starts nothing.
- `flai test` passes on the changed paths.

## Notes
