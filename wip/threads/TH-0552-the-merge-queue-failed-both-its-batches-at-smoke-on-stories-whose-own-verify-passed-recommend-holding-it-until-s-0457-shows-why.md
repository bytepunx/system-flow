---
id: TH-0552
title: "The merge queue failed both its batches at smoke on stories whose own verify passed: recommend holding it until S-0457 shows why"
anchor:
  path: wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md
  item: S-0455
status: resolved
participants: [orchestrator, agent-S-0455, alex]
created: 2026-10-21T06:25:15Z
updated: 2026-10-21T07:19:47Z
---

# TH-0552 The merge queue failed both its batches at smoke on stories whose own verify passed: recommend holding it until S-0457 shows why

On wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md.

## Entries

### 2026-10-21T06:25:15Z orchestrator
**Recommended:** run `flai queue hold`, or `flai serve disable merge-queue`, until S-0457 lands, so you can see what the batch's smoke step really fails on. While the queue is held I will go back to accepting stories directly, as I did before. Its accept refusal names the queue only while the queue batches.

## What happened

| Batch | Story | Smoke in the story's own verify | Batch's smoke, as quoted |
|-------|-------|---------------------------------|--------------------------|
| 20261021-014447 | S-0454 | passed at 01:31Z, same commit | "template render and check" |
| 20261021-061843 | S-0455 | passed at 06:07Z, same commit | "building bin/flai" |

- For S-0454, its agent traced the failure to the host's `litellm.env`, and S-0455 fixed that. The second batch still failed, this time on a different first line.
- Both quotes are only a first line, which is the gap I-0172 records. S-0457 is fixing it now: the thread will quote the full findings and keep each batch's logs. Until then neither the agents nor I can tell what fails.
- Each failure sends a sound story back to in-progress, costing its agent a close-out of 10 to 15 minutes.

## Why it is yours

Holding or disabling the queue is a host-level choice, and the guard leaves me only `queue add`. If you would rather keep the queue running, say so. I will keep enqueuing, and each story's agent will re-close-out as each batch sends it back.

Meanwhile S-0455 goes back to its agent per TH-0551.

### 2026-10-21T06:30:09Z agent-S-0455 (recommendation)
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

### 2026-10-21T06:59:42Z alex
Confirmed the recommendation of 2026-10-21T06:30:09Z agent-S-0455.

Source: design/system/flai-cli.md § The merge queue pass

### 2026-10-21T07:19:47Z alex
Resolved: S-0455 was accepted
