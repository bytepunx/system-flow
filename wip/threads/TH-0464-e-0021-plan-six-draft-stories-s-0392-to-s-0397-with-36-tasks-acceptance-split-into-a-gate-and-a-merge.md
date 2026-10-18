---
id: TH-0464
title: "E-0021 plan: six draft stories, S-0392 to S-0397, with 36 tasks; acceptance split into a gate and a merge"
anchor:
  path: wip/kanban/epics/E-0021-a-project-may-accept-a-story-through-a-pull-request-and-its-ci-instead-of-a-local-merge-with-the-checks-and-review-comments-read-back-to-the-story-s-agent.md
  item: E-0021
status: resolved
participants: [planner-E-0021, orchestrator]
created: 2026-10-16T04:42:53Z
updated: 2026-10-18T03:43:31Z
---

# TH-0464 E-0021 plan: six draft stories, S-0392 to S-0397, with 36 tasks; acceptance split into a gate and a merge

On wip/kanban/epics/E-0021-a-project-may-accept-a-story-through-a-pull-request-and-its-ci-instead-of-a-local-merge-with-the-checks-and-review-comments-read-back-to-the-story-s-agent.md.

## Entries

### 2026-10-16T04:42:53Z planner-E-0021
I drafted six stories for E-0021, all in the backlog, with their tasks. Cost of delay values and two criteria wait on TH-0463.

## Stories and order

| Story | What | After |
|-------|------|-------|
| S-0392 | `acceptance.mode`, the forge contract, the GitHub adapter, the host's token variable, `flai forge check` | none |
| S-0393 | flai serve pushes the branch and opens or updates the pull request on the move to review | S-0392 |
| S-0394 | checks and reviews kept, shown in `flai show`, `item_get`, and the dashboard; a failure sends the story back with a thread | S-0393 |
| S-0395 | the acceptance gate on checks and reviews, the orchestrator's included | S-0394 |
| S-0396 | the checks and reviewers waits in `metrics.md`, `flai stats`, and the waiting chart | S-0394 |
| S-0397 | `flai accept` merges through the remote, then archives and commits as today | S-0395 |

S-0395 and S-0396 can run together: they share only shared paths.

## Tasks and layers

- **S-0392:**
  - first: T-1602, T-1603, T-1604, T-1605
  - then: T-1606 and T-1607
  - then: T-1608
- **S-0393:**
  - first: T-1609, T-1610, T-1611
  - then: T-1612
  - then: T-1613
  - then: T-1614
- **S-0394:**
  - first: T-1615 and T-1616
  - then: T-1617, T-1618, T-1619
  - then: T-1620 and T-1621
  - then: T-1622
- **S-0395:**
  - first: T-1623 and T-1624
  - then: T-1625 and T-1626
  - then: T-1627
- **S-0396:** T-1628, then T-1629, then T-1630 and T-1631, then T-1632.
- **S-0397:**
  - first: T-1633 and T-1634
  - then: T-1635
  - then: T-1636
  - then: T-1637

Each story writes its own ADR first, and its docs last.

## Assumptions

1. The GitHub adapter uses the REST API with the standard library, not `gh`, so that it needs nothing on `PATH` and tests against `httptest`.
2. A new host action, `forge`, off by default, gates every push, request, and merge.
3. The token is read from the variable the host names, `GITHUB_TOKEN` by default, and is removed from every agent's environment.
4. The push uses the host's git credentials, as publishing does. The API calls use the token.
5. A story's pull request is found by its head branch, so nothing new is written to the story's front matter.
6. flai serve reads checks and reviews into a state under `.flai-cache`. `flai show`, `item_get`, the dashboard, and the gate read only that state, so agents, which have no token, see the same thing.
7. A failing required check or a requested change makes flai send the story back from review to in-progress, with a thread, and start its agent. That is a new maker of a transition that is the operator's alone today. The alternative, fixing the branch while the story stays in review, would skip the close-out.
8. The forge gate binds the operator's acceptance too, not only the orchestrator's, unlike an open thread.
9. Where the token is not at hand, as for the orchestrator, `flai accept` asks flai serve to make the merge.
10. The acceptance commit and the archive still reach the remote only when published (ADR-0067).

## Proposals

- **Split, done in the drafts.** The epic's acceptance outcome is two stories: S-0395, the gate, and S-0397, the merge through the remote. The gate needs no token and can be accepted on its own; the merge is the git-heavy part.
- **Not drafted: publishing to a main branch that refuses direct pushes.** Under pull-request acceptance the code reaches main through the merge, but the archive commit is published by `flai release --pending`. A protected main refuses that push. I recommend a later epic, for example one where publishing opens a pull request of its own. It is outside this epic's outcome, which keeps the acceptance commit as it is. Say if you want it drafted here instead.
- **Left alone:** `design/system/pushing-from-the-board.md`, which the epic lists. It records S-0052's finding as history, and the stories change nothing it describes.
- **Later epic:** GitLab, behind the contract S-0392 defines.

## Figures

| Story | Forecast | Delivery |
|-------|----------|----------|
| S-0392 | 1h | 2026-10-16T05:41Z |
| S-0393 | 1h5m | 07:01Z |
| S-0394 | 1h15m | 08:25Z |
| S-0395 | 28m | 08:56Z |
| S-0396 | 42m | 09:12Z |
| S-0397 | 45m | 09:46Z |

Five forecasts are raised from flai's, each for the reason its story's Planning notes give. Cost of delay values follow your answer to TH-0463 question 1.

### 2026-10-16T04:58:19Z planner-E-0021
A note on E-0022, whose planner is drafting S-0398 to S-0400 now. Both epics change acceptance.

**Paths both epics claim:**

- S-0399 and S-0392 both change these files:
  - `flai/internal/manifest/manifest.go` and `settings.go`
  - `flai/internal/config/config.go`
  - `flai/cmd/serve_actions.go`
  - `template/root/system-flow.yaml.tmpl`
- S-0399 also changes `flai/internal/serve/serve.go`, as S-0393 and S-0394 do.
- S-0399, S-0400, S-0395, and S-0397 change these files:
  - `flai/cmd/accept.go`
  - `flai/internal/preview/accept.go`
  - `flai/internal/hostapi/writes.go`
  - `design/system/workflow.md`

**Effect on the board:** the stories of the two epics will hold each other on those files while either is in progress, so they run one after the other, whichever starts first.

**Design:** a green batch lands each story "through the acceptance of today" (S-0400). Under pull-request acceptance, that becomes S-0397's merge through the remote, so the queue would have to wait on each story's checks too. Whichever epic lands second should say how the two compose.

- **Recommended:** finish E-0021's S-0397 before E-0022's S-0400 starts. S-0400 can then land a batch through either mode, and S-0395's gate stays the one gate. Name S-0397 in S-0400's `after` if you agree.
- I have not changed E-0022's stories: they are another planner's.

### 2026-10-18T03:43:31Z orchestrator
Resolved: S-0397 was accepted
