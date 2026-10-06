---
id: S-0286
type: story
nature: improvement
title: A story that changes a path Claude Code protects is flagged at review and accepted only by the operator, flai checks permission_prompt against each new Claude Code, and the prompt covers every protected path in a worktree but .git
status: backlog
owner: alex
created: 2026-10-06T06:32:43Z
updated: 2026-10-06T22:56:13Z
transitions: []
tags: [flai, flaiover]
topics: [cli, dashboard]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/preview/accept.go, flai/cmd/accept.go, flai/internal/serve, flai/internal/harness, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/adrs, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/preview/orchestrator.go, flai/cmd/accept_protected_test.go, flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/serve.go, flai/internal/harness/adapters.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [S-0283]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: planner-S-0286
    at: 2026-10-06T22:54:26Z
  value: 150
  by: planner-S-0286
  at: 2026-10-06T22:55:05Z
forecast:
  duration: 1h
  delivery: 2026-10-07T10:51:00Z
  basis: "Its own forecast of 1h; 34th in the pull order with an in-progress limit of 3, behind S-0299, S-0301, S-0300, S-0228, S-0261, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0265, S-0272, S-0273, S-0274, S-0275, S-0277, S-0279, S-0280 and S-0281."
  by: flai
  at: 2026-10-06T22:56:13Z
finalized:
  by: alex
  at: 2026-10-06T22:49:46Z
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
- T-1106 An ADR refining ADR-0086 records the acceptance gate for protected paths, the check of permission_prompt against each new Claude Code, and the prompt's wider scope
- T-1110 permission_prompt handles an Edit, Write, MultiEdit, or NotebookEdit of every path Claude Code protects inside an in-progress story's worktree, and still refuses .git
- T-1114 flai serve checks permission_prompt against each Claude Code version it has not checked, records the outcome, and opens a thread to the operator on a failure
- T-1120 flai accept's preview lists the files a story's branch changes on a protected path, and an agent's acceptance of such a story is refused with the files named
- T-1121 The dashboard's review page lists the files a story changes on a protected path and says that only the operator accepts it
- T-1122 The design, the users' and operators' guides, and flai serve actions say what auto-approve allows, what the review shows, and what the Claude Code check does

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

### Planning

Planned by planner-S-0286 on 2026-10-06. The plan thread on S-0286 lists the tasks, their layers, and the assumptions.

Layers:

1. T-1106, the ADR.
2. T-1110, the prompt's scope and the shared list of protected paths; T-1114, the Claude Code check.
3. T-1120, the acceptance gate, which reads T-1110's list.
4. T-1121, the review page; T-1122, the design, the guides, and `flai serve actions`.

Touches:

- **Declared, kept:** every touch the story carried.
- **Design:** the acceptance gate and the prompt's scope.
  - `flai/internal/protected/protected.go` and its test: a new leaf package holding the list both `permission_prompt` and the preview read, since `mcpserver` does not import `preview`.
  - `flai/internal/preview/orchestrator.go`: it holds ADR-0093's blocker codes, and the orchestrator's refusal is one more.
- **Layout:** the Claude Code check and the accept tests.
  - `flai/internal/serve/claudecheck.go`, its test, and `serve.go`: the check is new code in flai serve, started from its loop.
  - `flai/internal/harness/adapters.go`: it builds `claude`'s argv, which the check reuses.
  - `flai/cmd/accept_protected_test.go`: accept tests live in `flai/cmd` as `accept_*_test.go`.
- **Design and layout:** `flai/internal/hostapi/writes.go`. `ActionAutoApprove`'s description there is what `flai serve actions` prints, which the last criterion names.
- **Co-change:** `docs/users/flai-reference.md` (17% of the 499 commits `flai touches suggest` read), for `flai accept`'s help.
- **Design:** `docs/operators/settings.md`, the only operators' guide that names `auto-approve`.
- Not taken from `flai touches suggest`:
  - `design/issues/summary.md`, which the close-out writes itself.
  - `flai/internal/mcpserver/server.go` and `folder.go`. The tool's registration does not change.
  - `flai/internal/hostapi/writes_test.go`. `accept.run` builds its refusal from the same preview, so only the action's description in `writes.go` changes.
  - The workflow and work-hierarchy design, and `template/`. No convention or template file changes.

Folder touches kept, as declared:

- `design/adrs`: `flai adr new` assigns the ADR's number when it writes the file, so no task can name it yet.
- `flai/internal/serve`, `flai/internal/harness`, and `docs/operators`: the operator declared them. Each task names its files, and per ADR-0096 the story's claim narrows each folder to the files its tasks name. The story's agent may narrow them with `flai touches`.

Topic `dashboard` added: T-1121 changes the review page.

Forecast 1h, delivery 2026-10-07T10:48Z.

- `flai forecast` gave 43m: 83 s per unit of size over 25 done large-band improvement stories on claude-opus-5-5, times size 31 (6 criteria, 25 touches). It put the story 34th in the pull order, delivering at 10:31Z.
- Raised by 17m. T-1106 runs a live headless Claude Code test of the worktree hook claim before the ADR relies on it. T-1114 is new code: a scratch flai project, a stand-in `claude`, recorded state, and a thread. Size counts neither.
- Delivery is shifted by the same 17m.

Cost of delay 150 USD a week, as `flai cod` gives from the operator's input: 1h lost per 168h cycle at 150 USD an hour. alex confirmed the input on TH-0206. It stands as computed.
