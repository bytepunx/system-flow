---
id: D-0648
title: The template's design/workflows/board.md must equal flai's embedded baseline, flai/internal/workflow/baseline/board.md, byte for byte, so a story that changes the baseline board copies it to the template in the same story
kind: decision
sources: [S-0755, T-3260]
topics: [template, flai, go, code]
paths: [flai/internal/workflow/baseline/board.md, template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go]
date: 2026-11-16
status: absorbed
absorbed_by: "design/system/template.md#upgrading-a-project"
asked_again: []
cited: []
---

# D-0648 The template's design/workflows/board.md must equal flai's embedded baseline, flai/internal/workflow/baseline/board.md, byte for byte, so a story that changes the baseline board copies it to the template in the same story

S-0755 (T-3260) made the template ship flai's baseline workflow as `template/root/design/workflows/board.md`, as [ADR-0301](../adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md) § 7 asks, and held the two files equal with a test. `TestBaselineIsTheTemplatesBoard` in `flai/internal/workflow/baseline_test.go` fails when the template's file differs from the embedded baseline, front matter and body, and its message names the copy that fixes it, from `flai/internal/workflow/baseline/board.md` to `template/root/design/workflows/board.md`.

So a story that changes the baseline board, its states, transitions, rules, or timers, copies the file to the template in the same story. A project made from the template then runs what a project without the file runs. It is the same rule D-0433 gives for the baseline role files.

S-0755's narrative records the case at hand (`wip/archive/agents/S-0755.md`, § Decisions):

> S-0750 changes flai/internal/workflow/baseline/board.md (timers): whichever of the two is accepted second updates template/root/design/workflows/board.md to match, since T-3260's byte test holds them equal.

S-0755 was accepted first, on 2026-11-16, so S-0750, in progress then, makes the copy.

This repository has no `design/workflows/board.md` of its own yet. S-0755 decided it runs the baseline, and its upgrade adds the file. Once a project holds the file, `flai upgrade` never writes it again ([ADR-0307](../adrs/0307-flai-reads-the-workflow-file-of-the-checkout-it-opens-refuses-to-open-a.md) § 4), so the byte test binds only the template's copy.

Sources: S-0755's narrative, § Decisions; T-3260.
