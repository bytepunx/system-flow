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
- Every commit message names the story: subject line `<type>: [S-nnn] what changed`, or the epic for cross-story work. Types: `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `build`, `ci`, `chore`. Body says what and why, not how; a reader should not need the diff to understand it.
- Commit the documentation and work item updates with the code they describe. A commit that changes behavior and leaves `docs/` or `design/` stale is incomplete.
- Never force push. Never rewrite history that has been pushed. Never amend a commit that is not yours from this session.
- The only secrets, credentials, tokens, or private keys allowed in the repository are ones that work only against a local, disposable environment (for example compose defaults). If a real key lands in git, say so immediately; rotating it is the operator's decision.
- Do not commit generated output, build artifacts, caches, or dependency folders. If `git status` shows one, fix `.gitignore` in the same change.
- Verify the ignore rules do not hide files that must be committed. A fresh clone must build and pass `flai check --strict`.
- Branch names are `<type>-<story-id>-<slug>` for story work. Work on the default branch only when the operator says so.
- Pull requests follow `.github/pull_request_template.md`: the story ID and the definition of done, honestly ticked.
- Tags are the release mechanism and follow semver by delivery type: an epic completing is a major release; a `feature` story completing is a minor release; a `remediation` or `improvement` story, or a documentation-only change, is a patch release. `research` and `experiment` stories do not release. Every release gets a changelog entry, written by the release tooling where it exists and by hand in `CHANGELOG.md` otherwise. Creating a tag is fine; pushing it needs the operator's confirmation.
- Create remote repositories with the `gh` CLI when it is installed and authenticated. Ask the operator for the organization and whether the repository is public or private before creating it; never create a remote unasked. Without `gh`, give the operator the exact commands instead.
- Before committing, run `git status` and `git diff --stat` and read them. Unrelated changes are split out or explained.
- Attribution lines the operator or tooling requires go at the end of every commit message and pull request body, unchanged.

## When in doubt

- If a commit would need "and also" in its subject, it is two commits.
- If you are unsure whether something is safe to commit, it is not; ask.

<!-- system-flow:end-of-baseline -->

## Project additions
