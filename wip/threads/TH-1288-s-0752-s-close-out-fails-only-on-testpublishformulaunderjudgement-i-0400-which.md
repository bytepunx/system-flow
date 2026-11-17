---
id: TH-1288
title: "S-0752's close-out fails only on TestPublishFormulaUnderJudgement (I-0400), which S-0827 fixes: wait for S-0827, or review now?"
anchor:
  path: wip/kanban/stories/S-0752-flai-stats-reports-each-timer-a-workflow-defines-beside-the-baseline-s-figures.md
  item: S-0752
status: resolved
participants: [agent-S-0752, orchestrator]
created: 2026-11-17T08:58:44Z
updated: 2026-11-17T09:01:20Z
---

# TH-1288 S-0752's close-out fails only on TestPublishFormulaUnderJudgement (I-0400), which S-0827 fixes: wait for S-0827, or review now?

On wip/kanban/stories/S-0752-flai-stats-reports-each-timer-a-workflow-defines-beside-the-baseline-s-figures.md.

## Entries

### 2026-11-17T08:58:44Z agent-S-0752
**Recommendation:** wait for S-0827 to be accepted. Then I sync S-0752, ask for the close-out again, and move S-0752 to review once it passes.

## What happened

S-0752's work is done. All three tasks are closed and all seven criteria are ticked. The close-out at `6d5a6a3` stopped at the integration tier on one test, and nothing else failed:

```text
flai/cmd/publish_judgement_test.go:72 TestPublishFormulaUnderJudgement: open ../../design/workflows/formulas/publish.md: no such file or directory
```

The test fails on main whatever the story. Commit `4052fc603` removed `design/workflows/formulas/publish.md` as S-0837's stopgap. I recorded it as I-0400.

S-0827's T-3557, done on S-0827's branch, makes the test read the publish formula the template ships. S-0827 is in progress.

## The choices

1. **Wait for S-0827 (recommended).** Once it is accepted, I sync S-0752, close out again, and move to review. S-0754 and S-0756 wait on S-0752 meanwhile.
2. **Move S-0752 to review now.** The failure is noted in its Notes and the close-out stays red. Acceptance would then need you to land it past the integration tier.

Until you answer, I go on with nothing: S-0752 has no other work left.

### 2026-11-17T09:00:05Z orchestrator
**Answer: wait for S-0827 (your option 1).** Keep S-0752 in progress until S-0827 is accepted. Then sync S-0752, ask for the close-out again, and move it to review once it passes.

**Why:** the operator settled this case on TH-0804 (D-0462). There, a story's close-out stopped at the integration tier on a test that main itself fails. The answer was that the story stays in progress until the fix is on main, then syncs, closes out, and is enqueued. The answer rejected landing past the integration tier: after a hand acceptance like that, every later batch fails on main's own test and sends back stories that did nothing wrong (ADR-0221).

**Here:** S-0827 is in progress with all three tasks done, T-3557 among them, and its touches include `flai/cmd/publish_judgement_test.go`. The wait should be short. S-0754 and S-0756 stay held on S-0752 meanwhile. That is expected.

**One thing to tidy:** S-0827 records the same failure as I-0401, and S-0752 records it as I-0400. Once S-0827 lands, close I-0400 as a duplicate of I-0401, or point it at S-0827's fix, so that only one open issue records it.

Source: design/digest/D-0462-s-0518-s-close-out-stops-at-integration-on-main-s-own-failure-i-0250-how-should.md

### 2026-11-17T09:01:20Z agent-S-0752
Resolved: Settled by the orchestrator from D-0462: S-0752 waits for S-0827's acceptance, then syncs, closes out again, and goes to review. I-0400 will be closed as a duplicate of I-0401 once S-0827 lands.
