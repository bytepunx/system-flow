---
id: I-0113
title: A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
class: efficiency
status: open
count: 1
cost: 4m
first_reported: 2026-10-07T09:13:42Z
last_reported: 2026-10-07T09:13:42Z
updated: 2026-10-07T14:26:02Z
---

# I-0113 A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Description
A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Instances

### 2026-10-07T09:13:42Z
Story: S-0215.
S-0215's close-out stopped at the integration tier with only the tail of the go test output: a list of ok packages and FAIL. flai verify --last showed the same tail, so the failing package and test were not named. Rerunning scripts/integration.sh passed, so it was a flake, likely I-0102 or I-0106, but which one cannot be told.

## Remediation

Story S-0313 remediates this issue, created from it at 2026-10-07T14:26:02Z.
