---
id: T-2802
type: task
nature: improvement
title: Point the design documents at the files the split code moved to, and close I-0298
status: backlog
parent: S-0631
owner: alex
created: 2026-11-09T22:05:33Z
updated: 2026-11-09T22:05:33Z
transitions: []
stream: S-0631
tags: [cli, docs]
touches: [design/system/flai-cli.md, design/system/server-performance.md, design/system/design-system-tool.md, design/system/agent-adapters.md, design/system/strategic-agents.md, design/system/agent-coordination.md, design/system/flaiover-dashboard.md, design/system/agent-prose.md, design/system/release-signing.md, design/issues/I-0298-five-source-files-under-flai-are-over-2-000-lines.md, design/issues/summary.md]
after: [T-2801]
---
# T-2802 Point the design documents at the files the split code moved to, and close I-0298

## Work

The design documents name the six split files by path, often with the function they mean. Each reference that now points to a file the code left must name the file it moved to. Each split keeps its file's name, so a reference to code that stayed needs no change.

The plan found references in:

| Document | Lines (2026-11-09) |
|----------|--------------------|
| `design/system/flai-cli.md` | 62, 97–100, 111, 117, 142–143, 261, 323, 379, 693, 716, 720, 743 |
| `design/system/server-performance.md` | 274, 359, 360, 363, 546, 656 |
| `design/system/design-system-tool.md` | 90, 139, 180, 318, 451, 455 |
| `design/system/agent-adapters.md` | 44, 72, 137, 140, 803 |
| `design/system/strategic-agents.md` | 27, 474 |
| `design/system/agent-coordination.md` | 19, 29 |
| `design/system/flaiover-dashboard.md` | 261 |
| `design/system/agent-prose.md` | 487 |
| `design/system/release-signing.md` | 242 |

Find them again with a grep for each basename and its folder, since lines move. Bump `updated` on each document you change.

Leave alone what records the past:

- the accepted ADRs (0208, 0220, 0230, 0252, 0253);
- the digest;
- the analysis report;
- the other issues.

Add a short paragraph to `flai-cli.md`'s section on the code layout, if it has one, saying that no non-test Go file is over 2,000 lines and that `TestNoSourceFileIsOver2000Lines` holds it. Otherwise add it where the document describes the packages.

Then close the issue with `flai issue close I-0298 --reason "<reason>"`. The reason names:

- the six files split;
- the size test and its file;
- the nolint reason given in `flai/internal/watch/watch.go`.

It waits for T-2801, because the reason names its test, and the six splits are done by then, so the new file names are final.

## Done when

- A grep for each of the six basenames in `design/system/` finds no reference to code that is no longer in that file.
- I-0298 is closed, and `design/issues/summary.md` lists it so.
- `flai test` on the changed documents passes the markdown lint.

## Notes
