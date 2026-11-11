---
id: D-0477
title: "A story worktree's bin/flai also runs in the Alpine agent image, so scripts/flai-build.sh builds it with CGO_ENABLED=0, scripts/env.sh keeps tool caches in the home when .flai-cache is not writable, and the layer sets GOFLAGS=-buildvcs=false"
kind: decision
sources: [S-0510, I-0255, S-0520, S-0686]
topics: [docker, harness, go, code]
paths: [scripts/flai-build.sh, scripts/env.sh, flai/images/system-flow/Dockerfile, scripts/flai.sh, scripts/agent-container.sh, scripts/agent-vcs-test.sh]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0477 A story worktree's bin/flai also runs in the Alpine agent image, so scripts/flai-build.sh builds it with CGO_ENABLED=0, scripts/env.sh keeps tool caches in the home when .flai-cache is not writable, and the layer sets GOFLAGS=-buildvcs=false

S-0510's hand run of `scripts/flai.sh test` inside `flai-agent-sf-claude-code:local`, on its own worktree, hit three things. It fixed each one. From S-0510's Decisions (`wip/archive/agents/S-0510.md`):

```text
The in-container run met three things, each fixed here. `scripts/flai-build.sh` now builds `bin/flai` with `CGO_ENABLED=0`, because the worktree's binary is shared with the musl image. `scripts/env.sh` keeps the tool caches in the home when the main checkout's `.flai-cache` is not writable. The layer sets `GOFLAGS=-buildvcs=false` (I-0255).
```

- A dynamically linked `bin/flai` built on the host does not run on the image's musl libc. The worktree holds one binary for both, so it is built static.
- Inside the container the main checkout is mounted so that its `.flai-cache` may not be writable, so `env.sh` falls back to the home for the Go and pnpm caches.
- The third fix no longer holds, though the title still names it. S-0686 (2026-11-11) found I-0255's cause and dropped the layer's `GOFLAGS=-buildvcs=false`, so `go build` in the layer stamps VCS information as it does on the host. The 128 was git refusing the root-owned main checkout for dubious ownership. Go passes over a worktree's `.git` file and runs git in the main checkout, and the `safe.directory` the container harness gives since ADR-0253 admits it. `scripts/agent-vcs-test.sh`, in the smoke tier, fails without it. `design/tech/docker.md` and `docs/operators/runbooks/agent-container.md` now record this. Do not turn stamping off again; give the container the `safe.directory`.

Keep the first two in place when changing these scripts or the layer. The same split between host and image hits `flaiover/node_modules`, whose native bindings are libc-specific (I-0259).

S-0520 used the shared binary to check a branch's flai inside a container without rebuilding the images. From its Decisions (`wip/archive/agents/S-0520.md`):

```text
The hand check in T-2236 builds `bin/flai` in the worktree, with `CGO_ENABLED=0` (S-0510), rather than rebuilding the images. `scripts/flai.sh` inside the image, which has no `go`, runs that binary, so the check runs this branch's flai.
```

```text
The close-out was checked inside `flai-agent-claude-code:local`, as T-2236 says, not inside the layer image, which `scripts/agent-container.sh` does not run. A whole close-out inside the layer is S-0565's criterion 3.
```

So a hand check through `scripts/agent-container.sh` runs the harness image, without the project's layer and its `go`; a story's own flai reaches it only as the worktree's static `bin/flai`.

Sources: S-0510, I-0255, S-0520, S-0686.
