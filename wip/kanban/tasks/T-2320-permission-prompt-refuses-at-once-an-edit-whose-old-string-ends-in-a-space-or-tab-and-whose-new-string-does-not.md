---
id: T-2320
type: task
nature: remediation
title: permission_prompt refuses at once an edit whose old string ends in a space or tab and whose new string does not
status: in-progress
parent: S-0539
owner: alex
created: 2026-10-29T03:03:11Z
updated: 2026-10-29T03:58:09Z
transitions:
  - to: ready
    at: 2026-10-29T03:58:08Z
    by: agent-S-0539
  - to: in-progress
    at: 2026-10-29T03:58:09Z
    by: agent-S-0539
stream: S-0539
tags: [flai, mcp, permission-prompt]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go]
---
# T-2320 permission_prompt refuses at once an edit whose old string ends in a space or tab and whose new string does not

## Work

Claude Code drops trailing whitespace from an Edit's `new_string` before `permission_prompt` sees the call, and keeps it on `old_string`: TH-0712 shows an old string ending in `exit 0;` and a space, and a new string ending in `exit 0;` alone. flai hands `input` back as `updatedInput` unchanged, so the write it approves lacks the space, and fixing it costs a second permission round (I-0223).

In `flai/internal/mcpserver/permission.go`, once a request is a protected write inside an in-progress story's worktree, and before both the auto-approve branch and `askOperator`, refuse at once, opening no thread, an `Edit` whose `old_string` ends in a space or a tab while its `new_string` does not, and a `MultiEdit` any of whose `edits` does so. The refusal names the tool, the path, and for a `MultiEdit` the edit's index, says that Claude Code drops trailing whitespace from the new string, and tells the agent to end both strings on a character that is not a space. Add the refusal to `permissionPromptDescription`.

Before building it, confirm the cause from the S-0416 threads (TH-0712, TH-0713, TH-0720) and write it in the narrative's `## Decisions`. If they show something else, such as the model writing the new string without the space, ask on the story's plan thread before going on.

This task waits for no other: it is the first layer.

## Done when

- A test in `flai/internal/mcpserver/permission_test.go` reproduces S-0416's instance: an `Edit` of `.claude/settings.json` in an in-progress story's worktree whose old string ends in `exit 0;` and a space and whose new string ends in `exit 0;`. The test shows it refused at once, with no thread opened, with auto-approve both off and on.
- The same test shows a `MultiEdit` with one such edit refused, naming the edit's index.
- Tests show these edits still go on to the thread or to auto-approve as before: both strings ending on a non-space character, both ending in a space, and an old string ending on a non-space character with a new string that does not.
- `flai test flai/internal/mcpserver/permission.go flai/internal/mcpserver/permission_test.go` passes.

## Notes
