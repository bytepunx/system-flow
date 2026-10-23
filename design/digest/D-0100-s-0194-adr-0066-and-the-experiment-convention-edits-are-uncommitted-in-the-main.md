---
id: D-0100
title: "S-0194: ADR-0066 and the experiment convention edits are uncommitted in the main checkout; may they be committed to main so my branch can sync them: done, I have added the items and committed locally…"
kind: answer
sources: [TH-0060]
topics: [release, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/preview, flai/internal/release, flai/cmd/accept.go, flai/cmd/accept_research_test.go, flai/internal/experiment, flai/internal/check, flaiover/src, template, design/experiments, CLAUDE.md, design/system/workflow.md, design/system/work-hierarchy.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/repository-layout.md, docs/users, docs/operators, design/README.md]
date: 2026-10-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0100 S-0194: ADR-0066 and the experiment convention edits are uncommitted in the main checkout; may they be committed to main so my branch can sync them: done, I have added the items and committed locally…

The operator, alex, answered TH-0060 at 2026-10-02T10:11:37Z.

What was asked:

```text
The story is built on ADR-0066 and on your edits to `git.md` and `documentation.md` (here and in `template/root/`). Right now they exist only in the main checkout: ADR-0066 to ADR-0069 are untracked, and the conventions, `design/adrs/README.md`, and others are modified but not committed. My worktree branches from `main` HEAD (05b67f9), so it has neither the ADR nor the new convention text. That matters in two ways. First, the design documents I have to edit link ADR-0066, and `flai check` in the worktree would report those links as broken. Second, if I copy the ADR or the convention text onto…
```

The answer:

```text
done, I have added the items and committed locally in main
```

Source: TH-0060.
