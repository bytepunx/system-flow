---
id: T-2366
type: task
nature: remediation
title: The four whole-repository tests skip under -short, the close-out's go-test tier runs with -short, and the check and the markdown lint run once in a batch
status: in-progress
parent: S-0548
owner: alex
created: 2026-10-31T07:25:19Z
updated: 2026-10-31T08:15:59Z
transitions:
  - to: ready
    at: 2026-10-31T08:15:58Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T08:15:59Z
    by: agent-S-0548
stream: S-0548
tags: [cli]
touches: [flai/internal/mcpserver/prime_test.go, flai/internal/check/check_test.go, flai/internal/mdlint/repo_test.go, system-flow.yaml, scripts/smoke.sh, design/conventions/code-quality.md, design/conventions/tooling.md, docs/operators/index.md, flai/internal/verify/select_test.go]
after: [T-2361, T-2363]
---
# T-2366 The four whole-repository tests skip under -short, the close-out's go-test tier runs with -short, and the check and the markdown lint run once in a batch

## Work

`TestPrimeOnThisRepository` and `TestPrimeComesInPartsOnThisRepository` (`flai/internal/mcpserver/prime_test.go`, 131 s and 111 s under race) prime every open story of the monorepo; `TestMonorepoIsClean` (`flai/internal/check/check_test.go`, 54 s) and `TestRepositoryLintsClean` (`flai/internal/mdlint/repo_test.go`, 40 s) run `flai check` and the markdown lint over the whole repository. Make each skip under `testing.Short()`, saying it is the integration tier's own.

Give the `go-test` tier in `system-flow.yaml` `-short` again, beside the `-timeout` T-2363 gives it; `scripts/integration.sh` keeps running without `-short`. This reverses I-0201's choice that the close-out runs the changed packages without `-short`: say so in the narrative's `## Decisions`, citing I-0201.

Run `flai check --strict` and the markdown lint of the whole repository once in a batch: either drop `scripts/check.sh` and `scripts/lint-md.sh` from `scripts/smoke.sh`, leaving them to the two integration tests, or have those two tests skip in the integration tier and leave them to smoke (3 s and 21 s as scripts, against 54 s and 40 s under race). Choose by time and say why in `## Decisions`.

Update the testing tiers line in `design/conventions/code-quality.md` and the `flai test` line in `design/conventions/tooling.md`, both project additions below the marker, to say which tests are the integration tier's own and that the close-out runs `-short`; and the close-out's tiers paragraph in `docs/operators/index.md`.

Second layer: it waits for T-2361, which changes every `mcpserver` test file and `prime_test.go` among them, and for T-2363, which changes the same `go-test` line in `system-flow.yaml`.

## Done when

`go test -race -short ./internal/mcpserver ./internal/check ./internal/mdlint` skips the four tests and passes; `scripts/integration.sh` runs them; one batch's logs show the whole-repository check and lint once. The conventions and the operator docs say which tests are the integration tier's own. Story criteria 5 and 6.

## Notes

Drafted by the planner of S-0548 on 2026-10-31. It first drafted this task for the serve tests, which the story's agent's T-2362 covers, and rewrote it.
