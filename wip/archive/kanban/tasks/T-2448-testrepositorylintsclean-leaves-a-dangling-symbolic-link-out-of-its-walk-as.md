---
id: T-2448
type: task
nature: feature
title: TestRepositoryLintsClean leaves a dangling symbolic link out of its walk, as markdownlint-cli2 does, so the strategic agents' .claude/agents links fail it on no fresh checkout
status: done
parent: S-0572
owner: alex
created: 2026-11-03T00:00:01Z
updated: 2026-11-03T00:05:30Z
transitions:
  - to: ready
    at: 2026-11-03T00:00:35Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T00:00:36Z
    by: system-flow
  - to: done
    at: 2026-11-03T00:05:30Z
    by: system-flow
stream: S-0572
tags: []
touches: [flai/internal/mdlint/repo_test.go]
---
# T-2448 TestRepositoryLintsClean leaves a dangling symbolic link out of its walk, as markdownlint-cli2 does, so the strategic agents' .claude/agents links fail it on no fresh checkout

## Work

In `flai/internal/mdlint/repo_test.go`, `lintTree` walks the monorepo with `filepath.WalkDir` and reads every `.md` entry. Since S-0553 (ADR-0236), `.claude/agents/{analyzer,curator,orchestrator,planner}.md` are committed relative links into `.flai-cache/agents/`, which is gitignored and written only by `flai role definition --write`, so on a fresh checkout they dangle and `os.ReadFile` fails the test. Before reading an entry that is a symbolic link, `os.Stat` it; when the target is missing, skip it and log it, as the markdownlint-cli2 action in `system-flow-check.yml` skips it on the same tree. A link that resolves is read and linted as before. Waits for nothing: it changes a file no other task touches.

## Done when

- `go test -race -run TestRepositoryLintsClean ./internal/mdlint/` passes on this host, where the four links dangle
- a resolving link is still counted and linted: the test's file count and findings are unchanged for one

## Notes
