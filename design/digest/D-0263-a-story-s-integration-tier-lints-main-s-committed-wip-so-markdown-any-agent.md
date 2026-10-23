---
id: D-0263
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
kind: cause
sources: [I-0117]
topics: [testing, flai, cli, go, code, all]
paths: [flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/repo_test.go, scripts/lint-md.sh, scripts/README.md, design/system/flai-cli.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/summary.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0263 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

I-0117 was closed on 2026-10-08.

The cause:

```text
A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
```

The fix:

```text
S-0345: in a close-out (`CLOSE_OUT_STORY` set), `TestRepositoryLintsClean` (now in `flai/internal/mdlint/repo_test.go`) and `scripts/lint-md.sh` leave out the markdown under `wip/` that the story does not change. They do so as `TestMonorepoIsClean` and `scripts/check.sh` scope `flai check` (ADR-0085). A line another agent commits to main's `wip/` is a logged note in the integration tier and is left out of smoke's whole-repository lint. A `wip/` file the story changes, and every file outside `wip/`, still fails. `TestMainsWipOutsideTheStoryOfI0117` reproduces the case. `make lint-md` and CI do not set the variable, so they still lint all of main's `wip/`.
```

Source: I-0117.
