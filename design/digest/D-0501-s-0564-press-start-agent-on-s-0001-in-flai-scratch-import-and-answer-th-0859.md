---
id: D-0501
title: "S-0564: press Start agent on S-0001 in flai-scratch-import, and answer TH-0859 for S-0584: Recommended: start S-0001's agent in the flai-scratch-import project"
kind: answer
sources: [TH-0860]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md, design/issues/summary.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0501 S-0564: press Start agent on S-0001 in flai-scratch-import, and answer TH-0859 for S-0584: Recommended: start S-0001's agent in the flai-scratch-import project

The operator, alex, answered TH-0860 at 2026-11-07T08:45:32Z.

What was asked:

```text
**The import worked. Its instance has no agent.** The board's Import committed the import as `1e7935b` and made S-0001, in ready, at 08:25:48Z. You turned on `agent` at 08:26:06Z. flai serve starts an agent only when a story enters ready, so S-0001 was still in ready with no agent at 08:32Z. I recorded this as I-0300. `agent_start` refuses an agent flai serve started, so I cannot start it myself.
```

The answer:

````text
**Recommended:** start S-0001's agent in the flai-scratch-import project. Press Start agent on S-0001 there, or run:

```sh
cd /home/alex/git/flai-scratch-import && flai serve agent start S-0001
```

**Why it is yours:** starting an agent is a host action, and it spends. flai guard refuses it to agents, and I have no permission for another project's board.

**The other two asks are covered:**

- TH-0859: you confirmed it at 08:32:14Z. S-0584's agent started again at 08:32:16Z and is working T-2545. Nothing more is needed there.
- Acceptance: I accept S-0584 under `accept_reviews` once it reaches review and passes. S-0001 lives in flai-scratch-import, so its acceptance is yours there.

**Alternative:** decline. agent-S-0564 then records the import run as stopped at its instance and leaves criterion 1 unticked.

Source: design/workflows/formulas/import.md

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/workflows/formulas/import.md
````

Source: TH-0860.
