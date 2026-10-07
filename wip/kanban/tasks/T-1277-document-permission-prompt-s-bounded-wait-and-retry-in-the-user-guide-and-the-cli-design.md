---
id: T-1277
type: task
nature: improvement
title: Document permission_prompt's bounded wait and retry in the user guide and the CLI design
status: backlog
parent: S-0309
owner: alex
created: 2026-10-07T23:27:35Z
updated: 2026-10-07T23:27:35Z
transitions: []
stream: S-0309
tags: [docs, permission-prompt]
touches: [docs/users/flai.md, design/system/flai-cli.md]
after: [T-1274]
---
# T-1277 Document permission_prompt's bounded wait and retry in the user guide and the CLI design

## Work

It waits for T-1274 because it documents what the ADR decides, and links to it.

- `docs/users/flai.md` § Writes to paths Claude Code protects: the agent waits a few minutes, not until you answer; the thread stays open; your answer is taken when the agent makes the write again, after flai serve restarts it on the answer. Replace the paragraph's `.flai-cache/` and `cp` workaround with that.
- `design/system/flai-cli.md`, the `flai mcp` row's `permission_prompt` entry: the bound, the open thread, and the retry, with a link to the new ADR.

## Done when

- Both documents say what the ADR decides and link to it, and neither still says the agent waits until the operator answers.
- `flai test docs/users/flai.md design/system/flai-cli.md` passes.

## Notes
