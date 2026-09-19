---
id: S-0041
type: story
nature: feature
title: Review and acceptance from the dashboard
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-19T04:44:12Z
transitions:
  - to: ready
    at: 2026-09-19T01:56:14Z
    by: alex
  - to: in-progress
    at: 2026-09-19T04:07:13Z
    by: system-flow
  - to: review
    at: 2026-09-19T04:23:54Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:44:12Z
    by: alex
tags: [dashboard, cli]
---

# S-0041 Review and acceptance from the dashboard

## Goal
Stories in review can be read, discussed, and accepted or sent back from the dashboard: the designer sees the diff of the story branch against main, the ticked acceptance criteria, the narrative, and open threads, then accepts (which runs `flai accept`) or moves the story back with a reason.

## Acceptance criteria
- [x] A review page per story in review: acceptance criteria, narrative summary, open threads, the branch diff against main (files and hunks), and the release plan `flai release --dry-run` would produce
- [x] Accept runs `flai accept <id> --by <designer>` with the token holder as `--by`, streams progress, and shows the resulting tags; send back runs `flai move <id> in-progress --reason`
- [x] Failures (dirty tree, rebase conflict, release error) are shown verbatim from flai and leave the story in review
- [x] Tests with a fake flai; docs/users and docs/operators updated

## Tasks
- T-0185 Living design for the review page, the branch diff, and streamed acceptance
- T-0186 flai stream diff, and progress events from flai accept
- T-0187 Dashboard API: the branch diff, and acceptance as the designer with streamed progress
- T-0188 Review page: criteria, narrative, threads, diff, release plan, accept, and send back
- T-0189 User and operator documentation for review and acceptance
- T-0190 Verify: all tiers, and a review, a send back, and an acceptance through a real container

## Notes
Depends on S-0037 for the branch diff and S-0036 for authentication.

Decided when pulled, 2026-09-19. The branch diff is a flai command, `flai stream diff`, because flai is where git is run. "The token holder" is the manifest's owner, the same designer as for threads and board moves. Progress is newline-delimited JSON on the accept request's own response, fed by an info event flai logs per completed step. The release plan on the page is the acceptance preview, which is what `flai release --dry-run` computes. The board's and the item page's existing confirmation stays; the review page is the fuller path to the same acceptance.

Found while verifying, and fixed here (I-0021): `flai accept` refused a story with an unticked criterion or an open task only after it had merged the story's branch and removed its worktree. The rules for done are now checked on a copy of the item before anything changes, and show as a blocker in the preview, so the review page disables Accept and says why.

Verification, 2026-09-19. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint); flaiover lint, svelte-check 0 errors, 23 files and 147 tests, and a production build passed. A local image built from the branch ran as a second dashboard against a scratch git project whose manifest owner is `dana`; the operator's dashboard was not touched. By API: the diff of a story branch (one added file with its hunk); the preview; a send back with a reason, recorded in the story's notes with `by: dana`, and the story moved to review again; an acceptance streamed as `merged`, `done`, `archived`, `committed`, then `done` with the result, leaving the commit on main authored by the identity passed into the container, the story archived with `by: dana`, the worktree and branch gone; the diff of an accepted story answering with flai's reason. With the image rebuilt after the fix: a story with an unticked criterion showed the rule as a blocker in the preview, and a forced accept was refused with main, the worktree, the branch, and the status unchanged. In a browser against that container: the review page at 1440 pixels with criteria, narrative, the diff with its hunk open, the blocker, and Accept disabled; after the criterion was ticked, Accept clicked in the page showed the four steps and "accepted ... not pushed"; at 390 pixels a story without a branch showed flai's reason and nothing overflowed. Not exercised for real: a rebase conflict and a release error; their messages reach the page through the same path as the one that was, and the fake-flai test sends a conflict message through it verbatim. No tags were produced, because the scratch project has no components; the component test covers showing them.
