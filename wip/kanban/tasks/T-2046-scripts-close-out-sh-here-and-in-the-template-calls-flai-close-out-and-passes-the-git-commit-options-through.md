---
id: T-2046
type: task
nature: improvement
title: scripts/close-out.sh, here and in the template, calls flai close-out and passes the git commit options through
status: backlog
parent: S-0466
owner: alex
created: 2026-10-22T21:00:16Z
updated: 2026-10-22T21:00:16Z
transitions: []
stream: S-0466
tags: [cli, template]
touches: [scripts/close-out.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/README.md]
after: [T-2031]
---
# T-2046 scripts/close-out.sh, here and in the template, calls flai close-out and passes the git commit options through

## Work

Make `scripts/close-out.sh` a wrapper: it keeps its usage line and its branch check, then runs `scripts/flai.sh close-out "$story" "$@"`, passing the git commit options through, and exits with its status, so that `scripts/close-out.sh S-nnnn -m "<message>"` keeps working for the agents and the docs and the rounds, the hints, and the last line come from flai. Do the same in `template/root/scripts/close-out.sh` with the installed `flai`, and rewrite its header to say what flai does and that it needs the flai release that ships `flai close-out`, which the ADR of T-2030 and the changelog name; the template tier of `flai test --all` compares the two trees, so keep them in step as `scripts/README.md` and `template/root/scripts/README.md` describe them. Update the `close-out.sh` row of both READMEs. Keep the scripts POSIX `sh` under `set -eu`, as `tooling.md` says.

Waits for T-2031, which adds the command the scripts call. It does not wait for T-2033: the command's name and commit options are settled by T-2031, and the detached form changes nothing in the wrapper.

## Done when

- `scripts/close-out.sh S-nnnn -m "<message>"` in a story's worktree runs `flai close-out` and prints its last line; with no argument it prints the usage and exits 2; on a branch that is not the story's it refuses as today.
- `template/root/scripts/close-out.sh` does the same with `flai` from `PATH`, and `make flai-test`'s template tier passes.
- Both `scripts/README.md` rows describe the wrapper and name `flai close-out`; `flai check --strict` is clean.

## Notes

The third layer, with T-2033 (the detached run); the two share no path.
