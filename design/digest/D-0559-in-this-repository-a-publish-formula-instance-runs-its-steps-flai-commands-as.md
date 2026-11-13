---
id: D-0559
title: In this repository a publish formula instance runs its steps' flai commands as scripts/flai.sh, since the installed flai on PATH can lack the flags the steps use, as 1.86.6 lacked --judged (I-0332)
kind: decision
sources: [S-0658, I-0332, S-0704, S-0762, S-0766, S-0770]
topics: [formulas, release, flai]
paths: [design/workflows/formulas/publish.md, scripts/flai.sh, design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md]
date: 2026-11-10
status: absorbed
absorbed_by: "design/system/strategic-agents.md#working-an-instance-s-0422"
asked_again: []
cited: []
---

# D-0559 In this repository a publish formula instance runs its steps' flai commands as scripts/flai.sh, since the installed flai on PATH can lack the flags the steps use, as 1.86.6 lacked --judged (I-0332)

S-0658, an instance of the publish formula, ran its evaluate step, `flai release --evaluate --require-met --judged ""`, as written in the main checkout. The installed flai on `PATH`, 1.86.6, older than 1.86.10 in the project's history, refused it with exit 1, `unknown flag: --judged`, before evaluating.

From S-0658's Decisions (`wip/archive/agents/S-0658.md`):

```text
The `flai` commands of the formula's steps run as `scripts/flai.sh`, the tree's flai, in the main checkout, as `design/conventions/session-start.md` and `tooling.md` project additions ask. The installed flai on `PATH`, 1.86.6, refused the evaluate step with exit 1, `unknown flag: --judged`, before evaluating; that exit was the wrong binary's, not the step's, so the step was run again with the tree's flai, which exited 0.
```

Every step then exited 0, and the publish step tagged and pushed `flai/v1.86.11`. I-0332 records the cost: a step's exit depends on which flai the shell finds, and an agent that reads the formula literally blocks the run on a missing flag. S-0681 remediates it.

S-0704, the next publish instance, did the same from the start. From its Decisions (`wip/archive/agents/S-0704.md`):

```text
Each formula step's bare flai command runs as the main checkout's scripts/flai.sh, which builds bin/flai from main's tree, as session-start.md's project additions say: the installed flai is 1.86.12, behind 1.86.13 in the project's history (S-0681 records the bare flai on PATH problem).
```

Until S-0681 lands, a story's agent working a formula's instance here runs each step's `flai` as `scripts/flai.sh`. An exit that is the wrong binary's, such as `unknown flag`, is not the step's exit: run the step again with the tree's flai rather than block the task. D-0308 is the same split for `flai test` through MCP.

S-0762, a later publish instance, ran the same way with flai 1.86.17 installed. From its Decisions (`wip/archive/agents/S-0762.md`):

```text
Layers are the formula's: T-3265, then T-3266, then T-3267, each waiting on the one before; run in the main checkout with scripts/flai.sh, as S-0646's instance did, since the installed flai (1.86.17) is behind the tree.
```

S-0766, the next publish instance, did the same, with the installed flai further behind. From its Decisions (`wip/archive/agents/S-0766.md`):

```text
Formula steps run with the tree's flai, `scripts/flai.sh`, in the main checkout, as S-0729's publish did: the installed flai 1.86.17 is behind 1.88.0 in the project's history.
```

S-0770, the publish instance after it, now names the manifest's `flai.command` as where `scripts/flai.sh` comes from. From its Decisions (`wip/archive/agents/S-0770.md`):

```text
Each step's `flai` runs as the manifest's `flai.command`, `scripts/flai.sh`, in the main checkout, as the formula's steps say; the tasks' layers are the formula's own, fetch, then evaluate, then publish, each after the one before.
```

Sources: S-0658 (`wip/archive/agents/S-0658.md`, Decisions); I-0332; S-0704 (`wip/archive/agents/S-0704.md`, Decisions); S-0762 (`wip/archive/agents/S-0762.md`, Decisions); S-0766 (`wip/archive/agents/S-0766.md`, Decisions); S-0770 (`wip/archive/agents/S-0770.md`, Decisions).
