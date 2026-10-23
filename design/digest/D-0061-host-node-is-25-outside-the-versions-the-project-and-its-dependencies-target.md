---
id: D-0061
title: Host Node is 25, outside the versions the project and its dependencies target
kind: cause
sources: [I-0010]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flai/cmd]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0061 Host Node is 25, outside the versions the project and its dependencies target

I-0010 was closed on 2026-10-01.

The cause:

```text
Host Node is 25, outside the versions the project and its dependencies target
```

The fix:

```text
No longer applies: the agent shell's node is v24 (nvm), matching flaiover/.nvmrc; no instance since S-0107 (2026-09-24).
```

Source: I-0010.
