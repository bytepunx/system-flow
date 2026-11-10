---
id: T-2951
type: task
nature: remediation
title: The design, the user guide, and tooling.md in both copies say stream state writes Context, Decisions, and Open questions, and I-0277 is closed
status: done
parent: S-0668
owner: alex
created: 2026-11-10T08:51:10Z
updated: 2026-11-10T10:22:02Z
transitions:
  - to: ready
    at: 2026-11-10T10:13:08Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:13:09Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:22:02Z
    by: agent-S-0668
stream: S-0668
tags: [docs, template, narrative]
touches: [design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md, design/issues/I-0277-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md, design/issues/summary.md]
after: [T-2947, T-2949]
usage:
  source: log
  seconds: 496
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 19
      output: 208
      cache_read: 2671449
      cache_write: 18415
      cost: 0.9901
---
# T-2951 The design, the user guide, and tooling.md in both copies say stream state writes Context, Decisions, and Open questions, and I-0277 is closed

## Work

Write down the behaviour T-2944, T-2947, and T-2949 built, in the words their commands and tool now use:

- `design/system/agent-narrative.md` § Narrative structure: the paragraph after the example says `## Context`, `## Decisions`, and `## Open questions` are written by hand. Say that `flai stream state`, `stream_state`, and `stream.state` replace Context and append Decisions and Open questions bullets, keeping the threads block. Say a question so written holds the move to review, as a hand-written one does. Say that this is how an agent whose `wip/` is read-only, in a container (ADR-0253), records them. Bump `updated`.
- `design/system/flai-cli.md`: the `flai stream state` row of the command table, and `stream_state` in the `flai mcp` tool list.
- `docs/users/flai.md`: § Agent narratives, the paragraph that ends "are still edited by hand", and the `stream_state` row of the MCP table.
- `design/conventions/tooling.md` and `template/root/design/conventions/tooling.md`, line 17: name the new flags, and say the sections may still be edited by hand where `wip/` is writable. Change the two copies alike. Add one line under `## Unreleased` in `template/CHANGELOG.md`.
- Close I-0277 with `flai issue close I-0277 --reason "..."`, naming the widened `stream_state` and `flai stream state` and T-2949's test. That writes the issue file and `design/issues/summary.md`.

Waits for T-2947 and T-2949, whose flags and fields it describes.

## Done when

- `flai test` on the changed files passes, the markdown lint included.
- `flai check --strict` passes.
- I-0277 reads `status: closed` with the reason.

## Notes
