---
id: T-2865
type: task
nature: feature
title: flai kluster env prints a story's agent its own namespace's variables, with the namespace's own kubeconfig, while the namespace has a record
status: done
parent: S-0661
owner: alex
created: 2026-11-10T08:23:46Z
updated: 2026-11-13T12:59:17Z
transitions:
  - to: ready
    at: 2026-11-13T12:42:10Z
    by: agent-S-0661
  - to: in-progress
    at: 2026-11-13T12:42:11Z
    by: agent-S-0661
  - to: done
    at: 2026-11-13T12:59:17Z
    by: agent-S-0661
stream: S-0661
tags: [cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1026
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 86
      output: 460
      cache_read: 5402589
      cache_write: 128390
      cost: 2.0662
---
# T-2865 flai kluster env prints a story's agent its own namespace's variables, with the namespace's own kubeconfig, while the namespace has a record

## Work

ADR-0292 §2 decides the form: `flai kluster env S-nnnn` prints the namespace's variables as `export` lines, or one JSON object with `--json`, and `flai kluster exec S-nnnn -- <command>` runs a command with them. S-0763 built `exec` on `Cluster.AgentEnv`, at the record's commit, with `KUBECONFIG` the namespace's own kubeconfig and without `FLAI_KLUSTER_IMAGES` and `FLAI_KLUSTER_COMPONENT`. `env` gives exactly those variables, from the same code.

In `flai/cmd/kluster.go`, add `env` beside `exec`:

- It acts for the session's own story alone, as `exec` and `log` do, and refuses another story than `FLAI_STORY` or the story in `agent-S-nnnn`; the operator's shell names any.
- Like `exec`, it is open while the namespace has a record, a failed one included, and refuses a story with no record, saying to deploy it with `flai kluster deploy S-nnnn --job`. It needs no flai serve, no host action, and no `kluster` on `PATH`. It is not held to the head: the ADR makes `env` and `exec` reads for diagnosis.
- `export NAME='value'` lines, each value quoted for POSIX `sh` and zsh, in a stable order, so that `eval "$(flai kluster env S-nnnn)"` sets them. `--json` prints one object of name to value.
- When the namespace has no kubeconfig yet, it says so on standard error, as `exec` does, and prints the rest.
- Share `exec`'s code for the story, the record, and the variables, rather than copying it.

Tests in `kluster_test.go`, beside `exec`'s, with a faked record: the variables, `KUBECONFIG` the namespace's own and `FLAI_KLUSTER_NAMESPACE` among them, `--json`, a value that needs quoting, no record refused, and another story refused. It waits for no task: it reads S-0763's `Cluster.AgentEnv`, not T-2863's work, and touches only these two files.

## Done when

- `flai kluster env S-nnnn` prints the story's namespace's variables, `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE` among them, as `exec` gives them, and refuses a story with no record or another story than the session's own.
- `flai test flai/cmd/kluster.go flai/cmd/kluster_test.go` passes.

## Notes
