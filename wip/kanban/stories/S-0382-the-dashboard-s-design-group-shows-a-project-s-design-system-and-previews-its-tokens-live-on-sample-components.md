---
id: S-0382
type: story
nature: feature
title: The dashboard's Design group shows a project's design system and previews its tokens live on sample components
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:44:35Z
updated: 2026-10-13T17:44:35Z
transitions: []
tags: [dashboard]
topics: [design]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte.test.ts, flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/components/design/PreviewFrame.svelte, flaiover/src/lib/components/design/PreviewFrame.svelte.test.ts, flaiover/src/lib/components/design/Gallery.svelte, flaiover/src/lib/components/design/TokenSheet.svelte, flaiover/src/lib/components/design/TokenSheet.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0381]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0382 The dashboard's Design group shows a project's design system and previews its tokens live on sample components

## Goal

The operator opens a new top-level Design group in the dashboard's site menu and sees the project's design system as it stands: its tokens laid out, and a gallery of sample components styled by them in light and dark, so the effect of the design on the project's visual language is visible in the browser without building or deploying the project.

## Acceptance criteria

- [ ] The site menu has a fourth group, Design, beside Workflow, Status, and Host, whose first page, `/design`, keeps `?project=` as the other pages do.
- [ ] The page shows the project's design system: its colours as swatches with each one's contrast against the ground in both themes, its type scale as sample text, its spacing, radii, and shadows drawn to size, and its guidance rendered from the document.
- [ ] The page previews the tokens on a gallery of sample components, at least buttons, inputs, a form, a card, an alert, a table, navigation, and headings with body text, in light and dark, isolated so that the dashboard's own theme neither reaches the gallery nor changes with it.
- [ ] A project with no design system document is told so, with the command that writes one.
- [ ] The page follows a change to the document made elsewhere, as the other pages follow theirs.
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the Design group and its page.

## Tasks

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends an isolated preview (an iframe or a shadow root) whose CSS custom properties come from the tokens over components compiled with the dashboard, so that a changed token restyles the gallery with no build.
