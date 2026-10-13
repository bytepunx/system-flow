---
id: TH-0445
title: "S-0380 decision: where the design system lives, how it is previewed, and how you and agents change it (10 choices)"
anchor:
  path: wip/kanban/stories/S-0380-research-how-a-project-keeps-its-design-system-in-the-repository-how-the-dashboard-previews-it-live-and-how-agents-work-on-it-with-the-operator.md
  item: S-0380
status: open
participants: [agent-S-0380, orchestrator]
created: 2026-10-13T18:16:41Z
updated: 2026-10-13T18:17:05Z
---

# TH-0445 S-0380 decision: where the design system lives, how it is previewed, and how you and agents change it (10 choices)

On wip/kanban/stories/S-0380-research-how-a-project-keeps-its-design-system-in-the-repository-how-the-dashboard-previews-it-live-and-how-agents-work-on-it-with-the-operator.md.

## Entries

### 2026-10-13T18:16:41Z agent-S-0380
**Recommendation:** take all ten choices below as recommended. Reply "all as recommended", or give the numbers you want otherwise.

The finding is `design/system/design-system-tool.md`; each choice points at its section.

## What gets built, in order

1. **The document and its check** (S-0381). `design/system/design-system.md`: tokens in front matter, guidance in the body. A `design.tokens` rule in `flai check` makes a bad token refuse a dashboard save.
2. **The stylesheet** (S-0385). flai writes the tokens as a Tailwind v4 `@theme` file the manifest names, every time it writes the tokens. This comes before S-0383, whose save depends on it.
3. **The Design group and its live preview** (S-0382).
4. **Editing tokens on the page** (S-0383).
5. **Proposals and threads** (S-0384). This comes after S-0378 is accepted, since it uses S-0378's save on a story's branch.
6. **Agents' use** (S-0386). Prime and a baseline convention.

## Your choices (recommended first)

1. **Format and name** (§ Where a project keeps its design system).
   - **Recommended:** tokens in the front matter of `design/system/design-system.md`, guidance in its body.
   - Or: a YAML block in the body, a DTCG `.tokens.json`, or the UI's stylesheet as the source.
   - Why: only markdown fits `doc.get`, `doc.save`, prime, and `flai check` as they are. Figma and Style Dictionary can still get a generated DTCG export.
2. **Colour syntax.**
   - **Recommended:** six-digit hex and `oklch()`.
   - Or: hex only, or any CSS colour.
   - Why: flai can compute contrast for both with no new dependency.
3. **Colour roles.**
   - **Recommended:** a fixed set of required roles, with extra tokens allowed beside them.
   - The roles: `background`, `surface`, `text`, `text-muted`, `border`, `primary`, `on-primary`, `secondary`, `on-secondary`, `focus`, and `success`, `warning`, `danger`, `info`, each with a `-soft` tint.
   - Or: free names.
   - Why: sample components compiled into the dashboard can only style names they know.
4. **Contrast.**
   - **Recommended:** reported on the Design page and by `flai design`, not as a check finding.
   - Or: a check finding.
   - Why: any finding refuses a save, which would stop you saving a colour while you try it.
5. **Scope.**
   - **Recommended:** one design system per project now, and one document per UI sub-project later if you need it.
   - Or: one per UI sub-project from the start.
6. **Fonts in the preview.**
   - **Recommended:** system font stacks, plus Google Fonts loaded by family name only when the document asks for them.
   - Or: system stacks only.
   - Why: loading a Google font sends a request from your browser to Google.
7. **Preview isolation** (§ Previewing it live).
   - **Recommended:** a blank same-origin iframe that the page fills itself. The gallery is compiled with the dashboard, and the tokens are set as CSS variables, so a change shows at once with no build. If that fails, the fallback is a route outside the root layout.
   - Or: Tailwind's browser build (meant only for development), or your project's own build run on the host (large).
8. **Where your token save lands** (§ Working on it together).
   - **Recommended:** on main. flai writes the stylesheet in the same commit.
   - Or: only on a proposal story's branch, so that every code change goes through a story's close-out.
   - Why: it is the shortest path from your decision to the code. It refines ADR-0023 again.
9. **Working together.**
   - **Recommended:** edits on the page, proposals as stories previewed beside main's, a "request a change" form that makes a ready story, and threads on the document's headings.
   - Or: all of those plus a store of draft variants outside git.
   - Why: unsaved edits on the page already do what a store would, and agents cannot see a store.
10. **Agents' use** (§ How agents build from it).
    - **Recommended:** the generated stylesheet with `design.stale`, which names the document as well as the stylesheet. UI stories also get the whole document (not a brief, which lacks the tokens) and a baseline `design-system.md` convention, both selected by the document's `topics`. A lint for raw values in components is deferred.
    - Or: briefs and the convention only, or the lint now.

## What I do with your answer

- I record it in the finding's Decision section, with one ADR per decision.
- I name on TH-0438 what each of S-0381 to S-0386 must change for the planner to replan them. Two changes are known now: S-0383 must wait for S-0385, and S-0386 drops its separate list of governed sub-projects in favour of the document's `topics`.

Source: design/system/design-system-tool.md § Recommendation

### 2026-10-13T18:17:05Z orchestrator (recommendation)
**Recommendation:** reply "all as recommended". These are ten new design decisions, so they are yours, not mine.

I read the finding on `story/S-0380` (`design/system/design-system-tool.md`). Each choice follows from what flai and the dashboard do today:

- **1, 4, 8** rest on ADR-0023. `doc.save` takes only markdown, and it refuses any finding a save introduces. So tokens in markdown front matter work with no new host method. Contrast as a check finding would block every trial colour. A stylesheet left stale by a save would block the save too.
- **3 and 7:** the gallery is compiled with the dashboard, so it can only style role names it knows. A blank same-origin frame isolates it with nothing new in the container.
- **10:** the document's `topics` already decide which stories prime reaches. A second list of governed sub-projects in S-0386 could disagree with them.

**The one to weigh: choice 8.** A token save on main changes a code file, the UI's stylesheet, with no story and no close-out. ADR-0023 today commits only the one markdown path saved. If every code change must pass a story, pick the alternative: saves land only on a proposal's branch, and the page is for trying. I still recommend main. You review every change on this board anyway, and the stylesheet is generated, never hand-written.

**Choice 6:** a Google font is loaded only when the document names one. It sends a request from your browser to Google. Choose "system stacks only" if that matters to you.

Once you answer, S-0380's agent records the ADRs and names on TH-0438 what S-0381 to S-0386 must change. The planner then replans them, and I finalize those that come out consistent.

Source: design/adrs/0023-documents-are-saved-through-flai.md § Decision
