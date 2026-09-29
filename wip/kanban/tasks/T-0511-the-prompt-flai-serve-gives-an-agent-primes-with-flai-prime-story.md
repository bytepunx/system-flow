---
id: T-0511
type: task
nature: feature
title: The prompt flai serve gives an agent primes with flai prime --story
status: done
parent: S-0138
owner: alex
created: 2026-09-29T01:08:51Z
updated: 2026-09-29T01:12:30Z
transitions:
  - to: ready
    at: 2026-09-29T01:08:55Z
    by: agent-S-0138
  - to: in-progress
    at: 2026-09-29T01:11:29Z
    by: agent-S-0138
  - to: done
    at: 2026-09-29T01:12:30Z
    by: agent-S-0138
stream: S-0138
tags: [cli]
touches: [flai/internal/harness]
---
# T-0511 The prompt flai serve gives an agent primes with flai prime --story

## Work

Change `harness.Prompt` so the story prompt says to prime with `flai prime --story <id>` and to read what its catalog lists when needed, instead of `flai prime --cat`. Pin the new text in its test.

## Done when

- The prompt names `flai prime --story <id>`; `harness` tests pass.

## Notes
