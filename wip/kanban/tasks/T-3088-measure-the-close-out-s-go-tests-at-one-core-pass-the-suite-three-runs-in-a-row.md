---
id: T-3088
type: task
nature: improvement
title: Measure the close-out's Go tests at one core, pass the suite three runs in a row, and set the timeouts again
status: backlog
parent: S-0722
owner: alex
created: 2026-11-11T06:05:38Z
updated: 2026-11-11T06:05:38Z
transitions: []
stream: S-0722
tags: [tests, close-out]
touches: [system-flow.yaml, scripts/integration.sh, scripts/race.sh]
after: [T-3087]
---
# T-3088 Measure the close-out's Go tests at one core, pass the suite three runs in a row, and set the timeouts again

## Work

Measure the "after" figures as the first task measured the "before": for a change that touches `flai/internal/mcpserver`, the `go-test` tier's run and `scripts/integration.sh`, wall time, user CPU, and peak cores, under `/usr/bin/time -v` or with the load sampled. Write both sets side by side in the narrative's `## Decisions`. One close-out's Go tests should hold about one core; if they hold more, find what does (the compiler, a test that starts its own processes) and bound it.

Run the full suite at one core, `scripts/integration.sh`, three times in a row. A test that fails only at one core is fixed when the fix is small, most likely a timing test with short heartbeats or a held wait; name each in `## Decisions`. When a fix is not small, stop and say which test and why on a thread on S-0722; do not skip it and do not raise the cap. The files such a fix changes cannot be named yet: `flai task done` widens this task's touches with them.

Set the timeouts again from the measured figures: the `-timeout 20m` of the `go-test` tier in `system-flow.yaml` and of `scripts/integration.sh`, and the comment there citing S-0507 and I-0225, which sized them for `-race`. Keep `20m` in `scripts/race.sh`, which still runs `-race`, and move that comment there. Where a timeout stays, restate the reason from the new figures.

Waits for the second task: it measures the tiers as both earlier tasks leave them, and edits `race.sh`, which the second task adds.

## Done when

- The narrative holds the before and after figures, and after, one close-out's Go tests hold about one core.
- `scripts/integration.sh` passed three runs in a row at one core, or a thread on S-0722 names the test that does not and why.
- Each timeout is set from the figures or kept with its reason restated, and `flai check --strict` is clean.

## Notes
