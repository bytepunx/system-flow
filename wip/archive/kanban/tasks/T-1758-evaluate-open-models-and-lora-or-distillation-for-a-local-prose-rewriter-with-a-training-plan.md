---
id: T-1758
type: task
nature: research
title: Evaluate open models and LoRA or distillation for a local prose rewriter, with a training plan
status: done
parent: S-0413
owner: alex
created: 2026-10-16T05:46:21Z
updated: 2026-10-16T06:04:23Z
transitions:
  - to: ready
    at: 2026-10-16T05:46:56Z
    by: agent-S-0413
  - to: in-progress
    at: 2026-10-16T06:03:01Z
    by: agent-S-0413
  - to: done
    at: 2026-10-16T06:04:23Z
    by: agent-S-0413
stream: S-0413
tags: []
touches: [design/system/agent-prose.md]
after: [T-1755]
usage:
  source: log
  seconds: 82
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 22408
      cache_read: 2422734
      cache_write: 85061
      cost: 1.4586
---
# T-1758 Evaluate open models and LoRA or distillation for a local prose rewriter, with a training plan

## Work

Evaluate open-weight models that could be fine-tuned with LoRA or QLoRA, distilled, or quantized to run on consumer hardware as a rewriter: verbose agent prose in, short plain prose that passes the chosen linter out. Cover the candidate model families and sizes, the hardware each needs at inference and at training, local serving (llama.cpp, Ollama, vLLM, MLX) and how flai would call one, a training plan (data from this repository's own prose, rewrites by a large model filtered by the linter, the linter's score as a filter or a reward), how to measure the meaning kept, the costs, and the risks.

Waits for T-1755: the training plan scores against the linter it recommends.

## Done when

- `design/system/agent-prose.md` has a `## A local rewriting model` section with the model comparison, the hardware figures, a training and evaluation plan, and a recommendation.
- Acceptance criterion 2 is met.

## Notes
