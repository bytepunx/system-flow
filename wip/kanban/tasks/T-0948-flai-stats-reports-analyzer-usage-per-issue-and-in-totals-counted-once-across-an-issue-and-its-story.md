---
id: T-0948
type: task
nature: improvement
title: flai stats reports analyzer usage per issue and in totals, counted once across an issue and its story
status: backlog
parent: S-0227
owner: alex
created: 2026-10-05T05:45:50Z
updated: 2026-10-05T05:45:50Z
transitions: []
stream: S-0227
tags: [flai]
touches: [flai/internal/metrics, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0944]
---
# T-0948 flai stats reports analyzer usage per issue and in totals, counted once across an issue and its story

## Work

Make `flai stats` read the issues' strategic usage and report it, apart from the agents' figures. It waits for T-0944, because what counts once depends on the carry-over it makes, and through it for T-0940.

- `flai/internal/metrics` (`usage.go`, `strategic.go`): report each issue's strategic entries per issue, beside `ItemStrategic`.
- In the strategic totals, count an issue's entries only while no story has been made from it. Tell that from the stories that link the issue (`issues.Stories`, or the Remediation line `flai issue story` adds). After that, count the story's entry, which already sits under the items and their epic.
- `flai/cmd/stats.go`: print the per-issue figures in the summary and in `--json`, and update the command's help on strategic usage.
- `design/system/metrics.md` › Strategic usage: describe the per-issue figure and the count-once rule. Regenerate `docs/users/flai-reference.md` from the help, and say it in `docs/users/flai.md`.

## Done when

- Tests pin an issue with an `analyzer` entry counted in the totals, and that same issue counted once, as its story's, after a story is made from it.
- The agents' per-model figures leave out issue usage: a test pins it.
- `design/system/metrics.md` describes it, the reference is regenerated, and `scripts/flai-test.sh` and `flai check --strict` pass.

## Notes
