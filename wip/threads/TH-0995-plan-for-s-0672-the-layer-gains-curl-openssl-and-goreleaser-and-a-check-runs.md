---
id: TH-0995
title: "Plan for S-0672: the layer gains curl, openssl, and GoReleaser, and a --check runs smoke inside it; four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0672-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md
  item: S-0672
status: resolved
participants: [planner-S-0672, orchestrator]
created: 2026-11-10T08:48:52Z
updated: 2026-11-10T08:54:03Z
---

# TH-0995 Plan for S-0672: the layer gains curl, openssl, and GoReleaser, and a --check runs smoke inside it; four tasks in three layers

On wip/kanban/stories/S-0672-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md.

## Entries

### 2026-11-10T08:48:52Z planner-S-0672
**Plan.** Four tasks in three layers. Figures written: forecast 1h, delivery 12:45Z, cost of delay 75 USD a week. Nothing needs your answer. The story's Notes give the reasoning under `### Planning`.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2917: the layer adds `curl`, `openssl`, and GoReleaser at `install-tools.sh`'s pin (v2.18.1), and `agent-image.sh` passes `GORELEASER_VERSION` | `flai/images/system-flow/Dockerfile`, `scripts/agent-image.sh` | none |
| 2 | T-2919: `scripts/agent-image.sh --check` runs a tool check and `flai test smoke` inside the built layer. This is the test that reproduces I-0281: it fails on main's layer with "curl is required" | `scripts/agent-image.sh` | T-2917 |
| 3 | T-2922: `docker.md`, `docs/operators/index.md`, `docs/users/flai.md`, and `scripts/README.md` list the new tools and `--check` | those four files | T-2919 |
| 3 | T-2924: close I-0281, citing the passing `--check` | I-0281, `design/issues/summary.md` | T-2919 |

**Figures.**

- **Forecast:** `flai forecast` gave 22m. I raised it to 1h, because the layer is rebuilt and smoke runs inside the container at least twice. S-0630, which changed the same Dockerfile, took 30m.
- **Cost of delay:** `flai cod` gave 75 USD a week, from flai's 30m-per-cycle input. I left it as it is.
- **No `after`:** the story reads nothing another open story writes.

**Assumptions.**

1. **openssl is in scope.** I-0281 names only curl and GoReleaser. `install.sh` also refuses without `openssl`, and `install-test.sh` signs its test releases with it. Alpine's base has no `openssl` binary, so smoke would stop there next.
2. **GoReleaser does not come from `go install` in the layer.** v2.18.1 builds only with Go 1.27, and the layer pins Go 1.26.0 with `GOTOOLCHAIN=local`. T-2917 takes the checksummed release archive, or builds GoReleaser in a stage of its own on a Go 1.27 image. The story's agent records which.
3. **Only you can accept this story.** `flai/images/` and `scripts/agent-image.sh` are on the protected list (T-2157). The story's agent writes them itself, and you accept the story (ADR-0106).
4. **The build and the check need Docker on the host.** If the story's agent runs inside the container, it has no Docker. It will then ask you on a thread to run `scripts/agent-image.sh --layer` and `--check`.

**Not proposed:** a CI job that builds the layer. No workflow builds the agent images today (`docker.md`), and adding one is a story of its own.

### 2026-11-10T08:53:14Z orchestrator
Approved as written: T-2917, then T-2919, then T-2922 and T-2924 together; 1h, 75 USD a week.

- openssl belongs in scope. Criterion 1 asks that the smoke tier's install-test no longer fail inside the container, and `install.sh` stops on a missing openssl right after curl.
- `scripts/agent-image.sh --check` running smoke inside the built layer is the test that reproduces I-0281, and it fails on main's layer.
- GoReleaser from its checksummed release archive, or a Go 1.27 build stage, both keep the layer's Go 1.26.0 pin. Recording the choice in Decisions is right.
- The protected paths and the Docker the build needs are noted, with a thread to you if the agent has no Docker. Nothing here needs your answer before the story runs.

### 2026-11-10T08:54:03Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
