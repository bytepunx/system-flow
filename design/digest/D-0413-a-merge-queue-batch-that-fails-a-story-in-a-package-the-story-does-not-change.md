---
id: D-0413
title: A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way
kind: decision
sources: [S-0507, TH-0723, TH-0724, I-0225, I-0226]
topics: [merge-queue, testing, go, code]
paths: [scripts/integration.sh, flai/internal/mcpserver, flai/internal/serve/agents_test.go, flai/internal/serve/forgemerge_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0413 A merge queue batch that fails a story in a package the story does not change is fixed inside the story, not re-enqueued unchanged, since the next loaded batch fails the same way

The merge queue sent S-0507 back twice, both times for failures in packages S-0507 does not change. Each time its agent fixed the cause on the story's branch instead of enqueueing the same commit again. Its reason: "the next batch would likely fail the same way."

- TH-0723: batch 20261028-231632 failed the integration tier because `flai/internal/mcpserver` hit Go's 10-minute default test timeout (600.17 s). The package took 470 to 591 s under `-race` in every batch that day. T-2316 gave `scripts/integration.sh` `-timeout 20m`, and the script's comment gives the reason. The suite's slowness is still open as I-0225.
- TH-0724: the next batch passed mcpserver at 679 s. It then failed two `flai/internal/serve` tests that pass when run alone. Each checked for something before it had happened. T-2317 made one test check for one start and the other wait for the log line. I-0226 records it.

A later story sent back for a failure outside its own change should do the same. Find the cause first, in the batch log and in the package's times across batches. If the cause is timing or load, fix the timeout or the test in the story. Record the underlying cost as an issue for a story of its own. Enqueue again only after that.

Source: S-0507 (wip/archive/agents/S-0507.md, Decisions), TH-0723, TH-0724.
