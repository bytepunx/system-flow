---
id: S-0393
type: story
nature: feature
title: Under pull-request acceptance, a story's move to review has flai serve push its branch and open or update one pull request from the project's pull request template
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:34:50Z
updated: 2026-10-16T05:24:55Z
transitions: []
tags: [cli, template]
topics: [release]
touches: [flai/internal/forge/pullrequest.go, flai/internal/forge/pullrequest_test.go, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/serve.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/forge.go, flai/cmd/forge_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, design/adrs, design/system/workflow.md, design/system/flai-cli.md, design/conventions/git.md, template/root/design/conventions/git.md, template/root/.github/pull_request_template.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md]
after: [S-0392]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h5m
  delivery: 2026-10-16T07:45:00Z
  basis: "Its own forecast of 1h5m; 3rd in the pull order with an in-progress limit of 5, behind S-0241 and S-0392."
  by: flai
  at: 2026-10-16T05:24:55Z
---
# S-0393 Under pull-request acceptance, a story's move to review has flai serve push its branch and open or update one pull request from the project's pull request template

## Goal

A team that reviews on GitHub sees each story arrive as a pull request without anyone opening it. When a story under pull-request acceptance moves to review, `flai serve` on the host pushes the story's branch and opens one pull request for it, or updates the one it has. The request's body comes from the project's pull request template: the story, its criteria, its definition of done, and its last `flai verify` result. The push and the request use the host's credentials and the token the host names, never the dashboard's container and never the agent's own.

## Acceptance criteria

- [ ] Under `acceptance.mode: pull-request`, with the new host action `forge` on, `flai serve` pushes `story/S-nnnn` to `origin` for each story in review whose branch on the remote is missing or behind. It then opens a pull request against the main branch, titled `[S-nnnn] <title>`, or updates the one the branch already has, so a story has one pull request however often it returns to review. With the mode `local`, or with the host action off, nothing is pushed.
- [ ] `flai forge status S-nnnn` prints the story's pull request, its URL, and whether the remote branch matches the local one. When nothing was pushed, it says why: the mode, the host action, the token, or a failure.
- [ ] The pull request's body fills `.github/pull_request_template.md` with the story's ID, its acceptance criteria, the definition of done, and the result of the story's last `flai verify`, with its commit. A project without that file gets a body with the same four parts. The template's `pull_request_template.md` gains the places for them.
- [ ] A story branch that a sync rebased is pushed with `--force-with-lease` naming the head flai last pushed. No other branch is ever forced. This criterion stands as TH-0463 answers question 4.
- [ ] A push or a request that fails, for example an unset token, a refused push, or the forge being down, opens one thread on the story naming the reason and what clears it. flai serve tries again at its next look, and resolves the thread once the push and the request succeed.
- [ ] `flai forge push S-nnnn` does the same on the operator's word. `flai guard` refuses it to every agent. Every agent `flai serve` starts runs without the forge token's variable in its environment.
- [ ] An ADR records the push and the request on the move to review, the host action, and the exception for story branches. `design/system/workflow.md`, `design/conventions/git.md` and the template's, `design/system/flai-cli.md`, `docs/users/flai.md`, the reference, `docs/operators/index.md`, and `docs/operators/settings.md` describe them.

## Tasks
- T-1609 A story's pull request title and body are built from the project's pull request template, its criteria, the definition of done, and its last flai verify result
- T-1610 The forge host action exists, agents start without the forge token's variable, and flai guard refuses them flai forge push
- T-1611 An ADR records that flai serve pushes a story's branch and opens its pull request on the move to review, behind the forge host action, with the force-with-lease exception for story branches
- T-1612 flai serve pushes each story in review under pull-request acceptance and opens or updates its pull request, raising a failure on one thread on the story
- T-1613 flai forge push pushes a story and opens its pull request on the operator's word, and flai forge status says what flai serve did and why
- T-1614 The workflow, the git convention and the template's, and the user and operator docs describe the pull request opened on the move to review

## Notes

- Drafted by the planner for E-0021, after S-0392, whose contract and token this story uses.
- The pull request is found by its head branch, `story/S-nnnn`, so nothing about it is written to the story's front matter, and a clone on another host finds the same one.
- The push uses the host's git credentials, as `flai release --pending` does. The request uses the token the host names. Both happen in `flai serve`, on the host.
- `flai serve checks` names commands that the host runs on a story in review. They are not the forge's checks. Keep the words apart in the docs.

### Planning

- **Touches.** None were declared. Every touch is the planner's.
- **Forecast and cost of delay.** Written after `flai forecast` and `flai cod`.
