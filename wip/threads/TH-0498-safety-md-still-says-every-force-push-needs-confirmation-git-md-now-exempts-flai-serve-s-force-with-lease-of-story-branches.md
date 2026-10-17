---
id: TH-0498
title: "safety.md still says every force push needs confirmation; git.md now exempts flai serve's force-with-lease of story branches"
anchor:
  path: design/conventions/safety.md
  heading: Rules
status: open
participants: [orchestrator]
created: 2026-10-17T16:04:06Z
updated: 2026-10-17T16:04:06Z
---

# TH-0498 safety.md still says every force push needs confirmation; git.md now exempts flai serve's force-with-lease of story branches

On design/conventions/safety.md § Rules.

## Entries

### 2026-10-17T16:04:06Z orchestrator
**Recommendation:** a small follow-up story that adds the exception to the force-push line in `safety.md` § Rules, in both copies:
- `design/conventions/safety.md`
- `template/root/design/conventions/safety.md`

The line would read something like: "force pushes always do, except the force-with-lease push of a story's own branch that `flai serve` makes under pull-request acceptance (git.md, ADR-0166)."

**Why it matters:**
- **The conventions now disagree.** S-0393 put the exception into both copies of `git.md`, and you accepted it, so I take the exception as confirmed. It went out in template 1.0.97. `safety.md` still says force pushes "always" need confirmation, "per action and per context".
- **An agent reading both meets a contradiction.** It may refuse the push, or it may ask you every time.

**What I need:** say "draft it" and I will put it to the planner once there is an epic or a story to plan it under. Or say "leave it" if the wording in `git.md` and ADR-0166 is enough.

I raised this once on TH-0496 before the acceptance, and I will not raise it again.
