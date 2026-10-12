---
id: TH-0425
title: S-0367 is ready to accept but changes docs/users/flai-reference.md outside its touches
anchor:
  path: wip/kanban/stories/S-0367-perform-lint-fix-on-every-markdown-file-a-story-touches-as-part-of-validation.md
  item: S-0367
status: resolved
participants: [orchestrator, alex]
created: 2026-10-12T17:38:41Z
updated: 2026-10-12T17:44:30Z
---

# TH-0425 S-0367 is ready to accept but changes docs/users/flai-reference.md outside its touches

On wip/kanban/stories/S-0367-perform-lint-fix-on-every-markdown-file-a-story-touches-as-part-of-validation.md.

## Entries

### 2026-10-12T17:38:41Z orchestrator
**Recommendation:** add the generated reference to S-0367's touches. Then reply here, and I accept S-0367 at once:

```sh
flai touches S-0367 --add docs/users/flai-reference.md
```

## Where S-0367 stands

- `flai verify` passed every step at the branch head, 7bfd0d9a, run through the CLI on flai 1.43.0.
- The verifier matched all four criteria to the diff:
  1. The tier `fix_command`, with task done and sync fixing.
  2. Accept fixes the main checkout before it commits.
  3. Every rebase fixes.
  4. Accept refuses what the fix left.
- ADR-0139 and ADR-0140 are present and indexed.

## What stops it

The dry-run says:

```text
blocked: story/S-0367 changes docs/users/flai-reference.md, which is under none of the touches of S-0367
```

The file is the reference generated from the changed `accept.go` and `stream.go` help, added by the head commit 7bfd0d9a. It is in `claims.shared`, which keeps it from holding other stories, but acceptance still needs it in the story's touches (ADR-0093). Changing a story's touches is not among my permissions, so I leave it in review.

### 2026-10-12T17:44:30Z alex
Resolved: S-0367 was accepted
