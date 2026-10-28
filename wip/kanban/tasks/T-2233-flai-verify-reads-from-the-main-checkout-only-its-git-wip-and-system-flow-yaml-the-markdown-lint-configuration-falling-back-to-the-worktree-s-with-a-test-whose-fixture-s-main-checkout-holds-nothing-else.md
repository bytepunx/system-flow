---
id: T-2233
type: task
nature: remediation
title: flai verify reads from the main checkout only its .git, wip, and system-flow.yaml, the markdown lint configuration falling back to the worktree's, with a test whose fixture's main checkout holds nothing else
status: backlog
parent: S-0520
owner: alex
created: 2026-10-28T13:57:28Z
updated: 2026-10-28T13:57:28Z
transitions: []
stream: S-0520
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go, flai/cmd/verify_test.go]
after: [T-2232]
---
# T-2233 flai verify reads from the main checkout only its .git, wip, and system-flow.yaml, the markdown lint configuration falling back to the worktree's, with a test whose fixture's main checkout holds nothing else

## Work

Waits for T-2232, whose ADR settles the list of what the main checkout gives the container. In `flai/internal/check/check.go`, `markdown` loads the markdownlint configuration with `mdlint.Load(MainRoot)` and warns `markdown.config` when it cannot be read; change it to fall back to `mdlint.Load(c.repo.Root)` when the main checkout's cannot be read, so that inside a container, where the main checkout's root holds only what is mounted, the worktree's tracked configuration lints `wip/`, and warn only when neither reads. In the main checkout `Root` is `MainRoot`, so nothing changes on the host. Add the case to the markdown tests in `check_test.go`: a root with no configuration file and a worktree with one lints with the worktree's. In `flai/cmd/verify_test.go`, add the test that reproduces I-0218's second cause and guards the ADR's list: from `verifyFixture`, remove from the root everything but `.git`, `wip/`, and `system-flow.yaml` (no `.markdownlint.yaml`, no `design/`, no `docs/`, no `bin/`, no `scripts/`), keep the worktree whole, write the narrative's two sections in the root's `wip/agents/`, and run `flai verify S-004 --sync-only` and `flai verify S-004` with the fixture's tier in the worktree: both pass, the report lands under `<root>/.flai-cache/verify/S-004.json`, and `flai verify --last` in the worktree reads it. Keep the test beside S-0521's, which removes the root's manifest and expects the refusal; the two together state the rule: the main checkout is present when its manifest is, and present enough when `.git` and `wip/` are beside it.

## Done when

`go test ./internal/check/ ./cmd/` in `flai/` passes with the two new cases; the new verify test fails on main's code before this task (the markdown configuration read from the stripped root warns `markdown.config` and `--strict` fails) and passes after it; `gofmt`, `go vet`, and golangci-lint pass on the changed files.

## Notes

`flai verify` reads the lint configuration through `flai check`, so the fallback sits in `check.go`, not in `internal/verify`. If the stripped fixture finds another read from the root that the ADR's list does not name, say so on the plan thread and in the narrative: that read is either mounted in T-2234 and T-2235 or given a fallback here, and the ADR's list is corrected before T-2236 documents it.
