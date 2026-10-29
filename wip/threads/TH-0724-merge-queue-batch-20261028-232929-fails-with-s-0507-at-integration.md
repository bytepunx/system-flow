---
id: TH-0724
title: Merge queue batch 20261028-232929 fails with S-0507 at integration
anchor:
  path: wip/kanban/stories/S-0507-pilot-a-go-binary-in-the-bedrock-image-is-the-container-s-entrypoint-it-reads-a-start-request-starts-the-inner-harness-through-flai-s-adapters-relays-its-log-and-resumes-as-an-adr-records-beside-the-layered-images.md
  item: S-0507
status: resolved
participants: [flai serve]
created: 2026-10-28T23:41:14Z
updated: 2026-10-28T23:46:37Z
---

# TH-0724 Merge queue batch 20261028-232929 fails with S-0507 at integration

On wip/kanban/stories/S-0507-pilot-a-go-binary-in-the-bedrock-image-is-the-container-s-entrypoint-it-reads-a-start-request-starts-the-inner-harness-through-flai-s-adapters-relays-its-log-and-resumes-as-an-adr-records-beside-the-layered-images.md.

## Entries

### 2026-10-28T23:41:14Z flai serve
Merge queue batch `20261028-232929` fails with S-0507 at `integration`. The queue verified S-0507 alone on main at `229ac6a30909`, at `0be963c7433f`, the copy of its verified head `0be963c7433f` rebased in the batch, and `integration` failed:

```text
flai/internal/serve/agents_test.go:2682 TestAStorySentBackIsStartedAfreshOrBeforeARestart/no_run,_then_a_run_without_a_session: journal: [{At:2026-10-28T23:30:20Z Action:agent Method:serve.agent Project:t Root:/tmp/TestAStorySentBackIsStartedAfreshOrBeforeARestartno_run,_then_a1171893710/001 By:flai serve RequestID: Outcome:done Detail:started stub-agent (claude-code) again for S-0001 as builder-S-0001 in a new session, flai having sent it back from review on TH-0001 (pid 3607048); log /tmp/TestAStorySentBackIsStartedAfreshOrBeforeARestartno_run,_then_a1171893710/003/serve/agents/t-S-0001-20261028T233020Z.log} {At:2026-10-28T23:30:20Z Action:agent Method:serve.agent Project:t Root:/tmp/TestAStorySentBackIsStartedAfreshOrBeforeARestartno_run,_then_a1171893710/001 By:flai serve RequestID: Outcome:done Detail:opened TH-0002: S-0001's agent ended once in progress, and flai serve restarts it no more}]
flai/internal/serve/forgemerge_test.go:310 TestRequestMergeIsAnsweredByFlaiServesMergePass: the log of the merge, which must not hold the token:
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261028-232929`.

flai moved S-0507 back to in-progress and took it out of the merge queue. Fix it on `story/S-0507`, run flai verify S-0507, move it to review, and enqueue it again with flai queue add S-0507.

### 2026-10-28T23:46:37Z flai serve
Resolved: S-0507 is back in review at b1cb2776dedd, a new commit; it is batched again once it is enqueued
