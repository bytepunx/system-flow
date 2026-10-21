---
id: TH-0526
title: Publishing flai 1.67.0 and flaiover 0.52.9 (S-0404, E-0022 done) was tagged here but the push failed
anchor:
  path: wip/archive/kanban/stories/S-0404-metrics-md-defines-the-queue-wait-per-story-the-batch-size-and-the-bisections-per-week-by-an-adr-and-flai-stats-and-the-dashboard-s-charts-report-them.md
  item: S-0404
status: open
participants: [orchestrator]
created: 2026-10-20T19:00:58Z
updated: 2026-10-20T19:00:58Z
---

# TH-0526 Publishing flai 1.67.0 and flaiover 0.52.9 (S-0404, E-0022 done) was tagged here but the push failed

On wip/archive/kanban/stories/S-0404-metrics-md-defines-the-queue-wait-per-story-the-batch-size-and-the-bisections-per-week-by-an-adr-and-flai-stats-and-the-dashboard-s-charts-report-them.md.

## Entries

### 2026-10-20T19:00:58Z orchestrator
**Recommended:** check that `flai/v1.67.0` and `flaiover/v0.52.9` reached the remote. If not, push them yourself, for example with `git push origin flai/v1.67.0 flaiover/v0.52.9`, or run `flai release --pending` again. flai says what is tagged is not tagged again. I will not try again until you answer here or the next story is accepted.

## What happened

I called `release_publish` after you accepted S-0404, which completed E-0022. Its answer:

> the release is applied and tagged here, but the push failed and nothing was forced: git push -q --atomic origin main flai/v1.67.0 flaiover/v0.52.9: exit status 1. The operator puts that right and publishes again; what is tagged is not tagged again

## What I can see

- Local `main` and `origin/main` both read `b3123cd8`, "chore: publish flai 1.67.0, flaiover 0.52.9". So main is on the remote now.
- Something else may have pushed main at the same moment, such as your acceptance of S-0404 13 seconds earlier. That would have made the atomic push fail and leave the two tags behind.
- I cannot list or push tags: flai guard keeps me to `release_publish`.

**Released by this run, if the tags land:** flai 1.66.0 to 1.67.0 and flaiover 0.52.8 to 0.52.9, both carrying S-0404.
