---
id: T-1187
type: task
nature: feature
title: flai message send, reply, list, and show work the conversations, with --json
status: done
parent: S-0330
owner: alex
created: 2026-10-07T20:14:18Z
updated: 2026-10-07T20:48:38Z
transitions:
  - to: ready
    at: 2026-10-07T20:33:07Z
    by: agent-S-0330
  - to: in-progress
    at: 2026-10-07T20:33:08Z
    by: agent-S-0330
  - to: done
    at: 2026-10-07T20:48:38Z
    by: agent-S-0330
stream: S-0330
tags: [flai]
touches: [flai/cmd/message.go, flai/cmd/message_test.go, flai/cmd/root.go]
after: [T-1186]
usage:
  source: log
  seconds: 930
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 69
      output: 30114
      cache_read: 3983015
      cache_write: 130573
      cost: 2.1848
---
# T-1187 flai message send, reply, list, and show work the conversations, with --json

## Work

Add the `flai message` command over T-1186's package. It waits for T-1186, whose functions it calls. It shares no path with T-1188 and runs beside it.

- `flai message send <S-nnnn> "<text>" --from <S-nnnn> [--about <path>…]`; `--from` defaults to `FLAI_STORY` when it is set.
- `flai message reply <id> "<text>" --from <S-nnnn>`, `flai message list [--story S-nnnn] [--all]`, and `flai message show <id>`, each with `--json`; `list` says which side each open conversation awaits.
- Register the command in `flai/cmd/root.go`. The author is `FLAI_AGENT`, as for threads.

## Done when

- Tests in `flai/cmd/message_test.go` cover each subcommand, `--json`, and a refused send.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
