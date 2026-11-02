---
id: T-2156
type: task
nature: feature
title: flai serve agent harness container sets the image per inner harness, the network, and the login path on the host, shown by flai serve settings and indexed in settings.md
status: done
parent: S-0508
owner: alex
created: 2026-10-26T04:48:11Z
updated: 2026-11-02T07:10:02Z
transitions:
  - to: ready
    at: 2026-11-02T06:53:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T06:53:43Z
    by: system-flow
  - to: done
    at: 2026-11-02T07:10:02Z
    by: agent-S-0508
stream: S-0508
tags: [cli, docker]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md, flai/cmd/serve_settings.go]
usage:
  source: log
  seconds: 979
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 26852
      cache_read: 6131430
      cache_write: 203770
      cost: 3.135
---
# T-2156 flai serve agent harness container sets the image per inner harness, the network, and the login path on the host, shown by flai serve settings and indexed in settings.md

## Work

Extend `HarnessHost` in `flai/internal/config/config.go` with the container harness's settings: `Images` (a map from inner harness to image reference, `claude-code` to `flai-agent-claude-code:local` by default), `Network` (empty for Docker's default bridge), and `Login` (the folder holding the operator's Claude login, mounted read-only). Add the flags `--image <inner>=<ref>` (repeatable), `--network`, and `--login` to `newServeAgentHarnessCmd` in `cmd/serve_actions.go` for `container`, refused for another harness, and show them in `flai serve settings` and the agent settings view. Add `container` to `settable()`. Give each new key its row in `docs/operators/settings.md` so `settings_doc_test` passes, and regenerate `docs/users/flai-reference.md` with `make flai-reference`. This task waits for nothing: it changes no file the adapter task changes.

## Done when

`flai serve agent harness container --image claude-code=flai-agent-claude-code:local --login ~/.claude` is kept in the host configuration and shown by `flai serve settings`, `settings_doc_test` passes, and `flai test flai/internal/config flai/cmd` passes.

## Notes

`checkHarnesses` in `config.go` validates the harness hosts; make it refuse an `--image` for an inner harness `harness.Names()` lacks.
