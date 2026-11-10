---
id: TH-1047
title: "S-0667: should the zsh brace rule also go into the template's baseline tooling convention?"
anchor:
  path: wip/kanban/stories/S-0667-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md
  item: S-0667
status: open
participants: [agent-S-0667, orchestrator]
created: 2026-11-10T11:05:39Z
updated: 2026-11-10T11:08:58Z
---

# TH-1047 S-0667: should the zsh brace rule also go into the template's baseline tooling convention?

On wip/kanban/stories/S-0667-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md.

## Entries

### 2026-11-10T11:05:39Z agent-S-0667
**Recommended: yes, in a follow-up story.** Add "and write `${name}` before a colon" to the baseline zsh rule in `template/root/design/conventions/tooling.md`, then copy it here above the marker. S-0667 can go to review as it is.

**Why:** I-0276's empty mount was zsh, not Docker. In zsh, `"$ROOT:$ROOT:ro"` expands to `<ROOT>:<ROOT>o`, because `$ROOT:r` is a history modifier even inside double quotes. The baseline rule says "Quote every … variable", and quoting does not prevent this. Every project made from the template whose host shell is zsh can hit it. S-0667's Notes, under Cause, hold the commands and their output.

**What S-0667 did:** it added the rule only as a project addition below the marker in `design/conventions/tooling.md` (T-3078), since a story's agent does not edit a baseline rule. It also added a note to the runbook (T-2929) and closed I-0276. The close-out passed at c29bea5.

**Answers:**

1. Yes, in a follow-up story (recommended). I answer the narrative's open question with that, and move S-0667 to review unchanged. The operator makes the story with `flai story new`, or tells me to.
2. Yes, in S-0667. I add the sentence to the template's baseline and to `design/conventions/tooling.md` above the marker, with a `template/CHANGELOG.md` line, close out again, and move S-0667 to review.
3. No. The project addition is enough. I record that and move S-0667 to review.

### 2026-11-10T11:08:58Z orchestrator (recommendation)
**Recommended:** answer 1. Put the rule in the baseline through a follow-up story, and move S-0667 to review unchanged.

- S-0667 has met both of its criteria. The cause is found and recorded under its Notes › Cause, and I-0276 is closed. The baseline change is not one of its criteria.
- A baseline convention changes in `template/root/design/conventions/` first. It is then copied here above the marker, with a `template/CHANGELOG.md` line. Every project made from the template gets it on upgrade. That is a template release, and it deserves its own story and its own acceptance.
- Answer 2 widens a finished story's scope onto the template. Answer 3 leaves the trap in place for every zsh host that uses the template, and the story shows that quoting alone does not prevent it.
- Until the follow-up lands, the project addition in `design/conventions/tooling.md` covers this repository.

This changes a story's scope, so it is a recommendation for you to confirm, not an answer.

Source: design/conventions/decisions.md § Project additions

Source: design/conventions/decisions.md § Project additions
