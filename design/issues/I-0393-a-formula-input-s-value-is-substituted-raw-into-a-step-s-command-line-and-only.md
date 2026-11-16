---
id: I-0393
title: A formula input's value is substituted raw into a step's command line, and only ADR-0273's wording keeps the orchestrator's judgement shell-safe
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-16T06:16:44Z
last_reported: 2026-11-16T06:16:44Z
updated: 2026-11-16T06:16:44Z
usage:
  strategic:
    - kind: analyzer
      seconds: 98
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349359
          cache_write: 34305
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 35
          output: 14842
          cache_read: 2398930
          cache_write: 47317
          cost: 0.1557
---

# I-0393 A formula input's value is substituted raw into a step's command line, and only ADR-0273's wording keeps the orchestrator's judgement shell-safe

## Description
A formula input's value is substituted raw into a step's command line, and only ADR-0273's wording keeps the orchestrator's judgement shell-safe

## Instances

### 2026-11-16T06:16:44Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: no instance observed; a value with ; or $(...) from any caller of flai formula run would run in the agent's shell as the operator.

## Impact
flai/internal/formulas/formulas.go:458-465 Fill replaces {{name}} with the value as given and :477 applies it to s.Run.Flai and s.Run.Script, the command line the step's task tells the story's agent to run in its shell; names are validated at :58 and :749, values are not; flai/internal/mcpserver/role.go:60 builds the --input name=value text. ADR-0273 asks the orchestrator for one sentence with no double quote, backtick, or dollar sign. Kluster, tiers, and the close-out run as argv (serve/checks.go:424, kluster/deploy.go:46, serve/verifyjob.go:325).

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).
