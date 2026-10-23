---
id: T-2086
type: task
nature: feature
title: "The guides say how a tunnel or custom domain is admitted: the Host page's Addresses area, the key, and the restart that applies it"
status: done
parent: S-0489
owner: alex
created: 2026-10-23T18:09:04Z
updated: 2026-10-23T19:09:07Z
transitions:
  - to: ready
    at: 2026-10-23T19:07:19Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T19:07:20Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T19:09:07Z
    by: agent-S-0489
stream: S-0489
tags: [flaiover, docs]
touches: [docs/users/flaiover.md, docs/operators/index.md]
after: [T-2082, T-2085]
usage:
  source: log
  seconds: 107
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 6273
      cache_read: 1159302
      cache_write: 37100
      cost: 0.5771
---
# T-2086 The guides say how a tunnel or custom domain is admitted: the Host page's Addresses area, the key, and the restart that applies it

## Work

`docs/users/flaiover.md`: in § Host, describe the Addresses area between the Dashboard and flai host areas: what the list is for, that the request's own address is always accepted, that saving needs `flai serve enable settings --all-projects`, that flai refuses an entry that is not a host name, and that the running container keeps its old list until Restart in the Dashboard area; and in § Settings, say the Configuration group has the five settings pages and the Host page. `docs/operators/index.md` § The tunnel expectation: add the step that admits the tunnel's or the domain's host name, `flai config set dashboard.allowed_hosts <host>` or the Host page, then `flai dashboard restart`, and what fails without it, the live events' 403 and the page's "Reconnecting…"; and in § Security posture, note that an allowed host widens only which origins the browser may open the events socket from, not who may log in. The install runbook names no tunnel, so it stays as it is.

Waits for T-2082 and T-2085, which change the same pages' rail and Host sections and settle the words the area shows.

## Done when

- `docs/users/flaiover.md` § Host describes the Addresses area and the restart that applies a change, and § Settings names the Host page in the Configuration group.
- `docs/operators/index.md` § The tunnel expectation gives the admitting step and names what fails without it.
- The markdown lint passes on the changed pages.

## Notes

Planned by the planner on 2026-10-23. The settings index rows for the key and the variable are T-2084's, because its test enforces them.
