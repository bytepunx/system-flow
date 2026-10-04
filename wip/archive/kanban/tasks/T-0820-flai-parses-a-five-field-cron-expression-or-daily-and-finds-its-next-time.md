---
id: T-0820
type: task
nature: improvement
title: flai parses a five-field cron expression or daily and finds its next time
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:41:53Z
updated: 2026-10-04T20:49:31Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:47Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:42:49Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T20:49:31Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [flai/internal/cron]
usage:
  source: log
  seconds: 402
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 6201
      cache_read: 995922
      cache_write: 25170
      cost: 0.4742
---
# T-0820 flai parses a five-field cron expression or daily and finds its next time

## Work

A new package `flai/internal/cron`, standard library only (ADR-0084: no cron dependency). `Parse(spec string) (Schedule, error)` takes a five-field expression (minute, hour, day of month, month, day of week) with `*`, numbers, lists, ranges, and steps, or `daily`, which is `0 0 * * *`. `Schedule.Next(after time.Time) time.Time` is the first matching minute strictly after `after`, in UTC. Day of month and day of week combine as cron does: when both are restricted, either matches. Errors name the field and the value and say what is allowed. It waits for nothing: it is the first layer.

## Done when

- [ ] `Parse` accepts `daily` and valid expressions and refuses bad ones with an actionable error
- [ ] `Next` is right across hours, days, months, a year end, and the day-of-month/day-of-week rule
- [ ] Behavior tests in `flai/internal/cron` pass under `go test -race`

## Notes
