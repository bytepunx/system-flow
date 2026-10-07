---
id: ADR-0117
title: "The host API installs a release the dashboard names only when flai lists it as published, and MCP lists releases without installing one"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
---

# ADR-0117 The host API installs a release the dashboard names only when flai lists it as published, and MCP lists releases without installing one

## Context

The Updates page (`/host`) installs only the newest flai and the newest dashboard. `host.upgrade` and `dashboard.upgrade` take no version, image, or tag from the dashboard: flai installs only what the host's own configuration names (S-0081, S-0106). Going back to an earlier release takes a shell, `flai self-upgrade --version` or `flai dashboard upgrade --tag`, and nothing lists the releases there are to go back to. S-0298 lets the operator list the published releases and deploy a chosen one from the CLI, the host API, and the Updates page. The operator answered TH-0202: MCP lists versions and deploys nothing, and a dashboard version chosen on the Updates page applies once, with no pinning.

flai releases are GitHub releases tagged `flai/vX.Y.Z`, with assets. A dashboard release is a git tag `flaiover/vX.Y.Z`, for which `release-flaiover.yml` pushes the image tagged with the bare version, `X.Y.Z`; it makes no GitHub release.

## Decision

The host API installs a release the dashboard names only when flai lists it as published, and MCP lists the published releases without installing one.

1. **What is published.** A published flai release is a GitHub release of the configured releases repository whose tag is `flai/vX.Y.Z`, neither a draft nor a prerelease. A published dashboard release is a tag `flaiover/vX.Y.Z` of the same repository, with no suffix; its image tag is `X.Y.Z`. flai lists each newest first by semver, across every page the API answers (`selfupgrade.List` and `selfupgrade.ListTags`).
2. **The names.** The CLI lists with `flai self-upgrade --list`, `flai dashboard versions`, and `flai host versions`, each with `--json`, and installs a chosen release with `flai self-upgrade --version`, `flai host upgrade --version`, and `flai dashboard upgrade --tag`. The host API reads `host.versions` and `dashboard.versions`. MCP has the read tool `versions`.
3. **What the dashboard may name.** `host.upgrade` takes an optional `version` and `dashboard.upgrade` an optional `tag`, each a bare `X.Y.Z`. The host API refuses anything else before a command runs. The commands it runs refuse a release that is not published, naming the published ones, before anything is pulled or replaced: `flai self-upgrade --version` always, and `flai dashboard upgrade --tag` with `--published`, which the host API always passes. A tag the operator gives in a shell without `--published`, such as a mirror's, is used as it is. The image name is never taken from the dashboard.
4. **Like an upgrade.** Installing a chosen release goes the same way as installing the newest: the same host action (`host`, `dashboard`), the same progress and detached run, the same restart of the host and its processes, or the same swap and health check of the container.
5. **No pinning.** A dashboard version chosen on the Updates page applies to that container once. The configured `dashboard.tag` is not changed, so the next `flai dashboard restart` or `upgrade` without a tag uses it again. Pinning stays `flai config set dashboard.tag`.
6. **MCP lists only.** The `versions` tool answers both lists, marking the flai it runs, the newest of each, and any flai below the project's `flai.minimum`. No MCP tool installs or deploys a release.
7. **The minimum warns.** A flai release below a served project's `flai.minimum` is marked in every list and warned about when chosen, not refused: going back past it is the operator's call.

## Consequences

- An operator can go back from the Updates page, the CLI, or the host API, and an agent can see which releases there are, but only the operator deploys one.
- A dashboard request can install any published release, older ones included, where before it could install only the newest. Since flai lists what may be installed from the releases repository, and verifies a flai archive's signature as before ([ADR-0070](0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md)), a compromised dashboard can at most move the host to an earlier signed release, never to an image or binary of its choosing.
- Each listing calls the GitHub API: a page per 100 releases, or one call for the dashboard's tags. Unauthenticated, that counts against GitHub's hourly limit per address; a token in the configuration lifts it.
- A dashboard tag pushed before its image is built lists as published, and a pull of it fails, which `flai dashboard upgrade` already reports without replacing the container.

## Alternatives considered

- Listing dashboard releases from the image registry: it names what can be pulled, but needs the registry's token flow for each registry, and a mirror the operator configures would list differently from the releases repository the host API checks against.
- Taking any tag from the dashboard and letting the pull fail: the dashboard could then name any tag of the image, a mutable one such as `latest` included, which S-0081 kept it from.
- Pinning a version chosen on the page in `dashboard.tag`: the operator answered on TH-0202 that it applies once.
- An MCP tool that deploys: the operator answered on TH-0202 that MCP lists only.
