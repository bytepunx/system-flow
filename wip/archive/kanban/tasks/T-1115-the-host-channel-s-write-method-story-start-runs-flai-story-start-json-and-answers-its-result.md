---
id: T-1115
type: task
nature: improvement
title: The host channel's write method story.start runs flai story start --json and answers its result
status: done
parent: S-0274
owner: alex
created: 2026-10-06T22:53:44Z
updated: 2026-10-07T08:05:54Z
transitions:
  - to: ready
    at: 2026-10-07T08:00:04Z
    by: agent-S-0274
  - to: in-progress
    at: 2026-10-07T08:00:04Z
    by: agent-S-0274
  - to: done
    at: 2026-10-07T08:05:54Z
    by: agent-S-0274
stream: S-0274
tags: [hostapi, go]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, CLAUDE.md, design/conventions/session-start.md, design/conventions/work-management.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/conventions.md, docs/users/flai.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/CHANGELOG.md, template/root/CLAUDE.md.tmpl, template/root/design/conventions/session-start.md, template/root/design/conventions/work-management.md]
after: [T-1109]
usage:
  source: log
  seconds: 350
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 20693
      cache_read: 2986764
      cache_write: 95443
      cost: 1.5417
---
# T-1115 The host channel's write method story.start runs flai story start --json and answers its result

## Work

Add `story.start` to the write methods in `flai/internal/hostapi/writes.go`, next to `item.move`. It takes `story`, an optional `budget`, and the agent the channel acts for. Like its siblings, it builds the arguments `story start <S-nnnn> --json [--budget <b>]` and runs flai in the project's folder. A refusal exit maps to `Refused` with flai's reason, and the JSON result is answered as it is.

It waits for T-1109, the command it runs.

## Done when

- `writes_test.go` covers the arguments built with and without `budget`, a story ID that is not one refused before flai runs, and a refusal exit mapped to `Refused`.
- `scripts/flai-test.sh` passes.

## Notes
