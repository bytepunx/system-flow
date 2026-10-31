---
id: I-0225
title: flai/internal/mcpserver's tests take 8 to 10 minutes under -race, at Go's 10-minute default timeout, so the integration tier fails a merge queue batch on time alone
class: defect
status: closed
count: 1
cost: 30m
first_reported: 2026-10-28T23:28:42Z
last_reported: 2026-10-28T23:28:42Z
updated: 2026-10-31T09:50:10Z
---

# I-0225 flai/internal/mcpserver's tests take 8 to 10 minutes under -race, at Go's 10-minute default timeout, so the integration tier fails a merge queue batch on time alone

## Description
flai/internal/mcpserver's tests take 8 to 10 minutes under -race, at Go's 10-minute default timeout, so the integration tier fails a merge queue batch on time alone

## Instances

### 2026-10-28T23:28:42Z
Story: S-0507.
Merge queue batch 20261028-231632 sent S-0507 back on TH-0723: the integration tier's go test failed flai/internal/mcpserver at 600.17 s, "panic: test timed out after 10m0s", with TestWaitForWorkResumesYourOwnStory running for 0 s, so no test hung. The package took 470 to 591 s in every batch of 2026-10-28 (509, 501, 502, 470, 485, 470, 586, 591), and S-0507 changes nothing in it. S-0507 gives scripts/integration.sh -timeout 20m, which stops the false failure; the suite's time is still the cause, and is worth a story of its own to cut.

## Remediation
Closed 2026-10-31T09:50:10Z: S-0548: the mcpserver tests run in parallel (T-2361), 9m14s to 136 to 162 s under -race on this host at load 14 to 57, and the integration tier and the go-test tier pass -timeout 20m (T-2363), so the package no longer meets go's 10-minute default. What remains of its time is I-0234.
