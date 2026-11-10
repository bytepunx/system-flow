---
id: T-2918
type: task
nature: improvement
title: scripts/close-out.sh takes --fresh and --all before the story, or CLOSE_OUT_FRESH and CLOSE_OUT_ALL, and passes them to each round's flai verify
status: backlog
parent: S-0671
owner: alex
created: 2026-11-10T08:48:07Z
updated: 2026-11-10T08:48:07Z
transitions: []
stream: S-0671
tags: [close-out, scripts]
touches: [scripts/close-out.sh, scripts/README.md]
---
# T-2918 scripts/close-out.sh takes --fresh and --all before the story, or CLOSE_OUT_FRESH and CLOSE_OUT_ALL, and passes them to each round's flai verify

## Work

- In `scripts/close-out.sh`, read leading `--fresh` and `--all` options before the story, in any order, and `CLOSE_OUT_FRESH=1` and `CLOSE_OUT_ALL=1` from the environment. Options come before the story because every argument after it goes to `git commit`, and `--all` is a `git commit` option too.
- Pass `--fresh` and `--all`, as asked, to each round's `flai verify <story> --record-issues` run, rounds two and three included. Pass neither to the `--sync-only` check.
- Update the usage comment and the usage error to `scripts/close-out.sh [--fresh] [--all] S-nnnn [git commit options]`, and name the two variables in the comment.
- Update the script's row in `scripts/README.md` to say so.
- Waits for nothing: the interface is fixed in the story's `### Planning` note.

## Done when

- `flai test scripts/close-out.sh` passes, ShellCheck included.
- In a worktree at a commit the verify record holds a pass for, `scripts/close-out.sh S-nnnn` answers from the record, and `scripts/close-out.sh --fresh S-nnnn` and `CLOSE_OUT_FRESH=1 scripts/close-out.sh S-nnnn` run every selected tier; the outcome of each run is written in this task's Notes.
- `scripts/close-out.sh S-nnnn -m "<message>"` still commits with the message.

## Notes
