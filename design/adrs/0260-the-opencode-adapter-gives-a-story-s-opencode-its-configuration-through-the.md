---
id: ADR-0260
title: "The opencode adapter gives a story's OpenCode its configuration through the environment, flai's MCP entry in OPENCODE_CONFIG_CONTENT and the deny rules in OPENCODE_PERMISSION, so flai writes no file into the worktree and OpenCode merges the project's own opencode.json under them"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0130, ADR-0131]
---

# ADR-0260 The opencode adapter gives a story's OpenCode its configuration through the environment, flai's MCP entry in OPENCODE_CONFIG_CONTENT and the deny rules in OPENCODE_PERMISSION, so flai writes no file into the worktree and OpenCode merges the project's own opencode.json under them

## Context

A story's agent on OpenCode needs two things from flai that OpenCode reads from its configuration: a local MCP entry that runs `flai mcp --agent <name>`, and, when the host says `deny_protected` for the harness ([ADR-0131](0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md)), `deny` rules in its `permission` map for every path the protected list holds for `opencode` ([ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)). The project may have an `opencode.json` of its own, committed, with keys flai must keep.

The worktree must stay clean: a written `opencode.json` or a file under `.opencode/` left uncommitted refuses the story's move to review, and both are files the harness reads to decide what its agent may do, protected for that reason. S-0406's planning weighed two ways to keep it clean: `OPENCODE_CONFIG` naming a file under `.flai-cache`, or a git-ignored path in the worktree.

[agent-adapters.md § OpenCode](../system/agent-adapters.md#read-from-the-source-at-v1195-s-0406), read from OpenCode's source at v1.19.5, settles how OpenCode loads its configuration. Each source is deep-merged over the one before: the global files, `OPENCODE_CONFIG`, every project `opencode.json(c)` up to the worktree root, the config folders such as `.opencode/`, then `OPENCODE_CONFIG_CONTENT` (inline JSONC), the console organisation's and managed config, and last `OPENCODE_PERMISSION`, a JSON object merged into `permission`. Objects such as `mcp` and `permission` merge key by key; a string or an array is replaced whole.

## Decision

The opencode adapter gives a story's OpenCode its configuration through the environment, flai's MCP entry in `OPENCODE_CONFIG_CONTENT` and the deny rules in `OPENCODE_PERMISSION`, so flai writes no file into the worktree and OpenCode merges the project's own `opencode.json` under them.

1. `OPENCODE_CONFIG_CONTENT` holds one key, `mcp.flai`: `type: "local"`, `command` flai's executable with `mcp --agent <name>`, `enabled: true`, and `timeout: 1800000`, so that flai's tools that hold a call, `wait_for_events` and `permission_prompt`, run past OpenCode's default, as [ADR-0258](0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md) does for Claude Code. It names no provider, no model, and no key.
2. Under `deny_protected`, `OPENCODE_PERMISSION` holds an `edit` object with a `deny` rule for each path of `protected.For("opencode")`: a folder as `<name>/*` and `*/<name>/*`, a file as `<name>` and `*/<name>`, and a root path as itself. Without `deny_protected` the variable is not set.
3. Because a later object replaces a string whole, a project whose own `opencode.json` or `opencode.jsonc` at the worktree's root gives `permission`, or `permission.edit`, as one action keeps it: the adapter reads that action and writes it as the first rule, `"*"`, of the object it gives, so that what flai does not deny stays as the project said.
4. Every other key of the project's configuration is OpenCode's to merge, and flai reads nothing else of it.
5. The adapter sets the two variables in the start's environment; the start is logged and recorded whole, since neither carries a secret.

## Consequences

- The story's worktree gains no file from flai, and a project's `opencode.json` is never rewritten. A project's own `mcp.flai` entry is overridden by flai's, as the managed and organisation sources still override flai's.
- `edit` rules do not cover a write made through `bash`; until the guard plugin of a later E-0023 story, the operator's `deny_protected` is their word that the rest of their arguments cover it, as ADR-0131 says.
- OpenCode itself writes into a config folder it loads, a `.gitignore` and a plugin install, so a project that commits an `.opencode/` folder may see those in its worktree. That is OpenCode's, not flai's, and the guard story revisits it.
- The guard plugin, being code, cannot go in an environment variable; the later story that adds it chooses between `plugin` by absolute path in `OPENCODE_CONFIG_CONTENT` and an `OPENCODE_CONFIG_DIR` flai owns outside the worktree.

## Alternatives considered

- `OPENCODE_CONFIG` naming a file under the worktree's `.flai-cache`: the file must be written before each start, and it loads before the project's files, so a project's own `mcp.flai` or `permission` would override flai's.
- A git-ignored `opencode.json` in the worktree: a project with a committed `opencode.json` has nowhere to put it, `opencode.json` is protected, and an ignore rule that hides it hides the project's own.
- `OPENCODE_DISABLE_PROJECT_CONFIG` with a merged copy flai writes: flai would parse and merge OpenCode's JSONC itself, and drift from how OpenCode merges.
