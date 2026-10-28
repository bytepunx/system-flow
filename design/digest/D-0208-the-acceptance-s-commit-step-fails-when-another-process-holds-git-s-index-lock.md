---
id: D-0208
title: The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
kind: cause
sources: [I-0100, S-0307]
topics: [cli, git, flai, go, code, all]
paths: [flai/internal/storygit/indexlock.go, flai/internal/storygit/indexlock_test.go, flai/internal/preview/accept.go, flai/cmd/accept.go, flai/cmd/accept_lock_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0208 The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai

I-0100 was closed on 2026-10-07.

The cause:

```text
`flai accept` runs its steps in order: rebase and merge, move to done, archive, commit, then the overlap notices. The commit is a `git commit` in the main checkout, where flai serve's replans, planners' autocommits, and other acceptances also commit. When one of them holds `.git/index.lock` at that moment, the commit fails, and the acceptance exits with git's error. Everything before the commit has happened and stays: the branch is merged, the story is done, the items are moved to the archive and staged. Run again, the acceptance refuses ("already done"), and no other command finishes it. The operator, or someone watching, has to commit by hand.

The first lines of git's error are all the record there is, on the terminal of whoever ran the acceptance; flai serve's journal records the run as done because the move to done succeeded.
```

The fix:

```text
S-0307 fixed it: flai accept's git add and git commit run through storygit.RunPastIndexLock, which runs them again for about nine seconds while another process holds the index lock and never removes it; a story done and archived whose acceptance commit failed is finished by running flai accept again, which commits it under the usual subject and sends the overlap notices from the story's commits; and a commit that fails even so keeps git's whole output and names flai accept <id>.
```

S-0307's narrative gives the schedule: waits of 100 ms, 200 ms, 400 ms, 800 ms, 1.6 s, then 2 s three times, matching git's English message. D-0398 records how S-0307 read the rule against retries.

Sources: I-0100, S-0307.
