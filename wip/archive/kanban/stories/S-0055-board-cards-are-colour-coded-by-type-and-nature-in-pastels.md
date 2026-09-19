---
id: S-0055
type: story
nature: feature
title: Board cards are colour coded by type and nature in pastels
status: done
parent: E-0006
owner: alex
created: 2026-09-19T02:05:29Z
updated: 2026-09-19T06:34:49Z
transitions:
  - to: ready
    at: 2026-09-19T05:35:11Z
    by: alex
  - to: in-progress
    at: 2026-09-19T06:05:07Z
    by: system-flow
  - to: review
    at: 2026-09-19T06:22:56Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:25:40Z
    by: system-flow
  - to: review
    at: 2026-09-19T06:32:59Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:34:49Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/lib/components, flaiover/src/routes]
---

# S-0055 Board cards are colour coded by type and nature in pastels

## Goal
A card's kind is recognisable at a glance from its colour: whether it is an epic, a story, or a task, and its nature (feature, improvement, remediation, research, experiment), in pastel colours, with the text that says the same thing kept.

## Acceptance criteria
- [x] Every card on the board is colour coded by nature, with a distinct pastel for each of feature, improvement, remediation, research, and experiment, in the light and the dark theme
- [x] Type is also visible from colour or a coloured element, so an epic, a story, and a task of the same nature can be told apart when "epics and tasks too" is on; how the two factors share the card (for example background tint for nature and an edge stripe for type) is chosen with the operator from rendered options
- [x] The text indications stay: nature, type, and the BLOCKED flag are still written on the card. Colour blindness checks on the palette are waived by the operator for that reason, and the story's notes and the design document say so; the text on every tint still meets the contrast the theme test requires
- [x] Blocked cards and the dragging state remain obvious on every tint
- [x] A legend on the board names the colours, and the colours are theme tokens in `layout.css`, not literals in the component, so the item page and charts can reuse them later
- [x] Component tests for each nature and type; the theme test extended to the new tokens; the board looked at in a browser in both themes; `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the coding
- [x] Added by the operator at review on 2026-09-19: the same coding appears on the item page (nature and type in the header), on item hits in search, and in the charts, where each nature's series colour stays a saturated colour of the same hue as its tint and the item table shows the nature on its tint

## Tasks
- T-0206 Nature and type colour tokens in both themes with contrast tests, and rendered options for the operator to choose from
- T-0207 Apply the chosen coding to the card, add the board legend, and cover each nature and type in component tests
- T-0208 Verify blocked and dragging on every tint in a browser in both themes, and describe the coding in the design and the user docs
- T-0209 The same coding on the item page, in search results, and in the charts

## Notes
Raised by the operator on 2026-09-18: "introduce a type of color-coding to the cards based on various factors: type (epic or story) QoS (is it a feature, improvement, remediation, research) so they are readily recognized. I would prefer us to use pastel color coding (it's ok if they don't pass color blind tests since there is also text indication)."

The operator's "QoS" is the item's `nature` in the schema (`design/system/work-hierarchy.md`), which has five values; `experiment` was not in the operator's list and gets a colour too unless they say otherwise.

The dashboard's colours are tokens in `flaiover/src/routes/layout.css` with a light and a selected dark theme (S-0044), and `src/lib/theme.test.ts` checks contrast. The waiver covers telling the pastels apart, not reading the text on them.

Whether the same coding should appear on the item page, in search results, and in the charts' nature series is worth asking when refining; this story is the board.

Shares `BoardCard.svelte` with S-0054; see the note there.
- 2026-09-19T06:25:40Z: moved to in-progress: Operator extended the scope at review: the same colour coding on the item page, in search results, and in the charts.
