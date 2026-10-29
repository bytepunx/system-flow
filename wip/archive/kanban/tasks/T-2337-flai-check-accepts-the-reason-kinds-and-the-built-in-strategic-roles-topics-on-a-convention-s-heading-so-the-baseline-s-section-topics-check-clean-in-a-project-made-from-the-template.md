---
id: T-2337
type: task
nature: improvement
title: flai check accepts the reason kinds and the built-in strategic roles' topics on a convention's heading, so the baseline's section topics check clean in a project made from the template
status: done
parent: S-0537
owner: alex
created: 2026-10-29T15:24:25Z
updated: 2026-10-29T15:42:10Z
transitions:
  - to: ready
    at: 2026-10-29T15:39:45Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:39:46Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T15:42:10Z
    by: agent-S-0537
stream: S-0537
tags: []
touches: [flai/internal/check/check.go, flai/internal/check/vocabulary_test.go, flai/internal/check/check_test.go]
after: [T-2308]
usage:
  source: log
  seconds: 144
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 12418
      cache_read: 2481441
      cache_write: 80180
      cost: 1.2005
---
# T-2337 flai check accepts the reason kinds and the built-in strategic roles' topics on a convention's heading, so the baseline's section topics check clean in a project made from the template

## Work

ADR-0217 §2. T-2310 puts heading topics such as `all, review` and `all, analysis` on the baseline conventions here and in the template. `flai check` warns `doc.topic` on a heading word that is not in `Vocabulary` (`flai/internal/check/check.go`). In a project made from the template, no story declares `review` or `analysis`, so that check fails `--strict` there, and the smoke tier with it.

Add the words the reason pack's section rule reads to `Vocabulary`: the reason kinds and the built-in strategic roles' topics, as T-2308 exports them from `flai/internal/context`. Change only `Vocabulary` in `check.go`, which S-0418 also claims; the change is agreed in their conversation. Put the test in a new `vocabulary_test.go`, not in S-0418's test files.

Waits for T-2308, which exports the words.

## Done when

- A test in `vocabulary_test.go` finds every reason kind and every built-in role's topic in `Vocabulary` for a manifest with no stories.
- `flai check --strict` gives no `doc.topic` on the conventions T-2310 tagged.
- `flai test` passes on the changed paths.

## Notes
