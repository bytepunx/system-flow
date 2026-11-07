---
id: D-0499
title: "S-0564: press Import for flai-scratch-import on the board, and turn on agent for the new project: I see `S-0001` in the board for this newly imported project"
kind: recommendation
sources: [TH-0853]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md, design/issues/summary.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0853]
---

# D-0499 S-0564: press Import for flai-scratch-import on the board, and turn on agent for the new project: I see `S-0001` in the board for this newly imported project

The operator, alex, answered TH-0853 at 2026-11-07T08:26:53Z.

What was asked:

```text
**The scratch repository is ready for the import run.** It is `/home/alex/git/flai-scratch-import`, below `/home/alex/git`, a folder named for import. It holds a `README.md` and a `Makefile` whose `test` target passes, in one commit, `aedf008`. `flai serve import list` offers it as `flai-scratch-import`.
```

The recommendation it overruled:

````text
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
````

The operator's answer:

```text
I see `S-0001` in the board for this newly imported project
```

Source: TH-0853.
