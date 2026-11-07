---
id: T-2539
type: task
nature: improvement
title: PageSidebar and its state have no link entry, and a test fails if an entry takes an href again
status: done
parent: S-0559
owner: alex
created: 2026-11-07T05:08:22Z
updated: 2026-11-07T05:15:57Z
transitions:
  - to: ready
    at: 2026-11-07T05:11:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T05:11:55Z
    by: system-flow
  - to: done
    at: 2026-11-07T05:15:57Z
    by: agent-S-0559
stream: S-0559
tags: [flaiover]
touches: [flaiover/src/lib/pagesidebar.svelte.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/components/PageSidebar.svelte.test.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 242
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 9101
      cache_read: 3105926
      cache_write: 109176
      cost: 1.6387
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 21
          cache_read: 59927
          cache_write: 24605
          cost: 0.0207
---
# T-2539 PageSidebar and its state have no link entry, and a test fails if an entry takes an href again

## Work

Remove what S-0461 added for link entries, and only that, with the examples and the design that describe it, so that they change in one commit as `documentation.md` asks. It waits for nothing, since it is the change the rest of the story reports.

In `flaiover/src/lib/pagesidebar.svelte.ts`:

- Remove `href` from `SidebarEntry` and `'link'` from `SidebarKind`. `kindOf` returns the mode, or `label` without one.
- Remove `hashless` and its guard in `follow`. No page has only labels, and the item, review, and Docs pages all follow the hash.
- Correct the comments that name a link. A label is still never chosen, and the rest of the choice logic is unchanged.

In `flaiover/src/lib/components/PageSidebar.svelte`:

- Remove the link branch, its eslint comments, and `data-sidebar-link`.
- Keep the label, the `header` snippet, and the bound `current`. The Docs page uses the first two, and the item page binds the third.

In the two test files:

- Remove the `linked entries (S-0461)` tests and the `hashless` tests.
- Rewrite the tests that use `LINKS` only to reach a label or the bound `current`, over scroll or pane entries, so that those behaviours stay covered.
- Add the reproducing test: a `@ts-expect-error` on a `SidebarEntry` with an `href`, which `svelte-check` fails on if the field comes back.

In `design/components/flaiover/src/lib/components/PageSidebar.examples.json`, replace the example "Links to other pages, under labels" with one of scroll entries under labels, as the Docs page's **On this page** draws them.

In `design/system/flaiover-dashboard.md`, rewrite the PageSidebar paragraph: an entry is a pane, a scroll entry, or a label, the charts page's link entries are gone since this story, and `hashless` is not named. Bump its `updated`. The charts paragraph's history of S-0461 and S-0484 stays true and stays as it is.

## Done when

- `SidebarEntry` has no `href`, and no `link` kind, `hashless`, or `data-sidebar-link` is left under `flaiover/src`.
- The label over its children, the `header`, and the bound `current` each still have a test that passes.
- No example in `PageSidebar.examples.json` has an `href`, and `flai check --strict` reports no `components.examples` finding.
- `flai test` passes on the six files, the `flaiover-lint`, `vitest`, `svelte-check`, and markdown lint tiers among them.

## Notes
