---
id: T-1945
type: task
nature: remediation
title: The design and the contributors guide say the LiteLLM check falls back through its list of free models
status: backlog
parent: S-0456
owner: alex
created: 2026-10-21T06:03:41Z
updated: 2026-10-21T06:03:41Z
transitions: []
stream: S-0456
tags: [docs, tests]
touches: [design/system/agent-adapters.md, docs/contributors/index.md]
after: [T-1944]
---
# T-1945 The design and the contributors guide say the LiteLLM check falls back through its list of free models

## Work

Describe the fallback T-1944 builds where S-0455 described the probe:

- In `design/system/agent-adapters.md` § Checked against LiteLLM, say that the default is an ordered list of free models, which slugs are on it, that the check runs on the first the proxy answers, that an operator's own slug is not replaced, and that the check fails, never skips, when the whole list is refused. Cite I-0170 and S-0456. Bump `updated`.
- In `docs/contributors/index.md` § Gateway checks, change the `FLAI_TEST_LITELLM_HAIKU` row's default to the list and say when it falls back. Bump `updated`.

Waits for T-1944, because the list's slugs and the rules come from it.

## Done when

- Both documents name the list's slugs as T-1944 left them, and say what happens when the whole list is refused.
- `flai test design/system/agent-adapters.md docs/contributors/index.md` passes the markdown lint.

## Notes
