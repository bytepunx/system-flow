---
id: I-0227
title: flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests
class: defect
status: closed
count: 4
cost: 24m
first_reported: 2026-10-29T00:08:41Z
last_reported: 2026-10-31T10:05:21Z
updated: 2026-10-31T10:05:21Z
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

### 2026-10-31T10:05:21Z
Story: S-0551.
S-0551, 2026-10-31T10:05Z: flai test on flai/internal/mcpserver/work.go and work_test.go ran the package with -race for 10m18s and failed on go test's 10-minute default timeout, naming only the package. Five stories were in progress on the host. Earlier runs of the same package by S-0551's task sub-agent had shown the new wait_for_work and promotion-candidate tests passing.

## Remediation

Story S-0540 remediates this issue, created from it at 2026-10-29T05:59:43Z.
Closed 2026-10-31T09:50:14Z: S-0548: the go-test tier passes -timeout 20m and flai check warns, as manifest.tests-timeout, of a tier that runs go test with none (T-2363); the close-out's go-test runs -short, so the two prime-on-this-repository tests are the integration tier's own (T-2366); and the package runs in parallel, 9m14s to 136 to 162 s under -race (T-2361). The remaining time is I-0234.
