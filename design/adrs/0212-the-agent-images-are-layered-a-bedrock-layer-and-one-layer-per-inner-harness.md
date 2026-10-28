---
id: ADR-0212
title: "The agent images are layered, a bedrock layer and one layer per inner harness, and pilot, a second Go binary of the flai module, is their entrypoint in place of flai agent run, built locally within E-0029"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0203]
---

# ADR-0212 The agent images are layered, a bedrock layer and one layer per inner harness, and pilot, a second Go binary of the flai module, is their entrypoint in place of flai agent run, built locally within E-0029

## Context

[ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) lets flai run a story's agent in a container. It names one image, `ghcr.io/bytepunx/flai-agent`, with every inner harness baked in, built and signed by a release pipeline. It names `flai agent run` as the entrypoint: a flai subcommand that reads a start request, starts the inner loop, and writes neutral usage events that `flai serve` reads.

E-0029 makes two refinements. The operator asked for layered images: a bedrock image with what every agent needs, and one image per harness on top of it, or per local model's tooling. And the epic names the entrypoint `pilot`, a binary of its own rather than a flai subcommand.

S-0505 built the first two layers by hand: `flai/images/bedrock/Dockerfile` and `flai/images/claude-code/Dockerfile`, built by `scripts/agent-image.sh`, with `tini` as the entrypoint and a shell as the command. Claude Code ran in the second image, driven by hand.

Claude Code's adapter on the host already writes the argv, the environment, and the MCP configuration a run needs, and `flai serve` already has a usage reader, a stream reader, and an activity reader for its output ([ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)). What the container lacks is a process that turns a start request into that run.

## Decision

The agent images are layered, and pilot, a second Go binary of the flai module, is their entrypoint in place of `flai agent run`. E-0029 builds the images locally.

1. **Layers.** The bedrock layer is a small Linux base, Alpine, with `git`, `flai`, `pilot`, and `tini`, built from `flai/images/bedrock/Dockerfile`. Each inner harness is one layer over bedrock, built from `flai/images/<layer>/Dockerfile`; `claude-code` is the first. A project's own layer, for the tools its test tiers need, builds over a harness layer (S-0510). A harness layer inherits bedrock's entrypoint.
2. **pilot.** `flai/pilot/main.go` builds it, from the package `flai/internal/pilot`. The entrypoint is `tini -- pilot`.
3. **The start request.** pilot reads one JSON object on standard input: `harness.Request` plus these fields.

   | Field | Holds |
   |---|---|
   | `inner` | The inner harness to run, a key of `harness.Adapters` |
   | `mcp_url_env` | The name of the variable holding the host's HTTP MCP server URL |
   | `mcp_token_env` | The name of the variable holding the agent's bearer token |
   | `mcp_host_env` | Optional: the name of the variable holding the `Host` header the server listens on |
   | `cache` | The mounted cache folder |

   The request carries names of variables, never their values. `mcp_host_env` exists because the server's MCP SDK refuses a request whose `Host` is not loopback, and a container dials `host.docker.internal`.
4. **The run.** pilot starts the inner harness through its adapter in `harness.Adapters`, with the argv and environment the adapter gives. It writes the MCP client configuration the adapter gives under the cache folder, naming the token's variable rather than its value. A session the cache already records is resumed.
5. **Output.** pilot copies the inner harness's standard output to its own line by line, unchanged. `flai serve`'s readers (usage, stream, activity) stay per harness, as ADR-0130 has them. pilot's own diagnostics go to standard error, each line prefixed `pilot:`.
6. **Signals and exit.** pilot forwards `SIGTERM` to the inner harness and exits with its exit code. A malformed request or an unknown inner harness exits 2, with one line on standard error.
7. **No neutral events yet.** pilot writes no usage events of its own. No inner harness needs them yet; the loop of flai's own that ADR-0203 describes would be the first.
8. **Local builds.** `scripts/agent-image.sh` builds the images on the operator's machine, as `flai-agent-bedrock:local` and `flai-agent-claude-code:local`. They are unstamped and unsigned, and no workflow builds or publishes them. The signed, published image, its digest list, and its release pipeline are left to a later story.

The rest of ADR-0203 stands: the mounts, the host's MCP server over HTTP, the guard and the permission prompt, usage and cost, and secrets by name.

## Consequences

- `flai serve` talks to a container as it talks to a child process: a start request in, the harness's own log out, an exit code at the end.
- The image, not flai, decides which inner harness runs. A new harness is a new layer and an adapter; bedrock does not change.
- pilot ships only in the image. GoReleaser's flai release builds `flai` alone, and an installed flai has no `pilot`.
- The guard runs inside the container as the hook the worktree's `.claude/settings.json` names, with the image's `flai`, not one the agent built. pilot sets `FLAI_AGENT` and `FLAI_SESSION` from the request, so the guard and the MCP server see one name.
- The token is in the container's environment, as a key is in a child's on the host. The MCP configuration on the mounted cache names the variable, so the cache holds no secret.
- Readers stay per harness, so a second inner harness needs its own readers, as on the host. Neutral events are deferred, not refused.
- Until the signed image exists, `flai serve` cannot verify an agent image by digest as ADR-0203 requires. A local image is the operator's own build, trusted as their checkout is.
- `flai` and `pilot` are built from one tree in one image, so they cannot drift apart.

## Alternatives considered

- **`flai agent run`, a flai subcommand, as ADR-0203 named.** One binary, but every flai release would carry container-only code, and the entrypoint's surface would be flai's whole command tree. A separate binary keeps the entrypoint small and its imports explicit.
- **A shell entrypoint script.** No Go to write, but the script would re-implement each adapter's argv, environment, and MCP configuration, and drift from them. Signal forwarding, exit codes, and JSON parsing are fragile in `sh`.
- **Neutral usage events instead of relaying the log.** One reader for every harness, as ADR-0203 drew it. But Claude Code's readers exist, and pilot would need a translation from each harness's log to the events, kept in step with both. Relay costs nothing now; neutral events can follow when the first loop that needs them does.
