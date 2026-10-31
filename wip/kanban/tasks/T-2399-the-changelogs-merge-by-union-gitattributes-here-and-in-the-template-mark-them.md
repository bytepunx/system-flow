---
id: T-2399
type: task
nature: improvement
title: "The changelogs merge by union: .gitattributes here and in the template mark them, flai release names one Unreleased section, and the git convention says why"
status: in-progress
parent: S-0550
owner: alex
created: 2026-10-31T07:43:42Z
updated: 2026-10-31T08:17:32Z
transitions:
  - to: ready
    at: 2026-10-31T08:17:30Z
    by: agent-S-0550
  - to: in-progress
    at: 2026-10-31T08:17:32Z
    by: agent-S-0550
stream: S-0550
tags: [template, cli]
touches: [".gitattributes", template/root/.gitattributes, template/.gitattributes, flai/internal/release/release.go, flai/internal/release/release_test.go, flai/cmd/stream_sync_test.go, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md]
---
# T-2399 The changelogs merge by union: .gitattributes here and in the template mark them, flai release names one Unreleased section, and the git convention says why

## Work

Neither `.gitattributes` (root) nor `template/root/.gitattributes`, which is rendered into every project, has a `merge=` line. `template/CHANGELOG.md` sits outside `template/root/`, so in this monorepo the root file governs it.

- Add `CHANGELOG.md merge=union` to `.gitattributes` and `template/root/.gitattributes`. That covers `template/CHANGELOG.md` and the `<component>/CHANGELOG.md` that `flai release` writes. Check whether `flai template push` publishes `template/` as the template repository's root. If it does, add `template/.gitattributes` with the same line, so that the template's own changelog merges by union there. If it does not, leave that file out and say so in `## Decisions`.
- Union keeps both sides' lines. Today a story adds its changelog lines under a version heading it writes itself, such as `## 1.1.13 - 2026-10-31`, and two stories that both add one leave two sibling headings, which MD024 (`siblings_only`) refuses (I-0107, MS-0032, MS-0201). Make `## Unreleased` the one place a story writes:
  - `template/CHANGELOG.md` gains `## Unreleased` at its top.
  - `flai release` (`flai/internal/release/release.go`, where it writes `<component>/CHANGELOG.md`) renames `## Unreleased` to `## <version> - <date>`, puts its bullets there, and leaves a new empty `## Unreleased` above it. A changelog without the section is released as today.
  - Test it in `flai/internal/release/release_test.go`.
- Show in `flai/cmd/stream_sync_test.go` that the trial merge at `flai stream sync` reports no conflict for two branches that each add a line under `## Unreleased` with the attribute set. Check that `git merge-tree --write-tree` honours the attribute from the merged trees. If it does not, record what it reads in `## Decisions` and make the test show what does merge.
- In `design/conventions/git.md` and `template/root/design/conventions/git.md`, beside the changelog rule ("Every release gets a changelog entry"), say which files carry `merge=union` and why, and that a story writes its lines under `## Unreleased`, never under a version heading. This is a baseline change, made in both copies in the same story.

First layer: it waits for nothing. The documentation task waits for it, since both write `template/CHANGELOG.md`.

## Done when

Both `.gitattributes` files carry the line, and `template/.gitattributes` does too or `## Decisions` says why not. The release test shows `## Unreleased` becoming the version's section, with a fresh empty `## Unreleased` above it. The sync test shows two branches' changelog lines merging with no conflict. Both copies of `git.md` name the files and the reason. Story criterion 2.

## Notes

Drafted by the planner of S-0550 on 2026-10-31. The criterion also names the digest index. The digest has none (`design/digest/README.md`), and the story's ADR records that, so no task gives a digest file the attribute.
