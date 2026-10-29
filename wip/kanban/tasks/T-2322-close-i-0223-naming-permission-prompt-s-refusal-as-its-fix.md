---
id: T-2322
type: task
nature: remediation
title: Close I-0223 naming permission_prompt's refusal as its fix
status: backlog
parent: S-0539
owner: alex
created: 2026-10-29T03:03:28Z
updated: 2026-10-29T03:03:28Z
transitions: []
stream: S-0539
tags: [issues]
touches: [design/issues/I-0223-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md, design/issues/summary.md]
after: [T-2320]
---
# T-2322 Close I-0223 naming permission_prompt's refusal as its fix

## Work

Run `flai issue close I-0223 --reason "<reason>"` in the story's worktree. The reason says that Claude Code drops trailing whitespace from an Edit's new string, and that since S-0539 `permission_prompt` refuses at once, before any thread, an edit whose old string ends in a space or a tab while its new string does not, telling the agent to end both strings on a non-space character.

It waits for T-2320 because the issue is closed only once the fix it names is built and tested. It runs beside the documentation task: the two share no path.

## Done when

- I-0223's status is closed, with that reason, and `design/issues/summary.md` no longer lists it.
- Criterion 2 of S-0539 is ticked.

## Notes
