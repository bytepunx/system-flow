---
id: T-2051
type: task
nature: improvement
title: The user docs describe an example's answers, and I-0164 is closed with what fixed it
status: done
parent: S-0476
owner: alex
created: 2026-10-22T21:00:56Z
updated: 2026-10-22T22:08:44Z
transitions:
  - to: ready
    at: 2026-10-22T22:06:58Z
    by: agent-S-0476
  - to: in-progress
    at: 2026-10-22T22:06:59Z
    by: agent-S-0476
  - to: done
    at: 2026-10-22T22:08:44Z
    by: agent-S-0476
stream: S-0476
tags: [flai, flaiover]
touches: [docs/users/flai.md, docs/users/flaiover.md, design/system/component-library.md, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/summary.md]
after: [T-2048]
usage:
  source: log
  seconds: 105
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 17024
      cache_read: 3244234
      cache_write: 77962
      cost: 1.4499
---
# T-2051 The user docs describe an example's answers, and I-0164 is closed with what fixed it

## Work

- `docs/users/flai.md`: under § List a project's components › Component examples, add `answers` to the example's keys with a short example of one answer, and under § Render a component live, say that a framed component's request to the app's own server is answered from the chosen example's answers when one matches, and 501 otherwise.
- `docs/users/flaiover.md`: under § Components, say the same for the page, and under § Examples, describe the Answers panel: what it holds, that it is saved with the example, and that the event log shows each request with the status it got.
- `design/system/component-library.md` § flaiover's own components › What the pass found: the paragraph on requests to the app's own server now says the fourteen components carry examples with answers, with a link to the ADR of T-2035.
- Close the issue from the story's worktree: `flai issue close I-0164 --reason "<what fixed it>"`, naming the ADR and that examples carry the answers the harness gives; it rewrites `design/issues/summary.md`.

This task waits for T-2048, the last of the behaviour, so the docs describe what was built and the issue closes on it.

## Done when

- [ ] `docs/users/flai.md` and `docs/users/flaiover.md` describe `answers`, the matching rule in one sentence, and the page's Answers panel
- [ ] `design/system/component-library.md` § What the pass found no longer says the fourteen components cannot be shown
- [ ] I-0164 is closed with a reason that says what fixed it, and `design/issues/summary.md` is regenerated
- [ ] The markdown lint and `flai check --strict` pass

## Notes

Planned by planner-S-0476. The story's second acceptance criterion is this task's close of the issue.
