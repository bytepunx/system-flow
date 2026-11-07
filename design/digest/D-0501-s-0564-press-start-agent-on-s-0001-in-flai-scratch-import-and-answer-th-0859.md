---
id: D-0501
title: "S-0564: press Start agent on S-0001 in flai-scratch-import, and answer TH-0859 for S-0584: Recommended: start S-0001's agent in the flai-scratch-import project"
kind: answer
sources: [TH-0860, S-0589]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md, design/issues/summary.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0501 S-0564: press Start agent on S-0001 in flai-scratch-import, and answer TH-0859 for S-0584: Recommended: start S-0001's agent in the flai-scratch-import project

Starting a story's agent in another project flai serve serves is the operator's: a story's agent cannot do it, since `agent_start` refuses an agent flai serve started, and neither can the orchestrator, which has no permission for another project's board. Starting an agent is a host action, and it spends. So a story whose work needs an agent in another project, such as a formula's instance in a scratch repository it imported, asks the operator to press Start agent there, or to run `flai serve agent start S-nnnn` in that project.

The thread gave a wrong reason for why S-0001 had no agent: that flai serve starts an agent only when a story enters ready, and S-0001 was already in ready when the operator turned on `agent` (I-0300). S-0589 found that reason does not hold. flai serve starts a ready story that has had no run since it entered ready, whenever it entered ([ADR-0041](../adrs/0041-a-story-in-ready-gets-its-agent-whether-it-entered-ready-before-or-after-flai.md)). S-0001 waited because the imported project's default agent named no harness and no command was set on the host. It started once a default agent with a harness was set. See [D-0512](D-0512-a-ready-story-that-got-no-agent-after-the-agent-action-was-turned-on-was.md).

The operator, alex, answered TH-0860 at 2026-11-07T08:45:32Z, on S-0564, by taking the orchestrator's recommendation.

What was asked:

```text
**The import worked. Its instance has no agent.** The board's Import committed the import as `1e7935b` and made S-0001, in ready, at 08:25:48Z. You turned on `agent` at 08:26:06Z. flai serve starts an agent only when a story enters ready, so S-0001 was still in ready with no agent at 08:32Z. I recorded this as I-0300. `agent_start` refuses an agent flai serve started, so I cannot start it myself.
```

The answer:

```text
**Recommended:** start S-0001's agent in the flai-scratch-import project. Press Start agent on S-0001 there, or run: cd /home/alex/git/flai-scratch-import && flai serve agent start S-0001

**Why it is yours:** starting an agent is a host action, and it spends. flai guard refuses it to agents, and I have no permission for another project's board.

Acceptance: I accept S-0584 under `accept_reviews` once it reaches review and passes. S-0001 lives in flai-scratch-import, so its acceptance is yours there.
```

Source: TH-0860, on S-0564; the issue I-0300; the formula `design/workflows/formulas/import.md`; the Decisions of wip/archive/agents/S-0589.md, which corrected the reason.
