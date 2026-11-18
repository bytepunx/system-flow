---
id: T-3439
type: task
nature: improvement
title: This repository's CI jobs call flai test --at ci where one suite stands for what the job ran, and the CI design says which keep their own commands and why
status: done
parent: S-0812
owner: alex
created: 2026-11-16T23:47:33Z
updated: 2026-11-18T09:03:59Z
transitions:
  - to: ready
    at: 2026-11-18T08:57:02Z
    by: agent-S-0812
  - to: in-progress
    at: 2026-11-18T08:57:03Z
    by: agent-S-0812
  - to: done
    at: 2026-11-18T09:03:59Z
    by: agent-S-0812
stream: S-0812
tags: [flai]
touches: [".github/workflows/system-flow-check.yml", ".github/workflows/flai.yml", ".github/workflows/flaiover.yml", design/tech/ci.md, design/system/devex.md, design/conventions/code-quality.md]
after: [T-3437]
usage:
  source: log
  seconds: 416
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 10528
      cache_read: 2588908
      cache_write: 92328
      cost: 1.3034
---
# T-3439 This repository's CI jobs call flai test --at ci where one suite stands for what the job ran, and the CI design says which keep their own commands and why

## Work

ADR-0312 §6, criterion 3. This repository has no workflow file of its own, so the baseline's `ci: all` applies. Go through each job of `.github/workflows/system-flow-check.yml`, `flai.yml`, and `flaiover.yml`. Where one suite of this repository's manifest, as S-0805 labels them, runs what the job's steps run, replace the steps with `flai test --at ci --suite <suite>`. Build flai from the tree first, as `system-flow-check.yml` does. The smoke job is the likely one: `scripts/smoke.sh` is the `smoke` tier.

Keep a job's own commands where no suite runs the same thing, and say why. The likely cases:

- `flai.yml`'s test matrix runs `-race` on Ubuntu and runs on macOS. The `go-test` and `integration` tiers run on one core without `-race` (S-0722).
- The golangci-lint action builds with the workflow's Go.
- The cross-build matrix for linux, darwin, and windows is no tier.
- `flaiover.yml` runs `pnpm build` beside the `flaiover` tier's steps.

Keep every action pinned by SHA, as `scripts/ci-pins.sh` checks.

Write the note in `design/tech/ci.md`, naming each job, whether it calls `flai test --at ci`, and why. Bring the CI rows of `design/system/devex.md` and the CI sentence in the project additions of `design/conventions/code-quality.md` up to date.

It waits for T-3437, because the jobs call the command it adds.

## Done when

- Each job whose work one suite covers calls `flai test --at ci --suite <suite>`, and its old steps are gone.
- `design/tech/ci.md` names each job that keeps commands of its own and why.
- `design/system/devex.md` and `design/conventions/code-quality.md` match the workflows.
- `scripts/ci-pins.sh` passes.
- `flai test` on the changed paths passes.

## Notes
