---
id: S-0393
type: story
nature: feature
title: Under pull-request acceptance, a story's move to review has flai serve push its branch and open or update one pull request from the project's pull request template
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:34:50Z
updated: 2026-10-17T01:56:07Z
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
  delivery: 2026-10-17T04:43:00Z
  basis: "Its own forecast of 1h5m; 8th in the pull order with an in-progress limit of 5, behind S-0440, S-0441, S-0442, S-0444, S-0443, S-0241 and S-0392."
  by: flai
  at: 2026-10-17T01:56:07Z
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

- **Touches.** None were declared. Every touch is the planner's, and each is a file except `design/adrs`:
  - **Layout:**
    - `flai/internal/forge/pullrequest.go`, with its test, and `template/root/.github/pull_request_template.md`: criterion 3.
    - `flai/internal/serve/forge.go`, with its test, and `serve.go` to wire it into the minute loop: criteria 1, 4, and 5.
    - `flai/internal/hostapi/writes.go`, with its test: `ActionForge` in the `Actions` map, beside `ActionPush`.
    - `flai/internal/serve/start.go`, with its test: the agents' environment. It may prove to be built in `flai/internal/harness/adapters.go`, where `ANTHROPIC_AUTH_TOKEN` is set; T-1610 says to follow it there.
    - `flai/internal/guard/guard.go`, with its test: refusing `flai forge push` to agents.
    - `flai/cmd/forge.go`, with its test: `flai forge push` and `flai forge status`.
  - **Design:**
    - `design/system/workflow.md`, `design/conventions/git.md`, and `template/root/design/conventions/git.md`: what the move to review now does, and the force-push exception.
    - The CLI and operator docs, the reference, and `template/CHANGELOG.md`.
- **Folder touch kept.** `design/adrs`, for an ADR whose file name `flai adr new` makes. It is a shared path, so it holds nothing.
- **Co-change left out.**
  - `flai/cmd/serve_actions.go`: `flai serve enable` reads its names from `hostapi.Actions`, so the new action needs no change there.
  - `flaiover/src/lib/server/agent.ts` and the flaiover docs: this story adds no host method. S-0394 adds `forge.status`.
  - `flai/internal/mcpserver/*`: no tool changes.
- **Paths it shares with other stories.** `flai/internal/serve/serve.go` with S-0391, now done, and with E-0022's S-0399. `flai/internal/hostapi/writes.go` with S-0395, S-0397, and E-0022's S-0399.
- **Layers.**
  - First: T-1609, T-1610, and T-1611 together.
  - Second: T-1612, after T-1609 and T-1610.
  - Third: T-1613, after T-1612.
  - Fourth: T-1614, after T-1611 and T-1613.
- **Forecast: 1h5m.** `flai forecast` gave 50m: 96 s per unit over 61 done feature stories on claude-opus-5-5 in the large band, times size 31. Raised to 1h5m for git against a scratch remote with a lease, the guard, and the agents' environment, across six tasks in four layers. flai's replan sets the delivery from this duration.
- **Cost of delay: not yet set.** `flai cod` refuses until E-0021 has inputs, asked for on TH-0463.
- **Criterion 4 depends on TH-0463 question 4.** If the operator chooses merging main into story branches instead, criterion 4 and T-1611 and T-1612 change. The planner will ask before rewriting them once the story is finalized.
