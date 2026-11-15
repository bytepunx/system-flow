---
id: D-0634
title: The acceptance preview names children-closed and criteria-ticked in unexported constants of its own in flai/internal/preview/accept.go, not workflow.RuleChildrenClosed and RuleCriteriaTicked, so a rename of either rule changes both until a story switches the preview to the workflow package's
kind: decision
sources: [S-0744, T-3227]
topics: [flai, go, code]
paths: [flai/internal/preview/accept.go, flai/internal/workflow/rules.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0634 The acceptance preview names children-closed and criteria-ticked in unexported constants of its own in flai/internal/preview/accept.go, not workflow.RuleChildrenClosed and RuleCriteriaTicked, so a rename of either rule changes both until a story switches the preview to the workflow package's

S-0744 made `flai accept`'s preview report open tasks and unticked criteria only when the acceptance transition lists the `children-closed` and `criteria-ticked` rules (ADR-0301 §4), as `design/system/workflow.md` now says. Its narrative's Decisions (`wip/archive/agents/S-0744.md`) record where the preview gets the rules' names:

> T-3227's preview reports open children and unticked criteria only when the acceptance transition lists children-closed and criteria-ticked. It checks them before Repo.Move, in the same words and order, so the baseline's blockers are unchanged. The rule names are unexported constants in preview, to be switched to S-0738's once it is accepted.

S-0738 is accepted, and `flai/internal/workflow/rules.go` exports `RuleChildrenClosed` and `RuleCriteriaTicked`. But `flai/internal/preview/accept.go` still declares its own `ruleChildrenClosed` and `ruleCriteriaTicked`. So a story that renames either rule, or changes the preview's blockers, changes both files, or first switches the preview to the workflow package's constants. Where a rule's name and its check live is D-0631.

Source: S-0744's narrative, `wip/archive/agents/S-0744.md`, § Decisions (T-3227).
