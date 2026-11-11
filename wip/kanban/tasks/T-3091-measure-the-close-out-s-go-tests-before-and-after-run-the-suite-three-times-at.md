---
id: T-3091
type: task
nature: improvement
title: Measure the close-out's Go tests before and after, run the suite three times at one core, and set the timeouts again from the figures
status: backlog
parent: S-0722
owner: alex
created: 2026-11-11T06:06:49Z
updated: 2026-11-11T06:06:49Z
transitions: []
stream: S-0722
tags: []
touches: [scripts/integration.sh, system-flow.yaml]
after: [T-3090]
---
# T-3091 Measure the close-out's Go tests before and after, run the suite three times at one core, and set the timeouts again from the figures

## Work

- Measure on this host, for a change that touches `flai/internal/mcpserver`, the close-out's Go tests (the `integration` tier, which covers `go-test`): wall time, user CPU, and peak cores, before (`-race`, every core) and after (one core). Record the figures and the method in the narrative's Decisions.
- Run the whole Go suite at one core three times in a row. Fix a test that only passes with more cores when the fix is small; otherwise stop and ask on a thread.
- Set the 20-minute timeouts sized for `-race` again from the measured figures, or keep them with the reason restated, in `integration.sh`'s comment and the tiers.
- Waits for T-3090: it measures the commands T-3090 writes.

## Done when

- The narrative holds the before and after figures, and three passing runs at one core.

## Notes
