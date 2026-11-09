---
id: T-2780
type: task
nature: improvement
title: Describe the reference tests' tier in devex.md and tooling.md and close I-0308
status: backlog
parent: S-0634
owner: alex
created: 2026-11-09T21:58:38Z
updated: 2026-11-09T21:58:38Z
transitions: []
stream: S-0634
tags: [docs, testing]
touches: [design/system/devex.md, design/conventions/tooling.md, design/issues/I-0308-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md, design/issues/summary.md]
after: [T-2776]
---
# T-2780 Describe the reference tests' tier in devex.md and tooling.md and close I-0308

## Work

Two documents describe the `settings-index` tier as running only the `TestSettingsIndex` tests, selected by `docs/operators/settings.md`. Update both to what T-2773 and T-2776 built. The first is the Test tiers row of `design/system/devex.md` § Used by this monorepo. The second is the `flai test` line of the project additions in `design/conventions/tooling.md`, below the marker. Say that `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` run under `-short`, so `go-test` runs them on any change to `flai/cmd` or a package it imports, such as `flai/internal/manifest/settings.go`. Say that `settings-index` runs them too, selected by `docs/users/flai-reference.md` as well (I-0308). Bump each document's `updated`.

Then close the issue: `flai issue close I-0308 --reason "<what fixed it>"`, naming both changes and S-0634. It writes the issue file and regenerates `design/issues/summary.md`.

This task waits for T-2776: it describes the tier as T-2776 leaves it.

## Done when

- `design/system/devex.md` and `design/conventions/tooling.md` describe the tier and the two tests as built, and `updated` is today's date in both.
- `flai issue list` no longer lists I-0308 as open, and its file gives the reason.
- `flai test design/system/devex.md design/conventions/tooling.md` passes, and `flai check --strict` reports nothing on S-0634.

## Notes
