---
id: D-0374
title: "Plan for S-0467: three remedies for I-0150 in seven tasks over five layers: Recommended: yes, as its own story, not inside S-0467"
kind: answer
sources: [TH-0598]
topics: [orchestration, flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_own_thread_test.go, flai/cmd/accept_conflict_test.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md, design/issues/summary.md, design/adrs/0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-merge-queue-s-does.md, flai/cmd/accept_orchestrator_dirty_test.go, docs/users/flai-reference.md]
date: 2026-10-23
status: absorbed
absorbed_by: "design/adrs/0206-flai-serve-confirms-an-orchestrator-recommendation-the-operator-left-unanswered.md#decision"
asked_again: []
cited: []
---

# D-0374 Plan for S-0467: three remedies for I-0150 in seven tasks over five layers: Recommended: yes, as its own story, not inside S-0467

The operator, alex, answered TH-0598 at 2026-10-23T04:56:05Z.

What was asked:

```text
Plan by planner-S-0467. The story is finalized, so its words stand; I added tags `flai` and `template` (flai check requires a component tag), topic `orchestration`, 34 file touches, a forecast of 2h (flai's 54m raised, reason under `### Planning`), and the cost of delay value 14400 USD a week from the operator's 96h input.
```

The answer:

```text
Recommended: yes, as its own story, not inside S-0467. S-0467 stays as planned, and its remedy 3 still needs your A or B.

**A shape that fits what exists:**

- **A new permission,** off by default: `orchestration.permissions.confirm_after`, a duration such as `48h`.
- **What it covers.** When a recommendation I posted has waited that long with no entry from you, flai confirms it as the answer. The entry is marked as confirmed by timeout, so it shows as such in the thread, the log, and `flai stats`, beside the confirmations you make (ADR-0091).
- **What it never covers:**
  - money: cost of delay inputs, estimates, or spend;
  - host changes: installs, keys, host actions;
  - permission prompts;
  - `.claude/` writes;
  - a thread you have written on since the recommendation.

  Those keep waiting for you.
- **Who acts.** flai serve does the confirming, on a schedule. I would not resolve your threads myself, so the rule stays in one place that the guard and the stats can see.

**Why exclusions, not a blanket timer.** Today most threads awaiting you are money or host questions (TH-0467, TH-0532, TH-0533, TH-0587, TH-0590). The ones a timer would safely settle are design choices with a clear recommendation, like TH-0583, TH-0588, and this one.

It needs an ADR, since it changes who may answer a thread (it refines ADR-0090). Say "draft it" and I'll record it as an issue so that the issue step can make the story.

Source: design/adrs/0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md
```

Source: TH-0598.
