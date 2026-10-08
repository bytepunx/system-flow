---
title: Delegation
updated: 2026-10-08
audience: agent
order: 130
status: active
topics: [all]
roles: [story, explore, verify]
---

# Delegation

When and how the agent working a story hands work to a sub-agent and what it may do.

## Rules

- The template defines two sub-agents for harnesses that have them:
  - an explorer that finds and reads
  - a verifier that reviews a story's diff against its acceptance criteria and the conventions
  - neither can change files, work items, or threads
- Keep your own context for decisions and edits.
- Hand sub-agents work that you only need the conclusion from:
  - search across many files
  - long logs
  - reviews of a large diff
- Use the explorer to find and read code, designs, and logs
- Use the verifier to check a diff against the story's acceptance criteria and these conventions. Run the whole suite, the whole lint, and `flai check` yourself, through `flai verify`, as below ([ADR-0110](../adrs/0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md)).
- Never hand the explorer or the verifier an edit, and never hand any sub-agent a commit, a transition, or a question for the designer.
- While working:
  - run only the tests for what you changed, yourself, with `flai test` on the paths you changed, or the MCP tool `test` with those paths: it runs the test and lint tiers the manifest's `tests` declare that those paths select, cheapest first, and answers pass or the first findings
  - do not run `go test`, vitest, golangci-lint, or gofmt by hand and read their logs, and do not hand that run to a sub-agent
  - leave the whole suite, the whole lint, and `flai check` to the close-out before review
  - make independent edits and commands in one turn, as several tool calls in one message: consecutive edits to one file, reads of files you already know, commands that do not wait on each other
  - move a task you have just written to ready and in-progress in one command, `flai move T-nnnn ready && flai move T-nnnn in-progress`, since `flai move` refuses a task straight from `backlog` to `in-progress`
- Before moving a story to `review`:
  - commit what is outstanding, run `flai stream sync` again and resolve what it reports, as `git.md` says
  - run the close-out yourself, once, in the worktree: the project's close-out script, which runs `flai verify S-nnnn --record-issues` and then commits (`work-management.md`), or `flai verify S-nnnn`, or the MCP tool `verify`, where the project has no close-out script
  - run it as one command with the longest timeout the harness allows (Claude Code's Bash tool: 600000 ms), its exit status echoed after it on the same line: `scripts/close-out.sh S-nnnn -m "<message>"; echo "exit $?"`; never pipe its output, which loses the exit status, and never redirect it into a file
  - read its last line, which names the outcome and the step it stopped at, and the findings of the step that failed
  - when it stops, fix what it names, commit, and run it again; do not finish the steps by hand, and never hand the run to a sub-agent
  - its passing run is the story's run before review; do not repeat it. `flai verify S-nnnn --last` prints it again
  - check the diff against the story's acceptance criteria and these conventions, saying which criteria, by number, the diff meets and which it does not
  - hand that review to one fresh verifier when it pays, as when the diff is too large to read in your own context; tell it the run passed, and at which commit, and to read the result with `flai verify S-nnnn --last` rather than run the suite
  - if the review finds an issue, fix it, commit, and run the close-out again.
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

- Hand each task of the plan to a task sub-agent, so that its reading and editing stay out of your context. When a layer (`work-management.md`) holds more than one task, you may run them at once, one task sub-agent each, and then wait for all of them before you act on any. Launch each sub-agent in the foreground, a layer's in one message, so that each result comes back as the tool's result however long the sub-agent runs (in Claude Code, the Agent tool's `run_in_background` set to false), and go on from there. Never end your turn while a sub-agent runs in the background: a headless Claude Code session ends ten minutes after its turn ends, and the sub-agent with it. When only the designer's answer is left (`work-management.md`), end the session once every sub-agent is back. Never wait for a sub-agent with `wait_for_events` either, which reports work items and threads, not sub-agents: it is for a thread awaiting the designer, held by an agent that `flai serve` does not start again, and `flai guard` refuses it to a story's agent while a sub-agent of its session runs and no thread on the story or one of its tasks is open. That saves time only when the tasks are long beside what only you do: priming, reviewing each task, verifying, and closing out.
- Use a fork, which inherits your conversation, where the harness offers one. Otherwise give the sub-agent the task's ID, the worktree's path, the paths the task touches, and the conventions to read.
- Name the task's ID in each task sub-agent's description, the Agent tool's `description`, so that flai measures the task by its sub-agent's calls rather than by the time it was in progress. An explorer or verifier started for one task names that task too; one started for the story as a whole names none.
- A layer's tasks share the story's worktree when their `touches` have no path in common. Give each task its own worktree from the story's branch (`git worktree add -b task/T-nnnn <path> <story branch>`), and merge it back yourself, when one builds or tests what another changes, so that a half-done edit cannot fail a sibling's tests.
- A task's `touches` are a guess made before its code is read. Before you commit a task's work, compare the files it changed with its `touches`. When it changed a path another task of the layer touches, review both together, and redo the later one where their edits met.
- Review each one's work as your own before you accept it: read its diff against the task's `## Done when`, run `flai test` on the paths it changed, and fix or finish what falls short yourself. Then commit it and move the task. Tick each acceptance criterion the sub-agent named once your review has verified it (`work-management.md`).
- Only you commit, sync the stream, move items, keep the narrative, and talk to the designer. Record in the narrative what each task sub-agent did, and the interference it reported.

## As an explorer or a verifier

- You work for the agent that started you, not for the designer or the board. When your prompt names a story, prime with `prime` and your role (`flai prime --story <id> --role explore` or `--role verify`), reading only the sections you need.
- Use the worktree specified.
- Do not edit files.
- Do not move, create, or edit work items, write to threads or send messages, read the inbox, or wait for events or work; your tools leave them out; never work around that through the shell.
- Answer the question you were given. Do not do the story's work.
- As a verifier, review the diff against the story's acceptance criteria and the conventions, naming the criteria it meets by number. Do not run the tests, the lint, or `flai check` to learn whether they pass: read the story's last result with `flai verify S-nnnn --last`, and say so when it did not pass or its commit is not the branch's head.
- Run the close-out, or `flai verify`, only when your prompt asks. Then run it once, as one command with the longest timeout the harness allows (Claude Code's Bash tool: 600000 ms), its exit status echoed after it on the same line: `flai verify S-nnnn; echo "exit $?"`. Never pipe its output, which loses the exit status, and never redirect it into a file. Read its last line, which names the outcome and the step it stopped at.
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the agent that started you sees. Lead with the answer, then the evidence, then what you could not check.

## As a task sub-agent

- You work one task for the story's agent, in the worktree your prompt names. A fork already holds the conventions; otherwise prime with `prime` and the story, and read the conventions your prompt names.
- Edit only the paths your task touches. When the task needs a change outside them, stop and say so in your final message rather than make it.
- Run only the tests for what you changed, with `flai test` on the paths you changed or the MCP tool `test`, rather than running the test and lint tools by hand; leave the whole suite, the whole lint, and `flai check` to the story's agent.
- When a build or test fails in a path your task does not touch, another task's edit may be half done: do not fix it. Wait a minute and run it again, and say in your final message what failed, when, and for how long.
- Do not commit, stage, or otherwise write to git. Do not move, create, or edit work items, write to threads or send messages, read the inbox, wait for events or work, or sync the stream; the guard refuses them, and you do not work around it.
- Record a decision your task makes as an ADR through flai, never by copying what it writes or guessing a number: `flai adr new`, or the MCP tool `adr_new`, numbers it, names the file, and adds its index row, and `flai adr topics` sets its topics. Give none of them `--commit` or `--autocommit`, nor `adr_new` its `commit`: the guard refuses those, and the story's agent commits the ADR with the rest of your task ([ADR-0127](../adrs/0127-flai-guard-lets-a-story-s-sub-agent-write-an-adr-with-flai-adr-new-topics-and.md)).
- When you need the designer to decide something, stop and put the question in your final message, with your recommended answer first.
- Your final message is all the story's agent sees. Lead with done or not done, then the files you changed, the story's acceptance criteria your task's work meets, by their numbers in `flai criteria list`, the tests you ran and their results, the decisions a reader could have made differently, and what is left. You do not tick them; the story's agent does after its review.

## When in doubt

- If you need only the conclusion of a long piece of work, delegate it.
- If the work changes anything and is not a task you handed over, it is the story's agent's, not a sub-agent's.

<!-- system-flow:end-of-baseline -->

## Project additions
