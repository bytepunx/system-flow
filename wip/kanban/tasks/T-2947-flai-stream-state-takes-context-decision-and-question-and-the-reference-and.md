---
id: T-2947
type: task
nature: remediation
title: flai stream state takes --context, --decision, and --question, and the reference and flag index list them
status: in-progress
parent: S-0668
owner: alex
created: 2026-11-10T08:50:47Z
updated: 2026-11-10T09:55:09Z
transitions:
  - to: ready
    at: 2026-11-10T09:55:08Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:55:09Z
    by: system-flow
stream: S-0668
tags: [flai, narrative, cli]
touches: [flai/cmd/stream.go, flai/cmd/stream_state_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-2944]
---
# T-2947 flai stream state takes --context, --decision, and --question, and the reference and flag index list them

## Work

In `newStreamStateCmd` in `flai/cmd/stream.go`, add three flags beside `--current` and `--next`:

- `--context "<text>"` replaces `## Context`.
- `--decision "<text>"` appends one bullet to `## Decisions`. It may be given more than once.
- `--question "<text>"` appends one dated bullet to `## Open questions`, above the threads block. It may be given more than once.

Pass them to the widened `SetStreamState`. `-` reads standard input for `--context`, as for `--current`, and only one flag may be `-`. Giving none of the five flags is still an error, and its message names all five. Say what each flag does in `Use`, `Long`, and `Example`, and keep exit 4 for a refusal.

Add tests in `flai/cmd/stream_state_test.go`, in the style of `TestStreamStateEachFlagAloneAndJSON`. Cover each flag alone, a repeated `--decision`, and `--json`'s `written`.

Run `make flai-reference` to write `docs/users/flai-reference.md` and the flag rows of `docs/operators/settings.md` again.

Waits for T-2944, whose `SetStreamState` it calls.

## Done when

- `flai test flai/cmd/stream.go flai/cmd/stream_state_test.go docs/users/flai-reference.md docs/operators/settings.md` passes, `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` included.
- `flai stream state S-nnnn --decision "x"` on a story in progress adds `- x` under `## Decisions`, and every other section stays as it was.

## Notes
