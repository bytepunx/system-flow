---
id: T-0976
type: task
nature: improvement
title: flai touches adds with --add and removes with --remove, and its help says paths given alone replace the list
status: backlog
parent: S-0254
owner: alex
created: 2026-10-05T05:49:34Z
updated: 2026-10-05T05:49:34Z
transitions: []
stream: S-0254
tags: [cli]
touches: [flai/cmd/touches.go, flai/cmd/touches_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go]
---
# T-0976 flai touches adds with --add and removes with --remove, and its help says paths given alone replace the list

## Work

The fix I-0067 needs, proposed from its one instance: S-0206's agent ran `flai touches S-0206 <paths>` to add paths, and the story's five touches were replaced. Paths given alone keep replacing the list, because scripts and `flai stream sync`'s hint rely on that. The help says so, and two flags give the other edits:

- `--add`: the paths given join the list. One already there is kept once, and the order is kept.
- `--remove`: the paths given leave the list. A path not in it is not an error.
- `--add`, `--remove`, and `--clear` contradict each other, and `--add` or `--remove` with no paths is refused. Each refusal says what to run instead.
- Every path goes through `workitem.CleanTouches`, as a replacement does. An edit notice and the overlap report follow whenever the list changes.
- `Short` and a new `Long` say that paths given alone replace the list. The examples show `--add`, a replacement, `--clear`, and the show form.
- `flai stream sync` prints `flai touches <id> --add <outside paths>` instead of repeating every touch (`printOutside` in `flai/cmd/stream_sync.go`), and its test follows.

Tests in `flai/cmd/touches_test.go` reproduce the instance: adding to a story with touches keeps them. They also cover removing, a duplicate added, the contradicting flags, and a replacement that still replaces.

Waits for no task: it is the first layer.

## Done when

- `flai touches S-nnnn --add p` keeps the touches already there and adds `p`, and a test shows it
- `--remove` removes the paths given, and the contradicting flags are refused with what to run instead, each with a test
- `flai touches --help` says that paths given alone replace the list
- `flai stream sync` prints `flai touches <id> --add` with only the outside paths, and its test passes
- `scripts/flai-test.sh` passes

## Notes
