---
id: D-0390
title: Agents never read the operator's dashboard login token, and scripts/flaiover-dev.sh prints a login link carrying it, so an agent's dev server uses a throwaway token file
kind: decision
sources: [S-0048, S-0489]
topics: [flaiover, dashboard, host]
paths: [scripts/flaiover-dev.sh, design/system/devex.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0390 Agents never read the operator's dashboard login token, and scripts/flaiover-dev.sh prints a login link carrying it, so an agent's dev server uses a throwaway token file

S-0048's first dev server was started with the repository's dev script, which wrote a login link with the real token into its log. The agent deleted the log unread and restarted the server with a throwaway token file, so the token was never read or shown. `scripts/flaiover-dev.sh` still echoes `log in with:` and the login URL from `flai dashboard token --json`, so an agent does not read that script's output while it runs on the operator's token.

S-0489 held to the same line: "A real reproduction needs the operator's login token, which I do not read." It reproduced the fault in a test (T-2083) and left the live check to the operator.

Where the token lives is in `docs/operators/settings.md`; this entry records only that agents do not read it.

Sources: S-0048, S-0489.
