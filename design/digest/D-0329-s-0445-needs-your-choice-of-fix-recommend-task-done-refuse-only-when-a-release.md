---
id: D-0329
title: "S-0445 needs your choice of fix: recommend task done refuse only when a release changed its rules, recorded in each flai tag: Recommendation: A now, and B as S-0445"
kind: answer
sources: [TH-0506]
topics: [flai, cli, go, code, all]
paths: [design/adrs/README.md, flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/buildinfo_test.go, flai/cmd/version.go, flai/cmd/version_test.go, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/adrs/0180-flai-task-done-refuses-an-older-flai-only-when-the-newest-flai-release-raised.md, flai/internal/buildinfo/task-done-rules.txt, docs/users/flai-reference.md, flai/cmd/task_done.go]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0329 S-0445 needs your choice of fix: recommend task done refuse only when a release changed its rules, recorded in each flai tag: Recommendation: A now, and B as S-0445

The operator, alex, answered TH-0506 at 2026-10-20T21:12:38Z.

What was asked:

```text
**Recommended: B.** S-0445 needs one decision before I can draft its tasks. Which fix should it build?
```

The answer:

```text
**Recommendation:** A now, and B as S-0445.

- **A now.** Run `flai serve enable auto-upgrade` in the main checkout. It costs nothing to build, and it stops the refusals within one publish cycle today.
- **B as S-0445.** Make `task done` refuse only when a release changed its rules. That ends the refusals for good, including the minutes between a tag and its publication, which A cannot cover.
- **C is not needed with A on.** With auto-upgrade on, the upgrade threads stop piling up.

## A note on my part in this

I have published after every acceptance under the `judgement` policy: ten minor flai releases in about 28 hours. That cadence is what keeps the host behind ADR-0142's minor threshold.

If you prefer, I can batch releases instead, for example once per epic or once a day, with the reasons in the log. That would cut the churn without any change to flai. Say so here, and I will hold publishing to that rhythm from now on.

## Why you

ADR-0142's threshold is yours (TH-0401), and B refines it with a new ADR. Turning on a host action is also yours. Neither is mine to settle. The planner drafts S-0445's tasks once you answer.

Source: design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md § Decision
```

Source: TH-0506.
