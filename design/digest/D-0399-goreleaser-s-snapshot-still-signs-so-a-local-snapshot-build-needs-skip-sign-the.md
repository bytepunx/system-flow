---
id: D-0399
title: "GoReleaser's --snapshot still signs, so a local snapshot build needs --skip=sign: the signs entry's if is GoReleaser Pro only"
kind: decision
sources: [S-0232]
topics: [release, ci, go]
paths: [flai/.goreleaser.yaml, scripts/flai-snapshot.sh, flai/README.md, design/system/release-signing.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0399 GoReleaser's --snapshot still signs, so a local snapshot build needs --skip=sign: the signs entry's if is GoReleaser Pro only

S-0232 added `--skip=sign` to `scripts/flai-snapshot.sh` because "`--snapshot` does not skip signing, and the `signs` entry's `if` is GoReleaser Pro only". `flai/.goreleaser.yaml` says so in a comment. `design/system/release-signing.md` says snapshot builds are not signed, but not the flag that makes it so. Build a snapshot with `scripts/flai-snapshot.sh` (`make flai-snapshot`) rather than a bare `goreleaser release --snapshot`, which would try to sign with no key.

Source: S-0232.
