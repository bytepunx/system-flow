---
id: I-0332
title: A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-11-10T08:28:53Z
last_reported: 2026-11-10T08:28:53Z
updated: 2026-11-10T08:28:53Z
---

# I-0332 A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use

## Description
A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use

## Instances

### 2026-11-10T08:28:53Z
Story: S-0658.
S-0658, the publish formula's instance: its evaluate step names `flai release --evaluate --require-met --judged ""`. Run as written in the main checkout, the installed flai on PATH (1.86.6, older than 1.86.10 in the project's history) refused it with exit 1, `unknown flag: --judged`, before evaluating. Run as `scripts/flai.sh`, the tree's flai, as the project's tooling convention asks, it exited 0. A step's exit thus depends on which flai the agent's shell finds, and an agent that reads the formula literally blocks the run on a missing flag.

## Remediation
