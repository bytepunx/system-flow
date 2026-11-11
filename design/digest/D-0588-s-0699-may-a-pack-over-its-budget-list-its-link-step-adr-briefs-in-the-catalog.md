---
id: D-0588
title: "S-0699: may a pack over its budget list its link-step ADR briefs in the catalog, reversing TH-0032's \"briefs are never cut\": Recommended: Q1 no, Q2 yes"
kind: answer
sources: [TH-1023]
topics: [flai, cli, go, code, all]
paths: [flai/internal/topics/story.go, flai/internal/topics/story_test.go, flai/internal/context/context_test.go, flai/internal/mcpserver/prime_test.go, design/system/agent-context.md, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md]
date: 2026-11-11
status: absorbed
absorbed_by: "design/system/agent-context.md#delivered-in-parts"
asked_again: []
cited: [TH-1023, TH-1078]
---

# D-0588 S-0699: may a pack over its budget list its link-step ADR briefs in the catalog, reversing TH-0032's "briefs are never cut": Recommended: Q1 no, Q2 yes

The operator, alex, answered TH-1023 at 2026-11-11T08:59:12Z.

What was asked:

```text
**Recommendation: yes to both questions below.** Answering no to either drops that task. The story then fixes only the topic sources.
```

The answer:

```text
**Recommended:** Q1 no, Q2 yes.

**Q1, link-step ADR briefs to the catalog: no.** You chose "briefs are never cut" twice: on TH-0032, and again on TH-0805, recorded in ADR-0241 as option 1. Nothing in I-0261 shows a cost that you did not have before you then. Keep every brief. S-0699 then has no ADR task for it.

**Q2, a count of a brief's other reasons in place of `also`: yes.** It saves about 130 KB a pack, which is 39% of the brief bytes. It cuts no brief and no document. A loaded item keeps its full list. It changes one line of `design/system/agent-context.md`, not a decision you recorded. The alternative, keeping the list capped at three reasons, saves less and still repeats.

With Q1 no and Q2 yes, S-0699 has three tasks: T-2992, the `also` count task, and T-2996.

Source: design/adrs/0241-a-pack-whose-conventions-alone-exceed-its-budget-still-loads-what-is-named-and.md
```

Source: TH-1023.
