---
id: I-0159
title: TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-17T05:32:26Z
last_reported: 2026-10-17T05:32:26Z
updated: 2026-10-21T12:13:31Z
---

# I-0159 TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Description
TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Instances

### 2026-10-17T05:32:26Z
Story: S-0444.
S-0444's pass over flaiover's components rendered design/TokenEditor.svelte in the harness from its example "Colours, two edited". On a light page every dark cell's hex input shows dark text on a dark field, such as #0f172a for color.background, and on a dark page the light cells' fields are dark too. The component sets data-theme on each cell so that its swatch takes that theme's colours, and the input inherits the theme's field colours with the page's text colour. The Design page shows the same component, so the operator meets it there.

## Remediation

Story S-0474 remediates this issue, created from it at 2026-10-21T11:59:01Z.
Closed 2026-10-21T12:13:31Z: S-0474: TokenEditor marks each value cell's theme with data-cell-theme, which no stylesheet reads, instead of data-theme, which layout.css reads as the dark palette and Tailwind's dark variant on any element. The page's theme now styles every cell, and a test in TokenEditor.svelte.test.ts fails while any element in the editor carries data-theme.
