---
id: T-2092
type: task
nature: improvement
title: Validate and edit the manifest's projects list in flai's manifest package, and report manifest.projects in flai check
status: in-progress
parent: S-0491
owner: alex
created: 2026-10-23T20:37:56Z
updated: 2026-10-23T20:48:32Z
transitions:
  - to: ready
    at: 2026-10-23T20:48:31Z
    by: agent-S-0491
  - to: in-progress
    at: 2026-10-23T20:48:32Z
    by: agent-S-0491
stream: S-0491
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/projects.go, flai/internal/manifest/projects_test.go, flai/internal/manifest/settings.go, flai/internal/check/check.go, flai/internal/check/projects_test.go]
---
# T-2092 Validate and edit the manifest's projects list in flai's manifest package, and report manifest.projects in flai check

## Work

Give `flai/internal/manifest` the model and the editor every later layer builds on, in a new file `projects.go` beside `settings.go`:

- `Manifest.ProjectErrors()` reports, each by the entry's index, field, and name as `ComponentErrors` does (`manifest.go` § `ComponentErrors`), a `projects[]` entry whose `name` is empty or repeats another's, whose `path` is empty, absolute, not clean, leaves the repository, or repeats another's, and whose `kind` is empty. A kind flai does not know is not an error: kinds are open (`go`, `node`, `python`, `rust`, `sveltekit`, `svelte`, `template` today, `importer.go` and `ComponentKinds`). Given an `fs.FS` of the repository root, as `ComponentErrors` takes one, it also reports a `path` that is not a folder there.
- `EditProjects(data, fsys, op, Project)` returns the manifest bytes with one entry added, replaced by name (every key of the entry given, components and install kept when the change leaves them out), or removed by name, keeping every other line, the key order, and the comments, as `EditComponents` does (`settings.go` § `componentsLines`). It must handle a `projects` written on one line, `projects: []`, as the template writes it (`template/root/system-flow.yaml.tmpl`), by turning it into a block holding the new entry; `EditComponents` refuses that form today. It refuses an entry it cannot read back, as `componentsReadBack` does, and with a `*RefusedError` naming each field, any problem `ProjectErrors`, `ComponentErrors`, or `InstallErrors` finds in the manifest as it would be. `SetProjects(path, fsys, op, Project)` writes it, or nothing on a refusal.
- `settingsProblems` in `settings.go` gains `ProjectErrors`, so `flai manifest set` refuses a manifest with one, as it does for components and install.
- `flai check` reports each `ProjectErrors` finding as `manifest.projects`, beside `manifest.components` and `manifest.install` in `check.go` § `layout`, anchored at `keyLine(mf, "projects")`.

Tests in `projects_test.go` cover each refusal, the one-line `projects: []` case, add, set, and remove with comments kept, and the read-back check; `flai/internal/check/projects_test.go` follows `install_test.go` with `findings(t, repo, rule)`.

This task waits for no other: the CLI, the host API, and the MCP tool all call it.

## Done when

- `manifest.ProjectErrors` reports an empty or duplicate name, an empty, absolute, unclean, escaping, or duplicate path, an empty kind, and a path that is no folder, and `go test ./internal/manifest/...` shows each.
- `manifest.EditProjects` adds the first entry to a manifest that says `projects: []`, adds, replaces, and removes an entry in a block, and keeps every other line and comment, shown by tests.
- `flai check` on a manifest with a duplicate sub-project name reports `manifest.projects`, and `flai manifest set` refuses such a manifest.
- `scripts/flai-test.sh` passes for the paths changed.

## Notes

Precedents: `EditComponents` and `componentsLines` in `flai/internal/manifest/settings.go`, `AddShared` and `RemoveShared` in `shared.go`, `writeProjects` in `flai/cmd/import.go`, which replaces `projects: []` only at import.
