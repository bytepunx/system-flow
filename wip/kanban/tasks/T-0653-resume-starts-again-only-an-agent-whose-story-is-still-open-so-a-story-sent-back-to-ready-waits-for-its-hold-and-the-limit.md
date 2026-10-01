---
id: T-0653
type: task
nature: remediation
title: resume starts again only an agent whose story is still open, so a story sent back to ready waits for its hold and the limit
status: done
parent: S-0182
owner: arobson
created: 2026-10-01T09:16:24Z
updated: 2026-10-01T09:19:38Z
transitions:
  - to: ready
    at: 2026-10-01T09:16:34Z
    by: agent-S-0182
  - to: in-progress
    at: 2026-10-01T09:16:54Z
    by: agent-S-0182
  - to: done
    at: 2026-10-01T09:19:38Z
    by: agent-S-0182
stream: S-0182
tags: []
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 164
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 9881
      cache_read: 2836171
      cache_write: 34968
      cost: 1.0318
---
# T-0653 resume starts again only an agent whose story is still open, so a story sent back to ready waits for its hold and the limit

## Work

`resume()` starts again any agent that ended asking once its thread is answered, for a story in any state but done or cancelled. A story sent back to ready (ADR-0055) or to backlog meanwhile gets its agent past its hold and the in-progress limit. Resume only a story in progress or in review: it is open, so it holds others and no hold applies to it. Leave a story in ready to the launcher's own path, which judges its hold and the limit.

## Done when

A test answers the thread of an asked run whose story is back in ready and held, and the launcher starts nothing for it; the same with a story in progress still resumes; `go test ./internal/serve` passes.

## Notes
