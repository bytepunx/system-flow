---
id: S-0053
type: story
nature: improvement
title: "Research and experiment stories can be accepted: findings land on main and are pushed, with no release"
status: backlog
parent: E-0001
owner: alex
created: 2026-09-19T02:05:28Z
updated: 2026-09-19T02:05:28Z
transitions: []
tags: [cli, conventions]
touches: [flai/internal/release, flai/cmd, design/conventions, template]
---

# S-0053 Research and experiment stories can be accepted: findings land on main and are pushed, with no release

## Goal
A research story can be moved to done like any other. Its findings, which may be ADRs, design documents, and new stories, land on main and are pushed, and because a finding is not a deliverable of a component no release is cut. Today acceptance refuses the nature outright.

## Acceptance criteria
- [ ] `flai accept`, `flai move <story> done`, and the dashboard's preview and move all accept a `research` story in review: the branch is merged, the story is archived, the acceptance is committed, and the commit is pushed, with the plan reporting that research releases nothing and why
- [ ] No tag is created and no component version or changelog changes for a research story, whatever files its commits touched; if its commits touched a component's files, the plan says so plainly so the operator can see code is landing without a release
- [ ] The push happens when there is no release to cut, for research and for any other story whose plan is skipped, and a test pins it; from the dashboard the existing behaviour stands (accepted locally, push command shown) until S-0052 decides otherwise
- [ ] `experiment` is decided explicitly in this story, with the operator: the same as research, or left as it is, and the rule says which
- [ ] `design/conventions/git.md` in the template baseline and here replaces "research and experiment stories stay on a branch ... do not release from main" with the new rule; `design/system/workflow.md` and the release design agree; an ADR records the change because the tooling enforces it
- [ ] Tests: `release.LevelFor` and `Compute` for research, acceptance of a research story end to end with real git including the push to a scratch remote, the dashboard preview for a research story; `docs/users/flai.md` updated

## Tasks

## Notes
Raised by the operator on 2026-09-18 while queuing S-0052, itself a research story: "research stories should be able to get moved to done. their findings can produce ADRs, other stories, etc. and those should result in a push even if there is no release to cut."

Where the refusal is: `release.LevelFor` in `flai/internal/release/release.go` returns an error for `research` and `experiment` ("stay on a branch and do not release from main"). `release.Compute` calls it first, so both the dry run and the real acceptance fail, from the shell and from the board; `flai accept --no-release` is the only way through today. The rule comes from `design/conventions/git.md`.

What already works: `acceptItem` pushes whenever `--no-push` is not given, whether or not a release was planned (`flai/cmd/accept.go`), so once the nature stops being an error the push follows. The criterion asks for a test because nothing pins it.

The pre-release idea in `git.md` ("releasing research builds with semver pre-release suffixes ... is not yet defined") is left alone unless the operator wants it here.

The operator blocked S-0052 on 2026-09-19 with the reason "research stories cannot presently be moved to done"; this story clears that block, and the agent that finishes it unblocks S-0052 with `flai unblock`. Until then S-0052 could only be accepted from a shell with `--no-release`.
