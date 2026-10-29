---
id: I-0227
title: flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests
class: defect
status: open
count: 3
cost: 24m
first_reported: 2026-10-29T00:08:41Z
last_reported: 2026-10-29T05:14:13Z
updated: 2026-10-29T05:14:13Z
---

# I-0227 flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Description
flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Instances

### 2026-10-29T00:08:41Z
Story: S-0535.
S-0535 changed flai/internal/mcpserver. `flai test` on its paths failed go-test twice at 10m21s and 10m1.6s with a finding naming only the package. Run by hand with -timeout 560s, the package panicked "test timed out after 9m20s" while TestHeldWaitsSendProgress had run 1s, so no test hung. Without -race the package takes 130s on the story branch, against 76s for an exported copy of main's flai/. The gap is TestPrimeOnThisRepository (33s) and TestPrimeComesInPartsOnThisRepository (31s), which prime the real repository and skip outside a checkout; every other test took the same time on both. The host's load average was about 13 on 24 cores. The go-test tier in system-flow.yaml passes no -timeout, so go test's 10m default applies to the -race run, and the tier's finding names no test when the timeout panics. A sub-agent's earlier run of the same package passed in about 8.5 minutes. Remedy candidates: give the go-test tier a -timeout well above 10m, or report the tests still running when the timeout panics as the finding, or keep the two prime-on-this-repository tests out of -race.

### 2026-10-29T04:20:50Z
Story: S-0539.
S-0539 changed only permission.go and its test in flai/internal/mcpserver. flai test's go-test tier on the package took 9m58.6s and passed with 1.4 s to spare under Go's 10-minute default. In the close-out it took 9m44.7s, which pushed the close-out past Claude Code's 600-second Bash limit, so the run had to go on in the background and be waited for.

### 2026-10-29T05:14:13Z
Story: S-0417.
S-0417: flai test on mcpserver, serve, and workitem together (load average about 10) failed the go-test tier after 10m13s, naming only flai/internal/mcpserver. Run alone, the same package passed in 9m20s.

## Remediation
