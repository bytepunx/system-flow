---
id: I-0010
title: Host Node is 25, outside the versions the project and its dependencies target
class: efficiency
status: open
count: 4
cost: 3m
first_reported: 2026-09-17T06:01:24Z
last_reported: 2026-09-18T20:48:07Z
updated: 2026-09-18T20:48:07Z
---

# I-0010 Host Node is 25, outside the versions the project and its dependencies target

## Description
Host Node is 25, outside the versions the project and its dependencies target

## Instances

### 2026-09-17T06:01:24Z
S-0032: @prometheus-io/client declares node ^22 || ^24 || >=26 and flaiover/.npmrc has engine-strict, so pnpm add refused on the host's Node 25 until the engine check was relaxed for that install. CI and the image run Node 24. Corepack is also absent on 25, which is why pnpm is installed by npm into .flai-cache.

### 2026-09-18T18:18:10Z
S-0049: make flaiover-install refused in the story worktree on Node 25; rerun with npm_config_engine_strict=false for that install only

### 2026-09-18T19:51:45Z
S-0050: same refusal installing flaiover dependencies in the story worktree; engine-strict relaxed for the install

### 2026-09-18T20:48:07Z
S-0048: same refusal installing flaiover dependencies in the story worktree

## Remediation
