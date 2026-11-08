---
id: T-2664
type: task
nature: improvement
title: Measure an acceptance before and after in a scratch copy and record it in server-performance.md
status: backlog
parent: S-0609
owner: alex
created: 2026-11-07T23:49:51Z
updated: 2026-11-07T23:49:51Z
transitions: []
stream: S-0609
tags: [cli]
touches: [design/system/server-performance.md]
after: [T-2636, T-2649, T-2658]
---
# T-2664 Measure an acceptance before and after in a scratch copy and record it in server-performance.md

## Work

- Build flai at the story's merge base and at its head into a scratch location (`bin/` or the session scratchpad, never `~/.flai`).
- In a scratch copy of this repository, accept the same story with markdown in both the uncommitted set and the branch's set with each build, three runs each, with `markdownlint-cli2` installed by `scripts/install-tools.sh` and once without it, so the `npx` fallback is measured too. Time `flai accept` and one `flai task done`, and count the lint calls and the `git status --porcelain` calls, as cause 19 in `design/system/server-performance.md` counted them.
- Accept once with unfixable markdown in each set, with each build, and confirm the refusal is the same, as ADR-0140 says.
- Record the before and after, and how they were measured, in `design/system/server-performance.md`, in a section for S-0609 after the second pass, as the earlier stories' sections are written; mark cause 19 as addressed there.
- Waits for T-2636, T-2649, and T-2658: it measures what all three change.

## Done when

- `design/system/server-performance.md` holds the measured before and after for `flai accept` and `flai task done`, with the lint and git status counts, and the acceptance is at least 4 s faster with `markdownlint-cli2` installed, or the section says by how much it fell short and why.
- The refusal for unfixed markdown is the same before and after.
- `flai test design/system/server-performance.md` passes.

## Notes
