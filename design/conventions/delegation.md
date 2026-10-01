---
title: Delegation
updated: 2026-10-01
audience: agent
order: 130
status: active
topics: [all]
roles: [explore, verify]
---

# Delegation

When the agent working a story hands work to a sub-agent, what it hands over, and what a sub-agent may do. The template defines two sub-agents for harnesses that have them: an explorer that finds and reads, and a verifier that also runs the project's checks. Neither can change a file, a work item, or a thread.

## Rules

- Keep your own context for decisions and edits. Hand a sub-agent the work whose output you need only the conclusion of: search across many files, test and lint runs, long logs, and a check of your diff.
- Use the explorer to find and read code, design, and logs. Use the verifier to run the tests, lint, and `flai check`, and to check a diff. Never hand a sub-agent an edit, a transition, or a question for the designer.
- Before moving a story to `review`, have a fresh verifier check the diff in the worktree against the story's acceptance criteria and these conventions, and act on what it finds.
- A sub-agent starts with nothing but your prompt. Give it the worktree's path, the story and task IDs, the question, what you already know, and the shape of the answer you want.
- Ask for a summary, not raw output: the finding, with paths and line numbers, and failing output quoted only as far as it matters.
- A sub-agent's answer is evidence, not a verdict. Check what you act on, and say in the narrative what a sub-agent found when it decided something.
- Do it yourself when that is quicker: one file you know, one search whose answer you expect, one short command.
- Run sub-agents side by side when their questions are independent. Do not act on what an answer will cover until it arrives.
- A question a sub-agent returns for the designer is yours to ask, on the story or the task, as your own questions are.

## As a sub-agent

- You work for the agent that started you, not for the designer or the board. When your prompt names a story, prime with `prime` and your role (`flai prime --story <id> --role explore` or `--role verify`), and read only the sections you need.
- Work in the worktree your prompt names.
- Do not edit files. Do not move, create, or edit work items, write to threads, read the inbox, or wait for events or work; your tools leave them out, and you do not work around that through the shell.
- Answer the question you were given. Do not do the story's work.
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the agent that started you sees. Lead with the answer, then the evidence, then what you could not check.

## When in doubt

- If you need only the conclusion of a long piece of work, delegate it.
- If the work changes anything, it is the story's agent's, not a sub-agent's.

<!-- system-flow:end-of-baseline -->

## Project additions
- `.claude/settings.json` runs the guard from this tree, `scripts/flai.sh guard`, so it is the guard on this branch; a project made from the template runs the installed `flai guard` (ADR-0060). The agent definitions in `.claude/agents/` are copies of the template's; change them in `template/root/.claude/agents/` first.
