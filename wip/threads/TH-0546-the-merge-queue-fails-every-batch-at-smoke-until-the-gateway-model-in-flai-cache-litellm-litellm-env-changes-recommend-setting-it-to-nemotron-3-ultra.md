---
id: TH-0546
title: "The merge queue fails every batch at smoke until the gateway model in .flai-cache/litellm/litellm.env changes: recommend setting it to nemotron-3-ultra"
anchor:
  path: wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md
  item: S-0454
status: open
participants: [agent-S-0454, orchestrator]
created: 2026-10-21T01:54:25Z
updated: 2026-10-21T01:54:39Z
---

# TH-0546 The merge queue fails every batch at smoke until the gateway model in .flai-cache/litellm/litellm.env changes: recommend setting it to nemotron-3-ultra

On wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md.

## Entries

### 2026-10-21T01:54:25Z agent-S-0454
TH-0545's send-back of S-0454 is not caused by S-0454's change. The queue's smoke failed at the LiteLLM gateway check, the I-0171 failure (now I-0172 records why the thread showed only "template render and check").

**Why.** flai serve's environment has no `LITELLM_BASE_URL`, so `scripts/env.sh` loads `.flai-cache/litellm/litellm.env`. That file still sets `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free`, the retired model Claude Code 2.1.296 refuses. Agents' sessions set `LITELLM_BASE_URL`, so they never load the file, which is why S-0445 thought it did not matter.

**Evidence.** S-0454's close-out passed every tier, smoke included, at the batch's commit `6a4ae3b52fd6` with the variable set to nemotron-3-ultra for that run. The template render and check passes in a fresh worktree at that commit, and again with a copy of flai serve's config.

**Recommendation.** In `.flai-cache/litellm/litellm.env`, change that line to:

```text
FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free
```

That is the model S-0451's close-out passed smoke with. S-0455's T-1941 makes `scripts/litellm.sh` stop freezing the model into this file; until it lands, every batch fails smoke the same way. The file is the host's, so I have not changed it.

I am closing S-0454 out again and moving it back to review now. Enqueue it again once the line is changed. With no change, the next batch will send it back the same way.

### 2026-10-21T01:54:39Z orchestrator (recommendation)
**Recommendation:** in `.flai-cache/litellm/litellm.env` on the host, change the line to

```text
FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free
```

This is the host-side half of what you confirmed on TH-0534. I will not enqueue S-0454 again, or any other story, until you say it is done.

**Why this needs you.** The file is the host's, and no agent changes it.

**Why it blocks every story.** Under ADR-0176 a batch "that failed" is bisected, and the failing story is sent back to in-progress. flai serve loads this file because its environment has no `LITELLM_BASE_URL`. So every batch fails smoke on the retired model and sends its story back, however sound the story is. S-0454 just went through that at a commit that passed smoke in its own verify.

**Meanwhile, if you want:**
- Run `flai queue hold` so no batch runs until the line is fixed. I cannot hold the queue: the guard leaves me only `queue add`.
- S-0455's T-1941 is the lasting fix, and it is held behind S-0451. Resolving TH-0529, or accepting S-0451, unblocks it.

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision
