---
id: ADR-0118
title: A dashboard release chosen on the Updates page keeps running through restarts until an upgrade without a tag or a start after a stop
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0117]
---

# ADR-0118 A dashboard release chosen on the Updates page keeps running through restarts until an upgrade without a tag or a start after a stop

## Context

[ADR-0117](0117-the-host-api-installs-a-release-the-dashboard-names-only-when-flai-lists-it-as.md) §5 says a dashboard release chosen on the Updates page applies once and that the next `flai dashboard restart` or `upgrade` without a tag uses the configured `dashboard.tag` again. The operator's answer on TH-0202 was "apply once, no pinning", to a recommendation of "once, like `flai dashboard upgrade --tag`". Reading the code in S-0298 showed that `flai dashboard upgrade --tag` already behaves otherwise for a restart: `flai dashboard restart` reuses the running container's image reference, and `flai host`'s watch restarts the reference `flai dashboard` recorded ([ADR-0062](0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md)). Only `flai dashboard upgrade` without a tag and `flai dashboard` starting a container that is not running use the configured tag.

## Decision

A dashboard release chosen on the Updates page, or with `flai dashboard upgrade --tag`, keeps running through restarts until an upgrade without a tag or a start after a stop.

1. A restart, `flai dashboard restart` or `flai host`'s watch, keeps the release that runs.
2. `flai dashboard upgrade` without a tag, the Updates page's Upgrade, goes to the configured tag.
3. `flai dashboard` starting a container after `flai dashboard stop` uses the configured tag.
4. Nothing writes `dashboard.tag`: pinning stays `flai config set dashboard.tag`, as ADR-0117 §5 says.

This replaces only ADR-0117 §5's sentence on what the next restart uses.

## Consequences

- A rollback survives a crash or a host restart, so the dashboard does not move forward again on its own.
- The Updates page and the runbook say how long a chosen release lasts in these terms, not "until the next start".
- The behaviour is what `flai dashboard upgrade --tag` already did; no code changes for it.

## Alternatives considered

- Making a restart go back to the configured tag: a crash or a host restart would then undo a rollback without the operator asking, and the configured tag is `latest` by default, so the dashboard would move to the newest release on its own.
- Editing ADR-0117 §5: an accepted ADR is not edited.
