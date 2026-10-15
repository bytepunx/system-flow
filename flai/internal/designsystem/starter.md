---
title: Design system
updated: 2026-10-14
status: active
topics: [all]
---

# Design system

The visual language of this project's UI: when to use which colour, size, space, shape, and component. The values themselves live in the token file, `{{tokens}}`, which `flai check` validates and `flai design css` renders; this document says how to use them and never repeats them.

## Principles

Say in a few lines what the UI should feel like and what it puts first, such as legibility, density, or calm. A reader settling a choice the sections below do not cover should be able to decide from these.

## Colour

Say what each colour role is for: `background` and `surface` for the ground, `text` and `text-muted` for words, `primary` and `secondary` for actions with `on-primary` and `on-secondary` on them, `border` and `focus` for edges and the focus ring, and each state with its `-soft` tint for alerts. Name any extra colour token here with what it is for.

## Typography

Say which font stack is for what, `sans` for the interface and `mono` for code and figures, and which step of the type scale each kind of text uses, from `xs` to `4xl`.

## Spacing and layout

Say how the spacing step is used: the multiples for gaps inside a component, between components, and between sections, and how pages are laid out at each width.

## Shape and elevation

Say which radius goes on which kind of element, `sm`, `md`, or `lg`, and which shadow lifts what, from `sm` for a raised control to `lg` for a dialog.

## Components

List the components the UI is built from, and for each when to use it, its variants, and its states. Build each from the token utilities, never from a value of its own.

## Accessibility

Every role pair meets its contrast threshold in light and dark, as `flai design show` reports: 4.5:1 for text, 3:1 for borders and the focus ring. Say here what else the UI guarantees, such as a visible focus, a target size, and motion that respects reduced-motion settings.

## Changing the design system

Change a value in the token file, never in the code, and say here why it changed. Ask before changing what a role means, since every screen that uses the role changes with it.
