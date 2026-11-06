---
id: T-2470
type: task
nature: research
title: design/system/kluster.md describes how flai runs a project's kluster as the ADR decides, and the design index lists it
status: backlog
parent: S-0576
owner: alex
created: 2026-11-06T21:42:07Z
updated: 2026-11-06T21:42:07Z
transitions: []
stream: S-0576
tags: [kluster, design]
touches: [design/system/kluster.md, design/system/README.md]
after: [T-2469]
---
# T-2470 design/system/kluster.md describes how flai runs a project's kluster as the ADR decides, and the design index lists it

## Work

It waits for T-2469, whose accepted ADR it describes.

Write `design/system/kluster.md`, with front matter (`title`, `updated`, `status`) and the ADR's topics. Describe, without restating the ADR's reasons, and linking it for them:

- the project's cluster and its shared infrastructure;
- the `kluster:` block, key by key;
- a story's namespace and the default namespace;
- the build, load, deploy, and teardown steps and the variables each command gets;
- what happens when a story enters review, gets a new head, leaves review, or is accepted;
- the `kluster` host action;
- the Ops pages in the dashboard.

Leave a section per later story to fill in: the commands (S-0577, S-0578), flai serve (S-0579), and the dashboard (S-0580, S-0581). Each says what it will hold. Add the document to `design/system/README.md`.

## Done when

- `design/system/kluster.md` exists, matches the ADR, and links it.
- `design/system/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes
