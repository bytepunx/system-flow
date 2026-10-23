---
id: D-0309
title: "E-0020 plan: research first, then the format in flai, the Design page, editing, proposals on story branches, the stylesheet, and the convention (seven draft stories): Recommendation: start the…"
kind: answer
sources: [TH-0438]
topics: [design, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: []
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0309 E-0020 plan: research first, then the format in flai, the Design page, editing, proposals on story branches, the stylesheet, and the convention (seven draft stories): Recommendation: start the…

The operator, alex, answered TH-0438 at 2026-10-15T02:11:46Z.

What was asked:

```text
**Recommendation:** finalize S-0380 now and leave the other six drafts until S-0380's decision. Its last task, T-1531, names on this thread what each draft must change, and the planner replans them from that.
```

The answer:

```text
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
```

Source: TH-0438.
