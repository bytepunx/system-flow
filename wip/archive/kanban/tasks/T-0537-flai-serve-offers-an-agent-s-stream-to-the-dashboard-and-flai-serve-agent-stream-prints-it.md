---
id: T-0537
type: task
nature: feature
title: flai serve offers an agent's stream to the dashboard and flai serve agent stream prints it
status: done
parent: S-0142
owner: alex
created: 2026-09-29T05:44:46Z
updated: 2026-09-29T05:51:27Z
transitions:
  - to: ready
    at: 2026-09-29T05:44:51Z
    by: agent-S-0142
  - to: in-progress
    at: 2026-09-29T05:47:40Z
    by: agent-S-0142
  - to: done
    at: 2026-09-29T05:51:27Z
    by: agent-S-0142
stream: S-0142
tags: []
touches: [flai/internal/hostapi, flai/cmd, docs/users, design/system/flai-cli.md]
---
# T-0537 flai serve offers an agent's stream to the dashboard and flai serve agent stream prints it

## Work

Add the read-only channel method `agent.stream` (`story`, `after`) to `hostapi`, answered by `flai serve` through `hostapi.Host`, and the command `flai serve agent stream <story-id>` with `--follow` and `--json`. Document the method in `design/system/flai-cli.md` and the command in `docs/users/flai.md` and the generated reference.

## Done when

- A hostapi test asks `agent.stream` and gets the entries and the next offset; a bad story ID is refused as invalid params.
- `flai serve agent stream` prints the stream of a story's newest agent, and `--follow` goes on until the agent ends.
- `make test`, lint, and `make flai-reference` are clean.

## Notes
