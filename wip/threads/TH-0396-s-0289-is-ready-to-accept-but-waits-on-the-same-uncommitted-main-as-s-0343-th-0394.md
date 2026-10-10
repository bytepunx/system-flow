---
id: TH-0396
title: S-0289 is ready to accept, but waits on the same uncommitted main as S-0343 (TH-0394)
anchor:
  path: wip/kanban/stories/S-0289-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md
  item: S-0289
status: resolved
participants: [orchestrator, alex]
created: 2026-10-09T18:38:32Z
updated: 2026-10-10T18:58:58Z
---

# TH-0396 S-0289 is ready to accept, but waits on the same uncommitted main as S-0343 (TH-0394)

On wip/kanban/stories/S-0289-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md.

## Entries

### 2026-10-09T18:38:32Z orchestrator
**Recommendation:** clear main as TH-0394 says: commit or stash your work in progress, and run `flai accept S-0338 --by alex` (TH-0389). Then I accept S-0342, S-0343, and S-0289 at once. Or accept them yourself.

S-0289 passes everything I check:

- `flai verify` passed every step at the head, b044db7d, including integration and smoke.
- The verifier matched both criteria:
  1. `md056` in `flai/internal/mdlint/rules.go`, with `TestTableColumnCountOfI0077` and the `testdata/cases/tables.md` fixture reproducing I-0077's row.
  2. I-0077 is closed, its reason naming MD056.
- Every changed file is within its touches, apart from I-0134 and I-0135, which the close-out recorded.

The only blocker, from `flai accept S-0289 --by orchestrator --verified b044db7d --dry-run`:

```text
uncommitted outside wip: Makefile, docs/contributors/index.md, scripts/README.md, scripts/env.sh, scripts/litellm.sh; the real run refuses until they are committed or stashed, or --yes includes them
```

I will not use `--yes`. I leave S-0289 in review.

### 2026-10-10T18:58:58Z alex
Resolved: S-0289 was accepted
