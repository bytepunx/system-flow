---
id: T-3632
type: task
nature: improvement
title: flai check advises docs.unstarted-story on a design or docs line that names a story in backlog or ready, with an ADR and a test that reproduces I-0414's list
status: backlog
parent: S-0856
owner: alex
created: 2026-11-27T23:58:28Z
updated: 2026-11-28T00:01:09Z
transitions: []
stream: S-0856
tags: [flai, check]
touches: [flai/internal/check/check_docs.go, flai/internal/check/docs_stories.go, flai/internal/check/docs_stories_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, design/adrs/README.md, design/system/documentation-standard.md, design/system/flai-cli.md, docs/users/flai.md, design/system/agent-context.md]
---
# T-3632 flai check advises docs.unstarted-story on a design or docs line that names a story in backlog or ready, with an ADR and a test that reproduces I-0414's list

## Work

- Write the ADR with `flai adr new`: `flai check` advises `docs.unstarted-story` on a line of a markdown file under the design or docs folder, ADRs, `README.md` indexes, and the folders `documentation()` skips left out, that names a story whose status is `backlog` or `ready`, outside fenced code and inline code spans. It is advisory, so `--strict` passes over it, as `wip.after-reason` is ([ADR-0230](../../../design/adrs/0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md)). Its reason: a living document says what is, and a list of stories still to come is the line every later sibling rewrites (I-0414, I-0387). The ADR's number is not known until it is written, since T-3618 of S-0852 and T-3641 of S-0854 also plan an ADR, so its file is left out of the touches: the story's `design/adrs` folder touch covers it, and `flai task done` widens the touches to the file written.
- Add the rule in `flai/internal/check/docs_stories.go`, called from `documentation()` in `check_docs.go` for each design and docs file that is not an ADR, through `c.advise`. The message names the story, its status, and says to describe the present and leave the story's plan to the board.
- In `flai/internal/check/scope.go`, drop the finding from a check scoped to a story when it is outside the story, as `item.archive` and the others there are, so that a close-out records no issue for another story's document. Cover it in `scope_test.go`.
- In `docs_stories_test.go`, reproduce I-0414: a design document line "and, still to come, the acceptance's gates (S-0808)" with S-0808 in backlog gives the finding; the same line with S-0808 done, a code span, a fenced block, and an ADR give none.
- Run the rule on this repository and settle each finding: today `design/system/agent-context.md` names S-0241, a backlog story, as the repeat measurement still to run; say what was measured and leave the future to the board. A finding in `design/system/workflow.md` or `docs/operators/runbooks/change-the-board.md` is T-3627's to settle, not this task's.
- Document the rule: `design/system/documentation-standard.md` § Content rules, the rules of `flai check` in `design/system/flai-cli.md`, and `docs/users/flai.md` § Check the repository, where the count of warnings `--strict` passes over goes up by one.
- Waits for no task: it runs beside the rewrite of the hub sentence, and neither changes the other's files.

## Done when

- `flai check` on a repository whose design document names a backlog story outside code gives one advisory `docs.unstarted-story` warning on that line, and `--strict` passes.
- The new test fails without the rule and passes with it.
- `flai check --strict` on this repository, once T-3627 is done too, gives no `docs.unstarted-story` finding.
- `flai test` on the changed paths passes.

## Notes
