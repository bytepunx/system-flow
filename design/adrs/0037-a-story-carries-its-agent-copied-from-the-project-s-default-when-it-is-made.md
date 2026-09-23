---
id: ADR-0037
title: "A story carries its agent, copied from the project's default when it is made"
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: []
---

# ADR-0037 A story carries its agent, copied from the project's default when it is made

## Context

E-0008 has flai start an agent when a story reaches ready. S-0103 prepares for it: the operator asked for a default harness, model, and configuration that every story gets, and for each story to be able to say otherwise, from the CLI, MCP, and the dashboard. Where that lives decides what an older flai can still read, and what a change to the default does to stories already written.

## Decision

**An agent is a harness, a model, and options, kept in the manifest as the project's default and in a story's front matter as its own.** The operator set the shape in S-0103's criteria (2026-09-23).

- `system-flow.yaml` may hold `agent: { harness, model, config }`. `flai agent set` and `flai agent clear` write it, and nothing else in the file changes. `config` is a flat map of string keys to string values: options for the harness, such as `effort` or `max_turns`, which flai passes on and does not interpret.
- **A story made while a default is set gets a copy of it** in its front matter, with what `flai story new --harness/--model/--agent-config` or the dashboard's form gives laid over it. The copy is the story's: changing the default later changes new stories only.
- **A story's agent is one of its own words**, like tags and touches: `flai edit`, MCP's `item_edit`, and the dashboard's editor change it or remove it. An edit from the dashboard or MCP replaces the whole agent. `flai edit` merges the flags into the agent the story has, unless `--clear-agent` is given.
- Only a story carries an agent. An epic or a task with one is refused.
- A harness is a lower-case name, a model a name that may hold `.`, `:`, `/`, and `@`, and a config key a lower-case word. Values are one line. flai refuses anything else before writing, so no value can be read as a flag or break the YAML.
- **No key is written when there is nothing to write.** A project with no default gets stories with no `agent`. An older flai parses items strictly, so it goes on reading such a project. Once it has agents, the project needs the flai that knows them.

## Consequences

- S-0104's dispatcher reads the agent from the story alone. It needs neither the manifest nor the history of the default to know what to start.
- The front matter says who will work a story before anyone does. A reviewer sees the model a story was worked with in its file and its commits.
- The default is not a live setting. To move open stories to a new model, edit each one. This is on purpose: a story already in progress should not change model under its agent.
- An older flai refuses a repository that has any story with an agent, until it is upgraded. `flai self-upgrade` is the answer.

## Alternatives considered

- **Resolve the default when the agent starts, and store only what a story overrides.** Stories stay smaller, and a change to the default reaches every open story. But the criterion is that stories carry the default, and what a story will run would then depend on a file it does not name.
- **A separate file per story, or a table in the manifest keyed by story.** An older flai would read the items. But a story's agent would live away from the story: no rename, archive, or edit keeps it in step.
- **Agents as a list of named profiles, with stories naming one.** It is worth doing once there are several harnesses to choose between. A profile can be added later as one more key of the agent.
