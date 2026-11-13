---
id: D-0615
title: "S-0760 criterion 1: no flai command writes agent.image.args, so the manifest's key can only be written by hand: Recommended: take agent-S-0760's hand edit and keep S-0760's scope"
kind: answer
sources: [TH-1172, S-0760]
topics: [docker, template, release, flai, cli, go, conventions, code, all]
paths: [scripts/agent-image.sh, scripts/agent-image-test.sh, Makefile, system-flow.yaml, flai/images/system-flow/Dockerfile, docs/operators/runbooks/agent-container.md, docs/operators/runbooks/install.md, docs/operators/index.md, docs/users/flai.md, design/system/agent-adapters.md, design/tech/docker.md, template/root/docs/operators/index.md.tmpl, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, design/system/project-manifest.md, design/system/flai-cli.md, design/issues/I-0364-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md]
date: 2026-11-13
status: absorbed
absorbed_by: design/system/project-manifest.md
asked_again: []
cited: []
---

# D-0615 S-0760 criterion 1: no flai command writes agent.image.args, so the manifest's key can only be written by hand: Recommended: take agent-S-0760's hand edit and keep S-0760's scope

The operator, alex, answered TH-1172 at 2026-11-13T19:53:23Z.

What was asked:

```text
**Question.** Criterion 1 says `system-flow.yaml`'s `agent.image` names the layer's build arguments "written with the flai command that sets it". No flai command sets `agent.image.args`: `flai manifest set` lists no `agent.*` key, and `flai agent set` keeps `agent.image` as the file has it. `design/system/project-manifest.md` says the same of `agent.image.layer`: "only a hand edit changes it".
```

The answer:

```text
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
```

S-0760 took the recommendation. From its Decisions (`wip/archive/agents/S-0760.md`):

```text
agent.image.args written by hand: no flai command writes it (flai manifest set lists no agent key; flai agent set keeps agent.image as the file has it). Asked on TH-1172, recommending the hand edit and an issue for the missing command.
```

The missing command is I-0364. `design/system/project-manifest.md` and `design/system/flai-cli.md` now say the five `args` are written by hand, since no flai command writes them. The same hand-edit exception holds for the `tests` tiers (D-0449) and `projects[].install` (D-0593).

Sources: TH-1172, S-0760.
