---
id: T-0966
type: task
nature: feature
title: An ADR, the dashboard and manifest design, flai-cli.md, and the operator and user guides describe editing the strategic agents' settings from the dashboard
status: in-progress
parent: S-0229
owner: alex
created: 2026-10-05T05:47:14Z
updated: 2026-10-06T22:10:14Z
transitions:
  - to: ready
    at: 2026-10-06T22:10:14Z
    by: agent-S-0229
  - to: in-progress
    at: 2026-10-06T22:10:14Z
    by: agent-S-0229
stream: S-0229
tags: [flai, dashboard]
touches: [design/adrs, design/system/flaiover-dashboard.md, design/system/project-manifest.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/settings.md, docs/operators/index.md, docs/users/flaiover.md, docs/users/flai.md]
after: [T-0954, T-0964]
---
# T-0966 An ADR, the dashboard and manifest design, flai-cli.md, and the operator and user guides describe editing the strategic agents' settings from the dashboard

## Work

The fourth criterion asks for every setting described, and the change widens what the `settings` host action lets a dashboard token do, which is a decision of its own.

- A new ADR in `design/adrs`, refining ADR-0039: with `settings` on for a project, its dashboard may also change the manifest's `orchestration`, `planning` (but `currency`), and `analysis` keys, through `flai manifest set`, committed with the dashboard's trailer; a dashboard token can then grant the orchestrator permission to accept and publish; the guard reads the permissions on each call, so a change holds from the orchestrator's next call.
- `design/system/flaiover-dashboard.md`: the settings panels on the Workflow pages, `settings.manifest` under the host's settings beside `settings.default_agent`, the `strategic` block of `settings.get`, the refusal shown on its field, and the read-only state.
- `design/system/project-manifest.md`: `flai manifest set` beside `flai agent set` in the rule on which keys flai rewrites, and each strategic key's sentence as the panels show it.
- `design/system/strategic-agents.md`: "On the settings page" says the planning keys are edited on the Planner page; the orchestrator's section says permissions may change while it runs and hold from its next call.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai manifest set`, its keys, its refusals, and `--autocommit`.
- `docs/operators/settings.md`: each strategic key's "Changed with" names the dashboard page and `flai manifest set`. `docs/operators/index.md`: the `settings` host action's section lists the manifest's strategic keys among what it lets the dashboard change.
- `docs/users/flaiover.md`: each page's settings panel, what each permission allows and risks, and why a panel is read-only.

This task waits for T-0964 and T-0954, so the documents describe what was built.

## Done when

- the ADR is recorded with `flai adr new`, refines ADR-0039, and is accepted by the operator or left proposed with a thread asking
- each document above describes every setting the panels edit, as built
- `flai check --strict`, the markdown lint, and `go test ./cmd/` (the settings index test) pass

## Notes
