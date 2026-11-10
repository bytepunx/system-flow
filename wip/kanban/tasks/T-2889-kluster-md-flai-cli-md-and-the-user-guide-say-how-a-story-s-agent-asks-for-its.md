---
id: T-2889
type: task
nature: feature
title: kluster.md, flai-cli.md, and the user guide say how a story's agent asks for its deploy and when its namespace goes
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:28:09Z
updated: 2026-11-10T08:28:09Z
transitions: []
stream: S-0660
tags: [cli]
touches: [design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-2884, T-2887, T-2882]
---
# T-2889 kluster.md, flai-cli.md, and the user guide say how a story's agent asks for its deploy and when its namespace goes

## Work

Write what this story built into the design and the user guide:

- `design/system/kluster.md`: the row of "What happens when" that tears a story down on its way back to in-progress gives way to the ADR's table. The Commands table and its guard bullet name the request. The flai serve section's looks, failures, and journal say how a request is taken, refused, queued, and answered, that a failed asked-for deploy opens no thread, and who asked. Each part S-0659 marked not yet built is marked built.
- `design/system/flai-cli.md`: the kluster rows and the guard prose name the request.
- `docs/users/flai.md`: "Run a project's kluster" and "Deploy a story, tear it down, and promote main" say how a story's agent asks for its deploy and teardown, what it is told, each refusal, the limit, and how long the namespace lives.

It waits for T-2884, T-2887, and T-2882, since it describes what they built, under the names they built it with.

## Done when

- Each of criterion 9's documents says every criterion's behaviour, with a link to S-0659's ADR.
- `flai check --strict` and the markdown lint are clean on the three files.

## Notes
