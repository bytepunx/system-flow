---
id: TH-0410
title: "S-0366's head failed the integration tier in my verify run just before you accepted it: check main before the next release"
anchor:
  path: wip/archive/kanban/stories/S-0366-flai-check-finds-threads-archived-outside-the-story-at-close-out.md
  item: S-0366
status: open
participants: [orchestrator]
created: 2026-10-12T15:49:55Z
updated: 2026-10-12T15:49:55Z
---

# TH-0410 S-0366's head failed the integration tier in my verify run just before you accepted it: check main before the next release

On wip/archive/kanban/stories/S-0366-flai-check-finds-threads-archived-outside-the-story-at-close-out.md.

## Entries

### 2026-10-12T15:49:55Z orchestrator
**Recommendation:** run `make integration` on main once. If it passes, reply here and I publish S-0232 and S-0366. If it fails, its failing test is the thing to fix before a release.

## What I saw

At 15:43 I ran `flai verify S-0366` at its head, 8101e9f6. Every step passed until the integration tier, which failed:

```text
ok   github.com/bytepunx/system-flow/flai/tests/integration 3.030s
FAIL
exit status 1
```

The tail names no failing package, which is I-0113, S-0313's issue. The close-out's own run passed at ed64846c. The head after it added only the close-out's issue record, I-0139, and `design/issues/summary.md`. So this may be a flaky test, or one that broke from S-0365 or S-0232 landing on main meanwhile.

You accepted S-0366 at 15:48, while my run was going, so I could not run it again.

## What I am holding

- Pending release: S-0232 (2.63 USD a week) and S-0366 (12.5 USD a week). I publish once you answer here.
- S-0364 waits on TH-0409.
