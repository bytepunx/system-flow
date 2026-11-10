---
id: I-0291
title: "Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand"
class: efficiency
status: closed
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-10T05:52:36Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0291 Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand

## Description
Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: orchestration.release has value 1000 and no policy, read as judgement, which the publish formula's evaluate step can never meet, and release_publish refuses while the formula exists (strategic-agents.md § Publishing a release).

## Impact
serve journal push entries since 2026-11-02T05:31 (last mcp.release_publish by the orchestrator): six publish.run by alex, 2026-11-02 to 2026-11-07; T-2548 blocked 2026-11-07T08:06:26Z on flai release --evaluate --require-met exit 2; TH-0811 open since 2026-11-02 with one entry; 67 release-due events in .flai-cache/events.jsonl

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0625 remediates this issue, created from it at 2026-11-07T23:40:15Z.
Closed 2026-11-10T05:52:36Z: S-0625 (ADR-0273): flai release --evaluate takes --judged "<reason>", which under judgement meets the policy while a story is accepted and not yet released and whole_epics holds nothing back; the publish formula declares the input judgement, which its evaluate step passes as --judged; and the orchestrator, under judgement, runs flai formula run publish --ready --autocommit --input judgement="<its reasoning>". TestPublishFormulaUnderJudgement runs the real formula's step 2 and reproduces the cause (exit 2 with no judgement) and the fix (exit 0 with one); TestReleaseEvaluateRequireMet and TestEvaluateJudged cover the flag.
