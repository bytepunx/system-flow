---
id: T-3190
type: task
nature: feature
title: The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:31Z
updated: 2026-11-12T16:05:31Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go]
---
# T-3190 The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself

## Work

Add two fields to `HarnessHost` in `flai/internal/config/config.go` (lines 444 to 473), kept per project as the host's agent settings are (ADR-0261):

- `build`: an argument list, never a shell string;
- `auto_image`: a boolean, on when absent.

Validate them in the `checkHarnesses` area (lines 564 to 590): `build` only on `container`, a non-empty first word, no NUL. Add the setters, modelled on `WithMCPBind` (lines 110 to 120). Test both, with a configuration that names neither read as before.

It is the first layer.

## Done when

- `go test ./internal/config/` passes, with set, clear, refused on another harness, and absent.
- An existing host configuration round-trips unchanged.

## Notes
