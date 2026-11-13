---
id: T-2886
type: task
nature: feature
title: The dashboard's review page shows the namespace's last deploy and checks run beside the criteria, with a link to each log
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:27:48Z
updated: 2026-11-13T14:51:48Z
transitions:
  - to: ready
    at: 2026-11-13T14:27:55Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T14:27:55Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T14:51:48Z
    by: agent-S-0662
stream: S-0662
tags: [dashboard]
touches: [flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-2881]
usage:
  source: log
  seconds: 1433
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 167
      output: 64630
      cache_read: 12417920
      cache_write: 531777
      cost: 6.6801
---
# T-2886 The dashboard's review page shows the namespace's last deploy and checks run beside the criteria, with a link to each log

## Work

Make `verify.status` in `flai/internal/hostapi/writes_host.go` answer the report's `kluster` part, which T-2881 adds. In `flaiover/src/lib/review.ts`, type it and turn it into a view. In `flaiover/src/lib/components/Review.svelte`, show it beside the criteria: the deploy and the checks run, each with its commit, outcome, who asked, and when, whether it is at the head, and a link to its log. A deploy step's log opens the Deployments page's Logs pane at the namespace; a checks log opens where the review page shows checks logs today.

Show nothing for a project with the action off or no `deploy`.

It waits for T-2881, whose record it shows.

## Done when

- `writes_test.go` shows `verify.status` with the part.
- `review.test.ts` and `Review.svelte.test.ts` cover deployed and checked at the head, at an older commit, failed, and absent, and each link.
- `flai test` passes on the changed files, prettier, eslint, vitest, and svelte-check among them.

## Notes

- Criterion 4 of S-0662. `Review.svelte` mentions no kluster today.
