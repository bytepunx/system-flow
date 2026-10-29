---
id: I-0227
title: flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests
class: defect
status: open
count: 1
cost: 45m
first_reported: 2026-10-29T00:08:41Z
last_reported: 2026-10-29T00:08:41Z
updated: 2026-10-29T00:08:41Z
---

# I-0227 flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Description
flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Instances

### 2026-10-29T00:08:41Z
Story: S-0535.
S-0535 changed flai/internal/mcpserver. `flai test` on its paths failed go-test twice at 10m21s and 10m1.6s with a finding naming only the package. Run by hand with -timeout 560s, the package panicked "test timed out after 9m20s" while TestHeldWaitsSendProgress had run 1s, so no test hung. Without -race the package takes 130s on the story branch, against 76s for an exported copy of main's flai/. The gap is TestPrimeOnThisRepository (33s) and TestPrimeComesInPartsOnThisRepository (31s), which prime the real repository and skip outside a checkout; every other test took the same time on both. The host's load average was about 13 on 24 cores. The go-test tier in system-flow.yaml passes no -timeout, so go test's 10m default applies to the -race run, and the tier's finding names no test when the timeout panics. A sub-agent's earlier run of the same package passed in about 8.5 minutes. Remedy candidates: give the go-test tier a -timeout well above 10m, or report the tests still running when the timeout panics as the finding, or keep the two prime-on-this-repository tests out of -race.

## Remediation
