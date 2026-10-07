---
id: T-1242
type: task
nature: feature
title: flai dashboard versions marks each dashboard release with no digest list this binary can verify, in its table and its JSON
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:48Z
updated: 2026-10-07T22:54:00Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/cmd/dashboard_versions.go, flai/cmd/dashboard_versions_test.go]
after: [T-1239, T-1240]
---
# T-1242 flai dashboard versions marks each dashboard release with no digest list this binary can verify, in its table and its JSON

## Work

In `flai/cmd/dashboard_versions.go`, list dashboard releases from the flaiover GitHub releases T-1239 lists, and for each one check whether it carries a digest list and a signature this binary verifies, with T-1240's code. Add a field to `listedDashboard`, such as `deployable`, false for a release with no list, one that does not verify, or one signed by a key the binary does not know, such as a release from before signing, and show it in the table. The host API's `dashboard.versions` and flaiover's `/api/dashboard` pass the JSON through unchanged, so the Updates page receives the field without a change there.

Make `requirePublishedDashboard`, in the same file, refuse a release that is not deployable as it refuses one that is not published, naming the deployable ones, so that `flai dashboard upgrade --tag X.Y.Z --published`, which the host API always passes, cannot deploy it.

Cover a signed release, one with no list, one whose list does not verify, and one signed by an unknown key in `dashboard_versions_test.go`, with `releasesStandIn` answering releases with assets.

Waits for T-1239, whose listing it calls, and T-1240, whose verification it calls.

## Done when

- `flai dashboard versions` and `--json` mark each release with no digest list this binary can verify, and mark a signed one deployable.
- `requirePublishedDashboard` refuses a release that is not deployable, naming the deployable ones.
- The four cases are tested and `flai test flai/cmd` passes.

## Notes
