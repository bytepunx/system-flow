---
id: T-2927
type: task
nature: remediation
title: A hand-run container sees the main checkout's files, and refuses with the cause when a mount comes up empty
status: cancelled
parent: S-0667
owner: alex
created: 2026-11-10T08:48:52Z
updated: 2026-11-10T11:02:24Z
transitions:
  - to: cancelled
    at: 2026-11-10T11:02:24Z
    by: system-flow
stream: S-0667
tags: [docker, agent-container]
touches: [scripts/agent-container.sh]
after: [T-2926]
---
# T-2927 A hand-run container sees the main checkout's files, and refuses with the cause when a mount comes up empty

## Work

Waits for T-2926: it builds the fix T-2926 names in the story's Notes under `### Cause`.

Make the fix in `scripts/agent-container.sh`, the hand-run form of a container start, so that a container it starts from a story's agent's shell mounts the main checkout's paths with their files, as `flai serve`'s does. If T-2926 finds the fix lies outside the repository and the operator chose it on the story's thread, make only the check below.

Before it starts the agent or the command, the script checks that the mounts are not empty: in the same image and mounts, a file that must be there, such as the manifest at `$MANIFEST` and the worktree's `.git` file, is tested. When one is missing, it exits non-zero, runs nothing, and names the cause T-2926 found and its fix, with what it read (the docker context, `DOCKER_HOST`), on standard error.

Reproduce the cause as T-2926 recorded it, run the script, and record in the Notes that it refuses with that message; then clear the cause and record that it runs. Run `flai test scripts/agent-container.sh` for ShellCheck.

## Done when

- From a story's agent's shell on the host, `scripts/agent-container.sh S-nnnn -- ls "$ROOT/wip"` lists the main checkout's wip folder.
- With the cause T-2926 found in place, the script refuses with a message that names it, and runs nothing.
- `flai test scripts/agent-container.sh` passes.

## Notes
- 2026-11-10T11:02:24Z: moved to cancelled: T-2926 found the cause in zsh's $name:r modifier on a hand-typed -v spec; scripts/agent-container.sh builds its arguments in POSIX sh, which zsh never reads, so a check there guards nothing. T-3078 carries the fix.
