---
id: T-1895
type: task
nature: feature
title: The forge contract reports the approvals the base branch requires, and flai serve keeps them with the checks and reviews
status: done
parent: S-0395
owner: alex
created: 2026-10-17T17:36:24Z
updated: 2026-10-17T17:41:43Z
transitions:
  - to: ready
    at: 2026-10-17T17:36:39Z
    by: agent-S-0395
  - to: in-progress
    at: 2026-10-17T17:36:43Z
    by: agent-S-0395
  - to: done
    at: 2026-10-17T17:41:43Z
    by: agent-S-0395
stream: S-0395
tags: []
touches: [flai/internal/forge/forge.go, flai/internal/forge/github.go, flai/internal/forge/github_test.go, flai/internal/forge/state.go, flai/internal/forge/state_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forge_test.go]
usage:
  source: log
  seconds: 300
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 16089
      cache_read: 2831500
      cache_write: 87307
      cost: 1.41
---
# T-1895 The forge contract reports the approvals the base branch requires, and flai serve keeps them with the checks and reviews

## Work

ADR-0165 left the number of approvals the base branch requires out of the forge contract, for S-0395 to add if its gate needs them. The gate needs it for criterion 1's "too few approving reviews for the main branch's rule".

- Add `RequiredApprovals(ctx, base string) (int, error)` to `forge.Forge`: the most approving reviews the base branch's rulesets (`pull_request` rule, `required_approving_review_count`) or its classic protection (`required_pull_request_reviews.required_approving_review_count`) require, 0 when none. Either out of reach (403 or 404) does not fail the read, as `requiredChecks` does.
- Keep it in `forge.Status` as `required_approvals`, read by flai serve with the checks and reviews in `readPullRequest`.
- Add `Status.Decisions`: each reviewer's latest review that approves, requests changes, or was dismissed, as the forge counts a reviewer's decision, so that the gate and flai serve's send-back count reviews alike.

No dependency on another task: it is the first layer.

## Done when

- [ ] `forge.Forge` has `RequiredApprovals`, the GitHub adapter answers it from rulesets and classic protection, and a test covers both and an out-of-reach read.
- [ ] `forge.Status` keeps `required_approvals`, flai serve fills it at each read, and a test shows it kept.
- [ ] `Status.Decisions` gives each reviewer's latest decisive review, with a test.

## Notes

- Added by agent-S-0395 when it reviewed the planner's tasks: the kept state had no approvals count.
