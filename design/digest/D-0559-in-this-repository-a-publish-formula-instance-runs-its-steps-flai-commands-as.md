---
id: D-0559
title: In this repository a publish formula instance runs its steps' flai commands as scripts/flai.sh, since the installed flai on PATH can lack the flags the steps use, as 1.86.6 lacked --judged (I-0332)
kind: decision
sources: [S-0658, I-0332]
topics: [formulas, release, flai]
paths: [design/workflows/formulas/publish.md, scripts/flai.sh, design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md]
date: 2026-11-10
status: active
absorbed_by: ""
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

Until S-0681 lands, a story's agent working a formula's instance here runs each step's `flai` as `scripts/flai.sh`. An exit that is the wrong binary's, such as `unknown flag`, is not the step's exit: run the step again with the tree's flai rather than block the task. D-0308 is the same split for `flai test` through MCP.

Sources: S-0658 (`wip/archive/agents/S-0658.md`, Decisions); I-0332.
