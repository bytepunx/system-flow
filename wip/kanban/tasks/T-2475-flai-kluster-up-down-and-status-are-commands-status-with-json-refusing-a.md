---
id: T-2475
type: task
nature: feature
title: flai kluster up, down, and status are commands, status with --json, refusing a project with no kluster block with the block to add
status: backlog
parent: S-0577
owner: alex
created: 2026-11-06T21:42:52Z
updated: 2026-11-06T21:42:52Z
transitions: []
stream: S-0577
tags: [kluster, flai, cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/root.go]
after: [T-2474]
---
# T-2475 flai kluster up, down, and status are commands, status with --json, refusing a project with no kluster block with the block to add

## Work

It waits for T-2474, whose package the commands call.

Add `flai/cmd/kluster.go` with a `kluster` command and three subcommands, registered in `flai/cmd/root.go`:

- `up` brings the cluster up, with kluster's output streamed;
- `down` takes it down, asking first unless `--yes` (the global flag) is given;
- `status` prints the name, profile, provider, whether it runs, and its age, and with `--json` prints them as one object.

Each refuses a project with no `kluster:` block, naming the block and linking `design/system/kluster.md`. Each prints the package's missing-binary sentence as it is. Write the help text in the style of the other commands, with examples. Test each subcommand, the refusal, and `--json` in `flai/cmd/kluster_test.go` with a fake `kluster` on `PATH`.

## Done when

- The three subcommands work as above and are listed in `flai --help`.
- The refusal and the missing-binary sentence are tested.
- `flai test` on the three paths passes.

## Notes
