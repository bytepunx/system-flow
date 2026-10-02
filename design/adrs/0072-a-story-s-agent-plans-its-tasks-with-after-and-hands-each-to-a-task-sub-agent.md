---
id: ADR-0072
title: "A story's agent plans its tasks with after and hands each to a task sub-agent, and runs a layer's tasks at once only when they are long"
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0059]
---

# ADR-0072 A story's agent plans its tasks with after and hands each to a task sub-agent, and runs a layer's tasks at once only when they are long

## Context

S-0176, an experiment (ADR-0066), gave tasks `after`, the tasks of their story they wait for, showed the plan's layers in `flai show`, the board, and the story's page, and had the story's agent hand each task to a task sub-agent and run each layer's tasks at once. Its results document, [S-0176's results](../experiments/S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md), set three replays against their original runs. No story reached review sooner: the tasks took three to nine minutes, and most of a run is the story's agent's own priming, reviewing, verifying, and closing out. Cost fell 15 and 28% in two replays and rose 9% in the third, because the story's agent's context stayed small while sub-agents did the tasks' reading and editing. It recommended adapting: keep `after`, the plan's display, and the hand-off, and drop the expectation that layers at once make a story faster.

Two things stood in the way, both now done: the host's flai must carry `after` (flai 1.28.0, with `flai.minimum` 1.28.0), and per-task usage had to stop giving each of two tasks in progress at once the whole window (ADR-0071).

## Decision

The story's agent plans its tasks as it writes them, each with the paths it touches and the tasks of the story it waits for (`after`), and records the layers and why each task waits in the narrative's `## Decisions`. It works the plan layer by layer, handing each task to a task sub-agent whose description names the task's ID; it reviews, commits, and moves each task itself, and alone talks to the designer. Running a layer's tasks at once is the agent's choice, worth it only when the tasks are long beside the story's fixed costs. `work-management.md`, `delegation.md`, and the `claude-code` prompt say so, in this project and the template.

## Consequences

- A story's order of work is visible to the agent, the board, and the designer before it starts.
- A flai older than 1.28.0 refuses a task that carries `after`; `flai.minimum` keeps it off the host.
- Each task sub-agent starts with only its prompt and reads the code again, since headless sessions are not offered forks; S-0241 measures again when they are.
- A task's usage is its sub-agents' calls and a share of the agent's own (ADR-0071).

## Alternatives considered

- **Require the layers at once, as S-0176 tried.** It did not make stories of three to five short tasks faster, and a shared worktree lets two tasks overwrite each other's edits when their `touches` are wrong.
- **Drop the plan and the hand-off with the parallel layers.** It gives up the cost saving, which came from the hand-off, and the plan, which costs nothing.
