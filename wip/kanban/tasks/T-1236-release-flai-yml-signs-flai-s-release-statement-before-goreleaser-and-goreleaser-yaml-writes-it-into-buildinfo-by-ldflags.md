---
id: T-1236
type: task
nature: feature
title: release-flai.yml signs flai's release statement before GoReleaser, and .goreleaser.yaml writes it into buildinfo by ldflags
status: backlog
parent: S-0235
owner: alex
created: 2026-10-07T22:48:55Z
updated: 2026-10-07T22:48:55Z
transitions: []
stream: S-0235
tags: [cli, ci]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml]
after: [T-1234]
---
# T-1236 release-flai.yml signs flai's release statement before GoReleaser, and .goreleaser.yaml writes it into buildinfo by ldflags

## Work

The build half of criterion 1. It waits for T-1234, whose linker variable names the `ldflags` set.

- In `.github/workflows/release-flai.yml`, after the release-key check and before GoReleaser, write the statement `flai <version> <commit>` to a file with no trailing newline (`printf '%s'`), the version trimmed from the tag (`flai/v` removed) and the commit as GoReleaser's `{{ .ShortCommit }}` gives it (`git rev-parse --short HEAD` in the same checkout), and sign it with `cosign sign-blob --key env://COSIGN_PRIVATE_KEY --tlog-upload=false`, with the flags S-0232's `signs` uses. Hand the statement and the base64 signature to the GoReleaser step through its environment.
- In `flai/.goreleaser.yaml`, add two `-X` lines to `ldflags` for T-1234's variables, read with `envOrDefault` so that `goreleaser --snapshot` (`scripts/flai-snapshot.sh`, `make flai-snapshot`), which has neither, still builds and carries no stamp.
- Check that the statement GoReleaser's version and commit give and the one the workflow signed are the same line: a mismatch would make every release warn about itself.
- Run `scripts/flai-snapshot.sh` and confirm the snapshot's `flai version --json` shows no stamp.

## Done when

- `release-flai.yml` signs the line before GoReleaser runs and passes the statement and signature to it through the environment.
- `flai/.goreleaser.yaml` writes both into `buildinfo` by `ldflags`, and a snapshot build still succeeds and carries no stamp.

## Notes

Drafted by the planner. A real signed release cannot be made from the story; the first `flai/v*` tag after acceptance proves the workflow, which S-0233's verification depends on too.
