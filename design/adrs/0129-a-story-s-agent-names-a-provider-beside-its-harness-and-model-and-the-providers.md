---
id: ADR-0129
title: "A story's agent names a provider beside its harness and model, and the providers map, with each gateway's API, base URL, and the name of its key's variable, lives on the host"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0037, ADR-0038]
---

# ADR-0129 A story's agent names a provider beside its harness and model, and the providers map, with each gateway's API, base URL, and the name of its key's variable, lives on the host

## Context

A story's `agent` is a harness, a model, and a config map, copied from the project's default when the story is made ([ADR-0037](0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)), and `flai serve` starts it through the harness's adapter with the program and arguments the operator set on the host ([ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md)). Where the model's calls go is nowhere in flai: `claude` sends them to Anthropic unless the environment says otherwise, and the model name is passed through unread. Over a gateway, a LiteLLM proxy or OpenRouter, the three come apart: the same harness, another provider, and a model name that means something only on that provider. The finding of S-0339, [agent-adapters.md](../system/agent-adapters.md), sets this out under Abstractions.

## Decision

A story's agent names a provider beside its harness and model, and the providers map, with each gateway's API, base URL, and the name of its key's variable, lives on the host.

- `agent` gains `provider`, beside `harness`, `model`, `config`, and `roles`, in the manifest's default agent, under `planning`, `orchestration`, and `analysis`, and in a story's front matter, copied as ADR-0037 copies the rest. Absent, the harness's own default stands: today's behaviour.
- `provider` names an entry in a `providers` map on the host, in flai's configuration beside `agent.harnesses`, set with `flai serve agent provider <name>` and the dashboard's settings host action. An entry holds `api` (`anthropic-messages`, `openai-chat`, or `openai-responses`), `base_url`, `key_env`, the name of the environment variable that holds the key, and optional `models`, the harness's aliases on that provider. It never holds a key.
- `flai serve` reads the variable `key_env` names when it starts the agent, sets the harness's own variable in the child's environment, and logs neither; the start record keeps the name. It refuses a start whose provider names an unset variable, or an `api` the harness cannot speak, as it refuses an unknown harness.
- Each adapter derives its own settings from the entry: `claude-code` sets `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`, an empty `ANTHROPIC_API_KEY`, and the `ANTHROPIC_DEFAULT_*_MODEL` aliases; the `command` adapter gets `FLAI_PROVIDER`, `FLAI_PROVIDER_API`, `FLAI_PROVIDER_BASE_URL`, `FLAI_PROVIDER_KEY_ENV`, and `{provider}`.

## Consequences

- The manifest, the story schema, the host configuration, the host API's settings, the dashboard's agent fields, `docs/operators/settings.md` and its index test, and `docs/users/flai.md` change together in the story that builds it.
- A clone carries no endpoint and no key name: a project moved to another host names its providers again there, as it names its harness programs.
- A model name's meaning is the provider's; flai keeps validating it by pattern only.

## Alternatives considered

- A `providers` map in the manifest, committed, with a host override: a clone would carry it, but base URLs and variable names are host-bound, and anyone who edits the manifest could point a story at another endpoint, against ADR-0038's rule that what runs is the operator's.
- Provider settings inside `agent.config`: three keys per story, repeated, with the harness mapping each; the harness, not the story, knows what to set.
- The key itself in the host configuration: refused by the safety convention; the environment and the name of its variable suffice.
