---
id: TH-0410
title: "S-0366's head failed the integration tier in my verify run just before you accepted it: check main before the next release"
anchor:
  path: wip/archive/kanban/stories/S-0366-flai-check-finds-threads-archived-outside-the-story-at-close-out.md
  item: S-0366
status: resolved
participants: [orchestrator]
created: 2026-10-12T15:49:55Z
updated: 2026-10-12T15:57:13Z
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

### 2026-10-12T15:57:12Z orchestrator
Settled without needing your run. S-0364's last verify, at 3a66dca4 at 15:50, passed the integration tier on a base that holds S-0232, S-0365, and S-0366. So the failure at S-0366's head was not in main's code, most likely a flaky test. S-0313, in progress, makes the next such failure name its test.

I published S-0232, S-0364, and S-0366: flai 1.41.0 (`flai/v1.41.0`), flaiover 0.42.4, and template 1.0.87.

**What you need to do:** once `release-flai.yml` finishes for `flai/v1.41.0`, the first signed release, run `flai host upgrade`. Until then, the installed flai 1.40.1 rejects main's `{importers}` vet tier, so MCP `task_done`, `test`, and `verify` fail for the stories in progress: S-0234, S-0313, S-0337, and S-0349.

### 2026-10-12T15:57:13Z orchestrator
Resolved: S-0364's verify at 3a66dca4 passed integration on a base with S-0232, S-0365, S-0366; released as flai 1.41.0.
