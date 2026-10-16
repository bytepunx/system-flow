---
id: T-1747
type: task
nature: research
title: Evaluate open models that a LoRA or distillation can fit to consumer hardware as a prose rewriter, with a training plan
status: cancelled
parent: S-0413
owner: alex
created: 2026-10-16T05:45:45Z
updated: 2026-10-16T06:01:28Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:01:28Z
    by: agent-S-0413
stream: S-0413
tags: [prose, research, models]
touches: [design/system/agent-prose.md]
after: [T-1745]
---
# T-1747 Evaluate open models that a LoRA or distillation can fit to consumer hardware as a prose rewriter, with a training plan

## Work

Add a section on a local rewriting model to `design/system/agent-prose.md`, meeting the story's second criterion.

- Candidates: current open-weight instruction models from about 1B to about 14B parameters, such as the Qwen, Llama, Gemma, Phi, and Mistral families, and a larger one reduced by distillation or quantisation. Name each one's licence, size, and the memory it needs at 4-bit and 8-bit.
- Compare the ways to fit one to the task on consumer hardware: LoRA and QLoRA fine-tuning, distillation from a frontier model, preference tuning (DPO) with the linter's score as the signal, and prompting alone as the baseline. Name the tools (such as Unsloth, Axolotl, or PEFT) and the hardware and hours each takes.
- Write a training plan: the dataset, pairs of verbose agent text and a rewrite that passes T-1745's linter, made from the repository's own narratives, threads, and notes; how the linter filters and scores rewrites; how meaning is checked to be kept, such as IDs, paths, and figures surviving the rewrite; and the evaluation set and pass threshold.
- Say how flai would call the model when one is available: a local server with an OpenAI-compatible API (such as Ollama, llama.cpp, or vLLM) named as a provider as `design/system/agent-adapters.md` describes, where in flai the rewrite would run, and what happens when no model is running.
- Where it can be done in a few minutes in a scratch folder, try a small model with a prompt on part of T-1742's sample and score it with the linter, as the baseline any fine-tune must beat. Say so when it cannot be run here.

Waits for T-1745: the training plan needs the chosen linter, and it writes the same document.

## Done when

- The model section has the candidate table, the comparison of approaches with hardware and time, the training plan, and how flai would call the model.
- A recommendation names the model and approach to try first, or says why none suits yet.

## Notes
- 2026-10-16T06:01:28Z: moved to cancelled: Duplicates agent-S-0413's plan (T-1750, T-1755, T-1758, T-1760), as the planner recommends on TH-0472
