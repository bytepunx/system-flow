---
id: T-1242
type: task
nature: feature
title: flai dashboard versions marks each dashboard release with no digest list this binary can verify, in its table and its JSON
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:53:48Z
updated: 2026-10-12T18:17:24Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/cmd/dashboard_versions.go, flai/cmd/dashboard_versions_test.go]
after: [T-1239, T-1240]
---
# T-1242 flai dashboard versions marks each dashboard release with no digest list this binary can verify, in its table and its JSON

## Work

In `flai/cmd/dashboard_versions.go`, keep listing the dashboard versions from the `flaiover/v*` tags, and join them with the flaiover GitHub releases (`selfupgrade.List` with `DashboardTagPrefix`). For each version check whether its release carries `flaiover_<version>.digests` and a signature this binary verifies, with `selfupgrade.FetchSigned` and T-1240's parser, a few at a time as `signatureRefusals` in `flai/cmd/selfupgrade.go` does for flai. Mark a version with no release, no list, a list that does not verify or parse, or one signed by a key the binary does not know, such as one from before signing, with the fields `listedFlai` uses for the same mark: `unverifiable` (bool, omitempty), `signature` (`missing`, `bad`, or `unknown-key`, omitempty), and `reason` (omitempty), worded from `SignatureError.Summary()`. Show it in the table as `not deployable: <reason>`. The host API's `dashboard.versions` and flaiover's `/api/dashboard` pass the JSON through unchanged, so the Updates page receives the fields without a change there.

Make `requirePublishedDashboard`, in the same file, refuse a release that is unverifiable as it refuses one that is not published, naming the reason and the deployable ones, so that `flai dashboard upgrade --tag X.Y.Z --published`, which the host API always passes, cannot deploy it.

Cover a signed release, a tag with no release, a release with no list, one whose list does not verify, and one signed by an unknown key in `dashboard_versions_test.go`, with `releasesStandIn` answering releases with assets.

Waits for T-1239, whose listing and fetch it calls, and T-1240, whose parser it calls.

## Done when

- `flai dashboard versions` and `--json` mark each release with no digest list this binary can verify, and leave a signed one unmarked.
- `requirePublishedDashboard` refuses an unverifiable release, naming the reason and the deployable ones.
- The cases are tested and `flai test flai/cmd` passes.

## Notes
