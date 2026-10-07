---
id: T-1180
type: task
nature: improvement
title: A shared formatter shows a recorded UTC time in the browser's local zone, and the tests run in a fixed non-UTC zone
status: done
parent: S-0329
owner: alex
created: 2026-10-07T19:49:10Z
updated: 2026-10-07T20:26:12Z
transitions:
  - to: ready
    at: 2026-10-07T20:25:03Z
    by: agent-S-0329
  - to: in-progress
    at: 2026-10-07T20:25:04Z
    by: agent-S-0329
  - to: done
    at: 2026-10-07T20:26:12Z
    by: agent-S-0329
stream: S-0329
tags: [dashboard]
touches: [flaiover/src/lib/localtime.ts, flaiover/src/lib/localtime.test.ts, flaiover/vite.config.ts]
usage:
  source: log
  seconds: 68
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 29
      cache_read: 605804
      cache_write: 22547
      cost: 0.2808
---
# T-1180 A shared formatter shows a recorded UTC time in the browser's local zone, and the tests run in a fixed non-UTC zone

## Work

flaiover has no shared formatter for an absolute time: each component formats flai's UTC strings itself, most often with `s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC')` or `.slice(0, 10)`. The dashboard is a single-page app (`ssr = false` in `flaiover/src/routes/+layout.ts`), so every time is formatted in the browser, where its zone is known.

- Add `flaiover/src/lib/localtime.ts` with two functions that take a time as flai records it (`YYYY-MM-DDTHH:MM:SSZ`) and an optional IANA zone, the browser's when left out:
  - `localTime(at)`: the time to the minute in that zone with the zone's short name, as `2026-10-07 15:35 EDT`, keeping the shape the components show today with the zone in place of `UTC`.
  - `localDate(at)`: the calendar date in that zone, `YYYY-MM-DD`.
  - Both return what they were given when it does not parse, so a missing or odd value shows as it did.
- Add `flaiover/src/lib/localtime.test.ts`: a time across midnight in two zones, a zone with a half-hour offset, a summer and a winter time, and a value that does not parse.
- Pin the zone the tests run in to one that is not UTC, `America/New_York`, in `flaiover/vite.config.ts` (`env: { TZ: ... }` on both test projects, or a setup file), so that a component test proves the conversion rather than passing in UTC. Check that jsdom's `Date` and `Intl` follow it.

Nothing waits for this task; the three view tasks wait for it because they call it.

## Done when

- `flaiover/src/lib/localtime.ts` exports `localTime` and `localDate`, with tests that pass.
- The vitest projects run with a non-UTC zone, and `flai test flaiover/src/lib/localtime.ts flaiover/vite.config.ts` passes.

## Notes

- Pinning the zone may break an existing test that relies on the host running in UTC through a local `Date` method; fix it here, since the zone change is this task's.
- Recording stays UTC: nothing here changes what flai writes or what the API returns.
