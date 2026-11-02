---
id: T-2221
type: task
nature: remediation
title: The container adapter's Start mounts the main checkout's .git read-only and the story's registration read-write from storygit.GitMounts, and its argv test asserts nothing else under .git is writable
status: done
parent: S-0517
owner: alex
created: 2026-10-28T13:20:55Z
updated: 2026-11-02T18:39:10Z
transitions:
  - to: ready
    at: 2026-11-02T18:28:22Z
    by: agent-S-0517
  - to: in-progress
    at: 2026-11-02T18:28:23Z
    by: agent-S-0517
  - to: done
    at: 2026-11-02T18:39:09Z
    by: agent-S-0517
stream: S-0517
tags: [cli, docker]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-2219]
usage:
  source: log
  seconds: 646
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 16606
      cache_read: 3606948
      cache_write: 107363
      cost: 1.7605
---
# T-2221 The container adapter's Start mounts the main checkout's .git read-only and the story's registration read-write from storygit.GitMounts, and its argv test asserts nothing else under .git is writable

## Work

In `flai/internal/harness/container.go`, the adapter S-0508 wrote, have `Start` call `storygit.GitMounts` on the story's worktree and add its mounts to the `docker run` argv after the worktree's: `-v <common>:<common>:ro` for the main checkout's `.git`, then `-v <registration>:<registration>` read-write for `.git/worktrees/S-nnnn`, the shorter destination first so Docker nests the registration over the read-only parent. An error from the helper fails the start with its message, as the adapter refuses a missing image or login, so `flai serve` reports a worktree whose registration is absent rather than starting a container without git. A normal checkout, which the helper answers with no mounts, adds none. Keep every other rule of S-0508's `Start`: never the Docker socket, the home, `~/.flai`, or the main checkout's files, and the keys and token by `-e NAME` only.

In `container_test.go`, extend the argv test: on a fixture whose worktree is a real git worktree (or a `.git` file and registration the test writes, with `commondir`), the argv holds the common dir with `:ro` and the registration without, in that order, and no other `-v` whose source is under `.git`; the count of mounts the test asserts goes from three to five, and the test's name and words say why. Add the failing case: a worktree whose registration is absent makes `Start` return the helper's error, and the argv is empty. Update any sentence in the package's comments or `design/system/agent-adapters.md` S-0508 wrote that says "the three mounts" where T-2222 does not already own it.

## Done when

- `Start` adds the two git mounts, read-only and read-write, after the worktree's, and fails with the helper's error when the registration is absent
- `container_test.go` asserts the two mounts' sources, destinations, modes, and order, that nothing else under `.git` is mounted, and the failure, and `flai test flai/internal/harness` passes

## Notes

Layer 3: waits for T-2219, whose helper it calls. The adapter is S-0508's, so this story names S-0508 in `after`; if S-0508 put the argv elsewhere than `container.go`, change it where it is and widen the touches, saying so in the narrative.
