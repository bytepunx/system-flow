---
id: ADR-0166
title: "Under pull-request acceptance, flai serve pushes a story's branch and opens or updates its one pull request when the story is in review, behind the forge host action, with the token kept out of every agent's environment and a force-with-lease push of story branches as the one exception to never force pushing"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0032, ADR-0067, ADR-0165]
---

# ADR-0166 Under pull-request acceptance, flai serve pushes a story's branch and opens or updates its one pull request when the story is in review, behind the forge host action, with the token kept out of every agent's environment and a force-with-lease push of story branches as the one exception to never force pushing

## Context

[ADR-0165](0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md) lets a project choose pull-request acceptance and gives flai a forge contract and a token the host names. It does not say when a story's branch reaches the remote, or who opens its pull request. A team that reviews on GitHub has to see each story arrive as a pull request without anyone opening it by hand.

A story's agent moves its own story to review. That move runs in the agent's environment, with the agent's identity and whatever credentials it holds. The host's git credentials and the forge token belong to the operator. An agent should hold neither: with the token in its environment, an agent could push, open, or merge pull requests outside the workflow.

[ADR-0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md) and [ADR-0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) keep accepted work off the remote until it is published. A story branch under review is not accepted work. Under pull-request acceptance, though, it has to be on the remote before anyone can review it.

`flai stream sync` rebases `story/S-nnnn` onto main. Once the branch is on the remote, the next push after a sync is refused unless it is forced. `git.md` says never force push, and `safety.md` says a force push needs confirmation.

## Decision

Under pull-request acceptance, with the `forge` host action on, `flai serve` on the host pushes each story in review to `origin` and opens or updates its one pull request. The token stays out of every agent's environment, and a force-with-lease push of the story's own branch is the one exception to never force pushing.

- **Who pushes, and when.** Under `acceptance.mode: pull-request`, with the host action `forge` on, `flai serve` looks once a minute at each story in review. When `story/S-nnnn` on `origin` is missing or differs from the local head, it pushes the branch with the host's git credentials. It then opens a pull request against the main branch, titled `[S-nnnn] <title>`, or updates the one that branch already has. It does this on the host, never in the move itself, which runs in the agent's environment.
- **One pull request per story.** flai finds the pull request by its head branch, `story/S-nnnn`, and writes nothing about it into the story. A story that returns to review keeps its pull request, and a clone on another host finds the same one.
- **The body.** The body fills the project's `.github/pull_request_template.md` at places the template marks. It holds the story's ID, its acceptance criteria as the story has them ticked, the definition of done, and the story's last `flai verify` result with its commit. A project without that file gets a body with the same four parts under their own headings.
- **The host action.** `forge` is off by default and is turned on with `flai serve enable forge`. With it off, or with the mode `local`, flai pushes nothing and opens nothing.
- **The operator's word.** `flai forge push S-nnnn` does the same push and pull request for one story when the operator runs it. `flai forge status S-nnnn` says what flai did, and why it did nothing when it did nothing.
- **Failures.** A push or request that fails opens one thread on the story, naming the reason and what clears it. flai serve tries again at its next look and resolves the thread once the push and the request succeed.
- **Keeping the token from agents.** Every agent `flai serve` starts runs without the forge token's variable in its environment: story agents, their sub-agents, and strategic agents. `flai guard` refuses `flai forge push` to every agent.
- **The one exception to never force pushing.** When a sync has rebased a story branch, so that the remote head is no longer an ancestor of the local head, flai pushes `story/S-nnnn` with `--force-with-lease=story/S-nnnn:<head flai last pushed>`. flai keeps the last pushed head per story under `.flai-cache`. The lease refuses the push if someone else has pushed to the branch since. flai never forces `main` or any other branch. The planner recommended this on TH-0463 question 4, and the operator finalized S-0393 with it as a criterion. TH-0496 asks the operator to confirm it.

## Consequences

- Under pull-request acceptance, the remote carries every story branch in review. Accepted work still reaches the main branch on the remote only by publishing, or by the pull request's merge that S-0396 adds.
- `git.md`, in both the baseline and the template, names the force-with-lease push of story branches as the only exception to never force pushing. Agents still never force push.
- A pushed branch's history changes on the remote after every rebase. A reviewer's comments on lines a rebase moved may show as outdated.
- The operator must turn on `forge`, set the token's variable for `flai serve`, and give the host's git credentials push access to `origin`.

## Alternatives considered

- **Push in the move to review.** Rejected. The move runs as the agent, with the agent's environment, and would need the token there.
- **Under pull-request acceptance, sync merges main into story branches instead of rebasing them.** That needs no force push. But it makes history differ between the two modes, and it changes the trial merge, the conflict log, and the close-out's sync check. It is left to a story of its own if the operator chooses it on TH-0496.
- **Record the pull request's number in the story's front matter.** Rejected. Finding the request by its head branch needs no write to `wip/` from flai serve, and works the same from every clone.
