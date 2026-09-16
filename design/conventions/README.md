---
title: Agent conventions
updated: 2026-09-15
---

# Agent conventions

How agents work in this repository. Read every file below, in this order, at the start of every session and before any change. Each file is short. Together they replace the conversation that would otherwise re-establish standards, norms, and ways of working.

| Order | File | Governs |
|-------|------|---------|
| 10 | [session-start.md](session-start.md) | What to read before working, how to recover after a crash, how to end a session |
| 20 | [communication.md](communication.md) | How to report, when to ask, when to proceed, how to raise a concern |
| 30 | [work-management.md](work-management.md) | Pull, WIP, sizing, narratives, definitions of ready and done, blocking |
| 40 | [decisions.md](decisions.md) | What counts as a decision and where it is recorded |
| 50 | [documentation.md](documentation.md) | Where documents live, front matter, style, keeping docs true |
| 60 | [code-quality.md](code-quality.md) | Tests, lint, change size, dependencies |
| 70 | [git.md](git.md) | Commits, branches, history, pull requests |
| 80 | [safety.md](safety.md) | Secrets, destructive actions, untrusted content, sandbox |
| 90 | [tooling.md](tooling.md) | Using flai and the dashboard instead of hand edits |
| 100 | [continuous-improvement.md](continuous-improvement.md) | Recording recurring friction, defects, and blockers in `design/issues`, and reviewing them at delivery |
| 110 | [logging.md](logging.md) | What is logged, at which level, in what shape, and what never appears in a log (draft) |
| 120 | [telemetry.md](telemetry.md) | Metrics, traces, and health signals every service emits, and how they are named (draft) |

## How these files work

- Rules are imperative. Rationale, when it exists, is a link to `design/system` or an ADR.
- Everything above the `<!-- system-flow:end-of-baseline -->` line in each file is the template's baseline and is replaced by `flai upgrade`. Everything below it under `## Project additions` belongs to this project.
- Precedence when rules conflict: an explicit instruction from the operator in the current conversation, then project additions, then the baseline, then your own defaults. Log any conflict in the narrative. Never deviate silently and never edit a baseline rule; propose the edit instead.
- `flai prime` prints these files in order; `flai prime --cat` prints their content.
