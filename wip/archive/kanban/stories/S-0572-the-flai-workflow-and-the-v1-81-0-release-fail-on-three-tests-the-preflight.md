---
id: S-0572
type: story
nature: remediation
title: "The flai workflow and the v1.81.0 release fail on three tests: the preflight fixture's GNU-only sed -i on macOS, git 2.56's new gitfile message, and the dangling .claude/agents links the repository lint reads"
status: done
owner: alex
created: 2026-11-02T23:58:32Z
updated: 2026-11-03T00:33:38Z
transitions:
  - to: ready
    at: 2026-11-02T23:58:38Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T23:58:40Z
    by: system-flow
  - to: review
    at: 2026-11-03T00:10:35Z
    by: system-flow
  - to: done
    at: 2026-11-03T00:33:38Z
    by: alex
tags: [cli]
topics: [testing, ci]
touches: [flai/internal/mergequeue/preflight_test.go, flai/internal/storygit/mounts_test.go, flai/internal/mdlint/repo_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0572 The flai workflow and the v1.81.0 release fail on three tests: the preflight fixture's GNU-only sed -i on macOS, git 2.56's new gitfile message, and the dangling .claude/agents links the repository lint reads

## Goal

The `flai` workflow has failed on every push since 2026-11-02, and the `release flai` workflow failed for `flai/v1.81.0`, so no 1.81.0 binaries were published and the host's installed flai is behind the project's history (TH-0817). Three tests fail, each for its own cause, none a defect in flai's shipped behaviour:

1. On macOS, `TestThePreflightConfirmsALintFindingInAScratchWorktree/the_fix_clears_it` (`flai/internal/mergequeue/preflight_test.go`) fails since S-0549: the fixture's markdown tier fixes a file with `sed -i s/FIXABLE/fixed/`, which BSD sed reads as a backup suffix, so the fix never runs and the finding survives.
2. On ubuntu and in GoReleaser's `go test ./...` hook, `TestAWorktreeWithoutItsRegistrationIsNotARepository` (`flai/internal/storygit/mounts_test.go`, S-0517) expects git to say `not a git repository` for a worktree whose registration is gone. git 2.56.0 on the runner says `gitfile does not point to a valid repository`; git 2.47.3 on the host says the old text, so the test passed when it was written.
3. In the integration tier, `TestRepositoryLintsClean` (`flai/internal/mdlint/repo_test.go`) fails since S-0553: the four strategic agents' `.claude/agents/*.md` are committed relative links into `.flai-cache/agents/`, which is gitignored and written only by `flai role definition --write`. On any fresh checkout, CI's among them, the links dangle, and the walk reads them.

This story makes the three tests true on every host the workflows run on, changing no product code.

## Acceptance criteria
- [x] The preflight test's markdown tier fixes a file with no `sed -i` or other flag whose meaning differs between GNU and BSD sed, and `the fix clears it` passes on this host
- [x] `TestAWorktreeWithoutItsRegistrationIsNotARepository` names both of git's messages for a worktree without its registration, `not a git repository` from before 2.56 and `gitfile does not point to a valid repository` from 2.56 on, and passes on this host
- [x] `TestRepositoryLintsClean` passes on a checkout where the strategic agents' `.claude/agents` links dangle, leaving a dangling link out as markdownlint-cli2 does, and still lints a link that resolves
- [x] `go test ./...` passes in `flai/` on this host with the links dangling, and the three tests pass under `-race`

## Tasks
- T-2446 The preflight test's markdown tier fixes a file without sed -i, so the fix runs under BSD sed on macOS as under GNU sed
- T-2447 TestAWorktreeWithoutItsRegistrationIsNotARepository accepts git 2.56's gitfile message as well as the older not-a-git-repository one
- T-2448 TestRepositoryLintsClean leaves a dangling symbolic link out of its walk, as markdownlint-cli2 does, so the strategic agents' .claude/agents links fail it on no fresh checkout

## Notes

- The failing runs: the `flai` workflow at 41532192913 (macOS and ubuntu), `release flai` at 41532192894 for `flai/v1.81.0`. The macOS failure begins with the first run after 05df1a48a (S-0549); the two ubuntu failures with the first run after f37858dcf (S-0517) and c9d617abb (S-0553), both accepted after the 1.80.2 push.
- git's wording changed in `setup.c` between v2.47.3 (`not a git repository: %s`) and v2.56.0 (`gitfile does not point to a valid repository: %s`). `GitMounts` reads no git message, so only the test is affected.
- The links are ADR-0236's design and stay as they are; S-0553 left running `flai role definition --write` in the main checkout to the operator, which also clears the failed curator and orchestrator starts the inbox reports. This story does not run it.
- No flai changelog is written by hand; the release tooling writes it.
- Criteria 1 and 2 first promised a pass on macOS and on git 2.56, which no host here can run (the host has git 2.47.3 and no macOS). On the operator's word they were reworded to what the close-out verifies: the shape of the fix command and of the assertion, and a pass here. The macOS and git 2.56 runs themselves are the `flai` workflow's macOS and ubuntu jobs on the acceptance commit.
