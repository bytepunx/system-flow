---
id: T-3337
type: task
nature: remediation
title: The every-address test checks its listeners closed without counting on their ports staying free
status: done
parent: S-0779
owner: alex
created: 2026-11-13T21:30:45Z
updated: 2026-11-13T21:34:14Z
transitions:
  - to: ready
    at: 2026-11-13T21:30:57Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T21:30:57Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T21:34:14Z
    by: agent-S-0779
stream: S-0779
tags: []
touches: [flai/cmd/mcp_http_test.go]
usage:
  source: log
  seconds: 197
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1632
      cache_read: 706635
      cache_write: 8674
      cost: 0.2434
---

# T-3337 The every-address test checks its listeners closed without counting on their ports staying free

## Work

T-3333's load run found a third test of the same kind. In one of ten whole-package race runs beside the fifty-count run, `TestMCPOverHTTPListensOnEveryAddress` failed at `mcp_http_test.go:776`: `a stop closes every listener, but 127.0.0.1:40842 is still held: listen tcp 127.0.0.1:40842: bind: address already in use`. After the server stops, the test listens on each port it closed, counting on nobody else taking it meanwhile, which another test's connection or listener can.

Check that the server's listeners are gone without that count: an address is gone when it can be listened on or refuses a connection, and one that still accepts is asked again for a few seconds, since another test's listener may hold it for a moment. A listener the server failed to close accepts every time, so the test still catches that.

Waits for nothing in the story: T-3332, the last task on this file, is done.

## Done when

- `TestMCPOverHTTPListensOnEveryAddress` no longer fails when another process takes one of its ports after the server stops, and still fails when a listener stays open.
- `flai test` on `flai/cmd/mcp_http_test.go` passes.

## Notes
