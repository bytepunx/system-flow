---
id: I-0010
title: Host Node is 25, outside the versions the project and its dependencies target
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-09-17T06:01:24Z
last_reported: 2026-09-17T06:01:24Z
updated: 2026-09-17T06:01:24Z
---

# I-0010 Host Node is 25, outside the versions the project and its dependencies target

## Description
Host Node is 25, outside the versions the project and its dependencies target

## Instances

### 2026-09-17T06:01:24Z
S-0032: @prometheus-io/client declares node ^22 || ^24 || >=26 and flaiover/.npmrc has engine-strict, so pnpm add refused on the host's Node 25 until the engine check was relaxed for that install. CI and the image run Node 24. Corepack is also absent on 25, which is why pnpm is installed by npm into .flai-cache.

## Remediation
