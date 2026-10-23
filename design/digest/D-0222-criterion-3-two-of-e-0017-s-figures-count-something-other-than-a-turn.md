---
id: D-0222
title: "Criterion 3: two of E-0017's figures count something other than a turn: Recommendation: keep the classes as defined (alternative 1), and have the operator reword criterion 3 before it is ticked…"
kind: answer
sources: [TH-0278]
topics: [automation, conventions, mcp, hostapi, flai, cli, go, code, all]
paths: [flai/internal/usage/log.go, flai/internal/usage/usage.go, flai/internal/usage/log_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/usage/turns.go, flai/internal/usage/turns_test.go, flai/internal/usage/usage_test.go, flai/internal/workitem/fields_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/store_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/metrics/metrics.go, flai/internal/metrics/turns.go, flai/internal/metrics/turns_test.go, design/system/work-hierarchy.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0222 Criterion 3: two of E-0017's figures count something other than a turn: Recommendation: keep the classes as defined (alternative 1), and have the operator reword criterion 3 before it is ticked…

The operator, alex, answered TH-0278 at 2026-10-07T14:56:27Z.

What was asked:

```text
Recommendation: accept criterion 3 as met by the comparison in S-0293's notes, keeping the classes as defined. Two of the four classes agree within ten percent; the other two differ because the epic's figures count something other than a turn of that class.
```

The answer:

```text
Recommendation: keep the classes as defined (alternative 1), and have the operator reword criterion 3 before it is ticked, rather than tick it as written.

Why keep the classes: the empty-wake class is ADR-0105's definition, which ADR-0105 says S-0293 takes over. Counting every wait to reach 650 would contradict it. Counting quoted `make test` text as test runs would count commands that run nothing.

Why reword first: criterion 3 says the counts "agree with the epic's evidence within ten percent". Three of four classes do not: test runs −14%, empty wakes −70%, hand edits +70%. ADR-0089 has the agent tick a criterion only once it has verified it. A criterion that cannot be verified stays unticked, with the reason in the notes. Ticking this one as written would record agreement the numbers do not show.

Suggested wording for criterion 3: "Run over the logs to 2026-10-04, the counts are compared with the epic's evidence, and the story records the comparison and explains each class that differs by more than ten percent." The notes' "Comparison with E-0017's evidence" section already meets that.

Changing a finalized story's criterion is a change of scope, so this is the operator's call.

Source: design/adrs/0089-acceptance-criteria-are-ticked-through-flai-by-number-by-the-story-s-agent-once.md § Decision
```

Source: TH-0278.
