---
id: D-0415
title: "pilot keeps nothing of one harness's layout: the adapter returns the files to write in Start.Files, pilot marks its sessions with empty files of its own, and claude makes its own CLAUDE_CONFIG_DIR"
kind: decision
sources: [S-0507]
topics: [harness, docker, go, code]
paths: [flai/internal/pilot/run.go, flai/internal/harness/adapters.go, flai/internal/harness/harness.go, scripts/agent-container.sh]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0415 pilot keeps nothing of one harness's layout: the adapter returns the files to write in Start.Files, pilot marks its sessions with empty files of its own, and claude makes its own CLAUDE_CONFIG_DIR

S-0507 built pilot, the agent images' entrypoint (ADR-0212). Three of its decisions keep everything specific to one inner harness in that harness's adapter, so pilot stays the same for every layer:

- "The claude-code adapter stays pure: with the container fields set, its `Start` names `<cache>/mcp.json` in `--mcp-config`, sets `CLAUDE_CONFIG_DIR` to `<cache>/claude`, and returns the file's content in a new `Start.Files`. pilot writes those files; `flai serve` gets none and is unchanged."
- "pilot records each session it started as an empty file `<cache>/pilot/sessions/<session>`, a harness-neutral mark, rather than reading Claude Code's own session files."
- "pilot does not make `CLAUDE_CONFIG_DIR`: claude makes it when missing, and `scripts/agent-container.sh` links the login into it. Making it in pilot would put claude's layout in pilot."

`design/system/flai-cli.md` § "pilot, the agent images' entrypoint (S-0507)" describes how this works. This entry records why: a new inner harness, such as an OpenCode layer, adds an adapter and a layer. It changes nothing in pilot and does not teach pilot that harness's files.

Source: S-0507 (wip/archive/agents/S-0507.md, Decisions).
