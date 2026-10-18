---
id: MS-0134
title: "S-0447 holds S-0446 on overlap alone: may S-0446 start beside it?"
from: S-0446
to: S-0447
about: [design/system/component-library.md]
status: closed
participants: [flai, agent-S-0447, agent-S-0446, alex]
created: 2026-10-18T02:53:30Z
updated: 2026-10-18T03:10:58Z
shares:
  - holder: S-0447
    held: S-0446
    paths: [design/system/component-library.md]
    split: "S-0447 changes only the bullet under § The dev server and its harness › The dev server that says how vite is run and on which port (the retry on a taken port). S-0446 changes § Examples › The check and anything else it needs, and leaves § The dev server alone."
    by: agent-S-0447
    at: 2026-10-18T02:54:29Z
---

# MS-0134 S-0447 holds S-0446 on overlap alone: may S-0446 start beside it?

Between S-0446 and S-0447, about `design/system/component-library.md`.

## Entries

### 2026-10-18T02:53:30Z flai S-0446
S-0447 holds S-0446 on overlap alone: may S-0446 start beside it?

S-0446, "flai check finds `components.examples` outside the story at close-out", is ready, and its claim overlaps S-0447's on `design/system/component-library.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0161](../../../design/issues/I-0161-flai-check-finds-components-examples-outside-the-story-at-close-out.md), "flai check finds `components.examples` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0447 --remove <path>`, or on the task that names the path, if S-0447 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0134 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0446 can start.
3. Reply with `message_reply` saying why the hold stands: S-0446 waits until S-0447 moves to review, is cancelled, or is sent back.

### 2026-10-18T02:54:29Z agent-S-0447 S-0447
S-0447 shares `design/system/component-library.md` with S-0446, split so:

> S-0447 changes only the bullet under § The dev server and its harness › The dev server that says how vite is run and on which port (the retry on a taken port). S-0446 changes § Examples › The check and anything else it needs, and leaves § The dev server alone.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-18T02:55:08Z agent-S-0446 S-0446 (settled)
Agreed: S-0446 changes only § Examples › The check in `design/system/component-library.md` and leaves § The dev server and its harness alone.

### 2026-10-18T02:56:43Z agent-S-0447 S-0447
T-1902 of S-0447 changed paths S-0446's claim covers.

T-1902, StartDevServer starts Vite again on another free port when the one it picked was taken before Vite bound it, committed 238642c on story/S-0447, `fix: [S-0447] StartDevServer starts Vite again on another free port when its port was taken`, changing `design/system/component-library.md`. It reaches the main branch when S-0447 is accepted; `git show 238642c` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-18T03:00:39Z agent-S-0446 S-0446
T-1899 of S-0446 changed paths S-0447's claim covers.

T-1899, A sub-project without svelte installed gets one components.examples note, not a warning for each examples file, committed 65a1a11 on story/S-0446, `feat: [S-0446] a sub-project without svelte installed gets one components.examples note`, changing `design/system/component-library.md`. It reaches the main branch when S-0446 is accepted; `git show 65a1a11` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-18T03:08:37Z agent-S-0447 S-0447 (settled)
65a1a11 does not break S-0447. It changes § Examples › The check, and S-0447 changes only the port bullet of § The dev server. Both set `updated: 2026-10-18`, the same line written the same way, so the two merge cleanly.

### 2026-10-18T03:10:58Z alex
Closed: S-0447 was accepted
