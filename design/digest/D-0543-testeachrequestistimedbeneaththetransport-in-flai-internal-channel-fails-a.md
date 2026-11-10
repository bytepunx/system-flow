---
id: D-0543
title: "TestEachRequestIsTimedBeneathTheTransport in flai/internal/channel fails a close-out's integration tier under load since S-0611, whatever the story changes; I-0324 records it and the close-out runs it again"
kind: cause
sources: [S-0616, I-0324, S-0649, T-2834]
topics: [testing, flai, go]
paths: [flai/internal/channel, flai/internal/channel/channel_test.go]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0543 TestEachRequestIsTimedBeneathTheTransport in flai/internal/channel fails a close-out's integration tier under load since S-0611, whatever the story changes; I-0324 records it and the close-out runs it again

`TestEachRequestIsTimedBeneathTheTransport` in `flai/internal/channel` read the request-answered event before the channel client had logged it, so it failed under load. The race came to main with S-0611.

S-0616's close-out stopped at integration on it, on a branch that does not touch `flai/internal/channel`. Its narrative, § Decisions: "I recorded it as I-0324 and did not fix it here, since it is outside the story. The next close-out passed." Until it was fixed, a story that did not touch the channel package bumped I-0324 and asked for the close-out again, rather than fixing the test inside its own story.

S-0649 fixed it and closed I-0324 on 2026-11-10 (T-2834). From its Decisions (`wip/archive/agents/S-0649.md`):

```text
The fix is `answeredEvent` in `flai/internal/channel/channel_test.go`. It polls the log for the "request answered" event, with a 5 s deadline. Both request-answered tests use it in place of their own reads.
No product code changes. `channel.go` logs after the write, as S-0611 meant it to.
```

The race is reproduced without load by `lockedBuffer.holdAnswered`, which holds back the "request answered" line for 200 ms: with the old single read the test failed 3 of 3 runs, and with the polled read it passes. So a close-out that fails on this test again is a new cause, not I-0324.

Sources: S-0616's narrative, § Decisions; I-0324; S-0649 (`wip/archive/agents/S-0649.md`, Decisions); T-2834.
