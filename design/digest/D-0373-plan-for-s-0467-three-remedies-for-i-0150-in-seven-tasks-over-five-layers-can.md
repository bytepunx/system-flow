---
id: D-0373
title: "Plan for S-0467: three remedies for I-0150 in seven tasks over five layers: can we introduce an additional setting that would allow the orchestrator or another agent to make decisions on unresolved…"
kind: recommendation
sources: [TH-0598]
topics: [orchestration, flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_own_thread_test.go, flai/cmd/accept_conflict_test.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md, design/issues/summary.md, design/adrs/0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-merge-queue-s-does.md, flai/cmd/accept_orchestrator_dirty_test.go, docs/users/flai-reference.md]
date: 2026-10-23
status: absorbed
absorbed_by: "design/adrs/0206-flai-serve-confirms-an-orchestrator-recommendation-the-operator-left-unanswered.md#context"
asked_again: []
cited: []
---

# D-0373 Plan for S-0467: three remedies for I-0150 in seven tasks over five layers: can we introduce an additional setting that would allow the orchestrator or another agent to make decisions on unresolved…

The operator, alex, answered TH-0598 at 2026-10-23T04:46:55Z.

What was asked:

```text
Plan by planner-S-0467. The story is finalized, so its words stand; I added tags `flai` and `template` (flai check requires a component tag), topic `orchestration`, 34 file touches, a forecast of 2h (flai's 54m raised, reason under `### Planning`), and the cost of delay value 14400 USD a week from the operator's 96h input.
```

The recommendation it overruled:

```text
This story changes the rules on my own acceptances, so weigh my view with that in mind. Recommended for the operator: take remedies 1 and 2 as proposed, and take remedy 3 only in the stricter form.

1. **Dirty main: take it.** It does what the merge queue already does under ADR-0178. Uncommitted files stay out of the commit, and only a change that would overwrite one is refused. It widens nothing the queue does not already do.
2. **A rebase it cannot settle: take it.** Sending the story back with the paths starts its agent, which is the party that can resolve the conflict. Today the story sits in review until you step in.
3. **My own thread not gating: take the stricter alternative.** The rule would be: I must resolve my own thread before the acceptance, and a thread anyone else wrote on still holds the story. Under the proposed form, a concern I raise could never stop a story. TH-0589 is the case in point. It held S-0483 because the operator's instruction on TH-0578 had not been carried out, and that hold is what brought it to you. A thread I can resolve myself still ends the "gate it raised itself" wait, and it keeps a concern visible as a decision rather than letting it lapse at archive.

The value is 14,400 USD a week, from your 96h input, which puts this story near the top of anything ordered by value.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md
```

The operator's answer:

```text
can we introduce an additional setting that would allow the orchestrator or another agent to make decisions on unresolved threads after a time threshold?
```

Source: TH-0598.
