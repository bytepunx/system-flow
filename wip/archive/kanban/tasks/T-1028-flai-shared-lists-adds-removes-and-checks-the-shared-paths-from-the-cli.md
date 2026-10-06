---
id: T-1028
type: task
nature: feature
title: flai shared lists, adds, removes, and checks the shared paths from the CLI
status: done
parent: S-0295
owner: alex
created: 2026-10-06T12:15:45Z
updated: 2026-10-06T18:54:37Z
transitions:
  - to: ready
    at: 2026-10-06T18:37:34Z
    by: agent-S-0295
  - to: in-progress
    at: 2026-10-06T18:37:34Z
    by: agent-S-0295
  - to: done
    at: 2026-10-06T18:54:37Z
    by: agent-S-0295
stream: S-0295
tags: [flai]
touches: [flai/cmd/shared.go, flai/cmd/shared_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1025]
usage:
  source: log
  seconds: 1023
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 21484
      cache_read: 2894723
      cache_write: 84932
      cost: 1.4758
---
# T-1028 flai shared lists, adds, removes, and checks the shared paths from the CLI

## Work

Add the command T-1023's ADR names (proposed `flai shared`) in `flai/cmd/shared.go`, registered in `flai/cmd/root.go`:

- `list`: prints the patterns, one per line, and a JSON array with `--json`.
- `add <pattern>...` and `remove <pattern>...`: change the manifest through T-1025's edit function, refuse an invalid or duplicate pattern with the reason, and print what changed.
- `check <path>...`: says for each path or touches entry whether it is shared, and which pattern matched. With `--json` it prints an object per path. It also takes a story ID and reports each entry of its claim, so that an operator can see what a pattern would free before adding it.

Follow `flai board limit` for the shape of a command that writes a project file. Its writes do not commit, like the rest of flai's edits. Leave `docs/users/flai-reference.md` to T-1033, which regenerates it once every command is in.

This task waits for T-1025, whose matcher and edit function it calls.

## Done when

- `shared_test.go` covers each subcommand against a scratch project: list empty and full, add and remove with their refusals, check on a file, a folder, a non-matching path, and a story ID, with and without `--json`.
- `flai shared --help` and each subcommand's help say what they do and give the glob dialect.
- `scripts/flai-test.sh` passes.

## Notes
