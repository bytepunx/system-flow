---
id: TH-0455
title: S-0387 is ready to accept, blocked only by TH-0447 awaiting your answer on auto-upgrade
anchor:
  path: wip/kanban/stories/S-0387-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md
  item: S-0387
status: resolved
participants: [orchestrator]
created: 2026-10-15T02:21:49Z
updated: 2026-10-15T02:28:17Z
---

# TH-0455 S-0387 is ready to accept, blocked only by TH-0447 awaiting your answer on auto-upgrade

On wip/kanban/stories/S-0387-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md.

## Entries

### 2026-10-15T02:21:49Z orchestrator
**Recommendation:** confirm my recommendation on TH-0447 (full plan, `auto-upgrade` off by default and shell-only), or answer it otherwise. Once TH-0447 is resolved, I accept S-0387 at 626b4f55.

## Where it stands

- `flai verify` passed every step at the branch head, 626b4f55.
- The verifier matched both criteria:
  - 1: the minute check in `flai/internal/serve/upgrade.go`, the shell-only host action (`hostapi/writes.go`), and the version-aware tier refusals (`verify/manifest.go`, `mcpserver/test.go` and `verify.go`), with tests, ADR-0158, and the docs.
  - 2: I-0141 is closed.
- **The two points I asked for are in ADR-0158.**
  - Strategic and story agents are detached and are adopted by the new serve.
  - An upgrade waits while a dashboard write runs.
  - A refused release is asked again at most every 15 minutes. After an hour, flai opens one thread per release.

## The one blocker

- `flai accept --dry-run`: "thread TH-0447 on S-0387 is answered, not resolved".
- TH-0447 holds a decision not yet recorded by you: a host action that upgrades your host's flai on its own.

## Also for you

ADR-0158 is `proposed`. Accepting it is yours, as with the other proposed ADRs.

### 2026-10-15T02:28:17Z orchestrator
Resolved: TH-0447 is settled: the operator approved at 02:28Z. S-0387 is no longer blocked.
