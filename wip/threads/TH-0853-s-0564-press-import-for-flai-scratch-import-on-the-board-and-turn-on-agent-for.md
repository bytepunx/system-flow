---
id: TH-0853
title: "S-0564: press Import for flai-scratch-import on the board, and turn on agent for the new project"
anchor:
  path: wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md
  item: S-0564
status: resolved
participants: [agent-S-0564, orchestrator, alex]
created: 2026-11-07T08:03:52Z
updated: 2026-11-07T08:33:52Z
---

# TH-0853 S-0564: press Import for flai-scratch-import on the board, and turn on agent for the new project

On wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md.

## Entries

### 2026-11-07T08:03:52Z agent-S-0564
**The scratch repository is ready for the import run.** It is `/home/alex/git/flai-scratch-import`, below `/home/alex/git`, a folder named for import. It holds a `README.md` and a `Makefile` whose `test` target passes, in one commit, `aedf008`. `flai serve import list` offers it as `flai-scratch-import`.

**What I need from you, as the plan's question 2 recommended:**

1. On the board, press Import on `flai-scratch-import`. That runs `flai import --formula`, which starts the import formula's instance in the new project.
2. Turn on the `agent` host action for the new project, so that its instance gets an agent: `flai serve enable agent`, run in `/home/alex/git/flai-scratch-import`, or the dashboard's settings.
3. Reply here once both are done, with the instance's story ID if the board shows it.

**Decision taken without an answer, question 1:** you finalized S-0564 without choosing a focus. So the cadence ran with none, for all three, as the Notes and the recommendation say. Its instance is S-0584, in ready, committed on main at `aa01853e5`.

**What I do with your answer:** I follow the import instance to review and acceptance (T-2523), and S-0584 to acceptance (T-2526). If you would rather not import, say so: I will record the import run as not done and leave criterion 1 unticked.

### 2026-11-07T08:04:53Z orchestrator (recommendation)
**Recommended:** do both steps agent-S-0564 asks for, then reply here with the instance's story ID.

1. On the board, press Import on `flai-scratch-import`. The board's Import runs `flai import --formula`, which starts the import formula in the new project. That is the path S-0424 changed, and the one this run must prove.
2. Turn on the `agent` host action for the new project:

```sh
cd /home/alex/git/flai-scratch-import && flai serve enable agent
```

Or use the dashboard's settings.

**Why it is yours:** both are host changes. An import and a host action are the operator's, and flai guard refuses them to agents. This is the step you were asked about on TH-0845, question 2, and it has not been answered yet.

**Alternative:** decline the import. S-0564 then records the import run as not done and leaves criterion 1 unticked. Only the import waits on this. The cadence instance, S-0584, is already in progress.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/workflows/formulas/import.md

### 2026-11-07T08:26:53Z alex
I see `S-0001` in the board for this newly imported project

### 2026-11-07T08:32:02Z agent-S-0564
Thanks: the import is committed as `1e7935b`, and its instance is S-0001 (T-2514 is closed). S-0001 got no agent, because it entered ready before `agent` was on (I-0300). Starting it is asked on TH-0860.

### 2026-11-07T08:33:39Z orchestrator
Nothing is left to answer on this thread. alex pressed Import and turned on `agent`, and the board shows the instance as S-0001 (D-0499). The one remaining ask, starting S-0001's agent, is on TH-0860, where I recommended it to alex as theirs alone. agent-S-0564: resolve this thread, and wait on TH-0860 instead.

Source: design/digest/D-0499-s-0564-press-import-for-flai-scratch-import-on-the-board-and-turn-on-agent-for.md

Source: design/digest/D-0499-s-0564-press-import-for-flai-scratch-import-on-the-board-and-turn-on-agent-for.md

### 2026-11-07T08:33:52Z agent-S-0564
Resolved: The operator imported flai-scratch-import from the board (1e7935b, instance S-0001) and turned on agent for it. Starting S-0001's agent is asked on TH-0860.
