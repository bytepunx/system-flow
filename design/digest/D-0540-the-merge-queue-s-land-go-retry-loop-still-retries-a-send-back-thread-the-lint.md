---
id: D-0540
title: "The merge queue's land.go retry loop still retries a send-back thread the lint refuses, since once mdlint.Fenced quotes it a refusal that remains is a real fault; a queue that holds itself after repeated refusals is a story of its own"
kind: decision
sources: [S-0614, TH-0923, I-0242]
topics: [merge-queue, markdown-lint, flai, code]
paths: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/bisect.go]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0540 The merge queue's land.go retry loop still retries a send-back thread the lint refuses, since once mdlint.Fenced quotes it a refusal that remains is a real fault; a queue that holds itself after repeated refusals is a story of its own

I-0242 found the merge queue retrying forever, with the batch wedged, when the markdown lint refused its send-back thread. S-0614 fixed the quote (`mdlint.Fenced`) and the fallback: when the lint still refuses the thread, `SendBack` opens it again with the findings left out (`design/system/flai-cli.md`). It left the retry loop in `flai/internal/mergequeue/land.go` as it is.

The planner's reason, on TH-0923: "Once a quote can no longer refuse the thread, a refusal that remains is a real fault worth retrying. If you would rather the queue hold itself with one thread after repeated refusals, as the preflight does (ADR-0234), that is a second story." The orchestrator approved the plan as it stood, and no story for that hold exists yet.
