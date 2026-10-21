---
id: T-1946
type: task
nature: remediation
title: The gateway check passes at smoke on this host with the fallback, and I-0170 is closed
status: done
parent: S-0456
owner: alex
created: 2026-10-21T06:03:47Z
updated: 2026-10-21T07:50:14Z
transitions:
  - to: ready
    at: 2026-10-21T07:48:43Z
    by: agent-S-0456
  - to: in-progress
    at: 2026-10-21T07:48:44Z
    by: agent-S-0456
  - to: done
    at: 2026-10-21T07:50:14Z
    by: agent-S-0456
stream: S-0456
tags: [issues, tests]
touches: [design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, design/issues/summary.md]
after: [T-1944]
usage:
  source: log
  seconds: 90
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2917
      cache_read: 719006
      cache_write: 16966
      cost: 0.3199
---
# T-1946 The gateway check passes at smoke on this host with the fallback, and I-0170 is closed

## Work

Show the cause I-0170 describes no longer stops a close-out, then close it:

- Run the smoke tier's gateway check, `TestGatewayLiteLLM`, as the close-out runs it, from `.flai-cache/litellm/litellm.env` with no `FLAI_TEST_LITELLM_HAIKU` set by hand, and record the slug it ran on and the result in this task's notes.
- Close the issue with `flai issue close I-0170 --reason`, saying what fixed it: the probe from S-0455 and the fallback list from T-1944, and that the check fails, naming each slug, only when the whole list is refused. `flai issue close` keeps `design/issues/summary.md` current.

Waits for T-1944, whose fallback the run checks. It runs beside the documentation task, with no path in common.

## Done when

- `TestGatewayLiteLLM` passed through the proxy at smoke with the default list, recorded in this task's notes.
- I-0170 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean for the story.

## Notes

`TestGatewayLiteLLM` ran on 2026-10-21 through `scripts/gateway-smoke.sh`, as the smoke tier runs it, from `.flai-cache/litellm/litellm.env` with no `FLAI_TEST_LITELLM_HAIKU` set by hand: the key file names nemotron-3-ultra, on the list, so the check walked the list and passed on `openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free` in 146.3 s at $0 of key spend (T-1944's run). With the retired nemotron-3-super put first, it passed over it on its 404 and passed on nemotron-3-ultra in 147.9 s. The story's close-out runs the smoke tier again.
