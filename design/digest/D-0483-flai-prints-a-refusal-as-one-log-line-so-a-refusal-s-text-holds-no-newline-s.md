---
id: D-0483
title: "flai prints a refusal as one log line, so a refusal's text holds no newline: S-0577's refusal with a multi-line block example showed its line breaks as \\n"
kind: decision
sources: [S-0577, T-2475]
topics: [flai, cli, go, code]
paths: [flai/cmd/kluster.go, flai/cmd/serve_actions.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0483 flai prints a refusal as one log line, so a refusal's text holds no newline: S-0577's refusal with a multi-line block example showed its line breaks as \n

`flai kluster` and `flai serve enable kluster` refuse a project with no `kluster:` block, and the refusal names the block to add. S-0577 first wrote that block as a YAML example over several lines. flai prints a refusal as one log line, so the example came out with its line breaks as a literal `\n`. S-0577 rewrote the refusal as one line.

From S-0577's narrative, `## Decisions`:

```text
The no-block refusal was made one line: flai prints a refusal as one log line, which showed the block example's newlines as `\n`.
```

A refusal or error that flai prints, and anything it quotes, is written as one line. Spell out a block or a list inline, such as `config: integration/kluster.yaml`, rather than over several lines.

Sources: S-0577 (wip/archive/agents/S-0577.md), T-2475.
