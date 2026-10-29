---
id: D-0424
title: "Criterion 2: a real thread's pack is 31% of the general pack, under a quarter only for the test's small fixture thread: Recommended: trim the pack as S-0537's agent proposes, and tick criterion 2 on…"
kind: answer
sources: [TH-0753]
topics: [orchestration, cli, conventions, template, mcp, dashboard, flai, go, code, all]
paths: [flai/internal/context, flai/internal/mcpserver, flai/cmd, design/system/flai-cli.md, design/system/strategic-agents.md, design/adrs, docs, flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/reasons_figures.go, flai/internal/context/reasons_figures_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/context.go, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/agent-context.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, design/system/conventions.md, ".claude/agents/orchestrator.md", flai/internal/check/check.go, flai/internal/check/vocabulary_test.go, design/conventions/tooling.md, flai/internal/check/check_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0424 Criterion 2: a real thread's pack is 31% of the general pack, under a quarter only for the test's small fixture thread: Recommended: trim the pack as S-0537's agent proposes, and tick criterion 2 on…

The operator, alex, answered TH-0753 at 2026-10-29T16:31:31Z.

What was asked:

```text
Criterion 2's test passes. A pack primed for one thread is under a quarter of the general pack, measured over this repository's conventions and design. But its fixture thread has one entry, on a three-line story. Packs for real threads here are about a third. I have left criterion 2 unticked until you choose.
```

The answer:

```text
**Recommended: trim the pack as S-0537's agent proposes, and tick criterion 2 on a test over a real thread.**

- **The item a thread is on loads as its title, status, goal, and acceptance criteria.** Its notes, tasks, and planning come from `item_get` when a reply needs them. ADR-0209 §3 asks for "the thread, its item", not the item whole.
- **The digest in a reason pack is capped at a sixteenth of the budget.** That stays within ADR-0194's eighth.
- **The test measures a copy of a real thread of this repository on its story,** not a one-entry fixture.
- **Record it in a new ADR refining ADR-0217 §3 and §4.** `design/adrs` is already in S-0537's touches.

**Why:** criterion 2 is there so that a run for one thread reads much less than the general pack (ADR-0209 §3). A test that passes on a three-line story while real threads measure 31% to 33% meets the criterion's words but not its purpose. The trim reaches about 42 to 44 KB, under a quarter, and changes no other pack.

**Or:** keep the pack as it is. The agent then ticks criterion 2 on the test as written and notes the real-thread figures in S-0537's Notes.

The review pack, at 45%, is not criterion 2's concern. It can go to S-0536's follow-up story under E-0030 if it matters there.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § Decision
```

Source: TH-0753.
