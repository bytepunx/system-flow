---
id: S-0380
type: story
nature: research
title: Research how a project keeps its design system in the repository, how the dashboard previews it live, and how agents work on it with the operator
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:44:14Z
updated: 2026-10-13T17:46:16Z
transitions: []
tags: [dashboard, cli]
topics: [design]
touches: [design/system/design-system-tool.md, design/system/README.md, design/adrs, design/adrs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0380 Research how a project keeps its design system in the repository, how the dashboard previews it live, and how agents work on it with the operator

## Goal

Settle, before anything is built, how E-0020's Design menu works: where and in what format a project keeps its design system so that flai, the dashboard, and agents read one source; how the dashboard previews a project's tokens live on sample Svelte and Tailwind components without building the project; and how the operator and agents change the design system together. The finding and the operator's decision become the base the epic's other stories are checked against.

## Acceptance criteria

- [ ] A finding, `design/system/design-system-tool.md`, sets out where and in what format a project keeps its design system, its tokens (colours for light and dark, fonts, type scale, spacing, radii, shadows) and its guidance, so that flai can validate it and the dashboard and agents read the same file, with the options considered and how each fits the documentation standard and the host methods, which read and save only markdown.
- [ ] It sets out how the dashboard previews a project's tokens live on sample components without building the project, isolated from the dashboard's own theme, with the options considered, such as an iframe or a shadow root carrying CSS custom properties over components compiled with the dashboard, or Tailwind's browser build.
- [ ] It sets out how the operator and agents work on the design system together, such as a proposal on a story's branch previewed beside the main branch's before acceptance, a thread on the document, or an edit from the page, and how a story's agent is told to build its UI from the design system.
- [ ] It ends with a recommendation, and the operator's decision is recorded on a thread and in ADRs.
- [ ] E-0020's other stories are checked against the decision, and what each must change is named on the epic's plan thread.

## Tasks
- T-1527 Record what the dashboard's theme, document editor, and site menu, and the prime pack, offer a design system tool today
- T-1528 Set out where and in what format a project keeps its design system, and how the dashboard previews it live, with the options and a recommendation
- T-1529 Set out how the operator and agents change the design system together, and how a story's agent is told to build from it, with a recommendation
- T-1530 Report the options to the operator on a thread and record the decision in the finding and in ADRs
- T-1531 Check E-0020's other stories against the decision and name on the epic's plan thread what each must change

## Notes

- Drafted by the planner for E-0020. The epic's other stories assume the recommendation the planner expects (see the epic's plan thread); this story confirms or corrects that.
