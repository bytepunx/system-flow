---
id: I-0276
title: A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project
class: efficiency
status: open
count: 1
cost: 20m
first_reported: 2026-11-07T05:43:57Z
last_reported: 2026-11-07T05:43:57Z
updated: 2026-11-10T08:34:40Z
---

# I-0276 A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Description
A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Instances

### 2026-11-07T05:43:57Z
Story: S-0567.
S-0567, T-2500: a docker run from agent-S-0567's shell with -v /home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flow:ro mounted an empty folder, inside and outside Claude Code's sandbox, on Docker Desktop (docker-desktop context), while flai serve's planner container mounted the same path with its files. So the project's own guard hook could not be replayed by hand in the container; only the image's flai guard, which needs no checkout, could.

## Remediation

Story S-0667 remediates this issue, created from it at 2026-11-10T08:34:40Z.
