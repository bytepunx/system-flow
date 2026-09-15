---
title: Tooling
updated: 2026-09-15
audience: agent
order: 90
status: active
---

# Tooling

Use the project's tools for the project's data. The tools keep the standard true; hand edits drift.

## Rules

- Create and move work items with `flai`: `flai epic new`, `flai story new --epic`, `flai task new --story`, `flai move`, `flai block`, `flai unblock`. Do not hand-write front matter for something a command creates.
- Open and log narratives with `flai stream open` and `flai stream log`. The summary sections (`## Context`, `## Current state`, `## Next steps`, `## Decisions`, `## Open questions`) are edited by hand; the `## Log` is appended by the command.
- Hand-edit item bodies (goals, criteria, notes) freely. Hand-edit item front matter only for fields no command sets, and run `flai check` afterwards.
- Run `flai check --strict` before reporting a story as done and before any commit. Fix what it finds; do not explain it away.
- Use `flai board`, `flai show`, and `flai stats` to answer questions about work state instead of grepping files.
- Set `FLAI_AGENT` and `FLAI_SESSION` at session start so narrative entries record who wrote them.
- Archive with `flai archive`, never by moving files.
- Run the dashboard with `flai dashboard` when the operator wants to see the board or the charts; do not build ad hoc views.
- Use the project's build, test, and lint entry points (`Makefile`, `package.json` scripts, `go` commands as documented) rather than reconstructing them.
- When a tool is missing or the wrong version, say so and use a scratch install; do not silently fall back to hand edits or skip the step. Add or extend a script under `scripts/` so the operator can install it properly.
- If a command is missing for something you do repeatedly, propose it as a story instead of scripting around it.
- Put common tasks in `scripts/` as purpose-named shell scripts, so complex commands and command sequences have one name.
- The `Makefile` is the entry point for building, linting, quality checks, dependency updates, tests, and environment management. Make targets call the scripts in `scripts/` rather than inlining commands.
- Scripts and Make targets work the same locally and in CI.
- Local testing, previews, and validation run in Docker and Docker Compose; when a cluster is needed, use k3d via `bytepunx/kluster`.

## When in doubt

- If `flai` can do it, `flai` does it.
- If a hand edit was unavoidable, `flai check` is not optional.

<!-- system-flow:end-of-baseline -->

## Project additions
