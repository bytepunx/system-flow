---
id: T-3190
type: task
nature: feature
title: The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself
status: in-progress
parent: S-0759
owner: alex
created: 2026-11-12T16:05:31Z
updated: 2026-11-13T16:16:18Z
transitions:
  - to: ready
    at: 2026-11-13T16:16:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:16:18Z
    by: system-flow
stream: S-0759
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go]
---
# T-3190 The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself

## Work

ADR-0295 point 7 keeps both settings per project in the host's configuration, under `agent.projects.<root>` (ADR-0261), not on a harness entry. So, in `flai/internal/config/config.go`:

- `ProjectAgent` gains `Image ProjectImage` (`json:"image,omitzero"` or a pointer, so that an entry without it round-trips unchanged).
- `ProjectImage` holds `Command []string` (`json:"command,omitempty"`), the operator's build command as an argument list, never a shell string, and `Auto *bool` (`json:"auto,omitempty"`), nil meaning on.
- `ProjectAgent.IsZero` counts them, so an entry holding only an image setting is kept and one cleared of it is removed.
- `checkProjects` refuses an image command whose program is empty or any word holding a NUL, naming the key `agent.projects["<root>"].image.command` and the command that sets it, `flai serve agent image command --project <root> -- <program> [args...]`.
- Readers on `AgentStart`: `ImageCommand(root) []string` (a copy, nil when none) and `AutoImage(root) bool` (true unless set false). `ForProject` need not carry them.
- Setters are `WithProject` with the entry changed, as the CLI's other per-project writes do; add a helper only if it makes T-3191's code plainer.

Test set, clear, refused (empty program, NUL), absent, and that a host configuration naming neither loads and round-trips byte for byte as before.

It is the first layer: T-3191, T-3192, and T-3194 read these fields.

## Done when

- `flai test flai/internal/config` passes, with set, clear, refused, and absent.
- An existing host configuration round-trips unchanged.

## Notes
