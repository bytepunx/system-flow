---
id: T-1026
type: task
nature: improvement
title: The planner and the story's agent are told to declare file-level touches and to narrow a folder touch once tasks name its files
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:15:15Z
updated: 2026-10-06T12:15:38Z
transitions: []
stream: S-0295
tags: [flai, template]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/CHANGELOG.md]
after: [T-1023]
---
# T-1026 The planner and the story's agent are told to declare file-level touches and to narrow a folder touch once tasks name its files

## Work

A first guess at a story's touches claims whole folders (I-0087's S-0285). Tell the agents that write touches to name files where they can.

- **Planner.** In `strategic-agents.md` § As the planner (project and template copies), the planner predicts file-level touches for a story and its tasks. It keeps a folder touch only where the story may add files there that no task can name yet. It records which folder touches it kept, and why, under `### Planning`.
- **Story's agent.** In `work-management.md` (both copies), when it writes or reviews tasks, the story's agent narrows a folder touch of the story to the files its tasks name, or leaves it for T-1027's rule to narrow. It never declares a shared path to escape a hold, since shared paths are the operator's list.
- **Prompt.** Bring the planner's prompt in `flai/internal/harness/harness.go`, which names `flai touches suggest`, into line, and update its test.
- **Changelog.** Add a `template/CHANGELOG.md` entry for the convention change.

Baseline text changes in `template/` and in the project's copy together, as `CLAUDE.md` says. This task waits for T-1023, whose ADR the conventions cite.

## Done when

- Both conventions say what is above in the project and template copies, and the copies agree above the marker.
- The planner's prompt asks for file-level touches, and its test checks the sentence.
- `scripts/flai-test.sh`, the markdown lint, and `flai check --strict` pass.

## Notes
