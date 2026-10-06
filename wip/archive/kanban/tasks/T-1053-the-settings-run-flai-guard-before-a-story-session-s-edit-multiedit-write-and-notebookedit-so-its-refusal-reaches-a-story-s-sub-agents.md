---
id: T-1053
type: task
nature: remediation
title: The settings run flai guard before a story session's Edit, MultiEdit, Write, and NotebookEdit, so its refusal reaches a story's sub-agents
status: done
parent: S-0299
owner: alex
created: 2026-10-06T22:38:20Z
updated: 2026-10-06T22:53:04Z
transitions:
  - to: ready
    at: 2026-10-06T22:39:17Z
    by: agent-S-0299
  - to: in-progress
    at: 2026-10-06T22:39:17Z
    by: agent-S-0299
  - to: done
    at: 2026-10-06T22:53:04Z
    by: agent-S-0299
stream: S-0299
tags: []
touches: [".claude/settings.json", template/root/.claude/settings.json, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
after: [T-1035]
usage:
  source: log
  seconds: 827
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8199
      cache_read: 1135290
      cache_write: 36867
      cost: 0.6219
---
# T-1053 The settings run flai guard before a story session's Edit, MultiEdit, Write, and NotebookEdit, so its refusal reaches a story's sub-agents

## Work

T-1035's review found that the rule it added never runs in a story's session. The `PreToolUse` entry on `Edit|Write|NotebookEdit` in `.claude/settings.json` and in `template/root/.claude/settings.json` begins `[ "$FLAI_ROLE" = plan ] || [ "$FLAI_ROLE" = orchestrate ] || [ "$FLAI_ROLE" = analyze ] || exit 0;`. A story's agent's session has `FLAI_STORY` set and no `FLAI_ROLE`, so the hook exits before `flai guard` sees the sub-agent's write, and `permission_prompt` holds it as in I-0093. The planner's note that no settings change was needed was wrong.

- Add `|| [ -n "$FLAI_STORY" ]` to that condition in both files, and `MultiEdit` to the matcher, which ADR-0086 and the guard name.
- Update `TestTheSettingsRunTheGuardOnTheStrategicAgentsEdits` in `flai/cmd/guard_test.go`, which pins the prefix and says a story's edits are not guarded.
- Make the last paragraph of `flai guard --help` name a story's session among those that run the guard before file writes, and regenerate `docs/users/flai-reference.md`.
- The two settings files are `.claude/` writes: the story's agent makes them, after the rest of the story's work, through `permission_prompt`.

## Done when

- Both settings files run `flai guard` before a story session's `Edit`, `MultiEdit`, `Write`, and `NotebookEdit`, and still exit at once in an operator's own session.
- The settings test pins the new condition and matcher in both files.
- `flai guard --help` and `docs/users/flai-reference.md` say so, and `go test ./cmd/...` passes in `flai/`.

## Notes

Added by agent-S-0299 after T-1035's review. The guard costs about 18 ms a call, measured with `bin/flai`, so running it before every write in a story's session is cheap.
