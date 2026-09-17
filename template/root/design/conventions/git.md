---
title: Git
updated: 2026-09-15
audience: agent
order: 70
status: active
---

# Git

How history is made in this repository.

## Rules

- Commit when a story lands: when it moves to `review`, and again after the operator accepts and archives it. Not at task boundaries. Use conventional commit style. Always report commits in the chat.
- Every commit message names the story: subject line `<type>: [S-nnnn] what changed`, or the epic for cross-story work. Types: `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `build`, `ci`, `chore`. Body says what and why, not how; a reader should not need the diff to understand it.
- Commit the documentation and work item updates with the code they describe. A commit that changes behavior and leaves `docs/` or `design/` stale is incomplete.
- Never force push. Never rewrite history that has been pushed. Never amend a commit that is not yours from this session.
- The only secrets, credentials, tokens, or private keys allowed in the repository are ones that work only against a local, disposable environment (for example compose defaults). If a real key lands in git, say so immediately; rotating it is the operator's decision.
- Do not commit generated output, build artifacts, caches, or dependency folders. If `git status` shows one, fix `.gitignore` in the same change.
- Verify the ignore rules do not hide files that must be committed. A fresh clone must build and pass `flai check --strict`.
- Branch names are `<type>-<story-id>-<slug>` for story work. Work on the default branch only when the operator says so.
- Pull requests follow `.github/pull_request_template.md`: the story ID and the definition of done, honestly ticked.
- Tag on acceptance. When the operator moves an item to `done`, tag and push a semver release in the same step, public or private repository alike. If it is good enough to merge into main, it gets a matching release; no further approval is asked.
- The bump follows delivery type for the component the item delivers to: an epic done is a major release; a `feature` story done is a minor release; a `remediation` or `improvement` story, a documentation-only change, or a non-breaking dependency update is a patch release. Any other component the story touched incidentally with additive, non-breaking changes gets a patch; a breaking change is never incidental and needs its own story against that component.
- `research` and `experiment` stories stay on a branch and get their additional testing there; they do not release from main. Releasing research builds with semver pre-release suffixes (`1.3.0-rc.1`) is a pattern to adopt when needed and is not yet defined.
- Every release gets a changelog entry, written by the release tooling where it exists and by hand in `CHANGELOG.md` otherwise. `flai release` computes the bump, tags, and pushes once it exists; until then the agent does it by hand at acceptance.
- Create remote repositories with the `gh` CLI when it is installed and authenticated. Ask the operator for the organization and whether the repository is public or private before creating it; never create a remote unasked. Without `gh`, give the operator the exact commands instead.
- Before committing, run `git status` and `git diff --stat` and read them. Unrelated changes are split out or explained.
- Attribution lines the operator or tooling requires go at the end of every commit message and pull request body, unchanged.

## When in doubt

- If a commit would need "and also" in its subject, it is two commits.
- If you are unsure whether something is safe to commit, it is not; ask.

<!-- system-flow:end-of-baseline -->

## Project additions
