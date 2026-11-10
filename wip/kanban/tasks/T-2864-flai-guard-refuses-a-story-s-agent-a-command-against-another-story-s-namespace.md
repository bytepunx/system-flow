---
id: T-2864
type: task
nature: feature
title: flai guard refuses a story's agent a command against another story's namespace or the default namespace
status: backlog
parent: S-0661
owner: alex
created: 2026-11-10T08:23:35Z
updated: 2026-11-10T08:23:35Z
transitions: []
stream: S-0661
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go]
---
# T-2864 flai guard refuses a story's agent a command against another story's namespace or the default namespace

## Work

Read first the section of the ADR S-0659 accepted on what `flai guard` can tell of an agent's command against a namespace. Build the rule it decides there, no wider.

Beside `klusterWrites` in `flai/internal/guard/guard_rules.go`, add a rule for a session `flai serve` started for a story, and its sub-agents. It refuses a command that names a namespace other than the story's own, `s-nnnn` from `kluster.StoryNamespace`, or the block's default namespace. Such a command names it with `kubectl` or `helm` and `-n`, `--namespace`, or `--all-namespaces`, or sets `FLAI_KLUSTER_NAMESPACE` or `KUBECONFIG` to another value. Hook the rule where `refuse` and `script` judge a call. The refusal names the namespace, says the agent acts only on its own story's, and says what to run instead.

The operator's shell and the orchestrator are not refused. A command that names no namespace passes: the namespace's variables point it at the story's own.

Tests in `kluster_test.go` cover: the own namespace passing, and another story's and the default one refused, each by `-n`, `--namespace=`, and an environment assignment. They also cover `--all-namespaces` refused, and the operator not refused. It waits for no task: it touches only the guard, and is the first layer.

## Done when

- `flai guard` refuses a story's agent's `kubectl` or `helm` command against another story's namespace or the default one, with the reason and what to do, and passes one against its own.
- `go test ./internal/guard/...` passes.

## Notes
