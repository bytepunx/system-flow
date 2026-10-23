---
id: D-0335
title: "May S-0455's T-1943 run scripts/litellm.sh up on this host? It rewrites .flai-cache/litellm/admin.env and litellm.env to drop the retired model; recommend yes: Recommendation: yes, allow S-0455's…"
kind: answer
sources: [TH-0547]
topics: [flai, cli, go, code, all]
paths: [scripts/litellm.sh, scripts/README.md, flai/internal/serve/gateway_test.go, docs/contributors/index.md, design/system/agent-adapters.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0335 May S-0455's T-1943 run scripts/litellm.sh up on this host? It rewrites .flai-cache/litellm/admin.env and litellm.env to drop the retired model; recommend yes: Recommendation: yes, allow S-0455's…

The operator, alex, answered TH-0547 at 2026-10-21T05:57:02Z.

What was asked:

```text
**Recommendation: yes, allow it.**
```

The answer:

```text
**Recommendation:** yes, allow S-0455's T-1943 to run the fixed `scripts/litellm.sh up` on this host. One answer settles both this thread and TH-0546.

**Why this needs you.** It rewrites files on the host, so it is your call. It is also exactly the change TH-0546 asks for: `litellm.env` gets the current model, and `admin.env` drops the retired default that the old script froze in. Until then, flai serve's merge queue fails every batch at smoke, and under ADR-0176 each failed batch sends its story back to in-progress.

**Why yes is safe.**
- The agent lists everything the run keeps: the master key, salt key, Postgres password, upstream keys, port, image and budget.
- `config.yaml`, `compose.yaml` and `docker compose up -d` are what every `up` already does, and the containers stay as they are.
- The one value removed is the script's own default, not a choice you made.

If you would rather edit `litellm.env` yourself, as TH-0546 suggests, answer no here, and T-1943 records that the files were left alone.

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision
```

Source: TH-0547.
