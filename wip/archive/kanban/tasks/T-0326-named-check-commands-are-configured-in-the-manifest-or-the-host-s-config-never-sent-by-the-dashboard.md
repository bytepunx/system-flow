---
id: T-0326
type: task
nature: feature
title: Named check commands are configured in the manifest or the host's config, never sent by the dashboard
status: done
parent: S-0082
owner: alex
created: 2026-09-22T22:36:14Z
updated: 2026-09-22T22:44:31Z
transitions:
  - to: ready
    at: 2026-09-22T22:37:23Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T22:37:24Z
    by: system-flow
  - to: review
    at: 2026-09-22T22:44:31Z
    by: system-flow
  - to: done
    at: 2026-09-22T22:44:31Z
    by: system-flow
stream: S-0082
tags: []
---
# T-0326 Named check commands are configured in the manifest or the host's config, never sent by the dashboard

## Work
Add `Checks []NamedCommand` to `manifest.Manifest` (`flai/internal/manifest/manifest.go`), each `{name, command []string}` (yaml and json tags, matching `Project`'s double-duty). Add the same shape to `internal/config.Config` under a new `Checks` section: `Config.Checks.Commands []NamedCommand` and `Config.Checks.TimeoutMinutes int` (default 15 when zero) — neither is among `Keys()`, matching `Agent`. Write a small resolver (`internal/serve` is the natural home, beside the other host-side state): the host config's commands, when non-empty, are used; otherwise the manifest's. This is the operator's real choice of where to name them, not a merge of both.

`flai serve checks set --name <n> -- <argv...>`, `flai serve checks show`, `flai serve checks clear [<name>]` manage the host-config side, mirroring `flai serve agent set/show/clear` exactly (`{story}` and `{root}` substituted in an argument, nothing else interpreted, run never through a shell — reuse or mirror the `strings.NewReplacer` `agents.go` already does; this is the second use of that exact substitution, so pull it into one small shared helper both call, in whichever of `internal/serve` or `internal/manifest` doesn't create an import cycle). `flai serve enable checks` / `disable` need nothing new: `ActionChecks = "checks"` registered in `hostapi.Actions` is enough for the existing generic mechanism.

## Done when
`go test -race -short ./...` covers: the manifest parses and round-trips `checks:`; config's `checks` section round-trips and stays out of `Keys()`; `flai serve checks set/show/clear` (including the `{story}`/`{root}` substitution and quoting in `show`'s output); the resolver picks config over manifest when config is non-empty, and falls back correctly. `golangci-lint run ./...` clean.

## Notes
No manifest version bump: this is an additive, optional field, same as `Dashboard` was.
