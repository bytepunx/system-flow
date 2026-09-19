---
id: S-0055
type: story
nature: feature
title: Board cards are colour coded by type and nature in pastels
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T02:05:29Z
updated: 2026-09-19T02:05:29Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/lib/components, flaiover/src/routes]
---

# S-0055 Board cards are colour coded by type and nature in pastels

## Goal
A card's kind is recognisable at a glance from its colour: whether it is an epic, a story, or a task, and its nature (feature, improvement, remediation, research, experiment), in pastel colours, with the text that says the same thing kept.

## Acceptance criteria
- [ ] Every card on the board is colour coded by nature, with a distinct pastel for each of feature, improvement, remediation, research, and experiment, in the light and the dark theme
- [ ] Type is also visible from colour or a coloured element, so an epic, a story, and a task of the same nature can be told apart when "epics and tasks too" is on; how the two factors share the card (for example background tint for nature and an edge stripe for type) is chosen with the operator from rendered options
- [ ] The text indications stay: nature, type, and the BLOCKED flag are still written on the card. Colour blindness checks on the palette are waived by the operator for that reason, and the story's notes and the design document say so; the text on every tint still meets the contrast the theme test requires
- [ ] Blocked cards and the dragging state remain obvious on every tint
- [ ] A legend on the board names the colours, and the colours are theme tokens in `layout.css`, not literals in the component, so the item page and charts can reuse them later
- [ ] Component tests for each nature and type; the theme test extended to the new tokens; the board looked at in a browser in both themes; `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the coding

## Tasks

## Notes
Raised by the operator on 2026-09-18: "introduce a type of color-coding to the cards based on various factors: type (epic or story) QoS (is it a feature, improvement, remediation, research) so they are readily recognized. I would prefer us to use pastel color coding (it's ok if they don't pass color blind tests since there is also text indication)."

The operator's "QoS" is the item's `nature` in the schema (`design/system/work-hierarchy.md`), which has five values; `experiment` was not in the operator's list and gets a colour too unless they say otherwise.

The dashboard's colours are tokens in `flaiover/src/routes/layout.css` with a light and a selected dark theme (S-0044), and `src/lib/theme.test.ts` checks contrast. The waiver covers telling the pastels apart, not reading the text on them.

Whether the same coding should appear on the item page, in search results, and in the charts' nature series is worth asking when refining; this story is the board.

Shares `BoardCard.svelte` with S-0054; see the note there.
