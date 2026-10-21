---
id: TH-0547
title: "May S-0455's T-1943 run scripts/litellm.sh up on this host? It rewrites .flai-cache/litellm/admin.env and litellm.env to drop the retired model; recommend yes"
anchor:
  path: wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md
  item: S-0455
status: resolved
participants: [agent-S-0455, orchestrator, alex]
created: 2026-10-21T02:10:28Z
updated: 2026-10-21T06:00:53Z
---

# TH-0547 May S-0455's T-1943 run scripts/litellm.sh up on this host? It rewrites .flai-cache/litellm/admin.env and litellm.env to drop the retired model; recommend yes

On wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md.

## Entries

### 2026-10-21T02:10:28Z agent-S-0455
**Recommendation: yes, allow it.**

**What it does.** T-1943 runs the fixed `scripts/litellm.sh up` from S-0455's worktree, once T-1941 and T-1942 are done. On this host it will:

- remove `LITELLM_HAIKU_MODEL=openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free` from `.flai-cache/litellm/admin.env`. The script wrote that value itself on the first `up`; it is the retired default, not your choice.
- rewrite `.flai-cache/litellm/litellm.env` with `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`, keeping the virtual key.
- rewrite `config.yaml` and `compose.yaml` as every `up` does, and run `docker compose up -d` on the running proxy. With no change to the image or the config, the containers stay as they are.

The master key, salt key, Postgres password, upstream keys, port, image, and budget stay as they are.

**Why it needs you.** The files are the host's. This is also the host-side change TH-0546 asks you for. Once it is done, flai serve's merge queue batches stop failing smoke.

**Then.** I run the smoke tier with no `FLAI_TEST_LITELLM_HAIKU` set by hand, record the result in `design/system/agent-adapters.md`, and close I-0171.

**If you answer no**, or edit the line yourself first, I run the check against whatever the files name, and record that the host's files were left alone.

### 2026-10-21T03:25:07Z orchestrator (recommendation)
**Recommendation:** yes, allow S-0455's T-1943 to run the fixed `scripts/litellm.sh up` on this host. One answer settles both this thread and TH-0546.

**Why this needs you.** It rewrites files on the host, so it is your call. It is also exactly the change TH-0546 asks for: `litellm.env` gets the current model, and `admin.env` drops the retired default that the old script froze in. Until then, flai serve's merge queue fails every batch at smoke, and under ADR-0176 each failed batch sends its story back to in-progress.

**Why yes is safe.**
- The agent lists everything the run keeps: the master key, salt key, Postgres password, upstream keys, port, image and budget.
- `config.yaml`, `compose.yaml` and `docker compose up -d` are what every `up` already does, and the containers stay as they are.
- The one value removed is the script's own default, not a choice you made.

If you would rather edit `litellm.env` yourself, as TH-0546 suggests, answer no here, and T-1943 records that the files were left alone.

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision

### 2026-10-21T05:57:02Z alex
Confirmed the recommendation of 2026-10-21T03:25:07Z orchestrator.

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision

### 2026-10-21T06:00:53Z agent-S-0455
Resolved: Done on alex's confirmation: the fixed scripts/litellm.sh up dropped the retired slug from .flai-cache/litellm/admin.env and wrote nemotron-3-ultra into litellm.env, keeping keys, port, image, and budget; config.yaml and compose.yaml were unchanged. TestGatewayLiteLLM then passed in 145 s from the host's own litellm.env.
