---
id: ADR-0171
title: "A sub-project whose svelte is not installed gets one components.examples note, as a host without node does, and a check scoped to a story leaves out each such note unless the story changes a file under the examples folder it names"
status: accepted
date: 2026-10-18
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0164]
---

# ADR-0171 A sub-project whose svelte is not installed gets one components.examples note, as a host without node does, and a check scoped to a story leaves out each such note unless the story changes a file under the examples folder it names

## Context

[ADR-0164](0164-a-component-s-examples-are-named-sets-of-props-in-one-json-file-per-component.md) §5 has `flai check` report `components.examples`, a warning, on each examples file whose component's props flai could not read, once, with the reader's error. It passes over a project with no supported sub-project, and a host without `node`, with one note: a warning `flai check --strict` passes over.

I-0161, ``flai check finds `components.examples` outside the story at close-out``, has two instances, from the close-outs of S-0392 and S-0393 on 2026-10-17. Both have the same cause.

1. **A story's worktree has no `flaiover/node_modules`.** The props reader runs on the sub-project's own `svelte/compiler` ([ADR-0162](0162-the-component-library-reaches-sub-projects-of-kind-sveltekit-or-svelte-finds.md)), and `components.ReadProps` answers `ErrSvelteMissing`. Nothing installs the sub-project before the check: `scripts/flaiover-install.sh --if-needed` runs only in the test tiers, after it.
2. **The check warns on every examples file.** About 48 files, each with the same message: "svelte is not installed in flaiover … run npm install".
3. **None of them is the story's.** Neither story changed an examples file, so every warning was outside the story, a note, and `--record-issues` recorded all of them in I-0161. The closing agent can fix nothing: the files are right, and only the worktree's install is missing.

## Decision

A sub-project whose svelte is not installed gets one `components.examples` note, as a host without `node` does, and a check scoped to a story leaves out each such note unless the story changes a file under the examples folder it names.

1. **One note per sub-project without svelte.** When the props reader answers that svelte is not installed in a sub-project, `flai check` reports one `components.examples` note on that sub-project's examples folder, `design/components/<sub-project>`, saying that its examples' props are not checked and to install the sub-project. It is a warning `flai check --strict` passes over, as the note for a host without `node` is. The format and whether each file's component exists are still checked. Any other error of the reader stays a warning on each file, once, as ADR-0164 says.
2. **A check scoped to a story leaves such a note out unless the story changes a file under its folder.** `flai check --story S-nnnn` drops from its result the note for a sub-project without svelte, and the note for a host without `node`, whose folder is the examples folder, when no path the story's branch changes, or that is uncommitted in its worktree, is under the note's folder, and takes back the counts it added. When one is, the note is inside the story and stays: the story's own examples went unchecked, and its agent should know.
3. **`--record-issues` records neither note outside the story.** What is left out of the result is not recorded.
4. **Without `--story` both notes are reported**, so the main checkout's `flai check` still says which examples it could not check.

This refines ADR-0164 §5, whose reader error on each file is now one note when the cause is svelte not installed, and ADR-0085, whose rule that every finding outside the story is recorded in an issue now leaves out these notes, as [ADR-0122](0122-a-close-out-records-no-item-archive-and-a-check-scoped-to-a-story-leaves-out.md), [ADR-0133](0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md), [ADR-0136](0136-a-close-out-records-no-issues-no-story-and-a-check-scoped-to-a-story-leaves-out.md), and [ADR-0137](0137-a-close-out-records-no-threads-archived-and-a-check-scoped-to-a-story-leaves.md) do for other rules.

## Consequences

- A close-out in a worktree without the sub-project installed no longer records an issue for examples it did not change, and I-0161 can be closed.
- A story that changes an examples file in such a worktree sees one note, inside the story, that its examples' props went unchecked; `--strict` passes over it, so the test tiers, which install the sub-project, and the main checkout's check stay where a broken example is caught.
- A real format error, or a file whose component does not exist, is still a warning on its file, inside or outside the story, as before.

## Alternatives considered

- **Leave out every `components.examples` finding outside the story.** It would hide a real format error or a missing component on main from every close-out.
- **Install the sub-project before the check.** It gives `flai check` a side effect, a package install with network access, on every close-out and every save.
- **One note, but still recorded outside the story.** The closing agent cannot act on it, so the issue would count only how many close-outs ran in a worktree without `node_modules`.
