---
id: D-0520
title: "This repository's .opencode/plugins/flai-guard.js differs from the template's in its COMMAND constant alone: it runs scripts/flai.sh, as .claude/settings.json does, so a change to the plugin is made in both copies"
kind: decision
sources: [S-0407]
topics: [harness, guard, template, agents]
paths: [".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, scripts/flai.sh, ".claude/settings.json"]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0520 This repository's .opencode/plugins/flai-guard.js differs from the template's in its COMMAND constant alone: it runs scripts/flai.sh, as .claude/settings.json does, so a change to the plugin is made in both copies

S-0407 added flai's guard plugin for OpenCode in two copies. The template's, `template/root/.opencode/plugins/flai-guard.js`, runs the installed `flai`. This repository's, `.opencode/plugins/flai-guard.js`, differs in one constant, `COMMAND`, which runs `scripts/flai.sh` resolved from the plugin's own path, so the guard on a story branch here is the tree's own flai.

From S-0407's Decisions (`wip/archive/agents/S-0407.md`): "`.opencode/plugins/flai-guard.js` here differs from the template's in one constant: it runs `scripts/flai.sh`, as `.claude/settings.json` runs the guard, so the guard on this branch is the tree's (delegation.md project additions)."

What follows for a later story: change the template's plugin first and carry the change to this repository's copy, keeping its `COMMAND`, as `design/conventions/delegation.md` asks for `.claude/` under its project additions. Both paths are under `.opencode/`, which is protected on every harness (D-0514), so the story's agent writes them itself and the operator accepts the story.

Source: S-0407.
