---
id: T-1593
type: task
nature: feature
title: A design token colour converts from oklch to sRGB with clipping reported, and the colour roles' pairs get their contrast against ADR-0148's thresholds
status: backlog
parent: S-0381
owner: alex
created: 2026-10-14T23:53:18Z
updated: 2026-10-14T23:53:18Z
transitions: []
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/color.go, flai/internal/designsystem/color_test.go]
after: [T-1532]
---
# T-1593 A design token colour converts from oklch to sRGB with clipping reported, and the colour roles' pairs get their contrast against ADR-0148's thresholds

## Work

In `flai/internal/designsystem/color.go`, convert a colour of T-1532's types to sRGB: `srgb` from its components or `hex`, `oklch` by the fixed arithmetic through OKLab and linear sRGB, with no new dependency, clipping a colour outside sRGB to the gamut and saying it was clipped. Compute WCAG 2 relative luminance and the contrast ratio. `Contrast` takes one theme's colour tokens and answers each role pair ADR-0148 § Decision 4 names, with its ratio, its threshold, and whether it passes: `text`, `text-muted`, and each state on `background`, `surface`, or its `-soft` tint at 4.5:1; `on-primary` and `on-secondary` on their fill at 4.5:1; `border` and `focus` on `background` at 3:1. Test the conversion against known values, including a colour outside sRGB, and check the ratios against flaiover's own contrast test (`flaiover/src/lib/theme.test.ts`) so that flai and the dashboard agree. Contrast is reported, never a problem. It waits for T-1532, whose colour type it converts.

## Done when

- `go test ./internal/designsystem/` passes with conversion, clipping, ratio, and pair tests.
- The ratios match flaiover's contrast test on the same colours.

## Notes

- Layer 2, beside T-1592. Written by planner-S-0381.
