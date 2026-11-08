---
id: T-2684
type: task
nature: improvement
title: Measure warm threads.list, inbox.designer, and stats.get before and after, and record them in the design
status: in-progress
parent: S-0602
owner: alex
created: 2026-11-07T23:50:49Z
updated: 2026-11-08T03:23:39Z
transitions:
  - to: ready
    at: 2026-11-08T03:23:39Z
    by: agent-S-0602
  - to: in-progress
    at: 2026-11-08T03:23:39Z
    by: agent-S-0602
stream: S-0602
tags: [cli]
touches: [design/system/server-performance.md, design/system/flai-cli.md]
after: [T-2667, T-2671, T-2673, T-2676, T-2680]
---
# T-2684 Measure warm threads.list, inbox.designer, and stats.get before and after, and record them in the design

## Work

Time `threads.list`, `inbox.designer`, and `stats.get` warm on this repository, on the main branch's flai and on the story's, with their phases (`threads.read`, `messages.read`, and the reads of `stats.get`).

Each `flai hostapi --timing` call is a new process, so its reads are always cold. Take the warm figures from the second of two calls in one process: a Go program built with `go build -overlay` under `/tmp` that calls the method twice, as S-0596 did. Or take them from flai serve's log of the same methods. Note the host's load average beside the figures. Take no file under the repository for the program.

Then record the results:

- In `design/system/server-performance.md`, add a section after the S-0596 pass for cause 13. Give its cause, its fix, and the before and after of each method and phase, each read at least 80 ms less, and leave the S-0596 tables as they are: the shares on MS-0366, MS-0375, and MS-0391 give S-0602 its own cause 13 section only, kept in cause order beside the other stories' sections.
- In `design/system/flai-cli.md`, add to the Internal structure paragraph on the S-0156 store that threads, conversations, digest entries, issues, and activity documents are kept the same way, through `workitem/keptfiles.go`.

It waits for the five reader tasks, T-2667, T-2671, T-2673, T-2676, and T-2680, because the after figures need every reader moved.

## Done when

- `server-performance.md` records warm before and after figures for the three methods, each at least 80 ms less in its reads, with how they were taken.
- `flai-cli.md` names the kept files.
- `flai test design/system/server-performance.md design/system/flai-cli.md` passes.

## Notes
