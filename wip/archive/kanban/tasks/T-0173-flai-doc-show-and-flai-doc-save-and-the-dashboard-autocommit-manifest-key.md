---
id: T-0173
type: task
nature: feature
title: flai doc show and flai doc save, and the dashboard.autocommit manifest key
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:29Z
updated: 2026-09-19T02:17:58Z
transitions:
  - to: ready
    at: 2026-09-19T02:13:24Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:13:24Z
    by: alex
  - to: done
    at: 2026-09-19T02:17:58Z
    by: alex
stream: S-0040
tags: []
touches: [flai/cmd, flai/internal/manifest]
---

# T-0173 flai doc show and flai doc save, and the dashboard.autocommit manifest key

## Work
Add `flai doc show <path>` (JSON: path, content, hash, mode `full`, `body`, or `none`, and the reason when not `full`) and `flai doc save <path> --hash <sha256> [--message] [--trailer] [--no-commit]` reading the new content on stdin. Save: resolve the path inside design, docs, or wip as the MCP `doc_get` guard does (share it); refuse mode `none`; in mode `body` refuse when the front matter differs from the file's; compare the hash and on a mismatch fail with a conflict carrying the current content, its hash, and a unified diff between the current file and the submitted content (`git diff --no-index`, omitted without git); write atomically, run the check, and when it has an error finding on that path restore the old content and fail with the findings; for design and docs files whose front matter has `updated`, set it to today unless the designer changed it; then commit only that path unless `--no-commit` or the manifest's `dashboard.autocommit` is false: `docs:` for design and docs, `chore:` for wip, the item ID in brackets when the file is a work item or narrative. Add `Autocommit *bool` to the manifest's dashboard section, default true. Exit codes or error prefixes let the dashboard tell conflict and refusal from other failures. Tests with real git: save and commit with author and trailer, `--no-commit` and `autocommit: false`, body mode refusing a front matter change, mode none, conflict with diff, refusal with restore, the `updated` bump.

## Done when
- The new tests pass and cover every case listed
- `go test -race ./...` passes in `flai/` and golangci-lint is clean

## Notes
