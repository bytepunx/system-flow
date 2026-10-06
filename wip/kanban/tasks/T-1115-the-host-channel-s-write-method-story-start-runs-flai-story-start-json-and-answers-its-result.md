---
id: T-1115
type: task
nature: improvement
title: The host channel's write method story.start runs flai story start --json and answers its result
status: backlog
parent: S-0274
owner: alex
created: 2026-10-06T22:53:44Z
updated: 2026-10-06T22:53:44Z
transitions: []
stream: S-0274
tags: [hostapi, go]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1109]
---
# T-1115 The host channel's write method story.start runs flai story start --json and answers its result

## Work

Add `story.start` to the write methods in `flai/internal/hostapi/writes.go`, next to `item.move`. It takes `story`, an optional `budget`, and the agent the channel acts for. Like its siblings, it builds the arguments `story start <S-nnnn> --json [--budget <b>]` and runs flai in the project's folder. A refusal exit maps to `Refused` with flai's reason, and the JSON result is answered as it is.

It waits for T-1109, the command it runs.

## Done when

- `writes_test.go` covers the arguments built with and without `budget`, a story ID that is not one refused before flai runs, and a refusal exit mapped to `Refused`.
- `scripts/flai-test.sh` passes.

## Notes
