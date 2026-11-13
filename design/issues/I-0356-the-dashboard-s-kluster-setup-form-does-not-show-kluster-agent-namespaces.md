---
id: I-0356
title: The dashboard's kluster setup form does not show kluster.agent_namespaces
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T08:42:09Z
last_reported: 2026-11-13T08:42:09Z
updated: 2026-11-13T08:42:09Z
---

# I-0356 The dashboard's kluster setup form does not show kluster.agent_namespaces

## Description
The dashboard's kluster setup form does not show kluster.agent_namespaces

## Instances

### 2026-11-13T08:42:09Z
Story: S-0660.
S-0660 added kluster.agent_namespaces to the kluster: block and to flai manifest set's settings, as ADR-0292 §4 decides, but flaiover/src/lib/components/DeploySettings.svelte lists its keys in PLAIN_KEYS, which does not name it, so the setup form on the Klusters and Deployments pages does not show it. docs/operators/settings.md and kluster.md say it is changed with flai manifest set or by hand for now.

## Remediation
