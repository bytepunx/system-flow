---
id: T-1118
type: task
nature: improvement
title: flai-cli.md, the user guide, and the reference describe flai story start, story_start, and story.start
status: backlog
parent: S-0274
owner: alex
created: 2026-10-06T22:53:57Z
updated: 2026-10-06T22:53:57Z
transitions: []
stream: S-0274
tags: [docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1109, T-1111]
---
# T-1118 flai-cli.md, the user guide, and the reference describe flai story start, story_start, and story.start

## Work

Describe the new entry point where the CLI, MCP, and the host channel are documented:

- **`design/system/flai-cli.md`.** Add `flai story start` to the commands: what it composes, what it refuses (a story not ready or held), its result's keys, and that `story_start` over MCP and `story.start` on the host channel answer the same. Where the design describes pulling a story as `item_move` and then `flai stream open`, make it the one call.
- **`docs/users/flai.md`.** Begin the story loop with `flai story start`, add `story_start` to the MCP tools table, and add `story.start` where the host channel's write methods are listed.
- **`docs/users/flai-reference.md`.** Add the command, with its flags (`--budget`, `--json`) and its output.

It waits for T-1109 and T-1111, whose flags and keys it documents. It runs alongside the host channel and prompt tasks, which touch no file this one does.

## Done when

- The three documents describe `flai story start`, `story_start`, and `story.start`, with the same keys and refusals as the code.
- None of them still tells an agent to pull a story with `item_move` and then `flai stream open`.
- `flai check --strict` and the markdown lint pass.

## Notes
