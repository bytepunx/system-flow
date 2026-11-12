---
id: T-3151
type: task
nature: feature
title: Describe in design/system/flai-cli.md where the agent images come from as decided, marking each part not yet built
status: backlog
parent: S-0733
owner: alex
created: 2026-11-12T15:55:09Z
updated: 2026-11-12T15:55:09Z
transitions: []
stream: S-0733
tags: [cli]
touches: [design/system/flai-cli.md]
after: [T-3149]
---
# T-3151 Describe in design/system/flai-cli.md where the agent images come from as decided, marking each part not yet built

## Work

Add a section under § Internal structure of `design/system/flai-cli.md`, beside § The container harness (S-0508), that says what T-3149's ADR decides, with a link to it. Bump `updated`. It waits for T-3149 because it links the ADR by the number and file name `flai adr new` gives it, and describes what that ADR decides.

Cover:

- where each harness image comes from, and the signed list;
- the label that ties it to its flai, and what makes it stale;
- what a development tree runs;
- how the layer gets its build arguments;
- the host job that brings an image back, its limits and retry, and what the dashboard shows;
- whether a strategic agent waits;
- the operator's own build command, and the setting that turns the automatic build off.

Mark each part not yet built with the E-0036 story that builds it, as the plan thread on E-0036 names them, so a reader never takes the design for what is built. In § The container harness (S-0508), § The project's layer image, add one sentence pointing to the new section, and leave the rest as it describes what is built today.

## Done when

- `design/system/flai-cli.md` describes the ADR's design as decided, links the ADR, and marks each part not yet built with the story that builds it.
- `flai check --strict` and the markdown lint pass on it.

## Notes
