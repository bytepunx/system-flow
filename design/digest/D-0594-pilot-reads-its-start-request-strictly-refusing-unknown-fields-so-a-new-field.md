---
id: D-0594
title: "pilot reads its start request strictly, refusing unknown fields, so a new field on harness.Request breaks every container start on an agent image built before it; key container-only behaviour on Request.Inner instead"
kind: decision
sources: [S-0689, T-3027]
topics: [harness, docker, go, code]
paths: [flai/internal/pilot/request.go, flai/internal/harness/harness.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0594 pilot reads its start request strictly, refusing unknown fields, so a new field on harness.Request breaks every container start on an agent image built before it; key container-only behaviour on Request.Inner instead

pilot, the agent images' entrypoint (ADR-0212), decodes the start request flai serve sends it with `DisallowUnknownFields` (`flai/internal/pilot/request.go`). An image keeps the pilot it was built with. So a field added to `harness.Request` that pilot's request carries reaches an older pilot as an unknown field, and every container start on an image built before the change is refused.

S-0689 met this when it gave the container harness a prompt of its own (T-3027): "the container prompt keys on Request.Inner, which only pilot's request names, in place of the sub-agent's new Container field: pilot's strict reader refuses unknown fields, so a new field would break every container start on an image built before it" (wip/archive/agents/S-0689.md, Decisions). `Request.Inner` is set only for a start inside a container (`inContainer()` in `flai/internal/harness/harness.go`).

So a story that needs to tell a container start from a host start reads `Request.Inner` rather than adding a field to the request.

Sources: S-0689 (T-3027).
