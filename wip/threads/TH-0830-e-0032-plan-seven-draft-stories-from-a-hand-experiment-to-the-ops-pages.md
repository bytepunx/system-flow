---
id: TH-0830
title: "E-0032 plan: seven draft stories, from a hand experiment to the Ops pages"
anchor:
  path: wip/kanban/epics/E-0032-integrate-kluster-for-demo-and-testing.md
  item: E-0032
status: open
participants: [planner-E-0032, orchestrator]
created: 2026-11-06T21:53:02Z
updated: 2026-11-06T21:55:51Z
---

# TH-0830 E-0032 plan: seven draft stories, from a hand experiment to the Ops pages

On wip/kanban/epics/E-0032-integrate-kluster-for-demo-and-testing.md.

## Entries

### 2026-11-06T21:53:02Z planner-E-0032
E-0032 had no stories. I drafted seven, each with tasks in layers, touches, a forecast, and a cost of delay value. All are drafts in the backlog.

## Stories and order

| # | Story | Nature | after, and what it reads | Forecast | CoD/wk |
|---|-------|--------|--------------------------|----------|--------|
| 1 | S-0575 two authstar builds in two story namespaces of one kluster, by hand | experiment | none | 3h | 490.20 |
| 2 | S-0576 ADR: shared cluster, namespace per story, `kluster:` block, image load, acceptance, `kluster` host action | research | S-0575: its results document | 1h | 163.40 |
| 3 | S-0577 `flai kluster up/down/status`, the manifest block and its check, the `kluster` host action | feature | S-0576: the ADR's block, variables, and action name | 47m | 128.00 |
| 4 | S-0578 `flai kluster deploy/teardown/promote/deployments` | feature | S-0577: `flai/internal/kluster` and `flai/cmd/kluster.go` | 27m | 73.53 |
| 5 | S-0579 flai serve deploys a story at review, checks it, tears it down, promotes main on acceptance | feature | S-0578: deploy, teardown, promote, the record, and the lock | 1h | 163.40 |
| 6 | S-0580 Ops group in the rail, Klusters page | feature | S-0577: the package and the host action | 43m | 117.10 |
| 7 | S-0581 Ops › Deployments page, with demo access | feature | S-0578: the deploy functions and the record; S-0580: the Ops group, `hostapi/kluster.go`, and the routes | 42m | 114.38 |

S-0575 and S-0576 run one after the other. S-0577 follows them. After S-0577, S-0578 and S-0580 can run together. S-0579 and S-0581 follow S-0578.

The values are `flai cod`'s shares of your 1250 USD/week revenue input, split by forecast duration. S-0575 has the largest share only because it has the longest forecast.

## Tasks and layers

- S-0575: T-2466 → T-2467 → T-2468. They are one per layer, since all three write the one results document.
- S-0576: T-2469 (the ADR, accepted by you on a thread) → T-2470 (`design/system/kluster.md`) and T-2471 (tooling convention) together.
- S-0577: T-2472 (manifest and check) and T-2473 (host action) → T-2474 (package) → T-2475 (commands) → T-2476 (documents).
- S-0578: T-2477 (deploy) → T-2478 (teardown, promote, list) → T-2479 (commands) → T-2480 (documents).
- S-0579: T-2481 (deploy pass) → T-2482 (teardown and promote) and T-2483 (checks get the namespace) together → T-2484 (documents).
- S-0580: T-2485 (host API) and T-2486 (rail) → T-2487 (page and routes) → T-2488 (documents).
- S-0581: T-2489 (host API) and T-2490 (rail) → T-2491 (page and routes) → T-2492 (documents).

## Assumptions

- **The toggle is a host action.** The Features page lists host actions, and there is no manifest feature key. So the toggle is a `kluster` host action, off by default, which the page lists with no dashboard change.
- **flai only calls other tools.** It runs `kluster`, `docker` (or `k3d`) to load images, and the build, deploy, and teardown commands the project names in its `kluster:` block. flai needs no Helm or kubectl knowledge beyond deleting a namespace. S-0576's ADR settles the details.
- **kluster has gaps.** It has no image import, no deploy into a namespace, and no JSON status. flai works around them. The ADR lists the changes kluster could make as proposals for its own repository, not work items here.
- **The convention is reworded.** `design/conventions/tooling.md`'s project addition says kluster is not used here. S-0576 rewords it: system-flow itself runs no cluster, and flai drives a project's kluster.
- **The rail placement.** Ops sits between Design and Documentation, with Klusters, then Deployments.
- **The experiment runs outside this repository.** S-0575 runs kluster and authstar's scripts on the host, against a local cluster. It commits nothing in those repositories.

## Proposals

- **Not drafted, an addition elsewhere:** signet and authstar adopting this needs work in their own repositories. Each needs a `kluster:` block, build, deploy, and teardown scripts that take the namespace, a vhost per namespace, and a tenant per namespace. I recommend an epic in each once S-0576's ADR is accepted.
- **A merge you may prefer:** S-0575 and S-0576 could be one research story, if you want fewer hand-offs. I kept them apart so the ADR rests on measured results.
- **A split, if S-0577 runs long:** its host action task (T-2473) could be a story of its own. I kept it in S-0577 because it is two files.

Nothing here needs an answer before you finalize the drafts. Change any assumption above and I will plan it again.

### 2026-11-06T21:55:51Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan as drafted. Keep S-0575 and S-0576 apart, keep T-2473 in S-0577, and open the signet and authstar epics only once S-0576's ADR is accepted.

- **The plan fits.** Seven stories, ordered by `after` only where one reads another's output on main. S-0578 and S-0580 can run together after S-0577. Each story's tasks form layers that share no path. You approved E-0029's plan the same way (D-0384).
- **Keep S-0575 and S-0576 apart: yes.** The ADR should rest on measured results: namespace settings, load and deploy times, and kluster's gaps. A merged story would decide before measuring.
- **Split T-2473 out of S-0577: not now.** It is two files. Split it only if S-0577 runs long.
- **Epics in signet and authstar: yes, later.** They belong in those repositories, after S-0576's ADR fixes the `kluster:` block they fill in.

**What I did.** I finalized all seven drafts under `finalize_drafts`: S-0576 to S-0581 at 21:53Z, and S-0575 now. `flai promote --drafts` listed each complete and consistent. S-0575's host inputs exist: authstar's `integration/kluster.yaml` (profile `authstar`), `scenario-golden-path.sh`, and kluster v1.3.1. S-0576 to S-0581 wait behind S-0575 by `after`.

**Why it is yours.** The merge, the split, and the other repositories' epics are changes of scope.

**If you choose the merge.** Answer here before S-0575 starts, and planner-E-0032 can plan it again.

Source: design/digest/D-0384-plan-for-e-0029-six-draft-stories-their-order-their-tasks-and-the-assumptions.md
