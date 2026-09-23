---
id: ADR-0036
title: A folder that is not a project is served whole by flai mcp and flai dashboard
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: []
refines: [ADR-0030]
---

# ADR-0036 A folder that is not a project is served whole by flai mcp and flai dashboard

## Context

`flai mcp` and `flai dashboard` opened the project above the working directory first and failed without one (S-0101). An operator who keeps repositories in `~/git`, and starts an agent or the dashboard there, got "no system-flow.yaml found". Yet the dashboard's container, the login token, the agent credential, and `flai serve` are already per user, not per project ([ADR-0033](0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md)). And since S-0098 a folder can be named so that its repositories are offered for import ([ADR-0035](0035-repositories-under-folders-the-operator-names-can-be-imported-from-the-board.md)).

## Decision

**A folder that is not a project is served whole.** The operator decided (2026-09-23).

- **`flai mcp` there serves every system-flow project in the folder and below it**, up to three levels, not into hidden folders or build output. It looks again every few seconds, so a project created or imported while the agent works joins without a restart.
  - `inbox`, `wait_for_work`, and `wait_for_events` cover every project and name the project each thing is in. `wait_for_work` resumes the agent's own story in whichever project it is. Otherwise it pulls the first ready story in key order of projects, and pull order within one, of a project whose in-progress limit leaves room.
  - Every other tool takes `project`: a key `inbox` lists, or the project's folder. It is needed whenever there is more than one project. A single-project server takes it too and refuses any other.
  - With no project below the folder, the server still starts and says so.
  - This refines [ADR-0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md) for stdio only: `flai mcp` over HTTP still serves one project, and outside one it says what to run instead.
- **`flai dashboard` there starts or reports the shared dashboard and `flai serve` without a project of its own.**
  - It registers the projects below the folder, as `flai mcp` serves them there.
  - It names the folder for import, as `flai serve import add .` would.
  - It records the dashboard's address, so that `flai serve` offers it repositories to import even when no project is served.
  - Its other subcommands (`status`, `stop`, `token`, `restart`, `upgrade`, `check`, `logs`) work there with the host's configuration alone.
  - `stop` stops the container only once no project is left registered, and then forgets the recorded address.
- A manifest that does not parse is still an error, not "no project here".

## Consequences

- One agent can be started in `~/git` and work across every project in it, and one `flai dashboard` there shows them all and offers the rest for import.
- An agent in a folder must name the project on every per-project tool; the instructions and tool descriptions say so, and an unnamed call lists the projects.
- `flai stream open` and the other CLI steps still run in the project's own folder; `wait_for_work` answers with that folder.
- A cursor is still kept per project, so an agent that moves between a project and its parent folder does not see changes twice.
- Running `flai dashboard` in a folder names it for import, which is consent to importing, and testing, the repositories in it (ADR-0035). The output says so, and `flai serve import remove .` undoes it.

## Alternatives considered

- **Start, and point to the projects**: a `projects` tool, and every other tool failing outside a project. Smaller, but an agent would still have to be restarted in each project; the operator chose the folder.
- **An HTTP server for the folder too.** Its state, token, and address live in a project's `.flai-cache`, and `flai serve` already keeps one running per project (ADR-0034). A folder would need a place of its own for them, for no agent that stdio does not already serve.
- **`flai dashboard` in a folder only starts the container.** The projects below it would then not be on the board until `flai dashboard` ran in each; registering them is the folder's counterpart of registering the project a command runs in.
