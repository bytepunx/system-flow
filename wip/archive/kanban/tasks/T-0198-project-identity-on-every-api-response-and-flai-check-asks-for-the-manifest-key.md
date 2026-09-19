---
id: T-0198
type: task
nature: feature
title: Project identity on every API response, and flai check asks for the manifest key
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:07Z
updated: 2026-09-19T05:20:18Z
transitions:
  - to: ready
    at: 2026-09-19T05:18:00Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:18:00Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:20:18Z
    by: system-flow
stream: S-0043
tags: []
touches: [flaiover/src, flai/internal/check]
---

# T-0198 Project identity on every API response, and flai check asks for the manifest key

## Work
In `src/lib/server/respond.ts`, add `project: { name, key }` from the manifest to every JSON object body, success or error, leaving arrays as they are; in `hooks.server.ts`, set `X-Flai-Project-Key` and `X-Flai-Project-Name` (URI-encoded) on every `/api/*` and `/mcp` response, including streams and errors. In `flai/internal/check`, add the warning `manifest.key` when `system-flow.yaml` has no `key`, pointing at the file; confirm the template always renders one. Tests: object, array, error, and streamed responses; the check rule with and without a key.

## Done when
- The tests pass on both sides
- This repository and a rendered template pass `flai check --strict`

## Notes
