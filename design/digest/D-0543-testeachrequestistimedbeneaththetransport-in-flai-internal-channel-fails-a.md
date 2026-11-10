---
id: D-0543
title: "TestEachRequestIsTimedBeneathTheTransport in flai/internal/channel fails a close-out's integration tier under load since S-0611, whatever the story changes; I-0324 records it and the close-out runs it again"
kind: cause
sources: [S-0616, I-0324]
topics: [testing, flai, go]
paths: [flai/internal/channel]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0543 TestEachRequestIsTimedBeneathTheTransport in flai/internal/channel fails a close-out's integration tier under load since S-0611, whatever the story changes; I-0324 records it and the close-out runs it again

`TestEachRequestIsTimedBeneathTheTransport` in `flai/internal/channel` reads the request-answered event before the channel client has logged it, so it fails under load. The race came to main with S-0611.

S-0616's close-out stopped at integration on it, on a branch that does not touch `flai/internal/channel`. Its narrative, § Decisions: "I recorded it as I-0324 and did not fix it here, since it is outside the story. The next close-out passed."

[I-0324](../issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md) is open and counts each close-out it stops. A story that does not touch the channel package bumps I-0324 and asks for the close-out again, rather than fixing the test inside its own story.

Sources: S-0616's narrative, § Decisions; I-0324.
