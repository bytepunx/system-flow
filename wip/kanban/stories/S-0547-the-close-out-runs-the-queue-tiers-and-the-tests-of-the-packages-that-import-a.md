---
id: S-0547
type: story
nature: improvement
title: The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review
status: backlog
parent: E-0031
owner: alex
created: 2026-10-30T01:52:02Z
updated: 2026-10-30T01:52:02Z
transitions: []
tags: [flai]
topics: [testing, cli]
touches: [flai/internal/check, flai/internal/verify, flai/cmd, system-flow.yaml, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0547 The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review

## Goal

The close-out's go-test tier runs only the changed packages, and ADR-0174 leaves the `queue: true` tiers, integration and smoke, to the merge-queue batch. In the fortnight 24 of 104 queue verifications failed and 14 stories were sent back, 22 percent of landings, each costing the agent a restart, a fix, another close-out, a re-enqueue, and another 10-minute verification. I-0221 records a change that broke a caller package's tests and passed its close-out; I-0203 a stale generated file that failed every later batch at integration. This story makes the close-out find what the queue finds: the go-test tier runs the changed packages and every package that imports them (`go list -deps` in reverse), and the queue tiers run in the close-out too, in the worktree, when the story's touches reach what they test, with the batch keeping them for the landing check. With the close-out as a serve job, their time no longer hits a tool timeout.

## Acceptance criteria
- [ ] The go-test tier's `{packages}` expands to the changed packages and the packages that import them, and a test shows a change that breaks a caller's tests fails the close-out
- [ ] A tier with `queue: true` runs in the close-out when a `when:` the manifest gives it matches the story's touches, and in the batch as today; the manifest schema, `flai check`, and the reference describe `when`
- [ ] system-flow.yaml gives integration and smoke a `when` that reaches them from the flai and template paths, and I-0221 and I-0203 are closed by the story that proves it
- [ ] The testing convention, `flai-cli.md`, and the user docs say which tiers the close-out runs and which the batch repeats

## Tasks

The planner drafts them.

## Notes

Evidence: `flai stats` merge-queue figures for 2026-10-15 to 2026-10-29 and the serve journal's send-backs, in the operator's conversation on 2026-10-29.
