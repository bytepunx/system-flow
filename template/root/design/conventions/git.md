---
title: Git
updated: 2026-10-02
audience: agent
order: 70
status: active
topics: [all]
roles: [story]
---

# Git

How history is made in this repository.

## Rules

- Each story is worked on its branch, `story/S-nnnn`, in the worktree `flai stream open` creates under `.flai-cache/worktrees/`. Story commits land there; `wip/` is written in the main checkout and committed by `flai accept`.
- Use `flai stream sync` for the branch's git operations, never `git rebase` or `git merge` by hand: it rebases the story's branch onto `main`.
- Commit each task when it is done:
  - run `flai stream sync` first, resolve any conflicts it reports in the worktree, and run the tests for what the task changed
  - then commit the task's changes with its documentation and work item updates
- Before moving a story to `review`, run `flai stream sync` again, resolve any conflicts, and commit what is outstanding.
- Use conventional commit style. Commit messages name the story: subject line `<type>: [S-nnnn] what changed`, or the epic for cross-story work. Types: `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `build`, `ci`, `chore`. Body says what and why, not how; a reader should not need the diff to understand purpose.
- Commit messages end with a `Co-Authored-By:` trailer naming the model that authored the change, such as `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`, when an agent authored it.
- Avoid incomplete commits by including the documentation and work item updates with the associated code changes.
- Never force push. Never rewrite history that has been pushed. Never amend a commit that is not yours from this session.
- The only secrets, credentials, tokens, or private keys allowed in the repository are ones that work only against a local, disposable environment.
- If a real key lands in git, say so immediately; rotating it is the operator's decision.
- Do not commit generated output, build artifacts, caches, or dependency folders. If `git status` shows one, fix `.gitignore` in the same change.
- Verify the ignore rules do not hide files that must be committed. A fresh clone must build and pass `flai check --strict`.
- Never create a remote repository. Creating one, and choosing its organization and visibility, is the operator's.
- Pull requests follow `.github/pull_request_template.md`: the story ID and the definition of done, honestly ticked.
- Accepting an item merges it, moves it to done, archives it, and commits: nothing is tagged and nothing is pushed.
- Publishing is a deliberate step of its own, over everything accumulated since the last one, not tied to any single item: `git fetch`, then `flai release --pending` computes, tags, and pushes the branch and the tags, run by hand or from the board's Publish action, when the operator chooses to.
- A component's bump, once published, is the highest delivery type among everything accepted for it since its last release, not one release per item:
  - an epic contributes a major bump
  - a `feature` story a minor bump
  - a `remediation` or `improvement` story, a documentation-only change, or a non-breaking dependency update a patch
  - any other component a story touched incidentally with additive, non-breaking changes contributes a patch when it publishes
  - a breaking change is never incidental and needs its own story against that component
- A `research` story is accepted like any other, and its findings (ADRs, design documents, new stories) land on `main`, but it contributes nothing to any component's next release. Releasing research builds with semver pre-release suffixes (`1.3.0-rc.1`) is a pattern to adopt when needed and is not yet defined.
- An `experiment` story records its results in a document under `design/experiments`, and is accepted like any other, but contributes nothing to any component's next release.
- Every release gets a changelog entry, naming every item a batched publish bundles, written by the release tooling where it exists and by hand in `CHANGELOG.md` otherwise.
- Before committing, run `git status` and `git diff --stat` and read them. Unrelated changes are split out or explained.
- Attribution lines the operator or tooling requires go at the end of every commit message and pull request body, unchanged.

## When in doubt

- If a commit would need "and also" in its subject, it is two commits.
- If you are unsure whether something is safe to commit, it is not; ask.

<!-- system-flow:end-of-baseline -->

## Project additions
