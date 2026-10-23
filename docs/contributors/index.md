---
title: Contributors guide
updated: 2026-10-21
status: draft
---

# Contributing to system-flow

This monorepo builds itself with its own conventions. Start with the root `CLAUDE.md` and `wip/agents/index.md`, then the [design index](../../design/system/README.md).

| Part | Folder | Stack |
|------|--------|-------|
| Standard | `design/` | Markdown |
| Template | `template/` | Go text/template files, see the [template guide](template.md) |
| CLI | `flai/` | Go 1.26, see [go.md](../../design/tech/go.md) |
| Dashboard | `flaiover/` | SvelteKit 2, Svelte 5, Tailwind 4, see [sveltekit.md](../../design/tech/sveltekit.md) |

## Working a story

1. Pull a `ready` story from `wip/kanban/board.md` to `in-progress`.
2. Open its narrative under `wip/agents/`.
3. Keep the narrative's current state and next steps true as you go, with `flai stream state`.
4. Meet the definition of done in [workflow.md](../../design/system/workflow.md), move the story to `review`, open a pull request that names the story.

## Running the dashboard from source

This repository builds both sides from source: flai from `scripts/flai.sh`, with no release stamp, and the dashboard from `flai dashboard --build`, an image on no signed digest list. flai and the dashboard refuse such builds unless unsigned builds are allowed ([Allowing unsigned builds](../operators/index.md#allowing-unsigned-builds-s-0239)). Allow them in the configuration `scripts/flai.sh` uses, `.flai-cache/config.json`:

```sh
scripts/flai.sh config set dashboard.allow_unsigned true
scripts/flai.sh dashboard --build
```

When the `flai serve` that dials the dashboard runs with another configuration, such as the installed flai's `~/.flai/config.json`, set it there too: `flai config set dashboard.allow_unsigned true`. `flai serve` reads it before each dial and needs no restart. `--build` starts the container with `FLAIOVER_ALLOW_UNSIGNED=1`, so the dashboard accepts a flai with no stamp.

Expect it to be shown, never silent:

- A banner on every dashboard page, headed `unsigned allowed: flai`, `unsigned allowed: dashboard`, or `unsigned allowed: flai and dashboard`, and the same words beside the project in the switcher.
- `unsigned allowed: <sides>` in `flai serve status`, `flai host status`, and `flai dashboard status`.
- One `warn` per connection in `serve.log`, `connection allowed unsigned`.

## Tools and environment

Every script in `scripts/` sources `scripts/env.sh`, which sets this environment:

| What | Value |
|------|-------|
| `PATH` | The tree's own `bin/`, then the main checkout's `bin/` (from a story worktree under `.flai-cache/worktrees/`), then the main checkout's pnpm, then `~/go/bin` |
| `FLAI_CONFIG` | `.flai-cache/config.json` in the main checkout, unless already set |
| `FLAI_CACHE_DIR`, npm and pnpm caches | Under the main checkout's `.flai-cache/` |
| `LITELLM_BASE_URL`, `LITELLM_API_KEY`, `FLAI_TEST_LITELLM_HAIKU` | From `.flai-cache/litellm/litellm.env` when `scripts/litellm.sh up` has written it and the environment does not already set them: the local LiteLLM proxy and its budgeted virtual key, for the gateway tests and the agents `flai serve` starts; `up` writes `FLAI_TEST_LITELLM_HAIKU` as `LITELLM_HAIKU_MODEL` when that is set, and the script's default otherwise |

`FLAI_CONFIG` also decides which host actions the tree's flai sees, and the merge queue is one of them. While the merge queue can batch, the close-out leaves the `integration` and `smoke` tiers to the queue's batch ([ADR-0174](../../design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)). `scripts/close-out.sh` runs the tree's flai, so the host action must be on in this configuration too, not only in flai serve's. Turn it on from the main checkout:

```sh
scripts/flai.sh serve enable merge-queue
```

Without it, every close-out that changes `flai/` runs both tiers, about twelve minutes a round, which is longer than an agent's ten-minute tool call. `flai verify` then prints, before its outcome, the tiers it ran and the condition that is off (S-0466).

Run `make install-tools` once, from any tree: it installs the pinned golangci-lint v2 and GoReleaser into the main checkout's `bin/`, so every story worktree lints with them. `make flai-test` names the golangci-lint it runs and stops if it is not v2.

### Gateway checks

The gateway checks run `claude -p` through a gateway with the `claude-code` adapter's arguments, `flai guard`, and `permission_prompt` on, and log what the gateway charged beside Claude Code's `total_cost_usd`. Each run spends real money on the gateway's key. `scripts/gateway-smoke.sh`, the smoke tier's last step, runs the check of each gateway whose variables are set and prints one line for each it skips, naming the variables missing. The integration tier leaves them to it, and `make test` never runs them. A close-out whose branch changes `flai/` runs the smoke tier, so on a host with a key each such close-out spends one run.

| Variable | Default | What |
|----------|---------|------|
| `OPENROUTER_API_KEY` | none: the check skips | The OpenRouter key the check runs on |
| `FLAI_TEST_OPENROUTER_HAIKU` | `anthropic/claude-haiku-4.5` | The OpenRouter model `haiku` resolves to |
| `LITELLM_BASE_URL` | none: the check skips | The LiteLLM proxy's root, such as `http://127.0.0.1:4000`; the check runs through its unified `/v1/messages` |
| `LITELLM_API_KEY` | none: the check skips | The proxy's virtual key the check runs on |
| `FLAI_TEST_LITELLM_HAIKU` | The first that answers of `litellmHaikuDefaults` in `gateway_test.go`, in order, each as `openrouter-messages/<slug>`: `nvidia/nemotron-3-ultra-550b-a55b:free`, `minimax/minimax-m3.1-flash:free`, `nvidia/nemotron-3.5-lightning:free` | The proxy's model `haiku` resolves to on the unified route. Unset, or a slug on the list, the check probes the list in order, runs on the first slug the proxy answers 2xx, and logs each it passed over; when none answers, it fails, naming each slug and its answer. A slug not on the list is used alone, and fails the check when the proxy refuses it. The defaults are free OpenRouter models, so the check spends nothing. `scripts/litellm.sh up` sets it in the key file: `admin.env` keeps `LITELLM_HAIKU_MODEL` only when it was set, and drops it when it is a default the script has retired |
| `FLAI_TEST_LITELLM_PASSTHROUGH` | unset: the pass-through skips | `1` runs the check through the proxy's `/anthropic` pass-through too, on `claude-haiku-4-5`; it forwards to Anthropic only, so it needs the proxy's upstream `ANTHROPIC_API_KEY` and spends on it |
| `LITELLM_MASTER_KEY` | none: the spend log is not read | The proxy's master key, the only key that reads `/spend/logs`; without it the check logs the virtual key's spend from `/key/info` alone |
| `FLAI_TEST_GATEWAY_BUDGET` | `1.00` | The run's `max_budget_usd`, held to Claude Code's own count, which runs three to four times what OpenRouter charges; a run costs about $0.06 at OpenRouter |

Give the key a credit limit of its own on OpenRouter's keys page, so that a run gone wrong cannot spend past it, then run one check:

```sh
OPENROUTER_API_KEY=sk-or-... FLAI_TEST_GATEWAY_BUDGET=0.50 scripts/gateway-smoke.sh
```

For LiteLLM, `OPENROUTER_API_KEY=sk-or-... scripts/litellm.sh up` starts a local proxy that reaches OpenRouter's models as `openrouter-messages/<slug>`, in Anthropic's format through OpenRouter's Anthropic API, Anthropic's too when `ANTHROPIC_API_KEY` is set, and issues a virtual key with a `max_budget` (`LITELLM_MAX_BUDGET`, default 5 USD per 30 days), so a run gone wrong cannot spend past it. `scripts/env.sh` reads the key file it writes, so `scripts/gateway-smoke.sh` runs the LiteLLM check with no variable set by hand. The master key stays in `.flai-cache/litellm/admin.env`; set `LITELLM_MASTER_KEY` from it for one run when you want the spend log's rows:

```sh
LITELLM_MASTER_KEY="$(sed -n 's/^LITELLM_MASTER_KEY=//p' .flai-cache/litellm/admin.env)" scripts/gateway-smoke.sh
```

What each check found is in [agent-adapters.md](../../design/system/agent-adapters.md) § Checked against OpenRouter and § Checked against LiteLLM.

### Forge check

`flai/tests/integration/forge_github_test.go` reads a real GitHub repository through the forge adapter ([ADR-0165](../../design/adrs/0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md)). It writes nothing to the repository, and it is skipped, naming what is missing, unless the repository and the token are set.

| Variable | Default | What |
|----------|---------|------|
| `FLAI_FORGE_TEST_REPO` | none: the check skips | The repository, as `owner/name`, such as `bytepunx/system-flow` |
| `GITHUB_TOKEN` | none: the check skips | A token that can read the repository; `GITHUB_TOKEN="$(gh auth token)"` lends it `gh`'s for one run |
| `FLAI_FORGE_TEST_HEAD` | unset: checks and reviews are not read | A branch of the repository with an open pull request, whose checks and reviews the check reads |
| `FLAI_FORGE_TEST_API` | `https://api.github.com` | The API's root, for GitHub Enterprise |

An agent `flai serve` starts gets flai serve's environment without `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG`, plus `FLAI_AGENT`, `FLAI_STORY`, `FLAI_SESSION`, and `FLAI_STARTED_BY`. A `flai serve` it runs to try a change, with a `--config` of its own, does not reach the operator's host. `scripts/flaiover-dev.sh` gives the dev server the dashboard token file that `flai dashboard token --json` names.

## Changing a decision

Write a new ADR in `design/adrs` that supersedes the old one, then update `design/system` and `design/tech` in the same pull request.

## Releasing

Releases are per component and computed from what was accepted, never tagged by hand; the rules are in [git.md](../../design/conventions/git.md). Publishing, `flai release --pending` after a fetch or the board's Publish action, bumps each component that has accepted work since its last release, and is the one way accepted work reaches the remote ([ADR-0067](../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)):

- `flai`: tag `flai/vX.Y.Z`, GoReleaser publishes binaries.
- `flaiover`: tag `flaiover/vX.Y.Z`, the image workflow publishes to GHCR.
- Template: the version in `template/template.yaml` and an entry in `template/CHANGELOG.md`, then published to its own repository.

## The template

The template is developed in `./template` and published to [bytepunx/system-flow-template](https://github.com/bytepunx/system-flow-template), which is flai's default source. This repository keeps `template.repo: ./template` in `system-flow.yaml` so it develops against the working copy. How to edit, test, version, and publish it, here or in a fork: the [template guide](template.md).
