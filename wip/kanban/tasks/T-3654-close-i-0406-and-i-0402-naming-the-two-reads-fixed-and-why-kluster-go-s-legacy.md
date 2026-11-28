---
id: T-3654
type: task
nature: improvement
title: Close I-0406 and I-0402, naming the two reads fixed and why kluster.go's legacy send-back check keeps the baseline
status: backlog
parent: S-0851
owner: alex
created: 2026-11-28T00:03:57Z
updated: 2026-11-28T00:03:57Z
transitions: []
stream: S-0851
tags: [flai, issues]
touches: [design/issues/I-0406-flai-serve-still-reads-three-state-questions-by-the-baseline-workflow-directly.md, design/issues/I-0402-flai-serve-stops-the-dev-server-of-a-story-in-a-state-a-project-s-workflow-adds.md, design/issues/summary.md]
after: [T-3652, T-3653]
---
# T-3654 Close I-0406 and I-0402, naming the two reads fixed and why kluster.go's legacy send-back check keeps the baseline

## Work

I-0406 names three reads. T-3652 fixes `acceptance()` and T-3653 fixes `StoryLeft`. The third is `klusterKept` in `flai/internal/serve/kluster.go`. It reads `baselineFlow.Gives(..., workflow.EffectSendBack)` only for a send-back record with no `To`. Its comment says such a record was written by a flai serve that ran only the baseline, before the record kept its state. A record with `To` is already judged by that state, and `kluster_test.go` covers it on a renamed board (`doing`). So the baseline is deliberate there: it judges a legacy record by the board that wrote it. No code changes for it.

- Check that claim before closing: read where the send-back record gained `To` (`git log -S 'SentBack.To'` or the field's first commit in `flai/internal/serve/kluster.go`). If a flai that wrote records without `To` could already run a project's workflow file, say so on the plan's thread on S-0851 instead of closing. That read would then need `repo.Flow()` too.
- Run `flai issue close I-0406 --reason "<reason>"`. The reason names the two reads moved to `repo.Flow()`, with their task IDs, and why `klusterKept`'s legacy fallback and `builtinTriggersIn`'s translation from baseline names in `orchestrate_reasons.go` keep the baseline.
- Run `flai issue close I-0402 --reason "<reason>"`, naming T-3653's fix and its test. It is the defect the `StoryLeft` read caused.
- Let flai regenerate `design/issues/summary.md`.

It waits for T-3652 and T-3653, because the close reasons name what they changed and their tests.

## Done when

- I-0406 and I-0402 are closed, each with a reason naming the fix and its test.
- `flai check --strict` reports nothing for the two issue files or `design/issues/summary.md`.

## Notes
