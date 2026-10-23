---
id: TH-0600
title: "ADR-0186 is still proposed: accept it, and I publish S-0477's guard change"
anchor:
  path: design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md
status: open
participants: [orchestrator]
created: 2026-10-22T21:44:33Z
updated: 2026-10-22T22:37:44Z
---

# TH-0600 ADR-0186 is still proposed: accept it, and I publish S-0477's guard change

On design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md.

## Entries

### 2026-10-22T21:44:33Z orchestrator
**Recommended:** accept ADR-0186 with `flai adr accept 0186`, or say "publish" here. I then release S-0477 as flai 1.67.14.

## Why I am holding the release

S-0477 landed through the merge queue at 21:44Z. Its criteria are met, and flai verify passed. It changes what flai guard allows:

- `flai components list` becomes a read for every agent, the strategic agents included.
- `flai components serve` is allowed to a story's task sub-agent only. It is refused to read-only sub-agents and to the strategic agents, because it writes a Vite config and may install dependencies (ADR-0141).
- `flai components config` stays refused.

The approved plan made `serve` a plain read for every sub-agent. The story's agent tightened that and recorded it as ADR-0186, but left the ADR `proposed`. ADR-0183, 0184, and 0185 are all `accepted`. Under judgement I won't ship a guard permission change whose decision you have not accepted.

## If you would rather change it

Say what to change, and I'll leave S-0477 unreleased until a follow-up story adjusts the rule. Nothing else is pending release.

### 2026-10-22T22:37:44Z orchestrator
Update: S-0476 has landed too, with the component harness answers and ADR-0187, which is accepted. A release ships everything pending, so S-0476 now waits behind S-0477 on this answer. Accept ADR-0186 or say "publish", and both go out together.
