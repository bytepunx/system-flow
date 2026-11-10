---
id: T-3068
type: task
nature: improvement
title: I-0265 is closed with the reason that the operator's outside folders reach a story's agent and its sub-agents and Monitor is allowed
status: backlog
parent: S-0701
owner: alex
created: 2026-11-10T10:35:57Z
updated: 2026-11-10T10:35:57Z
transitions: []
stream: S-0701
tags: [issues]
touches: [design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, design/issues/summary.md]
after: [T-3067]
---
# T-3068 I-0265 is closed with the reason that the operator's outside folders reach a story's agent and its sub-agents and Monitor is allowed

## Work

Run `flai issue close I-0265 --reason "<reason>"`. The reason names the ADR from T-3048, `flai serve agent harness --outside`, the `--add-dir` and container mounts it gives a story's agent and its sub-agents, `Monitor` in the default `--allowedTools`, and the tests that reproduce the cause. flai rewrites `design/issues/summary.md`.

This task waits for T-3067, so that the issue closes once the fix and its documentation are complete.

## Done when

- I-0265 is closed, with that reason, and `design/issues/summary.md` lists it as closed.
- `flai check --strict` is clean.

## Notes
