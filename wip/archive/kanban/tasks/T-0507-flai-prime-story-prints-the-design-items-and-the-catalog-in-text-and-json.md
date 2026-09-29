---
id: T-0507
type: task
nature: feature
title: flai prime --story prints the design items and the catalog, in text and --json
status: done
parent: S-0137
owner: alex
created: 2026-09-29T00:36:08Z
updated: 2026-09-29T00:43:16Z
transitions:
  - to: ready
    at: 2026-09-29T00:36:14Z
    by: agent-S-0137
  - to: in-progress
    at: 2026-09-29T00:40:24Z
    by: agent-S-0137
  - to: done
    at: 2026-09-29T00:43:16Z
    by: agent-S-0137
stream: S-0137
tags: [cli]
touches: [flai/cmd/prime.go, flai/internal/context]
---
# T-0507 flai prime --story prints the design items and the catalog, in text and --json

## Work

- `primeStory` builds the story's sources and query, loads the documents, and adds the design selection to the pack.
- Text: items after the open issues, each headed with path, heading path, and reason; then the catalog; then what was left out. The header's size includes them.
- `--json`: every item with path, heading path, reason, and size.

## Done when

- `flai/cmd` tests cover the text and JSON output and pass.

## Notes
