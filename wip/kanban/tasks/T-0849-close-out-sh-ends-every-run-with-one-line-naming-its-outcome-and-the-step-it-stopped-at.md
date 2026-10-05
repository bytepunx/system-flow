---
id: T-0849
type: task
nature: improvement
title: close-out.sh ends every run with one line naming its outcome and the step it stopped at
status: backlog
parent: S-0266
owner: alex
created: 2026-10-05T00:18:50Z
updated: 2026-10-05T00:18:50Z
transitions: []
stream: S-0266
tags: [scripts, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md]
---
# T-0849 close-out.sh ends every run with one line naming its outcome and the step it stopped at

## Work

Make `scripts/close-out.sh`, and the template's copy in `template/root/scripts/close-out.sh`, end every run with one last line on standard output. The line names the story, the outcome, and the step, whether the run passes or stops. For example:

- `close-out: S-nnnn stopped at flai check --strict (exit 1)`
- `close-out: S-nnnn passed every step; ready to move to review`

Set a `step` variable before each step runs. Print the line from an `EXIT` trap that reads the exit status, so that a step stopped by `set -e` is named too. The script still exits with the failing step's status, and a usage error still exits 2. Keep the per-step messages the script prints now.

Update the `close-out.sh` row in `scripts/README.md` and `template/root/scripts/README.md` to say what the last line holds.

This task waits for none. It is the first layer, because the verifier's definition, the story agent's prompt, and the convention all describe the line it prints.

## Done when

- A run that stops ends with the line naming that step, and exits with that step's status. For example, a run where the narrative's `## Current state` is empty.
- A run that passes ends with the passing line and exits 0.
- A run with no story still prints its usage and exits 2.
- Both scripts print the same line for the same outcome.
- `scripts/template-test.sh` and the markdown lint pass.

## Notes

Drafted by the planner for S-0266. The template's script gets the line too, because the template's verifier definition will tell every project's verifier to read it.
