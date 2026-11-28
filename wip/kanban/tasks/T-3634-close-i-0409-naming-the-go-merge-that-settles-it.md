---
id: T-3634
type: task
nature: remediation
title: Close I-0409 naming the Go merge that settles it
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:34Z
updated: 2026-11-27T23:58:34Z
transitions: []
stream: S-0857
tags: [issues]
touches: [design/issues/I-0409-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md, design/issues/summary.md]
after: [T-3630, T-3633]
---
# T-3634 Close I-0409 naming the Go merge that settles it

## Work

Close I-0409 with `flai issue close I-0409 --reason`, saying that every rebase flai runs now settles a stop on a `.go` file where each side only adds whole top-level declarations at one place, naming the ADR of T-3623 and the tests of T-3628 and T-3630 that reproduce it. It waits for T-3630 and T-3633, since the reason names what they built and documented.

## Done when

- I-0409 is closed with that reason and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
