---
id: ADR-0170
title: "Under pull-request acceptance, flai accept merges the story's pull request on the forge by the project's merge method, through flai serve when the token is not at hand, then brings main up to origin's and accepts as under local acceptance"
status: accepted
date: 2026-10-18
supersedes: []
superseded_by: []
refines: [ADR-0032, ADR-0067, ADR-0165, ADR-0166, ADR-0169]
topics: [release, orchestration]
---

# ADR-0170 Under pull-request acceptance, flai accept merges the story's pull request on the forge by the project's merge method, through flai serve when the token is not at hand, then brings main up to origin's and accepts as under local acceptance

## Context

[ADR-0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md) has acceptance merge the story's branch into main on the host, move the story to done, archive it, and commit. [ADR-0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) keeps accepted work off the remote until it is published.

[ADR-0165](0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md) lets a project choose pull-request acceptance and a merge method. [ADR-0166](0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md) has flai serve push each story in review and open its pull request, and keeps the forge token out of every agent's environment. [ADR-0169](0169-under-pull-request-acceptance-every-acceptance-the-operator-s-and-the.md) gates every acceptance on the pull request's checks and reviews. It left the merge local until this story.

A team whose main branch takes changes only through pull requests cannot accept that way. The merge has to happen on the forge. Four things stand in the way:

- The token is where flai serve runs, and not where agents run. The orchestrator accepts as an agent.
- No channel leads from the command line to flai serve.
- The local main branch has commits the remote's main does not: the acceptance commits of stories not yet published, and the creation and moves of work items.
- A reviewer may merge the pull request on the forge before anyone runs `flai accept`.

## Decision

Under pull-request acceptance, `flai accept` merges the story's pull request on the forge by the project's merge method, through flai serve when the token is not at hand. It then brings main up to `origin`'s and accepts as under local acceptance.

- **Order.** Under `acceptance.mode: pull-request`, once ADR-0169's gate passes, `flai accept S-nnnn` does these in order, in place of the local merge:
  1. Refuses, naming the commits, when local main has unpublished commits outside `wip/`.
  2. Merges the pull request on the forge with `acceptance.merge_method`, at the head the gate passed at, and deletes `story/S-nnnn` on `origin`.
  3. Fetches, and brings local main up to `origin`'s main.
  4. Removes the story's local branch and worktree.
  5. Moves the story to done, archives it, and commits, as under `local`.
- **The merge.** `MergeStory` in `flai/internal/serve/forgemerge.go` finds the story's pull request. When it is merged already, it answers so. Otherwise it merges it with `Merge` and the project's method, and answers with the merge commit or the forge's reason for refusing. `Merge` refuses when the head has moved since the gate read it.
- **The contract.** `forge.Forge` gains `PullRequest(number)`: the pull request by its number, in any state. `FindPullRequest` sees only open ones, so it cannot tell a merged request from a missing one.
- **The token at hand.** Where the token's variable is set, `flai accept` calls `MergeStory` in its own process.
- **The hand-off.** Where it is unset, as for every agent, `flai accept` writes a request to `.flai-cache/forge/merge/S-nnnn.request.json`. A merge pass in flai serve looks for requests every second. It renames a request to `S-nnnn.taken.json`, merges with its token, and writes `S-nnnn.answer.json`. `flai accept` waits 30 seconds for the take and 5 minutes for the answer, and refuses when either runs out. The hand-off needs no new listener and no new secret.
- **The host action.** The merge is gated on the `forge` host action. While it is off, the acceptance refuses, naming the host action, whether `flai accept` or flai serve would make the merge.
- **The dashboard.** The host API's `accept.run` runs `flai accept` in flai serve's environment, where the token is at hand, so it merges directly. A `flai accept` run by the dashboard (`FLAI_AGENT=flaiover`) never hands off. Without the token there, flai serve has none to merge with either.
- **The remote branch.** It is deleted with `git push origin --delete story/S-nnnn` and the host's git credentials, as ADR-0166 pushes it. A failed delete is reported, and the merge stands.
- **Local main.** After the fetch, flai checks that the merge commit is on `origin`'s main. It rebases local main's unpublished `wip/` commits onto `origin`'s main in a temporary detached worktree, then moves main there with `git reset --keep`. Uncommitted edits in the main checkout survive. `git rebase` in the main checkout would refuse them, and `git stash` is shared by every worktree.
- **The local branch.** It is removed with `git branch -D`. A squash or rebase merge leaves it unmerged in git's eyes.
- **Merged by someone else.** flai serve's kept status records the pull request's state. A story whose pull request is merged is not pushed again, and gets no new pull request. ADR-0169's gate passes over it. `flai accept` makes no second merge: it fetches, and goes on from step 3.
- **Refusals.** A merge the forge refuses, for a conflict, a branch out of date, or a rule, stops the acceptance before anything local changes, naming the forge's reason. Unpublished commits on local main outside `wip/` stop it before the merge. They are code no pull request carried, and the next publish would take them to `origin`'s main past the review.
- **The dry run.** It says that it would merge through the remote and by which method, or that the request is already merged.
- **The same acceptance.** The acceptance commit's subject and the archive are the same in both modes, so metrics and history read the same.
- **Publishing.** The story's code reaches `origin`'s main at the merge. The acceptance commit and the archive reach it only when published, as ADR-0067 has it. A main branch that refuses every direct push refuses that publish too. This decision leaves that open. TH-0464 raised it and recommends a later epic.

## Consequences

- A team whose main branch takes changes only through pull requests can accept with flai, from the shell, the dashboard, or the orchestrator.
- Under pull-request acceptance, the story's code reaches the remote before it is published. Its release, tag, and changelog still wait for publishing.
- ADR-0166's Consequences name S-0396 for the merge of the pull request. It is this story, S-0397.
- Agents and the orchestrator can accept only while flai serve runs with the token and the `forge` host action on.
- A step that fails after the forge's merge leaves the request merged. Running `flai accept` again finds it merged and goes on.
- `git reset --keep` refuses when an uncommitted edit lies in a file the merge changed, and main stays where it was.
- Every acceptance under pull-request acceptance rewrites local main's unpublished `wip/` commits, so their hashes change.

## Alternatives considered

- **Give the agent the token.** Rejected. ADR-0166 keeps it out of every agent's environment.
- **A socket or HTTP endpoint in flai serve for the hand-off.** Rejected. It needs a new listener and a secret to guard it. A file under `.flai-cache` needs neither.
- **Merge locally, then push main.** Rejected. A main branch that takes changes only through pull requests refuses the push, and it is what the mode exists to avoid.
- **Rebase in the main checkout, or stash its edits first.** Rejected. `git rebase` refuses uncommitted changes, and the stash is shared by every worktree and every session.
- **Carry unpublished commits outside `wip/` onto `origin`'s main.** Rejected. They would reach the remote without review.
- **Tell a merged request from a missing one by `FindPullRequest`.** Rejected. It sees only open requests.
