---
id: T-3071
type: task
nature: improvement
title: Close I-0275 with what fixed it
status: backlog
parent: S-0666
owner: alex
created: 2026-11-10T10:36:14Z
updated: 2026-11-10T10:36:14Z
transitions: []
stream: S-0666
tags: [issues]
touches: [design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md, design/issues/summary.md]
after: [T-3066, T-3069]
---
# T-3071 Close I-0275 with what fixed it

## Work

Run `flai issue close I-0275 --reason "<reason>"`, the reason naming the new ADR, the fallback in `scripts/flai.sh` under `guard` (T-3066), the OpenCode plugin's 126 and 127 refusal (T-3069), and the test cases that reproduce a failed build. `flai issue close` rewrites `design/issues/summary.md` too.

It waits for T-3066 and T-3069, since the reason names what they fixed. It runs beside T-3070: the two share no path.

## Done when

- I-0275 is closed, its reason saying what fixed it, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` reports nothing new.

## Notes
