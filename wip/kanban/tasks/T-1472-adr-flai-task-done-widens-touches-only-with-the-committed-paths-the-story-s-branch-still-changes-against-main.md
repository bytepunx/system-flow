---
id: T-1472
type: task
nature: remediation
title: "ADR: flai task done widens touches only with the committed paths the story's branch still changes against main"
status: done
parent: S-0370
owner: alex
created: 2026-10-12T17:09:44Z
updated: 2026-10-12T20:06:02Z
transitions:
  - to: ready
    at: 2026-10-12T20:05:25Z
    by: agent-S-0370
  - to: in-progress
    at: 2026-10-12T20:05:26Z
    by: agent-S-0370
  - to: done
    at: 2026-10-12T20:06:02Z
    by: agent-S-0370
stream: S-0370
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md]
usage:
  source: log
  seconds: 36
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 73
      cache_read: 710432
      cache_write: 4463
      cost: 0.3249
---
# T-1472 ADR: flai task done widens touches only with the committed paths the story's branch still changes against main

## Work

Record the remedy for I-0121 as a new ADR with `flai adr new`, refining ADR-0107 (step 5, Touches) and ADR-0128 (point 4), and accept it with `flai adr accept`.

- Context: I-0121's instance. S-0232 agreed on MS-0014 to leave `scripts/flai-snapshot.sh` to S-0340, removed it from T-1223's and the story's touches, and reverted the file to main's version. Closing T-1223 again, `flai task done` added the path back to both, because the revert commit changed it, though the branch no longer differed from main there. The claim overlapped S-0340 again until it was narrowed by hand a second time.
- Decision: the touches step widens the task's and the story's touches with the paths the commit changed that the story's branch still changes against its merge base with the main branch (`git diff --name-only main...HEAD`, both sides of a rename), read after the sync step. A path the commit changed back to main's version is not added. A path left for another task is still never added (ADR-0128). Widening still only adds: it never removes a touch.
- Consequences: a story that hands a path to another story and reverts it keeps its narrowed claim. The tell step (S-0333) is unchanged: it still messages a story whose claim covers a path the commit changed, a revert included, since that story may want to know the branch no longer changes it. The `--commit` writes of `flai adr new` and `flai issue` (S-0275) are unchanged, since they commit only new content they wrote.

`flai adr new` names the ADR file; `flai task done` widens this task's touches with it when it commits it.

Waits for nothing: it is the first layer, and the code and docs tasks build what it decides.

## Done when

- The ADR is accepted under `design/adrs/`, refines ADR-0107 and ADR-0128, and is listed in `design/adrs/README.md`.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

Drafted by the planner for S-0370. The ADR's number is not named: S-0363 also adds an ADR, and whichever story writes second takes the next number.
