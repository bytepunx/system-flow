---
id: E-0031
type: epic
nature: improvement
title: "Workflow tuning: the close-out finds what the merge queue finds, a batch is built against a clean main, hub files merge without conflict, and stories run in parallel with a trial merge instead of a hold"
status: in-progress
owner: alex
created: 2026-10-30T01:48:50Z
updated: 2026-10-31T07:22:05Z
transitions:
  - to: ready
    at: 2026-10-31T07:21:29Z
    by: alex
  - to: in-progress
    at: 2026-10-31T07:21:39Z
    by: agent-S-0548
tags: [flai]
topics: [testing, git, orchestration, planning]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 538
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 100
          output: 46177
          cache_read: 5457420
          cache_write: 180629
          cost: 3.4605
    - kind: orchestrator
      seconds: 117
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 350
          cache_read: 1417808
          cache_write: 237169
          cost: 0.3741
---
# E-0031 Workflow tuning: the close-out finds what the merge queue finds, a batch is built against a clean main, hub files merge without conflict, and stories run in parallel with a trial merge instead of a hold

## Outcome

An analysis of the fortnight 2026-10-15 to 2026-10-29, from `flai stats`, the serve journal, and 261 story agent sessions, found that the system pays for repetition and for serialization that neither adds value nor catches anything early. Per story at the median: 264 turns, 33M context tokens, 4 compactions and 5 re-primes, 11 test runs, 2 close-outs. Of 267 close-out runs, 64 died at the harness's 600 s tool timeout and 146 ended with no summary; only 88 reported every step passing, and each retry re-ran every tier. The close-out runs only the changed packages and leaves the integration and smoke tiers to the merge queue, so 24 of 104 queue verifications failed and 14 stories were sent back (22 percent of landings), each costing a restart, another close-out, and another 10-minute verification. Batches averaged 1.0 story: 80 verifications passed for 49 landings, and S-0541 alone went through 16 batches and 28 failed landings because one uncommitted note in the main checkout had a bad Markdown link that nothing checked before the batch was built. The mcpserver package takes 8 to 10 minutes under the race detector against a 10-minute default, and timing-dependent serve tests fail batches on their own (I-0225, I-0226, I-0227, I-0233).

The worktrees give isolation that the scheduler takes away. With an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never. The journal records 124 overlap holds and 113 `after` holds against 62 real conflicts, most on the same hub files: the issues summary, the flai user docs, `flai-cli.md`, the ADR index, and the template changelog, which also drew 248 overlap notices over 97 acceptances and a `stream sync` for every open story each time. The `after` chains were planner-made (S-0417 to S-0422), as the 2026-10-15 analysis already noted in I-0152.

When this epic is done, a story's close-out finds what the merge queue would find, in the worktree, as a job the agent waits on rather than a tool call that times out; a queue verification takes under five minutes; a batch is checked against the main checkout's uncommitted files before it is built, so a dirty main costs seconds once; the hub files merge without conflict and raise no overlap notice; a story starts on a trial merge instead of an overlap hold, with `after` kept for real dependencies; and acceptances land through the queue in batches of more than one. The measures above, recorded in `flai stats` and the serve journal, are the epic's before; each story records its after in the close-out.

## Stories

- S-0546 The close-out runs as a flai serve job the story's agent waits on by event, so it never dies at the harness's tool timeout and is never re-run from the start
- S-0547 The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review
- S-0548 A queue verification takes under five minutes: the slow and network-bound tests are fixed, tests run in parallel, the gateway tests leave smoke, and every go test tier sets its timeout
- S-0549 The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch
- S-0550 Hub files merge without conflict: the ADR index, the issues summary, and the changelog take a union merge or are generated, and shared paths raise no overlap notice
- S-0551 A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do
- S-0552 Acceptances land only through the merge queue, and the queue waits a window to fill a batch, so one verification lands several stories and an acceptance by hand drops no batch

## Notes

Asked by the operator in conversation on 2026-10-29 after the analysis, which lives in that conversation; the figures above are its evidence. The orchestrator's own wake cost is E-0030's and ADR-0209's, not this epic's.
- 2026-10-31T07:21:29Z: moved to ready: follows S-0548, which moved to ready
- 2026-10-31T07:21:39Z: moved to in-progress: follows S-0548, which moved to in-progress
