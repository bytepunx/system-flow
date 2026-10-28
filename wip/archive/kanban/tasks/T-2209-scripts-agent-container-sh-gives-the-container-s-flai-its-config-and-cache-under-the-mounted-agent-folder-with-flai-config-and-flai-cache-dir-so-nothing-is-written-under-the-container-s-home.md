---
id: T-2209
type: task
nature: remediation
title: scripts/agent-container.sh gives the container's flai its config and cache under the mounted agent folder with FLAI_CONFIG and FLAI_CACHE_DIR, so nothing is written under the container's home
status: done
parent: S-0521
owner: alex
created: 2026-10-28T13:02:47Z
updated: 2026-10-28T13:11:18Z
transitions:
  - to: ready
    at: 2026-10-28T13:08:13Z
    by: agent-S-0521
  - to: in-progress
    at: 2026-10-28T13:08:15Z
    by: agent-S-0521
  - to: done
    at: 2026-10-28T13:11:18Z
    by: agent-S-0521
stream: S-0521
tags: [docker, scripts]
touches: [scripts/agent-container.sh, docs/operators/runbooks/agent-container.md]
usage:
  source: log
  seconds: 183
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7116
      cache_read: 1304235
      cache_write: 74927
      cost: 0.8653
---
# T-2209 scripts/agent-container.sh gives the container's flai its config and cache under the mounted agent folder with FLAI_CONFIG and FLAI_CACHE_DIR, so nothing is written under the container's home

## Work

`scripts/agent-container.sh` passes `CLAUDE_CONFIG_DIR`, `FLAI_AGENT`, `FLAI_SESSION`, and `FLAI_STORY` to the container and nothing for flai's own config, so the image's `flai` creates `$HOME/.flai/config.json` under `/home/agent` on its first run (I-0219; `config.Load` creates the file when `FLAI_CONFIG` names none, ADR-0008). Add `-e FLAI_CONFIG=$AGENT_DIR/flai/config.json` and `-e FLAI_CACHE_DIR=$AGENT_DIR/flai/cache` to the `docker run` line, with `mkdir -p "$AGENT_DIR/flai"` beside the `mkdir -p "$AGENT_DIR"` the script does, so that the config flai writes and the cache it uses land in the story's mounted agent folder, which is the container's to write, and the container's home stays untouched. Keep the script's header comment and the runbook true: `docs/operators/runbooks/agent-container.md` § What it mounts and its table of variables name the two variables and the folder they point under. Check by hand once the image is built with the tree's flai: `scripts/agent-container.sh S-nnnn -- flai config show` (or `flai --version` followed by a look at the folder) writes `.flai-cache/agents/S-nnnn/flai/config.json` and nothing under `/home/agent`, which `scripts/agent-container.sh S-nnnn -- ls -A /home/agent/.flai` shows absent. ShellCheck passes (`flai test scripts/agent-container.sh`). This task waits for nothing.

## Done when

- The `docker run` in `scripts/agent-container.sh` passes `FLAI_CONFIG` and `FLAI_CACHE_DIR` naming paths under the mounted agent folder, and the folder exists before the start
- A flai command run through the script writes its config under `.flai-cache/agents/S-nnnn/flai/` and nothing under the container's `/home/agent`, checked by hand and recorded in the narrative
- `docs/operators/runbooks/agent-container.md` names the two variables and where they point
- ShellCheck and the markdown lint pass on the two files

## Notes

pilot (S-0507) and the container adapter (S-0508) start the container in place of this script later; the plan thread asks them to set the same two variables under the cache folder the start request names.
