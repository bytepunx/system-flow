---
id: S-0381
type: story
nature: feature
title: A project's design tokens are a DTCG file the manifest names, beside its guidance, and flai validates them, renders them as CSS, reports their contrast, and serves them to the dashboard
status: in-progress
parent: E-0020
owner: alex
created: 2026-10-13T17:44:24Z
updated: 2026-10-15T00:29:49Z
transitions:
  - to: ready
    at: 2026-10-14T23:56:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-14T23:56:32Z
    by: agent-S-0381
tags: [cli, template]
topics: [design]
touches: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go, flai/internal/designsystem/resolver.go, flai/internal/designsystem/resolver_test.go, flai/internal/designsystem/color.go, flai/internal/designsystem/color_test.go, flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go, flai/internal/designsystem/starter.md, flai/internal/designsystem/starter.tokens.json, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/cmd/root.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go, design/system/flai-cli.md, design/system/project-manifest.md, design/system/documentation-standard.md, design/conventions/documentation.md, template/root/design/conventions/documentation.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/designsystem/testdata/full.tokens.json, flai/internal/designsystem/testdata/full.css.golden, flai/internal/designsystem/testdata/full.tailwind.css.golden, flai/internal/designsystem/testdata/show-missing.json.golden, flai/internal/designsystem/testdata/show-problems.json.golden, flai/internal/designsystem/testdata/show.json.golden, flai/internal/designsystem/testdata/themes.css.golden, flai/internal/designsystem/testdata/themes.tailwind.css.golden, flai/internal/designsystem/testdata/themes.tokens.json]
after: [S-0380]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1977
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-14
      hand_edits: 1
      work: 18
    - day: 2026-10-15
      hand_edits: 1
      work: 10
  models:
    - model: claude-opus-5-5
      input: 316
      output: 1520
      cache_read: 12070174
      cache_write: 924055
      cost: 6.0002
  strategic:
    - kind: planner
      seconds: 334
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 27523
          cache_read: 3647039
          cache_write: 144111
          cost: 2.433
    - kind: orchestrator
      seconds: 194
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 1074
          cache_read: 14080036
          cache_write: 32137
          cost: 3.5419
cost_of_delay:
  value: 319.15
  by: planner-S-0381
  at: 2026-10-14T23:54:31Z
forecast:
  duration: 1h15m
  delivery: 2026-10-15T01:19:00Z
  basis: "flai forecast's 95 s per unit of size over 54 done large feature stories on claude-opus-5-5 gave 56m at size 35; raised to 1h15m for a new package with a line-aware JSON reader, a resolver subset, and colour arithmetic across ten tasks in five layers."
  by: planner-S-0381
  at: 2026-10-14T23:54:02Z
finalized:
  by: orchestrator
  at: 2026-10-14T23:56:15Z
---
# S-0381 A project's design tokens are a DTCG file the manifest names, beside its guidance, and flai validates them, renders them as CSS, reports their contrast, and serves them to the dashboard

## Goal

A project keeps its design tokens in a W3C Design Tokens Community Group (DTCG) file, format 2025.10, that the manifest names, `design/system/design-system.tokens.json` by default, and its guidance beside it in the markdown document `design/system/design-system.md` ([ADR-0146](../../../design/adrs/0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md)). Tokens are not front matter. flai owns the subset of the format it reads: it writes both files from starters, validates the tokens ([ADR-0148](../../../design/adrs/0148-a-design-token-s-colour-is-in-the-srgb-colour-space-with-its-hex-or-in-oklch.md)), resolves light and dark through the Resolver's `theme` modifier, renders them as CSS custom properties and as a Tailwind v4 `@theme` block, reports the colour roles' contrast, and answers all of it to the dashboard, so that the dashboard and agents read the same tokens. Reading and saving the token file from the dashboard, `file.get` and `file.save` ([ADR-0147](../../../design/adrs/0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)), is the next story's, not this one's.

## Acceptance criteria

- [ ] The manifest key `design.tokens` names the token file, `design/system/design-system.tokens.json` when unset, and `design/system/project-manifest.md` and `docs/operators/settings.md` describe it.
- [ ] `flai design init` writes the token file and `design/system/design-system.md` from starters: every required colour role in light and dark, the dark values through the Resolver's `theme` modifier with the contexts `light` and `dark`, the `sans` and `mono` font stacks, the type scale `xs` to `4xl`, the spacing step, the radii and shadows `sm`, `md`, and `lg`, in groups named after Tailwind v4's theme namespaces, and the guidance's headings; it refuses when either file exists.
- [x] `flai check` reports as `design.tokens`, at the token's path and line: a token file that does not parse, a type flai does not read, a colour space other than `srgb` with its `hex` or `oklch`, alpha outside shadows and extra tokens, a missing role or step, and a resolver other than the `theme` modifier with `light` and `dark`; it passes a project with no token file.
- [ ] `flai design css` prints the tokens as CSS custom properties, one block per theme, and with `--tailwind` as a Tailwind v4 `@theme` block; `flai design show --json` gives the token file's path, the tokens of each theme, the guidance's headings, the problems, both renderings, and each role pair's contrast in light and dark against ADR-0148's thresholds, saying when an `oklch` colour was clipped to sRGB. Contrast is never a finding.
- [ ] The host method `design.get` answers what `flai design show --json` prints, and says when the project has no token file.
- [x] `design/system/documentation-standard.md` and both copies of `documentation.md`, the template's and this repository's, carry the one exception to "Everything is markdown": a project's design token file, in DTCG format, named by the manifest.
- [ ] `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated reference say how.

## Tasks
- T-1532 flai/internal/designsystem reads the DTCG token file and the guidance's headings, and reports each token it does not read at its path and line
- T-1533 The design system's tokens render as CSS custom properties per theme and as a Tailwind v4 @theme block, and Show gives the tokens with both renderings
- T-1534 flai check reports a design system token that does not parse as design.tokens with its key and line
- T-1535 flai design init writes the token file and the guidance from starters, and flai design show and flai design css print them
- T-1536 The host method design.get answers the dashboard what flai design show --json prints
- T-1537 The flai design and user docs and the reference describe the token file, the guidance, flai design, design.tokens, and design.get
- T-1590 The manifest key design.tokens names the design token file, design/system/design-system.tokens.json when unset
- T-1591 The documentation convention and standard allow one exception to Everything is markdown: a project's design token file, in DTCG format, named by the manifest
- T-1592 The design token file's light and dark values resolve through the DTCG Resolver's theme modifier, and a resolver other than theme with light and dark is a problem at its path and line
- T-1593 A design token colour converts from oklch to sRGB with clipping reported, and the colour roles' pairs get their contrast against ADR-0148's thresholds

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends tokens in the front matter of one markdown document under `design/system`; the story is checked against S-0380's decision before it is finalized.

### Planning

- **Replanned by planner-S-0381 on 2026-10-14** against S-0380's decision (ADR-0146 to ADR-0148, `design/system/design-system-tool.md` § Epic › S-0381) and the operator's answer on TH-0439: the tokens are a DTCG file, not front matter. The title, goal, and criteria were rewritten to match, since the story is a draft. Tag `template` added, for T-1591.
- **Touches.** None declared before the first plan. All are the planners', every one a file:
  - `flai/internal/designsystem/` `tokens.go`, `resolver.go`, `color.go`, `css.go`, each with its test, `starter.md`, `starter.tokens.json`: layout, a new package, and design (§ Epic names `resolver.go`, `color.go`, and `starter.tokens.json`). Criteria 2 to 4.
  - `flai/internal/manifest/manifest.go` with its test, `design/system/project-manifest.md`, `docs/operators/settings.md`: design (§ Epic), and co-change (`settings.md` at 19% in `flai touches suggest`; its test fails on a manifest key with no row). Criterion 1.
  - `flai/internal/check/design.go` with its test, and `check.go`, whose `Run` calls each rule: layout. Criterion 3.
  - `flai/cmd/design.go` with its test, and `flai/cmd/root.go`, which registers each command group: layout and co-change. Criteria 2 and 4.
  - `flai/internal/hostapi/design.go` with its test, and `hostapi.go`, whose method table merges each group's methods: layout. Criterion 5.
  - `design/system/documentation-standard.md`, `design/conventions/documentation.md`, `template/root/design/conventions/documentation.md`: design (ADR-0146 § Decision 5). `template/CHANGELOG.md`: co-change (8%), since every template change gets an entry. Criterion 6.
  - `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: design, as `documentation.md` asks for a new command. Criterion 7.
- **Co-change left out.** `docs/operators/index.md`, the dashboard docs, `design/adrs/README.md`, and `flai/internal/hostapi/writes.go`: the story adds no runbook, no dashboard view, no ADR, and no write.
- **No folder touch.** If T-1592 finds the spec needs a resolver document beside the token file, its starter adds one file under `flai/internal/designsystem/`, which `flai task done` widens.
- **Layers.** T-1532, T-1590, and T-1591; then T-1592 and T-1593; then T-1533 and T-1534; then T-1535 and T-1536; then T-1537.
- **Forecast: 1h15m, delivery 2026-10-15T01:19:00Z.** `flai forecast` gave 56m: 95 s per unit of size over 54 done large feature stories on claude-opus-5-5, times size 35 (5 criteria, 30 touches), and about 59m with the seven criteria now written. Raised to 1h15m: a new package with a JSON reader that keeps each token's line, a subset of a new spec's Resolver module, and colour arithmetic checked against flaiover's test, across ten tasks in five layers, which a rate per touch does not weigh. The delivery adds the 19m to flai's 01:00Z.
- **Cost of delay: 319.15 USD a week.** `flai cod`: 1h15m of the 4h42m forecast over E-0020's six open stories, of the epic's 1200 USD a week (8h lost per cycle). It stands. It falls when the file story is drafted and takes its share.
