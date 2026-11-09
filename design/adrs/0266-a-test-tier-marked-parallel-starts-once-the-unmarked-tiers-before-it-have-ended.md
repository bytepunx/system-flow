---
id: ADR-0266
title: "A test tier marked parallel starts once the unmarked tiers before it have ended, beside the marked ones, and the run's verdict stays a serial run's"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
---

# ADR-0266 A test tier marked parallel starts once the unmarked tiers before it have ended, beside the marked ones, and the run's verdict stays a serial run's

## Context

`flai verify` runs the tiers it selects one after another, cheapest first, and stops at the first that fails (`verify.Run`). So in this repository's close-out the `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both, though none reads what another writes. Passed close-outs in the week to 2026-11-07 took p50 247 s and p90 510 s; `integration` took a mean 238 s and `flaiover` 176 s (cause 20 in [server-performance.md](../system/server-performance.md#where-the-time-goes), S-0610). A project needs a way to say which tiers may run beside others, without losing what the serial order gives: the cheap tiers fail fast before the expensive ones start, and the verdict names the first tier that fails.

## Decision

A test tier in the manifest's `tests` may carry `parallel: true`.

1. **When a tier starts.** A tier marked `parallel` starts once every selected tier before it in the list that is not marked has ended, without waiting for the marked tiers before it. A tier not marked waits for every selected tier before it, as now. With no tier marked, the run is the serial run it was.
2. **What stays as a serial run gives it.** The tiers are listed and started in the list's order, cheapest first. A tier another selected tier `covers` is skipped, as before. The run stops at the first tier in list order that fails: no tier after it starts, a tier after it that is running is stopped and listed `not-reached`, and a tier after it that had already ended is listed `not-reached` too. A tier before it that is running finishes, so the tier `stopped_at` names is the first in list order that fails. The verdict, the tier it names, and the findings are those of a serial run over the same outcomes. Findings are capped in list order.
3. **What a tier is told.** Each tier is told the tiers it waited for that passed, so that what it may skip does not depend on which tiers happened to end before it started.
4. **Validation.** `manifest.TestErrors` refuses `parallel` on a `manual` tier, which runs alone, by its name. Any other tier may carry it, with any other key.
5. **An older flai ignores the key.** `manifest.Load` reads the manifest without refusing unknown keys, so a flai that does not know `parallel` runs every tier one after another, as before; its `flai manifest set tests=<JSON>`, which is strict, refuses a tier that carries it, as it refuses `queue` and `manual`.

## Consequences

- A project marks its expensive tiers that read nothing of each other, such as `integration` and `flaiover` here, and leaves its cheap per-path tiers unmarked, so the expensive ones still start only once those pass.
- Marked tiers share the host's cores. On a loaded host a tier can take longer beside another than alone; the saving is measured, not assumed.
- Tiers that run at once share the worktree. A project marks only tiers that do not write what another reads.
- `verify.Tiers`, the merge queue's batch, runs its tiers one at a time and may keep doing so; the key is a permission, never a requirement.

## Alternatives considered

- **`after: [<tier>]`, naming the tiers a tier waits for.** It says exactly what each tier waits for, and would let `smoke` wait for `template` and `markdown` only. It needs a cycle check, a meaning for an `after` that names a tier the run did not select, and a way to tell an empty list, wait for none, from no key, wait for all, which `omitempty` YAML does not keep. No tier here needs that precision.
- **Run every tier at once.** Simplest, but the expensive tiers would start before a gofmt or vet finding stops the run, and the host would run every tier for a change that a cheap one rejects.
