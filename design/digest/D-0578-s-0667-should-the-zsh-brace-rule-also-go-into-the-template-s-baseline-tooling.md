---
id: D-0578
title: "S-0667: should the zsh brace rule also go into the template's baseline tooling convention: Recommended: answer 1. Put the rule in the baseline through a follow-up story, and move S-0667 to review…"
kind: answer
sources: [TH-1047]
topics: [all]
paths: [docs/operators/runbooks/agent-container.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md, design/conventions/tooling.md, design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1053]
---

# D-0578 S-0667: should the zsh brace rule also go into the template's baseline tooling convention: Recommended: answer 1. Put the rule in the baseline through a follow-up story, and move S-0667 to review…

The operator, alex, answered TH-1047 at 2026-11-10T19:28:14Z.

What was asked:

```text
**Recommended: yes, in a follow-up story.** Add "and write `${name}` before a colon" to the baseline zsh rule in `template/root/design/conventions/tooling.md`, then copy it here above the marker. S-0667 can go to review as it is.
```

The answer:

```text
**Recommended:** answer 1. Put the rule in the baseline through a follow-up story, and move S-0667 to review unchanged.

- S-0667 has met both of its criteria. The cause is found and recorded under its Notes › Cause, and I-0276 is closed. The baseline change is not one of its criteria.
- A baseline convention changes in `template/root/design/conventions/` first. It is then copied here above the marker, with a `template/CHANGELOG.md` line. Every project made from the template gets it on upgrade. That is a template release, and it deserves its own story and its own acceptance.
- Answer 2 widens a finished story's scope onto the template. Answer 3 leaves the trap in place for every zsh host that uses the template, and the story shows that quoting alone does not prevent it.
- Until the follow-up lands, the project addition in `design/conventions/tooling.md` covers this repository.

This changes a story's scope, so it is a recommendation for you to confirm, not an answer.

Source: design/conventions/decisions.md § Project additions

Source: design/conventions/decisions.md § Project additions
```

Source: TH-1047.
