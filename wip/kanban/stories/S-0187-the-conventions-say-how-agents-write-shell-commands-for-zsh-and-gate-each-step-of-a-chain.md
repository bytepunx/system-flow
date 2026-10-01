---
id: S-0187
type: story
nature: improvement
title: The conventions say how agents write shell commands for zsh and gate each step of a chain
status: backlog
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:00:32Z
transitions: []
tags: [template]
topics: [conventions]
touches: [design/conventions/, template/, scripts/]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0187 The conventions say how agents write shell commands for zsh and gate each step of a chain

## Goal

Two behavioural issues that no convention addresses:

- **I-0006** (seven occurrences): the agent shell is zsh, and bash idioms fail or misbehave silently: unquoted globs abort with "no matches found", `==` in `[`, word-splitting of unquoted variables. Scripts under `scripts/` are POSIX `sh` and are fine; ad hoc commands are not.
- **I-0012** (two occurrences): a close-out chain was not gated on each step's exit code (a pipe or `;` hid a failure), so a commit went out with an empty narrative. `item_move` now refuses a story with uncommitted changes, but nothing stops a chain from carrying on after a failed step.

## Acceptance criteria
- [ ] The tooling convention (this repository and `template/`, below the baseline marker where it is project-specific, in the baseline where it holds for every project) has a short shell rule: the host shell may be zsh, so quote globs, URLs, and variables, use `[ a = b ]`, and put any sequence of more than a few commands in a script under `scripts/`
- [ ] The same rule says each step of a chain that must succeed is joined with `&&` or run under `set -e`, never `;` or a pipe that hides its exit code, and a step that produces what a commit records is checked before the commit
- [ ] The close-out steps agents run before review (lint, tests, `flai check --strict`, narrative, commit) are a script under `scripts/` that stops at the first failure, and the work-management convention points to it
- [ ] I-0006 and I-0012 are closed with what fixed them

## Tasks

## Notes
