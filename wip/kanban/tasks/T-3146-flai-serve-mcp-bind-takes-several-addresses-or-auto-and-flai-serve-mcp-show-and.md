---
id: T-3146
type: task
nature: improvement
title: flai serve mcp bind takes several addresses or auto, and flai serve mcp show and flai serve settings print them
status: backlog
parent: S-0732
owner: alex
created: 2026-11-12T15:54:33Z
updated: 2026-11-12T15:54:53Z
transitions: []
stream: S-0732
tags: [cli]
touches: [flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/cmd/serve_view.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md]
after: [T-3144, T-3145]
---
# T-3146 flai serve mcp bind takes several addresses or auto, and flai serve mcp show and flai serve settings print them

## Work

Let `flai serve mcp bind` in `flai/cmd/serve_mcp.go` take one or more addresses, or `auto`, and save them through `WithMCPBind`. Have `mcpBind` and `mcpBindOf` answer the list, with `auto` resolved through T-3145's lookup. Have `flai serve mcp show`, `flai serve settings` (`flai/cmd/serve_settings.go`), and the settings view (`flai/cmd/serve_view.go`) print every address. Under `--json`, keep `bind` as the first address, for the readers that take a string, and add `binds` with the whole list. When `auto` is set, print `auto` with the addresses it resolves to now. Rewrite the command's `Long` and `Short` help to say what `mcp.bind` holds and what `auto` does, and regenerate `docs/users/flai-reference.md` with `make flai-reference`. Update `TestServeMCPBind` and the settings JSON in `flai/cmd/serve_actions_test.go`.

It waits for T-3144, whose list and `auto` value it sets and prints. It also waits for T-3145, whose lookup resolves `auto` for the host and for what it prints.

## Done when

- `flai serve mcp bind 127.0.0.1 172.17.0.1` saves both, and `flai serve mcp show` prints both.
- `flai serve mcp bind auto` saves `auto` and shows what it resolves to; `flai serve mcp bind auto 127.0.0.1` is refused.
- A configuration with one address prints as before.
- `TestReferenceIsCurrent` passes.
- `flai test` on the changed files passes.

## Notes
