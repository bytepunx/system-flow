---
id: S-0661
type: story
nature: feature
title: A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:09Z
updated: 2026-11-10T08:09:09Z
transitions: []
tags: [cli]
after: [S-0660]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables

## Goal

A story's agent runs the project's checks against its story's namespace before review, and runs its own commands with the namespace's variables, so that it verifies a criterion that needs the whole system and ticks it.

## Acceptance criteria

- [ ] A story's agent runs the project's checks against its own story's namespace, deployed at its branch's head, with one flai command and one MCP tool, and gets each check's result and log
- [ ] The run gets the namespace's variables, `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE` among them, as a checks run in review does, and is refused when the namespace is not deployed at the head, saying to deploy first
- [ ] The agent gets the same variables for a command of its own, such as a request to a service in its namespace, in the way the ADR decides
- [ ] `flai guard` refuses an agent's command against another story's namespace or the default namespace, as far as the ADR says it can tell
- [ ] Checks runs against the cluster still run one at a time per project, after the job queue
- [ ] A failed run asked for by the agent opens no thread to the operator
- [ ] The run is recorded with its commit, its outcome, and who asked, where the next story reads it
- [ ] `design/system/kluster.md`, `docs/users/flai.md`, and its reference say all of it

## Tasks

## Notes

- `KlusterEnv` in `flai/internal/serve/kluster.go` decides the variables for every checks run today.
