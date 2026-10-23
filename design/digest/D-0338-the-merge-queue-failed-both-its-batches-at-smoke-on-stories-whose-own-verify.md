---
id: D-0338
title: "The merge queue failed both its batches at smoke on stories whose own verify passed: recommend holding it until S-0457 shows why: Recommendation: hold the queue (`flai queue hold`) until one of the…"
kind: answer
sources: [TH-0552]
topics: [flai, cli, go, code, all]
paths: [scripts/litellm.sh, scripts/README.md, flai/internal/serve/gateway_test.go, docs/contributors/index.md, design/system/agent-adapters.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0338 The merge queue failed both its batches at smoke on stories whose own verify passed: recommend holding it until S-0457 shows why: Recommendation: hold the queue (`flai queue hold`) until one of the…

The operator, alex, answered TH-0552 at 2026-10-21T06:59:42Z.

What was asked:

```text
**Recommended:** run `flai queue hold`, or `flai serve disable merge-queue`, until S-0457 lands, so you can see what the batch's smoke step really fails on. While the queue is held I will go back to accepting stories directly, as I did before. Its accept refusal names the queue only while the queue batches.
```

The answer:

```text
**Recommendation:** hold the queue (`flai queue hold`) until one of the two fixes below is in. S-0455 is sound: the batch would fail with any story.

**Cause, reproduced in S-0455's worktree (I-0174).**
- `scripts/check.sh`, a step of the smoke tier, scopes `flai check --strict` to the story only when `CLOSE_OUT_STORY` is set.
- The batch's verification sets no `CLOSE_OUT_STORY` (flai-cli.md, The merge queue pass), so smoke runs an unscoped strict check over the whole repository.
- Unscoped, it fails on six warnings outside any batched story:

| Finding | What clears it |
|---------|----------------|
| `item.archive` S-0452, cancelled but not archived | `flai archive` in the main checkout |
| `threads.archived` TH-0526, open on archived S-0404 | resolve TH-0526 |
| `issues.no-story` I-0143, I-0144, I-0145 | a story for each, or close them |
| `components.examples` for flaiover, svelte not installed | install flaiover's dependencies |

- `env -u CLOSE_OUT_STORY scripts/check.sh` exits 1. With `CLOSE_OUT_STORY=S-0455` the same check passes. The rest of smoke (template render, markdown lint, installer, LiteLLM gateway) passed with flai serve's PATH and keys.
- The batch's quote, "building bin/flai", is only the first line of `flai.sh`'s stderr (I-0172), which is why it pointed nowhere.

**Fixes, either one:**
1. Lasting: the batch verification scopes the check to its stories, by setting `CLOSE_OUT_STORY` or by `check.sh` taking the batch's stories. That is a story from I-0174, a natural neighbour of S-0457.
2. Today: clear the findings in the table on main. A batch then passes until main carries another finding a scoped check passes over.

I am closing S-0455 out again, at a new commit that only records I-0174, and moving it back to review. Enqueue it once the queue can pass.

Source: design/system/flai-cli.md § The merge queue pass
```

Source: TH-0552.
