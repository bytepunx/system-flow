---
id: TH-0802
title: Merge queue batch 20261102-035614 fails with S-0518 at markdown
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [flai serve, alex]
created: 2026-11-02T03:56:41Z
updated: 2026-11-02T04:16:18Z
---

# TH-0802 Merge queue batch 20261102-035614 fails with S-0518 at markdown

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-02T03:56:41Z flai serve
Merge queue batch `20261102-035614` fails with S-0518 at `markdown`. The queue verified S-0518 alone on main at `77285a4a1ffa`, at `aec32de7338a`, the copy of its verified head `444b35bfb995` rebased in the batch, and `markdown` failed:

```text
markdown: lint-md: markdownlint-cli2 0.20.0
    npm warn Unknown env config "store-dir". This will stop working in the next major version of npm. See `npm help npmrc` for supported config options.
    markdownlint-cli2 v0.20.0 (markdownlint v0.40.0)
    Finding: design/adrs/0235-the-guard-hook-refuses-a-call-when-flai-guard-cannot-start-exit-126-or-127-and.md design/adrs/README.md design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md design/issues/summary.md design/system/agent-adapters.md docs/operators/settings.md docs/users/flai-reference.md docs/users/flai.md scripts/README.md template/CHANGELOG.md
    Linting: 11 file(s)
    Summary: 1 error(s)
    template/CHANGELOG.md:26 error MD024/no-duplicate-heading Multiple headings with the same content [Context: "Unreleased"]
    exit status 1
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261102-035614`.

flai moved S-0518 back to in-progress and took it out of the merge queue. Fix it on `story/S-0518`, run flai verify S-0518, move it to review, and enqueue it again with flai queue add S-0518.

### 2026-11-02T04:16:18Z alex
Resolved.
