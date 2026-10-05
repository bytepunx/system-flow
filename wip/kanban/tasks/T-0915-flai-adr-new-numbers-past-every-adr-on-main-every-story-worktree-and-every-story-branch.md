---
id: T-0915
type: task
nature: remediation
title: flai adr new numbers past every ADR on main, every story worktree, and every story branch
status: backlog
parent: S-0245
owner: alex
created: 2026-10-05T05:44:13Z
updated: 2026-10-05T05:44:13Z
transitions: []
stream: S-0245
tags: [flai, adr]
touches: [flai/internal/adr, flai/cmd/adr.go, flai/cmd/adr_test.go, docs/users/flai-reference.md]
---
# T-0915 flai adr new numbers past every ADR on main, every story worktree, and every story branch

## Work

Proposed fix, from I-0063's instances. `adr.NextNumber` (`flai/internal/adr/adr.go`) reads only the `design/adrs` of the checkout it runs in. `adr.New` already takes a runner. Pass that runner to the numbering, and number one past the highest `NNNN-*.md` among three sets of names:

- the current checkout's own files;
- the names `storygit.FolderNames` returns for `design/adrs`, the helper S-0252 added for issues. It covers the main checkout, every linked worktree (uncommitted files included), the main branch, and every local `story/*` branch.

Leave the checks of `--supersedes` and `--refines` reading the local files. An ADR can only supersede or refine an ADR its own branch holds.

Say in `flai adr new`'s help that the number is past every open story's ADRs. Then regenerate the reference with `make flai-reference`.

It waits for no task.

## Done when

- A test in `flai/cmd/adr_test.go` reproduces I-0063, as S-0252's test does for issues. The first story branch has ADR-0002 committed, and the second has ADR-0003 uncommitted in its worktree. `flai adr new` run in a third story's worktree gives ADR-0004. The test fails without the change.
- A project with no story branches still numbers from its own `design/adrs`, as before, and the existing `TestAdrNew*` tests pass.
- `docs/users/flai-reference.md` matches the help text, and `scripts/flai-test.sh` passes.

## Notes
