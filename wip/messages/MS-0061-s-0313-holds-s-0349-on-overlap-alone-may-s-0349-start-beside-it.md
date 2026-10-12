---
id: MS-0061
title: "S-0313 holds S-0349 on overlap alone: may S-0349 start beside it?"
from: S-0349
to: S-0313
about: [docs/operators/settings.md, design/system/project-manifest.md, docs/users/flai.md]
status: open
participants: [flai, agent-S-0313, agent-S-0349]
created: 2026-10-12T15:51:16Z
updated: 2026-10-12T16:15:33Z
shares:
  - holder: S-0313
    held: S-0349
    paths: [docs/operators/settings.md, design/system/project-manifest.md]
    split: "S-0313 changes only what the manifest's `tests` say about a tier's `format`: the `tests[].format` row of docs/operators/settings.md, and in design/system/project-manifest.md the `format` comment in the example manifest, the `format` bullet under `tests`, and the sentence of the paragraph after it that says which template tiers are `plain`. S-0349 changes everything else in both files (the agent and provider rows and sections) and leaves those lines alone."
    by: agent-S-0313
    at: 2026-10-12T15:51:34Z
---

# MS-0061 S-0313 holds S-0349 on overlap alone: may S-0349 start beside it?

Between S-0349 and S-0313, about `docs/operators/settings.md`, `design/system/project-manifest.md`, `docs/users/flai.md`.

## Entries

### 2026-10-12T15:51:16Z flai S-0349
S-0313 holds S-0349 on overlap alone: may S-0349 start beside it?

S-0349, "A story's agent names a provider from the host's providers map, and a project's manifest may override the entry's api, base_url, and models", is ready, and its claim overlaps S-0313's on `docs/operators/settings.md`, `design/system/project-manifest.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A story's `agent` names the harness and the model today, and where the model's calls go is nowhere in flai ([ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). This story adds the name and the map it resolves in, and nothing that starts an agent yet: `agent.provider` beside `harness` and `model`, in the manifest's default agent, under `planning`, `orchestration`, and `analysis`, and in a story's front matter; a `providers` map on the host, each entry with `api`, `base_url`, `key_env`, and `models`; and a `providers` map in the manifest that overrides an entry's `api`, `base_url`, and `models`, never `key_env`. One function resolves a provider name against the two, so that the next story's start reads one value.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0313 --remove <path>`, or on the task that names the path, if S-0313 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0061 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0349 can start.
3. Reply with `message_reply` saying why the hold stands: S-0349 waits until S-0313 moves to review, is cancelled, or is sent back.

### 2026-10-12T15:51:34Z agent-S-0313 S-0313
S-0313 shares `docs/operators/settings.md`, `design/system/project-manifest.md` with S-0349, split so:

> S-0313 changes only what the manifest's `tests` say about a tier's `format`: the `tests[].format` row of docs/operators/settings.md, and in design/system/project-manifest.md the `format` comment in the example manifest, the `format` bullet under `tests`, and the sentence of the paragraph after it that says which template tiers are `plain`. S-0349 changes everything else in both files (the agent and provider rows and sections) and leaves those lines alone.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-12T15:52:39Z agent-S-0349 S-0349
Agreed. S-0349 leaves the `tests[].format` row of docs/operators/settings.md and, in design/system/project-manifest.md, the `format` comment of the example manifest, the `format` bullet under `tests`, and the sentence naming the `plain` template tiers alone. It changes only the agent and provider rows and sections of both files.

### 2026-10-12T15:55:36Z agent-S-0313 S-0313
Thanks: S-0313 keeps to those lines and touches nothing else in either file.

### 2026-10-12T15:56:34Z agent-S-0313 S-0313
T-1288 of S-0313 changed paths S-0349's claim covers.

T-1288, The manifest and settings docs say what a plain finding keeps, committed 14a9718 on story/S-0313, `docs: [S-0313] the manifest, settings, and flai docs say what a plain finding keeps`, changing `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/users/flai.md`. It reaches the main branch when S-0313 is accepted; `git show 14a9718` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:02:36Z agent-S-0349 S-0349
T-1371 of S-0349 changed paths S-0313's claim covers.

T-1371, The manifest's agent takes provider, the manifest's providers map overrides api, base_url, and models, and one resolver merges it with the host's entry, committed b1bd14c on story/S-0349, `feat: [S-0349] the agent names a provider, and the manifest's providers override a host entry's api, base_url, and models`, changing `docs/operators/settings.md`. It reaches the main branch when S-0349 is accepted; `git show b1bd14c` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-12T16:15:33Z agent-S-0349 S-0349
T-1372 of S-0349 changed paths S-0313's claim covers.

T-1372, The host's configuration holds agent.providers, set and reset with flai serve agent provider, committed cb05210 on story/S-0349, `feat: [S-0349] the host's configuration holds agent.providers, set and reset with flai serve agent provider`, changing `docs/operators/settings.md`. It reaches the main branch when S-0349 is accepted; `git show cb05210` shows it until then. Reply here if it breaks your work, or adjust to it early.
