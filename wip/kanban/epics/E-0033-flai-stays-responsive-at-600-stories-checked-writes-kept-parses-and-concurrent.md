---
id: E-0033
type: epic
nature: improvement
title: "flai stays responsive at 600 stories: checked writes, kept parses, and concurrent reads"
status: in-progress
owner: alex
created: 2026-11-07T22:38:45Z
updated: 2026-11-09T07:36:47Z
transitions:
  - to: ready
    at: 2026-11-07T23:49:08Z
    by: alex
  - to: in-progress
    at: 2026-11-07T23:49:17Z
    by: agent-S-0598
tags: [cli, dashboard]
topics: [server-side, back-end]
usage:
  source: sum
  seconds: 65813
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 3
      work: 57
    - day: 2026-11-08
      ceremony: 28
      test_runs: 5
      hand_edits: 27
      work: 939
    - day: 2026-11-09
      ceremony: 13
      test_runs: 13
      hand_edits: 13
      work: 565
  models:
    - model: claude-opus-5-5
      input: 5450
      output: 1241827
      cache_read: 940492805
      cache_write: 22663077
      cost: 367.1884
  strategic:
    - kind: planner
      seconds: 1958
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 338
          output: 38199
          cache_read: 14666696
          cache_write: 822422
          cost: 1.048
        - model: claude-opus-5-5
          input: 2077
          output: 137305
          cache_read: 177386605
          cache_write: 2725272
          cost: 52.2198
    - kind: orchestrator
      seconds: 4466
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 599
          output: 5055
          cache_read: 26130818
          cache_write: 2906057
          cost: 7.6121
        - model: claude-sonnet-5-5
          input: 84
          output: 3113
          cache_read: 700757
          cache_write: 450660
          cost: 0.8488
---
# E-0033 flai stays responsive at 600 stories: checked writes, kept parses, and concurrent reads

## Outcome

With about 600 stories and 3,200 work items, creating or editing an item takes about 2 s rather than 7 to 18 s, an acceptance's replan lands in seconds rather than a minute or more, a board answers in under 300 ms under load, and the charts answer warm in about 200 ms at any window. Every fix keeps the files under `wip/` and `design/` as they are; what it keeps is in memory or in a cache under `.flai-cache` that can be rebuilt.

S-0596 measured the causes in [server-performance.md § A second pass at 596 stories (S-0596)](../../../design/system/server-performance.md#a-second-pass-at-596-stories-s-0596). Each story below is one cause, 8 to 21, with its numbers and its proposed fix.

## Stories
- S-0598 A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s
- S-0599 flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s
- S-0600 A flai process parses work items on a pool of workers and keeps parses and lint results under .flai-cache, so a cold list takes under 0.3 s, not 1.1 to 1.5 s
- S-0601 flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s
- S-0602 flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms
- S-0603 agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes about 250 ms, not 0.8 to 2.4 s
- S-0604 board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing
- S-0605 settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms
- S-0606 item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms
- S-0607 stats.get runs its reads at once and its sections at once, and answers identical requests in flight once
- S-0608 The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does
- S-0609 flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance
- S-0610 The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out
- S-0611 The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB
- S-0636 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606

## Notes

- Causes 8 and 9 cost the operator most: every dashboard write and every acceptance's replan. Cause 12 costs flai serve most: about 30% of its wall time.
- Each story remeasures its cause on the same host and repository, as E-0012's did, and records it in server-performance.md.
- 2026-11-07T23:49:08Z: moved to ready: follows S-0598, which moved to ready
- 2026-11-07T23:49:17Z: moved to in-progress: follows S-0598, which moved to in-progress
