---
id: S-0392
type: story
nature: feature
title: A project chooses local or pull-request acceptance in system-flow.yaml, and flai reaches its GitHub remote through a forge contract with a token the host names
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:33:00Z
updated: 2026-10-17T03:34:53Z
transitions: []
tags: [cli, template]
topics: [release]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/forge/forge.go, flai/internal/forge/forge_test.go, flai/internal/forge/github.go, flai/internal/forge/github_test.go, flai/cmd/forge.go, flai/cmd/forge_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/tests/integration/forge_github_test.go, design/adrs, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h
  delivery: 2026-10-17T04:40:00Z
  basis: "Its own forecast of 1h; 4th in the pull order with an in-progress limit of 5, behind S-0442, S-0444, S-0443 and S-0241."
  by: flai
  at: 2026-10-17T03:34:53Z
finalized:
  by: alex
  at: 2026-10-16T05:26:51Z
---
# S-0392 A project chooses local or pull-request acceptance in system-flow.yaml, and flai reaches its GitHub remote through a forge contract with a token the host names

## Goal

An operator whose team gates main with pull requests can tell flai so. They set the acceptance mode in `system-flow.yaml`, name on the host the variable that holds a forge token, and run `flai forge check` to see whether flai can reach the repository. Nothing yet opens or merges a pull request. This story lays the setting and the forge contract, with GitHub behind it, that every later story of E-0021 builds on. A project that sets nothing sees no change.

## Acceptance criteria

- [ ] `system-flow.yaml` takes `acceptance.mode`, which is `local` (the default) or `pull-request`, plus `acceptance.forge` (`github`) and `acceptance.merge_method` (`merge`, `squash`, or `rebase`, default `merge`). Each is validated at load and settable with `flai manifest set`. A manifest without them loads and behaves exactly as today.
- [ ] The host names the variable that holds a forge's token in flai's host configuration, beside the providers (ADR-0129), with `flai serve forge github --token-env <NAME>`. The default name is `GITHUB_TOKEN`. flai reads the token from that variable only when it calls the forge, and never writes it to the configuration, a log, a work item, or a start record.
- [ ] `flai/internal/forge` defines the contract that acceptance uses: find a pull request by its head branch, open one or update it, read its required checks and its reviews, and merge it by a method. A GitHub adapter implements the contract over the REST API with the standard library, and is tested against a fake server.
- [ ] `flai forge check` prints the forge, the repository it reads from the `origin` remote, the variable it reads the token from, and whether the token can read the repository and push to it. It names the variable but never the token, and exits non-zero with the reason when any of these fails.
- [ ] An integration test against a real repository runs only where `FLAI_FORGE_TEST_REPO` and the token are set. Otherwise it is skipped, naming what is missing, so `make integration` and a close-out pass without them.
- [ ] An ADR records the setting, the contract, and where the token lives. `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, the generated reference, `docs/operators/settings.md`, and the template's `system-flow.yaml`, as a commented default, describe them.

## Tasks
- T-1602 The manifest takes acceptance.mode, acceptance.forge, and acceptance.merge_method, validated at load and settable with flai manifest set
- T-1603 flai/internal/forge defines the forge contract, and a GitHub adapter implements it over the REST API, tested against a fake server
- T-1604 The host names the variable that holds a forge's token with flai serve forge, kept beside the providers and never the token itself
- T-1605 An ADR records pull-request acceptance as a project's choice, the forge contract, and where the forge token lives
- T-1606 flai forge check says whether the token the host names can read and push to the repository behind origin
- T-1607 An integration test drives the GitHub adapter against a real repository only where FLAI_FORGE_TEST_REPO and the token are set, and is skipped saying so otherwise
- T-1608 The design, user, and operator docs, the reference, and the template's manifest describe the acceptance setting, the forge, and its token

## Notes

- Drafted by the planner for E-0021. GitLab is a later epic: the contract must hold it without any change to acceptance. Keep GitHub specifics, such as the REST paths, the `x-access-token` user, and check runs against commit statuses, inside `github.go`.
- The REST API, not `gh`, so that the adapter needs nothing on the host's `PATH` and can be tested against `httptest`. The epic allows either.

### Planning

- **Touches.** None were declared. Every touch is the planner's, and each is a file except `design/adrs`:
  - **Layout:**
    - `flai/internal/manifest/manifest.go` and `settings.go`, with their tests: the `Manifest` struct and `buildCatalog`. Criterion 1.
    - `flai/internal/config/config.go`, with its test: the host's `Config`, beside `Agent.Providers`. Criterion 2.
    - `flai/cmd/serve_actions.go`, with its test: where `flai serve agent provider --key-env` lives. Criterion 2.
    - `flai/internal/forge/forge.go` and `github.go`, with their tests: a new package, with the fake forge the later stories test against. Criterion 3.
    - `flai/cmd/forge.go`, its test, and `flai/cmd/root.go`: criterion 4.
    - `flai/tests/integration/forge_github_test.go`: criterion 5, after the skip in `mcp_http_test.go`.
  - **Design:** the docs, the generated reference, and `template/root/system-flow.yaml.tmpl`, for criterion 6.
  - **Co-change:** `docs/operators/index.md`, the most frequent co-change at 16%. The token's setup belongs there.
- **Folder touch kept.** `design/adrs`, because the ADR's file name is not known until `flai adr new` makes it. It is a shared path (`claims.shared`), so it holds no story.
- **Co-change left out.**
  - `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flai/internal/hostapi/writes.go`, and the harness files: this story adds no dashboard method and no host action.
  - `template/template.yaml`: the version is bumped at publish, not by a story.
- **Shared paths.** `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/adrs`, and `template/CHANGELOG.md` hold nothing.
- **Overlap with E-0022.** S-0399 also changes `manifest.go`, `settings.go`, `config.go`, `serve_actions.go`, and the template's manifest. Whichever starts second waits, as TH-0464 says.
- **Layers.**
  - First: T-1602, T-1603, T-1604, and T-1605 together.
  - Second: T-1606, after T-1602, T-1603, and T-1604; and T-1607, after T-1603.
  - Third: T-1608, after T-1605 and T-1606.
- **Forecast: 1h.** `flai forecast` gave 48m: 96 s per unit of size over 61 done feature stories on claude-opus-5-5 in the large band, times size 30. Raised a quarter to 1h for a new package with an HTTP adapter and a fake forge the later stories reuse, plus an ADR, across seven tasks in three layers. flai's replan now sets the delivery from this duration.
- **Cost of delay: not yet set.** `flai cod` refuses: neither E-0021 nor this story has inputs. They are asked for on TH-0463.
