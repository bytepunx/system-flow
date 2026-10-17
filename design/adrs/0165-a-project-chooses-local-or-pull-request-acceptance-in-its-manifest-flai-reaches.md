---
id: ADR-0165
title: "A project chooses local or pull-request acceptance in its manifest, flai reaches the forge through a contract with one adapter per forge, and the host names the variable that holds the forge's token"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0129]
---

# ADR-0165 A project chooses local or pull-request acceptance in its manifest, flai reaches the forge through a contract with one adapter per forge, and the host names the variable that holds the forge's token

## Context

Acceptance merges a story's branch into main on the host ([ADR-0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md)). A team whose main branch is protected, whose CI is the gate, or whose reviewers work on GitHub cannot accept that way. E-0021 gives a project the choice of accepting through a pull request instead. Every later story of the epic needs three things settled first: where the project makes the choice, how flai talks to the forge, and where the forge's token lives. GitLab is to follow in a later epic, so the way flai talks to the forge must not be GitHub's shape.

## Decision

A project chooses local or pull-request acceptance in its manifest, flai reaches the forge through a contract with one adapter per forge, and the host names the variable that holds the forge's token.

- **The setting.** `system-flow.yaml` takes `acceptance.mode`, `local` (the default) or `pull-request`; `acceptance.forge`, `github` (the default, and the only forge today); and `acceptance.merge_method`, `merge` (the default), `squash`, or `rebase`. `flai manifest set` sets each. The manifest's load refuses an unknown value, naming the key and the values it takes, rather than falling back to `local`: a project that meant to gate on a pull request must not have stories merged locally because of a typing error. A manifest without the block loads and accepts as before. The forge and the merge method are allowed under `local`, and unused.
- **The contract.** `flai/internal/forge` defines `Forge`, bound to one repository when it is made: `Access` (can the token read the repository and push to it), `FindPullRequest` by head branch, `OpenPullRequest` (open one, or update the open one for that head), `Checks` on the head commit (each with whether the base branch requires it, a required check that has not reported being pending), `Reviews`, and `Merge` by a method, refused with `ErrMergeRefused` when the forge will not merge or the head has moved. Callers test errors with `ErrNotFound`, `ErrUnauthorized`, `ErrForbidden`, and `ErrMergeRefused`. No method takes an owner or a repository name, and what is a forge's own, such as GitHub's REST paths, its check runs beside commit statuses, its rulesets beside branch protection, and its remote URLs, stays in that forge's adapter. GitHub's adapter, `github.go`, speaks the REST API with the standard library, so that it needs nothing on the host's `PATH` and is tested against an `httptest` server.
- **The token.** The host's configuration takes `forges.<forge>.token_env`, the name of the environment variable that holds the forge's token, set with `flai serve forge github --token-env <NAME>` beside the providers' `key_env` ([ADR-0129](0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). Unset, it is `GITHUB_TOKEN`. An adapter gets the token from a `TokenFunc` at each call and keeps nothing it answers, so the variable is read only when flai calls the forge. The token is never written to the configuration, a log, a work item, an error, or a start record; the commands that report on it, `flai serve forge` and `flai forge check`, name the variable and say whether it is set.
- **The check.** `flai forge check` reads the repository from the `origin` remote and says whether the token can read it and push to it, exiting non-zero with the reason when it cannot.

## Consequences

- The later stories of E-0021 build on the contract without a change to the setting: S-0393 opens or updates the pull request, S-0394 reads checks and reviews back, S-0395 gates acceptance on them, and S-0397 merges by the project's method.
- A GitLab adapter implements `Forge` and adds `gitlab` to the forges the manifest and the host configuration take; acceptance does not change.
- A project moved to another host names its forge's variable again there, as it names its providers. A clone carries the mode and the method, never a token or a variable's name.
- A bad acceptance value stops every command that loads the manifest until it is fixed or removed by hand.
- The contract does not yet say how many approvals the base branch requires or who is asked to review; S-0395 extends it if its gate needs them.

## Alternatives considered

- **The `gh` command line rather than the REST API.** The epic allows either. `gh` would need installing and authenticating on the host and could not be tested against a fake server.
- **The token's variable in the manifest.** A clone would carry a name that is the host's, against ADR-0129's split of what is the project's and what is the host's.
- **The forge settings inside `agent.providers` in the host configuration.** The forge is not where an agent's model calls go, and `flai serve forge` is not an agent command, so the entries sit beside the providers at the top of the configuration.
- **An unknown mode reported only by `flai check`, as the strategic agents' settings are.** Acceptance would go on locally in the meantime, which is what a project that set `pull-request` chose against.
