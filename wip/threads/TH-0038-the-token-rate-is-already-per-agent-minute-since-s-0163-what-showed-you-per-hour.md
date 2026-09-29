---
id: TH-0038
title: "The token rate is already per agent minute since S-0163: what showed you per hour?"
anchor:
  path: wip/kanban/stories/S-0165-tokens-rate-chart-should-express-token-use-in-minutes.md
  item: S-0165
status: resolved
participants: [agent-S-0165, alex]
created: 2026-09-29T22:38:20Z
updated: 2026-09-29T22:40:43Z
---

# TH-0038 The token rate is already per agent minute since S-0163: what showed you per hour?

On wip/kanban/stories/S-0165-tokens-rate-chart-should-express-token-use-in-minutes.md.

## Entries

### 2026-09-29T22:38:20Z agent-S-0165
S-0163 (accepted 21:02 UTC) already changed the token rate to tokens per agent minute, per ADR-0053 point 4:

- `flaiover/src/lib/viz/charts.ts` `tokenRate` charts `tokens_per_minute`, axis "tokens per agent minute"; the table and summary line say "tokens per agent minute".
- The running dashboard container (image revision 5fe2577, flai 1.26.0) serves a bundle with "tokens per agent minute" and none with "tokens per agent hour".
- `flai stats` prints "per agent minute" and "tokens/min".
- `tokens_per_hour` is still in `flai stats --json` beside `tokens_per_minute`, for dashboards written to flai 1.25, as ADR-0053 decided; nothing charts it.

So I cannot find anything that shows tokens per agent hour. Two guesses:

1. A browser tab opened before the dashboard was restarted at 22:32 UTC, still running the older bundle (flaiover 0.30, S-0143's point per item in tokens per agent hour). A reload would show minutes.
2. "Tokens per hour", the title of the Tokens per day chart when the bucket is an hour: that is tokens spent in each clock hour, not per agent hour.

Which did you see, or where else? If it is (1), I would close this story with nothing to change (or cancel it; your call). If it is (2) or you want `tokens_per_hour` removed from `flai stats --json` (the flai/cmd claim suggests that), say so and I will do that, with an ADR since metrics.md changes only with one.

### 2026-09-29T22:40:24Z alex
good job, this is solved, the dashboard was still showing the old chart (not sure why) but a hard refresh did fix it.

### 2026-09-29T22:40:39Z agent-S-0165
Thanks. The likely why: a SvelteKit tab keeps the client code it loaded, and moving between pages inside it never fetches a new bundle, so a tab left open across the 22:32 restart went on drawing the S-0143 chart until a full reload. SvelteKit can poll for a new version (`kit.version.pollInterval`) and reload on the next navigation when there is one; flaiover sets neither. If you want that, it is its own story; I have not made one. S-0165 goes to review with nothing changed but its narrative and the ticked criterion.

### 2026-09-29T22:40:43Z agent-S-0165
Resolved: The chart was already per agent minute since S-0163; the operator's tab held the old bundle and a hard refresh showed minutes.
