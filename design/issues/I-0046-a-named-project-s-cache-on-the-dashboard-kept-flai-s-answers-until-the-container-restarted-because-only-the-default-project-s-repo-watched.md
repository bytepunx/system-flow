---
id: I-0046
title: A named project's cache on the dashboard kept flai's answers until the container restarted, because only the default project's Repo watched
class: defect
status: open
count: 1
first_reported: 2026-09-26T03:31:44Z
last_reported: 2026-09-26T03:31:44Z
updated: 2026-09-26T03:31:44Z
---

# I-0046 A named project's cache on the dashboard kept flai's answers until the container restarted, because only the default project's Repo watched

## Description

flaiover keeps one cache of flai's answers per project, a `Repo` made by `repo()` for the project key a request names. A cache is only forgotten when its `Repo` listens for flai's `change`, `connected`, and `gone` events, and only some `Repo`s were ever told to listen: the default project's, at startup, and the one `/api/events` or `/api/stats` was asked for. A page that asked `/api/board?project=<key>` without either of those having reached that project saw what flai said the first time until the container restarted.

## Instances

### 2026-09-26T03:31:44Z
/api/board?project=sf listed S-0028, S-0115, and S-0116 in done after they were accepted and archived on another machine and pulled here, while /api/board, flai board, and wip/kanban/stories showed none. repo() made a Repo per project key, but only hooks.server.ts's startup watch() on the default Repo, and /api/events and /api/stats for the project they were asked about, ever started one listening. Workaround: docker restart flaiover. Fixed in S-0117: repo() calls watch() on every Repo it makes.

## Remediation

S-0117: `repo()` calls `watch()` on every `Repo` it makes, so a `Repo` listens from the moment it exists, whichever route made it. `project-events.test.ts` pins it with a real connection: a named project's board is asked, its flai reports a removed story file, and the board is asked of flai again.
