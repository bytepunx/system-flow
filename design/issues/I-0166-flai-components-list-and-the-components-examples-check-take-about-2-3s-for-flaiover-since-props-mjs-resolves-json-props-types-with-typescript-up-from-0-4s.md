---
id: I-0166
title: flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
class: efficiency
status: closed
count: 1
cost: 2s
first_reported: 2026-10-19T14:20:51Z
last_reported: 2026-10-19T14:20:51Z
updated: 2026-10-23T07:02:23Z
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
Closed 2026-10-23T07:02:23Z: Fixed by S-0478 (T-2068): props.mjs keeps the skeletons it resolves in a cache under .flai-cache/component-props/<sub-project path>/, keyed by the script, the sub-project, its typescript, and the components' typings, and stamped by every path typescript looked at. A run with nothing changed takes them from the cache without loading typescript: flaiover's 58 components take 0.43s, against 2.7s cold. TestListTakesSkeletonsFromTheCacheUntilAFileTypescriptReadChanges reproduces it.
