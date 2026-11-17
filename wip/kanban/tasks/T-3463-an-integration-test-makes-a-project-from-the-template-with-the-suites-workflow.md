---
id: T-3463
type: task
nature: feature
title: An integration test makes a project from the template with the suites workflow file and checks the close-out, the acceptance from integration, and the publish
status: backlog
parent: S-0815
owner: alex
created: 2026-11-16T23:50:02Z
updated: 2026-11-16T23:50:21Z
transitions: []
stream: S-0815
tags: [cli, template]
touches: [flai/cmd/suites_e2e_test.go, flai/cmd/testdata/e2e-suites-board.md]
---
# T-3463 An integration test makes a project from the template with the suites workflow file and checks the close-out, the acceptance from integration, and the publish

## Work

Write `testdata/e2e-suites-board.md`: the template's `design/workflows/board.md` with the move to `review` given `verify: [static, behavior]`, a state `integration` between `review` and `done` whose `runs:` gives `suites: [integration, smoke]` and `kluster: [deploy, checks]`, transitions from `review` into it, from it to `done` with the `accept` effect and `verify: [static, behavior]`, and its `send-back`, and `verify_at.publish: [integration, smoke]`.

Write `suites_e2e_test.go`, skipped under `-short` so the integration tier runs it, as `TestRenderPrototypeTemplate` in `new_test.go` is. Render `../../template` with `runCLI` into a temporary project. Copy the fixture over the project's workflow file; never edit `template/root/design/workflows/board.md`, which must equal flai's baseline (D-0648). Point each tier's command at a stub script that records its suite in a file and fails when a marker file exists. Then check:

- `flai check --strict` passes on the project.
- A story's close-out, `flai verify` on its branch, runs the `static` and `behavior` tiers and no `integration` or `smoke` tier.
- A story in `integration` with a passing record is accepted from it while main gains a commit that only the `integration` or `smoke` tiers select.
- `flai release --pending` against a local bare remote runs the `integration` and `smoke` tiers at main's head before it tags, and tags nothing when one fails.

The serve job on entering `integration` is T-3466's; this task leaves it out.

## Done when

The test passes in the integration tier against the stories in the story's `after`, and fails when the close-out runs a slow suite, when the acceptance is held for a slow suite's path, or when the publish runs no slow suite.

## Notes
