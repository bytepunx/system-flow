---
id: T-1535
type: task
nature: feature
title: flai design init writes the token file and the guidance from starters, and flai design show and flai design css print them
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:41Z
updated: 2026-10-15T00:54:45Z
transitions:
  - to: ready
    at: 2026-10-15T00:30:27Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:30:28Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:54:45Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/cmd/design.go, flai/cmd/design_test.go, flai/cmd/root.go, flai/internal/designsystem/starter.md, flai/internal/designsystem/starter.tokens.json, flai/internal/designsystem/starter.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-1533, T-1590]
usage:
  source: log
  seconds: 1457
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 122
      output: 63393
      cache_read: 5954625
      cache_write: 287875
      cost: 4.0403
---
# T-1535 flai design init writes the token file and the guidance from starters, and flai design show and flai design css print them

## Work

Add the command group `flai design` in `flai/cmd/design.go`, registered in `root.go`. `init` writes the token file at the path the manifest's `design.tokens` names, from `flai/internal/designsystem/starter.tokens.json`, and `design/system/design-system.md` from `starter.md`, both embedded. It refuses when either exists, and writes neither. The token starter has every required role in light and dark, the dark values in the layout T-1592 settled, the `sans` and `mono` stacks, the type scale `xs` to `4xl`, the spacing step, and the radii and shadows `sm`, `md`, and `lg`. The guidance starter has front matter with `topics` and the guidance's headings. `show` prints a summary with the contrast, and with `--json` what `designsystem.Show` returns. `css` prints the custom properties, and with `--tailwind` the `@theme` block. Tests run each command against a temporary project. It waits for T-1533, whose `Show` and renderings it prints, and T-1590, whose key names the path.

## Done when

- `go test ./cmd/ -run Design` passes for init (new, token file existing, guidance existing), show (text and JSON, with and without a token file), and css (both forms).
- The starters pass `flai check --strict`, and every role pair of the starter passes its contrast threshold in both themes.

## Notes

- Layer 4, beside T-1536. Rewritten by planner-S-0381 from the front matter plan, as TH-0453 approved.
