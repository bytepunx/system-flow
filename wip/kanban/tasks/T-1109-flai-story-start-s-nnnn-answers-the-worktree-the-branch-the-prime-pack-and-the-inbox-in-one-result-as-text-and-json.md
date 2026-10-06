---
id: T-1109
type: task
nature: improvement
title: flai story start S-nnnn answers the worktree, the branch, the prime pack, and the inbox in one result, as text and --json
status: backlog
parent: S-0274
owner: alex
created: 2026-10-06T22:53:30Z
updated: 2026-10-06T22:53:30Z
transitions: []
stream: S-0274
tags: [cli, go]
touches: [flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/cmd/items.go, flai/cmd/prime.go, flai/cmd/move.go]
after: [T-1099, T-1103]
---
# T-1109 flai story start S-nnnn answers the worktree, the branch, the prime pack, and the inbox in one result, as text and --json

## Work

Add `flai story start <S-nnnn>` in `flai/cmd/story_start.go`, registered on the `story` command in `newItemCmd` (`flai/cmd/items.go`), with `--budget` as `flai prime` takes it. It calls the storystart function under `FLAI_AGENT`, then the exported inbox function under the same name, and prints one result:

- As text: the move (and the epic's, when it followed), the worktree and branch as `flai stream open` prints them, the pack as `flai prime --story` prints it (through `printPack` in `prime.go`), then the inbox's threads awaiting you, ready stories, and changes.
- With `--json`: one object with `story`, `followed`, `worktree`, `branch`, `from`, `pack`, and `inbox`.

A refusal exits non-zero with the reason, the same refusal `flai move` gives for a story that is not ready or is held. Share its wording with `move.go` rather than writing it twice.

It waits for T-1099 and T-1103, whose functions it calls.

## Done when

- `story_start_test.go` covers a ready story started (text and `--json`, every key present), a backlog story refused, and a held story refused with its hold's reason, each refusal leaving the story where it was.
- `flai story start --help` says what it does and that it replaces `flai move`, `flai stream open`, `flai prime --story`, and `inbox` at the start of a story.
- `scripts/flai-test.sh` passes.

## Notes
