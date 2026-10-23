---
id: T-2093
type: task
nature: improvement
title: Add flai subproject list, add, set, and remove, creating the folder an added sub-project lacks, and document the command
status: backlog
parent: S-0491
owner: alex
created: 2026-10-23T20:39:00Z
updated: 2026-10-23T20:39:00Z
transitions: []
stream: S-0491
tags: [cli]
touches: [flai/cmd/root.go, flai/cmd/subproject.go, flai/cmd/subproject_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [T-2092]
---
# T-2093 Add flai subproject list, add, set, and remove, creating the folder an added sub-project lacks, and document the command

## Work

Add the command group `flai subproject` in `flai/cmd/subproject.go`, registered in `root.go` beside `components` and `shared`. The name says sub-project, as the design and the host API's `sub_project` parameter do, and leaves `flai serve project` alone.

- `flai subproject list [--json]` prints each `projects[]` entry: name, path, kind, tags, whether its folder exists, and whether it has `components` or `install`.
- `flai subproject add <name> --path <path> --kind <kind> [--tags a,b] [--autocommit] [--trailer]` adds an entry through `manifest.SetProjects` from T-2092. When `--path` is not given, it is the name. When the folder does not exist under the repository, the command creates it with a `README.md` that names the sub-project and its kind, so that git carries the folder, and prints the paths it wrote; it refuses a path that exists and is a file.
- `flai subproject set <name> [--path] [--kind] [--tags] [--name <new>]` changes the keys given and keeps the rest; `flai subproject remove <name>` removes the entry and leaves the folder alone, saying so.
- Refusals come from the manifest package as `refusedSettings` prints them for `components config`; `--autocommit` commits `system-flow.yaml` and the folder's `README.md` alone, as `editShared` in `shared.go` does, with `--trailer`; `--json` answers with the entry, the paths written, and the commit.

Document it: a section in `docs/users/flai.md` beside "Shared paths" and "List a project's components", saying how an operator or an agent adds a sub-project without editing `system-flow.yaml`, that the folder is created when missing, and that `flai components config` and `install` set the rest of an entry; a row in `design/system/flai-cli.md` § Commands; and `make flai-reference`, which rewrites `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` (`flai/cmd/reference_test.go` fails when it is stale).

It waits for T-2092, whose editor and validator it calls.

## Done when

- `flai subproject add web --kind sveltekit` in a project whose manifest says `projects: []` writes the entry as a block, creates `web/README.md`, and `flai subproject list` shows it with its folder present.
- `flai subproject set`, `remove`, and a refused add (duplicate name, path leaving the repository) behave as `flai/cmd/subproject_test.go` shows, writing nothing on a refusal.
- `--autocommit` commits the manifest and the README in one commit with the trailer given.
- `docs/users/flai.md`, `design/system/flai-cli.md`, `docs/users/flai-reference.md`, and the flag index in `docs/operators/settings.md` name the command, and `go test ./cmd/ -run 'TestFlagIndexIsCurrent|TestSettingsIndexIsComplete'` passes.
- `scripts/flai-test.sh` passes for the paths changed.

## Notes

Precedents: `flai/cmd/components.go` § `config` and `flai/cmd/shared.go` § `editShared`. `flai import` creates the layout folders in `flai/cmd/import.go`.
