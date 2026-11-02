---
id: MS-0326
title: "S-0510 holds S-0554 on overlap alone: may S-0554 start beside it?"
from: S-0554
to: S-0510
about: [docs/operators/index.md]
status: open
participants: [flai]
created: 2026-11-02T20:31:42Z
updated: 2026-11-02T20:31:42Z
---

# MS-0326 S-0510 holds S-0554 on overlap alone: may S-0554 start beside it?

Between S-0554 and S-0510, about `docs/operators/index.md`.

## Entries

### 2026-11-02T20:31:42Z flai S-0554
S-0510 holds S-0554 on overlap alone: may S-0554 start beside it?

S-0554, "A role file's start prompt and its triggers replace the planner's, the orchestrator's, the analyzer's, and the curator's start prompts and schedules in Go, so the built-in roles are defined by their files alone", is ready, and its claim overlaps S-0510's on `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> After S-0419 and S-0553, each built-in strategic role's definition comes from its role file. Two things still live in Go. The start prompts, `planPrompt`, `orchestratePrompt`, `analyzePrompt`, `curatePrompt`, and `definedPrompt` in `flai/internal/harness/harness.go`, say what each run is for. The rules for when the planner, the analyzer, and the curator run are in `flai/internal/serve/replan.go`, `analysis_schedule.go`, and `curate_schedule.go`. Their role files' `trigger` entries are validated but never read, and the conditions `stale-figures`, `replan-agent`, and `each-ready-story` are named but not implemented. This story moves both into the role files: a start-prompt section templated over the run's fields, under a new ADR refining ADR-0207 § 4, and the triggers read through the event and schedule launchers S-0418 built. Then the four built-in roles are defined by their files alone, as the operator asked on TH-0760, and a project changes what a built-in role is told, or when it runs, with no flai release.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0510 --remove <path>`, or on the task that names the path, if S-0510 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0326 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0554 can start.
3. Reply with `message_reply` saying why the hold stands: S-0554 waits until S-0510 moves to review, is cancelled, or is sent back.
