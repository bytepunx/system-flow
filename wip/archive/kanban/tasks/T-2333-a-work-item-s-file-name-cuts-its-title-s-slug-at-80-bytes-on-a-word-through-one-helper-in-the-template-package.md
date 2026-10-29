---
id: T-2333
type: task
nature: improvement
title: A work item's file name cuts its title's slug at 80 bytes on a word, through one helper in the template package
status: done
parent: S-0542
owner: alex
created: 2026-10-29T15:23:38Z
updated: 2026-10-29T15:28:38Z
transitions:
  - to: ready
    at: 2026-10-29T15:26:01Z
    by: agent-S-0542
  - to: in-progress
    at: 2026-10-29T15:26:02Z
    by: agent-S-0542
  - to: done
    at: 2026-10-29T15:28:38Z
    by: agent-S-0542
stream: S-0542
tags: [flai]
touches: [flai/internal/template/render.go, flai/internal/template/template_test.go, flai/internal/workitem/repo.go, flai/internal/workitem/create_test.go]
usage:
  source: log
  seconds: 156
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 5495
      cache_read: 818657
      cache_write: 45022
      cost: 0.5678
---
# T-2333 A work item's file name cuts its title's slug at 80 bytes on a word, through one helper in the template package

## Work

Add a helper beside `Slug` in `flai/internal/template/render.go` that cuts a slug at a given number of bytes on a word, as `flai/internal/adr/adr.go` and `flai/internal/digest/write.go` cut theirs at 80 today: cut at the limit, back to the last dash when that dash is past half the limit, and trim the trailing dash. Test it in `flai/internal/template/template_test.go`: a slug under the limit is unchanged, a long one is cut on a word, and one with no dash in range is cut at the limit.

Make `workitem.FileName` in `flai/internal/workitem/repo.go` cut its slug at 80 bytes with the helper. Every work item name comes from it: `flai story new`, `task new`, `epic new`, a retitle in `itemedit`, and an experiment's expected results document.

Reproduce I-0231 in `flai/internal/workitem/create_test.go`: create a story with S-0536's title, whose name today is 230 bytes, and assert that its file name is at most 92 bytes and that the name with Claude Code's `.tmp.<pid>.<hex>` suffix, about 30 bytes, fits in 255.

Waits for nothing: the other tasks build on this helper.

## Done when

- The helper and its test are in `flai/internal/template`.
- A work item created with a title whose full slug is over 80 bytes gets a file name of its ID, a dash, at most 80 bytes of slug, and `.md`, and the test that shows it passes.
- `flai test` on the four files passes.

## Notes
