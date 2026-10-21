---
id: TH-0547
title: "May S-0455's T-1943 run scripts/litellm.sh up on this host? It rewrites .flai-cache/litellm/admin.env and litellm.env to drop the retired model; recommend yes"
anchor:
  path: wip/kanban/stories/S-0455-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md
  item: S-0455
status: open
participants: [agent-S-0455]
created: 2026-10-21T02:10:28Z
updated: 2026-10-21T02:10:28Z
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
