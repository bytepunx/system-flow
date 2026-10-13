---
id: T-1410
type: task
nature: feature
title: flai-cli.md describes the spend readers, when flai serve reads them, and priced_by
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:51:10Z
updated: 2026-10-13T20:15:38Z
transitions:
  - to: ready
    at: 2026-10-13T20:06:59Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T20:07:00Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T20:15:38Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [design/system/flai-cli.md, design/system/project-manifest.md, docs/operators/index.md, flai/cmd/serve_usage.go, docs/users/flai-reference.md]
after: [T-1409]
usage:
  source: log
  seconds: 518
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 27187
      cache_read: 3337987
      cache_write: 101512
      cost: 1.8218
---
# T-1410 flai-cli.md describes the spend readers, when flai serve reads them, and priced_by

## Work

- `design/system/flai-cli.md`, where usage is measured: the `Spend` contract, the OpenRouter and LiteLLM readers and how each joins a run, how the reader is chosen by base URL, `spend_key_env`, the spend record beside each log and when it is written, when `flai serve` asks a gateway (each ended run, once, with its timeout), the apportioning over models, the fallback and its `spend not read` log line, the spend keys kept out of agents' environments, the tasks scaled by the story's ratio, and `priced_by` with its sum; cite ADR-0132.
- `design/system/project-manifest.md`: the host's provider entry gains `spend_key_env`, the host's alone like `key_env`.
- `docs/operators/index.md`, beside a provider's key: the spend key, what it is for, and that agents never get it.
- `flai serve agent usage`'s help in `flai/cmd/serve_usage.go`: runs through a provider are priced from the gateway's spend log, and `--json` gains `unread`; then `make flai-reference` regenerates `docs/users/flai-reference.md` (never hand-edit it).

It waits for T-1409, so that it describes the pricing as built.

## Done when

- The sections match the code.
- `flai test design/system/flai-cli.md design/system/project-manifest.md docs/operators/index.md flai/cmd/serve_usage.go docs/users/flai-reference.md` passes.

## Notes

Layer 4 of S-0357. `metrics.md` changes with the next story, which shows `priced_by`. Meets criterion 6.
