---
id: T-3620
type: task
nature: improvement
title: flai accept --by orchestrator keeps a recorded removal with --keep-removal and --reason while keep_removals is on and the narrative's Decisions name the path
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:57:42Z
updated: 2026-11-27T23:57:42Z
transitions: []
stream: S-0852
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, docs/users/flai-reference.md]
after: [T-3619]
---
# T-3620 flai accept --by orchestrator keeps a recorded removal with --keep-removal and --reason while keep_removals is on and the narrative's Decisions name the path

## Work

`keptRemovals` in `flai/cmd/accept.go` refuses every `--keep-removal` with `--by orchestrator` (ADR-0155). Change it as the ADR of T-3618 says. With the permission T-3619 adds on, the orchestrator's keep passes when it gives a reason and the story's narrative on its branch names the path under `## Decisions`. A keep the narrative does not name is refused, naming the path and the section to write. With the permission off, the refusal stays, and it now names the permission. `keptHeld` still refuses a path the record does not hold. `keptNotes` records `by: orchestrator` under `### Removals kept`, beside the orchestrator's evidence. The preview in `flai/internal/preview/accept.go` lists a path the orchestrator kept as kept, not as the blocker `undoes_main`. Update the flag's help, which no longer says operator only, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add a test in `flai/cmd/accept_undone_test.go` that reproduces I-0415. A branch rewrites a line main added, as S-0811 did. The orchestrator's acceptance with the keep, the permission on, and the path in the narrative's Decisions lands the story. The same acceptance is refused with the permission off, and refused for a path the narrative does not name.

It waits for T-3619, whose permission it reads.

## Done when

- The I-0415 reproduction test passes, and so do the two refusals.
- `flai accept --by orchestrator --dry-run --json` shows a kept path as kept, with its reason.
- `docs/users/flai-reference.md` is current, and `flai test` passes on the changed paths.

## Notes
