---
id: T-2321
type: task
nature: improvement
title: Document permission_prompt's refusal of an edit whose new string lost the old string's trailing space
status: backlog
parent: S-0539
owner: alex
created: 2026-10-29T03:03:22Z
updated: 2026-10-29T03:03:22Z
transitions: []
stream: S-0539
tags: [docs, permission-prompt]
touches: [docs/users/flai.md, design/system/agent-adapters.md]
after: [T-2320]
---
# T-2321 Document permission_prompt's refusal of an edit whose new string lost the old string's trailing space

## Work

Say in the documents what T-2320 made `permission_prompt` refuse at once, and why, in the same words as its refusal:

- `docs/users/flai.md`: the `permission_prompt` row of the MCP tools table, and the paragraph of the protected-writes section that lists what it refuses at once with no thread.
- `design/system/agent-adapters.md`: the `permission_prompt` section, where it lists what is denied at once, naming I-0223.

It waits for T-2320 because it describes the refusal that task builds, in that task's words.

## Done when

- Both documents say that an `Edit`, or an edit of a `MultiEdit`, whose old string ends in a space or a tab while its new string does not is refused at once, with no thread and under auto-approve too, because Claude Code drops trailing whitespace from the new string.
- `flai test docs/users/flai.md design/system/agent-adapters.md` passes.

## Notes
