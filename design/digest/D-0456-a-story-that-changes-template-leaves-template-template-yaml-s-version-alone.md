---
id: D-0456
title: "A story that changes template/ leaves template/template.yaml's version alone: flai release bumps it at publish, and the story's changelog lines go under ## Unreleased"
kind: decision
sources: [S-0421, S-0524, T-1773]
topics: [template, planning, git]
paths: [template/template.yaml, template/CHANGELOG.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0456 A story that changes template/ leaves template/template.yaml's version alone: flai release bumps it at publish, and the story's changelog lines go under ## Unreleased

A story, or a task a planner writes for it, does not bump the `version` line of `template/template.yaml`. `flai release` sets it when it publishes. The story writes its `template/CHANGELOG.md` lines under `## Unreleased`, which `flai release` then names with the version ([ADR-0228](../adrs/0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md) § 4, `design/conventions/git.md`).

- S-0421, against T-1773's Work, which asked for a bump: "`template/template.yaml` is not bumped, though T-1773's Work asked for it: `flai release` bumps it at publish (the last five bumps are publish commits), and the changelog line goes under `## Unreleased` (git.md)."
- S-0524: "`template.yaml` is bumped at publish, not here."

Every recent commit that changes `template/template.yaml` is a `chore: publish ...` commit. A planner should not write a version bump into a task's Work. `docs/contributors/template.md` § Versioning says the release tooling bumps the version in this monorepo. The convention does not yet say a story leaves it alone, which is why D-0410 recorded earlier stories doing both.

Sources: the Decisions of wip/archive/agents/S-0421.md and S-0524.md; T-1773.
