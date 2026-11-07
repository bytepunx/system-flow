---
id: D-0505
title: "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue"
kind: answer
sources: [TH-0500]
topics: []
paths: []
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0505 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt: file this as an issue

The operator, alex, answered TH-0500 at 2026-11-07T09:27:11Z.

What was asked:

```text
**Recommendation:** a small hardening fix. In `sendBackReason`, quote the check name and the reviewer login, for example in a code span, and strip control characters, the way the thread body already does. File it as an issue, or fold it into S-0395, whose agent touches this send-back path next.
```

The answer:

```text
file this as an issue
```

Source: TH-0500.
