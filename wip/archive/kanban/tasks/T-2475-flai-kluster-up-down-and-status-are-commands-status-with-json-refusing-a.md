---
id: T-2475
type: task
nature: feature
title: flai kluster up, down, and status are commands, status with --json, refusing a project with no kluster block with the block to add
status: done
parent: S-0577
owner: alex
created: 2026-11-06T21:42:52Z
updated: 2026-11-07T01:17:20Z
transitions:
  - to: ready
    at: 2026-11-07T01:12:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T01:12:04Z
    by: system-flow
  - to: done
    at: 2026-11-07T01:17:20Z
    by: agent-S-0577
stream: S-0577
tags: [kluster, flai, cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
after: [T-2474]
usage:
  source: log
  seconds: 316
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 16945
      cache_read: 7568854
      cache_write: 84616
      cost: 2.4495
---
# T-2475 flai kluster up, down, and status are commands, status with --json, refusing a project with no kluster block with the block to add

## Work

It waits for T-2474, whose package the commands call.

Add `flai/cmd/kluster.go` with a `kluster` command and three subcommands, registered in `flai/cmd/root.go`:

- `up` brings the cluster up, with kluster's output streamed;
- `down` takes it down, asking first unless `--yes` (the global flag) is given;
- `status` prints the name, profile, provider, whether it runs, and its age, and with `--json` prints them as one object.

Each refuses a project with no `kluster:` block, naming the block and linking `design/system/kluster.md`. Each prints the package's missing-binary sentence as it is. Write the help text in the style of the other commands, with examples. Test each subcommand, the refusal, and `--json` in `flai/cmd/kluster_test.go` with a fake `kluster` on `PATH`.

In `flai/cmd/serve_actions.go`, refuse `flai serve enable kluster` for a project whose manifest has no `kluster:` block, with the same sentence naming the block to add, as `design/system/kluster.md` § The `kluster` host action says. The dashboard's Features toggle runs `flai serve enable` (`flai/internal/hostapi/settings.go`), so this one refusal covers both. `--all-projects` and `disable` are not refused. Test it in `flai/cmd/serve_actions_test.go`.

## Done when

- The three subcommands work as above and are listed in `flai --help`.
- The refusal and the missing-binary sentence are tested.
- `flai serve enable kluster` is refused without a block and allowed with one, as tested.
- `flai test` on the changed paths passes.

## Notes
