---
id: T-2900
type: task
nature: remediation
title: Probe the rebuilt image on S-0682, check its narrative is untouched, and close I-0224
status: done
parent: S-0682
owner: alex
created: 2026-11-10T08:46:31Z
updated: 2026-11-10T09:25:44Z
transitions:
  - to: ready
    at: 2026-11-10T09:08:43Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T09:08:44Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T09:25:44Z
    by: agent-S-0682
stream: S-0682
tags: [containers, issues]
touches: [design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md, design/issues/summary.md, design/system/flai-cli.md, docs/operators/runbooks/agent-container.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-2897, T-2898]
usage:
  source: log
  seconds: 1020
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 76
      output: 14448
      cache_read: 8711941
      cache_write: 155071
      cost: 3.0749
---
# T-2900 Probe the rebuilt image on S-0682, check its narrative is untouched, and close I-0224

## Work

Waits for T-2897 and T-2898: the trial runs the probe pilot and the script now make, in an image built from them.

- Rebuild the images from the story's worktree with `make agent-images`, so pilot carries the probe.
- Run the free check first: `scripts/agent-container.sh --probe S-0682 -- env`, and record in this task's Notes that `FLAI_AGENT=probe-S-0682` is set and `FLAI_STORY` is not.
- The capped trial calls a paid model. Ask on a thread on S-0682 before it, recommending one run with `FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-0682`. With the operator's yes, run it, and check that S-0682's narrative's Current state, Next steps, and Log are what they were before it and that no item changed (`flai stream state` and `inbox` show nothing by `probe-S-0682`). Record what it reported. Without docker or a Claude login on the host, record that instead and rest criterion 1 on T-2897's test.
- Tick criterion 1 with `flai criteria tick S-0682 1` once T-2897's test passes and the checks above hold.
- Run `flai issue close I-0224 --reason "<reason>"`: the reason names `scripts/agent-container.sh --probe` and the start request's `probe`, which runs the image as `probe-S-nnnn` with a prompt that works no item and no `FLAI_STORY`, held by T-2897's test. Then tick criterion 2.

## Done when

- The free check, and the capped trial or the reason it did not run, are recorded in this task's Notes.
- I-0224 is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- Both of S-0682's criteria are ticked.

## Notes

- Images rebuilt from story/S-0682 with `make agent-images` at 2026-11-10T09:13Z. The first try failed on Docker Desktop's credential helper (I-0339); a scratch `DOCKER_CONFIG` with no credential store built them, leaving `~/.docker/config.json` as it was.
- Free check, `FLAI_CLAUDE_LOGIN=~/.claude-me scripts/agent-container.sh --probe S-0682 -- env`, exit 0: `FLAI_AGENT=probe-S-0682`, a new `FLAI_SESSION`, and no `FLAI_STORY`.
- Capped probe, on the operator's yes on TH-1000, started 2026-11-10T09:14:40Z: `FLAI_CLAUDE_LOGIN=~/.claude-me FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-0682`, exit 0, 0.0098 USD. It called `inbox` (as `probe-S-0682`), `board`, `item_get S-0704`, `doc_get design/conventions/README.md`, `doc_search`, and one Bash call for the three commands, and no writing tool. S-0682's narrative, story file, and five task files had the same sha256 after it as before, and `.flai-cache/events.jsonl` gained no line.
- The probe reported `flai --version` failing with `unknown flag: --version`: the command is `flai version`. The prompt, its test, the runbook, and flai-cli.md now say `flai version`.
