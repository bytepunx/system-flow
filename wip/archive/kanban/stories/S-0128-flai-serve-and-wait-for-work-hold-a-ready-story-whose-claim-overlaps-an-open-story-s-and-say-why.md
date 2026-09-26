---
id: S-0128
type: story
nature: feature
title: flai serve and wait_for_work hold a ready story whose claim overlaps an open story's, and say why
status: done
parent: E-0009
owner: alex
created: 2026-09-26T07:59:21Z
updated: 2026-09-26T17:47:43Z
transitions:
  - to: ready
    at: 2026-09-26T08:03:50Z
    by: alex
  - to: in-progress
    at: 2026-09-26T08:04:04Z
    by: agent-S-0128
  - to: review
    at: 2026-09-26T08:17:46Z
    by: agent-S-0128
  - to: done
    at: 2026-09-26T17:47:43Z
    by: alex
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

- [x] A story's claim is its `touches` plus those of its tasks that are not done or cancelled, with a sub-project's name or tag from `system-flow.yaml` read as its path; two claims overlap by the `/`-prefix rule `flai check` uses, and an empty claim overlaps every claim. Behaviour tests cover each case.
- [x] The flai serve launcher does not start a held story's agent and starts the next ready story in pull order that is not held, within the in-progress limit; the held story keeps its place and is started first once clear.
- [x] `wait_for_work` offers the first ready story that is not held, and when every ready story is held it answers waiting with each held story and its reason.
- [x] `agent.status` reports a held ready story as waiting with a reason of the form `held (overlap): touches <path>, inside <path> which S-nnnn (in progress) touches; starts when S-nnnn is accepted, cancelled, or sent back`, and `held (no-touches): …` for an empty claim, whether or not the story ever had an agent.
- [x] `flai board` (text and `--json`), and the MCP `board` and `inbox` ready lists, mark a held story with `held` and its reason.
- [x] `flai move <story> in-progress` and MCP `item_move` warn on a held story and still move it; `flai serve agent start` still starts a held story's agent.
- [x] `design/system/workflow.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` describe the hold as built; `make test` and lint pass.

## Tasks
- T-0464 Compute each ready story's hold from claims, and mark held cards on the board
- T-0465 flai serve's launcher skips held stories and agent.status says why
- T-0466 wait_for_work offers the first story not held, and moves warn on a held story
- T-0467 Describe the hold in the design and user docs

## Notes

- Overlap is `pathsOverlap` in `flai/internal/check/check.go`; the launcher's classification is `launcher.look` in `flai/internal/serve/agents.go`, whose per-story skip reasons are kept only in memory today and need to reach `Activity()`; `wait_for_work` is `flai/internal/mcpserver/work.go`.
- The `after:` field (ADR-0046) is a separate story; leave room for a second hold reason.
- The survey is `design/system/agent-coordination.md`.
- Lint, as verified: golangci-lint v2.5.0 reports 0 issues, gofmt and go vet are clean. `make lint-md` reports six findings, all in files S-0128 did not write (archived S-0122 and S-0124 files, TH-0012), which TH-0017 covers.
