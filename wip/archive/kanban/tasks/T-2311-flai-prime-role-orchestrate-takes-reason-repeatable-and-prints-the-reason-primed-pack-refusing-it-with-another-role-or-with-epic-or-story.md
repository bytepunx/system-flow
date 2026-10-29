---
id: T-2311
type: task
nature: improvement
title: flai prime --role orchestrate takes --reason, repeatable, and prints the reason-primed pack, refusing it with another role or with --epic or --story
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:07:30Z
updated: 2026-10-29T16:12:03Z
transitions:
  - to: ready
    at: 2026-10-29T15:52:16Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:52:17Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T16:12:03Z
    by: agent-S-0537
stream: S-0537
tags: [flai, cli, prime]
touches: [flai/cmd/prime.go, flai/cmd/prime_test.go]
after: [T-2308, T-2309]
usage:
  source: log
  seconds: 1186
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 11598
      cache_read: 2317534
      cache_write: 74884
      cost: 1.1212
---
# T-2311 flai prime --role orchestrate takes --reason, repeatable, and prints the reason-primed pack, refusing it with another role or with --epic or --story

## Work

Criterion 1, the CLI half. In `flai/cmd/prime.go`, add `--reason`, a repeatable string flag, to `flai prime`. With `--role orchestrate` it calls the reason-primed pack of T-2308 in place of `ForStrategic`, in text and `--json`, with `--budget` and `--part` as today. With another role, or beside `--epic` or `--story`, it is refused, naming why. An unknown reason is refused with T-2308's message naming the known kinds, and exits non-zero. Without `--reason` the general pack prints as today, byte for byte.

Waits for T-2308 and T-2309, which build the pack for every reason the tests cover.

## Done when

- `prime_test.go` primes with each reason kind once, asserting the item each loads and that no design brief or catalog does, and covers the refusals: an unknown reason, `--reason` with `--role plan`, and `--reason` with `--story`.
- A test without `--reason` keeps the general pack unchanged.
- `flai test` passes on the changed paths.

## Notes
