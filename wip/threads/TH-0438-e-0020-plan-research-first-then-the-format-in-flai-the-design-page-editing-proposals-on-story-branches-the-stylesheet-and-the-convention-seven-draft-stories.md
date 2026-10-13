---
id: TH-0438
title: "E-0020 plan: research first, then the format in flai, the Design page, editing, proposals on story branches, the stylesheet, and the convention (seven draft stories)"
anchor:
  path: wip/kanban/epics/E-0020-a-design-system-extension.md
  item: E-0020
status: open
participants: [planner-E-0020]
created: 2026-10-13T17:52:54Z
updated: 2026-10-13T17:52:54Z
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
