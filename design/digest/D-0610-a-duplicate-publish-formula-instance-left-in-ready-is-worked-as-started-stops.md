---
id: D-0610
title: A duplicate publish formula instance left in ready is worked as started, stops cleanly at its evaluate step once the other instance has released, and is accepted as a stopped instance, not cancelled
kind: decision
sources: [S-0772, TH-1161, S-0771]
topics: [release, formula, orchestration]
paths: [design/workflows/formulas/publish.md, wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0610 A duplicate publish formula instance left in ready is worked as started, stops cleanly at its evaluate step once the other instance has released, and is accepted as a stopped instance, not cancelled

On 2026-11-13 two instances of the publish formula, S-0771 (made by system-flow at 13:41:40Z) and S-0772 (made by the orchestrator's run at 13:42:43Z, commit 17d68ecb7), sat in ready for the same batch, S-0661. The orchestrator opened TH-1161 recommending "cancel S-0772", and said: "When I ran the formula, I had not checked for an instance already open." Cancelling a story is not in its permissions, so it moved neither.

S-0772's agent recorded its decision in its narrative (`wip/archive/agents/S-0772.md`): "Cancelling is the operator's; I work the formula as started, and its evaluate step stops it cleanly if S-0771 already released everything." So it went: `flai release --evaluate --require-met` exited 2, count 0 of 1, since S-0771 had already published flai 1.92.0 and template 1.1.37; `flai formula stop S-0772 --step evaluate --exit 2` closed T-3304 and skipped T-3305; nothing was tagged or pushed.

The orchestrator then resolved TH-1161: "It is accepted as a formula instance that stopped cleanly, not cancelled", and accepted S-0772 at 15:22:35Z, as [ADR-0271](../adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md) allows. The duplicate cost one agent run and nothing else. Before running `flai formula run publish`, check that no publish instance is already open for the batch.

Sources: S-0772, TH-1161, S-0771.
