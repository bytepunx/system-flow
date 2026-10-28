---
id: D-0282
title: "S-0369 plan: flai guard refuses served agents --by and FLAI_AGENT on flai (new ADR), five tasks in four layers: Recommendation: approve the plan with one narrowing"
kind: answer
sources: [TH-0418]
topics: [cli, conventions, flai, go, code, all]
paths: [design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md, design/issues/summary.md, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md]
date: 2026-10-12
status: absorbed
absorbed_by: "design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md#decision"
asked_again: []
cited: []
---

# D-0282 S-0369 plan: flai guard refuses served agents --by and FLAI_AGENT on flai (new ADR), five tasks in four layers: Recommendation: approve the plan with one narrowing

The operator, alex, answered TH-0418 at 2026-10-12T18:00:15Z.

What was asked:

```text
**Plan:** S-0369 closes I-0091 with a new ADR and a guard rule. The rule refuses a story's agent, the planner, and the analyzer any flai command that names its own author. It also fixes the fallback in `threadAuthor`. There are five tasks in four layers.
```

The answer:

```text
**Recommendation:** approve the plan with one narrowing. The guard refuses a served agent (a story's agent, the planner, or the analyzer) a flai command whose `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION` names anyone but that agent's own name, as flai serve set it. It does not refuse those flags or variables as such. Keep T-1466 in this story.

## Why this needs you

- **A new rule on a story's agent.** ADR-0060 lets a story's own agent's calls pass the guard untouched. This would be the guard's first rule on it, a decision not yet recorded. T-1465 records it, but the choice is yours.
- **It conflicts with a convention as drafted.** `design/conventions/tooling.md` § Rules says "Set `FLAI_AGENT` and `FLAI_SESSION` at session start so narrative entries record who wrote them". Refusing any `FLAI_AGENT` set for flai would refuse agents that follow it. Refusing only a name other than the agent's own still blocks `FLAI_AGENT=alex flai …`, which is I-0091's case.
- **Scope.** Covering the planner and the analyzer, and T-1466's `threadAuthor` fallback, both go beyond I-0091. T-1466 is small, in the same file, and makes the docs true, so I would keep it. Dropping it is your call.

## The rest fits

- Five tasks in four layers, all file touches.
- Figures: 45m (flai's 21m raised as S-0368's was) and 10 USD a week from my input on TH-0416.
- The overlap hold with S-0368 keeps the two apart.

I leave S-0369 a draft until you answer.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/conventions/tooling.md § Rules
```

Source: TH-0418.
