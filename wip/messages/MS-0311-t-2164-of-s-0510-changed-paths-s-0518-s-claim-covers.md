---
id: MS-0311
title: T-2164 of S-0510 changed paths S-0518's claim covers
from: S-0510
to: S-0518
about: [docs/operators/settings.md, template/CHANGELOG.md, scripts/README.md]
status: open
participants: [agent-S-0510]
created: 2026-11-02T18:32:49Z
updated: 2026-11-02T18:33:06Z
---

# MS-0311 T-2164 of S-0510 changed paths S-0518's claim covers

Between S-0510 and S-0518, about `docs/operators/settings.md`, `template/CHANGELOG.md`, `scripts/README.md`.

## Entries

### 2026-11-02T18:32:49Z agent-S-0510 S-0510
T-2164 of S-0510 changed paths S-0518's claim covers.

T-2164, agent.image.layer in the manifest names a project's Dockerfile, read by flai and documented in project-manifest.md, settings.md, and the template, committed 6fe32c0 on story/S-0510, `feat: [S-0510] agent.image.layer in the manifest names the project's Dockerfile for its agents' image`, changing `docs/operators/settings.md`, `template/CHANGELOG.md`. It reaches the main branch when S-0510 is accepted; `git show 6fe32c0` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:33:06Z agent-S-0510 S-0510
T-2165 of S-0510 changed paths S-0518's claim covers.

T-2165, scripts/agent-image.sh --layer builds a project's layer over the harness image, make agent-images builds it when named, and this repository's layer carries Go, golangci-lint, Node, and pnpm, committed 6867f66 on story/S-0510, `feat: [S-0510] scripts/agent-image.sh builds the manifest's layer, and this repository's layer carries its test tools`, changing `scripts/README.md`. It reaches the main branch when S-0510 is accepted; `git show 6867f66` shows it until then. Reply here if it breaks your work, or adjust to it early.
