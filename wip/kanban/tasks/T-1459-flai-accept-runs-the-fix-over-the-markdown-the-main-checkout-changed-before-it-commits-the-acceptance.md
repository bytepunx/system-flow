---
id: T-1459
type: task
nature: improvement
title: flai accept runs the fix over the markdown the main checkout changed, before it commits the acceptance
status: done
parent: S-0367
owner: alex
created: 2026-10-12T16:53:10Z
updated: 2026-10-12T17:19:37Z
transitions:
  - to: ready
    at: 2026-10-12T17:09:55Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T17:09:55Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:19:37Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/cmd/accept_markdown.go, flai/cmd/branch.go]
after: [T-1470]
usage:
  source: log
  seconds: 582
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 17013
      cache_read: 2688621
      cache_write: 74414
      cost: 1.34
---
# T-1459 flai accept runs the fix over the markdown the main checkout changed, before it commits the acceptance

## Work

- Before `commitAcceptance`'s `git add -A` in the main checkout, run `verify.CheckoutFix` with the main checkout's tiers over the paths it has not committed, the wip items, narratives, threads, and messages flai wrote there (criterion 2, confirmed on TH-0413). A fix that fails is logged and the acceptance stands.
- The operator decided on TH-0413 that markdown the fix cannot fix blocks the acceptance, with its files, lines, and errors:
  - At preflight, before anything changes: fix the main checkout's uncommitted files, then `verify.CheckoutUnfixed` over them; when it fails, refuse with its findings.
  - After the story branch's rebase and its fix (`syncStoryBranch`), before the merge: `verify.CheckoutUnfixed` in the worktree over the files the branch changes against the main branch; when it fails, refuse with its findings and leave the story in review.
- Waits for T-1470, which makes `verify.CheckoutUnfixed`.

## Done when

- [ ] An acceptance whose wip markdown a fix command rewrites commits it fixed, under a test.
- [ ] An acceptance whose branch or main checkout holds markdown a tier with a fix command still fails on is refused, naming the files, lines, and messages, and changes nothing on the main branch, under a test.

## Notes
