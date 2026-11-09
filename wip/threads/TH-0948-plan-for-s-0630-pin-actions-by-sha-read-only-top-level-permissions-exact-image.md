---
id: TH-0948
title: "Plan for S-0630: pin actions by SHA, read-only top-level permissions, exact image tags, one lockfile, and a check that reproduces I-0297; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0630-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md
  item: S-0630
status: open
participants: [planner-S-0630, orchestrator]
created: 2026-11-09T21:59:42Z
updated: 2026-11-09T22:08:07Z
---

# TH-0948 Plan for S-0630: pin actions by SHA, read-only top-level permissions, exact image tags, one lockfile, and a check that reproduces I-0297; cost of delay inputs needed

On wip/kanban/stories/S-0630-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md.

## Entries

### 2026-11-09T21:59:42Z planner-S-0630
**Question: cost of delay inputs.** S-0630 has no cost of delay inputs and no epic, so `flai cod S-0630` cannot work out a value.

**Recommended:** `penalty_per_week: 10` (USD).

Why: I-0297 is a low-severity impression with no observed instance and no time lost. The real exposure is a moving action tag or base image running in the release workflows that hold the cosign key and `contents`/`packages`/`id-token: write`. That is a supply-chain risk, and it is smaller than S-0627's script injection (`penalty_per_week: 100`) and S-0629's refused writes (`20`).

Alternatives:

- `penalty_per_week: 25`, if a compromise of the release path weighs more for you.
- `time_lost_per_cycle: 0`, which gives a value of 0 and leaves it last in a `cod` order.

I will write the value once the inputs are set.

**Plan.** S-0630's `### Planning` notes give the proposed solution, where each touch came from, and the figures. The story's goal asked for the solution to be proposed before it is built, so the operator can change it here.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2774 SHA pins for every action; top-level `permissions: contents: read` everywhere; the release jobs name their own writes; GoReleaser `v2.18.1`; the template workflow gets the `permissions:` block | six `.github/workflows/*.yml` files, `template/root/.github/workflows/system-flow-check.yml` | none |
| 1 | T-2778 exact image tags, plus a digest for `flaiover/Dockerfile`; `otel-lgtm` and `litellm` by release tag; dependabot's `docker` entry covers the `flai/images` folders | `flaiover/Dockerfile`, `flai/images/bedrock/Dockerfile`, `flai/images/system-flow/Dockerfile`, `flaiover/compose.yaml`, `scripts/litellm.sh`, `.github/dependabot.yml` | none |
| 1 | T-2779 remove `flaiover/package-lock.json` and the root `package-lock.json` | those two files | none |
| 2 | T-2781 `scripts/ci-pins.sh`, the test: it fails at the story's base on each instance and passes on the tree; `scripts/smoke.sh` runs it | `scripts/ci-pins.sh`, `scripts/smoke.sh`, `scripts/README.md` | T-2774, T-2778, T-2779 |
| 3 | T-2783 `design/tech/{ci,docker,node,observability}.md` and `design/system/release-signing.md` match; `flai issue close I-0297` | those five files, the I-0297 file, `design/issues/summary.md` | T-2774, T-2778, T-2779, T-2781 |

The three layer-1 tasks share no path, so they run together.

Figures: the forecast is 1h, raised from `flai forecast`'s 43m for the SHA and digest lookups and the integration and smoke tiers that the template and `flai/images` paths select. The cost of delay value waits on the question above.

**Assumptions.**

1. Actions are pinned by commit SHA, because tags can be moved and the release jobs hold the signing key. Images are pinned by exact tag, with a digest only where a signed release builds from the file. The agent images are built by no workflow, so dependabot moving exact tags is enough there.
2. The template's workflow gets the `permissions:` block but keeps its version tags. A project made from the template ships no `dependabot.yml`, so SHA pins there would go stale.
3. The caret ranges in `flaiover/package.json` stay. The frozen pnpm lockfile is the pin.
4. The test is a script that the smoke tier and CI run, not a new tier in `system-flow.yaml`. A story that changes only workflows therefore meets it in CI rather than in its own close-out. If you want it as a tier selected by `.github/**` and the Dockerfiles, say so, and T-2781 will add the tier to `system-flow.yaml`.
5. Out of scope: signing with `--tlog-upload=false` and `--insecure-ignore-tlog`, which ADR-0070 decided.

### 2026-11-09T22:08:07Z orchestrator (recommendation)
**Recommended:** `penalty_per_week: 10` USD on S-0630, as the planner recommends. Keep the plan as written.

**The cost of delay input.** I-0297 has no observed instance and no time lost. The exposure is a moved tag or image in the release workflows that hold the signing key. That is real but smaller than S-0627's script injection (100) and S-0629's refused writes (20), so 10 fits the planner's scale.

- No revenue rides on it, so a penalty fits.
- The `throughput` policy orders ready by forecast. The value decides no order today.
- Alternative: `penalty_per_week: 25`, if you weigh a compromise of the release path higher.

**Why I recommend rather than set it.** S-0630 was in the backlog when the planner asked, but you have since finalized it and moved it to ready. `plan_backlog_stories` lets me set only a backlog story's first inputs. Once you confirm, or the recommendation is confirmed after `confirm_after`, I set the input under `set_cost_of_delay`. The planner then works out the value.

**The plan fits S-0630.**

- Layer 1's three tasks (T-2774, T-2778, T-2779) share no path and run together. T-2781's script reproduces each instance for criterion 1. T-2783 closes I-0297 for criterion 2.
- The touches name files only. The template workflow keeps version tags, which is right: a project made from it has no dependabot to move a SHA.
- Assumption 4: keep `scripts/ci-pins.sh` in the smoke tier and CI rather than add a tier. Smoke already runs at the close-out of a story that changes the template or `flai/images`, and a new tier is a manifest change the story does not need.
- 1h over flai's 43m gives its reason: the SHA and digest lookups.

Source: design/conventions/strategic-agents.md § As the orchestrator
