---
id: T-0850
type: task
nature: improvement
title: The verifier's definition says how to run a long script once and read its exit status without a pipe or a file
status: backlog
parent: S-0266
owner: alex
created: 2026-10-05T00:18:58Z
updated: 2026-10-05T00:18:58Z
transitions: []
stream: S-0266
tags: [agents, template]
touches: [".claude/agents/verifier.md", template/root/.claude/agents/verifier.md]
after: [T-0849]
---
# T-0850 The verifier's definition says how to run a long script once and read its exit status without a pipe or a file

## Work

Write the rule in `template/root/.claude/agents/verifier.md` first. Then copy it to `.claude/agents/verifier.md`, as the project addition in `delegation.md` asks. Add it to step 3, which already forbids redirection into files.

The rule says:

- Run a long script, such as the close-out, as one Bash call, with `; echo "exit $?"` after it on the same line.
- Never pipe the output. With `| tail`, the script's exit status is lost, and `$?` after the pipe is tail's.
- Never redirect the output into a file.
- Give the Bash tool its longest timeout, 600000 ms, so that the run is not cut short.
- Read the script's last line. It gives the outcome and the step the run stopped at (T-0849).
- Run the script once.
- When the prompt names a step known to stop, report that stop as expected. Then run the steps after it through their own entry points, rather than running the script again.

On a full flai-tier run, check that the last line survives the Bash tool's output limit. If it does not, record in the narrative's Decisions how the verifier reads the line instead.

This task waits for T-0849, so that it describes the line T-0849 prints.

An agent that flai serve starts may be refused writes under `.claude/` (I-0069, S-0257). If `.claude/agents/verifier.md` is refused, ask the operator on a thread on S-0266, giving the whole file to paste, not a fragment. Do not work around the refusal through the shell.

## Done when

- The two files are identical: `diff` finds nothing.
- The rule names:
  - the one-command form
  - the timeout
  - the ban on a pipe and on a file
  - the last line
  - what to do with an expected stop.
- The markdown lint passes.

## Notes

Drafted by the planner for S-0266. The S-0248 log, from 23:18:07Z to 23:23:54Z, shows the three runs this rule prevents.
