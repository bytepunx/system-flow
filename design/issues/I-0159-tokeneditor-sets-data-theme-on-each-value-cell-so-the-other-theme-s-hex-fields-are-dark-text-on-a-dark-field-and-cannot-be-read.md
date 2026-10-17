---
id: I-0159
title: TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-17T05:32:26Z
last_reported: 2026-10-17T05:32:26Z
updated: 2026-10-17T05:32:26Z
---

# I-0159 TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Description
TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Instances

### 2026-10-17T05:32:26Z
Story: S-0444.
S-0444's pass over flaiover's components rendered design/TokenEditor.svelte in the harness from its example "Colours, two edited". On a light page every dark cell's hex input shows dark text on a dark field, such as #0f172a for color.background, and on a dark page the light cells' fields are dark too. The component sets data-theme on each cell so that its swatch takes that theme's colours, and the input inherits the theme's field colours with the page's text colour. The Design page shows the same component, so the operator meets it there.

## Remediation
