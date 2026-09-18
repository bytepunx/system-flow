---
id: T-0149
type: task
nature: improvement
title: User docs and dashboard wording for the new ready rule
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:30Z
updated: 2026-09-18T18:15:03Z
transitions:
  - to: ready
    at: 2026-09-18T18:14:32Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:14:32Z
    by: alex
  - to: done
    at: 2026-09-18T18:15:03Z
    by: alex
stream: S-0049
tags: []
touches: [docs/users, flaiover]
---

# T-0149 User docs and dashboard wording for the new ready rule

## Work
Update `docs/users/flai.md` (the `flai move S-0001 ready` comment and any prose on the definition of ready), the flaiover guide under `docs/users`, and `design/system/flai-cli.md` if it restates the rule. Search `flaiover/src` for the old rule text; the board shows the message flai returns, so change dashboard code or tests only where they pin the old wording.

## Done when
- No user-facing document says ready needs a task; the review rule is documented where the ready rule was
- flaiover tests pass if any flaiover file changed

## Notes
