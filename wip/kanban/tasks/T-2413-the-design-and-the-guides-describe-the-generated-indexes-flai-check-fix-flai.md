---
id: T-2413
type: task
nature: improvement
title: The design and the guides describe the generated indexes, flai check --fix, flai adr index, the union-merged changelogs, and the acceptance notice that skips the shared paths
status: done
parent: S-0550
owner: alex
created: 2026-10-31T07:44:43Z
updated: 2026-10-31T10:20:42Z
transitions:
  - to: ready
    at: 2026-10-31T10:11:03Z
    by: agent-S-0550
  - to: in-progress
    at: 2026-10-31T10:11:05Z
    by: agent-S-0550
  - to: done
    at: 2026-10-31T10:20:41Z
    by: agent-S-0550
stream: S-0550
tags: [cli, template]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/continuous-improvement.md, design/system/workflow.md, design/conventions/decisions.md, template/root/design/conventions/decisions.md, template/CHANGELOG.md, docs/contributors/template.md, flai/cmd/check.go]
after: [T-2409, T-2397, T-2399, T-2402]
usage:
  source: log
  seconds: 576
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 129
      output: 45613
      cache_read: 7836167
      cache_write: 368156
      cost: 4.5135
---
# T-2413 The design and the guides describe the generated indexes, flai check --fix, flai adr index, the union-merged changelogs, and the acceptance notice that skips the shared paths

## Work

Bring the documents in line with what the other tasks built and the story's ADR decided.

- `design/system/flai-cli.md`: the `flai check` row gains `--fix` and the two generated files it rewrites. The `flai adr` rows gain `flai adr index`, and say that `new`, `accept`, and a supersede regenerate the index. The `overlapped` event and the shared-paths text say that the acceptance notice leaves out the shared paths, while the trial merge and `flai task done`'s message still report them.
- `docs/users/flai.md`: the same in "Check the repository", "Record a decision", "Record recurring friction", and "Shared paths". Then regenerate `docs/users/flai-reference.md` with `make flai-reference`, never by hand.
- `design/system/continuous-improvement.md` § Summary: `flai check` compares the summary with its folder, and `--fix` rewrites it.
- `design/system/workflow.md` § Branches and collisions: the acceptance notice and the shared paths, the generated indexes, and the union-merged changelogs, below the before figures the measuring task wrote there.
- `design/conventions/decisions.md` and `template/root/design/conventions/decisions.md`: where the convention says flai adds the row to the index, say the index is generated from the ADR files and never edited by hand. This is a baseline change, made in both copies in the same story.
- `template/CHANGELOG.md`: under `## Unreleased`, one entry for this story's template changes (`.gitattributes`, `git.md`, `decisions.md`, `design/adrs/README.md`), naming the ADR and the flai release it needs.

Fourth layer: it waits for the ADR index check task, the notice task, the changelog task (`template/CHANGELOG.md`), and the measuring task (`workflow.md`), so that it describes what was built.

## Done when

Every document named above says what the code does, with no hand edit to `flai-reference.md`. `flai test` on the changed paths passes, the markdown lint included, and `flai check --strict --story S-0550` reports nothing new. Story criteria 1 to 3, their documentation.

## Notes

Drafted by the planner of S-0550 on 2026-10-31.
