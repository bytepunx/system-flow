---
id: T-2182
type: task
nature: improvement
title: A test in flai/internal/verify asserts this repository's manifest runs shellcheck and the script tests for a change under scripts/
status: done
parent: S-0503
owner: alex
created: 2026-10-26T05:23:40Z
updated: 2026-10-26T05:39:15Z
transitions:
  - to: ready
    at: 2026-10-26T05:37:45Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-26T05:37:46Z
    by: agent-S-0503
  - to: done
    at: 2026-10-26T05:39:15Z
    by: agent-S-0503
stream: S-0503
tags: [flai, tests]
touches: [flai/internal/verify/select_test.go]
after: [T-2181]
usage:
  source: log
  seconds: 89
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 7347
      cache_read: 1338604
      cache_write: 40047
      cost: 0.6788
---
# T-2182 A test in flai/internal/verify asserts this repository's manifest runs shellcheck and the script tests for a change under scripts/

## Work

- Add `TestRepositoryManifestRunsScriptTestsForScriptChanges` beside `TestRepositoryManifestRunsChangedGoPackagesWithoutShort` in `flai/internal/verify/select_test.go`, reading the real manifest with `CheckoutTiers("../../..")` and skipping only when the monorepo is not present, not under `-short`, so `make test` catches a manifest that opens the gap again. It reproduces I-0199's instance: `SelectStory` for the paths S-0481 changed, `scripts/flaiover-install.sh` and `scripts/flaiover-install-test.sh`, names the `shellcheck` tier with those files as `{files}` and the `flaiover` tier.
- Then, for each `scripts/*-test.sh` on disk, assert that `SelectStory` for that one path names a tier other than `shellcheck` whose command runs it: the tier's command script names the test script, directly or through the `scripts/*.sh` it names in turn, followed by their text. A script test added later that no tier both selects and runs fails the test and names the script. `scripts/install-published-test.sh` is the one exception, named in the test with its reason: it needs the network and a token, and CI's `install-published.yml` runs it.
- Keep the comment style of the two tests above it: the issue, what the manifest must do, and why the test is not gated.

Waits for T-2181: it reads the manifest that task changes, and fails against the manifest as it stands today.

## Done when

- `go test ./internal/verify -run TestRepositoryManifest -count=1` passes in `flai/`.
- With the `scripts/flaiover-*.sh` pattern taken out of the `flaiover` tier once, in the worktree and put back, the test fails and names `scripts/flaiover-install-test.sh`.
- gofmt, `go vet`, and golangci-lint are clean on the package.

## Notes

Drafted by planner-S-0503; the second bullet tightened by agent-S-0503: with `smoke` selecting `scripts/**`, every script test is selected by some tier, so selection alone would pass for a script no tier runs. The test is the story's "test that reproduces it where one fits": it reads the manifest in the checkout, as the tests for I-0105 and I-0201 in the same file do.
