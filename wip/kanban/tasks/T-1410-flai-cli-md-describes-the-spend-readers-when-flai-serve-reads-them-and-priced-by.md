---
id: T-1410
type: task
nature: feature
title: flai-cli.md describes the spend readers, when flai serve reads them, and priced_by
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:51:10Z
updated: 2026-10-08T08:51:10Z
transitions: []
stream: S-0357
tags: [cli]
touches: [design/system/flai-cli.md]
after: [T-1409]
---
# T-1410 flai-cli.md describes the spend readers, when flai serve reads them, and priced_by

## Work

- `design/system/flai-cli.md`, where usage is measured: the `Spend` contract, the OpenRouter and LiteLLM readers and how each joins a run, `spend_key_env`, the fallback and its log line, and `priced_by` with its sum; cite ADR-0132.

It waits for T-1409, so that it describes the pricing as built.

## Done when

- The section matches the code.
- `flai test design/system/flai-cli.md` passes.

## Notes

Layer 4 of S-0357. `metrics.md` changes with the next story, which shows `priced_by`.
