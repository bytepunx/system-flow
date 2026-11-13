---
id: TH-1172
title: "S-0760 criterion 1: no flai command writes agent.image.args, so the manifest's key can only be written by hand"
anchor:
  path: wip/kanban/stories/S-0760-this-repository-builds-its-agent-images-through-flai-s-general-path-with.md
  item: S-0760
status: resolved
participants: [agent-S-0760, orchestrator, alex]
created: 2026-11-13T19:48:03Z
updated: 2026-11-13T20:00:57Z
---

# TH-1172 S-0760 criterion 1: no flai command writes agent.image.args, so the manifest's key can only be written by hand

On wip/kanban/stories/S-0760-this-repository-builds-its-agent-images-through-flai-s-general-path-with.md.

## Entries

### 2026-11-13T19:48:03Z agent-S-0760
**Question.** Criterion 1 says `system-flow.yaml`'s `agent.image` names the layer's build arguments "written with the flai command that sets it". No flai command sets `agent.image.args`: `flai manifest set` lists no `agent.*` key, and `flai agent set` keeps `agent.image` as the file has it. `design/system/project-manifest.md` says the same of `agent.image.layer`: "only a hand edit changes it".

**Recommended:** I write the five arguments under `agent.image.args` by hand (`GO_VERSION`, `GOLANGCI_LINT_VERSION`, `SHELLCHECK_VERSION`, `PNPM_VERSION`, `GORELEASER_VERSION`, the values `scripts/agent-image.sh --layer` passes today), run `flai check`, and record the missing command as an issue on S-0760 for you to make a story of at acceptance.

**The other way:** I add `agent.image.args` to `flai manifest set` in this story, with its validation, docs, and tests. That widens S-0760 into `flai/cmd` and adds roughly an hour.

I am going on with the hand edit and the other three tasks meanwhile; if you choose the command, I replace the hand edit with it.

### 2026-11-13T19:49:31Z orchestrator (recommendation)
**Recommended:** take agent-S-0760's hand edit and keep S-0760's scope.

- Write the five arguments under `agent.image.args` by hand, with the values `scripts/agent-image.sh --layer` passes today.
- Keep every other line of `system-flow.yaml` as it is, and run `flai check`, which reports a bad block as `manifest.agent`.
- File the missing command as an issue on S-0760, to be made a story later.
- Reword criterion 1 so it no longer says "written with the flai command that sets it". The criterion is the operator's words, so the change needs your confirmation.

**Why:**

- `design/system/project-manifest.md` says `flai agent set` and `flai agent clear` keep `agent.image` as the file has it, "so only a hand edit changes it". No command writes `agent.image.args` today, so the criterion asks for a command that does not exist.
- The digest records the same exception twice, both handled by a hand edit validated with `flai check`: D-0449 for the `tests` tiers and D-0593 for `projects[].install`.
- The other way adds `agent.image.args` to `flai manifest set`. That widens S-0760 into `flai/cmd`, with its own validation, docs, and tests, for about an hour more. It is a change to a command that deserves its own story and criteria.

**Alternative:** answer "command" here, and agent-S-0760 adds `agent.image.args` to `flai manifest set` in this story and replaces the hand edit with it.

Source: design/system/project-manifest.md § Project manifest

### 2026-11-13T19:53:23Z alex
Confirmed the recommendation of 2026-11-13T19:49:31Z orchestrator.

Source: design/system/project-manifest.md § Project manifest

### 2026-11-13T20:00:57Z agent-S-0760
Resolved: Done as confirmed: agent.image.args written by hand and checked with flai check and flai image status (4c3a215), the missing command recorded as I-0364, and criterion 1 and the Planning note reworded to say so.
