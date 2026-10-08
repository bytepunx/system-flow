---
title: Agent adapters for LiteLLM and OpenRouter
updated: 2026-10-08
status: active
topics: [cli, template]
---

# Agent adapters for LiteLLM and OpenRouter

The finding of S-0339. The operator asked for research toward an epic for agent adapters for LiteLLM and OpenRouter, and for the abstractions that free flai's intent, start an agent on a story, keep it to what the operator allows, and measure what it spent, from the specifics of one vendor: what Anthropic's SDK and Claude Code give against what OpenAI's SDK and the harnesses built on it give. This document says where flai is tied to Anthropic and Claude Code today, what LiteLLM and OpenRouter offer and the harnesses that drive them, the abstractions that would let another harness stand where Claude Code stands, and ends with a recommendation. Nothing here is built. The operator's decision and the epic are recorded at the end.

Every statement about this repository comes from its files, read on 2026-10-08, and names the file or ADR it comes from. Design documents are linked; code is named by its path.

## Today

flai has one harness it knows how to drive, `claude-code`, and one escape hatch, `command`, which runs whatever the operator writes. Everything below `flai serve` that reads an agent's output, guards its calls, or answers its permission prompts reads Claude Code's formats. The MCP server, `flai mcp`, is the one piece every harness could share as it stands: it speaks the Model Context Protocol over stdio and knows nothing of who called it. Each subsection names what is neutral already and what only Claude Code provides.

### Starting and resuming an agent

An agent is a harness, a model, and a flat `config` map, kept in `system-flow.yaml` as the project's default and copied into a story's front matter when it is made ([ADR-0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)). `flai serve` starts the story's agent through an adapter for its harness ([ADR-0038](../adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md)).

The adapter boundary is `flai/internal/harness/harness.go`:

| Piece | What it is |
|-------|------------|
| `Adapter` | Two methods: `DefaultHost() Host`, the program and arguments when the operator set none, and `Start(req Request, host Host) (Start, error)`, which builds the command and refuses what the story says that the harness does not take. |
| `Request` | What an agent is started for: `Story`, `Root`, `Project`, `Agent` (the `manifest.Agent`), `Name` (the `FLAI_AGENT` it works under), `Flai` (this executable, the agent's MCP server), `Role`, `Item`, `Focus`, `Session`, `Answered`, `Restart`, `AutoRestart`, `Commit`, `Started`, `Past`, `Begun`. |
| `Host` | The operator's say on the host for one harness: `Program` and `Args`, set with `flai serve agent harness <name>` and the dashboard's `settings.harness` (`flai/internal/hostapi/settings.go`). |
| `Start` | `Harness`, `Argv` (program first), `Env` added to `flai serve`'s own. Run as it stands, never through a shell. |
| `Adapters` | A map of two: `ClaudeCode` (`"claude-code"`) and `Command` (`"command"`). `For` picks one by the story's `agent.harness`, falls back to `command` when the story names none and the operator set one, and refuses any other name with `flai cannot start harness %q; it can start claude-code, command`. |
| `Prompt` | The start prompt. Only the `claude-code` adapter sends it; the `command` adapter has no prompt and is told its context in environment variables. |

The two adapters are in `flai/internal/harness/adapters.go`:

| | `claude-code` | `command` |
|-|---------------|-----------|
| Program | `claude`, or the operator's | Whatever the operator set with `flai serve agent set -- <program> [args...]`; no default |
| Argument list | `claude -p <prompt> --output-format stream-json --verbose --mcp-config <json> --strict-mcp-config --name "<project> <story>"`, then `--session-id <uuid>` or `--resume <uuid>`, then `--model <model>`, the config flags, `--agents <json>`, `--agent <definition>` for a strategic agent, then `--permission-prompt-tool mcp__flai__permission_prompt` and the operator's arguments, by default `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` | The operator's arguments with `{story}`, `{root}`, `{model}`, and `{harness}` replaced |
| Prompt | Positional, after `-p` | None |
| Model | `--model`, passed through unread | `{model}` and `FLAI_MODEL` |
| Config | Three keys, each checked against a pattern: `effort` (`low`, `medium`, `high`, `xhigh`, `max`) as `--effort`, `max_budget_usd` as `--max-budget-usd`, `fallback_model` as `--fallback-model`; any other key refuses the start | Not interpreted: `FLAI_AGENT_CONFIG`, a JSON object |
| Roles | `--agents`: the project's `.claude/agents/explorer.md` and `verifier.md` read as JSON definitions with the role's model laid over them ([ADR-0065](../adrs/0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md)); a role on another harness, or with config, refuses the start | `FLAI_AGENT_ROLES`, a JSON object |
| MCP | `--mcp-config` with one server, `flai` on stdio, `<this flai> mcp --agent <name>`, and `--strict-mcp-config` so the project's `.mcp.json` and the operator's own servers are not loaded | Nothing: the command reaches flai however it likes |
| Session | A random UUID `flai serve` makes (`newSession` in `flai/internal/serve/agents.go`, "the form Claude Code's sessions take"), passed as `--session-id` on a first start and as `--resume` when the agent is started again after its question was answered | Ignored; `FLAI_ANSWERED` carries the thread's ID instead |
| Environment | `FLAI_ROLE`, and `FLAI_ITEM` for the planner; `flai serve` adds `FLAI_AGENT`, `FLAI_SESSION`, `FLAI_STORY`, and `FLAI_STARTED_BY=flai-serve` | The same, plus `FLAI_MODEL`, `FLAI_HARNESS`, `FLAI_AGENT_CONFIG`, `FLAI_AGENT_ROLES`, `FLAI_ANSWERED`, `FLAI_COMMIT`, `FLAI_FOCUS` |

The manifest checks only shape (`flai/internal/manifest/agent.go`): a harness matches `^[a-z][a-z0-9._-]{0,63}$`, a model `^[A-Za-z0-9][A-Za-z0-9._:/@+-]{0,127}$`, and a config key is a lower-case word with a one-line value. Its error text names `claude-code` and `claude-opus-5-5` as examples, and the known roles are `explore` and `verify`. Whether a harness exists is decided only in `harness.For`, and only there is the list `claude-code, command`. No model name is validated anywhere: what the story writes goes to the harness as it stands.

`flai serve` reads the agent's standard output as Claude Code's stream-json, one event per line (`flai/internal/serve/stream.go`): a `system` event with subtype `init` is the session's start, `assistant` and `user` events are turn calls and tool results, `result` is the run's end, and any line that is not an event is kept as plain output. The same log feeds the live stream the dashboard shows and the usage measurement below. It is written under flai serve's state as `<key>-<story>-<start>.log`, and strategic runs as `<key>-planner-<start>.log`, `<key>-orchestrator-<start>.log`, and `<key>-analyzer-<start>.log` ([ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md), [ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md), [ADR-0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md)).

Neutral already: the `Adapter` interface, `Request`, `Host`, and `Start`; the manifest's agent shape and its role map; the operator's per-harness program and arguments; the MCP server the agent is given; the `command` adapter's environment contract. Claude Code's: every flag in the argument list, the three config keys, the sub-agent definitions as `--agents` JSON, the UUID session and `--resume`, and the stream-json reader every later step rests on. The `command` adapter gives flai no session to resume, no events to read, and no permission handler: an agent started through it is a black box that ends.

### Permissions

Claude Code run headless has no one to ask when a tool call needs a person's approval, so `flai serve` names flai's own MCP tool as the one Claude Code asks: `--permission-prompt-tool mcp__flai__permission_prompt`, passed by `permissionArgs` in `flai/internal/harness/adapters.go` ahead of the operator's arguments ([ADR-0086](../adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md)).

The tool is `flai/internal/mcpserver/permission.go`. Its input and output are Claude Code's permission-prompt-tool contract:

| Direction | Shape |
|-----------|-------|
| In | `tool_name` (required), `input` (the tool call's input), `tool_use_id`, and flai's own `project` |
| Out | `{"behavior":"allow","updatedInput":<input>}` or `{"behavior":"deny","message":"<reason>"}` |

It approves one thing: an `Edit`, `Write`, `MultiEdit`, or `NotebookEdit` of a path Claude Code protects, inside an in-progress story's worktree, never a path in `.git`. Those paths are flai's copy of Claude Code's list, `flai/internal/protected/protected.go`: the folders `.git`, `.config/git`, `.vscode`, `.idea`, `.husky`, `.cargo`, `.devcontainer`, `.yarn`, `.mvn`, `.claude`, and files such as `.gitconfig`, `.bashrc`, `.npmrc`, `.pre-commit-config.yaml`, `.mcp.json`, and `.claude.json`. The same list drives the acceptance preview that leaves a story changing a protected path to the operator alone ([ADR-0106](../adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). Everything else is denied at once: the session's own `--allowedTools` and `--permission-mode` decide what a story's agent may do, as ADR-0038 has it.

Unless the project's `auto-approve` host action is on, the tool opens a thread on the story, `Allow <tool> <path>?`, showing the change, and waits at most four minutes (`permissionWait`) for an answer whose first word is `allow`, `yes`, `approve`, `approved`, or `ok` from the story's owner or the project's owner ([ADR-0097](../adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md)). Unanswered, it refuses the write and leaves the thread open, and takes the answer when the same write is made again ([ADR-0124](../adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md)). The four-minute bound exists because of Claude Code's idle timeout on an MCP call.

Because a Claude Code release once changed what it accepts from the tool, `flai serve` checks each `claude` version it has not seen (`flai/internal/serve/claudecheck.go`): it reads `claude --version`, and when the version has no record in `claude-checks.json` beside its configuration, it builds a scratch project with one story in progress and auto-approve on, runs one headless `claude -p` on the model `haiku` with the adapter's own permission arguments, asks it to `Write` `.claude/flai-check.md`, and records `passed` or `failed` with Claude Code's output. A failure opens a thread on `system-flow.yaml` of every served project whose agents run `claude-code`.

Neutral already: the decision itself, which files may be written, by whom, after whose answer, on which thread, and the `protected` list as a policy. Claude Code's: the tool being called by the harness rather than by the model, its `behavior`/`updatedInput`/`message` shape, the tool names it keys on, the list of protected paths being Claude Code's list, the four-minute bound, and the version check, which has no meaning for a harness that does not run `claude`.

### The guard

`flai guard` is a Claude Code hook ([ADR-0060](../adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)). `template/root/.claude/settings.json`, and this repository's `.claude/settings.json` which runs `"$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard` instead of the installed `flai guard`, register it four times:

| Hook | Matcher | Runs when | Does |
|------|---------|-----------|------|
| `PreToolUse` | `Bash\|mcp__flai__.*` | Every session | Refuses with exit 2 and the reason on standard error; exit 0 passes |
| `PreToolUse` | `Edit\|MultiEdit\|Write\|NotebookEdit` | `FLAI_ROLE` is `plan`, `orchestrate`, or `analyze`, or `FLAI_STORY` is set | The same |
| `SubagentStart` | none | `FLAI_STORY` is set | Records the sub-agent; prints nothing, exit 0 |
| `SubagentStop` | none | `FLAI_STORY` is set | Removes it; prints nothing, exit 0 |

The hook's input is Claude Code's hook JSON on standard input, read by `Event` in `flai/internal/guard/guard.go`: `hook_event_name`, `session_id`, `tool_name`, `tool_input`, `agent_id`, and `agent_type`. From `tool_input` it reads `command` for `Bash`, `file_path` and `notebook_path` for the file tools, and the fields of flai's own MCP tools, `id`, `to`, `type`, `draft`, `commit`, `recommendation`, `source`, and `cost_of_delay`. A sub-agent is recognised by one thing: `agent_id` is not empty. The story's agent, the planner, the orchestrator, and the analyzer are told apart by `FLAI_STORY` and `FLAI_ROLE` in the environment, which `flai serve` set.

What it keys on is Claude Code's tool vocabulary: `Bash`, `Edit`, `MultiEdit`, `Write`, `NotebookEdit`, and the MCP tool names under `MCPPrefix`, `mcp__flai__`. Within a `Bash` command it splits the shell line and judges each `git` and `flai` command by name. It refuses a sub-agent every flai MCP tool but the reads (`board`, `doc_get`, `doc_search`, `item_get`, `prime`, `thread_get`, `who_touches`, and a few more), every `git` and `flai` command that writes, and, while `auto-approve` is off, any write to a file under a `.claude/` folder ([ADR-0102](../adrs/0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md)). It holds a strategic agent's own calls to its role's permissions. It records each session's running sub-agents in `.flai-cache/guard/<session_id>.json` from the `SubagentStart` and `SubagentStop` hooks and refuses the story agent's `wait_for_events` while one runs and no thread is open ([ADR-0092](../adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)). It fails open: input it cannot read passes.

Neutral already: the policy, which caller may call which flai tool or command, and the shell-line parsing of `git` and `flai` commands. Claude Code's: the hook events and their JSON, the `settings.json` that registers them, `CLAUDE_PROJECT_DIR`, the built-in tool names, the `mcp__<server>__<tool>` naming, `agent_id` as the sign of a sub-agent, and the `SubagentStart` and `SubagentStop` record. A harness with no hooks has no guard at all: the guard runs only because Claude Code calls it before each tool.

### Sub-agents and strategic agents

The template ships five agent definitions in `template/root/.claude/agents/`, copied to this repository's `.claude/agents/`, each a markdown file with Claude Code's sub-agent front matter, `name`, `description`, `tools`, and `model`, and the prompt as its body:

| Definition | Role | Tools | Model | Decided in |
|------------|------|-------|-------|------------|
| `explorer.md` | `explore`: read-only search for a story's agent | `Read`, `Grep`, `Glob`, and seven flai reads | `haiku` | [ADR-0059](../adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) |
| `verifier.md` | `verify`: reviews the diff against the criteria | The explorer's and `Bash` | `sonnet` | ADR-0059 |
| `planner.md` | `plan`: plans one epic or story | Reads, `Bash`, `Agent`, and the item and thread tools | `inherit` | [ADR-0075](../adrs/0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md), [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md) |
| `orchestrator.md` | `orchestrate`: keeps the work moving | Reads, `Bash`, `Agent`, and the orchestration tools | `inherit` | ADR-0075, [ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md) |
| `analyzer.md` | `analyze`: writes one report under `design/analysis/` | Reads, `Edit`, `Write`, `Bash`, `Agent`, and the thread and issue tools | `inherit` | ADR-0075, [ADR-0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md) |

The `tools` list is an allowlist in Claude Code's vocabulary, built-in names and `mcp__flai__*` names, and the `model` is a Claude Code alias (`haiku`, `sonnet`) or `inherit`. The explorer and verifier run as Claude Code sub-agents inside the story's session, started by the Agent tool; `flai serve` passes their definitions as `--agents` with the story's role model laid over them (ADR-0065). The three strategic agents are sessions of their own that `flai serve` starts with `--agent planner`, `--agent orchestrator`, or `--agent analyzer`, so the definition is the source of their tools and instructions (`claudeCodeStrategic` in `flai/internal/harness/adapters.go`).

The conventions lean on the Agent tool directly. [delegation.md](../conventions/delegation.md) tells a story's agent to launch each task sub-agent in the foreground, "in Claude Code, the Agent tool's `run_in_background` set to false", to name the task's ID in the Agent tool's `description` so that usage is apportioned, and to "use a fork, which inherits your conversation, where the harness offers one". A sub-agent is told apart from its parent in the log by `parent_tool_use_id`, the ID of the `Agent` (or `Task`) `tool_use` that started it (ADR-0059).

Neutral already: the roles, what each may do, and their prompts, which are plain markdown. Claude Code's: the file format and location (`.claude/agents/`), the `tools` and `model` vocabularies, the Agent tool and its `description` and `run_in_background` fields, the fork, `--agents` and `--agent`, and `parent_tool_use_id` as the sign of a sub-agent's call.

### Usage and cost

Work items carry `usage`, tokens and cost per model, measured from the agents' logs ([ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). The reader is `flai/internal/usage/log.go`, and everything it reads is Claude Code's stream-json:

| Event | Fields read | Used for |
|-------|-------------|----------|
| `system` with `subtype: init` | `session_id` | A run's start |
| `assistant` | `message.id`, `message.model`, `message.usage` (`input_tokens`, `output_tokens`, `cache_read_input_tokens`, `cache_creation_input_tokens`), `message.content[].type`, `name`, `id`, `input.description`, `input.prompt`, and `parent_tool_use_id` | Each call to a model, which sub-agent made it, which turn it belongs to |
| `user` | `message.content[]` `tool_result` with `tool_use_id`, `is_error`, and `task_description` | Pairing a tool call with its result; empty wakes |
| `result` | `total_cost_usd`, `usage`, `modelUsage` keyed by model with `inputTokens`, `outputTokens`, `cacheReadInputTokens`, `cacheCreationInputTokens`, `costUSD` | A session's reported totals |

A session's newest `result` is its total. Calls after it, a run still going or one that died, are counted from their `assistant` events and priced at the model's blended rate, `Rates`, reported cost over reported tokens across the logs flai keeps, and marked estimated. flai has no price table of its own: with no `result` for a model ever, an estimate's cost stays at nothing. `flai/internal/metrics/` (`usage.go`, `spend.go`, `strategic.go`) sums the items' `usage` into cost per item, per agent hour, and spend over time, and names no model and no price.

A task is given each session's totals in the share of its calls' input and cache tokens that are its own: the calls of the sub-agents started for it, matched by the story's task ID in the `Agent` call's `description` or prompt, and an even share of the story agent's other calls while it was in progress ([ADR-0071](../adrs/0071-a-task-s-usage-is-the-calls-of-the-sub-agents-started-for-it-and-an-even-share.md)). An empty wake is a `tool_use` named `mcp__flai__wait_for_events` without `parent_tool_use_id` whose `tool_result` says `timed_out` with no events ([ADR-0105](../adrs/0105-a-story-s-empty-wakes-the-wait-for-events-calls-of-its-agents-that-timed-out.md)). A turn is one `assistant` message without `parent_tool_use_id` that calls a tool, gathered across Claude Code's repeats of its message ID, and classed in `flai/internal/usage/turns.go` as `test_runs`, `hand_edits`, `empty_wakes`, `ceremony`, or `work` by the tool names and shell commands it carries ([ADR-0116](../adrs/0116-when-flai-measures-a-story-s-usage-from-its-logs-it-classifies-each-turn-of-the.md)).

Neutral already: the `usage` schema on items (`models`, `input`, `output`, `cache_read`, `cache_write`, `cost`, `seconds`, `estimated`, `empty_wakes`, `turns`), the summing up the hierarchy, the metrics on top, and the five turn classes as ideas. Claude Code's: every field name above, `result` as the only source of a reported cost, `modelUsage` and `costUSD`, `parent_tool_use_id` for attribution, the repeated-message quirk, and the `Agent`/`Task` tool names the task apportionment matches on.

### Model names and files the harness reads

| Where | What names Claude |
|-------|-------------------|
| `system-flow.yaml` | This project's default agent is `claude-code` with `model: claude-opus-5-5`. A project made from the template has none until the operator sets one. |
| `flai/internal/manifest/agent.go` | Validates a model only by pattern. Its examples and error text say `claude-code`, `claude-opus-5-5`, `verify=sonnet`. |
| `flai/cmd/agent.go`, `flai/cmd/manifest.go`, `flai/internal/manifest/settings.go`, `flai/internal/hostapi/settings.go` | Help text and refusals use `claude-code`, `claude-opus-5-5`, `claude-sonnet-5`, `effort=high` as examples. The dashboard's `settings.default_agent` and `settings.manifest` accept `harness`, `model`, `config`, `roles` and refuse `command` as a harness to set in `settings.harness`. Nothing restricts the model. |
| `flaiover/src/lib/components/AgentFields.svelte` | Placeholders `claude-code` and `claude-opus-5-5` when the project has no default. |
| `flaiover/src/lib/viz/palette.ts` | `MODEL_SLOT` and the mark shapes key on `opus`, `sonnet`, `haiku`, `fable` to colour usage charts by model family. |
| `flai/internal/serve/claudecheck.go` | `claudeCheckModel = "haiku"`, the cheapest model for the permission check. |
| `template/root/.claude/agents/*.md` | `model: haiku`, `model: sonnet`, `model: inherit`. |
| `template/root/design/conventions/git.md` and `flai/cmd/issue.go` | The commit trailer `Co-Authored-By: Claude ... <noreply@anthropic.com>`. |
| `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/settings.md` | Examples: `claude-opus-5-5`, `claude-sonnet-5`, `claude-haiku-4-5`, `haiku`. |

The files a session reads are Claude Code's. `CLAUDE.md` at the root is the map every session loads; the harness prompt says `Follow CLAUDE.md, or AGENTS.md where there is no CLAUDE.md` (`flai/internal/harness/harness.go`), and `flai upgrade` merges the template's `CLAUDE.md` as a marker file (`flai/cmd/upgrade.go`). [agent-context.md](agent-context.md) surveys the equivalents, a root `AGENTS.md`, Copilot's and Cursor's rules, and keeps flai's answer in `flai prime`, which is neutral. `.mcp.json` names the `flai` server as `flai mcp` for the operator's own session; a session `flai serve` starts ignores it under `--strict-mcp-config`. The `.claude/` folder holds `settings.json` (the hooks) and `agents/`; both are protected paths, written only through `permission_prompt`.

### Summary

| Concern | Neutral already | Claude Code's |
|---------|-----------------|---------------|
| Starting and resuming | `Adapter`, `Request`, `Host`, `Start`; the manifest's agent and roles; per-harness program and arguments; the MCP server; the `command` adapter's environment | Every flag of `claude -p`; `effort`, `max_budget_usd`, `fallback_model`; `--agents`; UUID sessions and `--resume`; the stream-json reader |
| Permissions | The decision, the owners, the thread, the protected list as policy | `--permission-prompt-tool`; `behavior`/`updatedInput`/`message`; Claude Code's protected paths; the four-minute bound; the `claude --version` check |
| The guard | The policy per caller; the `git` and `flai` command parsing | `PreToolUse`, `SubagentStart`, `SubagentStop`; `settings.json`; the hook JSON; built-in tool names; `mcp__flai__*`; `agent_id` |
| Sub-agents and strategic agents | The roles and their prompts | `.claude/agents/` front matter; `tools` and `model` vocabularies; the Agent tool, `description`, `run_in_background`, the fork; `--agent`; `parent_tool_use_id` |
| Usage and cost | The `usage` schema, the hierarchy sums, the metrics, the turn classes | Every stream-json field; `result`, `modelUsage`, `costUSD`; the repeated-message quirk; `Agent`/`Task` names |
| Model names and files | The model pattern; `flai prime`; `.mcp.json` as MCP | Example and placeholder names; `haiku` for the check; chart colours by family; `CLAUDE.md`; `.claude/`; the commit trailer |

## LiteLLM and OpenRouter

What each gateway offers and which harnesses can run a story's agent over it, from the vendors' and harnesses' public documentation read on 2026-10-08. Every claim links its page. A claim the documentation does not settle is marked **to check**.

### LiteLLM

LiteLLM is two things: a Python SDK and a self-hosted proxy, which its documentation calls the AI gateway ([docs.litellm.ai](https://docs.litellm.ai/docs/)). Both are open source; the proxy runs from Docker or `litellm --config config.yaml`.

| Concern | SDK | Proxy (AI gateway) |
|---------|-----|--------------------|
| Hosting | A library in the caller's Python process | Self-hosted; an enterprise tier adds features, the gateway itself is free ([simple_proxy](https://docs.litellm.ai/docs/simple_proxy)) |
| APIs served | `completion()` and `responses()` in Python, OpenAI message shape in and out | OpenAI `/chat/completions`; OpenAI Responses API at `/v1/responses`, bridged to `/chat/completions` for providers without a native one, with `previous_response_id` for multi-turn ([response_api](https://docs.litellm.ai/docs/response_api)); an Anthropic-compatible `/v1/messages` that reaches every provider LiteLLM supports, `openai`, `anthropic`, `bedrock`, `vertex_ai`, `gemini`, `azure` and so on ([anthropic_unified](https://docs.litellm.ai/docs/anthropic_unified)); a pass-through `/anthropic/v1/messages` forwarded to Anthropic untranslated, Anthropic models only ([anthropic pass-through](https://docs.litellm.ai/docs/pass_through/anthropic_completion)) |
| Streaming | `stream=True` | `"stream": true` on every endpoint above, server-sent events; a WebSocket mode for the Responses API |
| Model names | `<provider>/<model>`: `anthropic/claude-sonnet-5`, `openai/gpt-5.6-terra`, `bedrock/us.anthropic.claude-sonnet-5`, `openrouter/anthropic/claude-sonnet-4` ([openrouter provider](https://docs.litellm.ai/docs/providers/openrouter)) | The names in `config.yaml`'s `model_list`: a `model_name` the caller sends, mapped to one or more `litellm_params.model` deployments in the SDK's form; `aliases` on a key map one name to another ([virtual_keys](https://docs.litellm.ai/docs/proxy/virtual_keys)) |
| Routing and fallbacks | `Router` with retries and fallbacks across deployments | `router_settings`: `fallbacks: [{"model-a": ["model-b"]}]`, `context_window_fallbacks`, `content_policy_fallbacks`, `num_retries`, `allowed_fails`, `cooldown_time`; load balancing across deployments of one `model_name`; `disable_fallbacks` per request or key ([reliability](https://docs.litellm.ai/docs/proxy/reliability)) |
| Tool calls | OpenAI `tools` and `tool_choice` translated to each provider | Passed through on every endpoint; on `/v1/messages` Anthropic's `tools` and `stop_reason: tool_use` ([anthropic_unified](https://docs.litellm.ai/docs/anthropic_unified)) |
| MCP | None | An MCP gateway: servers registered in `config.yaml` or the UI over streamable HTTP, SSE, or `transport: stdio` with `command` and `args`; clients reach them at `/mcp/` with `x-litellm-api-key` and per-server `x-mcp-<alias>-<header>` headers; with `require_approval: never` the proxy calls MCP tools itself inside a `/chat/completions` request ([mcp](https://docs.litellm.ai/docs/mcp)) |
| Authentication | Provider keys in the process environment, `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `OPENROUTER_API_KEY` | A master key, `LITELLM_MASTER_KEY` or `general_settings.master_key`, `sk-` prefixed; virtual keys from `POST /key/generate` with `models`, `max_budget`, `tpm_limit`, `rpm_limit`, `budget_duration`; provider keys stay in `config.yaml` as `litellm_params.api_key` or `os.environ/ANTHROPIC_API_KEY`; the caller sends its virtual key as `Authorization: Bearer` or `x-api-key` ([virtual_keys](https://docs.litellm.ai/docs/proxy/virtual_keys)) |
| Usage and cost | `usage` in the response; `completion_cost()` prices it from LiteLLM's model cost map | `usage` in the body; headers `x-litellm-response-cost`, `x-litellm-key-spend`, `x-litellm-call-id`, `x-litellm-model-id`, `x-litellm-attempted-fallbacks` ([response_headers](https://docs.litellm.ai/docs/proxy/response_headers)); every call written to `LiteLLM_SpendLogs` with key, user, team, model, tokens, and spend, read at `/spend/logs`, `/spend/logs/v2` (filters `api_key`, `user_id`, `model`), `/key/info`, `/user/daily/activity`; custom prices per deployment in `model_info` ([cost_tracking](https://docs.litellm.ai/docs/proxy/cost_tracking)) |

LiteLLM documents Claude Code against the proxy: `ANTHROPIC_BASE_URL` at the proxy root for the unified `/v1/messages`, or at `<proxy>/anthropic` for the pass-through, and `ANTHROPIC_AUTH_TOKEN` a master or virtual key, which bounds the models Claude Code may name; Bedrock upstreams want `CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS=1`, and non-Anthropic base URLs fill the context with MCP tools unless `ENABLE_TOOL_SEARCH=true` ([claude_responses_api](https://docs.litellm.ai/docs/tutorials/claude_responses_api)).

### OpenRouter

OpenRouter is a hosted service at `https://openrouter.ai/api/v1`; nothing runs on the operator's host ([quickstart](https://openrouter.ai/docs/quickstart)).

| Concern | OpenRouter |
|---------|------------|
| Hosting | Hosted only. Credits are bought up front; a 5.5% fee on credit purchases, and the provider's per-token rate ([Claude Code tutorial](https://openrouter.ai/blog/tutorials/claude-code-openrouter/)) |
| APIs served | OpenAI `/api/v1/chat/completions` and `/api/v1/completions`; an OpenAI-shaped Responses API at `/api/v1/responses` ([create a response](https://www.openrouter.ai/docs/api/api-reference/responses/create-a-response)); an Anthropic Messages endpoint at `/api/v1/messages` that takes any model on OpenRouter, with `models` and `fallbacks` arrays ([create a message](https://openrouter.ai/docs/api/api-reference/anthropic-messages/create-a-message)); `/api/v1/models`, `/api/v1/generation`, `/api/v1/key` |
| Streaming | `"stream": true` on each; server-sent events ending in `[DONE]`; usage in the final chunk ([usage accounting](https://openrouter.ai/docs/use-cases/usage-accounting)) |
| Model names | `<author>/<model>`: `anthropic/claude-sonnet-4`, `openai/gpt-5`; `~anthropic/claude-sonnet-latest` with a tilde is an alias that follows the family's newest version; suffixes `:nitro` (throughput) and `:floor` (price) ([provider routing](https://openrouter.ai/docs/features/provider-routing)) |
| Routing and fallbacks | A `provider` object per request: `order`, `allow_fallbacks` (default true), `only`, `ignore`, `sort` (`price`, `throughput`, `latency`), `require_parameters`, `data_collection`, `zdr`, `quantizations`; a `models` array for model fallbacks; BYOK endpoints are tried first ([provider routing](https://openrouter.ai/docs/features/provider-routing)) |
| Tool calls | OpenAI `tools` and `tool_choice`, passed to the provider; Anthropic `tools` on `/api/v1/messages` |
| MCP | None server-side. OpenRouter is a model API; the client connects to MCP servers and converts their tool definitions to OpenAI tools, by hand or with the `@openrouter/mcp` package ([MCP servers](https://openrouter.ai/docs/guides/coding-agents/mcp-servers)). A separate OpenRouter MCP server exposes OpenRouter's own data, models, prices, credits, to an editor ([mcp-server](https://openrouter.ai/docs/mcp-server)) |
| Authentication | `Authorization: Bearer <OPENROUTER_API_KEY>`; keys made at `openrouter.ai/keys` with a name and a credit limit ([authentication](https://openrouter.ai/docs/api-reference/authentication)). Provider keys, when the operator brings their own, live in OpenRouter's workspace integrations settings, not on the host; BYOK requests cost 5% of OpenRouter's list price for the same call, after a monthly allowance ([BYOK](https://openrouter.ai/docs/use-cases/byok)) |
| Usage and cost | Every response carries `usage` with `prompt_tokens`, `completion_tokens`, `total_tokens`, `cost` (credits charged), `cost_details.upstream_inference_cost` (BYOK only), `prompt_tokens_details.cached_tokens` and `cache_write_tokens`, `completion_tokens_details.reasoning_tokens`; the `usage: {include: true}` opt-in is deprecated because usage is now always included. `GET /api/v1/generation?id=<id>` returns the same after the fact. On `/api/v1/messages`, Anthropic's `input_tokens`, `output_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens` plus OpenRouter's `cost` and `cost_details` ([usage accounting](https://openrouter.ai/docs/use-cases/usage-accounting), [create a message](https://openrouter.ai/docs/api/api-reference/anthropic-messages/create-a-message)). `GET /api/v1/key` returns the key's `limit`, `limit_remaining`, `usage`, `usage_daily`, `usage_weekly`, `usage_monthly` ([limits](https://openrouter.ai/docs/api-reference/limits)) |

OpenRouter documents Claude Code and the Agent SDK against its Anthropic endpoint: `ANTHROPIC_BASE_URL="https://openrouter.ai/api"`, `ANTHROPIC_AUTH_TOKEN` the OpenRouter key, `ANTHROPIC_API_KEY=""` set empty rather than unset so Claude Code does not fall back to Anthropic, and `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL`, `ANTHROPIC_DEFAULT_HAIKU_MODEL`, `CLAUDE_CODE_SUBAGENT_MODEL` naming OpenRouter slugs such as `~anthropic/claude-opus-latest`. It says the integration "is only guaranteed to work with the Anthropic first-party provider"; non-Anthropic models are not supported through Claude Code ([Claude Code tutorial](https://openrouter.ai/blog/tutorials/claude-code-openrouter/), [Anthropic Agent SDK](https://openrouter.ai/docs/guides/community/anthropic-agent-sdk)).

### Claude Code over a gateway

Claude Code speaks the Anthropic Messages format to whatever `ANTHROPIC_BASE_URL` names, calling `/v1/messages?beta=true` and, optionally, `/v1/messages/count_tokens`; it also accepts a gateway in Bedrock's format (`ANTHROPIC_BEDROCK_BASE_URL`, `CLAUDE_CODE_USE_BEDROCK=1`, `CLAUDE_CODE_SKIP_BEDROCK_AUTH=1`), Vertex's (`ANTHROPIC_VERTEX_BASE_URL`, `CLAUDE_CODE_USE_VERTEX=1`, `CLAUDE_CODE_SKIP_VERTEX_AUTH=1`), and Foundry's (`ANTHROPIC_FOUNDRY_BASE_URL`, `ANTHROPIC_FOUNDRY_API_KEY`) ([gateway compatibility guide](https://code.claude.com/docs/en/llm-gateway-protocol), [connect to a gateway](https://code.claude.com/docs/en/llm-gateway-connect)). The credential goes in `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`), `ANTHROPIC_API_KEY` (`x-api-key`), or an `apiKeyHelper` command in settings (both headers, cached five minutes). A gateway must forward `anthropic-beta` and `anthropic-version` unchanged and the streaming events as they come; it may consume `x-claude-code-session-id`, `x-claude-code-agent-id`, and `x-claude-code-parent-agent-id`, which attribute a request to a session and a sub-agent. Anthropic "doesn't support routing Claude Code to non-Claude models through any gateway" ([other LLM gateways](https://code.claude.com/docs/en/llm-gateway)). Gateway model discovery, `CLAUDE_CODE_ENABLE_GATEWAY_MODEL_DISCOVERY=1`, reads `/v1/models` and keeps only IDs containing `claude` or `anthropic`.

| Concern | Claude Code |
|---------|-------------|
| Headless, resume | `claude -p "<prompt>"`, `--session-id`, `--resume <id>`, `--continue`; `--bare` skips hooks, `.mcp.json`, agents, and `CLAUDE.md`, loading only what flags name ([run programmatically](https://code.claude.com/docs/en/headless)) |
| MCP | `--mcp-config <file-or-json>` with stdio servers, `--strict-mcp-config`; `system/init` reports `mcp_servers` and `mcp_server_errors` |
| Sub-agents, own model | `.claude/agents/*.md` or `--agents <json>`; `model` is `haiku`, `sonnet`, `opus`, `fable`, a full ID, or `inherit`; `CLAUDE_CODE_SUBAGENT_MODEL` sets the default ([subagents](https://code.claude.com/docs/en/sub-agents)) |
| Tool-refusal hook | `PreToolUse` in `settings.json`: exit 2, or JSON `hookSpecificOutput.permissionDecision: deny`; input has `session_id`, `tool_name`, `tool_input`, `agent_id`; `SubagentStart` and `SubagentStop` ([hooks](https://code.claude.com/docs/en/hooks)) |
| Permission handler | `--permission-prompt-tool mcp__<server>__<tool>`, `--permission-mode`, `--allowedTools`, `--permission-prompts none` for unattended runs; `permissionDecision: ask` from a hook escalates |
| Usage stream | `--output-format stream-json --verbose`: `assistant` events with `message.usage`, `parent_tool_use_id`; `result` with `total_cost_usd` and `modelUsage`; "client-side estimates" that can differ from the bill ([run programmatically](https://code.claude.com/docs/en/headless)) |
| Instruction files | `CLAUDE.md` at the root and in `~/.claude/`; agents in `.claude/agents/` |

What changes over a gateway: the model names in `--model`, the role models, and `--fallback-model` must be names the gateway serves, and `haiku`, `sonnet`, `opus` resolve only through the `ANTHROPIC_DEFAULT_*_MODEL` variables; `total_cost_usd` is Claude Code's own estimate from Anthropic's prices, not what LiteLLM or OpenRouter charged; and the models stay Claude.

### OpenAI Codex CLI

Codex talks to any OpenAI-shaped endpoint through `[model_providers.<id>]` in `~/.codex/config.toml` or `.codex/config.toml`: `base_url`, `env_key` naming the environment variable that holds the key, `wire_api = "chat"` or `"responses"`, `http_headers`, `env_http_headers`, `request_max_retries`; the session picks it with `model_provider` and `model` ([config](https://learn.chatgpt.com/docs/config-file/config-advanced)). LiteLLM's `/v1/responses` or `/chat/completions` and OpenRouter's `/api/v1/responses` or `/api/v1/chat/completions` both fit; which `wire_api` each gateway serves well is **to check**. Codex has no Anthropic-format client.

| Concern | Codex |
|---------|-------|
| Headless, resume | `codex exec "<prompt>"` or stdin; `codex exec resume <SESSION_ID>` or `--last`; `--sandbox read-only`, `workspace-write`, `danger-full-access`; `-c key=value` overrides; `--output-last-message`, `--output-schema`; a required MCP server that fails to start ends the run ([non-interactive mode](https://learn.chatgpt.com/docs/non-interactive-mode)) |
| MCP | `[mcp_servers.<name>]` with `command`, `args`, `env` (stdio) or `url`, `bearer_token_env_var`; `enabled`, `startup_timeout_sec`, `tool_timeout_sec`, `enabled_tools`, `disabled_tools`, `default_tools_approval_mode`; `codex mcp add <name> -- <command>` ([MCP](https://learn.chatgpt.com/docs/extend/mcp?surface=cli)) |
| Sub-agents, own model | On by default; built-in `default`, `worker`, `explorer`; custom agents as TOML in `.codex/agents/` or `~/.codex/agents/` with `name`, `description`, `developer_instructions`, and optional `model`, `model_reasoning_effort`, `sandbox_mode`, `mcp_servers`; `agents.default_subagent_model`, `agents.max_concurrent_threads_per_session` ([subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents)) |
| Tool-refusal hook | `PreToolUse` in `hooks.json` (`~/.codex/`, `.codex/`) or `[hooks]` in `config.toml`; covers Bash, `apply_patch` edits, and MCP tools; denies by exit 2 or `"permissionDecision": "deny"`, with `allow` and `ask`; input `session_id`, `tool_name`, `tool_input`, `cwd`; also `SubagentStart`, `SubagentStop`. Non-managed hooks need trust via `/hooks` or `--dangerously-bypass-hook-trust` ([hooks](https://learn.chatgpt.com/docs/hooks)) |
| Permission handler | `approval_policy`: `never`, `on-request`, `untrusted`, `on-failure`, or a `granular` table; headless, approvals follow the sandbox and policy, with no documented hook into an external approver. `ask` from a `PreToolUse` hook in `codex exec` is **to check** |
| Usage stream | `--json`: JSONL with `thread.started`, `turn.started`, `item.started`, `item.completed`, `turn.completed` carrying `usage` with `input_tokens`, `output_tokens`, `cached_input_tokens`, `turn.failed`, `error`. Tokens, no cost, no per-model split |
| Instruction files | `AGENTS.md` from the repository root down, `AGENTS.override.md`, `~/.codex/AGENTS.md`; `project_doc_fallback_filenames` can add `CLAUDE.md`; 32 KiB cap `project_doc_max_bytes` ([AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md)) |

### OpenCode

OpenCode reaches a gateway through `provider.<id>` in `opencode.json`: `npm: "@ai-sdk/openai-compatible"` for `/v1/chat/completions`, `npm: "@ai-sdk/openai"` for `/v1/responses`, `options.baseURL`, `options.apiKey` with `{env:VAR}`, `options.headers`, and a `models` map; OpenRouter is a first-class provider with per-model `provider.order` and `allow_fallbacks` options; models are named `provider/model` ([providers](https://opencode.ai/docs/providers/), [config](https://opencode.ai/docs/config/)). An Anthropic-compatible provider entry exists through the `@ai-sdk/anthropic` package; whether it accepts a custom `baseURL` toward LiteLLM's `/v1/messages` is **to check**.

| Concern | OpenCode |
|---------|----------|
| Headless, resume | `opencode run "<prompt>"`, `--session <id>`, `--continue`, `--fork`, `--agent`, `--model provider/model`, `--format json`, `--auto`; `--attach` to a running server avoids MCP startup on each call ([CLI](https://opencode.ai/docs/cli/)) |
| MCP | `mcp.<name>` with `type: "local"`, `command: [...]`, `environment`, `cwd`, `timeout`, or `type: "remote"` with `url`, `headers`, `oauth` ([MCP servers](https://opencode.ai/docs/mcp-servers/)) |
| Sub-agents, own model | `agent.<name>` in `opencode.json` or markdown in `.opencode/agents/` and `~/.config/opencode/agents/`; `mode: subagent`, `model`, `prompt`, `tools`, `permission`, `steps`; a subagent without `model` inherits the caller's; the Task tool gated by `permission.task` ([agents](https://opencode.ai/docs/agents/)) |
| Tool-refusal hook | A plugin in `.opencode/plugins/` or `~/.config/opencode/plugins/`: `tool.execute.before` throws to block the call; `permission.asked` and `permission.replied` are events ([plugins](https://opencode.ai/docs/plugins/)). Whether a plugin can tell a sub-agent's call from the primary's is **to check** |
| Permission handler | `permission` map: `allow`, `ask`, `deny` per tool, with glob rules per command or path, and per agent; `opencode run --auto` approves whatever is not denied, otherwise an `ask` in headless mode has no answerer ([permissions](https://opencode.ai/docs/permissions/)). No documented external approval callback for the CLI; the server API may offer one, **to check** |
| Usage stream | `--format json`: JSONL with `step_start`, `text`, `tool_use`, `step_finish`; `step_finish` carries tokens (input, output, reasoning, cache read and write) and `cost` in USD, priced from OpenCode's own model table (third-party write-ups; the official page lists only the formats, so the fields are **to check**) |
| Instruction files | `AGENTS.md` up the tree and `~/.config/opencode/AGENTS.md`; falls back to `CLAUDE.md` and `~/.claude/CLAUDE.md`; `instructions` in `opencode.json` adds paths, globs, or URLs ([rules](https://opencode.ai/docs/rules/)) |

### Goose (Block)

Goose has providers for both gateways: `GOOSE_PROVIDER=litellm` with `LITELLM_HOST`, `LITELLM_BASE_PATH`, `LITELLM_API_KEY`, `LITELLM_CUSTOM_HEADERS`; `GOOSE_PROVIDER=openrouter` with `OPENROUTER_API_KEY`; `openai` with `OPENAI_HOST` for any OpenAI-compatible endpoint; `anthropic` with `ANTHROPIC_HOST`; and declarative JSON providers for OpenAI-, Anthropic-, or Ollama-shaped APIs; `GOOSE_MODEL` names the model ([providers](https://goose-docs.ai/docs/getting-started/providers)).

| Concern | Goose |
|---------|-------|
| Headless, resume | `goose run -t "<text>"` or `-i <file>` (`-` for stdin), `--recipe`, `-n <name>`, `-r` to resume the named session, `--no-session`, `--max-turns`, `--quiet`, `--provider`, `--model` ([running tasks](https://goose-docs.ai/docs/guides/running-tasks), [man page](https://www.mankier.com/1/goose-run)). Sessions are named, not UUIDs |
| MCP | `--with-extension "ENV=value <command> <args>"` adds a stdio server for the run; `--with-streamable-http-extension <url>`; `--with-builtin` |
| Sub-agents, own model | Subagents spawned on request or by recipe; `GOOSE_SUBAGENT_PROVIDER` and `GOOSE_SUBAGENT_MODEL` set defaults, a recipe's `settings` override them; extensions inherited from the parent unless restricted ([subagents](https://goose-docs.ai/docs/guides/context-engineering/subagents)) |
| Tool-refusal hook | Hooks in `hooks/hooks.json` under `~/.agents/plugins/<name>/` or `<project>/.agents/plugins/<name>/`, the Open Plugins specification: `PreToolUse`, `BeforeShellExecution`, `PostToolUse` and others; stdin JSON with `event`, `session_id`, `tool_name`, `tool_input`; deny by exit 2 or `{"decision":"block","reason":"..."}` ([hooks](https://goose-docs.ai/docs/guides/context-engineering/hooks/)). No `agent_id`: whether a hook can tell a subagent's call is **to check** |
| Permission handler | `GOOSE_MODE`: `auto`, `approve`, `smart_approve`, `chat`; per-tool permissions in `permission.yaml`. What `approve` does under `goose run` with no terminal is **to check** ([permissions](https://goose-docs.ai/docs/guides/managing-tools/goose-permissions)) |
| Usage stream | `--output-format json` (one document at the end) or `stream-json` (events as they happen); token totals in the session store `~/.local/share/goose/sessions/sessions.db` and `goose session export`; no cost ([logs](https://goose-docs.ai/docs/guides/logs)). The event fields are **to check** |
| Instruction files | `AGENTS.md` then `.goosehints` at each directory level, `~/.config/goose/.goosehints` globally; `CONTEXT_FILE_NAMES` (default `["AGENTS.md", ".goosehints"]`) can add `CLAUDE.md` ([goosehints](https://goose-docs.ai/docs/guides/context-engineering/using-goosehints)) |

### aider

aider calls models through the LiteLLM Python SDK, so every LiteLLM model name works: `openrouter/anthropic/claude-sonnet-4` with `OPENROUTER_API_KEY`, or `openai/<model>` with `OPENAI_API_BASE` and `OPENAI_API_KEY` against a LiteLLM proxy or OpenRouter's `/api/v1` ([other LLMs](https://aider.chat/docs/llms/other.html), [OpenRouter](https://aider.chat/docs/llms/openrouter.html), [OpenAI compatible](https://aider.chat/docs/llms/openai-compat.html)).

| Concern | aider |
|---------|-------|
| Headless, resume | `aider --message "<text>"` or `--message-file`, `--yes-always`, `--no-auto-commits`, `--dry-run`; `--restore-chat-history` reloads the chat history file, not a session by ID ([scripting](https://aider.chat/docs/scripting.html), [options](https://aider.chat/docs/config/options.html)) |
| MCP | Not in aider's documentation. Third-party pages describe an `mcp-server` list in `.aider.conf.yml` with `command`, `args`, `env`; **to check** against the current release |
| Sub-agents, own model | None. `--weak-model` and `--editor-model` give one session a second and third model for commit messages and edits, not sub-agents |
| Tool-refusal hook | None. aider's tools are its own edit formats and shell suggestions, not a tool loop |
| Permission handler | `--yes-always` or the interactive confirmation; nothing an external process answers |
| Usage stream | `/tokens` reports the context's size; `--llm-history-file` logs the raw exchange; the per-message token and cost line aider prints is **to check**; no structured event stream ([commands](https://aider.chat/docs/usage/commands.html)) |
| Instruction files | None by default; `--read CONVENTIONS.md` or `read:` in `.aider.conf.yml` ([conventions](https://aider.chat/docs/usage/conventions.html)) |

### An agent loop of flai's own

A Go loop in flai calling the gateway's OpenAI-compatible (`/chat/completions` or `/v1/responses`) or Anthropic-compatible (`/v1/messages`) API, with flai's MCP tools and a few file and shell tools, would make every concern flai's to decide, and flai's to build:

- An MCP client over stdio and streamable HTTP, listing tools and calling them, with `flai mcp` as one server among the project's.
- The tool loop: tool definitions in each API's shape, streamed responses parsed for `tool_use` or `tool_calls`, results appended, a turn budget, and context compaction when the window fills.
- A session store, so a run can be resumed by ID after a thread is answered, and a log in a shape `flai/internal/usage` can read.
- A guard and permission handler as function calls, not hooks: the loop asks the policy before each tool, and asks `permission_prompt` on a thread when the policy says ask.
- Sub-agents as nested loops with their own model and tool allowlist, tagged in the log so usage is apportioned.
- Pricing: `usage` from the gateway, plus `cost` from OpenRouter's body or LiteLLM's `x-litellm-response-cost` header, so no price table of its own.
- Instruction loading, `CLAUDE.md` or `AGENTS.md` and `flai prime`, and the system prompt that Claude Code supplies today.

What it would not get: Claude Code's prompt caching defaults, system prompt, built-in tools, and the years of harness behaviour the conventions assume.

### Harnesses against the concerns

| Harness | Start/resume headless | MCP (stdio) | Sub-agents with own model | Tool-refusal hook | Permission handler | Usage stream | Instruction file |
|---------|-----------------------|-------------|---------------------------|-------------------|--------------------|--------------|------------------|
| Claude Code over a gateway | yes, `--session-id`, `--resume` | yes, `--mcp-config` | yes, `.claude/agents`, `--agents` | yes, `PreToolUse` with `agent_id` | yes, `--permission-prompt-tool` | yes, stream-json with cost (estimate) | `CLAUDE.md` |
| Codex CLI | yes, `exec`, `exec resume` | yes, `[mcp_servers]` | yes, `.codex/agents/*.toml` | yes, `PreToolUse`, needs hook trust | partial, `approval_policy`; no external approver | partial, JSONL tokens, no cost | `AGENTS.md` |
| OpenCode | yes, `run --session` | yes, `mcp.<name>` local | yes, `.opencode/agents` | partial, plugin `tool.execute.before` throws | partial, `permission` map, `--auto`; no external approver | yes, JSONL tokens and cost (**to check**) | `AGENTS.md`, falls back to `CLAUDE.md` |
| Goose | yes, `run -n`, `-r` | yes, `--with-extension` | yes, `GOOSE_SUBAGENT_*`, recipes | yes, `PreToolUse` hook; no sub-agent ID | partial, `GOOSE_MODE`; headless `approve` **to check** | partial, json/stream-json tokens, no cost | `AGENTS.md`, `.goosehints` |
| aider | partial, `--message`; no session ID | no (**to check**) | no | no | no | no | `--read` |
| flai's own loop | to build | to build | to build | to build (a function, not a hook) | to build | to build, cost from the gateway | to build |

What the documentation does not settle:

- Whether OpenRouter's or LiteLLM's `/v1/messages` carries every beta header and body field Claude Code sends, which the compatibility guide says a gateway must forward byte for byte; LiteLLM's unified endpoint translates, so some fields may be dropped. **To check** by running one `claude -p` through each.
- Whether Codex's `PreToolUse` `ask` and OpenCode's `permission.asked` can be answered by an outside process in a headless run, as `permission_prompt` is answered on a thread today.
- Whether Codex's and Goose's hooks tell a sub-agent's tool call from the main agent's, as Claude Code's `agent_id` does; Codex has `SubagentStart` and `SubagentStop`, Goose has neither.
- The exact fields of OpenCode's `step_finish` and Goose's `stream-json` events, read only from third-party write-ups.
- Which `wire_api` Codex should use against each gateway, and whether OpenCode's Anthropic provider accepts a custom base URL.
- Whether aider's current release is an MCP client; its own documentation does not say.
- What any harness reports as cost over a gateway: Claude Code and OpenCode price from their own tables, Codex and Goose report tokens only. Only the gateway knows what was charged: OpenRouter's `usage.cost` and `/api/v1/generation`, LiteLLM's `x-litellm-response-cost` and `/spend/logs`.

## Abstractions

What flai needs from any agent, stated without Anthropic's or OpenAI's API in it. Each subsection gives the contract, how each harness in [LiteLLM and OpenRouter](#litellm-and-openrouter) could meet it, and what flai must refuse or degrade when one cannot. Nothing here is built; the shapes are proposals for the epic.

### Harness, provider, and model

A story's `agent` runs three things together today: `harness`, the program that runs the agent loop with its tools and sub-agents; the provider its model calls go to, which is nowhere in flai, since `claude` sends them to Anthropic unless the environment `flai serve` happens to run in says otherwise; and `model`, a name flai passes through unread (`flai/internal/manifest/agent.go`). Over a gateway the three come apart: the same harness, a different provider, and a model name that means something only on that provider (`anthropic/claude-opus-5-5` on LiteLLM, `~anthropic/claude-opus-latest` on OpenRouter).

| Thing | What it is | Today | Proposed |
|-------|------------|-------|----------|
| Harness | The program running the loop, its tools, its sub-agents | `agent.harness`, a name `harness.For` knows | Unchanged; more adapters |
| Provider | Where the model calls go: Anthropic direct, a LiteLLM proxy, OpenRouter | The harness's own default | `agent.provider`, the name of an entry in a `providers` map; absent means the harness's own default |
| Model | The model's name on that provider | `agent.model`, pattern-checked only | Unchanged; its meaning is the provider's |

The provider entry holds what a harness needs to reach the gateway and nothing secret:

```yaml
providers:
  litellm:
    api: anthropic-messages      # anthropic-messages | openai-chat | openai-responses
    base_url: http://127.0.0.1:4000
    key_env: LITELLM_VIRTUAL_KEY # the NAME of the variable holding the key; never the key
    models:                      # optional: the harness's aliases on this provider
      haiku: anthropic/claude-haiku-4-5
      sonnet: anthropic/claude-sonnet-5
      opus: anthropic/claude-opus-5-5
```

The key stays in the operator's environment under the name `key_env` gives. `flai serve` reads that variable when it starts the agent and sets the harness's own variable in the child's environment, logging neither; the harness's `Start` record keeps the name, not the value. `flai serve` refuses a start whose provider names a variable that is unset, as it refuses an unknown harness.

Where each piece belongs follows how `Host` works today: the story names, the host says what the name is ([ADR-0038](../adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md)).

| Level | Holds | Why |
|-------|-------|-----|
| Host, `~/.flai/config.json` | `agent.harnesses.<name>.program` and `.args` (today); `agent.providers.<name>` with `api`, `base_url`, `key_env`, `models` (new), set with `flai serve agent provider <name>` and `settings.provider` | What runs on this host and where its network goes are the operator's; whoever edits a story in the dashboard must not choose them |
| Manifest, `system-flow.yaml` | `agent.harness`, `agent.provider`, `agent.model`, `agent.config`, `agent.roles`; the same under `planning`, `orchestration`, `analysis` | The project's default, copied into stories ([ADR-0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)) |
| Story front matter | The same fields | The story's own; changing it restarts a ready story's agent |

An alternative is a `providers` map in the manifest, committed, so a clone carries it. Base URLs are host-bound (`127.0.0.1:4000`) and so is the key's variable name, which is why the host is recommended; a manifest entry would be a default the host overrides, as `agent.harnesses` has no manifest half today. **To decide** by the operator.

What each adapter derives from the provider:

| Harness | Derived from `providers.<name>` |
|---------|---------------------------------|
| `claude-code` | Requires `api: anthropic-messages`. `ANTHROPIC_BASE_URL=<base_url>`, `ANTHROPIC_AUTH_TOKEN=$<key_env>`, `ANTHROPIC_API_KEY=` set empty, `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL`, `ANTHROPIC_DEFAULT_HAIKU_MODEL`, and `CLAUDE_CODE_SUBAGENT_MODEL` from `models`; `--model <model>`; `--fallback-model` must be a provider name. With no provider, nothing: today's behaviour |
| Codex | Requires `api: openai-chat` or `openai-responses`. `-c model_providers.<name>.base_url=<base_url> -c model_providers.<name>.env_key=<key_env> -c model_providers.<name>.wire_api=chat\|responses -c model_provider=<name> -c model=<model>`; Codex reads the key itself from `env_key`, so flai sets nothing secret |
| OpenCode | A generated configuration file named by `OPENCODE_CONFIG`: `provider.<name>.npm` from `api` (`@ai-sdk/openai-compatible`, `@ai-sdk/openai`, or `@ai-sdk/anthropic`, the last **to check**), `options.baseURL`, `options.apiKey: "{env:<key_env>}"`; `--model <name>/<model>` |
| Goose | `GOOSE_PROVIDER` from `api` and `base_url` (`litellm` with `LITELLM_HOST`, `openrouter`, `openai` with `OPENAI_HOST`, `anthropic` with `ANTHROPIC_HOST`), the provider's key variable set from `$<key_env>`, `GOOSE_MODEL=<model>` |
| `command` | `FLAI_PROVIDER`, `FLAI_PROVIDER_API`, `FLAI_PROVIDER_BASE_URL`, `FLAI_PROVIDER_KEY_ENV`, and `{provider}` in its arguments; the command reads the key itself |
| flai's own loop | Calls `base_url` in the `api` shape with the key from `$<key_env>` |

### Starting and resuming an agent

The contract is today's `Adapter` with two additions: the provider, and a statement of what the harness can do, so `flai serve` refuses or degrades before a story is started rather than after it fails.

```go
type Adapter interface {
    DefaultHost() Host
    Capabilities() Capabilities      // what flai serve may count on
    Start(req Request, host Host, prov *Provider) (Start, error)
    Reader() usage.Reader            // reads this harness's log; see Usage and cost
}
type Capabilities struct{ Resume, Guard, Permission, RoleModels, StrategicAgents, Cost bool }
```

| Harness | Start | Resume after a thread is answered |
|---------|-------|-----------------------------------|
| Claude Code | `claude -p`, `--session-id <uuid>` flai makes | `--resume <uuid>` (today) |
| Codex | `codex exec "<prompt>"` | `codex exec resume <id>`; the ID comes from the `thread.started` event, so flai reads it from the log rather than choosing it |
| OpenCode | `opencode run --session <id>` | `--session <id>` again, or `--continue` |
| Goose | `goose run -n <name>` | `-r -n <name>` |
| aider | `aider --message` | None |
| flai's loop | Its own session store | Its own |

What a harness without resume loses: `flai serve` starts a fresh session with the answered thread named in the prompt, as `FLAI_ANSWERED` tells the `command` harness today. The agent primes again and has lost its context; the work goes on at a higher cost. Degrade, do not refuse.

### Refusing a sub-agent's write

Today `flai guard` is a hook whose input is Claude Code's JSON ([ADR-0060](../adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)). The contract is the decision under it: before each tool call, given who calls and what, allow or refuse with a reason.

```go
type Call struct {
    Session  string
    Subagent bool              // a sub-agent's call, not the session's own agent's
    Role     string            // FLAI_ROLE: plan, orchestrate, analyze, or empty
    Kind     CallKind          // Shell, FileWrite, FlaiTool
    Command  string            // Shell: the line, split and judged as today
    Path     string            // FileWrite
    Tool     string            // FlaiTool: the tool's name without the harness's prefix
    Args     map[string]any
}
func Decide(c Call, running SubagentRecord) Verdict
```

`flai guard` keeps one `Decide` and gains one input reader per harness, each mapping the harness's tool names to `Kind` and its sub-agent sign to `Subagent`.

| Harness | Hook | Sub-agent sign | Registration |
|---------|------|----------------|--------------|
| Claude Code | `PreToolUse`, exit 2 | `agent_id` | `.claude/settings.json` (today) |
| Codex | `PreToolUse`, exit 2 or `permissionDecision: deny` | None in the input; a record from `SubagentStart` and `SubagentStop`, as ADR-0092 keeps, **to check** | `.codex/hooks.json`; needs hook trust, `--dangerously-bypass-hook-trust` or a managed hook |
| OpenCode | A plugin's `tool.execute.before` throws | **To check** | `.opencode/plugins/flai-guard.js`, a few lines that exec `flai guard` |
| Goose | `PreToolUse` in `hooks.json`, exit 2 or `decision: block` | None; every call looks like the agent's own | `.agents/plugins/flai/hooks/hooks.json` |
| aider | None | | |
| flai's loop | `Decide` called before each tool | The loop knows | None |

A second line the hook does not need: `flai mcp` could refuse a sub-agent itself if it knew the caller. It does not: a sub-agent shares its parent's MCP connection in every harness above. Shell and file writes never reach flai at all. The hook, or the loop, is the only place.

What a harness without a hook loses: flai cannot tell a sub-agent's write from the agent's. The story's agent may run; its roles may not. Refuse a role on such a harness, as a role on another harness is refused today, unless the operator sets `agent.harnesses.<name>.guard: none` on the host and takes the risk. Goose, until its hooks carry a sub-agent sign, is such a harness.

### Asking the owner before a protected write

Today Claude Code calls `permission_prompt` as its permission handler ([ADR-0086](../adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md)). The contract: when the policy says ask, hold the call, put the question to the story's owner or the project's owner on a thread ([ADR-0097](../adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md)), and answer allow or deny within a bound ([ADR-0124](../adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md)); operations `Ask(call) (allow|deny, reason)` and the bound.

| Harness | Way to ask | Bound |
|---------|------------|-------|
| Claude Code | `--permission-prompt-tool` (today), or the `PreToolUse` hook holding the call | The MCP idle timeout, four minutes; a hook's own timeout |
| Codex | No external approver. A `PreToolUse` hook that holds the call and returns `allow` or `deny` itself: `flai guard --ask` | The hook timeout, **to check** |
| OpenCode | A plugin's `tool.execute.before` awaiting the answer | **To check** |
| Goose | A `PreToolUse` hook that holds | **To check** |
| flai's loop | `Ask` called by the loop | flai's own |

The common form is the guard holding the call and asking, with the MCP tool as Claude Code's variant of it. The four-minute bound becomes the harness's. The protected list stays flai's policy (`flai/internal/protected/protected.go`) and grows one entry per harness: the files that change what the harness's agent may do, `.codex/`, `.opencode/`, `.agents/plugins/`, `AGENTS.md`, beside `.claude/` and `.mcp.json`.

What a harness without a hook loses: a protected write is denied outright by the harness's own permission map, or allowed outright under its `--auto` or `danger-full-access`. Refuse to start a story's agent on such a harness unless `auto-approve` is on for the project or the host's arguments deny those paths in the harness's own vocabulary; say which in the refusal.

### A model per sub-agent role

The contract is [ADR-0065](../adrs/0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md)'s: `roles.<role>.model` laid over the role's definition, on the story's provider. A role's model is a name on that provider; the provider's `models` aliases resolve `haiku` and `sonnet`.

| Harness | How the role's model reaches it | Per role |
|---------|-------------------------------|----------|
| Claude Code | `--agents` JSON with `model` (today) | Yes |
| Codex | A TOML definition with `model`, written to `.codex/agents/<role>.toml` in the worktree, or a `CODEX_HOME` flai prepares; `-c agents.default_subagent_model` for the rest | Yes, **to check** |
| OpenCode | `agent.<role>.model` in the generated configuration | Yes |
| Goose | `GOOSE_SUBAGENT_MODEL`, one for every sub-agent; a recipe per role for more | One model for all roles |
| flai's loop | A nested loop with its own model | Yes |

What a harness without it loses: the role's model is ignored. Refuse the start when a story's role names a model the harness cannot honour, as today, rather than run the story on a model the story did not name.

### Definitions and the instructions file

A role definition is a name, a description, a prompt, a tool allowlist, and a model. Today the five live in `.claude/agents/` in Claude Code's format and vocabulary. Options:

| Option | Source | Cost |
|--------|--------|------|
| 1 | Keep `.claude/agents/*.md` as the source; each adapter translates at start into its own format (`--agents` JSON, `.codex/agents/*.toml`, OpenCode `agent.<name>`), with a tool-name table per harness (`Read`, `Grep`, `Glob`, `Bash`, `mcp__flai__*` to each harness's names) | One table per harness; the files stay where Claude Code reads them unchanged |
| 2 | Move the source to a neutral folder in the template and generate each harness's files with `flai upgrade` | A new folder in every project; `.claude/agents/` becomes generated and must not be hand-edited |
| 3 | A hand-kept copy per harness | Three copies that drift |

The instructions file: Codex, OpenCode, and Goose read `AGENTS.md`, and Codex (`project_doc_fallback_filenames`) and Goose (`CONTEXT_FILE_NAMES`) can be told to read `CLAUDE.md`; OpenCode falls back to it on its own. Keep `CLAUDE.md` the map and point each harness at it from the adapter, as [agent-context.md](agent-context.md) already leans; add `AGENTS.md` only if a harness that cannot be pointed is adopted.

What a harness without definitions loses: no roles, and the strategic agents run from flai's `Prompt` alone, which already carries the role (`FLAI_ROLE`), without the definition's tool allowlist. Refuse roles; let a strategic agent run only where the guard holds it to its role.

### Usage and cost

The contract is a reader per harness producing neutral events the measurement in `flai/internal/usage` already works from, and a source of price.

```go
type Reader interface{ Read(line []byte) (Event, bool) }
type Event struct {
    Kind       EventKind // SessionStart, Call, ToolResult, End
    Session    string
    Model      string
    Tokens     Tokens    // In, Out, CacheRead, CacheWrite
    CostUSD    *float64  // nil when the harness reports none
    ParentCall string    // the call that started this sub-agent; empty for the agent's own
    Tool, Description string
}
```

| Harness | Tokens per call | Per model | Cost | Sub-agent attribution |
|---------|-----------------|-----------|------|-----------------------|
| Claude Code | yes | yes | its own estimate from Anthropic's prices | `parent_tool_use_id` |
| Codex | `turn.completed` | no | none | **To check** |
| OpenCode | `step_finish` | **to check** | its own table, **to check** | **To check** |
| Goose | session store, `stream-json` **to check** | no | none | none |
| flai's loop | yes | yes | the gateway's | the loop's own |

Where a price comes from when the log gives none, or gives an estimate from the wrong price list:

| Source | How | Fits |
|--------|-----|------|
| The harness | `total_cost_usd`, OpenCode's `cost` | Today's `result`; wrong over a gateway |
| The gateway, per call | LiteLLM `x-litellm-response-cost`; OpenRouter `usage.cost` | Only flai's own loop sees the response |
| The gateway, after the fact | LiteLLM `/spend/logs` by virtual key; OpenRouter `/api/v1/key` and `/api/v1/generation` | A key per project, or per agent, joins spend to a story: a `spend` reader per provider `api`; per-call joining needs the gateway to keep a session tag, LiteLLM's reading of `x-claude-code-session-id` **to check** |
| A price table | LiteLLM's public model cost map, or flai's own | ADR-0051 chose measured over tabulated; a table dates |
| Estimated | `Rates` from reported totals (today) | Nothing to rate when no run ever reported a cost |

Recommended: tokens and attribution from the harness's reader; cost from the gateway's spend log when the story has a provider, from the harness when it has none, estimated otherwise, with the source named on the item's `usage` (`priced_by: harness | gateway | estimate`). That changes `design/system/metrics.md` and needs an ADR. What a harness that logs tokens only loses: cost stays estimated, marked so as today, and nothing when no rate exists.

### Settings each option adds

| Option | Host | Manifest and story | Also |
|--------|------|--------------------|------|
| All | `agent.providers.<name>.api`, `.base_url`, `.key_env`, `.models.<alias>`; `agent.harnesses.<name>.guard` | `agent.provider`, and under `planning`, `orchestration`, `analysis` | `settings.provider` in the host API; the Settings page |
| A, Claude Code over a gateway | None more | `fallback_model` and role models as provider names | |
| B, Codex | `agent.harnesses.codex.program` (default `codex`), `.args` (default sandbox and approval policy) | `config` keys the adapter takes, such as `reasoning_effort` | `.codex/hooks.json` in the template |
| C, OpenCode | `agent.harnesses.opencode.program`, `.args` | `config` keys | `.opencode/plugins/flai-guard.js` in the template |
| D, Goose | `agent.harnesses.goose.program`, `.args` | `config` keys | `.agents/plugins/flai/hooks/hooks.json` in the template |
| E, flai's loop | Tool and turn limits | `config` keys | New dependencies in `design/tech/` |

No key is stored anywhere: `key_env` is a name. Documentation each changes: `docs/operators/settings.md`, one row per key, and its test `TestSettingsIndexIsComplete` in `flai/cmd/settings_doc_test.go`; `docs/operators/index.md` under authentication, naming the key variable; `docs/users/flai.md`, "Starting an agent", "Sub-agents", and "Writes to paths Claude Code protects", the last renamed; `docs/users/flai-reference.md` through `make flai-reference`; `design/system/flai-cli.md`.

### Options compared

| Option | Met now | Partial | To build | Effort | Unlocks | Risk | ADRs |
|--------|---------|---------|----------|--------|---------|------|------|
| A, Claude Code over a gateway | Start, resume, guard, permission, roles, definitions | Cost (estimate, not the bill) | Provider split; gateway spend reader | Small | Self-hosted LiteLLM in front of Anthropic, Bedrock, Vertex; OpenRouter credits; one key per project | Beta headers dropped by LiteLLM's unified endpoint; Claude models only | Refines 0037, 0038, 0051 |
| B, Codex CLI | Start, resume, roles, MCP | Guard (sub-agent sign), cost (tokens only) | Adapter, JSONL reader, hooks.json, hold-and-ask in the guard, TOML definitions, hook trust | Medium | OpenAI models and anything OpenAI-shaped through either gateway, self-hosted included | Hook trust in headless runs; `ask` with no approver | Refines 0060, 0065, 0086, 0092, 0106, 0124 |
| C, OpenCode | Start, resume, roles, MCP, cost | Guard (plugin), permission | Adapter, JSONL reader, plugin, definitions | Medium | Any model either gateway serves | Plugin API stability; sub-agent sign unknown | As B |
| D, Goose | Start, resume, MCP | Roles (one model), guard (no sub-agent sign), cost (tokens) | Adapter, reader, hooks, provider mapping | Medium | Both gateways natively | No sub-agent sign: roles refused | As B, and 0071 |
| E, flai's loop | None | | Everything: MCP client, tool loop, session store, compaction, guard as a function, nested loops, pricing | Large | Any model, any API, exact gateway cost per call | Rebuilds a harness; loses Claude Code's prompt, caching, tools | Supersedes 0060 and 0086 for its own harness |

### Recommendation

Build in this order:

1. **A spike story first**: run `claude -p` with the adapter's own arguments through LiteLLM's unified `/v1/messages`, its `/anthropic` pass-through, and OpenRouter's `/api/v1/messages`, with `permission_prompt` and `flai guard` on, as `claudecheck.go` already does against Anthropic. It settles the beta-header question, what `total_cost_usd` says against the gateway's spend log, and whether the aliases resolve. One or two days; every later story depends on its answers.
2. **A, the provider split and Claude Code over a gateway.** The smallest change that reaches both gateways, and every concern stays met. It adds `providers` on the host, `agent.provider` on the manifest and story, the environment the adapter derives, and a `spend` reader per gateway API. Claude models only, which is what every project runs today.
3. **The neutral contracts, inside A**: `Capabilities`, `Reader`, `Call` and `Decide`, the guard's hold-and-ask, and the protected list per harness. They cost little while one adapter exists and make the next one an adapter, not a rewrite.
4. **B, Codex CLI.** The harness nearest Claude Code in shape: headless exec and resume, `PreToolUse` with `SubagentStart` and `SubagentStop`, MCP in configuration, sub-agent definitions with a model. It unlocks every OpenAI-shaped model either gateway serves, self-hosted ones included. Its open items, hook trust and the sub-agent sign, are checked in its first story.
5. **C, OpenCode**, if a second non-Anthropic harness is wanted; it reports cost and reads `CLAUDE.md`, but its guard is a plugin and its permission handling is less known.
6. **Not now**: D, Goose, until its hooks can tell a sub-agent's call, since without that no role may run; E, flai's own loop, which rebuilds what every harness gives and is justified only if B and C cannot hold the guard.

What the recommendation does to the current ADRs:

| ADR | Fate | Why |
|-----|------|-----|
| [0037](../adrs/0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md) | Refined | `agent` gains `provider`; the copy and merge rules stand |
| [0038](../adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) | Refined | The host gains `providers` beside `harnesses`; adapters declare capabilities; what runs stays the operator's |
| [0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md) | Refined | Cost may be read from the gateway's spend log; `usage` names its price source; still measured, never tabulated |
| [0059](../adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) | Untouched | The roles and their work are the same on any harness |
| [0060](../adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) | Refined by A and B; superseded by E for its harness | The policy stands; the input shapes and registrations multiply; a loop calls `Decide` with no hook |
| [0065](../adrs/0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md) | Refined | Each adapter lays the role's model over the definition in its own format; a role's model is a provider name |
| [0071](../adrs/0071-a-task-s-usage-is-the-calls-of-the-sub-agents-started-for-it-and-an-even-share.md) | Untouched by A and B; refined by D | Attribution needs a parent call and the task's ID in the sub-agent's start; Goose has neither |
| [0086](../adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) | Refined | The MCP permission tool becomes Claude Code's form of the guard's hold-and-ask |
| [0092](../adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md) | Refined by B | The running-sub-agent record comes from each harness's start and stop events |
| [0097](../adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) | Untouched | Who answers does not depend on the harness |
| [0102](../adrs/0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md) | Untouched | The rule applies through whichever guard runs |
| [0105](../adrs/0105-a-story-s-empty-wakes-the-wait-for-events-calls-of-its-agents-that-timed-out.md) | Untouched | An empty wake is a tool name and a result, both in the neutral event |
| [0106](../adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md) | Refined | The protected list becomes flai's, with each harness's own files in it |
| [0116](../adrs/0116-when-flai-measures-a-story-s-usage-from-its-logs-it-classifies-each-turn-of-the.md) | Untouched | Turn classes come from tool names and shell commands the reader maps |
| [0124](../adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md) | Refined | The bound is Claude Code's MCP idle timeout; each harness has its own |
| [0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md), [0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md), [0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md) | Untouched by A; refined by B and C | A strategic agent is started from its definition with `--agent`; Codex's and OpenCode's equivalents are **to check** |
| [0108](../adrs/0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md), [0064](../adrs/0064-a-story-in-ready-or-in-progress-with-no-agent-run-on-this-host-is-started-here.md) | Untouched | Restarts and begun-elsewhere starts are about outcomes, not harnesses |

Questions only the operator can answer, each with the recommended answer first:

1. **The split.** Add `provider` beside `harness` and `model`, with the `providers` map on the host; or put `providers` in the manifest with a host override. Recommended: the host, as `harnesses` is.
2. **Which adapters, in what order.** Recommended: the spike, then A with the neutral contracts, then B; C on demand; D and E not now.
3. **A harness without a guard hook or a permission handler.** May it run a story's agent, and with what limits? Recommended: the story's agent may run, its roles may not, unless `agent.harnesses.<name>.guard: none` is set on the host; a protected write on such a harness is refused unless `auto-approve` is on.
4. **Cost when a harness logs none.** The gateway's spend log, a price table, or an estimate. Recommended: the gateway's spend log per provider, keyed by a virtual key per project, named on `usage` as its source; a price table is not added.
5. **A paid trial before the epic.** The spike needs a LiteLLM proxy, free, and an OpenRouter key with a few dollars of credit. Recommended: yes, both, in the spike story.

## Decision

The five questions at the end of [Recommendation](#recommendation) were put to the operator on TH-0368 on 2026-10-08, recommendation first. The first four are drafted as proposed ADRs, to be accepted as the operator answers; the fifth, the paid trial, is spend and the operator's alone, recorded in the epic.

| Question | ADR | Status |
|----------|-----|--------|
| 1, the harness, provider, and model split, with `providers` on the host | [ADR-0129](../adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md) | proposed |
| 2, the neutral contracts every harness meets flai through | [ADR-0130](../adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) | proposed |
| 3, a harness without a guard hook or a permission handler | [ADR-0131](../adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md) | proposed |
| 4, where a story's cost comes from, `priced_by` on `usage` | [ADR-0132](../adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md) | proposed |
| 5, a paid trial in the spike story | none; the epic's notes | awaiting the operator |

The order the epic builds in, question 2, is the recommendation's: the spike, then Claude Code over a gateway with the provider split and the neutral contracts, then Codex CLI; OpenCode on demand; Goose and a loop of flai's own not now.
