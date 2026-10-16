---
id: T-1858
type: task
nature: feature
title: An ADR records the dev server on a loopback port framed cross-origin, the components host action, the idle stop, and the page-to-harness protocol
status: backlog
parent: S-0438
owner: alex
created: 2026-10-16T21:33:30Z
updated: 2026-10-16T21:33:30Z
transitions: []
stream: S-0438
tags: [cli, dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/component-library.md]
---
# T-1858 An ADR records the dev server on a loopback port framed cross-origin, the components host action, the idle stop, and the page-to-harness protocol

## Work

- Write and accept the ADR with `flai adr new`, on TH-0483 answers 1 and 5. It records the dev server on 127.0.0.1 and a port flai picks, and the frame being cross-origin to the dashboard. It refines ADR-0029 and ADR-0031: the browser reaches a port flai opens on the host, and the container still holds nothing. It records the `components` host action and `components.idle`, and the protocol's messages and fields, each side checking the other's origin. Say that a proxy through the channel is the later way for a dashboard off the host.
- Add the dev server and the protocol to `design/system/component-library.md`.

Waits for nothing.

## Done when

The ADR is accepted and indexed, the design document has the protocol as a table, and `flai check --strict` passes.

## Notes

`design/adrs` is a folder touch for the new ADR's file, in `claims.shared`.
