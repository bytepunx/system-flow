---
id: S-0383
type: story
nature: feature
title: The operator changes the design system's tokens on the Design page, sees the preview follow before saving, and saves through flai
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:44:47Z
updated: 2026-10-13T17:44:47Z
transitions: []
tags: [dashboard]
topics: [design]
touches: [flaiover/src/lib/design/tokens.ts, flaiover/src/lib/design/tokens.test.ts, flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0382]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0383 The operator changes the design system's tokens on the Design page, sees the preview follow before saving, and saves through flai

## Goal

The operator tries a decision on the design system, a colour, a font, a step of the type scale, and sees the sample components respond at once on the Design page, before anything is saved. When it looks right, they save it, and flai writes it into the design system document, so the decision is recorded where agents read it.

## Acceptance criteria

- [ ] Each token on the Design page can be changed in place: a colour with a picker and a text field, a length or a scale step with a number and its unit, a font family with a text field.
- [ ] The preview follows each change at once, before saving, and a colour pair that falls below 4.5:1 for text or 3:1 for controls is flagged beside the token, in the theme it fails in.
- [ ] Saving writes the tokens into the document's front matter through flai (`doc.save`, ADR-0023), leaving its body and the rest of its front matter as they were, and shows flai's refusal, such as a `design.tokens` finding or a conflict with a change made meanwhile, without losing the edits.
- [ ] Leaving the page with unsaved changes asks first, as the document editor does, and the changes can be discarded back to the saved tokens.
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say how tokens are changed and saved.

## Tasks

## Notes

- Drafted by the planner for E-0020. It assumes the save goes through the existing `doc.show` and `doc.save` host methods, so that flai's strict check, with S-0381's `design.tokens`, refuses invalid tokens and no new write is needed.
