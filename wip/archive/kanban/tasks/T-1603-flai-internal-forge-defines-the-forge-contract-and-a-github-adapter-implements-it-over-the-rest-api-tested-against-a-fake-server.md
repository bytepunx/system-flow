---
id: T-1603
type: task
nature: feature
title: flai/internal/forge defines the forge contract, and a GitHub adapter implements it over the REST API, tested against a fake server
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:26Z
updated: 2026-10-17T07:02:30Z
transitions:
  - to: ready
    at: 2026-10-17T06:42:32Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T06:42:33Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:02:30Z
    by: agent-S-0392
stream: S-0392
tags: [cli]
touches: [flai/internal/forge/forge.go, flai/internal/forge/forge_test.go, flai/internal/forge/github.go, flai/internal/forge/github_test.go]
usage:
  source: log
  seconds: 1197
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 12756
      cache_read: 1571496
      cache_write: 53616
      cost: 0.8758
---
# T-1603 flai/internal/forge defines the forge contract, and a GitHub adapter implements it over the REST API, tested against a fake server

## Work

- In `forge.go`, define the contract as an interface with neutral types, holding no GitHub name:
  - the repository, read from a remote URL in SSH or HTTPS form
  - `Find` a pull request by its head branch, including one already merged or closed
  - `Open` one, or `Update` its title and body
  - `Status`: the head commit; each required check with its name, state, and log URL; each review with its author, state, and comment URLs; and the number of approving reviews the base branch requires
  - `Merge` it by `merge`, `squash`, or `rebase`, and delete its head branch
  - `Access`: whether the token can read the repository and push to it
- Add a registry by forge name, so that a later adapter, such as GitLab, plugs in without any change to the callers.
- In `github.go`, implement the contract over the GitHub REST API with `net/http`:
  - Read the required checks and approvals from the branch protection, or from the ruleset, of the base branch. Read the checks from check runs and from commit statuses.
  - Take the base URL as a value, so that GitHub Enterprise and tests can set it.
  - Take the token as a value: the package never reads the environment.
- Add a fake forge for tests, an `httptest` server other packages can start, so that the later stories of E-0021 test against it.
- In `github_test.go`, test every method against it: success, a 401, a 404, a missing pull request, a failing required check, a requested change, and a merge refused.
- Never log the token, and never put it in an error.
- First in its layer: it waits for no task.

## Done when

The contract, the GitHub adapter, and the fake forge exist, and every method is tested against the fake, the failures included.

## Notes
