---
id: T-3183
type: task
nature: feature
title: An ImageJob component shows one image job's state, step, live output, and a failure's cause with a Retry button
status: backlog
parent: S-0757
owner: alex
created: 2026-11-12T16:03:51Z
updated: 2026-11-12T16:03:58Z
transitions: []
stream: S-0757
tags: [dashboard]
touches: [flaiover/src/lib/components/ImageJob.svelte, flaiover/src/lib/components/ImageJob.svelte.test.ts, design/components/flaiover/src/lib/components/ImageJob.examples.json]
after: [T-3181]
---
# T-3183 An ImageJob component shows one image job's state, step, live output, and a failure's cause with a Retry button

## Work

Write `flaiover/src/lib/components/ImageJob.svelte`. It takes one job in the shape `images.get` gives (T-3181) and shows:

- its images, its way, its step, and its times;
- its output through `LogFeed.svelte`, following a running job;
- on failure, its cause and step in the error style of the design system, with a Retry button that raises an event.

Use the design tokens' classes only. Write its tests in `ImageJob.svelte.test.ts` for running, passed, failed, and the retry event. Write its examples, `design/components/flaiover/src/lib/components/ImageJob.examples.json`, one per state, as `design/components/README.md` says.

It waits for T-3181, whose job shape it shows. It shares no file with T-3182, so the two run together.

## Done when

- vitest passes on the component, and `flai check` reports no `components.examples` finding for it.
- prettier, eslint, and svelte-check pass.

## Notes
