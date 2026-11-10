---
id: T-2987
type: task
nature: remediation
title: "Add template/** to the integration and smoke tiers' paths, with a manifest test that a template-only change runs both"
status: backlog
parent: S-0696
owner: alex
created: 2026-11-10T10:29:40Z
updated: 2026-11-10T10:29:40Z
transitions: []
stream: S-0696
tags: [flai, testing, manifest]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
---
# T-2987 Add template/** to the integration and smoke tiers' paths, with a manifest test that a template-only change runs both

## Work

- In `system-flow.yaml`, add `"template/**"` to the `paths` of the `integration` tier and of the `smoke` tier. Keep their `when` as it is.
- In `flai/internal/verify/select_test.go`, add a repository manifest test beside `TestRepositoryManifestRunsTheGatewayChecksOnlyByName`. It reads the manifest with `CheckoutTiers` and takes one template-only path, such as `template/root/design/conventions/code-quality.md`. It asserts that `SelectStory` selects `integration` and `smoke`. It also asserts that `SelectStoryLeavingQueue`, with that path and an empty claim, runs both rather than queueing them, and that each `RanFor` names `template/**`.
- Do not gate the test under `-short`, so that `go-test` catches a manifest that drops the pattern, as `TestRepositoryManifestTypeChecksChangedFlaioverFiles` does.
- First check that the test fails on the manifest as it is, then make the manifest change.
- This task waits for no other task: it is the fix the others describe and close.

## Done when

- The new test fails without the manifest change and passes with it.
- `flai test system-flow.yaml flai/internal/verify/select_test.go` passes.

## Notes
