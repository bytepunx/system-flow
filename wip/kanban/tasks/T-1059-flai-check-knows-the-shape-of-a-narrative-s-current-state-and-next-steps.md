---
id: T-1059
type: task
nature: improvement
title: flai check knows the shape of a narrative's Current state and Next steps
status: backlog
parent: S-0271
owner: alex
created: 2026-10-06T22:50:01Z
updated: 2026-10-06T22:50:01Z
transitions: []
stream: S-0271
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go]
after: [T-1054]
---
# T-1059 flai check knows the shape of a narrative's Current state and Next steps

## Work

The goal's last sentence: "`flai check` then knows the sections' shape". Today `check.go` checks only that `narrativeSections` are present.

- Add a rule, such as `narrative.state`, through the parser T-1054 exports. It flags the narrative of a story in `in-progress` or `review` in either of these cases:
  - `## Current state` is empty, or still the template's placeholder.
  - `## Next steps` is empty, the placeholder, or not a list.
- Use the same test of "written" as `scripts/close-out.sh` uses, so that the two agree.
- Give the finding the line of the section, and say to fix it with `flai stream state`.
- Choose the severity so that no narrative active on the main branch today fails `flai check --strict`. Check them before choosing.

## Done when

- [ ] Tests in `check_test.go` show that:
  - an empty section and a placeholder section are flagged;
  - `## Next steps` that is not a list is flagged;
  - a written narrative, and a story in another state, are not flagged.
- [ ] `flai check --strict` passes on the main checkout.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner for S-0271. The shape is the planner's reading of the goal. The plan thread names it as an assumption.
