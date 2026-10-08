---
id: I-0113
title: A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
class: efficiency
status: open
count: 2
cost: 7m
first_reported: 2026-10-07T09:13:42Z
last_reported: 2026-10-08T04:55:56Z
updated: 2026-10-08T04:55:56Z
---

# I-0113 A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Description
A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Instances

### 2026-10-07T09:13:42Z
Story: S-0215.
S-0215's close-out stopped at the integration tier with only the tail of the go test output: a list of ok packages and FAIL. flai verify --last showed the same tail, so the failing package and test were not named. Rerunning scripts/integration.sh passed, so it was a flake, likely I-0102 or I-0106, but which one cannot be told.

### 2026-10-08T04:55:56Z
Story: S-0320.
S-0320's second close-out failed the integration tier. Its findings hold only `ok` lines from `internal/protected` through `tests/integration`, then `FAIL`, so the failing package (one before `protected`, since every listed package passed) and its test are not named. The same code had passed integration in the first close-out minutes earlier; only issue files and main's S-0316 commits changed in between.

## Remediation

Story S-0313 remediates this issue, created from it at 2026-10-07T14:26:02Z.
