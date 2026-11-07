---
id: T-2452
type: task
nature: experiment
title: The Go proof story is worked to review in its container on this repository's image layer, its close-out passing every tier the manifest's tests select inside the container, and each limitation met is an issue
status: in-progress
parent: S-0565
owner: alex
created: 2026-11-03T01:25:10Z
updated: 2026-11-07T06:18:21Z
transitions:
  - to: ready
    at: 2026-11-07T06:18:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T06:18:21Z
    by: system-flow
stream: S-0565
tags: [docker]
touches: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/issues/summary.md]
after: [T-2451]
---
# T-2452 The Go proof story is worked to review in its container on this repository's image layer, its close-out passing every tier the manifest's tests select inside the container, and each limitation met is an issue

## Work

Waits for T-2451: the Go run comes once the first run has shown the harness works, so a limitation the first one found is known before the second meets it, and both write the results document and the issues summary.

Hold `wait_for_events` until the Go proof story reaches review, or its agent fails past its automatic restarts. Check that its container ran this repository's image layer, built from `flai/images/system-flow/Dockerfile` and named in the manifest as S-0510 has it. Read its close-out's verdict and tiers with `flai verify S-nnnn --last`.

Under `flai serve` that close-out is a job on the host (ADR-0229), not a run inside the container. So unless the answer on TH-0821 named a story that moves the job into the container, also run the close-out once by hand inside the layered image on the Go story's worktree, with S-0520's mounts in place: `scripts/agent-container.sh S-nnnn -- scripts/close-out.sh S-nnnn`. Compare the tiers that ran with those the manifest's `tests` select for the story's changed paths. A tier that was skipped, left to the merge queue, or run on the host is not a pass inside the container: record which.

Then run the host's flaiover tier on the main checkout once more, through `flai test`, to see whether I-0259's musl bindings came back. Record the runs, the tiers, and their outcomes in the results document. File each limitation met with `flai issue new`, as T-2451 does, and bump an issue met again with `flai issue bump`.

## Done when

- The Go proof story is in review, worked by its agent in its container on this repository's layer.
- `scripts/close-out.sh` passed inside the container with every tier the manifest's `tests` select, and the results document lists those tiers and their outcome, apart from the host job's.
- Each limitation met is an issue under `design/issues`, named by ID in the results document.

## Notes
