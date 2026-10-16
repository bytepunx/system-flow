---
id: E-0021
type: epic
nature: feature
title: A project may accept a story through a pull request and its CI instead of a local merge, with the checks and review comments read back to the story's agent
status: ready
owner: alex
created: 2026-10-16T04:27:54Z
updated: 2026-10-16T05:25:43Z
transitions:
  - to: ready
    at: 2026-10-16T05:25:43Z
    by: alex
tags: [cli, dashboard, template]
topics: [cli, orchestration, release, dashboard]
touches: [flai/internal/serve, flai/internal/hostapi, flai/internal/workitem, flai/cmd, flaiover/src, design/system/workflow.md, design/system/pushing-from-the-board.md, design/system/metrics.md, design/conventions/git.md, docs/operators/index.md, docs/users/flai.md, template]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 18818
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 98
          output: 1593
          cache_read: 38628779
          cache_write: 20867
          cost: 9.689
---
# E-0021 A project may accept a story through a pull request and its CI instead of a local merge, with the checks and review comments read back to the story's agent

## Outcome

Acceptance today merges a story's branch into main on the host, moves it to done, archives it, and commits; nothing is tagged or pushed, and publishing is a later, batched step ([ADR-0032](../../../design/adrs/0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md), [workflow.md § Transitions and who makes them](../../../design/system/workflow.md#transitions-and-who-makes-them)). `git.md` names a pull request template, but nothing in flai opens one. A team whose main branch is protected, whose CI is the gate, or whose reviewers work on GitHub cannot adopt system-flow as it stands: Symphony lands pull requests with CI status as proof of work, Vibe Kanban opens them with descriptions, OpenCode ships a GitHub workflow. This epic gives a project the choice, and leaves local acceptance the default.

- [ ] A project chooses in `system-flow.yaml` whether a story is accepted by a local merge, as today, or through a pull request on its remote. The default is the local merge, and a project that does not choose sees no change; the setting is indexed in `docs/operators/settings.md`.
- [ ] With pull-request acceptance, a story's move to review pushes its branch and opens one pull request per story, or updates the one it has, from `.github/pull_request_template.md`: the story ID, its acceptance criteria, its definition of done, and its last `flai verify` result. The push and the request go through `flai serve` on the host with a token from the environment it runs in, never from the dashboard's container and never from the agent's own credentials ([pushing-from-the-board.md](../../../design/system/pushing-from-the-board.md), [ADR-0031](../../../design/adrs/0031-the-dashboard-s-container-holds-nothing-of-the-operator-s-the-host-pushes.md)).
- [ ] The pull request's checks and reviews are read back: the story's page in the dashboard and `flai show` carry each required check's state and each review's, and a failing check or a requested change reaches the story's agent as a thread on the story that names the check or the comment and links its log, so the agent fixes it and pushes again without the operator relaying.
- [ ] `flai accept` under pull-request acceptance merges through the remote, with the merge, squash, or rebase method the project sets, then does what it does today: done, archive, commit, with main fetched first. It refuses, naming them, while a required check fails or a requested review is outstanding, and the orchestrator's acceptance under [ADR-0093](../../../design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-fl.md) obeys the same gate. The acceptance commit's subject and the archive stay as they are, so metrics and history read the same in both modes.
- [ ] Time a story waits on its checks and on its reviewers is a wait of its own in `metrics.md`, defined by an ADR, reported by `flai stats` per week and in the dashboard's waiting chart beside the thread and review waits, so that a slow CI shows up where the analyzer looks.
- [ ] GitHub is the first forge, through `gh` or its REST API; the forge is an adapter behind a contract so that GitLab follows in a later epic without touching acceptance.
- [ ] Integration tests run only where a forge token and a test repository are configured and are skipped, saying so, where they are not, so no close-out or CI run without them fails for want of them.
- [ ] `docs/operators/index.md`, `docs/users/flai.md`, `design/system/workflow.md`, and `design/conventions/git.md` describe both modes, and the template carries the setting and the pull request template.

## Stories

To be drafted by the planner. Order: the setting and the forge contract first, with GitHub behind it; then opening and updating the pull request on the move to review; then reading checks and reviews back as threads; then acceptance through the remote and the gate, the orchestrator's included; then the wait in metrics and the dashboard; documentation and the template with each story.
- S-0392 A project chooses local or pull-request acceptance in system-flow.yaml, and flai reaches its GitHub remote through a forge contract with a token the host names
- S-0393 Under pull-request acceptance, a story's move to review has flai serve push its branch and open or update one pull request from the project's pull request template
- S-0394 A story's pull request checks and reviews show in flai show and on its dashboard page, and a failing check or a requested change reaches its agent as a thread that sends the story back
- S-0395 Under pull-request acceptance, flai accept and the orchestrator's acceptance refuse a story while a required check has not passed or a required review is outstanding, naming each
- S-0396 The time a story waits on its pull request's checks and on its reviewers is a wait of its own in metrics.md, reported by flai stats per week and charted beside the thread and review waits
- S-0397 Under pull-request acceptance, flai accept merges the story's pull request through the remote by the project's merge method, then fetches main and moves to done, archives, and commits as today

## Notes

- From the comparison of system-flow with Symphony, Vibe Kanban, Conductor, Paperclip, and OpenCode on 2026-10-16: the absence of a pull request path is the largest adoption blocker for any team that is not one operator on one machine.
- The planner reads first: `design/system/workflow.md` § Transitions and who makes them and § What waits on an acceptance; `design/system/pushing-from-the-board.md`; `design/conventions/git.md`; ADR-0032, ADR-0093, ADR-0031; `flai/internal/serve` for the push host action.
- Open to the designer: whether a story in review under pull-request acceptance still holds the review column's limit while it waits on CI, or whether a wait on the forge counts apart; whether the orchestrator may approve a pull request's review itself when `accept_reviews` is on, or only merge one a person approved.
- Secrets: the forge token lives in the variable the host names, as a provider's key does (ADR-0129), and is never written to a log, an item, or the start record.
