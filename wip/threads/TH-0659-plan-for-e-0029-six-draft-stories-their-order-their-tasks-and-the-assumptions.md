---
id: TH-0659
title: "Plan for E-0029: six draft stories, their order, their tasks, and the assumptions"
anchor:
  path: wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md
  item: E-0029
status: resolved
participants: [planner-E-0029, orchestrator, alex]
created: 2026-10-26T04:58:13Z
updated: 2026-10-26T12:28:11Z
---

# TH-0659 Plan for E-0029: six draft stories, their order, their tasks, and the assumptions

On wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md.

## Entries

### 2026-10-26T04:58:13Z planner-E-0029
Six draft stories deliver the epic, in this order (`after`):

1. S-0505 (experiment, waits for nothing): the spike ADR-0203 names, with E-0029's layered images built by hand and locally (`flai/images/bedrock`, `flai/images/claude-code`, `scripts/agent-image.sh`, `make agent-images`), one documentation-only story worked inside the container over the HTTP MCP server, and the results settling the login, the bind mount's speed, the HTTP server's tool coverage, the guard inside, the network mode, and the tiers that cannot run. Tasks, one per layer: T-2149 (images), T-2159 (container script, login, network), T-2166 (the story worked inside, tools and guard recorded, issues filed), T-2173 (mount speed), T-2177 (results, index, agent-adapters.md, docker.md).
2. S-0506 (after S-0505): the HTTP MCP server serves a containerised agent as stdio does. Layers: T-2150 (tool parity test), T-2151 (ADR on `permission_prompt` under auto-approve over HTTP), T-2152 (`mcp.bind`) together; then T-2160 (`FLAI_MCP_URL`, `FLAI_MCP_TOKEN`); then T-2167 (docs).
3. S-0507 (after S-0505): pilot, the Go binary that is the image's entrypoint, and the ADR recording layered images and pilot. Layers: T-2153 (ADR), T-2154 (`harness.Request` container fields, HTTP MCP file) together; T-2161 (request); T-2168 (run, relay, resume, signals); T-2174 (images build pilot); T-2178 (docs).
4. S-0508 (after S-0506 and S-0507): the `container` harness adapter and `flai serve`'s start path. Layers: T-2155 (adapter), T-2156 (host settings), T-2157 (protected list) together; T-2162 (serve start, refusals, stop, resume); T-2169 (one story worked end to end under flai serve, docs).
5. S-0509 (after S-0508): the four strategic agents in containers. Layers: T-2158 (ADR on the main checkout's mount); T-2163 (adapter and pilot for a role); T-2170 (the four start paths); T-2175 (a planner run in a container, docs).
6. S-0510 (after S-0508): a project's image layer named in the manifest, `flai check agent.image`, and this repository's layer with Go, golangci-lint, Node, and pnpm. Layers: T-2164 (manifest key), T-2165 (build script and this repository's Dockerfile) together; T-2171 (check), T-2172 (adapter runs the layer) together; T-2176 (a Go-changing story worked inside, docs).

S-0509 and S-0510 can run together once S-0508 lands; S-0506 and S-0507 can run together once S-0505 lands.

Assumptions:

- pilot is the container's entrypoint in place of ADR-0203's `flai agent run`, as the epic's notes ask, built as a second `main` of the flai module (`flai/pilot/main.go`, `flai/internal/pilot`); S-0507's ADR records this and the layered images.
- The images live under `flai/images/<layer>/Dockerfile`, tagged `flai-agent-<layer>:local`, built by `scripts/agent-image.sh` from the repository root; nothing under `.github/workflows`. The signed, published image and its release pipeline are outside this epic.
- `harness: container` names the adapter and `config.inner` the harness pilot runs inside (`claude-code` by default). pilot relays the inner harness's log unchanged, so the adapter's `Capabilities` and `Reader` are the inner harness's and no neutral usage event is built here.
- Only Claude Code is baked in. ADR-0203's loop of flai's own, local models, the `policy` value of `priced_by`, Kubernetes, and signet are outside this epic.
- The container's MCP configuration is passed with `--mcp-config` and `--strict-mcp-config` from the story's cache folder, never by writing the worktree's protected `.mcp.json`.
- For the strategic agents the main checkout is mounted read-write at its host path (recommended); S-0509's ADR decides.
- `permission_prompt` over HTTP honours the auto-approve host action for an agent `flai serve` started and named (recommended); S-0506's ADR decides.
- The proof stories (S-0505's and S-0508's documentation-only story, S-0510's Go-changing story) are chosen with you when each is worked.
- Every story keeps a folder touch, `design/adrs` or `design/issues`, for the files `flai adr new` and `flai issue new` name at run time; each story's Planning note says so.

Proposals (nothing created for them):

- Add, as a follow-up epic rather than here: flai's own agent loop for local models in a role, with the `policy` rate and `priced_by` per model from ADR-0203, once S-0508 shows the container path works.
- Add later: the signed `ghcr.io/bytepunx/flai-agent` image and its release workflow, built as `flaiover`'s is, once the local images are stable; and an OpenCode layer when E-0023 lands.
- Drop from the spike: ADR-0203's one measurement of Docker Agent as the inner loop, since ADR-0202 is rejected; say so if you want it kept and I will add a task to S-0505.

Cost of delay values wait on TH-0657 (the epic's inputs); forecasts are written on each story with their basis.

### 2026-10-26T04:58:28Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan as drafted, and take the planner's three proposals.

- **The plan fits.** There are six stories, ordered by `after` only where one needs another's change on main. S-0506 and S-0507 can run beside each other after S-0505, and S-0509 and S-0510 after S-0508. Each story's tasks form layers that share no path. Each ADR question, `permission_prompt` over HTTP and the strategic agents' mount, is left to its story's ADR with a stated recommendation.
- **Proposal 1, a follow-up epic for flai's own loop with local models: yes, later.** Defer it until S-0508 shows the container path works.
- **Proposal 2, the signed `ghcr.io` image and its release workflow, and an OpenCode layer: yes, later.** Defer them until the local images are stable and E-0023 lands.
- **Proposal 3, drop the Docker Agent measurement from S-0505: yes.** ADR-0202 was rejected, so measuring Docker Agent as the inner loop answers a question already closed.

**Why it is yours.** Proposals 1 and 3 change scope, and finalizing E-0029's drafts waits on TH-0657's input.

**What I do after.** Once TH-0657 is answered and the planner writes the values, I finalize under finalize_drafts each draft that `flai promote --drafts` lists as complete and consistent.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-26T12:23:01Z alex
Confirmed the recommendation of 2026-10-26T04:58:28Z orchestrator.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-26T12:26:58Z alex
TH-0657 is answered

### 2026-10-26T12:28:11Z alex
Resolved.
