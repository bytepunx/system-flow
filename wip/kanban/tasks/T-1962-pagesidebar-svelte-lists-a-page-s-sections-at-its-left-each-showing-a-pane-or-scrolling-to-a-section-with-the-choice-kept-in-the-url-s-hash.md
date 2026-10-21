---
id: T-1962
type: task
nature: feature
title: PageSidebar.svelte lists a page's sections at its left, each showing a pane or scrolling to a section, with the choice kept in the URL's hash
status: done
parent: S-0460
owner: alex
created: 2026-10-21T07:03:41Z
updated: 2026-10-21T09:23:52Z
transitions:
  - to: ready
    at: 2026-10-21T09:10:29Z
    by: agent-S-0460
  - to: in-progress
    at: 2026-10-21T09:10:29Z
    by: agent-S-0460
  - to: done
    at: 2026-10-21T09:23:52Z
    by: agent-S-0460
stream: S-0460
tags: [dashboard]
touches: [flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/components/PageSidebar.svelte.test.ts, flaiover/src/lib/pagesidebar.svelte.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json]
usage:
  source: log
  seconds: 803
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 98
      output: 73468
      cache_read: 5068056
      cache_write: 233017
      cost: 3.8691
---
# T-1962 PageSidebar.svelte lists a page's sections at its left, each showing a pane or scrolling to a section, with the choice kept in the URL's hash

## Work

- Write `flaiover/src/lib/pagesidebar.svelte.ts`. It holds the entries (`{ id, label, count?, mode: 'pane' | 'scroll', children? }`) and the current id as rune state, and keeps that id in step with `location.hash` both ways. For scroll entries it marks the section in view current through an `IntersectionObserver`, and falls back to the first entry where there is none, as under jsdom.
- Write `flaiover/src/lib/components/PageSidebar.svelte`:
  - a `<nav>` with an accessible name the page gives, and a list of entries with `aria-current` on the current one;
  - a pane entry sets the current pane, and the page renders only that pane through a snippet per id;
  - a scroll entry scrolls its section into view, respecting `prefers-reduced-motion`;
  - from `md` it is sticky beside the content, and below `md` it is a row above the content that scrolls within itself;
  - every colour and size comes from the token utilities.
- Test pane switching, the hash both ways, the counts, the children, and keyboard use.
- Write `PageSidebar.examples.json` with a set for each of these: panes with counts, scroll entries, and nested children.
- Waits for nothing in this story, which waits for S-0458.

## Done when

Both test files pass, and `flai check` reports no `components.examples` finding for it.

## Notes
