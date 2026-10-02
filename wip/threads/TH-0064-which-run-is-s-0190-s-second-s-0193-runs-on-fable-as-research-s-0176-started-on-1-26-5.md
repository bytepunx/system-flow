---
id: TH-0064
title: "Which run is S-0190's second: S-0193 runs on Fable as research, S-0176 started on 1.26.5"
anchor:
  path: wip/kanban/stories/S-0190-measure-what-delegating-costs-with-cheaper-sub-agents-and-one-verifier-run-before-review.md
  item: S-0190
status: resolved
participants: [agent-S-0190, alex]
created: 2026-10-02T12:13:37Z
updated: 2026-10-02T12:39:40Z
---

# TH-0064 Which run is S-0190's second: S-0193 runs on Fable as research, S-0176 started on 1.26.5

On wip/kanban/stories/S-0190-measure-what-delegating-costs-with-cheaper-sub-agents-and-one-verifier-run-before-review.md.

## Entries

### 2026-10-02T12:13:37Z agent-S-0190
Recommendation: measure S-0194 now, and take as the second run the next story that runs on claude-opus-5-5 and changes code (S-0178, S-0191, or S-0192, all held until S-0193 ends). Report S-0193 and S-0176 apart, in the section's prose.

Why: two stories have run on a flai with S-0189, and neither compares with S-0188's comparables, which all ran on Opus and changed code.

- S-0194 (feature, flai 1.26.6) qualifies. It cost 7.32 US dollars. Its sub-agents cost 0.78: explorer on haiku 0.29, two verifiers on sonnet 0.50. Its agent ran no whole-suite test run itself; both close-out runs were the verifier's.
- S-0193 is running now on 1.27.0, but its agent is claude-fable-5-1, not Opus. As research, its result is threads and stories, not code, so it has no size to compare and no non-delegating research run to compare against. Its cost also includes time spent waiting on your answers.
- S-0176 used the S-0189 definitions (haiku and sonnet sub-agents), but flai serve started it at 11:36:28Z on 1.26.5, seven minutes before 1.26.6 was installed. So it had S-0189's models without S-0189's prompt. It was also its own experiment, with Opus task sub-agents (6.91 of its 23.44) and blind reviews.

Alternatives:
(b) Measure S-0194 and S-0193 now, with S-0193 caveated as Fable and research. This finishes sooner, but the second row would not answer the question.
(c) Measure S-0194 and S-0176 now, with S-0176 caveated. Its total cost is dominated by work that the S-0189 changes don't address.

With (a), I'd block S-0190 until that story ends. That frees the in-progress slot, and flai restarts me when I'm unblocked or this thread is answered.

### 2026-10-02T12:29:20Z alex
a

### 2026-10-02T12:29:41Z agent-S-0190
Going with (a). I've blocked S-0190 until the next story that runs on claude-opus-5-5 and changes code has ended. Its log will be the second run. S-0194 is already measured; the numbers are in the narrative's Decisions. S-0193 and S-0176 will be reported apart.

### 2026-10-02T12:39:40Z alex
Resolved.
