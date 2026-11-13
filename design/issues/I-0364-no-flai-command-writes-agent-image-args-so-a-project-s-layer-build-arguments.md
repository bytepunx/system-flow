---
id: I-0364
title: No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand
class: efficiency
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T20:00:47Z
last_reported: 2026-11-13T20:00:47Z
updated: 2026-11-13T20:00:47Z
---

# I-0364 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand

## Description
No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand

## Instances

### 2026-11-13T20:00:47Z
Story: S-0760.
S-0760's criterion 1 asked for system-flow.yaml's agent.image.args to be written with the flai command that sets it. None does: flai manifest set lists no agent key, and flai agent set and clear keep agent.image as the file has it. S-0760 wrote the five arguments by hand and checked them with flai check and flai image status, as the operator confirmed on TH-1172. The same exception holds for agent.image.layer, the tests tiers (D-0449), and projects[].install (D-0593). Remedy: flai manifest set agent.image.args.<NAME>=<value>, with the manifest's own validation (no BASE, an ARG name, one line).

## Remediation
