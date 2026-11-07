---
id: T-2584
type: task
nature: remediation
title: The design and the user and operator docs say every story has a default agent and harness, and that host agent settings are per project
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:12:49Z
updated: 2026-11-07T11:35:09Z
transitions:
  - to: ready
    at: 2026-11-07T11:29:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T11:29:53Z
    by: system-flow
  - to: done
    at: 2026-11-07T11:35:09Z
    by: agent-S-0587
stream: S-0587
tags: [docs]
touches: [design/system/agent-adapters.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/template.md, design/system/strategic-agents.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/operators/runbooks/install.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, flai/cmd/serve_actions.go]
after: [T-2578, T-2579, T-2580, T-2581, T-2582, T-2583]
usage:
  source: log
  seconds: 316
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 30119
      cache_read: 19651559
      cache_write: 169950
      cost: 5.5689
---
# T-2584 The design and the user and operator docs say every story has a default agent and harness, and that host agent settings are per project

## Work

It waits for every code task, so the docs describe what shipped. In the same story, as `documentation.md` asks, rewrite where these docs say a project "has none until the operator sets one", or that a story with no command "is not started":

- `design/system/agent-adapters.md`;
- `design/system/project-manifest.md` (`agent:`);
- `design/system/flai-cli.md` (`flai agent`, `flai serve agent`);
- `design/system/template.md`;
- `docs/operators/index.md` § "Starting an agent when a story becomes ready";
- `docs/operators/runbooks/install.md`;
- `docs/users/flai.md`;
- `docs/users/flaiover.md` (Settings, Agents, Harnesses).

They should now say:

- the order a story's agent resolves in: the story's own agent, then the manifest's, then the built-in;
- that the template ships a default agent;
- that the host command and harness programs are per project with a host-wide fallback and `--all-projects`;
- that a refusal names the settings page.

Bump each document's `updated`.

## Done when

- No document under `design/system` or `docs` still says a new project has no default agent, or that the host command is the same for every project.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
