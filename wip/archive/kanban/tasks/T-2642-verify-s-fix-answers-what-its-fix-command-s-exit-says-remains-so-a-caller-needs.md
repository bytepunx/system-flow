---
id: T-2642
type: task
nature: improvement
title: verify's fix answers what its fix command's exit says remains, so a caller needs no check after a --fix
status: done
parent: S-0609
owner: alex
created: 2026-11-07T23:48:45Z
updated: 2026-11-08T02:26:52Z
transitions:
  - to: ready
    at: 2026-11-08T02:17:14Z
    by: agent-S-0609
  - to: in-progress
    at: 2026-11-08T02:17:16Z
    by: agent-S-0609
  - to: done
    at: 2026-11-08T02:26:52Z
    by: agent-S-0609
stream: S-0609
tags: [cli]
touches: [flai/internal/verify/fix.go, flai/internal/verify/fix_test.go]
usage:
  source: log
  seconds: 576
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 8686
      cache_read: 3657804
      cache_write: 55082
      cost: 1.2596
---
# T-2642 verify's fix answers what its fix command's exit says remains, so a caller needs no check after a --fix

## Work

- Add to `flai/internal/verify/fix.go` a fix that also answers what remains: for each tier it fixes, the fix command's exit, and on a non-zero exit its findings parsed in the tier's `format`, as `Run` parses a tier's command. Where a tier's fix exits 0, nothing remains for it.
- Give it a checkout form beside `CheckoutFix`, reading the tiers once. Keep `Fix` and `CheckoutFix` as they are for their other callers (`flai/internal/taskdone/taskdone.go`, `flai/internal/mergequeue/preflight.go`, `flai/cmd/branch.go`).
- A fix command that cannot start, or a tier whose fix output cannot be parsed, answers that its result is unknown, so the caller falls back to `CheckoutUnfixed` and the ADR-0140 refusal is never weakened.
- Test it in `flai/internal/verify/fix_test.go` with a fake `Proc`: exit 0, exit 1 with findings, and a command that cannot start.
- Waits for no task: it changes no file another task changes, and the script's change in the first task does not change the markdown tier's exit codes.

## Done when

- The new fix answers, for each tier it ran, its changed paths and what remains, from the fix command's own run, with no second command run.
- An unknown result is distinct from a pass.
- `flai test flai/internal/verify/fix.go flai/internal/verify/fix_test.go` passes.

## Notes
