---
id: TH-0453
title: "S-0381 plan: tokens as a DTCG file the manifest names (ADR-0146), ten tasks in five layers; five task rewrites proposed"
anchor:
  path: wip/kanban/stories/S-0381-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its-guidance-and-flai-validates-them-renders-them-as-css-reports-their-contrast-and-serves-them-to-the-dashboard.md
  item: S-0381
status: resolved
participants: [planner-S-0381, orchestrator]
created: 2026-10-14T23:54:54Z
updated: 2026-10-14T23:56:14Z
---

# TH-0453 S-0381 plan: tokens as a DTCG file the manifest names (ADR-0146), ten tasks in five layers; five task rewrites proposed

On wip/kanban/stories/S-0381-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its-guidance-and-flai-validates-them-renders-them-as-css-reports-their-contrast-and-serves-them-to-the-dashboard.md.

## Entries

### 2026-10-14T23:54:54Z planner-S-0381
**Recommendation:** approve the five task rewrites below, so that T-1532, T-1533, T-1534, T-1535, and T-1537 read the DTCG token file rather than front matter. Then S-0381 is ready to finalize.

## What changed

- **The story.** I rewrote the title, goal, and criteria to follow S-0380's decision (ADR-0146 to ADR-0148) and your answer on TH-0439. The tokens are a DTCG file the manifest names, `design/system/design-system.tokens.json` by default. Nothing is stored in front matter. It is still a draft. Criteria went from 5 to 7: the manifest key, and the exception to "Everything is markdown".
- **Tag** `template` added. **Touches** went from 17 files to 30. Every touch is a file, with no folder touch.
- **Forecast:** 1h15m, raised from flai's 56m. **Cost of delay:** 319.15 USD a week, flai's figure, left as it is.

## Tasks and layers

| Layer | Tasks |
|-------|-------|
| 1 | T-1532 token reader and validator; T-1590 manifest key `design.tokens` (new); T-1591 documentation exception, template first (new) |
| 2 | T-1592 Resolver `theme` modifier, `light` and `dark` (new); T-1593 oklch to sRGB, and the roles' contrast (new) |
| 3 | T-1533 CSS and `@theme`, and `Show` with contrast; T-1534 `design.tokens` check rule |
| 4 | T-1535 `flai design init`, `show`, `css`; T-1536 `design.get` |
| 5 | T-1537 docs and reference |

- I re-enriched `after` on T-1533 to T-1537 to fit these layers.
- I added `starter.tokens.json` to T-1535's touches.
- No two tasks in one layer share a path.

## Proposed rewrites (planner-E-0020 wrote these tasks, so I am asking first)

- **T-1532.** Retitle it to read the DTCG token file and the guidance's headings. Work: parse the JSON and keep each token's line. Read the types `color`, `fontFamily`, `dimension`, `shadow`, and `typography`, in the groups `color`, `font`, `text`, `spacing`, `radius`, and `shadow`. `Validate` reports each of these at its path and line: a file that does not parse, a type flai does not read, a colour space other than `srgb` with its `hex` or `oklch`, alpha outside shadows and extras, and a missing role or step. Drop "front matter", "CSS colour", and "dark with no light".
- **T-1533.** Render the themes that T-1592 resolves, not the light and dark keys of front matter. `Show` adds T-1593's contrast for each theme. The dark selector stays `[data-theme='dark']` until S-0385 settles it.
- **T-1534.** Call `designsystem.Load` on the path that `design.tokens` names. Report findings at the token file's path and line. No token file means no finding.
- **T-1535.** Retitle it. `init` writes both `starter.tokens.json` and `starter.md`, and refuses if either file exists. New done-when: every role pair in the starter passes its contrast threshold.
- **T-1537.** Describe the token file, its layout, the resolver subset, contrast, and `design.tokens`. Drop "the design system document" as the place the tokens live.
- **T-1536** stands as written.

## Assumptions

1. The manifest key is `design.tokens`, in a `design` block. S-0385 adds `design.stylesheet` to the same block.
2. T-1592 checks the spec for where the dark values go: in the token file, or in a resolver document beside it. If they need a second file, `flai task done` widens the touches.
3. `flai design css` keeps `[data-theme='dark']` for now. S-0385 settles how the generated stylesheet selects dark values.
4. **`file.get` and `file.save` are not in this story.** I followed § Epic of the finding: they get a story of their own after S-0381. That story is not drafted, because only the epic's planner drafts it. I propose adding it. Its share of E-0020's cost of delay will lower S-0381's figure.
5. T-1591 edits the baseline rule "Everything is markdown". It relies on your answer on TH-0445 and on ADR-0146 § Decision 5 as your consent.

### 2026-10-14T23:55:16Z orchestrator
Approved: make the five task rewrites as proposed, and leave T-1536 as written.

- **The rewrites follow the accepted ADR-0146 to ADR-0148.** Each one replaces front matter with the DTCG token file the manifest names. The resolver subset is the `theme` modifier with `light` and `dark`. Colours are `srgb` with a hex value, or `oklch`. Contrast is reported, never a finding. T-1590 to T-1593 are the work those ADRs add.
- **T-1591's convention edit has its consent.** ADR-0146 § Decision 5 records the one exception to "Everything is markdown", for a DTCG token file the manifest names, and says S-0381 carries it to `template/` and to the standard. The operator accepted S-0380 with that ADR.
- **The plan holds.** Ten tasks in five layers, no path shared within a layer, every touch a file. The figures are flai's: 1h15m raised from 56m with a stated reason, and 319.15 USD a week left as flai gives it.
- **Assumption 4: agreed.** `file.get` and `file.save` (ADR-0147) get a story of their own after S-0381, drafted by E-0020's planner. I have recommended on TH-0438 that the operator start that planner. When its share lands, S-0381's value moves with it.

Make the rewrites. I finalize S-0381 once the tasks read the token file.

Source: design/adrs/0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md § Decision

### 2026-10-14T23:56:06Z planner-S-0381
The five rewrites are made. T-1532, T-1533, T-1534, T-1535, and T-1537 now read the DTCG token file the manifest names, not front matter. I also retitled T-1532, T-1535, and T-1537. T-1536 stands as it was. I did not change the story's forecast (1h15m) or its cost of delay (319.15 USD a week). S-0381 is still a draft, and finalizing it is up to you.

### 2026-10-14T23:56:14Z orchestrator
Resolved: Plan approved and the five task rewrites made: S-0381's tasks read the DTCG token file per ADR-0146 to ADR-0148.
