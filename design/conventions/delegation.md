---
title: Delegation
updated: 2026-10-02
audience: agent
order: 130
status: active
topics: [all]
roles: [explore, verify]
---

# Delegation

When and how the agent working a story hands work to a sub-agent and what it may do.

## Rules

- The template defines two sub-agents for harnesses that have them:
  - an explorer that finds and reads
  - a verifier that runs the project's checks
  - neither can change files, work items, or threads
- Keep your own context for decisions and edits.
- Hand sub-agents work that you only need the conclusion from:
  - search across many files
  - test and lint runs
  - long logs
  - diff checks
- Use the explorer to find and read code, designs, and logs
- Use the verifier to run tests, lint, `flai check`, and checking diffs.
- Never hand a sub-agent an edit, a transition, or a question for the designer.
- While working:
  - run only the tests for what you changed
  - leave the whole suite tests, linting, and `flai check` to the verifier
- Before moving a story to `review`:
  - commit what is outstanding and have one fresh verifier run the whole suite, the lint, and `flai check` in the worktree through the project's close-out script where it has one
  - and check the diff against the story's acceptance criteria and these conventions
  - if an issue is found, fix it, commit, and run a fresh verifier to confirm the fixes
  - a verifier's passing run is the story's run before review; do not repeat it.
- Fix what a verifier finds yourself; never delegate a fix to a sub-agent.
- Sub-agents:
  - start with nothing but your prompt
  - are given the worktree's path, the story and task IDs, the question, what you already know, and the shape of the answer you want.
  - ask for summary vs. raw output: the finding, with paths and line numbers, and failing output quoted only as far as it matters.
  - their answer is evidence, not verdict.
  - check what you act on, record in narrative the finding if deciding something.
  - don't delegate if it's slower: one file you know, one search whose answer you expect, one short command.
  - run sub-agents in parallel when their questions are independent.
  - do not act on what an answer will cover until it arrives.
  - a sub-agent question for the designer is yours to ask, on the story or the task.

## As a sub-agent

- You work for the agent that started you, not for the designer or the board. When your prompt names a story, prime with `prime` and your role (`flai prime --story <id> --role explore` or `--role verify`), reading only the sections you need.
- Use the worktree specified.
- Do not edit files.
- Do not move, create, or edit work items, write to threads, read the inbox, or wait for events or work; your tools leave them out; never work around that through the shell.
- Answer the question you were given. Do not do the story's work.
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the agent that started you sees. Lead with the answer, then the evidence, then what you could not check.

## When in doubt

- If you need only the conclusion of a long piece of work, delegate it.
- If the work changes anything, it is the story's agent's, not a sub-agent's.

<!-- system-flow:end-of-baseline -->

## Project additions
- `.claude/settings.json` runs the guard from this tree, `scripts/flai.sh guard`, passing on only its refusals as the template's does, so it is the guard on this branch; a project made from the template runs the installed `flai guard` (ADR-0060). The agent definitions in `.claude/agents/` are copies of the template's; change them in `template/root/.claude/agents/` first.
