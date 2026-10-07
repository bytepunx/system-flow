---
id: T-1139
type: task
nature: remediation
title: The design and the users' guide say that acceptance waits out a held index lock and that a failed acceptance commit is finished with flai accept
status: backlog
parent: S-0307
owner: alex
created: 2026-10-07T01:12:53Z
updated: 2026-10-07T01:12:53Z
transitions: []
stream: S-0307
tags: [flai, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1138]
---
# T-1139 The design and the users' guide say that acceptance waits out a held index lock and that a failed acceptance commit is finished with flai accept

## Work

Say what T-1138 changed, where acceptance is described. It waits for T-1138, since the help it regenerates and the behaviour it describes are that task's.

- `design/system/flai-cli.md`, the `flai accept` row: the commit step runs through the index-lock retry in `flai/internal/storygit/indexlock.go`, how long it waits, and that it never removes the lock; an item done and archived with its acceptance commit missing resumes at the commit, with the overlap notices taken from the story's commits; the failure message.
- `docs/users/flai.md`, under Accept and release: what an operator sees when another process holds the lock, and that running `flai accept <id>` again finishes an acceptance whose commit failed.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference` for the new `flai accept` help. No flag changes, so `docs/operators/settings.md` should not change; if it does, find out why before keeping it.

## Done when

- The `flai accept` row in `design/system/flai-cli.md` and the Accept and release section of `docs/users/flai.md` describe the retry and the resumed commit.
- `docs/users/flai-reference.md` matches `flai accept --help`, and `make flai-reference` leaves nothing to change.
- The markdown lint passes on the three files.

## Notes
