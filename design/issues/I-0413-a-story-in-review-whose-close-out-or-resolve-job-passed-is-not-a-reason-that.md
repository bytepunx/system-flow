---
id: I-0413
title: A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event
class: efficiency
status: open
count: 1
cost: 48h
first_reported: 2026-11-27T22:46:35Z
last_reported: 2026-11-27T22:46:35Z
updated: 2026-11-27T22:46:35Z
---

# I-0413 A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event

## Description
A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event

## Instances

### 2026-11-27T22:46:35Z
Report: design/analysis/2026-11-27-all.md.
S-0839 and S-0843 passed their close-outs on 2026-11-18, were refused by the dirty-main lint (I filed that apart), and once the lint was fixed on 2026-11-21 and S-0843's resolve job passed at 18:26Z, nothing started the orchestrator for them: a verified event and a resolve job's end are not among its triggers, its later runs that day acted on their own reasons and ended (ADR-0209), and its own threads on the stories awaited the operator. Both waited six more days for the operator's acceptance.

## Impact
Serve journal: S-0843's rebase-and-resolve job "passed in 12m46s; stayed in review" at 2026-11-21T18:26:54Z; S-0839 passed its close-out at 2026-11-18T10:01:34Z; after 2026-11-21T18:27Z no journal entry names either story until the operator accepted both at 2026-11-27T22:33Z. The orchestrator ran five more times on 2026-11-21 (18:33 to 18:58) for reasons naming S-0811, S-0812, and S-0834 only. design/workflows/roles/orchestrator.md trigger: asked, thread awaits-orchestrator, moved to review/backlog, edited draft-changed, planned, release-due, ready-empty, reordered, schedule; no verified event and no resolve-job end. flai stats waiting.longest: S-0839 822,674 s and S-0843 818,829 s in review, awaited alex. .flai-cache/events.jsonl holds 536 verified events.

- time_lost_per_cycle: 48h

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
