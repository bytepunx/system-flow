---
id: T-2109
type: task
nature: improvement
title: The manifest names, beside flai.minimum, the command that runs the flai the tree builds, with its settings row, and this repository sets it to scripts/flai.sh
status: done
parent: S-0495
owner: alex
created: 2026-10-25T20:32:42Z
updated: 2026-10-25T20:45:04Z
transitions:
  - to: ready
    at: 2026-10-25T20:36:39Z
    by: agent-S-0495
  - to: in-progress
    at: 2026-10-25T20:37:39Z
    by: agent-S-0495
  - to: done
    at: 2026-10-25T20:45:04Z
    by: agent-S-0495
stream: S-0495
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md, design/system/project-manifest.md, system-flow.yaml]
after: [T-2108]
usage:
  source: log
  seconds: 445
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 6737
      cache_read: 982995
      cache_write: 52688
      cost: 0.6619
---
# T-2109 The manifest names, beside flai.minimum, the command that runs the flai the tree builds, with its settings row, and this repository sets it to scripts/flai.sh

## Work

Add to the manifest's `flai` block in `flai/internal/manifest/manifest.go`, beside `Minimum` (`flai.minimum`), the key the ADR of T-2108 names (`flai.command` as recommended): a string, optional, a command relative to the repository root that runs the flai the tree builds, with a doc comment saying that task done's check step runs it in a story's worktree in place of the running flai's check, and that only a project that builds flai sets it. `Load` refuses an absolute path or one that leaves the repository (`..`), naming the key, as it refuses a bad `flai.minimum`; `manifest_test.go` covers a set key, an unset one, and a refused one.

Give the key its row in `docs/operators/settings.md` § Project manifest, which `TestSettingsIndexIsComplete` in `flai/cmd/settings_doc_test.go` requires for every key `manifest.Manifest` carries, and describe it in `design/system/project-manifest.md` beside `flai.minimum`, with the ADR linked.

Set it in this repository's `system-flow.yaml` to `scripts/flai.sh` with `flai manifest set flai.command=scripts/flai.sh` from the story's worktree; when the command does not take the key, hand-edit the manifest and run `flai check`, as `tooling.md` allows for a field no command sets. `scripts/flai.sh` builds `bin/flai` from the tree it is run in (`scripts/env.sh` sets `ROOT` from the script's own location), so in a worktree it is that branch's flai.

Waits for T-2108, which fixes the key's name and shape. The check-step task waits for this one, since it reads the key.

## Done when

`go test ./internal/manifest/... ./cmd/...` passes without `-short`, `TestSettingsIndexIsComplete` among them; the manifest of the story's worktree loads with the key set to `scripts/flai.sh`; a manifest with an absolute path is refused with a message naming the key; and `flai check --strict --story S-0495` reports nothing.

## Notes

The template's manifest does not get the key: a project made from the template installs flai releases and never carries a format its flai does not know.
