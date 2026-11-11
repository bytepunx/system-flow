---
id: T-3078
type: task
nature: remediation
title: "The tooling convention says to brace a shell variable before a colon, since zsh reads $name:r as a modifier"
status: done
parent: S-0667
owner: alex
created: 2026-11-10T11:02:17Z
updated: 2026-11-10T11:03:14Z
transitions:
  - to: ready
    at: 2026-11-10T11:03:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T11:03:02Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:03:14Z
    by: agent-S-0667
stream: S-0667
tags: [docker, agent-container]
touches: [design/conventions/tooling.md]
after: [T-2926]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 829
      cache_read: 426071
      cache_write: 12598
      cost: 0.2026
---
# T-3078 The tooling convention says to brace a shell variable before a colon, since zsh reads $name:r as a modifier

## Work

Waits for T-2926: it writes the fix T-2926 named in S-0667's Notes under `### Cause`.

In `design/conventions/tooling.md`, below the marker under `## Project additions`, add a rule beside the baseline's zsh rule. In zsh, `$name:` followed by `r`, `h`, `t`, `e`, or another modifier letter is a history modifier, even inside double quotes. So a hand-written command writes `${name}` before a colon, as in `-v "${ROOT}:${ROOT}:ro"`, and names I-0276. Bump the file's `updated`. The baseline above the marker stays as it is: S-0667's narrative proposes the baseline edit in `## Open questions`.

## Done when

- `design/conventions/tooling.md`'s project additions say to brace a name before a colon in a shell command, why, and the example.
- `flai test design/conventions/tooling.md` passes.

## Notes
