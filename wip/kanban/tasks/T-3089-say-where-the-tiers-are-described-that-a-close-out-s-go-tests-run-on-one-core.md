---
id: T-3089
type: task
nature: improvement
title: Say where the tiers are described that a close-out's Go tests run on one core without -race, and when to run flai test race
status: backlog
parent: S-0722
owner: alex
created: 2026-11-11T06:05:46Z
updated: 2026-11-11T06:05:46Z
transitions: []
stream: S-0722
tags: [tests, close-out, docs]
touches: [design/system/workflow.md, design/system/project-manifest.md, design/system/devex.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/tech/go.md, scripts/README.md, docs/operators/index.md]
after: [T-3087]
---
# T-3089 Say where the tiers are described that a close-out's Go tests run on one core without -race, and when to run flai test race

## Work

Say in each place that describes this repository's test tiers what a close-out's Go tests now run: no `-race`, `-p 1`, one core, by design, so that several close-outs can run at once; that the race detector is the manual `race` tier, `flai test race` or `make race`, which `make flai-test` runs and no close-out or batch does; and when to run it: after a change to concurrent code (goroutines, locks, channels, shared caches) and before a release.

- `design/system/workflow.md`: what a close-out runs, near the close-out and merge queue paragraphs.
- `design/system/project-manifest.md`: the sentence that quotes this repository's `go-test` command, and the timeout paragraph that cites `-race` (I-0227).
- `design/system/devex.md`: the test tiers row.
- `design/conventions/code-quality.md` and `design/conventions/tooling.md`, below the marker only: the tiers they list, with `race` beside `gateway`.
- `design/tech/go.md`: the `go test` row.
- `scripts/README.md`: `test.sh` and `integration.sh` as they now run, and a row for `race.sh`.
- `docs/operators/index.md`: the close-out paragraph that lists the tiers.

The template's copies do not mention `-race` and its `test.sh` and `integration.sh` are placeholders: check that this still holds, and change the template's copy of anything changed here that it carries. Say what you found in `## Decisions`.

Waits for the second task: it describes the `race` tier that task adds. It runs beside the third task, which touches none of these files.

## Done when

- Each file above says what a close-out's Go tests run, that it is one core by design, and when to run `flai test race`.
- The template is checked, and changed where it carries a copy.
- `flai test` on the changed files passes and `flai check --strict` is clean.

## Notes
