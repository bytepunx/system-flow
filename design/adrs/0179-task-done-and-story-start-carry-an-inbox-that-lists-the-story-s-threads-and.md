---
id: ADR-0179
title: "task_done and story_start carry an inbox that lists the story's threads and counts the rest, and each answer fits 40,000 bytes of JSON"
status: accepted
date: 2026-10-20
supersedes: []
superseded_by: []
refines: [ADR-0104, ADR-0107]
topics: [mcp, inbox, flai]
---

# ADR-0179 task_done and story_start carry an inbox that lists the story's threads and counts the rest, and each answer fits 40,000 bytes of JSON

## Context

`task_done` and `story_start`, over MCP, and `flai task done --json` and `flai story start --json` each carry the agent's inbox as the MCP tool `inbox` answers it ([ADR-0107](0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md)). That inbox lists every open thread awaiting the agent, the operator's threads on other items among them. Claude Code saves a tool result over 50,000 characters to a file and shows only a preview ([ADR-0104](0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)).

[I-0167](../issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md) records four stories, S-0398, S-0399, S-0401, and S-0404, whose `story_start` answered 87,106 to 101,139 characters and whose `task_done` answered 83,000 to 84,000. Their inboxes carried 31 to 33 threads awaiting the operator and none on the story. Claude Code saved each answer to a file, and the agent read it back with `jq`. S-0450's own `story_start` answered 80,220 characters: the inbox 70,676 bytes, about 43,000 of them threads and about 27,000 changes.

`story_start`'s fit (`flai/internal/mcpserver/story_start.go`) trims only the inbox's changes and the pack. `task_done` and both CLI `--json` forms have no fit.

## Decision

The inbox that `task_done` and `story_start` carry lists in full only the open threads on the story and its tasks, counts the rest in `threads_omitted`, and each whole answer fits 40,000 bytes of JSON.

- **The story's threads.** The inbox in these answers is read with the inbox's story filter set to the story. It lists in full the open threads on the story and its tasks, awaiting the agent as `inbox` lists them, and `awaiting_you` counts those listed. The open threads on other items that the inbox would otherwise have listed are counted, not listed, in a new field `threads_omitted`: threads on other items left out; call `inbox` to list them. `messages`, `ready`, `can_pull`, and `changes` are as before.
- **The standalone inbox.** The MCP tool `inbox` and `flai thread list` answer as before. `threads_omitted` is omitted when zero, and only a story-scoped read sets it.
- **The fit.** Each whole answer, over MCP and from the CLI's `--json` alike, is fitted to 40,000 bytes of JSON, `context.PartLimit`, the limit ADR-0104 sets for `prime`'s parts. First the inbox's changes are cut from the oldest, counted in `changes_omitted` as now. For `story_start`, the pack is part 1 when the answer fits with it, else only its header, as now. Then, if the answer is still over, the last entry of each listed thread, and its pending recommendation's text, is shortened to a fixed length with an ellipsis, counted in a new field `entries_shortened`, and the agent reads them whole with `thread_get`.

This refines ADR-0104, whose limit now bounds these answers too, and ADR-0107, whose inbox step answers the story's inbox rather than the agent's whole one.

## Consequences

- `task_done` and `story_start` answer within Claude Code's limit, and the agent reads them without `jq`. A test reproduces I-0167's instances: an inbox with many threads awaiting the operator on other items and none on the story.
- A story's agent no longer sees the operator's threads on other items in these answers. It answers only its own story's threads anyway, and calls `inbox` for the rest when `threads_omitted` says there are some.
- A shortened entry is read whole with `thread_get`. `entries_shortened` says how many were shortened.
- The CLI's `--json` and the MCP answer stay the same result, built by `internal/taskdone` and `internal/storystart`, so the host channel's `task.done` and `story.start` get the fit too.

## Alternatives considered

- **Leave the inbox out of these answers.** Rejected: the story's own threads, its messages, and the changes are what the agent needs at a task transition.
- **Cap the number of threads listed.** Rejected: the cap would cut the story's own threads as readily as the operator's on other items.
- **Fit only the MCP answers.** Rejected: `flai task done --json` and `flai story start --json` are the same result, and an agent reading them through Bash meets the same limit.
