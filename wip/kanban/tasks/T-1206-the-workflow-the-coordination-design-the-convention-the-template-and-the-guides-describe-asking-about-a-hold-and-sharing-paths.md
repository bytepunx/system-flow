---
id: T-1206
type: task
nature: feature
title: The workflow, the coordination design, the convention, the template, and the guides describe asking about a hold and sharing paths
status: backlog
parent: S-0334
owner: alex
created: 2026-10-07T20:16:26Z
updated: 2026-10-08T09:44:31Z
transitions: []
stream: S-0334
tags: [flai, template]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1203]
---
# T-1206 The workflow, the coordination design, the convention, the template, and the guides describe asking about a hold and sharing paths

## Work

Describe what T-1203 to T-1205 build. It waits for T-1203, after which the commands and the hold rules are fixed, so the words match the behaviour; it shares no path with T-1205 and runs beside it.

- `design/system/workflow.md` § Branches and collisions: a **Shared by agreement** entry beside **Shared paths** and **Held**, and the request flai serve sends.
- `design/system/agent-coordination.md` § Decision: the ADR as a refinement.
- `work-management.md` in both copies: when asked about a hold, narrow your touches if you will not change the paths, share them with a split if the work divides, or say why the hold stands; when started on a share, keep to the split. `template/CHANGELOG.md` records it.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai message share` and `message_share`; regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- Each document names the share and links the ADR.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
