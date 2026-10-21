---
id: T-1979
type: task
nature: improvement
title: flai move to done prints a refused acceptance as data with exit 4 under --json, and item.move answers it as Refused
status: done
parent: S-0465
owner: alex
created: 2026-10-21T07:36:21Z
updated: 2026-10-21T08:16:46Z
transitions:
  - to: ready
    at: 2026-10-21T07:39:03Z
    by: agent-S-0465
  - to: in-progress
    at: 2026-10-21T07:39:03Z
    by: agent-S-0465
  - to: done
    at: 2026-10-21T07:54:28Z
    by: agent-S-0465
stream: S-0465
tags: [flai]
touches: [flai/cmd/move.go, flai/cmd/move_refused_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai.md, design/system/flai-cli.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 925
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 72
      output: 408
      cache_read: 2585274
      cache_write: 277951
      cost: 1.0635
---
# T-1979 flai move to done prints a refused acceptance as data with exit 4 under --json, and item.move answers it as Refused

## Work

`flai move <story> done` from review runs `a.acceptItem` in `flai/cmd/move.go` and returns its error as it is. So under `--json` a refusal for lint the fix left is plain text with exit 1. `flai accept` sends the same error through `a.acceptRefusal` in `flai/cmd/accept.go`, which prints `{"refused": {...}}` and exits 4 (`exitDocRefused`, S-0378). That refusal carries the `findings` with each file's `checkout` and `story`, and `undoes_main` (ADR-0155).

- In `flai/cmd/move.go`, return `a.acceptRefusal(err)` when `acceptItem` fails, so that a move to done refuses exactly as `flai accept` does. Without `--json`, the text and the exit code stay as they are.
- In `flai/internal/hostapi/writes.go`, give `item.move` `exits: map[int]int{4: Refused}`, as `accept.run` has, so that the channel passes the refused object to the dashboard as the error's data. Check that the channel runs `item.move` with `--json`.
- Add `flai/cmd/move_refused_test.go`. Reuse the fixture behind `acceptRefusal` in `flai/cmd/accept_fix_test.go`: a story whose fix leaves lint, on its branch and in the main checkout. Assert that `flai move S-0001 done --json` exits 4 and prints the same `refused` object, findings and checkouts, as `flai accept --json`. Add a case to `flai/internal/hostapi/writes_test.go` asserting that `item.move` maps exit 4 to `Refused` with the data.
- Document the `--json` refusal of a move to done in `docs/users/flai.md` § Accept and release, beside the `flai accept --json` one, and in the `flai move` row of `design/system/flai-cli.md`.

This task waits for no other: the dashboard tasks read the shape `flai accept --json` already gives.

## Done when

- `flai move <story> done --json` exits 4 and prints `{"refused": {"id", "reason", "findings"}}` for lint the fix left, and `{"refused": {"id", "reason", "undoes_main"}}` for a branch that undoes main's change, the same objects as `flai accept --json`.
- `item.move` answers exit 4 as `Refused` with the refused object as data.
- The new tests pass, and `flai test` passes on the changed paths.
- `docs/users/flai.md` and `design/system/flai-cli.md` say so.

## Notes
