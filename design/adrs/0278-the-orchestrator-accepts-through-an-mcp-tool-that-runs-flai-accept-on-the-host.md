---
id: ADR-0278
title: "The orchestrator accepts through an MCP tool that runs flai accept on the host and verifies through the close-out job, and its shell's flai accept and tier-running flai verify refuse in a role's container"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0093, ADR-0247]
---

# ADR-0278 The orchestrator accepts through an MCP tool that runs flai accept on the host and verifies through the close-out job, and its shell's flai accept and tier-running flai verify refuse in a role's container

## Context

[ADR-0247](0247-a-strategic-agent-s-container-mounts-the-main-checkout-read-write-at-its-host.md) runs the orchestrator in a container that mounts the main checkout and the role's cache folder, `.flai-cache/agents/orchestrate`. The container's flai reads its configuration from `FLAI_CONFIG`, which `flai/internal/harness/container.go` sets to `<root>/.flai-cache/agents/<role>/flai/config.json`. The host's `~/.flai/config.json` is never mounted, so that file holds no host action. [I-0256](../issues/I-0256-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md) records what follows.

- **`flai accept --by orchestrator`** decides between landing a story and enqueuing it from the host actions: while the merge queue batches, it enqueues ([ADR-0225](0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md)). `hostActionOn` in `flai/cmd/verify.go` reads a missing configuration as every action off. In the container, the orchestrator's acceptance therefore lands the story by hand while the queue batches, and can drop a batch.
- **`flai verify`** leaves the tiers marked `queue` to the batch only while the queue can batch ([ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)). In the container it reads the queue as off and runs every tier, and the image has neither Go nor Node for them.

Two routes already run on the host. The MCP tool `verify` asks flai serve for the close-out job, which runs in the story's worktree on the host and records its verdict ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)); `flai verify --last` only reads that record. The MCP tool `release_publish` runs `flai release --pending` on the host through `hostapi.Publisher`, under the host's configuration ([ADR-0094](0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md)). The container's agent reaches the host's HTTP MCP server with `FLAI_MCP_URL` and `FLAI_MCP_TOKEN`, naming itself in `X-Flai-Agent`. The host's dashboard runs the orchestrator's acceptance as `accept.run` (`flai/internal/hostapi/writes_items.go`), which builds `flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -`.

One more gap shows in the code. The MCP server learns its caller's role from its own `FLAI_ROLE` (`flai/internal/mcpserver/server.go`, `newServer`). Over stdio, the orchestrator's `flai mcp` inherits it. The HTTP server is a child of `flai host`, so its `FLAI_ROLE` is empty, and `queue_add` and `release_publish` refuse an orchestrator that calls them from a container.

## Decision

The orchestrator accepts through an MCP tool, `accept`, that runs `flai accept --by orchestrator` on the host, and verifies through the close-out job; its shell's `flai accept` and tier-running `flai verify` refuse in a role's container, naming those tools.

1. **The tool `accept`.** It takes the story, `verified` (the commit the orchestrator's verifier checked), `evidence` (as [ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md) takes it), and `dry_run`. It runs the MCP server's own flai, on the host and with the host's configuration, as `hostapi.Publisher` runs `flai release --pending`: `flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence -`, the evidence on standard input, with arguments built as `accept.run` builds the orchestrator's. With `dry_run` it runs `flai accept <S-nnnn> --by orchestrator --dry-run --json` and answers the preview, its blockers and the orchestrator's among them. The acceptance therefore sees the host's actions: it enqueues in place of landing while the merge queue batches, and needs the `forge` host action under pull-request acceptance, as `accept.run` does. It answers flai's result, or flai's refusal as the tool's error, and journals the run under the method `mcp.accept`.
2. **Who calls it.** Only the orchestrator, and only while `orchestration.permissions.accept_reviews` is on. `flai guard` refuses it to every other role and every sub-agent, as it refuses `queue_add`. The handler refuses it too, as `queue_add`'s does, to any caller the server does not know as the orchestrator.
3. **The orchestrator over HTTP.** The HTTP MCP server knows a caller as the orchestrator when its `X-Flai-Agent` names the orchestrator's run that flai serve records for the project and that is still running, as [ADR-0243](0243-over-http-permission-prompt-honours-the-auto-approve-host-action-only-for-the.md) knows a story's agent. The name alone is not enough, since every agent holds the same token. This makes `accept`, `queue_add`, and `release_publish` reach the orchestrator in a container. Over stdio, `FLAI_ROLE` still says it.
4. **Verifying.** The orchestrator verifies a story with the MCP tool `verify`, which runs the close-out job on the host or answers from the verify record, and reads the verdict with `flai verify <S-nnnn> --last`. The commit it gives `accept` is the commit that verdict names.
5. **The shell fails closed.** When `FLAI_ROLE` is `orchestrate` and the configuration flai reads (`--config`, else `FLAI_CONFIG`) lies under the project's `.flai-cache/agents/`, `flai accept`, its `--dry-run` included, and a `flai verify` that runs tiers refuse before they read or change anything. The refusal says that the host's actions cannot be read in a role's container, and names the MCP tools `accept` and `verify`. `flai verify --last` and `flai verify --job` still run: the one reads the record, the other asks flai serve. The operator's runs, a story's agent's, and the orchestrator's on the host are unchanged.

## Consequences

- An orchestrator in a container accepts as one on the host does: through the queue while it batches, by hand otherwise, and under the same conditions.
- A shell acceptance in the container can no longer land a story while the queue batches, and no tier runs in an image that lacks its tools.
- The orchestrator's role text, in its three copies, and the strategic-agents convention, in both, name `accept` and `verify` in place of the shell's `flai accept` and tier-running `flai verify`. `docs/users/flai.md` lists the tool.
- The HTTP MCP server reads flai serve's agent state to know the orchestrator, as it does for a story's agent under auto-approve.
- A story's agent in a container still runs `flai verify` there ([ADR-0253](0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md), [ADR-0274](0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). With no host actions, it runs the `queue` tiers too: more work, not a wrong acceptance. This ADR leaves that alone.

## Alternatives considered

- **Mount or copy the host's actions into the role's cache.** The cache is mounted read-write, so a shell write in the container could turn the copy's actions on or off without the guard seeing it, against [ADR-0220](0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md). Mounting `~/.flai` itself would give the container the operator's whole configuration, which ADR-0247 keeps out.
- **Have the guard refuse the shell's `flai accept` in a container.** The guard sees only the calls a hook reports, and a story's agent's guard is the image's flai. flai's own refusal holds wherever flai runs.
- **Pass `FLAI_ROLE` to the HTTP server in a header.** Any agent holding the token could then claim to be the orchestrator. The served run's name is the one flai serve gave it.
