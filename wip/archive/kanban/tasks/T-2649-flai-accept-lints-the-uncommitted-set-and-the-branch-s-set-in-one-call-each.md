---
id: T-2649
type: task
nature: improvement
title: flai accept lints the uncommitted set and the branch's set in one call each, with no check after a fix whose exit says what remains
status: done
parent: S-0609
owner: alex
created: 2026-11-07T23:49:12Z
updated: 2026-11-08T02:43:11Z
transitions:
  - to: ready
    at: 2026-11-08T02:31:08Z
    by: agent-S-0609
  - to: in-progress
    at: 2026-11-08T02:31:09Z
    by: agent-S-0609
  - to: done
    at: 2026-11-08T02:43:11Z
    by: agent-S-0609
stream: S-0609
tags: [cli]
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/branch.go, flai/cmd/accept_fix_test.go, flai/cmd/move_refused_test.go, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md]
after: [T-2642]
usage:
  source: log
  seconds: 722
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 87
      output: 33183
      cache_read: 13973037
      cache_write: 210418
      cost: 4.8118
---
# T-2649 flai accept lints the uncommitted set and the branch's set in one call each, with no check after a fix whose exit says what remains

## Work

- Today an acceptance runs the markdown tier five times, one after another: `refuseUnfixedUncommitted` fixes and then checks the uncommitted set (`flai/cmd/accept_markdown.go`), the shared-set commit and `commitAcceptance` fix it again (`flai/cmd/accept.go`), and the rebase fixes the branch's set and `refuseUnfixedBranch` checks it again (`flai/cmd/branch.go`).
- Lint the uncommitted set once: `refuseUnfixedUncommitted` takes what remains from T-2642's fix and checks again with `CheckoutUnfixed` only when that result is unknown. The later fixes before the shared-set commit and before `commitAcceptance` run only over paths whose content changed since that fix, and not at all when none did.
- Lint the branch's set once: the fix after the rebase answers what remains, and `refuseUnfixedBranch` refuses from it, checking again only when it is unknown.
- Keep the refusal's text, findings, and `--json` form as ADR-0140 and S-0378 have them; what the fix leaves still refuses the acceptance before anything changes or before the merge.
- Extend `flai/cmd/accept_fix_test.go` to count the tier's runs in an acceptance and to show the same refusal for unfixed markdown.
- Waits for T-2642, whose fix answers what remains.

## Done when

- An acceptance of a story with markdown in both sets runs the markdown tier once for the uncommitted set and once for the branch's set, as the test counts.
- Unfixable markdown in either set refuses the acceptance with the same text and findings as before.
- `flai test flai/cmd/accept_markdown.go flai/cmd/accept.go flai/cmd/branch.go flai/cmd/accept_fix_test.go` passes.

## Notes
