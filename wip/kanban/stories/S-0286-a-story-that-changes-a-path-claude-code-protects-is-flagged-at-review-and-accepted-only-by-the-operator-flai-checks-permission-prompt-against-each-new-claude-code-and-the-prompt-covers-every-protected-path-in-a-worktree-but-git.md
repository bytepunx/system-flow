---
id: S-0286
type: story
nature: improvement
title: A story that changes a path Claude Code protects is flagged at review and accepted only by the operator, flai checks permission_prompt against each new Claude Code, and the prompt covers every protected path in a worktree but .git
status: backlog
owner: alex
created: 2026-10-06T06:32:43Z
updated: 2026-10-06T06:32:52Z
transitions: []
tags: [flai, flaiover]
topics: [cli]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/preview/accept.go, flai/cmd/accept.go, flai/internal/serve, flai/internal/harness, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/adrs, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators]
after: [S-0283]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0286 A story that changes a path Claude Code protects is flagged at review and accepted only by the operator, flai checks permission_prompt against each new Claude Code, and the prompt covers every protected path in a worktree but .git

## Goal

An agent's write under `.claude/` in its story's worktree goes through with no person in the middle of the story, and the operator still sees every such change before it takes effect. [ADR-0086](../../../design/adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) gave flai serve's agents `permission_prompt` and the `auto-approve` host action. This story adds what makes auto-approve safe to leave on, and what tells the operator when Claude Code stops accepting the prompt's answer.

The reasoning: every refused write so far was in a story's worktree, in its `.claude/` or its `template/root/.claude/`. flai serve starts the agent in the main checkout, and Claude Code reads a session's hooks, settings, and agent definitions from there. A write in the worktree changes nothing the running agent obeys until the story is accepted. So acceptance is where a person checks it, and nothing but a person may accept it.

Three parts:

1. **Acceptance is the gate.** The review of a story lists the files its branch changes on a path Claude Code protects. Such a story is accepted by the operator only: an agent's acceptance of it is refused, whatever lets the agent accept others.
2. **flai checks the prompt against Claude Code.** A Claude Code update broke `permission_prompt`'s answer the day after it shipped (S-0283), and the first sign was a story that hung. When the Claude Code flai serve runs has a version flai has not checked, flai makes one real write through `permission_prompt` with it, in a scratch folder, and tells the operator when it fails.
3. **The prompt covers what Claude Code protects.** Today it approves only a path with a `.claude` folder. Claude Code protects more, `.mcp.json` among them, which the template also ships. The prompt handles every protected path inside an in-progress story's worktree except `.git`, by the same rules: auto-approve, or a thread.

## Acceptance criteria
- [ ] An ADR refining ADR-0086 records the three parts and why acceptance is the gate
- [ ] `flai accept`'s preview and the dashboard's review page list the files a story's branch changes on a protected path, and say that only the operator accepts it
- [ ] An agent's acceptance of such a story is refused with the files named: by an agent on its own name, and by the orchestrator under `accept_reviews` once S-0221 has that path; a test covers each path there is
- [ ] When flai serve's Claude Code has a version not yet checked, flai runs one write under a `.claude/` folder of a scratch project through `permission_prompt` with it, records the version and the outcome, and on a failure opens a thread to the operator that quotes Claude Code's error; a version already checked is not checked again; a test runs it with a stand-in for `claude`
- [ ] `permission_prompt` handles an Edit, Write, MultiEdit, or NotebookEdit of any path Claude Code protects inside an in-progress story's worktree, and still refuses `.git`, the main checkout, and every path outside the worktree, with tests
- [ ] The design, the users' and operators' guides, and `flai serve actions` say what auto-approve allows, what the review shows, and what the check does

## Tasks

## Notes

The operator's decisions, alex, 2026-10-06, in conversation with Claude:

- The check on a `.claude/` write moves from the middle of the story to acceptance. Auto-approve is to be on for this project.
- A story that changes a protected path is never accepted by an agent.
- flai checks `permission_prompt` against a new Claude Code version and tells the operator when it fails.
- The prompt's scope widens to Claude Code's other protected paths in the worktree, never `.git`.

Not chosen, and why:

- `--permission-mode bypassPermissions` lifts the protection everywhere, the main checkout's settings and `.git` included, with no review.
- Moving the template's files out of `.claude/` changes the template's layout and still needs this repository's own copy generated. Take it up only if the prompt breaks again.
- A flai tool that writes the files itself goes around Claude Code's protection instead of through `--permission-prompt-tool`.

Facts found while designing it:

- Claude Code's list of protected paths is in its documentation, "Protected paths" on the permission modes page. It says allow rules in settings do not pre-approve such a write, and that `acceptEdits` prompts for one, which headless goes to the permission prompt tool.
- Against Claude Code 2.1.290, an allow answered as one text block lets a Write under `.claude/agents/` through in a headless `acceptEdits` session. The run cost about 0.02 USD on Haiku, which is what one check of part 2 costs.
- That the running agent reads its configuration from the main checkout rests on Claude Code's documentation and on the agent's working directory, which is the main checkout. It was not tested with a hook changed in a worktree; test it before the ADR relies on it.
- flai has no `auto-accept` code today: the board's text names the tag and nothing enforces it.

Waits for S-0283: the check of part 2 needs an answer Claude Code accepts, and both change `permission.go`. S-0284 changes who answers a thread, in the same file; this story leaves that to it.
