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
- Never hand the explorer or the verifier an edit, and never hand any sub-agent a commit, a transition, or a question for the designer.
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
  - start with nothing but your prompt, unless a fork, which holds your conversation
  - are given the worktree's path, the story and task IDs, the question, what you already know, and the shape of the answer you want.
  - ask for summary vs. raw output: the finding, with paths and line numbers, and failing output quoted only as far as it matters.
  - their answer is evidence, not verdict.
  - check what you act on, record in narrative the finding if deciding something.
  - don't delegate if it's slower: one file you know, one search whose answer you expect, one short command.
  - run sub-agents in parallel when their questions are independent.
  - do not act on what an answer will cover until it arrives.
  - a sub-agent question for the designer is yours to ask, on the story or the task.

## Tasks worked by sub-agents

- When a layer of the plan (`work-management.md`) holds more than one task, run them at once, one task sub-agent each, and wait for all of them before you act on any. Work a layer of one task yourself.
- Use a fork, which inherits your conversation, where the harness offers one. Otherwise give the sub-agent the task's ID, the worktree's path, the paths the task touches, and the conventions to read.
- A layer's tasks share the story's worktree when their `touches` have no path in common. Give each task its own worktree from the story's branch (`git worktree add -b task/T-nnnn <path> <story branch>`), and merge it back yourself, when one builds or tests what another changes, so that a half-done edit cannot fail a sibling's tests.
- A task's `touches` are a guess made before its code is read. Before you commit a task's work, compare the files it changed with its `touches`. When it changed a path another task of the layer touches, review both together, and redo the later one where their edits met.
- Review each one's work as your own before you accept it: read its diff against the task's `## Done when`, run the tests for what it changed, and fix or finish what falls short yourself. Then commit it and move the task.
- Only you commit, sync the stream, move items, keep the narrative, and talk to the designer. Record in the narrative what each task sub-agent did, and the interference it reported.

## As an explorer or a verifier

- You work for the agent that started you, not for the designer or the board. When your prompt names a story, prime with `prime` and your role (`flai prime --story <id> --role explore` or `--role verify`), reading only the sections you need.
- Use the worktree specified.
- Do not edit files.
- Do not move, create, or edit work items, write to threads, read the inbox, or wait for events or work; your tools leave them out; never work around that through the shell.
- Answer the question you were given. Do not do the story's work.
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the agent that started you sees. Lead with the answer, then the evidence, then what you could not check.

## As a task sub-agent

- You work one task for the story's agent, in the worktree your prompt names. A fork already holds the conventions; otherwise prime with `prime` and the story, and read the conventions your prompt names.
- Edit only the paths your task touches. When the task needs a change outside them, stop and say so in your final message rather than make it.
- Run only the tests for what you changed; leave the whole suite, the lint, and `flai check` to the story's agent.
- When a build or test fails in a path your task does not touch, another task's edit may be half done: do not fix it. Wait a minute and run it again, and say in your final message what failed, when, and for how long.
- Do not commit, stage, or otherwise write to git. Do not move, create, or edit work items, write to threads, read the inbox, wait for events or work, or sync the stream; the guard refuses them, and you do not work around it.
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the story's agent sees. Lead with done or not done, then the files you changed, the tests you ran and their results, the decisions a reader could have made differently, and what is left.

## When in doubt

- If you need only the conclusion of a long piece of work, delegate it.
- If the work changes anything and is not a task you handed over, it is the story's agent's, not a sub-agent's.

<!-- system-flow:end-of-baseline -->

## Project additions
