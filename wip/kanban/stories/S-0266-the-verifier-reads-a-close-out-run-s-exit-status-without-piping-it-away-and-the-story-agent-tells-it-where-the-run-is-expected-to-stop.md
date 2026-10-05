---
id: S-0266
type: story
nature: improvement
title: The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop
status: backlog
owner: alex
created: 2026-10-05T00:06:34Z
updated: 2026-10-05T00:06:34Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0266 The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop

## Goal

S-0248's first verifier ran `scripts/close-out.sh` three times, about two minutes each: the first run piped the output through `tail`, which lost the script's exit status; the second printed `$?` after the pipe, which was tail's; the third redirected the output into a file, which the verifier's own rules forbid. The verifier's definition (`.claude/agents/verifier.md`, and its copy in `template/root/.claude/agents/verifier.md`) bans writing files but does not say how to read a long script's status, and `scripts/close-out.sh` prints its success line only when every step passes, so a run that stops gives nothing to tell a stop from an interrupted tail. The story agent's prompt (`flai/internal/harness/harness.go`) also hands the verifier the run without naming the step it is known to stop at, I-0057's `flai check --strict`, so the verifier treated the expected stop as something to reproduce. One run, read once, is the aim: the agent log `~/.flai/serve/agents/sf-S-0248-20261004T231346Z.log` from 23:18:07Z to 23:23:54Z shows the waste.

## Acceptance criteria
- [ ] The verifier definition says how to run a long script and read its exit status in one command without a pipe or a file, with the timeout to give it, in `.claude/agents/verifier.md` and `template/root/.claude/agents/verifier.md`
- [ ] `scripts/close-out.sh` ends every run with one line that says its outcome and the step it stopped at, so that the verifier never needs to run it again to learn why it stopped
- [ ] The story agent's prompt, when it hands a close-out run to the verifier, names any step it already knows will stop and why, so that the verifier reports it and goes on to the steps after it rather than re-running
- [ ] `design/conventions/delegation.md` and the template's copy say the same, and the design (`design/system/devex.md` or where close-out is described) records it

## Tasks

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04. Companion to S-0249, which removes the expected stop itself.
