---
id: MS-0110
title: "S-0378 holds S-0357 on overlap alone: may S-0357 start beside it?"
from: S-0357
to: S-0378
about: [docs/operators/settings.md]
status: open
participants: [flai]
created: 2026-10-13T18:51:22Z
updated: 2026-10-13T18:51:22Z
---

# MS-0110 S-0378 holds S-0357 on overlap alone: may S-0357 start beside it?

Between S-0357 and S-0378, about `docs/operators/settings.md`.

## Entries

### 2026-10-13T18:51:22Z flai S-0357
S-0378 holds S-0357 on overlap alone: may S-0357 start beside it?

S-0357, "A story's agent run through a provider is priced from the gateway's spend log, and an item's usage records priced_by, summed up the hierarchy as its least certain source", is ready, and its claim overlaps S-0378's on `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Over a gateway, Claude Code's `total_cost_usd` is its own estimate at Anthropic's prices, not what the gateway charged. [ADR-0132](../../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md) takes the cost of a run through a provider from the gateway's spend log, keeps the harness's figure when there is no provider, and estimates otherwise, naming the source on the item's `usage` as `priced_by`. This story builds the spend readers, OpenRouter's from `/api/v1/generation` and `/api/v1/key` and LiteLLM's from `/spend/logs`, the pricing when `flai serve` measures a run, and `priced_by` in the front matter. It joins a run's calls to the gateway as S-0351 and S-0356 found they can be joined. Showing `priced_by` in `flai stats` and the dashboard is the next story.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0378 --remove <path>`, or on the task that names the path, if S-0378 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0110 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0357 can start.
3. Reply with `message_reply` saying why the hold stands: S-0357 waits until S-0378 moves to review, is cancelled, or is sent back.
