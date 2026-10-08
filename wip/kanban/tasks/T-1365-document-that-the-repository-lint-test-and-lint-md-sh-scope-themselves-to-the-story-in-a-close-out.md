---
id: T-1365
type: task
nature: remediation
title: Document that the repository lint test and lint-md.sh scope themselves to the story in a close-out
status: backlog
parent: S-0345
owner: alex
created: 2026-10-08T08:40:03Z
updated: 2026-10-08T08:40:03Z
transitions: []
stream: S-0345
tags: [flai]
touches: [scripts/README.md, design/system/flai-cli.md]
after: [T-1363, T-1364]
---
# T-1365 Document that the repository lint test and lint-md.sh scope themselves to the story in a close-out

## Work

Say what T-1363 and T-1364 changed where builders and agents look for it.

- `scripts/README.md`: the `lint-md.sh` row says that with `CLOSE_OUT_STORY` set and no files given it leaves out the `wip/` files the story does not change, as the `check.sh` row says of `flai check`.
- `design/system/flai-cli.md`, the `flai verify` row: where it says `scripts/check.sh` and the tests that check the repository read `CLOSE_OUT_STORY` to scope themselves, name `scripts/lint-md.sh` and `TestRepositoryLintsClean` beside them, with S-0345.
- `docs/users/flai.md` already says the checks inside the tiers count only what is the story's; leave it unless the change makes it untrue.

Waits for T-1363 and T-1364, so it describes what they built rather than what was planned.

## Done when

- Both documents name the two scoped checks and agree with the code.
- `flai test scripts/README.md design/system/flai-cli.md` passes.

## Notes

Drafted by the planner for S-0345.
