---
id: ADR-0089
title: "Acceptance criteria are ticked through flai by number, by the story's agent once it has verified each"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
topics: [all]
---

# ADR-0089 Acceptance criteria are ticked through flai by number, by the story's agent once it has verified each

## Context

A story's acceptance criteria are checkboxes under its `## Acceptance criteria` heading. `flai move` refuses a story to done while one is unticked, and the operator reads the ticks at review as what the agent verified. Until S-0282 nothing in flai ticked a box: an agent hand-edited the item's file, which meant knowing the syntax and touching a file flai owns the rest of, or left the boxes unticked and the operator ticked them at acceptance. Ticks were inconsistent between agents, and an unticked box at review said nothing about whether the criterion had been checked. Task sub-agents do much of a story's work, but the guard refuses them every flai write (ADR-0059, ADR-0060), and work-management.md says never to tick a criterion you did not verify.

## Decision

A criterion is ticked through flai, by the story's agent, once it has verified it.

1. **One change, three ways in.** `flai criteria tick|untick <id> <n>...` on the CLI, the MCP tool `criteria_tick`, and the host API action `item.criteria` behind the dashboard's `POST /api/items/[id]/criteria` all make the same change: the boxes named, and no other byte of the item. `flai criteria list <id>` prints them. Each is an edit of the item's body through `itemedit`, with `flai edit`'s checks, its refusal of an archived or closed item, its hash guard against a change made meanwhile, and its notice to agents that the item's criteria changed.
2. **Numbers.** A criterion is named by its number, from 1, in the order of the boxes under `## Acceptance criteria`, nested ones included, as `flai criteria list` prints them. A number with no box is refused and nothing is written.
3. **Who ticks.** The story's agent, after verifying the criterion, as work-management.md requires (TH-0162). A task sub-agent says in its final message which criteria its task meets; the story's agent reviews the work and ticks them. A verifier says which criteria the diff meets. The guard lets a sub-agent run `flai criteria list`, a read, and refuses it `tick`, `untick`, and `criteria_tick`, as it refuses every other write. The planner and the orchestrator are refused them too.
4. **At review.** `flai move` to review warns, naming the count, when a story has an unticked criterion, and does not refuse: a criterion that cannot be verified stays unticked with the reason in the story's notes. The move to done still refuses one.

## Consequences

- Every agent ticks the same way, and a tick at review means the story's agent verified the criterion.
- An agent no longer hand-edits an item's file to tick a box, so the edit is checked and the other agents are told.
- A task sub-agent cannot tick, so the story's agent spends a call per tick, or one call for several.
- Numbers shift when someone adds or removes a criterion between the read and the tick; the hash guard refuses that tick instead of ticking the wrong box.
- The dashboard has the route but no checkbox on the story page yet; the operator still ticks at acceptance by editing the body or through the route.

## Alternatives considered

- **Let a task sub-agent tick its own story's criteria.** It is the one that did the work, but the guard would have to tell a task sub-agent from an explorer or a verifier, and the tick would come before the story's agent's review, which is where a criterion is verified (TH-0162).
- **Refuse the move to review with an unticked criterion.** It would force a tick on a criterion that could not be verified in the agent's environment, which work-management.md tells it to leave unticked and explain.
- **Name criteria by their text.** Text is long, may repeat, and changes with an edit; a number with the hash guard is short and safe.
- **One `flai criteria <id> --tick n --untick m` command.** The guard tells a read from a write by subcommand, as it does `flai thread list` from `flai thread reply`; a flag that makes the same command write is a form it has to parse.
