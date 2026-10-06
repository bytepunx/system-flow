---
id: T-1045
type: task
nature: improvement
title: The host API reads host.versions and dashboard.versions, and its upgrades install a named release only when it is published
status: backlog
parent: S-0298
owner: alex
created: 2026-10-06T21:45:25Z
updated: 2026-10-06T21:45:40Z
transitions: []
stream: S-0298
tags: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1043]
---
# T-1045 The host API reads host.versions and dashboard.versions, and its upgrades install a named release only when it is published

## Work

In `flai/internal/hostapi/writes.go` add the reads `host.versions` (`flai host versions --json`) and `dashboard.versions` (`flai dashboard versions --json`). Let `host.upgrade` decode an optional `version` and `dashboard.upgrade` an optional `tag`, each a bare semver; build `host upgrade --version <v>` and `dashboard upgrade --published --tag <t>` with it, and still never an image name. The commands refuse an unpublished release themselves (T-1040, T-1043); refuse anything that is not a semver here, before a command runs. Rewrite the comment that says neither takes a tag from the dashboard to cite the ADR, and have `describeHost` and `describeDashboardUpgrade` name the release asked for in the journal. Waits for T-1043 for `flai host versions` and `--version`, and through it for T-1040's `flai dashboard versions` and `--published`.

## Done when

- Both reads answer their lists, and both upgrades pass a valid version or tag and refuse a malformed one, tested in `writes_test.go`, including the detached upgrade that survives its connection dying.
- `scripts/flai-test.sh` passes for `flai/internal/hostapi`.

## Notes
