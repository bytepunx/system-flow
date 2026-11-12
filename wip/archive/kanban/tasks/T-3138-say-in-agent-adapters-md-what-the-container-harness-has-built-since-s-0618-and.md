---
id: T-3138
type: task
nature: remediation
title: Say in agent-adapters.md what the container harness has built since S-0618 and S-0628
status: done
parent: S-0725
owner: alex
created: 2026-11-12T15:51:36Z
updated: 2026-11-12T15:56:15Z
transitions:
  - to: ready
    at: 2026-11-12T15:55:27Z
    by: agent-S-0725
  - to: in-progress
    at: 2026-11-12T15:55:27Z
    by: agent-S-0725
  - to: done
    at: 2026-11-12T15:56:15Z
    by: agent-S-0725
stream: S-0725
tags: [docs]
touches: [design/system/agent-adapters.md]
usage:
  source: log
  seconds: 48
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1133
      cache_read: 749548
      cache_write: 21852
      cost: 0.3474
---
# T-3138 Say in agent-adapters.md what the container harness has built since S-0618 and S-0628

## Work

It waits for no task: it changes `design/system/agent-adapters.md` alone, which no other task of this story touches.

In § A container harness of flai's own, the paragraph that begins "What is built.", the last sentence says "Still not done: the close-out in the image (S-0510, S-0520), and the `permission_prompt` check of the image's Claude Code version, which `flai serve` only warns of." Both are done:

- The close-out of a container story runs its verify in a container of its own from the story's image, since S-0628 ([ADR-0274](../adrs/0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). Read ADR-0274 and `design/system/flai-cli.md` for what that container mounts and is passed.
- ADR-0106's `permission_prompt` check runs inside the image beside the start, since S-0618 (`considerImage` and `checkImage` in `flai/internal/serve/claudecheck.go`; `design/system/flai-cli.md` § The check inside an agent image (S-0618)).

Replace the sentence with one that says each is built, naming its story and ADR or design section, and link rather than repeat what `flai-cli.md` says. Keep what is still not built, flai's own loop, in the sentence earlier in the paragraph. Look for any other "not done" or "to build" claim in that section and its tables about the close-out or the image check, and correct it the same way. Bump the front matter's `updated`.

## Done when

- § A container harness of flai's own no longer says the close-out in the image or the image's `permission_prompt` check is not done, and says each is built with its story and its ADR or design section.
- `flai test design/system/agent-adapters.md` passes, the markdown lint among its tiers.

## Notes
