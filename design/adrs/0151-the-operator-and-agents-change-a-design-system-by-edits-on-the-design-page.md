---
id: ADR-0151
title: "The operator and agents change a design system by edits on the Design page, proposals as stories previewed beside main's, a request form that makes a ready story, and threads on the guidance's headings, with no store of variants"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
topics: [design, dashboard, cli]
---

# ADR-0151 The operator and agents change a design system by edits on the Design page, proposals as stories previewed beside main's, a request form that makes a ready story, and threads on the guidance's headings, with no store of variants

## Context

The operator decides a project's design system; agents propose changes and apply them. The finding, [design-system-tool.md](../system/design-system-tool.md), weighed edits on the Design page, proposals as stories, threads on the guidance, and a store of draft variants for how each is reviewed, recorded, and undone, and recommended choice 9 below. On TH-0445, on 2026-10-14, the operator answered only choice 1, so it stands as recommended.

## Decision

The operator and agents change a design system by edits on the Design page, proposals as stories previewed beside main's, a request form that makes a ready story, and threads on the guidance's headings, and there is no store of variants.

1. **Edits on the page.** The operator tries values with the preview following each change, and saves on main ([ADR-0150](0150-a-token-save-from-the-design-page-lands-on-main-and-flai-writes-the-stylesheet.md)). Unsaved values are kept in the browser across a reload.
2. **Proposals.** What an agent proposes or applies is a story whose branch changes the token file, the guidance, or both. The Design page lists the open proposals and previews a proposal's tokens, read from its worktree with `file.get` and `story`, beside main's. The operator tweaks them on the branch with `file.save` and `story` ([ADR-0145](0145-a-dashboard-save-may-name-a-story-and-is-then-committed-on-that-story-s-branch.md), [ADR-0147](0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)) while the story is open. Review, history, and undo are the board's.
3. **The request form.** A form on the Design page makes a story with `item.new`, as the board's new-item form does, and moves it to ready. Its goal is the operator's words; its criteria are fixed: the tokens as asked, the stylesheet in step, the UI's tests passing. Its touches are the token file and the stylesheet, and its body may carry the values the operator tried. The operator wrote it, so it is not a draft.
4. **Threads.** The Design page lists the guidance document's threads by heading and opens one on a heading, with the existing `thread.new`. A request that needs an agent becomes a story, which the thread names.
5. **No store of variants.** Nothing outside git holds alternative token sets.

## Consequences

- A design story claims the token file, so a second ready story that claims it is held until the first is accepted ([ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)).
- The stylesheet lies inside the UI sub-project, so a UI story that claims the sub-project by name overlaps every design story. Narrow touches keep them apart.
- A thread on the guidance discusses and records; it starts no agent and changes no file.

## Alternatives considered

- **A store of draft variants under `.flai-cache` or `wip/`.** It needs a store, methods, a picker, and a way to promote a variant. A store under `.flai-cache` is lost to history, and agents see neither. Unsaved edits on the page give what it would.
