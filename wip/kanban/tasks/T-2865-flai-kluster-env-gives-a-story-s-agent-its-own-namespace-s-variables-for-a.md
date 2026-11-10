---
id: T-2865
type: task
nature: feature
title: flai kluster env gives a story's agent its own namespace's variables for a command of its own, refusing one not deployed at the head
status: backlog
parent: S-0661
owner: alex
created: 2026-11-10T08:23:46Z
updated: 2026-11-10T08:23:46Z
transitions: []
stream: S-0661
tags: [cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go]
after: [T-2863]
---
# T-2865 flai kluster env gives a story's agent its own namespace's variables for a command of its own, refusing one not deployed at the head

## Work

Build the way the ADR S-0659 accepted gives the agent the namespace's variables for a command of its own. This plan assumes `flai kluster env S-nnnn [--json]`. It prints `export NAME=value` lines a shell can `eval`, or one JSON object with `--json`. If the ADR decides another form, such as `flai kluster exec S-nnnn -- <command>`, build that one and say so in the story's `## Decisions`.

In `flai/cmd/kluster.go`, add the command beside `deploy` and `deployments`. It reads the variables with `serve.KlusterEnv`, and refuses with the reason T-2863 gives it when the namespace is not deployed at the branch's head, saying to deploy first. It writes nothing and needs no flai serve, so a story's agent and the operator both run it. It refuses another story than the agent's own (`FLAI_STORY`) for a session `flai serve` started. It names `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE` among the variables, as `Cluster.Env` gives them.

Tests in `kluster_test.go`, with a faked record and branch: the variables at the head, `--json`, each refusal, and another story refused. It waits for T-2863, which gives `KlusterEnv` the reasons this command refuses with.

## Done when

- `flai kluster env S-nnnn` prints the story's namespace's variables, `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE` among them, when it is deployed at the head, and refuses otherwise, saying to deploy first.
- `go test ./cmd/ -run Kluster` passes.

## Notes
