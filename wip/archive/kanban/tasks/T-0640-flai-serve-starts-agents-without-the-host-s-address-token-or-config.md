---
id: T-0640
type: task
nature: remediation
title: flai serve starts agents without the host's address, token, or config
status: done
parent: S-0183
owner: arobson
created: 2026-10-01T08:34:17Z
updated: 2026-10-01T08:35:31Z
transitions:
  - to: ready
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
  - to: in-progress
    at: 2026-10-01T08:34:25Z
    by: agent-S-0183
  - to: done
    at: 2026-10-01T08:35:31Z
    by: agent-S-0183
stream: S-0183
tags: []
touches: [flai/internal/serve/agents.go]
usage:
  source: log
  seconds: 66
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5349
      cache_read: 1116082
      cache_write: 19871
      cost: 0.4729
---
# T-0640 flai serve starts agents without the host's address, token, or config

## Work
Build each agent's environment from flai serve's with FLAI_HOST_URL, FLAI_HOST_TOKEN, and FLAI_CONFIG removed, then the story's variables and the harness's. Test the environment a started agent sees.

## Done when
A test starts an agent from a serve whose environment holds the three and finds none of them in the agent's environment; the tests pass.

## Notes
