---
id: I-0166
title: flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
class: efficiency
status: open
count: 1
cost: 2s
first_reported: 2026-10-19T14:20:51Z
last_reported: 2026-10-19T14:20:51Z
updated: 2026-10-21T11:59:35Z
---

# I-0166 flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s

## Description
flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s

## Instances

### 2026-10-19T14:20:51Z
Story: S-0449.
S-0449's T-1907 resolves each json prop's type with the sub-project's typescript in one program per run, to give it a skeleton (TH-0515). Measured on flaiover's 54 components: props.mjs went from 0.38-0.48s to 2.0-2.4s, and `bin/flai components list --project flaiover --json` takes 2.4s. Most of it is building the program (about 1.3s, since the components' imports pull in 567 files) and the checker (0.3s). Every components.list pays it: the Components page asks on arrival and every minute, and flai check's components.examples rule lists the props too. Caching the skeletons by the files the program read, or keeping a watch program in flai serve, would bring it back down.

## Remediation

Story S-0478 remediates this issue, created from it at 2026-10-21T11:59:35Z.
