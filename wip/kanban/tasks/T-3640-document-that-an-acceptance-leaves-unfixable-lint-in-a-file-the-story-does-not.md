---
id: T-3640
type: task
nature: improvement
title: Document that an acceptance leaves unfixable lint in a file the story does not change uncommitted, and that digest absorb and compact lint, then close I-0411
status: backlog
parent: S-0853
owner: alex
created: 2026-11-27T23:59:05Z
updated: 2026-11-27T23:59:05Z
transitions: []
stream: S-0853
tags: [flai, docs]
touches: [design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, template/CHANGELOG.md, design/issues/I-0411-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md, design/issues/summary.md]
after: [T-3636, T-3638, T-3639]
---
# T-3640 Document that an acceptance leaves unfixable lint in a file the story does not change uncommitted, and that digest absorb and compact lint, then close I-0411

## Work

Write down what T-3636, T-3638, and T-3639 built, linking T-3631's ADR:

- `design/system/flai-cli.md` § The acceptance gate and § The merge queue pass: what the acceptance's lint refuses for, what it leaves uncommitted and names, and when the preflight holds the queue.
- `design/system/workflow.md` § What waits on an acceptance (S-0296): an uncommitted file on main with unfixable lint that the story does not change no longer waits for the operator.
- `design/system/strategic-agents.md` § The curator: `absorb` and `compact` lint what they rewrite, as `new` and `edit` do.
- `docs/users/flai.md` § Accept and release, § When the orchestrator accepts, and § Absorb and compact the digest.
- `template/CHANGELOG.md`: one line under `## Unreleased`.

Then close I-0411 with `flai issue close I-0411 --reason` naming the ADR and the three changes, which regenerates `design/issues/summary.md`.

It waits for T-3636, T-3638, and T-3639: it documents what they built and closes the issue once they are done.

## Done when

- The four documents and the changelog say what the code does, with no sentence left that says unfixable lint in any uncommitted file blocks every acceptance.
- I-0411 is closed with its reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed paths.

## Notes
