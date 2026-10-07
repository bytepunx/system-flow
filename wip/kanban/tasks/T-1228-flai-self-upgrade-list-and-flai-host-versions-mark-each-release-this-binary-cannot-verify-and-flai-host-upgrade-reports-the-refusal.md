---
id: T-1228
type: task
nature: feature
title: flai self-upgrade --list and flai host versions mark each release this binary cannot verify, and flai host upgrade reports the refusal
status: backlog
parent: S-0233
owner: alex
created: 2026-10-07T22:40:31Z
updated: 2026-10-07T22:40:31Z
transitions: []
stream: S-0233
tags: [cli]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/signature.go, flai/internal/selfupgrade/signature_test.go, flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, flai/cmd/host.go, flai/cmd/host_versions_test.go]
after: [T-1226]
---
# T-1228 flai self-upgrade --list and flai host versions mark each release this binary cannot verify, and flai host upgrade reports the refusal

## Work

Give each published flai release a mark for whether this binary can verify its signature, from what T-1226 built: a release with no `checksums.txt.sig` asset, such as one from before signing, or one signed by a key the binary does not know, cannot be verified. Add the mark to `listedFlai` in `flai/cmd/selfupgrade.go`, in its JSON (for example `"unverifiable": true` with a `reason`) and as a mark in `writeFlaiReleases`' text.

`flai host versions` prints the same list through `printHostVersions` in `flai/cmd/host.go`, so it carries the mark in text and `--json`. Check that `flai host upgrade --version` with such a release ends with self-upgrade's refusal, naming the release and the reason, and that the host does not restart.

Update the help text of `self-upgrade` and `host upgrade` to say that a release is verified before it is installed.

Waits for T-1226: the mark uses its verification, and both change `flai/internal/selfupgrade/selfupgrade.go` and `signature.go`.

## Done when

- `flai self-upgrade --list` and `flai host versions`, in text and with `--json`, mark each release whose signature this binary cannot verify, and say why.
- `flai host upgrade --version` with such a release is refused with the release and the reason, and nothing is replaced.
- Tests in `flai/cmd/selfupgrade_test.go` and `flai/cmd/host_versions_test.go` cover the mark, and `flai test` passes on the paths changed.

## Notes
