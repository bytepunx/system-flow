---
id: T-1348
type: task
nature: remediation
title: scripts/flai-snapshot.sh builds a host-only release with a plain version for the local release server
status: done
parent: S-0340
owner: alex
created: 2026-10-08T08:05:33Z
updated: 2026-10-08T08:24:12Z
transitions:
  - to: ready
    at: 2026-10-08T08:16:38Z
    by: agent-S-0340
  - to: in-progress
    at: 2026-10-08T08:16:39Z
    by: agent-S-0340
  - to: done
    at: 2026-10-08T08:24:12Z
    by: agent-S-0340
stream: S-0340
tags: [cli]
touches: [scripts/flai-snapshot.sh, flai/.goreleaser.yaml]
usage:
  source: log
  seconds: 453
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 11636
      cache_read: 1698993
      cache_write: 50630
      cost: 0.8948
---
# T-1348 scripts/flai-snapshot.sh builds a host-only release with a plain version for the local release server

## Work

Give `scripts/flai-snapshot.sh` a mode for the smoke test, such as `--local X.Y.Z`. It builds only the host's platform, skips GoReleaser's `before` hooks, and stamps the plain version given.

- The `before` hooks run `go mod tidy` and `go test ./...`. The integration tier has already run the tests by the time smoke runs.
- `selfupgrade.List` keeps only tags that are a plain `X.Y.Z` (`buildinfo.Bare`), and `install.sh` names the archive from the tag. A snapshot's own `0.0.0-<commit>` version and its `flai_snapshot_<os>_<arch>` archive name match neither.
- Change `flai/.goreleaser.yaml` so that a local build's version, its `buildinfo.Version`, and its archive name follow the version given. A tagged release and a plain `scripts/flai-snapshot.sh` build as they do now.
- GoReleaser's `release` has no `--single-target`. Limit the targets with the environment or a generated configuration overlay.

Waits for no task. S-0232, in progress, touches both files too: that holds the story, not this task.

## Done when

- The local mode writes one archive for the host to `flai/dist`, `flai_<X.Y.Z>_<os>_<arch>.tar.gz`, and a `checksums.txt` that names it.
- The `flai` in that archive answers `flai <X.Y.Z>` to `flai version`.
- `scripts/flai-snapshot.sh` with no arguments, and `goreleaser check`, give what they gave before.

## Notes

Drafted by the planner from S-0340's criterion 1 and the code as it stands: `flai/internal/selfupgrade/selfupgrade.go` (`List`, `plainVersion`), `install.sh` (`ARCHIVE`), and `flai/.goreleaser.yaml` (`snapshot`, `archives.name_template`).
