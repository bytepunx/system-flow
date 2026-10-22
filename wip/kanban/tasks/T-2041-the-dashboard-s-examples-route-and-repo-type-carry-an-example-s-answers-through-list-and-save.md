---
id: T-2041
type: task
nature: improvement
title: The dashboard's examples route and repo type carry an example's answers through list and save
status: backlog
parent: S-0476
owner: alex
created: 2026-10-22T20:59:41Z
updated: 2026-10-22T21:01:35Z
transitions: []
stream: S-0476
tags: [flaiover]
touches: [flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/lib/server/repo.ts]
after: [T-2035]
---
# T-2041 The dashboard's examples route and repo type carry an example's answers through list and save

## Work

- In `flaiover/src/lib/server/repo.ts`, add `answers?: ComponentAnswer[]` to `ComponentExample`, with a `ComponentAnswer` type (`path`, `method?`, `status?`, `headers?`, `body?`, `events?: { event?: string; data: unknown }[]`) as T-2035's ADR defines it, so `components()` passes what the host lists.
- In `flaiover/src/routes/api/components/examples/+server.ts`, `exampleOf` rebuilds a saved example from `name`, `props`, `theme`, and `width` only, so `answers` is dropped on save: carry it when the request gives an array, and write it after `width` with the file's two-space indentation. Leave the key out when the array is empty or absent, as `theme` and `width` are left out.
- Tests in `examples.test.ts`: a PUT with answers writes them to the file and keeps the other examples' answers; a PUT without leaves none.

This task waits for T-2035 for the key's name and shape; it shares no file with T-2038 and T-2040 and runs beside them.

## Done when

- [ ] `ComponentExample` in `repo.ts` types `answers`
- [ ] `PUT /api/components/examples` keeps an example's `answers` in the file and leaves the key out when there are none
- [ ] `examples.test.ts` covers both, and vitest, prettier, and eslint pass on the changed files

## Notes

Planned by planner-S-0476.
