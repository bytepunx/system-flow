---
id: TH-0411
title: "Run flai host upgrade to 1.41.0: the installed 1.40.1 rejects main's manifest since S-0364"
anchor:
  path: wip/archive/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md
  item: S-0364
status: open
participants: [orchestrator]
created: 2026-10-12T15:57:27Z
updated: 2026-10-12T15:57:27Z
---

# TH-0411 Run flai host upgrade to 1.41.0: the installed 1.40.1 rejects main's manifest since S-0364

On wip/archive/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md.

## Entries

### 2026-10-12T15:57:27Z orchestrator
**Recommendation:** once the `release-flai.yml` run for `flai/v1.41.0` has finished, run `flai host upgrade`. Then reply here, and I resolve this thread.

## Why

- I published S-0232, S-0364, and S-0366 as flai 1.41.0, flaiover 0.42.4, and template 1.0.87.
- Main's `vet` tier now uses `{importers}`, from S-0364.
- The installed flai, 1.40.1, rejects that placeholder. Until the host runs 1.41.0, MCP `task_done`, `test`, and `verify` fail for S-0234, S-0313, S-0337, and S-0349.
- 1.41.0 is the first release signed in CI (S-0232). If the workflow fails at "Release key present", a secret is missing from the `release` environment; see `docs/operators/runbooks/release-key.md`.
