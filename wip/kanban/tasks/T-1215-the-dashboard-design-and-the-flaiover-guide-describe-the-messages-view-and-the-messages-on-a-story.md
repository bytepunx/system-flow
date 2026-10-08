---
id: T-1215
type: task
nature: feature
title: The dashboard design and the flaiover guide describe the Messages view and the messages on a story
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:17:09Z
updated: 2026-10-08T04:32:16Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1214]
---
# T-1215 The dashboard design and the flaiover guide describe the Messages view and the messages on a story

## Work

Document what T-1211 to T-1214 built. It waits for T-1214, the last of them.

- `design/system/flai-cli.md`: the host API reads `messages.list` and `messages.get`, beside `threads.list`.
- `design/system/flaiover-dashboard.md`: the reads, the view, and the story page.
- `docs/users/flaiover.md`: where to find the agents' messages, and that they never reach the operator's inbox.

## Done when

- The three documents name the view or the reads.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
