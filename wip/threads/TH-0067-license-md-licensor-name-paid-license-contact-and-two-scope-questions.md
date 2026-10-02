---
id: TH-0067
title: "LICENSE.md: licensor name, paid-license contact, and two scope questions"
anchor:
  path: wip/kanban/stories/S-0231-author-license-md-file.md
  item: S-0231
status: answered
participants: [agent-S-0231, alex, claude]
created: 2026-10-02T12:48:42Z
updated: 2026-10-02T15:40:01Z
---

# TH-0067 LICENSE.md: licensor name, paid-license contact, and two scope questions

On wip/kanban/stories/S-0231-author-license-md-file.md.

## Entries

### 2026-10-02T12:48:42Z agent-S-0231
I am drafting LICENSE.md as an adaptation of the PolyForm Shield License 1.0.0 (the PolyForm Project asks that adaptations not carry the PolyForm name, so it is titled the system-flow Shield License 1.0). Four things only you can settle; I am proceeding with the defaults in brackets and will adjust when you answer.

1. Licensor and copyright line. Who holds the copyright for the `Required Notice:` line and the licensor definition? [Default: "Copyright 2026 Bytepunx" with the repository URL `https://github.com/bytepunx/system-flow`.]
2. Paid licenses. The license says commercial use outside the free classes needs a paid license from the licensor. Should it name a contact (email or URL) for obtaining one? [Default: it says to contact the licensor at the address in the Required Notice, so the repository URL.]
3. Individuals. The goal lists four free classes and none covers an individual using the software personally and noncommercially (a hobbyist, a personal repo). Strictly read, such use also needs a paid license. Keep it strict, or add personal noncommercial use as a fifth free class? [Default: strict, as the goal is written.]
4. Nonprofits. "501(c)(3)" is a US tax status. Should charities with equivalent status under another country's law also qualify? [Default: US 501(c)(3) only, as the goal is written.]

The hosted-service and licensing prohibition applies to everyone, free class or not, and I am writing it so that no license under these terms lifts it; only a separate written agreement with the licensor could.

### 2026-10-02T13:14:33Z agent-S-0231
S-0231 is moving to review with the defaults above in the text: licensor "Copyright 2026 Bytepunx" with the repository URL, the paid-license pointer to that address, no personal noncommercial class, US 501(c)(3) only. One more phrasing to confirm, point 5: the educational class reads "an educational institution funded by a state or national government", my reading of the goal's "State-funded". If you mean US states only, or public institutions anywhere, say which and I will adjust. If any answer differs from a default, send the story back from review and I will change LICENSE.md, its copy under flai/internal/license, the README, and the docs together.

### 2026-10-02T15:35:31Z alex
1 - Copyright 2026 Bytepunx LLC, alex@robsonandmilligan.com. 2 - alex@robsonandmilligan.com 3 - cover individuals using it for commercial work that grosses less than $1M USD annually 4 - for now, keep it to the US tax status.

### 2026-10-02T15:40:01Z claude
Applied on story/S-0231 by the operator's session, since the story's agent had ended: the Required Notice reads "Copyright 2026 Bytepunx LLC, alex@robsonandmilligan.com" with the repository URL; paid licenses are obtained at alex@robsonandmilligan.com; a fifth permitted purpose, Individual use, covers an individual's own work, commercial included, when their gross revenue from work using the software was under USD 1,000,000 in their most recent tax year, with use on behalf of a company left to the company's class; nonprofits stay US 501(c)(3). The embedded copy under flai/internal/license, README.md, and docs/users/flai.md say the same. Point 5 (educational: 'funded by a state or national government') is unanswered and unchanged.
