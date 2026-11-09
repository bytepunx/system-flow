---
id: I-0307
title: flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing
class: efficiency
status: closed
count: 1
cost: 20m
first_reported: 2026-11-07T22:43:39Z
last_reported: 2026-11-07T22:43:39Z
updated: 2026-11-09T09:06:42Z
---

# I-0307 flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing

## Description
flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing

## Instances

### 2026-11-07T22:43:39Z
Story: S-0596.
S-0596's three exploration sub-agents (T-2607, T-2608, T-2609) were asked to time methods with `flai hostapi --timing`, which docs/users/flai.md § Request timing names as the way to time one method. flai guard refused it to each, since hostapi is not on its read list (flai/internal/guard/guard.go), and refused `git clone` for a scratch copy too. Each wrote and built an in-process Go program with `go build -overlay` instead, or a tar copy of the main checkout, which took about 20 minutes each and leaves their figures less comparable with S-0152's. Read-only hostapi methods could be allowed to sub-agents, with the write methods still refused.

## Remediation

Story S-0612 remediates this issue, created from it at 2026-11-07T23:38:52Z.
Closed 2026-11-09T09:06:42Z: S-0612: flai guard lets every agent run flai hostapi alone and with a read method on guard.HostAPIReads, --timing included, so a sub-agent can time a method as docs/users/flai.md § Request timing says. It refuses every other method by name, since the host API's write methods are the story's agent's (ADR-0060). flai/internal/hostapi/guardreads_test.go keeps the list in step with the method table.
