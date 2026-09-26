---
id: S-0128
type: story
nature: feature
title: flai serve and wait_for_work hold a ready story whose claim overlaps an open story's, and say why
status: backlog
parent: E-0009
owner: alex
created: 2026-09-26T07:59:21Z
updated: 2026-09-26T07:59:21Z
transitions: []
tags: [cli]
touches: [flai/internal/workitem, flai/internal/serve, flai/internal/mcpserver, flai/cmd, docs, design/system]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0128 flai serve and wait_for_work hold a ready story whose claim overlaps an open story's, and say why

## Goal

Build the hold of [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) in flai: a ready story whose claim overlaps the claim of a story in progress or in review is not started by flai serve and not offered by `wait_for_work`, and flai says why and what clears it. The dashboard's yellow card is the next story.

## Acceptance criteria

- [ ] A story's claim is its `touches` plus those of its tasks that are not done or cancelled, with a sub-project's name or tag from `system-flow.yaml` read as its path; two claims overlap by the `/`-prefix rule `flai check` uses, and an empty claim overlaps every claim. Behaviour tests cover each case.
- [ ] The flai serve launcher does not start a held story's agent and starts the next ready story in pull order that is not held, within the in-progress limit; the held story keeps its place and is started first once clear.
- [ ] `wait_for_work` offers the first ready story that is not held, and when every ready story is held it answers waiting with each held story and its reason.
- [ ] `agent.status` reports a held ready story as waiting with a reason of the form `held (overlap): touches <path>, inside <path> which S-nnnn (in progress) touches; starts when S-nnnn is accepted, cancelled, or sent back`, and `held (no-touches): …` for an empty claim, whether or not the story ever had an agent.
- [ ] `flai board` (text and `--json`), and the MCP `board` and `inbox` ready lists, mark a held story with `held` and its reason.
- [ ] `flai move <story> in-progress` and MCP `item_move` warn on a held story and still move it; `flai serve agent start` still starts a held story's agent.
- [ ] `design/system/workflow.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` describe the hold as built; `make test` and lint pass.

## Tasks

## Notes

- Overlap is `pathsOverlap` in `flai/internal/check/check.go`; the launcher's classification is `launcher.look` in `flai/internal/serve/agents.go`, whose per-story skip reasons are kept only in memory today and need to reach `Activity()`; `wait_for_work` is `flai/internal/mcpserver/work.go`.
- The `after:` field (ADR-0046) is a separate story; leave room for a second hold reason.
- The survey is `design/system/agent-coordination.md`.
