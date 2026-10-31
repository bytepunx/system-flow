---
id: D-0442
title: "S-0549 plan: seven tasks in five layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 20m` on S-0549, as the planner recommends"
kind: answer
sources: [TH-0775]
topics: [testing, git, cli, dashboard, orchestration, planning, flai, go, flaiover, sveltekit, code, all]
paths: [flai/internal/serve, flai/internal/mergequeue, flai/internal/check, design/system/flai-cli.md, design/adrs, docs/users/flai.md, design/adrs/README.md, flai/cmd/accept_markdown.go, flai/internal/verify/findings.go, flai/internal/verify/findings_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/operators/index.md]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0442 S-0549 plan: seven tasks in five layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 20m` on S-0549, as the planner recommends

The operator, alex, answered TH-0775 at 2026-10-31T09:50:06Z.

What was asked:

```text
**Question: cost of delay inputs.** Neither S-0549 nor E-0031 has inputs, so `flai cod S-0549` refuses and I wrote no value.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 20m` on S-0549, as the planner recommends:

```bash
flai edit S-0549 --time-lost-per-cycle 20m
```

## The figure, and not 12m or 1h

- **The goal is time lost.** About 40 minutes a fortnight of finding and fixing a dirty main by hand. The verification time wasted is the host's, not money.
- **The story gives no reason for an alternative.** It does not count the stories a failing batch holds back, which would argue for 1h.
- **On the story, not E-0031**, as on TH-0773 and TH-0774.

## Why this needs you

S-0549 is ready. You moved it there at 2026-10-31T07:24:09Z. flai refuses the orchestrator inputs on a story past the backlog: it refused S-0547's just now with "S-0547 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.** T-2382 and T-2383 cover criterion 1, T-2386 criteria 1 to 3, T-2384 criterion 3, T-2381, T-2385, and T-2387 criterion 4.
- **Assumptions: right.** Holding on error-level findings only, a hold apart from `flai queue hold`, landing again only while main is the verified base: each keeps the queue safe. T-2381's ADR records them.
- **Proposal, `flai queue` text output: leave it out.** No criterion asks for it, `flai queue --json` and `flai board --json` carry the hold, and S-0552 changes that file.
- **Touches and forecast fit.** The 5h forecast over flai's 16m states its reason.

The thread stays open for the cost of delay input only.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
````

Source: TH-0775.
