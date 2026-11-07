---
id: T-2583
type: task
nature: remediation
title: The story page and the host agent notice render the refusal's settings pages as links
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:12:44Z
updated: 2026-11-07T10:52:10Z
transitions:
  - to: ready
    at: 2026-11-07T10:32:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:32:07Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:52:10Z
    by: agent-S-0587
stream: S-0587
tags: [dashboard]
touches: [flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/HostAgentNotice.svelte, flaiover/src/lib/components/HostAgentNotice.svelte.test.ts]
after: [T-2577]
usage:
  source: log
  seconds: 1162
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 52
      output: 357
      cache_read: 8303203
      cache_write: 36086
      cost: 2.9912
---
# T-2583 The story page and the host agent notice render the refusal's settings pages as links

## Work

It waits for T-2577: ADR-0261 § 4 fixes the form, a `/settings/<page>` path written out in the refusal's or the waiting reason's words, such as `/settings/harnesses` and `/settings/agents`. It needs no code of T-2582's: the tests use messages of that form.

- Extend `reasonParts` in `flaiover/src/lib/activity.ts`, which links only story IDs today. It should also link `/settings/<page>` paths that name a settings page in `settingspages.ts`, resolved with `resolve('/settings/[page]', …)`. A path naming no page stays text; a trailing comma or full stop is not part of the path.
- In `flaiover/src/lib/components/StoryAgent.svelte`, render the Start agent refusal and the "flai serve has started no agent" waiting reason through it.
- In `flaiover/src/lib/components/HostAgentNotice.svelte`, render the waiting reasons through it too.

## Done when

- The tests show that a refusal naming `/settings/agents` or `/settings/harnesses` renders a link to that page: `activity.test.ts`, `StoryAgent.svelte.test.ts`, and `HostAgentNotice.svelte.test.ts`.
- `flai test` on the changed paths passes the flaiover tiers.

## Notes
