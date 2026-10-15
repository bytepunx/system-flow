---
title: Design system
updated: 2026-10-15
audience: agent
order: 150
status: active
topics: [all]
roles: [story, verify]
---

# Design system

How a story's agent builds UI from the project's design system and proposes a change to it.

## Rules

- These rules apply when your pack holds the design system: its guidance, `design/system/design-system.md`, and its token file, at the path the manifest's `design.tokens` names. `flai prime` gives both, and this convention, only to the stories the guidance's `topics` reach, and to none in a project without a design system ([ADR-0152](../adrs/0152-agents-build-ui-from-a-tailwind-v4-stylesheet-flai-generates-from-the-tokens.md)).
- When the pack briefs them instead of holding them whole, read the guidance with `doc_get` and the tokens with `flai design show` before you build.
- Build UI from the tokens, through the utilities of the stylesheet `flai design css --write` generates at the path the manifest's `design.stylesheet` names.
- Import that stylesheet whole from the UI's own stylesheet.
- Never edit the generated stylesheet by hand. Change the token file and run `flai design css --write` again; until you do, `flai check` warns `design.stale`.
- Add no colour, font, type size, spacing, radius, or shadow outside the tokens: no arbitrary value such as `bg-[#3b82f6]`, and no hex.
- When the UI needs a value the tokens lack, change the design system on the story's branch, never on main:
  - add or change the token in the token file
  - widen the story's touches to the token file and the stylesheet
  - run `flai design css --write` and commit the stylesheet with the token file
  - say why in the guidance
- The operator previews the branch's design system on the Design page beside main's before accepting the story ([ADR-0151](../adrs/0151-the-operator-and-agents-change-a-design-system-by-edits-on-the-design-page.md)).
- Ask on a thread before a change that alters what a role means, such as what `primary` or `danger` is for: every screen that uses the role changes with it.
- Record the reason for a design decision in the guidance, under the heading it concerns, not only in the narrative or a commit.
- As a verifier, check the diff for a colour, font, type size, spacing, radius, or shadow outside the tokens, and for a hand edit of the generated stylesheet, and name each with its path and line.

## When in doubt

- If a value is not in the tokens, it is a change to the design system, not a detail of the UI: make it on the branch, or ask.
- If the guidance does not say which token fits, ask on a thread rather than pick one silently.

<!-- system-flow:end-of-baseline -->

## Project additions
