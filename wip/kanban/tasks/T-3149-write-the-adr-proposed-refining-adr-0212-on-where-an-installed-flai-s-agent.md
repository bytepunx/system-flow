---
id: T-3149
type: task
nature: feature
title: Write the ADR, proposed, refining ADR-0212 on where an installed flai's agent images come from and how flai keeps them
status: backlog
parent: S-0733
owner: alex
created: 2026-11-12T15:55:01Z
updated: 2026-11-12T15:55:28Z
transitions: []
stream: S-0733
tags: [cli]
touches: [design/adrs/README.md]
---
# T-3149 Write the ADR, proposed, refining ADR-0212 on where an installed flai's agent images come from and how flai keeps them

## Work

Record the decision from the story's worktree with `flai adr new "<the decision, as a sentence>" --refines 0212 --body-stdin --commit`, status proposed. `--commit` writes the next number, the file, and the index, and adds the new file to the story's touches, so no task names it now.

Read first:

- ADR-0203 § Decision (The image) and § Security concerns (the supply chain);
- ADR-0212 decision 8 and § Consequences;
- `design/system/flai-cli.md` § Verifying a dashboard release before it is run, the signed digest list a pull by digest would reuse;
- `design/system/release-signing.md` § The digest list;
- `design/tech/docker.md` § Agent images;
- E-0036's § What the rebuild of 2026-11-11 showed.

The ADR answers, one section each:

1. Where the bedrock and claude-code images come from for an installed flai: pull by a signed digest, build from Dockerfiles carried in the flai binary, or pull first and build as the fallback. Name the registry, the image names, the signed list, and where `flai` and `pilot` for Linux come from for a build.
2. How an image is tied to its flai (the `dev.system-flow.flai` label or another), and what makes it stale.
3. What a development tree runs, where the images hold the tree's flai.
4. How the project's layer gets its build arguments in any project: a manifest key, the Dockerfile's `ARG` defaults, or both.
5. What `flai serve` does on a missing or stale image: a host job, one at a time, holding the start, with its output in the dashboard; the first build's consent; the limit on a step with no progress; and the retry of a network failure.
6. Whether a strategic agent waits for the layer's build.
7. Where the operator's own build command is named, and the setting that turns the automatic build off. The protected list (`flai/internal/protected`) already holds `scripts/agent-image.sh`.

Recommend in the ADR, for the operator to confirm in T-3152:

- pull first and build as the fallback, which fits ADR-0203 and works on a host with no access to the registry;
- no consent asked, since the build is reported as it runs and the setting turns it off;
- strategic agents start once the harness image is there.

It is the first layer and waits for no task.

## Done when

- A proposed ADR answering points 1 to 7 is committed on `story/S-0733`, refining ADR-0212.
- `design/adrs/README.md` lists it, as `flai adr new` writes it.

## Notes
