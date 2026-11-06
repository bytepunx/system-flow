---
id: T-2471
type: task
nature: research
title: "The tooling convention's project addition no longer says kluster is not used here: system-flow runs no cluster itself, and flai drives a project's kluster behind the kluster host action"
status: backlog
parent: S-0576
owner: alex
created: 2026-11-06T21:42:12Z
updated: 2026-11-06T21:42:12Z
transitions: []
stream: S-0576
tags: [kluster, conventions]
touches: [design/conventions/tooling.md]
after: [T-2469]
---
# T-2471 The tooling convention's project addition no longer says kluster is not used here: system-flow runs no cluster itself, and flai drives a project's kluster behind the kluster host action

## Work

It waits for T-2469, whose accepted ADR names the host action and the block this line points to.

In `design/conventions/tooling.md`, below the `system-flow:end-of-baseline` marker, the project addition reads "No cluster: flaiover runs as a single Docker container via `flai dashboard`; k3d and `bytepunx/kluster` are not used here." Reword it so that it says:

- system-flow itself runs and tests with no cluster, and flaiover is still one Docker container run by `flai dashboard`;
- flai drives a project's kluster, from its `kluster:` block and behind the `kluster` host action, and links `design/system/kluster.md` and the ADR.

Change nothing above the marker. Bump `updated`.

## Done when

- The project addition says both points and links the design document and the ADR.
- The baseline above the marker is unchanged.
- `flai check --strict` and the markdown lint pass on the file.

## Notes
