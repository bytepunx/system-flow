---
id: ADR-0039
title: "A settings host action, turned on only in a shell, lets the dashboard change the host's settings"
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: []
refines: [ADR-0029]
---

# ADR-0039 A settings host action, turned on only in a shell, lets the dashboard change the host's settings

## Context

Since S-0078 every host action has been off until the operator enabled it by name in a shell on the host, and nothing a dashboard could ask for changed that ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). A test held that guarantee: no method ran `flai serve` or `flai config`. The same held for everything else kept in the host's configuration: the agent's command and harnesses, the checks, the import folders, and the tokens. S-0105 asks for all of it to be manageable per project from the dashboard.

## Decision

**One more host action, `settings`, is the operator's consent to changing the host's settings from the dashboard.** The operator decided on 2026-09-23: the shell switch, all four areas, and commands too.

- `settings` is turned on and off only in a shell (`flai serve enable settings`, `flai serve disable settings`). No method changes it, and `settings.action` refuses to.
- With `settings` on for a project, that project's dashboard may change what belongs to the project:
  - its other host actions (push, agent, checks, dashboard);
  - its default agent, written to `system-flow.yaml` and committed with the dashboard's trailer;
  - its MCP server's token.
- With `settings` on for every project (`--all-projects`), the dashboard may also change what the host keeps for every project:
  - the agent's name, attended minutes, and command;
  - each harness's program and arguments;
  - the checks and their time limit;
  - the import folders;
  - the dashboard token.

  One project's consent does not stretch over the others.
- **Commands are editable too**, as argument lists: one argument a line on the page, run as they stand, never through a shell.
- Each change is a flai command that hostapi builds from values it checks, as every dashboard write is (`serve enable`, `agent set --replace --autocommit`, `serve agent set -- …`, `serve agent harness`, `serve checks set`, `serve import add`, `mcp token --rotate`, `dashboard token --rotate --no-restart`). Each is journalled as the command it ran. `flai config` and the journal stay out of reach.
- The dashboard token is rotated without restarting the dashboard. The dashboard takes the new token from the answer, and the asking session gets a cookie with it, so the operator stays logged in while every other session and agent is shut out. The new login link is shown once. No token is ever in `settings.get`.
- `flai serve` passes its configuration path on to the flai it runs, so a change lands in the file it reads, `--config` included.

## Consequences

- With `settings` on everywhere, the dashboard token is worth a shell on the host as the operator: whoever holds it can set the agent command, or a check, to anything. The action's text, the page, and the operator guide say so. The shell remains the only way to turn it off, so a stolen token cannot keep itself in or widen what the operator allowed.
- ADR-0029's guarantee becomes: nothing a dashboard asks for changes the host's settings unless `settings` is on, and nothing changes `settings` itself. The test now holds that.
- The trial found that rotating the MCP token restarted a server flai serve had started as one that would outlive flai serve. It now keeps `--exit-with` across a restart.

## Alternatives considered

- **No switch: the dashboard token alone.** Simpler, but the token would become shell-equivalent the moment this shipped, on every host, whether the operator wanted it or not.
- **Commands shown, not edited.** A stolen token could then only run what the operator already chose. The operator chose editing.
- **Per-project copies of the host-wide settings.** Each project would have its own agent command and checks, and one project's consent would suffice. It is a larger change to the configuration, and nothing asks for different commands per project yet.
