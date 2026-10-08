---
id: T-1310
type: task
nature: remediation
title: flai serve lets a project's requests in flight finish when its manifest stops loading
status: done
parent: S-0321
owner: alex
created: 2026-10-08T00:23:22Z
updated: 2026-10-08T08:07:05Z
transitions:
  - to: ready
    at: 2026-10-08T07:58:42Z
    by: agent-S-0321
  - to: in-progress
    at: 2026-10-08T07:58:43Z
    by: agent-S-0321
  - to: done
    at: 2026-10-08T08:07:05Z
    by: agent-S-0321
stream: S-0321
tags: [flai, serve]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/channel/channel.go]
usage:
  source: log
  seconds: 502
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 60
      output: 262
      cache_read: 2723652
      cache_write: 133151
      cost: 1.2806
---
# T-1310 flai serve lets a project's requests in flight finish when its manifest stops loading

## Work

When `reconcile` in `flai/internal/serve/serve.go` finds a served project unavailable, it calls `running.halt`, which cancels the context the project's `channel.Client` runs under (`flai/internal/channel/channel.go`). That kills every request in flight. On 2026-10-07 it killed the `publish.run` that had just raised the minimum, which the journal recorded as failed with `flai exited with -1` although the tags had been pushed (I-0107).

A project dropped because its manifest stopped loading, such as a `manifest.TooOldError`, should stop taking new requests and let the ones in flight finish, within a bound, before its client stops. Each finished request reports its real exit. A project removed from the registry, or served under another key, still stops at once, as it does now.

Waits for nothing: it shares no path with the release task, so the two run together.

## Done when

- A test in `flai/internal/serve/serve_test.go` writes a `flai.minimum` above the running flai into a served project's `system-flow.yaml` while a request runs. The request finishes and answers its real result, and then the project is dropped with the reason logged once.
- A project removed from the registry still stops its requests at once, as `TestServesRegisteredProjectsAndFollowsTheRegistry` expects.
- `flai test` passes on the changed paths.

## Notes
