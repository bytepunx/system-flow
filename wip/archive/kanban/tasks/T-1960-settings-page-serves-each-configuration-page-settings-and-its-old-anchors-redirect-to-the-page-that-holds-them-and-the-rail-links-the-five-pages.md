---
id: T-1960
type: task
nature: improvement
title: "/settings/<page> serves each Configuration page, /settings and its old anchors redirect to the page that holds them, and the rail links the five pages"
status: done
parent: S-0459
owner: alex
created: 2026-10-21T07:03:19Z
updated: 2026-10-21T09:11:11Z
transitions:
  - to: ready
    at: 2026-10-21T09:05:01Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T09:05:02Z
    by: system-flow
  - to: done
    at: 2026-10-21T09:11:11Z
    by: agent-S-0459
stream: S-0459
tags: [dashboard]
touches: [flaiover/src/routes/settings/+page.svelte, flaiover/src/routes/settings/+page.ts, "flaiover/src/routes/settings/[page]/+page.svelte", "flaiover/src/routes/settings/[page]/+page.ts", flaiover/src/routes/settings/settings.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts]
after: [T-1959]
usage:
  source: log
  seconds: 369
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 77
      output: 31033
      cache_read: 4147257
      cache_write: 140449
      cost: 2.2981
---
# T-1960 /settings/<page> serves each Configuration page, /settings and its old anchors redirect to the page that holds them, and the rail links the five pages

## Work

- Add `flaiover/src/routes/settings/[page]/+page.ts`, which throws a 404 for a slug that `settingspages.ts` does not hold. Add `+page.svelte`, which renders the page's title and `<SettingsPanel page={slug} />`.
- Make `/settings` go to `/settings/features`. An anchor such as `#harnesses` is read in the browser, where the server never sees it, so `/settings/+page.svelte` reads the hash and goes to the page that `settingspages.ts` maps it to, with `replaceState` so that back skips it.
- In `flaiover/src/lib/sitemenu.ts`, point the five Configuration entries at `resolve('/settings/[page]', …)`. Update `sitemenu.test.ts` so that `locate` names each entry on its page.
- Write `settings.svelte.test.ts`: each slug renders, an unknown slug is a 404, and the redirect from `/settings` works with and without an anchor.
- Waits for T-1959, the `page` prop and the map these routes use.

## Done when

The new tests and `sitemenu.test.ts` pass, and each Configuration entry in the rail opens its own page as the current page.

## Notes
