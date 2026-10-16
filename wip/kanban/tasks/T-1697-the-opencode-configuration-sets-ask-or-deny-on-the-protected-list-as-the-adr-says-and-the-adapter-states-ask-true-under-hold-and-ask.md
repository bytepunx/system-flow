---
id: T-1697
type: task
nature: feature
title: The opencode configuration sets ask or deny on the protected list as the ADR says, and the adapter states Ask true under hold-and-ask
status: backlog
parent: S-0408
owner: alex
created: 2026-10-16T05:16:31Z
updated: 2026-10-16T05:23:19Z
transitions: []
stream: S-0408
tags: [cli]
touches: [flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness_test.go]
after: [T-1695]
---
# T-1697 The opencode configuration sets ask or deny on the protected list as the ADR says, and the adapter states Ask true under hold-and-ask

## Work

- In `flai/internal/harness/opencodeconfig.go`, write the `permission` rules for the protected list as the ADR says: `ask` under hold-and-ask, or `deny` with `--auto` passed in `opencode.go`.
- Under hold-and-ask, set `Ask: true` in the adapter's `Capabilities`, so `flai serve` starts a story on `opencode` without `auto-approve` or `--deny-protected`; update `TestEachAdapterStatesItsCapabilities`.
- Keep S-0406's `--deny-protected` rules working as before.
- Waits for T-1695: the ADR chooses the rules.

## Done when

The configuration and the capability follow the ADR, and the harness tests pass.

## Notes
