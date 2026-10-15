---
id: S-0383
type: story
nature: feature
title: The operator changes the design system's tokens on the Design page, sees the preview follow before saving, and saves through flai
status: in-progress
parent: E-0020
owner: alex
created: 2026-10-13T17:44:47Z
updated: 2026-10-15T01:59:06Z
transitions:
  - to: ready
    at: 2026-10-15T01:10:39Z
    by: alex
  - to: in-progress
    at: 2026-10-15T01:56:59Z
    by: agent-S-0383
tags: [dashboard]
topics: [design]
touches: [flaiover/src/lib/design/tokens.ts, flaiover/src/lib/design/tokens.test.ts, flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0382]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1269
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-15
      work: 22
  models:
    - model: claude-opus-5-5
      input: 174
      output: 777
      cache_read: 9198161
      cache_write: 684063
      cost: 4.5714
  strategic:
    - kind: orchestrator
      seconds: 479
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 43
          output: 693
          cache_read: 9262776
          cache_write: 15505
          cost: 2.3302
cost_of_delay:
  value: 115.38
  by: planner-E-0020
  at: 2026-10-13T17:51:49Z
forecast:
  duration: 30m
  delivery: 2026-10-15T02:31:00Z
  basis: "Its own forecast of 30m; 2nd in the pull order with an in-progress limit of 5, behind S-0387 and S-0385."
  by: flai
  at: 2026-10-15T01:56:53Z
finalized:
  by: alex
  at: 2026-10-15T00:09:29Z
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
- T-1543 An edit model holds changed tokens apart from the saved ones, renders them as flai does for the preview, and writes them back into the DTCG token file, leaving the rest as it was
- T-1544 The preview flags each colour pair below 4.5:1 for text or 3:1 for controls, in the theme it fails in
- T-1545 A token editor changes each kind of token in place and marks a bad value and a failing contrast beside it
- T-1546 The Design page edits tokens with the preview following before saving, saves through doc.save, shows a refusal without losing the edits, and asks before leaving
- T-1547 The dashboard's design and user docs say how tokens are changed on the Design page and saved

## Notes

- Drafted by the planner for E-0020. It assumes the save goes through the existing `doc.show` and `doc.save` host methods, so that flai's strict check, with S-0381's `design.tokens`, refuses invalid tokens and no new write is needed.

### Planning

- **Touches.** None declared. All are the planner's, every one a file:
  - `flaiover/src/lib/design/tokens.ts` with its test: new, the edit model and the front matter round trip. Criteria 1 and 3.
  - `flaiover/src/lib/design/preview.ts` with its test: S-0382's, extended with the contrast check, in the way `theme.test.ts` holds the dashboard's palette. Criterion 2.
  - `flaiover/src/lib/components/design/TokenEditor.svelte` with its test: new. Criteria 1 and 2.
  - `flaiover/src/routes/design/+page.svelte` with its page test: S-0382's page, which gains the Edit mode, the save, and the leave guard of `routes/edit/[...path]/+page.svelte`. Criteria 3 and 4.
  - `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`: design. Criterion 5.
- **Not touched.** `flaiover/src/routes/api/docs/edit/+server.ts` and `api/docs/file/+server.ts` are called as they are, so no endpoint and no host write is added. S-0378, in progress, changes them; T-1546 says to read them as they are when the story starts.
- **Co-change.** `flai touches suggest` lists hub docs and `flaiover/src/lib/server/repo.ts`. None is added: the save uses the document endpoints, not a new read.
- **No folder touch.**
- **Layers.** T-1543 and T-1544 together; then T-1545; then T-1546; then T-1547.
- **Forecast: 30m, delivery 2026-10-13T23:11:00Z.** `flai forecast` gave 24m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 15. Raised a quarter to 30m for a YAML round trip that must keep the rest of the front matter and new Svelte components with component tests. The delivery adds the 6m to flai's 23:05Z.
- **Cost of delay: 115.38 USD a week.** `flai cod`: 30m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.
