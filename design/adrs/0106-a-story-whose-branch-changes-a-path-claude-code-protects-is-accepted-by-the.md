---
id: ADR-0106
title: "A story whose branch changes a path Claude Code protects is accepted by the operator only, flai checks permission_prompt against each new Claude Code, and the prompt covers every protected path in a story's worktree but .git"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0086]
---

# ADR-0106 A story whose branch changes a path Claude Code protects is accepted by the operator only, flai checks permission_prompt against each new Claude Code, and the prompt covers every protected path in a story's worktree but .git

## Context

[ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) gave the agents flai serve starts `permission_prompt` as Claude Code's permission prompt tool. It approves an Edit, Write, MultiEdit, or NotebookEdit of a file in a `.claude/` folder inside an in-progress story's worktree, on the owner's answer on a thread or at once under the `auto-approve` host action. The operator wants auto-approve on for this project, so that such a write needs no person in the middle of a story. Three things stand in the way of leaving it on.

1. **Nothing reviews the write.** With auto-approve on, a story's agent can change hooks, settings, and agent definitions with no one reading them, and an agent may accept a story: the orchestrator under `accept_reviews` ([ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)), or a story's agent running `flai accept` on its own name.
2. **A Claude Code update can break the prompt silently.** A Claude Code release changed the answer it accepts from the tool the day after `permission_prompt` shipped (S-0283, I-0082). The first sign was a story whose agent hung.
3. **The prompt covers less than Claude Code protects.** Claude Code asks for more than `.claude/`. The template ships `.mcp.json`, which is one of them, and a story that changes it is refused today.

Every refused write so far was in a story's worktree, in its `.claude/` or its `template/root/.claude/`. flai serve starts the agent in the main checkout, and Claude Code reads a session's settings and hooks from the folder it starts in. That was tested for this ADR against Claude Code 2.1.290 on 2026-10-07: a headless `claude -p` started in a scratch git repository ran Bash with a `PreToolUse` hook in `wt/.claude/settings.json`, a folder below it, and the hook did not run; the same hook in the repository's own `.claude/settings.json` did. So a write in a worktree changes nothing the running agent obeys until the story's branch is merged, which is when the story is accepted.

Claude Code's documentation lists the protected paths under "Protected paths" on its permission modes page. Writes to them are prompted in `default` and `acceptEdits` modes, headless through the permission prompt tool, and `permissions.allow` rules do not pre-approve them. As of Claude Code 2.1.290 they are:

| Kind | Paths |
|------|-------|
| Folders | `.git`, `.config/git`, `.vscode`, `.idea`, `.husky`, `.cargo`, `.devcontainer`, `.yarn`, `.mvn`, `.claude` (except a few of Claude Code's own, such as `.claude/worktrees/`) |
| Files | `.gitconfig`, `.gitmodules`; `.bashrc`, `.bash_profile`, `.bash_login`, `.bash_aliases`, `.bash_logout`, `.zshrc`, `.zprofile`, `.zshenv`, `.zlogin`, `.zlogout`, `.profile`, `.envrc`; `.npmrc`, `.yarnrc`, `.yarnrc.yml`, `.pnp.cjs`, `.pnp.loader.mjs`, `.pnpmfile.cjs`, `bunfig.toml`, `.bunfig.toml`; `.bazelrc`, `.bazelversion`, `.bazeliskrc`; `.pre-commit-config.yaml`, `lefthook.yml`, `lefthook.yaml`, `.lefthook.yml`, `.lefthook.yaml`; `gradle-wrapper.properties`, `maven-wrapper.properties`; `.devcontainer.json`; `.ripgreprc`, `pyrightconfig.json`; `.mcp.json`, `.claude.json` |

The operator decided the three parts below on 2026-10-06, in the story's Notes.

## Decision

**A story whose branch changes a path Claude Code protects is accepted by the operator only; flai checks `permission_prompt` against each Claude Code version flai serve runs that it has not checked; and `permission_prompt` handles every protected path inside an in-progress story's worktree except `.git`.** This refines ADR-0086.

### The list

flai keeps one list of the protected paths above, in a package of its own that `permission_prompt` and the acceptance preview both read, and the review page shows what the preview reports. A path is protected when one of its folders is a protected folder, `.config/git` counting as the two folders in order, or when its file name is a protected file. The list follows Claude Code's documentation; a release of Claude Code that adds to it is followed by a change to flai's list.

### 1. Acceptance is the gate

- The acceptance preview lists the files the story's branch changes on a protected path, beside what it already reports, and says that only the operator accepts the story. The dashboard's review page shows the same list and the same sentence.
- The operator is the story's owner or the project's owner, the manifest's `owner`, as [ADR-0097](0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) names who answers a permission thread. With neither named, the operator is anyone but the orchestrator.
- `flai accept` refuses an acceptance of such a story by anyone else, naming the files: a story's agent accepting on its own name, and the orchestrator under `accept_reviews`, for which the files are one more blocker of ADR-0093's preview. The refusal comes before anything is merged, whatever lets that agent accept other stories.
- An agent that runs `flai accept --by <operator>` on the operator's word, as this project's conventions allow, accepts as the operator.

Acceptance is the gate because it is the first moment the write can take effect: until then it sits on a branch the running agents do not read. A person reads the change there, in the diff the review page shows, and nothing but a person may let it through. With that gate, auto-approve lets a story's agent write the files during the story without asking, and loses no review.

### 2. The check against each new Claude Code

- When flai serve starts, and when it starts a claude-code agent, it reads the version of the `claude` it runs. When flai has no record of that version, it checks it once, beside the agent: the check never delays or stops an agent's start.
- The check makes one real write through `permission_prompt`, with the adapter's own permission arguments and the cheapest model: in a scratch project in a temporary folder, with one story in progress, its worktree, and auto-approve on for it, a headless `claude -p` is asked to Write a file in the worktree's `.claude/` folder. It passes when the file holds what was asked, and fails otherwise.
- flai records the version, when it was checked, and the outcome, with Claude Code's error on a failure, in its host state beside its configuration. A version with a record is not checked again, whatever its outcome.
- On a failure flai opens a thread to the operator on the manifest, `system-flow.yaml`, of each project it serves with claude-code agents, quoting Claude Code's error, naming the version, and saying that protected-path writes will not go through until it is fixed. The agents still run: a story that never writes a protected path does not need the prompt.
- One check costs about 0.02 USD on Haiku.

### 3. The prompt covers every protected path but `.git`

- `permission_prompt` approves an Edit, Write, MultiEdit, or NotebookEdit of any protected path inside an in-progress story's worktree, by ADR-0086's rules: at once under auto-approve, otherwise on the owner's answer on a thread.
- It still refuses a path with a `.git` folder or file in it: the worktree's `.git` is git's link to the repository, and the repository's history is never an agent's to write. It still refuses the main checkout, every path outside a story's worktree, a path a symbolic link takes out of it, and a story not in progress.
- A path that is not protected is refused as before: Claude Code does not ask for it.

## Consequences

- Auto-approve can be left on. A story's agent writes its `.claude/` files, `.mcp.json`, and the template's copies without a thread, and the operator reads them at acceptance.
- A story that changes a protected path waits in review for the operator, however the orchestrator judges it. The orchestrator says so on a thread, as ADR-0093 has it do for every blocker.
- A Claude Code update that breaks the prompt is found when flai serve first runs it, on a thread, rather than by a story that hangs.
- Each new Claude Code version costs one short headless run.
- The list is copied from Claude Code's documentation and can fall behind it. A path Claude Code adds is refused by the prompt until flai's list follows, as every protected path outside `.claude/` was before.
- The claim that a worktree's configuration does not reach the running agent rests on one test of settings and hooks. Claude Code also discovers skills and `CLAUDE.md` files in folders below the one it starts in when it works with files there, so a story's agent that reads its own worktree may load them. That is no different from what it reads in any other file of the worktree, and is not a protected-path write.

## Alternatives considered

- **`--permission-mode bypassPermissions`.** It lifts the protection everywhere, the main checkout's settings and `.git` included, with no review at all.
- **Moving the template's files out of `.claude/`.** It changes the template's layout, and this repository's own copy would still need generating into `.claude/`. To be taken up only if the prompt breaks again.
- **A flai tool that writes the files itself.** It goes around Claude Code's protection rather than through `--permission-prompt-tool`, which Claude Code offers for this.
- **Asking the operator on a thread during the story, as ADR-0086 does without auto-approve.** It keeps a person in the middle of every such story, and the write cannot take effect before acceptance anyway.
