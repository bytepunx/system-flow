---
id: ADR-0174
title: "While the merge queue can batch, flai verify and the close-out leave the test tiers the manifest marks queue to the batch, which always runs them, and flai verify --all runs them"
status: accepted
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0110, ADR-0156]
---

# ADR-0174 While the merge queue can batch, flai verify and the close-out leave the test tiers the manifest marks queue to the batch, which always runs them, and flai verify --all runs them

## Context

A story's close-out runs `flai verify`, and `flai verify` runs every test tier the branch's changes select ([ADR-0110](0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md)). In this repository that includes the `integration` and `smoke` tiers, about twelve minutes together. When main gains commits while the close-out runs, the close-out syncs and verifies again, and main's new paths send those tiers back ([ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md)). [I-0148](../issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md) records such a close-out: two rounds, about 24 minutes, past the agent's 10-minute tool limit.

The merge queue ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)) verifies the stories it batches once more, rebased on main, with the tiers their changes select. While it runs, a story's own run of the most expensive tiers is paid twice. The operator chose, on TH-0465 question 5, to leave them to the batch.

## Decision

While the merge queue can batch, `flai verify` and the close-out leave the test tiers the manifest marks `queue` to the batch, which always runs them, and `flai verify --all` runs them.

- **The key.** A tier in `system-flow.yaml`'s `tests` may carry `queue: true`. It is a boolean and may go on any tier. This project marks `integration` and `smoke`.
- **When the queue can batch.** For a project, when all of these hold, the conditions under which `flai serve` builds a batch:
  - the `merge-queue` host action is on for the project, in the configuration the running flai reads;
  - `orchestration.permissions.accept_reviews` is on;
  - acceptance is local.

  A hold does not count: a held queue resumes and batches later. With any condition off, no batch would run the tiers, so `flai verify` runs them as before.
- **What `flai verify` does.** While the queue can batch, it selects the tiers as before, then leaves out each selected tier marked `queue`. Each tier it leaves out is a step of its own, in the state `queued`, named as left to the merge queue. The run passes without them. A tier that a `queue` tier covers ([S-0342](../../wip/archive/kanban/stories/)) is not skipped for it, since the covering tier does not run. The MCP `verify` tool and the close-out, which runs `flai verify`, do the same.
- **`--all`.** `flai verify --all`, and the MCP `verify` tool's `all`, leave no tier to the queue. `flai test --all` runs every tier, as before.
- **The batch.** The queue's verification of a batch never leaves a tier out: it runs every tier its paths select, those marked `queue` included. So no story the queue lands lands without them.

This refines ADR-0110: the close-out a story's agent runs before review no longer runs every tier while the queue can batch. It refines ADR-0156: across a clean rebase, `flai verify` reuses and sends back only the tiers it runs, and leaves the `queue` tiers to the batch in every round.

## Consequences

- A close-out that main outruns no longer runs `integration` and `smoke` twice, and fits the agent's tool limit. I-0148 closes once this runs with the queue on.
- A story reaches review without the `queue` tiers having run on it. Eligibility for the queue is unchanged: its last `flai verify` passed at its head.
- `flai accept` by hand is unchanged (ADR-0173) and runs no test tier. A story accepted by hand while the queue can batch lands without the `queue` tiers, unless someone runs `flai verify --all` on it first. Enqueue it instead.
- The decision is read from the configuration the running flai reads. A flai run with another configuration, such as this repository's `scripts/flai.sh` with `.flai-cache/config.json`, sees the action only when that configuration turns it on.
- A report's fingerprint changes once, since the tiers carry the new key, so the first run after the upgrade reuses nothing.

## Alternatives considered

- **Leave the close-out as it is.** Rejected: the queue would save only the orchestrator's verify, and I-0148 would stay open.
- **Leave the tiers out whenever the host action is on.** Rejected: with `accept_reviews` off or under pull-request acceptance, no batch runs them.
- **Leave the tiers out only for a story already enqueued.** Rejected: a story is enqueued after review, once its close-out has run.
- **Show a left-out tier as `skipped`.** Rejected: `skipped` means another selected tier covers it, and the report would hide that nothing ran it.
