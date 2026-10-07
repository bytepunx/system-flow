---
id: T-1260
type: task
nature: feature
title: The dashboard accepts a flai without a valid stamp only when FLAIOVER_ALLOW_UNSIGNED is set, and records which side is unsigned
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:08Z
updated: 2026-10-07T23:09:08Z
transitions: []
stream: S-0239
tags: [dashboard]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts]
---
# T-1260 The dashboard accepts a flai without a valid stamp only when FLAIOVER_ALLOW_UNSIGNED is set, and records which side is unsigned

## Work

The dashboard's half of criterion 2. It waits for nothing in this story: S-0237's stamp check in `hello` is on main by the time the story is pulled, and the flai side shares no path with it.

- `flaiover/src/lib/server/release.ts`: read `FLAIOVER_ALLOW_UNSIGNED` once at start; `1` allows, anything else or unset refuses, as S-0237 built it.
- `flaiover/src/lib/server/agent.ts`: with the allowance, a `hello` whose stamp is missing, unverifiable, or names another component is let through to the credential proof instead of closing with `4403`. The credential proof is never skipped. Each connection records which side is unsigned: `flai` when flai's stamp failed, `dashboard` when this build carries no valid stamp of its own, or both. A signed pair records nothing, with the allowance on or off.
- Tests, in `agent.test.ts` and `release.test.ts`: with the variable unset, an unsigned flai is still closed with `4403`; with it set, an unsigned flai is accepted and recorded as the unsigned side; an unsigned dashboard build is recorded as such; a signed pair with the variable set records nothing; a wrong credential is still refused with the variable set.

## Done when

- The dashboard accepts a flai without a valid stamp only when `FLAIOVER_ALLOW_UNSIGNED=1`, and records per connection which side is unsigned.
- A signed pair records nothing with the variable set.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. If S-0237 placed the stamp check outside `release.ts` and `agent.ts`, the story's agent widens the touches to where it is.
