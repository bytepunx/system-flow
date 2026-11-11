---
id: T-3100
type: task
nature: remediation
title: "flai-build.sh links the built tree's own commit and date and builds with -buildvcs=false"
status: backlog
parent: S-0724
owner: alex
created: 2026-11-11T08:54:50Z
updated: 2026-11-11T08:54:50Z
transitions: []
stream: S-0724
tags: [go, build, worktree]
touches: [scripts/flai-build.sh, scripts/README.md]
after: [T-3093]
---
# T-3100 flai-build.sh links the built tree's own commit and date and builds with -buildvcs=false

## Work

Make `scripts/flai-build.sh` stamp `bin/flai` with the commit of the tree it builds, not the main checkout's.

- Take the commit with `git -C "$ROOT" rev-parse HEAD` and its commit time in UTC, in RFC 3339 as `vcs.time` gives it. git resolves a worktree's `.git` file, so in a story's worktree these are the story branch's.
- Pass them as `-ldflags "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=… -X …buildinfo.Date=…"`, and build with `-buildvcs=false`, so Go's lookup no longer reaches the main checkout's `.git` folder.
- When git cannot answer, as in a tree with no `.git`, build without the two values, as now. `flai version` then prints `unknown` for them.
- Leave `buildinfo.Version` unset. Since `OwnStamp()` reads the raw `Commit`, check that `flai version` on a tree build still reports no release stamp, and say so in the task's Notes.
- Change the `flai-build.sh` row of `scripts/README.md` to say it links the tree's commit and date.

It waits for T-3093: its test is what shows the fix works.

## Done when

- `scripts/flai-build-stamp-test.sh` passes.
- `bin/flai version` built in a story's worktree prints the worktree's HEAD as its commit, and built in the main checkout prints the main checkout's HEAD.
- `go version -m bin/flai` lists no `vcs.revision` for a tree build.
- The flai tier and the shell lint tier pass.

## Notes
