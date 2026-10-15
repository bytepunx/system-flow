---
id: TH-0438
title: "E-0020 plan: research first, then the format in flai, the Design page, editing, proposals on story branches, the stylesheet, and the convention (seven draft stories)"
anchor:
  path: wip/kanban/epics/E-0020-a-design-system-extension.md
  item: E-0020
status: resolved
participants: [planner-E-0020, agent-S-0380, orchestrator, alex]
created: 2026-10-13T17:52:54Z
updated: 2026-10-15T05:02:16Z
---

# TH-0438 E-0020 plan: research first, then the format in flai, the Design page, editing, proposals on story branches, the stylesheet, and the convention (seven draft stories)

On wip/kanban/epics/E-0020-a-design-system-extension.md.

## Entries

### 2026-10-13T17:52:54Z planner-E-0020
**Recommendation:** finalize S-0380 now and leave the other six drafts until S-0380's decision. Its last task, T-1531, names on this thread what each draft must change, and the planner replans them from that.

## Stories and order

| # | Story | Nature | after | Forecast | CoD/week |
|---|-------|--------|-------|----------|----------|
| 1 | S-0380 Research: where and in what format the design system lives, how the dashboard previews it live, how agents and the operator work on it | research | none | 1h | 230.77 |
| 2 | S-0381 The design system document: tokens in its front matter. flai validates it (`design.tokens`), renders it as CSS and as a Tailwind v4 `@theme`, and serves it (`design.get`) | feature, cli | S-0380 | 45m | 173.08 |
| 3 | S-0382 The Design group in the site menu, with a `/design` page: tokens laid out, plus a live, isolated gallery of sample components in light and dark | feature, dashboard | S-0381 | 50m | 192.31 |
| 4 | S-0383 Tokens edited on the page, with the preview following before saving and contrast flagged. Saved through `doc.save` | feature, dashboard | S-0382 | 30m | 115.38 |
| 5 | S-0384 Proposals: an agent's change is a story whose branch changes the document. The page previews it beside main's and offers a form that turns your request into a story | feature, both | S-0383 | 1h5m | 250.00 |
| 6 | S-0385 `design.stylesheet` in the manifest: `flai design css --write` keeps the project's Tailwind theme file, and `design.stale` catches drift | feature, cli | S-0381 | 27m | 103.85 |
| 7 | S-0386 Agents are primed with the design system for the sub-projects it governs. A baseline `design-system.md` convention says to build from the tokens and propose changes on the branch | feature, cli, template | S-0384, S-0385 | 35m | 134.62 |

Each value is flai's share of your 1200 USD a week (8h lost per cycle at 150 USD an hour), by forecast. Each Planning note gives the reason for each raise.

## Tasks and layers

- **S-0380:** T-1527 → T-1528 → T-1529 → T-1530 → T-1531. One per layer, since all edit the finding.
- **S-0381:** T-1532 → {T-1533, T-1534} → {T-1535, T-1536} → T-1537.
- **S-0382:** {T-1538, T-1539, T-1540} → T-1541 → T-1542.
- **S-0383:** {T-1543, T-1544} → T-1545 → T-1546 → T-1547.
- **S-0384:** {T-1548, T-1549, T-1550} → {T-1551, T-1552} → T-1553 → T-1554 → T-1555.
- **S-0385:** {T-1556, T-1557} → {T-1558, T-1559} → T-1560.
- **S-0386:** {T-1561, T-1562, T-1563} → {T-1564, T-1565}.

Every touch names a file. The only folder touch is `design/adrs` on S-0380, kept for ADRs not yet numbered.

## Assumptions (S-0380 confirms or corrects them)

1. A project's design system is one markdown document, `design/system/design-system.md`. Its tokens sit in front matter under Tailwind v4's theme namespaces, and its guidance sits in the body. That keeps it inside the documentation standard and the existing `doc.show` and `doc.save`, which take only markdown. Prime already briefs `design/system` by topic, so no new folder is needed.
2. The preview is an isolated frame (an iframe or a shadow root). Components compiled with the dashboard read CSS custom properties taken from the tokens, so a changed token restyles the gallery with no build of the project.
3. Working with agents uses the board itself: a proposal is a story whose branch changes the document, previewed before you accept it. There is no second store of drafts.
4. Tokens reach code through one generated stylesheet that flai checks, and agents are told so by a baseline convention.
5. The gallery is generic: buttons, inputs, a form, a card, an alert, a table, navigation, and type. It is not your project's own components.

## Proposed, not drafted

- **Add: a Patterns page.** It would preview whole sample screens (a form page, a list page, a dashboard) built from the gallery. I recommend waiting until S-0383 shows whether the component gallery is enough.
- **Add: this repository's own design system.** It would seed `design/system/design-system.md` from flaiover's `layout.css` theme and generate that theme with `design.stylesheet`. This dogfoods the epic. I recommend adding it after S-0385.
- **Drop S-0380 instead?** You could skip the research and finalize S-0381 onward on the assumptions above. That saves about 1h and an answer from you, but the format and preview choices would then be made in code with no ADR. I recommend keeping S-0380.

This plan needs no input from you before S-0380 can start.

### 2026-10-14T20:45:18Z agent-S-0380
S-0380's decision (TH-0445, ADR-0146 to ADR-0152) changes all six drafts. The tokens move to a DTCG file that the dashboard reads and saves through a new `file.get` and `file.save`. One new story builds that pair, and the order changes. The full reasoning is in § Epic of `design/system/design-system-tool.md` on `story/S-0380`.

## What each draft must change

### S-0381: the token file, the guidance, and their check

- Title and goal: the tokens are a DTCG file the manifest names, by default `design/system/design-system.tokens.json`. The guidance sits beside it in `design/system/design-system.md`. Tokens are no longer front matter.
- `init` writes both files from starters. Dark values come through the resolver's `theme` modifier.
- `design.tokens` reports, at the token's path and line:
  - a file that does not parse;
  - a colour space other than `srgb` with hex, or `oklch`;
  - a missing role or step;
  - a bad resolver.
- `show --json` adds each role's contrast. Contrast is reported, never a finding.
- New criterion and task: the exception to "Everything is markdown", in both copies of `documentation.md` and in `documentation-standard.md`. Add the `template` tag.
- Touches to add: `designsystem/resolver.go`, `color.go`, `starter.tokens.json`, `internal/manifest/manifest.go`, `project-manifest.md`, `docs/operators/settings.md`, and the three convention files.

### New story: flai file show and flai file save, and the host methods file.get and file.save

- It builds ADR-0147, with `--story` on both commands.
- It comes after S-0381, and it sets the file size limit that ADR-0147 leaves to it.
- Touches: `flai/cmd/file.go`, `root.go`, `internal/docedit/`, `hostapi/writes.go` and its tests, `dashboard-host-channel.md`, `flai-cli.md`, `docs/users/flai.md`, and the reference.

### S-0382: the Design page and its preview

- Contrast is shown for the roles' pairs, as `design.get` reports them, not for each colour against the ground. Extra tokens show as swatches and specimens.
- The preview is a blank same-origin frame that the page fills. The gallery's stylesheet is imported with `?inline`, the gallery is mounted with `mount`, and the tokens are set as `--ds-*` values. The fallback is a route outside the root layout.
- T-1540 checks first that `?inline` works and that events work inside the frame.
- New criterion: Google Fonts load only when a font token names a family. Add `flaiover/src/lib/design/gallery.css`.

### S-0383: editing on the page

- It saves through `file.save` on main, and flai writes the stylesheet in the same commit. Drop `doc.save` and the front-matter round trip from T-1543 and T-1546.
- Colours are entered as `srgb` with hex, or as `oklch`. The browser's contrast is tested against the same fixture as flai's.
- New criterion: unsaved values survive a reload.
- Add `api/design/tokens/+server.ts`, which answers 501 when the host's flai is older, and `agent.ts`.

### S-0384: proposals and threads

- A proposal's branch changes the token file or the guidance.
- The request form makes a story with fixed criteria, with the token file and the stylesheet as touches, and moves it to ready. The story is not a draft.
- New criteria, one task each:
  - tweak a proposal's tokens on its branch, through `file.save` with `story`;
  - list the guidance's threads by heading, and open one with `thread.new`.

### S-0385: the stylesheet

- `design.stale` names the token file as well as the stylesheet.
- The story settles how the generated file selects dark values.
- New criterion and task: every write of the token file, by `file.save` or a `flai design` command, writes the stylesheet in the same commit. `file.save` can reach the stylesheet's path.

### S-0386: agents' use

- Drop `init --project`, the governs list, and T-1561. The guidance's `topics` are the only list of which stories it reaches.
- Prime loads the guidance and the token file whole for the stories those topics reach, and briefs them above ADR-0050's threshold. The token file is prime's first read of a file that is not markdown.
- The convention gets `roles: [story, verify]` and ADR-0152's wording.
- Drop the `flai/cmd/design.go` touches. Add `context/story.go` and `conventions/conventions.go`.

## The order after the decision

1. S-0381, after S-0380.
2. The file story, after S-0381.
3. S-0382, after S-0381.
4. S-0385, after S-0381 and the file story.
5. S-0383, after S-0382, S-0385, and the file story.
6. S-0384, after S-0383.
7. S-0386, after S-0385 and the file story. It no longer waits for S-0384.

## The two proposed stories

- **Patterns page:** still wait for S-0383. It builds more on the same gallery and needs no new decision.
- **This repository's own design system:** do it after S-0385 and S-0386. It is larger than first proposed:
  - flaiover's `--t-*` colours must map onto the required roles;
  - `theme.test.ts` must read the generated file;
  - the dark selection must serve `[data-theme='dark']`.

### 2026-10-14T20:47:40Z orchestrator (recommendation)
**Recommendation:** start the planner on E-0020 (Plan on the dashboard, or `flai plan E-0020`). Have it replan S-0381 to S-0386 and the new file story from the list above. Have it also draft one story for the component design you asked for on TH-0445 at 20:34Z, which S-0380 did not cover.

## Why this needs you

- **I cannot start this planner.** `flai plan --candidates` lists an epic only while it is in the backlog with no stories, or when all its stories are done or cancelled. E-0020 has six open drafts. None of them has been edited since its plan, so no replan triggers on its own.
- **Component design is a change of scope.** You asked for component design after the ten choices, at 20:34Z. You wanted tools to evolve how components look, recorded so agents can refer to them, and previewed before implementation. S-0380 was accepted at 20:46Z without it. Its finding defers the project's own components: "Not now" (§ The project's own components). Its § Proposed, not drafted lists only a Patterns page and this repository's own design system.

## What I recommend for component design

- **A cheap first step, which the finding already names.** Component specs, as HTML snippets with the design system's classes, are kept beside the guidance. The Design page renders them in the same isolated frame as the gallery (ADR-0149). Agents are primed with them, like the guidance (ADR-0152). This needs no host build. It fits after S-0382 and the file story, beside the Patterns page.
- **The project's own built components** (a Storybook-like workshop on the host) need a host action and a build per change. They stay a later decision.
- **Or accept the deferral as it stands.** Say so here, and the planner replans only the six drafts and the file story.

I leave S-0381 to S-0386 as drafts until the planner replans them. Then I finalize those that come out consistent with ADR-0146 to ADR-0152.

Source: design/system/design-system-tool.md § The project's own components

### 2026-10-15T02:11:46Z alex
Confirmed the recommendation of 2026-10-14T20:47:40Z orchestrator.

Source: design/system/design-system-tool.md § The project's own components

### 2026-10-15T02:16:58Z orchestrator
Noted: you confirmed replanning E-0020 with a `file.get`/`file.save` story (ADR-0147) and a component-design story (component specs as snippets beside the guidance, previewed in the Design page's frame, primed to agents).

**One step is yours.** A confirmation does not start a planner, and `flai plan --candidates` does not list E-0020, so I cannot start it. Run `flai plan E-0020`, or use Plan on the dashboard. Its planner then drafts the two stories. S-0383 is waiting on the first one (TH-0454).

I settle the new drafts' plan threads and finalize the consistent ones once they come.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-15T05:02:16Z alex
Resolved: S-0386 was accepted
