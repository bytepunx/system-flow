---
id: T-3279
type: task
nature: feature
title: flai kluster exec runs a command in the story's worktree with its namespace's variables and its own kubeconfig while the namespace has a record
status: done
parent: S-0763
owner: alex
created: 2026-11-13T09:03:28Z
updated: 2026-11-13T12:06:16Z
transitions:
  - to: ready
    at: 2026-11-13T11:47:19Z
    by: agent-S-0763
  - to: in-progress
    at: 2026-11-13T11:47:20Z
    by: agent-S-0763
  - to: done
    at: 2026-11-13T12:06:16Z
    by: agent-S-0763
stream: S-0763
tags: [cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3276, T-3277]
usage:
  source: log
  seconds: 1136
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 11572
      cache_read: 2351864
      cache_write: 89494
      cost: 1.2566
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 15
          cache_read: 81709
          cache_write: 2328
          cost: 0.0242
---
# T-3279 flai kluster exec runs a command in the story's worktree with its namespace's variables and its own kubeconfig while the namespace has a record

## Work

Criterion 3. In `flai/cmd/kluster.go`, add `flai kluster exec S-nnnn -- <command> [args...]`:

- It reads the story's namespace's record, and refuses when there is none, saying to deploy first with `flai kluster deploy S-nnnn --job`. A failed record is no refusal: the agent diagnoses with it.
- It runs the command as an argument list, never through a shell, in the story's worktree, with flai's own environment and `Cluster.AgentEnv` over it, so `FLAI_KLUSTER_IMAGES` is absent and `KUBECONFIG` is `.flai-cache/kluster/agents/s-nnnn.kubeconfig`. It passes standard input, output, and error through, and exits with the command's exit code.
- It refuses a story with no worktree, naming `flai stream open`, a project with no `kluster:` block, and, for a session with a story of its own, another story, as `--job` does (`sessionStory`).
- When the namespace's kubeconfig is absent, as after a deploy that failed before its kubeconfig step, it says so on standard error and runs the command anyway, since a project's own script may not need `kubectl`.
- It needs no flai serve and writes nothing: what the command does is recorded by the agent, not by flai.

Tests in `kluster_test.go`, with a faked record and worktree and a command that prints its environment and working folder: the variables, `KUBECONFIG` once, no `FLAI_KLUSTER_IMAGES`, the worktree as its folder, the exit code passed through, a failed record allowed, and each refusal. Run `scripts/flai-reference.sh` again for the reference and the flag index.

Waits for T-3276, which gives `Cluster.AgentEnv` and the kubeconfig's path, and for T-3277, which changes `flai/cmd/kluster.go`, `kluster_test.go`, and the generated reference in the layer before.

## Done when

- `flai kluster exec S-nnnn -- env` prints the namespace's variables, with `KUBECONFIG` the namespace's kubeconfig and no `FLAI_KLUSTER_IMAGES`, run in the story's worktree, while the namespace has a record, a failed one included.
- It exits with the command's exit code and refuses a namespace with no record, saying to deploy first.
- `go test ./cmd/` passes, `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` among them.

## Notes
