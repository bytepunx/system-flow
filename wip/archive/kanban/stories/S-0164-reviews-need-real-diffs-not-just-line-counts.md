---
id: S-0164
type: story
nature: feature
title: Reviews need real diffs, not just line counts
status: done
parent: E-0014
owner: alex
created: 2026-09-29T21:10:42Z
updated: 2026-09-29T22:18:34Z
transitions:
  - to: ready
    at: 2026-09-29T21:10:44Z
    by: alex
  - to: in-progress
    at: 2026-09-29T21:11:17Z
    by: agent-S-0164
  - to: review
    at: 2026-09-29T21:22:54Z
    by: agent-S-0164
  - to: done
    at: 2026-09-29T22:18:34Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
usage:
  source: log
  seconds: 741
  models:
    - model: claude-fable-5-1
      input: 96
      output: 39175
      cache_read: 5516982
      cache_write: 171726
      cost: 6.7735
---
# S-0164 Reviews need real diffs, not just line counts

## Goal

Each file diff currently shows the path and the lines added/removed.

Add to each of these a toggle control to expand and collapse the detail line to reveal a control that shows the removed lines (- sign in the margin, dim red background behind each line with a brighter red outline around consecutive lines removed) and the added lines (+ sign in the margin, dim green background behind each line with a brighter green outline around consecutive lines added)

## Acceptance criteria
- [x] a toggle control on the diff summary line reveals or hides the diff panel
- [x] a `-` in the margin of every subtracted line
- [x] a `+` in the margin of every added line
- [x] a dim red background behind removed lines and a brighter red outline surrounding consecutive removed lines
- [x] a dim green background behind added lines and a brighter green outline surrounding consecutive added lines

## Tasks
- T-0585 A patch is read as runs of consecutive added and removed lines, each line's sign apart from its text
- T-0586 Each file's summary line has a toggle control, and its panel shows signs in a margin, dim backgrounds, and an outline round each run
- T-0587 The design and the user guide say what the review's diff shows, every tier passes, and the panel is looked at in light and dark

## Notes

The work is in the dashboard alone. flai is not changed: `flai stream diff` and the method `stream.diff` have carried each file's hunks since S-0041, and the panel needs nothing more. The review page already opened a file's hunks when its line was clicked, but nothing on the line said so, which is why a file read as a path and two counts.

How each criterion was verified:

| Criterion | By test | In a browser |
|-----------|---------|--------------|
| A toggle control reveals or hides the panel | `DiffView.svelte.test.ts`: the mark, `aria-expanded`, the panel present and absent, a panel per file | Closed, all open, and one closed again: two panels of three left |
| `-` in the margin of every removed line | `DiffView.svelte.test.ts` and `review.test.ts`: the sign apart from the text, for every removed line | Read in the panel of `flaiover/src/lib/review.ts` |
| `+` in the margin of every added line | The same, for every added line | The same |
| Dim red behind removed lines, a brighter red outline round consecutive ones | `DiffView.svelte.test.ts`: one block per run, its lines' background and its border | Computed in Chrome: background `danger-soft`, a 1px solid border in `danger`, in light and dark |
| Dim green behind added lines, a brighter green outline round consecutive ones | The same, for added runs | Computed in Chrome: background `good-soft`, a 1px solid border in `good`, in light and dark |

The browser was Chrome, headless, on a scratch page that mounts `DiffView.svelte` with the dashboard's stylesheet and this story's own diff as `flai stream diff S-0164 --json` gave it. It was looked at closed and open in light and dark, and in a panel narrower than its longest line, where the outline goes round the whole run when scrolled. The scratch page and its images are in `.flai-cache/e2e-S-0164/`, which git ignores. Not done: the review page was not opened against the host's running flai serve, so the first time the panel is seen there is this story's own review, once the dashboard runs this branch's code.

The colours are the theme's: `good-soft` and `danger-soft` behind a line, `good` and `danger` for the outline and the text. In dark the outline is the brighter of the two. In light it is the stronger: a dark red or green round a pale one.

`scripts/flaiover-test.sh` passes: prettier, eslint, svelte-check, and vitest (628 tests). Of `make flai-test`, gofmt, vet, golangci-lint, the behaviour and integration tiers, and the smoke tier's template step pass. `flai check --strict` has 0 errors and 5 warnings, none from this story: E-0003, E-0010, E-0011, and E-0012 are done and not archived, and TH-0032 is answered on an archived story. The markdown lint has one error, in `wip/threads/TH-0035` on main (I-0027). Both stop the smoke tier.
