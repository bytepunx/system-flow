---
id: T-3049
type: task
nature: improvement
title: roleStart mounts the main checkout's .git/config read-only, with tests that a strategic container's config is not writable
status: done
parent: S-0688
owner: alex
created: 2026-11-10T10:35:00Z
updated: 2026-11-10T19:43:24Z
transitions:
  - to: ready
    at: 2026-11-10T19:32:47Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:32:49Z
    by: system-flow
  - to: done
    at: 2026-11-10T19:43:24Z
    by: agent-S-0688
stream: S-0688
tags: [flai, container, security]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/container_test.go, flai/internal/serve/container.go]
after: [T-3045]
usage:
  source: log
  seconds: 635
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4374
      cache_read: 1734896
      cache_write: 51085
      cost: 0.8214
---
# T-3049 roleStart mounts the main checkout's .git/config read-only, with tests that a strategic container's config is not writable

## Work

In `roleStart` (`flai/internal/harness/container.go`), add `-v <root>/.git/config:<root>/.git/config:ro` beside the `.git/hooks` mount, and say so in its comment with the ADR T-3045 wrote. Refuse the start, as `mountable` refuses a missing path, when `.git/config` is not a regular file, so that a missing or replaced file is not mounted silently.

Reproduce the cause in the tests. Extend `TestAStrategicAgentsContainerMountsTheCheckout` (`flai/internal/harness/container_test.go`) to assert the `:ro` config mount, and add a case for the refusal. Extend `roleContainerStarted` (`flai/internal/serve/container_test.go`), which the planner, orchestrator, analyzer, and curator container start tests share, to check the config mount for every role. Where the test setup allows, add one test that runs the built `docker run` arguments' mount list against a checkout and shows a write to `.git/config` is refused, or say in the task's Notes why a mount-argument test is as far as it goes without Docker.

It waits for T-3045, whose ADR is the decision this implements.

## Done when

- A strategic container's `docker run` mounts `.git/config` read-only for all four roles, and the tests above pass with `flai test` on the changed files.
- A start whose `.git/config` is missing or not a regular file is refused with a reason naming it.

## Notes
