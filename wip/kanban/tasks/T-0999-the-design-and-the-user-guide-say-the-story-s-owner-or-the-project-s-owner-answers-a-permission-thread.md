---
id: T-0999
type: task
nature: remediation
title: The design and the user guide say the story's owner or the project's owner answers a permission thread
status: backlog
parent: S-0284
owner: alex
created: 2026-10-06T06:24:19Z
updated: 2026-10-06T06:24:19Z
transitions: []
stream: S-0284
tags: [docs]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-0997]
---
# T-0999 The design and the user guide say the story's owner or the project's owner answers a permission thread

## Work

Say who answers a `permission_prompt` thread as T-0997's ADR decides, linking the ADR:

- `design/system/flai-cli.md`, the `flai mcp` row's `permission_prompt` text, which says "an entry after the question by the story's owner (anyone but the agent when the story has none)".
- `docs/users/flai.md`, the `permission_prompt` row of the MCP tools table ("once the story's owner answers `allow`") and the section on `.claude/` writes ("Reply as the story's owner").

Say it in the reader's terms: reply as the story's owner or as the project's owner, the `owner` in `system-flow.yaml` and the name the dashboard replies as. Leave `design/conventions/delegation.md` as it is: it says "asks the operator", which stays true.

## Done when

- both documents name the project's owner beside the story's owner and link the ADR
- `flai check --strict` and the markdown lint are clean on both

## Notes
