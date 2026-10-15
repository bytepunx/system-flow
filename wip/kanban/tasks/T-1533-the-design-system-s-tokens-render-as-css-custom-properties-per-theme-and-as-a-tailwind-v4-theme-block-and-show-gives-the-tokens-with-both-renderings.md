---
id: T-1533
type: task
nature: feature
title: The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:32Z
updated: 2026-10-15T00:29:25Z
transitions:
  - to: ready
    at: 2026-10-15T00:22:31Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:22:32Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:29:25Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go, flai/internal/designsystem/testdata/full.css.golden, flai/internal/designsystem/testdata/full.tailwind.css.golden, flai/internal/designsystem/testdata/show-missing.json.golden, flai/internal/designsystem/testdata/show-problems.json.golden, flai/internal/designsystem/testdata/show.json.golden, flai/internal/designsystem/testdata/themes.css.golden, flai/internal/designsystem/testdata/themes.tailwind.css.golden, flai/internal/designsystem/testdata/themes.tokens.json, flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go]
after: [T-1592, T-1593]
usage:
  source: log
  seconds: 412
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 44
      output: 199
      cache_read: 1852432
      cache_write: 122899
      cost: 0.9121
---
# T-1533 The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings

## Work

In `flai/internal/designsystem/css.go`, `CSS` renders the themes T-1592 resolves as custom properties: a `:root` block for light and a `[data-theme='dark']` block for dark, the selector S-0385 may change when it settles how the generated stylesheet selects dark values. `Tailwind` renders them as a Tailwind v4 `@theme` block, each group named after its namespace (`--color-*`, `--font-*`, `--text-*` with `--line-height`, `--spacing`, `--radius-*`, `--shadow-*`). An `oklch` colour is written as `oklch()`, an `srgb` one as its hex. Both are deterministic, ordered by group then name. `Show` returns the shape `flai design show --json` and `design.get` answer: the token file's path, the tokens of each theme, the guidance's headings, the problems, both renderings, and T-1593's contrast for each theme, or that the project has no token file. Golden tests in `css_test.go`. It waits for T-1592, whose themes it renders, and T-1593, whose contrast `Show` reports.

## Done when

- `go test ./internal/designsystem/` passes with golden tests of both renderings, and of `Show` with a valid file, with a file that has problems, and with none.

## Notes

- Layer 3, beside T-1534, with which it shares no file. Rewritten by planner-S-0381 from the front matter plan, as TH-0453 approved.
