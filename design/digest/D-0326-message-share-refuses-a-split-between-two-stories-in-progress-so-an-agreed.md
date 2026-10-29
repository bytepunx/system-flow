---
id: D-0326
title: message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
kind: cause
sources: [I-0169]
topics: [flai, cli, go, code, all]
paths: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md, design/issues/summary.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, docs/contributors/index.md, docs/operators/index.md, flai/internal/serve/gateway_test.go, scripts/litellm.sh]
date: 2026-10-20
status: absorbed
absorbed_by: "design/adrs/0213-a-checked-creation-of-a-work-item-refuses-only-a-new-finding-inside-the-new.md#decision"
asked_again: []
cited: []
---

# D-0326 message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

I-0169 was closed on 2026-10-20.

The cause:

```text
message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
```

The fix:

```text
fixed by S-0451: the refusal came from itemnew.Create, behind flai task new and item_new, which refused the advisory wip.overlap a new task brings, while flai edit and flai touches keep it as a warning (ADR-0019). It now keeps it as a warning too, so two stories in progress that agree a split need no share, which lifts only a ready story's hold (ADR-0134)
```

Source: I-0169.
