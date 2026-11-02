---
id: ADR-0242
title: "A story's agent whose harness is none runs elsewhere, and flai serve starts nothing for it"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0108, ADR-0130]
topics: [harness, flai]
---

# ADR-0242 A story's agent whose harness is none runs elsewhere, and flai serve starts nothing for it

## Context

I-0214 records two instances on S-0505 (2026-10-28). S-0513 was made to be worked by hand in the container image `flai-agent-claude-code:local`, and its notes said so. When it entered ready, `flai serve` started a host agent for it within a second. That agent read the notes and stood down without a change. When it ended with the story in progress, `flai serve` started a second as automatic restart 1 of 2 (ADR-0108), while the container agent was working the story; it stood down too.

Nothing could have stopped either start. A story's agent, its harness, model, and config (ADR-0037), is the only place that says how its agent is started, and only `flai serve`'s adapters read it (ADR-0130). `harness.For` knows two names, `claude-code` and `command`, and starts a story that names none with the host's command when one is set. The story's notes are prose no start path reads. The container agent could not pull the story itself either: `flai serve` takes a ready story first, and `flai guard` refuses one story's agent another's name (ADR-0143).

## Decision

A reserved harness name, `none`, on a story's agent says that the story's agent runs elsewhere, by hand or in a container nobody has made a harness for yet, and `flai serve` starts nothing for it.

- `flai serve` starts no agent for such a story: not when it enters ready or sits there, not as an automatic restart when a run ends with it in progress (ADR-0108), not on a resume, and not when it is sent back by the operator, the merge queue, or an acceptance. It names the story in `waiting` with the reason, `its agent runs elsewhere (harness none)`, and counts that as no failure. It opens no restart-limit thread for it.
- `flai serve agent start`, `restart`, and `commit`, and the dashboard's writes that reach them, refuse such a story, saying that its agent runs elsewhere, that `flai serve` starts nothing for it, and to start it with `flai story start` or name a harness `flai serve` can start.
- `harness.For` refuses `none` with the same words, so that no path above the adapters starts it by mistake. `none` is no adapter: it is not in `Adapters`, and the unknown-harness refusal does not list it.
- `flai story start` and the MCP tool `story_start` work as they do, since they start no agent: the agent that runs elsewhere starts the story with them.
- The manifest's shape check accepts `none` as it accepts any name of the right shape, so `flai edit --harness none` and `flai agent set --harness none` set it today. A project's default agent may name it, and every story made then carries it.

## Consequences

- A story meant for an agent `flai serve` cannot start is said to be so on the story itself, where `flai serve` reads it, and stays alone through ready, in progress, and every send back.
- The operator sees why such a story has no agent in `waiting`, on the story's page, and in the refusal of a start they ask for.
- Once S-0508 makes the container a harness `flai serve` starts, a story worked in a container names that harness instead, and `none` stays as the way to say "by hand".
- A strategic agent whose harness is `none` fails its start at `harness.For` with the refusal's words; no role is expected to carry it, so no strategic start path checks it apart.

## Alternatives considered

- A new field on the story, such as `agent.elsewhere: true`, beside its harness. It says the same thing in two places that can disagree, a harness that `flai serve` can start beside a flag that says not to, and it changes the work item schema.
- A host setting listing the stories `flai serve` must leave alone. It says on the host what is a fact about the story, is lost on another host, and is not seen on the story's page.
- Waiting for S-0508's container harness. It removes the container case only, leaves the hand-run case, and S-0508 is in the backlog.
