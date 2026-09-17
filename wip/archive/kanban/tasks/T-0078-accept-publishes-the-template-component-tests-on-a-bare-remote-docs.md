---
id: T-0078
type: task
nature: feature
title: "accept publishes the template component; tests on a bare remote; docs"
status: done
parent: S-0021
owner: alex
created: 2026-09-17T03:34:08Z
updated: 2026-09-17T03:37:57Z
transitions:
  - to: ready
    at: 2026-09-17T03:37:56Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:37:56Z
    by: agent
  - to: done
    at: 2026-09-17T03:37:57Z
    by: agent
stream: S-0021
tags: [cli, docs]
---

# T-0078 accept publishes the template component; tests on a bare remote; docs

## Work
flai accept publishes template components after tagging (push with --tag at the new version; --no-publish and --no-push opt out; a publish failure is reported without undoing the acceptance); docs/users/flai.md push section; docs/contributors sync procedure replaced by the command; design template.md and flai-cli.md; real dry run against the published repo.

## Done when
Docs match; dry run against bytepunx/system-flow-template reports the expected delta.

## Notes
