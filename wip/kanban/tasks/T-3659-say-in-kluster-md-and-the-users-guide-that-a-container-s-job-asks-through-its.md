---
id: T-3659
type: task
nature: remediation
title: Say in kluster.md and the users guide that a container's --job asks through its cache folder, and close I-0405
status: backlog
parent: S-0849
owner: alex
created: 2026-11-28T00:05:18Z
updated: 2026-11-28T00:05:18Z
transitions: []
stream: S-0849
tags: [docs, kluster]
touches: [design/system/kluster.md, docs/users/flai.md, design/issues/I-0405-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md, design/issues/summary.md]
after: [T-3657, T-3658]
---
# T-3659 Say in kluster.md and the users guide that a container's --job asks through its cache folder, and close I-0405

## Work

Update the documents that say a container's `--job` does not reach flai serve, then close the issue.

- `design/system/kluster.md` § "The requests (S-0660, S-0661)": the request file goes in the main checkout's `requests/` on the host, and in `.flai-cache/agents/S-nnnn/kluster-requests/` in a story's container. flai serve takes from both, refuses a request for another story in a story's folder, and writes the answer and the job state where the asker reads them. Replace the sentence that names I-0405 and says the container's agent asks with the MCP tools. Link the new ADR.
- `design/system/kluster.md` § "The agent's copy (S-0834)": say the requests folder is the cache folder's second kluster use, beside the copy and outside it, so a teardown's removal of the copy leaves it.
- `docs/users/flai.md`, the "In a container" item under "A story's agent works in its namespace", and the `--job` section: the three commands work from the container as on the host. Drop the line that they do not reach flai serve.
- Bump `updated` on both documents.
- Close the issue with `flai issue close I-0405 --reason "<what fixed it>"`, naming the ADR, the folder, and S-0849. Let flai regenerate `design/issues/summary.md`.

It waits for T-3657 and T-3658, since the documents describe what those two built.

## Done when

- Neither document says a container's `--job` fails to reach flai serve, and both say where its request goes.
- I-0405's status is closed, with a reason naming the fix, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` and the markdown lint pass on the four paths.

## Notes
