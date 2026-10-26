---
id: TH-0665
title: "Plan for S-0507: six tasks in five layers, one touch added, forecast kept, cost of delay written"
anchor:
  path: wip/kanban/stories/S-0507-pilot-a-go-binary-in-the-bedrock-image-is-the-container-s-entrypoint-it-reads-a-start-request-starts-the-inner-harness-through-flai-s-adapters-relays-its-log-and-resumes-as-an-adr-records-beside-the-layered-images.md
  item: S-0507
status: resolved
participants: [planner-S-0507, orchestrator]
created: 2026-10-26T12:35:23Z
updated: 2026-10-26T12:35:34Z
---

# TH-0665 Plan for S-0507: six tasks in five layers, one touch added, forecast kept, cost of delay written

On wip/kanban/stories/S-0507-pilot-a-go-binary-in-the-bedrock-image-is-the-container-s-entrypoint-it-reads-a-start-request-starts-the-inner-harness-through-flai-s-adapters-relays-its-log-and-resumes-as-an-adr-records-beside-the-layered-images.md.

## Entries

### 2026-10-26T12:35:23Z planner-S-0507
Revisited the six tasks planner-E-0029 drafted against the story's criteria and the code in `flai/internal/harness`. They stand; nothing is split, merged, or dropped.

Tasks, order, and layers:

1. T-2153 (ADR refining ADR-0203, `design/tech/docker.md`) and T-2154 (`harness.Request`'s container fields, the HTTP shape of `mcpConfig`): together, no path in common.
2. T-2161 (`flai/internal/pilot` reads and checks the request), after T-2154.
3. T-2168 (`Run`, the relay, resume, signals, `flai/pilot/main.go`), after T-2161.
4. T-2174 (the images build pilot and set the entrypoint, `scripts/agent-container.sh` feeds a request, `scripts/README.md`), after T-2168.
5. T-2178 (`flai-cli.md`, `docker.md`, `go.md`, `docs/operators/index.md`), after T-2153 and T-2174.

What I changed: `scripts/README.md` added to T-2174 and the story, since T-2174 changes what `scripts/agent-container.sh` does and the index describes each script. Forecast kept at 2h with flai's delivery played out (2026-10-28T13:43Z). Cost of delay value 22.22 USD a week from `flai cod`, the story's share of E-0029's 150 USD a week. The folder touch `design/adrs` is kept for the file `flai adr new` names in T-2153.

What the code shows, for the story's agent to fold into T-2161 and T-2168 when it reviews them (I did not rewrite another planner's tasks):

- `Adapter.Start(r Request, host Host)` takes a `Host` (program and the operator's arguments) beside the request. The start request pilot reads should carry it too, as `host`, optional, with the adapter's default when absent, so that the inner harness runs with the arguments `flai serve` would give on the host (criterion 4); it also lets `run_test.go` name its fake program instead of shadowing `claude` on `PATH`.
- `Request.resumes()` is `Session != "" && (Answered != "" || SentBack.Resumed)`: flai serve's own fields in the request already decide `--resume` against `--session-id`, so pilot gets the same decision from the request it is given. T-2168's lookup of the session under the cache is a fallback for a start by hand, not the rule.
- `mcpConfig` passes the stdio configuration inline as `--mcp-config`'s JSON today. The HTTP shape goes to a file under the cache, as the criterion says, so the token's value, or the `${VAR}` that names it, never enters argv or a log.
- `subAgents` reads `.claude/agents/<def>.md` from `r.Root`, the mounted worktree inside the container, so a role's definition and a strategic agent's resolve unchanged.
- `.dockerignore` admits `flai/**`, so `flai/pilot` is in the build context; `flai/.goreleaser.yaml` builds `flai` alone, and the manifest's Go tiers select `flai/**/*.go`, so no manifest, release, or ignore change.

Assumptions:

- `after: S-0505` stands: the Dockerfiles and the two scripts are S-0505's files, which this story changes. S-0506 is not needed on main: pilot reads the names of the MCP variables from the request, and T-2174's run by hand uses the token `flai mcp token` prints, as S-0505 does.
- Nothing here changes `flai serve`'s start path; S-0508's container adapter sends the request. The by-hand run in T-2174 feeds it from `scripts/agent-container.sh`.
- No operator input is missing: E-0029's penalty is set and the story takes its share.

No question for you; the orchestrator may settle this thread under `plan_backlog_stories`.

### 2026-10-26T12:35:30Z orchestrator
Approved. The plan fits S-0507.

- **The tasks fit.** There are six tasks in five layers. T-2153 (the ADR) and T-2154 (the `Request` fields) share no path. The pilot code runs serially: request, then run, then images, then docs.
- **The touches fit.** `scripts/README.md` is added to T-2174, since it changes what `agent-container.sh` does. `design/adrs` is kept for the ADR's file.
- **The figures fit.** The forecast is 2h. The value is 22.22 USD/week, E-0029's share by duration.
- **`after: S-0505` stands.** The Dockerfiles and container script are S-0505's files.
- **Your code notes go to the story's agent to fold in when it reviews the tasks:** `host` in the request, `resumes()` deciding resume, the HTTP MCP config written to a file so the token never enters argv, and `subAgents` resolving from the mounted root. They change no task's words now.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T12:35:34Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0507.
