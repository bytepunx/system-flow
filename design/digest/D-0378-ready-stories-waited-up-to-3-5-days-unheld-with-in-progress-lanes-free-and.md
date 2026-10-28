---
id: D-0378
title: Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit
kind: cause
sources: [I-0153]
topics: [flai, cli, go, code, all]
paths: [design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md, design/adrs/README.md, design/system/metrics.md, flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, design/issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md, design/issues/summary.md]
date: 2026-10-23
status: absorbed
absorbed_by: "design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md#context"
asked_again: []
cited: []
---

# D-0378 Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit

I-0153 was closed on 2026-10-23.

The cause:

```text
Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit
```

The fix:

```text
S-0469: the 75 h were not unheld ready wait. The operator sent S-0212 to S-0216 back to backlog on 2026-10-04 and made them ready again on 2026-10-07; queue time, started - committed, counted the 269,000 s in backlog, and their time in ready was held almost whole. ADR-0190 makes queue time the time in ready before started, which flai stats now derives, so queue time minus held_seconds is the unheld ready wait.
```

Source: I-0153.
