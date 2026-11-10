---
id: T-2946
type: task
nature: improvement
title: The design and the user guide say what a container run records and logs
status: backlog
parent: S-0665
owner: alex
created: 2026-11-10T08:50:46Z
updated: 2026-11-10T08:50:46Z
transitions: []
stream: S-0665
tags: [docs, containers]
touches: [design/system/flai-cli.md, docs/users/flai.md, design/system/strategic-agents.md]
after: [T-2943, T-2945]
---
# T-2946 The design and the user guide say what a container run records and logs

## Work

Say what a container run now records and logs, where `harness_version` is documented:

- `design/system/flai-cli.md`, "The container harness (S-0508)", the paragraph on the start in `flai serve`: the `Start` field the adapter fills, the run's `image`, `container`, and `mounts`, and the journal's start line.
- `docs/users/flai.md`, "In a container: the container harness": that `agents.json` names each container run's image, container, and mounts, and the journal names the image and container, so a run can be shown afterwards to have used the project's layer image.
- `design/system/strategic-agents.md`, "In a container (S-0509)" › "The live run (S-0567)": the I-0274 bullet says this story fixed it.

It waits for T-2943 and T-2945, whose fields and line it describes.

## Done when

- The three sections describe the recorded fields and the journal line as built.
- The markdown lint passes on the three files.

## Notes
