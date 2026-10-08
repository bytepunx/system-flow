---
id: T-1363
type: task
nature: remediation
title: TestRepositoryLintsClean scopes itself to the story in a close-out, so main's committed wip is a note, not a failure
status: in-progress
parent: S-0345
owner: alex
created: 2026-10-08T08:39:49Z
updated: 2026-10-08T09:09:24Z
transitions:
  - to: ready
    at: 2026-10-08T09:09:24Z
    by: agent-S-0345
  - to: in-progress
    at: 2026-10-08T09:09:24Z
    by: agent-S-0345
stream: S-0345
tags: [flai]
touches: [flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/repo_test.go]
---
# T-1363 TestRepositoryLintsClean scopes itself to the story in a close-out, so main's committed wip is a note, not a failure

## Work

`TestRepositoryLintsClean` in `flai/internal/mdlint/mdlint_test.go` lints every markdown file the branch holds, `wip/` included, and fails on any finding. A story branch holds `wip/` as main last committed it, so a line any agent commits there fails every story's integration tier (I-0117). `TestMonorepoIsClean` in `flai/internal/check/check_test.go` already scopes itself when `CLOSE_OUT_STORY` is set (S-0249, ADR-0085); do the same here.

- Move the test to a new `flai/internal/mdlint/repo_test.go` in the external package `mdlint_test`: `workitem` imports `mdlint`, so the internal test package cannot import `workitem` or `storygit` without a cycle.
- With `CLOSE_OUT_STORY` set, take the story's changed paths with `storygit.StoryChanges`, as `TestMonorepoIsClean` does. A finding on a file under `wip/` that the story does not change is logged with `t.Logf` as outside the story. Every other finding still fails, so the test keeps guarding mdlint's parity with markdownlint on the design, docs, and code trees.
- Unset, as in CI and `make`, every finding fails as before.
- Factor the walk and the split into a helper that takes the root and the set of changed paths. Add a test that reproduces I-0117 on a temporary tree: a bare `www.` line in `wip/agents/x.md` is a note when the story does not change that file and a failure when it does, and the same line in `design/x.md` fails either way.

Runs in the first layer: it waits for no task and shares no path with the `lint-md.sh` task.

## Done when

- `TestRepositoryLintsClean` lives in `flai/internal/mdlint/repo_test.go` and passes in the full `go test` run.
- With `CLOSE_OUT_STORY` set, a markdown finding in a `wip/` file the story does not change is logged, not failed; the new test shows it, and shows that a `wip/` file the story changes and a file outside `wip/` still fail.
- `flai test flai/internal/mdlint` passes.

## Notes

Drafted by the planner for S-0345.
