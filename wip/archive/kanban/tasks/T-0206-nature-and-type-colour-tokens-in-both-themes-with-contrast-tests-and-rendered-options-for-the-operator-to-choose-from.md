---
id: T-0206
type: task
nature: feature
title: Nature and type colour tokens in both themes with contrast tests, and rendered options for the operator to choose from
status: done
parent: S-0055
owner: alex
created: 2026-09-19T06:05:57Z
updated: 2026-09-19T06:17:24Z
transitions:
  - to: ready
    at: 2026-09-19T06:06:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:06:22Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:17:24Z
    by: system-flow
stream: S-0055
tags: []
---

# T-0206 Nature and type colour tokens in both themes with contrast tests, and rendered options for the operator to choose from

## Work
Add theme tokens to `flaiover/src/routes/layout.css` for the five natures (a pastel tint each) and the three types, in the light and the dark block, exposed through `@theme inline` so Tailwind classes exist for them. Extend `src/lib/theme.test.ts`: every new token is present in both themes, and `ink`, `muted`, and `danger` read at 4.5:1 on every nature tint, because the card's title, details, and BLOCKED flag sit on it. Telling the pastels apart is waived by the operator; reading the text on them is not. Then render real cards from the board, with the S-0054 divider, in each way the two factors can share a card (tint for nature with an edge stripe for type; tint for nature with the ID as a type-coloured chip; tint for type with a stripe for nature), in both themes, send the images to the operator, and ask which they want. Do not pick for them. Ask in the same question whether `experiment`, which was not in their list, gets a colour.

## Done when
- The tokens are in `layout.css` for both themes and the theme test covers them and passes
- The options were sent to the operator and their choice is recorded in the narrative's Decisions

## Notes
