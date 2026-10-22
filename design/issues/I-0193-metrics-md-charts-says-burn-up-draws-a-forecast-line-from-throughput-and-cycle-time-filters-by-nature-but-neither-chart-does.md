---
id: I-0193
title: "metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-22T21:09:37Z
last_reported: 2026-10-22T21:09:37Z
updated: 2026-10-22T21:09:37Z
---

# I-0193 metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does

## Description
metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does

## Instances

### 2026-10-22T21:09:37Z
Story: S-0485.
Writing S-0485's chart explainers from metrics.md: its § Charts row for Burn-up lists a forecast line from throughput, but burnUp in flaiover/src/lib/viz/charts.ts draws only scope and done (docs/users/flaiover.md agrees with the chart). Its Cycle Time row says "Filter by nature", but the chart is filtered by epic and coloured by nature. Also, metrics.md is silent on the dashed all-models line of Tokens / Min and Tokens / $ and on Strategic Use's mean waiting, which the dashboard works out from wait_threads_seconds plus wait_review_seconds. The explainers follow what the charts draw; metrics.md changes only with an ADR, so S-0485 left it as it is.

## Remediation
