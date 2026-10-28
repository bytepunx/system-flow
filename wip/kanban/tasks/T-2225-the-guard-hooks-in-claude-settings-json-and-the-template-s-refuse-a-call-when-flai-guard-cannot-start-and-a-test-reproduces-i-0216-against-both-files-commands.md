---
id: T-2225
type: task
nature: remediation
title: The guard hooks in .claude/settings.json and the template's refuse a call when flai guard cannot start, and a test reproduces I-0216 against both files' commands
status: backlog
parent: S-0518
owner: alex
created: 2026-10-28T13:34:09Z
updated: 2026-10-28T13:34:09Z
transitions: []
stream: S-0518
tags: [cli, template]
touches: [".claude/settings.json", template/root/.claude/settings.json, flai/cmd/guard_test.go, flai/internal/harness/harness_test.go, template/CHANGELOG.md]
after: [T-2224]
---
# T-2225 The guard hooks in .claude/settings.json and the template's refuse a call when flai guard cannot start, and a test reproduces I-0216 against both files' commands

## Work

Change the two `PreToolUse` commands in `template/root/.claude/settings.json` and `.claude/settings.json` so that the guard's exit code is kept and judged: 0 allows; 2 refuses with the guard's output, as today; 126 and 127 refuse with `flai guard did not run (exit N): <output>` and one line saying no call passes until flai is installed and on `PATH`; any other exit allows, as ADR-0060 has it. One shape for both files, the template's running `flai guard` and this repository's `"$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard`, and the `Edit|MultiEdit|Write|NotebookEdit` command keeps its `editGuardCondition` prefix unchanged, since `checkGuardHooks` pins it; leave the `SubagentStart` and `SubagentStop` hooks as they are. Both files are protected paths: make each write yourself, not through a sub-agent, and the write asks the operator on a thread through `permission_prompt` unless `auto-approve` is on (ADR-0086); make the same write again once it is answered. In `flai/cmd/guard_test.go`, extend `checkGuardHooks` to run each file's Bash hook command under `sh -c` with a fake guard on a temporary `PATH`, and assert: a guard exiting 0 gives 0; one exiting 2 with text gives 2 and the text on standard error; one exiting 1 gives 0; no guard at all gives 2 with `did not run (exit 127)`. For this repository's command the fake is reached through a copy of `scripts/flai.sh`, `scripts/env.sh`, and `scripts/flai-build.sh` under a temporary `CLAUDE_PROJECT_DIR` with an empty `flai/` folder and no `bin/`, and `PATH` without `go`: the case that reproduces I-0216, which exits 0 on the old command and 2 on the new one when `flai` is absent too, and runs the fake when it is present (T-2224). Adapt the substring checks in `checkGuardHooks` and in `TestTheTemplatesAnalyzerIsHeldToItsReport` (`flai/internal/harness/harness_test.go`) if the new command no longer contains `guard 2>&1`. Add the template's entry to `template/CHANGELOG.md`: the hook refuses a call when `flai guard` cannot start. Waits for T-2224, whose fallback the reproducing case exercises.

## Done when

Both settings files carry the new commands, `go test ./cmd -run 'TestTheSettingsRunTheGuard|TestGuard'` and `go test ./internal/harness -run TestTheTemplatesAnalyzer` pass with the new cases, the reproducing case fails on the old command and passes on the new, `template/CHANGELOG.md` names the change, and `flai check --strict` is clean.

## Notes

A story whose branch changes a path Claude Code protects is accepted by the operator alone (ADR-0106); say so when the story goes to review.
