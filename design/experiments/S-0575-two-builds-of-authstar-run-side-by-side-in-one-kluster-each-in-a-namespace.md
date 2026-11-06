---
title: Two builds of authstar side by side in one kluster, a namespace per story
updated: 2026-11-06
status: active
story: S-0575
---

# Two builds of authstar side by side in one kluster, a namespace per story

E-0032 wants a story in review deployed to a local kluster for its integration or smoke tier. This experiment tried, by hand on authstar, the shape its outcome guesses at: one shared cluster, one namespace per story, and each story's own settings over the shared signet, RabbitMQ, Postgres, and ClickHouse.

## Hypothesis

One kluster brought up with authstar's `authstar` profile can host several builds of authstar at once, each in a namespace named for a story (`s-0901`, `s-0902`). The shared services need only per-namespace settings, not a cluster or an install each:

- a RabbitMQ vhost
- Postgres and ClickHouse databases with the namespace's prefix
- signet bundles under a signet namespace equal to the Kubernetes namespace, since SPIRE names each workload `spiffe://<trust domain>/ns/<namespace>/sa/<service>` and signet grants a workload its own namespace and service by convention

A story's deploy then costs a build, an image load, and a rollout, minutes rather than the cluster's bring-up. Deleting a story's namespace leaves the others and the shared services working.

## Success measure

- One `kluster up` with authstar's `integration/kluster.yaml` brings up signet, RabbitMQ, Postgres, and ClickHouse answering, and its time is recorded.
- Two builds from two commits of `~/git/bytepunx/authstar` run at once in `s-0901` and `s-0902`. Each has its own vhost, databases, and signet bundles, and authstar's golden-path scenario runs against each, passing or failing at a named step for a named reason.
- Deleting `s-0901` leaves `s-0902` and the shared services answering, and build A is then promoted to the default `authstar` namespace.
- The time of each build, load, deploy, and teardown is recorded, and the per-story cost is small beside the bring-up's.
- Nothing is committed in the authstar, signet, or kluster repositories.

## What was done

Host: `squatchship`, WSL2 on Linux 6.18, 24 cores, 62 GB. kluster 1.3.1, signet CLI 0.10.3. Another kluster cluster, `marketeer`, was running throughout and was left alone.

Scratch and copies:

- `/tmp/S-0575` holds the logs, the kubeconfig, and the two exported source trees.
- The story worktree's git-ignored `bin/s-0575/` holds what the experiment wrote: the scripts below, a copy of authstar's `integration/scripts`, and a copy of `kluster-lib`.
- Build A is authstar `e02e430`, main today. Build B is `af55c04`, main before S-0009 changed chronicle's token check. Each was exported with `git archive <commit> | tar -x -C /tmp/S-0575/src-<A|B>`, so the authstar checkout was only read.
- The host has no `nc`, which authstar's scripts call as `nc -z` to wait for a port-forward. A shim in `bin/s-0575/path/nc` opens `/dev/tcp/<host>/<port>` instead.
- The kubeconfig was written with `kluster kubeconfig --name authstar-integration --output /tmp/S-0575/kubeconfig` and used through `KUBECONFIG`, never merged into `~/.kube/config`.

### Bring-up

```sh
cd /tmp/S-0575
cp ~/git/bytepunx/authstar/integration/kluster.yaml .   # name authstar-integration, profile authstar
kluster --config ./kluster.yaml up
```

authstar's `integration/scripts/provision.sh` was not run. It adds a deploy key to the GitHub repository `bytepunx/authstar-integration` with `gh`, and commits and pushes secrets there for signet to sync. kluster 1.3.1's `authstar` profile no longer needs it: `kluster-lib/profile/authstar_services.go` seeds signet through its admin API, creates the databases, and deploys all eight services into the `authstar` namespace itself.

The shared services were checked like this:

```sh
kubectl -n postgres exec postgresql-0 -- env PGPASSWORD=postgres psql -U postgres -tAc "select datname from pg_database"
kubectl -n rabbitmq exec rabbitmq-0 -- rabbitmqctl list_vhosts -q
kubectl -n clickhouse exec clickhouse-shard0-0 -- clickhouse-client --password clickhouse -q "show databases"
kubectl -n signet port-forward deploy/signet 18444:8444 &
signet --server localhost:18444 status --token "$(kubectl create token signet-admin -n signet --audience signet --duration=1h)"
```

A baseline golden-path run against kluster's own `authstar` namespace followed. It used the operator tokens kluster writes to `~/.config/kluster/authstar-integration-authstar-tokens.yaml`, passed in as `TOWER_OPERATOR_TOKEN` and `HERALD_OPERATOR_TOKEN`.

### Builds

`bin/s-0575/build.sh` builds each service's own `Dockerfile` from the exported tree, tagged with the namespace it is for:

```sh
bash bin/s-0575/build.sh /tmp/S-0575/src-A s-0901 /tmp/S-0575/logs   # docker.io/bytepunx/authstar-<service>:s-0901
bash bin/s-0575/build.sh /tmp/S-0575/src-B s-0902 /tmp/S-0575/logs   # docker.io/bytepunx/authstar-<service>:s-0902
```

### Story namespaces

`bin/s-0575/ns-up.sh <namespace> <tag> <source> <logs>` makes one story namespace in six timed steps. The scripts are written out under [Scripts](#scripts).

1. **Load.** `docker save` each of the eight images and `ctr -n k8s.io images import` it into the k3d node, as authstar's `deploy-local.sh` does. kluster has no image import.
2. **Namespace.** The commit's `manifests/00-namespace.yaml`, renamed: the namespace and its eight service accounts. SPIRE's `ClusterSPIFFEID` from kluster's `spire` profile selects every namespace but six system ones, so the pods get `spiffe://authstar.cluster.local/ns/<namespace>/sa/<service>` with nothing more.
3. **RabbitMQ.** `rabbitmqctl add_vhost <namespace>`, `add_user <namespace> <random password>`, and `set_permissions` on that vhost alone.
4. **Databases.** Postgres `CREATE DATABASE <prefix>tower`, `<prefix>keep`, and `<prefix>herald`, and ClickHouse `CREATE DATABASE <prefix>chronicle`, where `s-0901`'s prefix is `s0901_`. The commit's `01-db-init-job.yaml` is not used: it creates the fixed names.
5. **Signet.** `bin/s-0575/spec.sh` builds each service's bundle from the commit's own `config/authstar/<service>.yaml` and `scripts/lib-secrets.sh`, with the namespace's values put in. `storens`, a scratch command over kluster-lib's exported `addon.SeedSignetBundle` and `addon.EnsureSignetPolicy`, pushes them through signet's admin API under signet namespace `<namespace>`, with the four cross-service grants authstar's `provision.sh` makes: herald puts into tower, keep, and portcullis, and tower into portcullis.
6. **Deploy.** The commit's `manifests/1[0-7]-*.yaml`, renamed, then `kubectl rollout status` for each of the eight Deployments.

Every manifest passes through one `sed` rename:

| Rule | Why |
|------|-----|
| `namespace: authstar`, `name: authstar`, and `value: authstar` become the namespace | The objects' namespace and each service's `SIGNET_NAMESPACE` |
| `.authstar.svc.cluster.local` becomes `.<namespace>.svc.cluster.local` | steward's and web's upstreams |
| `image: bytepunx/authstar-<service>:latest` becomes `docker.io/bytepunx/authstar-<service>:<tag>`, and `imagePullPolicy: Always` becomes `IfNotPresent` | The loaded build, never a pull |
| `signet.signet.svc.cluster.local:8444` becomes `:8443` | A drift in authstar, below |

The golden path ran from a copy of the commit's `integration/scripts` in `bin/s-0575/scripts`, with two changes made through variables:

- `lib.sh`'s `AUTHSTAR_NAMESPACE="authstar"` became `"${AUTHSTAR_NAMESPACE:-authstar}"`.
- `scenario-golden-path.sh`'s `psql -d tower` and `-d herald` became `-d "${DB_PREFIX:-}tower"` and `"${DB_PREFIX:-}herald"`.

`bin/s-0575/golden.sh <namespace> <tokens> <log>` sets both variables and the namespace's operator tokens, runs the scenario, and kills the port-forwards it leaves behind.

### Teardown and promotion

`bin/s-0575/ns-down.sh s-0901 <logs>` deletes the namespace and waits for it to go, deletes its RabbitMQ vhost and user, and drops its three Postgres databases and its ClickHouse database. It leaves the namespace's signet bundles and policies, for want of a delete (below). The golden path was then run against `s-0902` again.

`bin/s-0575/promote.sh s-0901 /tmp/S-0575/src-A <logs>` promotes build A as an acceptance would. It deletes the Deployments and CronJobs kluster made in `authstar`, since `kubectl apply` would only merge into those, then runs `ns-up.sh authstar s-0901 …`. In `authstar` the database prefix is empty, so it reuses kluster's `tower`, `keep`, `herald`, and `chronicle` databases and the `authstar` vhost, with a new RabbitMQ user `authstar`. The golden path then ran against `authstar` and against `s-0902`.

## Results

### Bring-up

`kluster up` took 691 seconds, 11m31s, from 22:02:19Z to 22:13:50Z. The Go build of the seeding tool (below) ran on about five cores for two of those minutes, so an idle host may be a little faster.

| Step | Time |
|------|------|
| k3d cluster | 23s |
| cert-manager | 32s |
| SPIRE | 53s |
| traefik-tls | 5s |
| signet, unpinned, the latest chart | 2m50s |
| RabbitMQ | 2m01s |
| RabbitMQ topology operator | 25s |
| Dex | 12s |
| Postgres | 56s |
| ClickHouse | 1m04s |
| Configure SPIRE, signet | 25s, under 1s |
| Configure authstar: vhost, provisioner Job, databases, signet bundles, eight Deployments | 1m44s |

At the end every shared service answered:

| Service | Namespace | What answered |
|---------|-----------|---------------|
| Postgres | `postgres` | Databases `tower`, `keep`, `herald` |
| RabbitMQ | `rabbitmq` | Vhosts `/` and `authstar`; users `rmq-…` (administrator), `tower`, `keep` |
| ClickHouse | `clickhouse` | Database `chronicle` |
| signet | `signet` | `State: UNSEALED` |

All eight authstar services were running in `authstar`, from Docker Hub's `bytepunx/authstar-<service>:latest`.

The baseline golden path failed at its first step, in 8 seconds, before any story namespace existed:

```text
FAIL: tenant provisioning (HTTP 500): 14 UNAVAILABLE: No connection established.
Last error: SSL routines:tls_validate_record_header:wrong version number
```

Tower answered `POST /tenants` with a TLS call to a port that answers in plain text. So the `latest` images kluster pulls do not match the signet chart it installs, unpinned, today. kluster's own deployment of authstar is not a known-good baseline. Each story namespace therefore deploys a local build. The same failure came back with build A in `s-0901`, so its cause is not the images; it is the first drift listed under [What it took](#what-it-took).

### Builds

| Build | Commit | Total | Per service |
|-------|--------|-------|-------------|
| A, tag `s-0901` | `e02e430` | 1097s | portcullis 793s (Rust, cold), tower 74s, keep 48s, herald 50s, web 47s, steward 30s, chronicle 29s, lookout 26s |
| B, tag `s-0902` | `af55c04` | 41s | chronicle 30s; the other seven 1s to 2s each from Docker's layer cache |

Only chronicle differs between the two commits, so build B rebuilt chronicle alone. A story's build costs what its change rebuilds, once the host's cache is warm. A cold portcullis costs 13 minutes.

### Per-namespace settings

Each namespace had these of its own, all made by `ns-up.sh`:

| Setting | `s-0901` | `s-0902` | How it was made |
|---------|----------|----------|-----------------|
| Kubernetes namespace and eight service accounts | `s-0901` | `s-0902` | The commit's `00-namespace.yaml`, renamed |
| SPIFFE IDs | `…/ns/s-0901/sa/<service>` | `…/ns/s-0902/sa/<service>` | kluster's `ClusterSPIFFEID`, unchanged |
| RabbitMQ vhost and user | `s-0901`, user `s-0901` | `s-0902`, user `s-0902` | `rabbitmqctl`, permissions on its own vhost only |
| Postgres databases | `s0901_tower`, `s0901_keep`, `s0901_herald` | `s0902_…` | `CREATE DATABASE`, on the shared superuser |
| ClickHouse database | `s0901_chronicle` | `s0902_chronicle` | `CREATE DATABASE`; chronicle's `clickhouseDatabase` names it |
| Signet bundles | eight, under signet namespace `s-0901` | eight, under `s-0902` | `storens`, from the commit's own config and secrets |
| Signet policies | four, for `…/ns/s-0901/sa/herald` and `/sa/tower` | four, for `s-0902` | `storens`, through `EnsureSignetPolicy` |
| Operator tokens | fresh `towerOperatorToken`, `heraldOperatorToken` | fresh | `openssl rand -hex 32`, kept in `/tmp/S-0575/logs/<namespace>.tokens` |
| Signet keys | fresh `tenant-acme-session-key` and `-internal-key` (EC P-256), identity hash key, peppers, salt | fresh | `spec.sh`: `openssl ecparam -name prime256v1 -genkey` for each EC key, random values for the rest |
| authstar tenant | `acme`, on host `acme.authstar.app` | `acme`, on the same host | The golden path provisions it; each namespace's own tower and portcullis hold it |

The tenant is the one setting authstar does not let a namespace own. Its host, `acme.authstar.app`, is fixed in portcullis's config, in Dex's redirect URI, and in the golden path's `Host:` header. The scenario reaches each namespace through port-forwards, so the two never met. A browser demo of two namespaces at once would need a host per namespace, such as `acme.s-0901.authstar.app`, in portcullis's config and in Dex's client.

Both namespaces ran at once, with all eight services in each `Running` and no restarts. Each used only its own shared-service settings:

| Check | `authstar` | `s-0901` | `s-0902` |
|-------|------------|----------|----------|
| RabbitMQ queues in its vhost | 6 | 6 | 6 |
| RabbitMQ connections to its vhost, and their user | 7, kluster's `rmq-…` | 7, `s-0901` | 7, `s-0902` |
| `tenant_event` rows in its tower database | 1 | 3 | 3 |
| Tables in its chronicle database | 1 | 1 | 1 |

### Load and deploy times

From the last, clean run of each, in seconds:

| Step | `s-0901` (build A) | `s-0902` (build B) |
|------|--------------------|--------------------|
| Load eight images into the node | 31 | 32 |
| Namespace and service accounts | 1 | 0 |
| RabbitMQ vhost, user, permissions | 34 | 34 |
| Postgres and ClickHouse databases | 8 | 8 |
| Signet bundles and policies | 4 | 4 |
| Apply manifests and wait for eight rollouts | 11 | 12 |
| **Total** | **89** | **90** |

About a minute and a half per story namespace, against 11m31s for the cluster. The RabbitMQ step is four `rabbitmqctl` calls in one `kubectl exec`, and each call starts an Erlang VM. RabbitMQ's HTTP management API, or the topology operator's `Vhost` and `Permission` resources that kluster itself uses, were not timed.

### Golden path

Final runs, both after the changes under [What it took](#what-it-took):

| Step | `s-0901` (build A) | `s-0902` (build B) |
|------|--------------------|--------------------|
| 1. Tenant provisioning | PASS (201) | PASS (201) |
| 2. portcullis and Dex login | PASS | PASS |
| 2. tower resolves the account | PASS | PASS |
| 3. keep billing webhook | PASS | PASS |
| 4. herald key rotation | PASS | PASS |
| 5. herald send and verify | PASS: the invite decrypted; the send fails as expected, with no send adapter | PASS, the same |
| 6. chronicle refuses the report without a token (401) | PASS | PASS |
| 6. chronicle report read with a forwarded token | SKIP: the scenario defers it to authstar S-0017 | SKIP |
| Exit, time | 0, 7s | 0, 8s |

Build B passes step 6 as build A does. The 401 check does not tell chronicle's old token check from its new one, and the step that would is skipped.

### What it took

The golden path passed in both namespaces only after six changes. None is about namespaces: each would stop the default `authstar` namespace too. Each was made outside the three repositories.

| # | What failed | Where the fault is | What the experiment did |
|---|-------------|--------------------|-------------------------|
| 1 | Tenant provisioning, HTTP 500, `SSL routines:tls_validate_record_header:wrong version number` | authstar: `integration/manifests/1x-*.yaml` set `SIGNET_ADMIN_ADDR` to signet's admin port 8444. Since ADR 0103 tower writes through workload mTLS (`tower/src/lib/signet-admin.ts`), which signetd serves on 8443, and kluster serves 8444 without TLS. | Renamed `:8444` to `:8443` in the manifests |
| 2 | tower crash-loops: `signet bundle missing provisioning.defaultProvider.scope` | authstar: `integration/config/authstar/tower.yaml` gives `scope: openid email`, a string; tower wants a list | `spec.sh` writes `scope: [openid, email]` |
| 3 | steward and web crash-loop: their signet bundle lacks `towerUrl` or `towerBaseUrl` | authstar's README says both read plain environment variables; both are in Signet-Configured Mode now. signet: a bundle with config and no secrets fails `SyncBundle` with `walk secrets dir: lstat …/secrets: no such file or directory`, and `extractTarGz` in `internal/api/bundle.go` keeps regular files only, so an empty folder entry does not help. | Patched the scratch copy of kluster-lib's `tarGzMultiFile` to add an empty `secrets/.keep` to every bundle; signetd's walk skips it as not `.yaml` |
| 4 | Step 2 stops the scenario with no message: Dex answers `Unregistered redirect_uri ("https://acme.authstar.app/oauth2/dex/callback")` | kluster: `addon/dex.go`'s static client lists only `http://localhost:8080/auth/callback` and `http://localhost:3000/callback` | Added the tenant's callback to the cluster's Dex config secret and restarted Dex |
| 5 | Step 2 again: Dex answers `Invalid Email Address and password` | kluster: `addon/dex.go`'s comment says its hash is bcrypt of `password`, but Python's `bcrypt.checkpw` rejects `password`, `Password`, and `admin` | Put a fresh bcrypt hash of `password` in the cluster's Dex config |
| 6 | Step 5 times out: no `communication.sent` or `communication.failed` within 10s | authstar: `scenario-golden-path.sh` reads Postgres databases `tower` and `herald` by name | `DB_PREFIX` in the scratch copy |

Two more cost runs:

- The scenario's port-forwards outlive it. `port_forward` `disown`s each one, so the EXIT trap's `kill $(jobs -p)` never reaches them. A second run then reached the first run's namespace on the same local ports, and passed or failed against the wrong build. `golden.sh` kills them after each run, and two namespaces' scenarios cannot run at once on the default ports.
- Any curl transport error ends the scenario under `set -e`, with no `FAIL` line and no summary.

### Teardown

Deleting `s-0901` with its eight services running took 69 seconds, and 70 in the run before. Two earlier teardowns, after deploys that had stopped before any pod started, took 40 and 37 seconds, 7 and 6 of them for the namespace. The last run:

| Step | Seconds |
|------|---------|
| Delete the namespace and wait for it to go | 37 |
| Delete the RabbitMQ vhost and user | 28 |
| Drop three Postgres and one ClickHouse database | 4 |

Afterwards `s-0901`, its vhost, and its databases were gone, and `s-0902`'s and the default namespace's remained. Every pod in the cluster was `Running` or `Completed`. The golden path against `s-0902` passed again, exit 0 in 11 seconds.

Signet still holds `s-0901`'s four policies, as `signet policy list` showed after the teardown, each once though two of the namespace's deploys ensured them. Its eight bundles stay too, since nothing deleted them. Only a workload whose SPIFFE ID names `s-0901` could use them. But a story serial is never reused, so they only accumulate.

tower and herald in `s-0902` had restarted three and two times, with exit 0. That is their design: each golden-path run's herald key rotation changes their signet bundles, and each logs `signet bundle changed -- acquired restart lock, shutting down for Kubernetes to restart` about 10 seconds after `rotate-keys`.

### Promotion

| Step | Seconds |
|------|---------|
| Delete kluster's Deployments and CronJobs in `authstar` | 0 |
| Load eight images | 33 |
| Namespace, RabbitMQ user, databases, signet | 35 |
| Apply manifests and wait for eight rollouts | 19 |
| **Total** | **87** |

All eight services in `authstar` then ran build A's images, `docker.io/bytepunx/authstar-<service>:s-0901`. The first golden path against `authstar` failed one step:

```text
FAIL: herald send-and-verify -- unexpected outcome (communication.failed):
payload decryption failed: error:02000079:rsa routines::oaep decoding error
```

The same run's step 4 had rotated herald's key, and tower took the new key only on its restart about 10 seconds later. The scenario sends the invite within a second of the rotation, so tower encrypted it to the old key. A second run 45 seconds later passed, exit 0 in 8 seconds, and `s-0902` passed alongside it, exit 0 in 7 seconds. `s-0901` and `s-0902` never showed this because the scenario's earlier runs there had rotated the key already. It is a race in the scenario's first run after any fresh deploy, not a cost of the namespace.

### What kluster lacks

For flai to do the above with no hand work, as S-0577 to S-0579 mean it to:

| # | Lack | What the experiment did instead |
|---|------|---------------------------------|
| 1 | An image import. There is no `kluster load` or equivalent; `kluster --help` lists `charts`, `down`, `kubeconfig`, `setup`, `status`, `up`, and `use`. | `docker save` into `ctr -n k8s.io images import` in the k3d node container, 31s to 34s for eight images |
| 2 | A deploy into a namespace. The `authstar` profile deploys authstar itself, in Go, into the constant `authStarNamespace = "authstar"`, from Docker Hub's `latest`. | The commit's own manifests, renamed, through `kubectl apply` |
| 3 | Per-namespace shared-service setup: a vhost and user, databases, signet bundles and policies for a namespace it was not built for. kluster-lib has `SeedSignetBundle` and `EnsureSignetPolicy`, but no command offers them. | `rabbitmqctl`, `psql`, and `clickhouse-client` through `kubectl exec`, and `storens` over kluster-lib |
| 4 | A teardown of a namespace and what it made in the shared services, signet's bundles and policies included | `ns-down.sh`, with signet's left behind |
| 5 | Machine-readable status. `kluster status` has no flag beyond `--help` and prints a table of `NAME`, `RUNNING`, and `AGE`; `kluster up` reports progress in prose. Neither says whether an addon answers. | `kubectl` and each service's own client, by hand |
| 6 | A pinned signet. `kluster charts list` pins every addon but signet; `up` warned `no pinned version for addon "signet"` and installed chart 0.11.3. The `authstar` profile's `latest` images do not provision a tenant against it. | Local builds |
| 7 | A Dex that works for authstar's tenant: no `https://acme.authstar.app/oauth2/dex/callback` redirect URI, and a password hash that does not match `password` | Both changed in the cluster's Dex secret |

Three lacks are outside kluster. In signet, a bundle with config and no secrets fails `SyncBundle`, so kluster-lib's `SeedSignetBundle` cannot seed steward or web. Also in signet, the CLI cannot delete a bundle's secrets or config, though the admin API has `DeleteSecret` and `DeleteConfig`. In authstar, the harness has drifted from its services: the four changes under [What it took](#what-it-took), the README's claim that steward and web do not read signet, the scenario's leaked port-forwards, and the tenant host fixed at `acme.authstar.app`.

### What was read and what was changed

Nothing was committed in authstar, signet, or kluster. At the end each stood at the commit it had at the start, with no commit made since 21:57Z on any branch: authstar `e02e430`, signet `8389600`, kluster `6d7a87c`. No file in any of the three changed after the story started. authstar's 26 and signet's one uncommitted changes, all under `design/digest` and `wip`, predate it (last written 21:13Z and 17:29Z), and kluster's checkout is clean.

| Repository | Read | Run | Changed, in a copy outside it |
|------------|------|-----|-------------------------------|
| authstar | Sources at `e02e430` and `af55c04`, by `git archive`; `integration/` | Each service's `Dockerfile`; the scenario, from the copy | `lib.sh` and `scenario-golden-path.sh` in `bin/s-0575/scripts`, by variables; the manifests and config, renamed on their way into the cluster |
| signet | The server's sync, bundle, and listener code; the chart's values | The installed CLI 0.10.3, for `status` | Nothing |
| kluster | The `authstar`, `signet`, `spire`, and `dex` profiles and addons | The installed `kluster` 1.3.1, for `up`, `kubeconfig`, and `charts list` | `kluster-lib` in `bin/s-0575/kluster-lib`: one function, and the added `storens` |

## Recommendation

**Adapt.** A namespace per story over one shared kluster works for authstar. Two builds ran at once, each passed the golden path on its own vhost, databases, and signet bundles, and one was torn down and another promoted with the rest still working. A story costs about 90 seconds to deploy and 70 to tear down, plus whatever its change rebuilds, against 11m31s for the cluster. So S-0576's ADR should take one cluster per project and a namespace per story, named for the story's serial, with these per-namespace settings:

- a RabbitMQ vhost and user of the namespace's name
- databases named with the namespace's prefix
- signet bundles and policies under a signet namespace equal to the Kubernetes namespace, with fresh keys and operator tokens

Adapt rather than adopt, for three reasons:

- **The deploy is the project's, not kluster's.** What made a namespace was authstar's manifests and config with the namespace substituted, and per-service settings only authstar knows: which databases, which secrets are keys, which services grant which writes. The ADR should have a project name a deploy command that takes the namespace, its image tag, and the shared services' endpoints, as S-0578 plans, and flai should run it. kluster should keep the cluster and the shared services, and gain an image import, a namespace setup and teardown for the shared services, signet's included, and a JSON status, lacks 1 and 3 to 5 above. Those changes are kluster's, made in kluster.
- **A baseline has to pass before any namespace means anything.** kluster's own `authstar` deployment failed at step 1, and the golden path passed only after six fixes to drift between authstar, signet, and kluster. Those go to the three repositories as issues. Until they land, flai's review check against a namespace would fail for reasons that have nothing to do with the story.
- **Two things a namespace cannot own today.** The tenant's host is fixed at `acme.authstar.app`, so a browser demo of two namespaces at once needs a host per namespace in portcullis's config and Dex's client. And the scenario's fixed local ports and leaked port-forwards stop two namespaces' checks running at once. S-0579's review checks should run one namespace at a time, or give each its own ports, until authstar changes.

## Scripts

What the experiment ran, from the worktree's git-ignored `bin/s-0575/`. Logs and secrets stayed in `/tmp/S-0575`. Tabs are shown as four spaces.

### bin/s-0575/build.sh

```sh
#!/usr/bin/env bash
# S-0575: build authstar's eight service images from an exported commit,
# tagged for a story namespace, timing each build.
#   build.sh <source dir> <tag> <log dir>
set -euo pipefail
src="$1" tag="$2" logs="$3"
times="$logs/build-$tag.times"
: > "$times"
for svc in tower keep herald lookout chronicle steward web portcullis; do
  ref="docker.io/bytepunx/authstar-${svc}:${tag}"
  s=$(date +%s)
  if (cd "$src/$svc" && docker build -t "$ref" .) > "$logs/build-$tag-$svc.log" 2>&1; then
    echo "$svc ok $(( $(date +%s) - s ))s" >> "$times"
  else
    echo "$svc FAILED $(( $(date +%s) - s ))s" >> "$times"
  fi
done
echo "done $(date -u +%FT%TZ)" >> "$times"
```

### bin/s-0575/ns-up.sh

```sh
#!/usr/bin/env bash
# S-0575: load one build into the kluster and deploy it into a story namespace
# over the shared services, timing each step into <logs>/ns-<namespace>.times.
#   ns-up.sh <namespace> <image tag> <source dir> <logs dir>
# Needs KUBECONFIG pointing at the authstar-integration cluster.
set -euo pipefail
ns="$1" tag="$2" src="$3" logs="$4"
here="$(cd "$(dirname "$0")" && pwd)"
prefix="${ns//-/}_"                      # s-0901 -> s0901_
[ "$ns" = authstar ] && prefix=""        # main's namespace keeps kluster's names
services=(tower keep herald lookout chronicle steward web portcullis)
times="$logs/ns-$ns.times"
: > "$times"
step() { local name="$1"; shift; local s; s=$(date +%s); "$@"; echo "$name $(( $(date +%s) - s ))s" >> "$times"; }

load() {
  for svc in "${services[@]}"; do
    docker save "docker.io/bytepunx/authstar-${svc}:${tag}" \
      | docker exec -i k3d-authstar-integration-server-0 ctr -n k8s.io images import - > /dev/null
  done
}

rename() {
  sed -e "s/^\(\s*namespace:\) authstar$/\1 ${ns}/" \
      -e "s/^\(\s*name:\) authstar$/\1 ${ns}/" \
      -e "s/^\(\s*value:\) authstar$/\1 ${ns}/" \
      -e "s/\.authstar\.svc\.cluster\.local/.${ns}.svc.cluster.local/g" \
      -e "s#image: bytepunx/authstar-\([a-z]*\):latest#image: docker.io/bytepunx/authstar-\1:${tag}#" \
      -e "s/imagePullPolicy: Always/imagePullPolicy: IfNotPresent/" \
      -e "s/signet\.signet\.svc\.cluster\.local:8444/signet.signet.svc.cluster.local:8443/"
  # The last rule: since ADR 0103 the services write through workload mTLS,
  # which signetd serves on 8443; the committed manifests still name the
  # admin port, which kluster serves without TLS.
}

namespace() { rename < "$src/integration/manifests/00-namespace.yaml" | kubectl apply -f - > /dev/null; }

rabbitmq() {
  local pw; pw="$(openssl rand -hex 24)"
  kubectl -n rabbitmq exec rabbitmq-0 -c rabbitmq -- sh -c "
    rabbitmqctl -q list_vhosts | grep -qx '$ns' || rabbitmqctl -q add_vhost '$ns'
    if rabbitmqctl -q list_users | cut -f1 | grep -qx '$ns'; then
      rabbitmqctl -q change_password '$ns' '$pw'
    else
      rabbitmqctl -q add_user '$ns' '$pw'
    fi
    rabbitmqctl -q set_permissions -p '$ns' '$ns' '.*' '.*' '.*'" > /dev/null
  printf 'amqp://%s:%s@rabbitmq.rabbitmq.svc.cluster.local:5672/%s' "$ns" "$pw" "$ns" > "$logs/$ns.amqp"
}

databases() {
  local psql=(kubectl -n postgres exec postgresql-0 -- env PGPASSWORD=postgres psql -U postgres -tAq)
  for db in tower keep herald; do
    if [ -z "$("${psql[@]}" -c "SELECT 1 FROM pg_database WHERE datname='${prefix}${db}'")" ]; then
      "${psql[@]}" -c "CREATE DATABASE ${prefix}${db}" > /dev/null
    fi
  done
  kubectl -n clickhouse exec clickhouse-shard0-0 -- clickhouse-client --password clickhouse -q "CREATE DATABASE IF NOT EXISTS ${prefix}chronicle"
}

signet() {
  local tt ht
  tt="$(openssl rand -hex 32)" ht="$(openssl rand -hex 32)"
  printf 'TOWER_OPERATOR_TOKEN=%s\nHERALD_OPERATOR_TOKEN=%s\n' "$tt" "$ht" > "$logs/$ns.tokens"
  bash "$here/spec.sh" "$src" "$ns" "$prefix" "$(cat "$logs/$ns.amqp")" "$tt" "$ht" > "$logs/$ns.spec.json"
  "$here/storens" -kubeconfig "$KUBECONFIG" < "$logs/$ns.spec.json" 2> "$logs/$ns.storens.log"
}

deploy() {
  for f in "$src"/integration/manifests/1[0-7]-*.yaml; do
    rename < "$f" | kubectl apply -f - > /dev/null
  done
  for svc in "${services[@]}"; do
    kubectl rollout status "deployment/$svc" -n "$ns" --timeout=300s > /dev/null
  done
}

step load load
step namespace namespace
step rabbitmq rabbitmq
step databases databases
step signet signet
step deploy deploy
echo "done $(date -u +%FT%TZ)" >> "$times"
```

### bin/s-0575/spec.sh

```sh
#!/usr/bin/env bash
# S-0575: print the storens spec for one story namespace, from one commit's
# own integration/config/authstar/*.yaml and integration/scripts/lib-secrets.sh,
# with the namespace's own values in place of the authstar namespace's:
#   spec.sh <source dir> <namespace> <db prefix> <amqp url> <tower token> <herald token>
# - every *.authstar.svc.cluster.local host becomes *.<namespace>.svc.cluster.local
# - dbConnectionString points at <db prefix><service> on kluster's Postgres
# - rabbitmqUrl is the namespace's own user on its own vhost
# - operatorToken is the namespace's own, given on the command line
# - every EC private key, dev-placeholder pepper, salt, or key, and hex secret
#   is generated fresh, so no two namespaces share key material
# - chronicle's clickhouseDatabase is <db prefix>chronicle
set -euo pipefail
src="$1" ns="$2" prefix="$3" amqp="$4" tower_token="$5" herald_token="$6"
h="$src/integration"

# shellcheck source=/dev/null
source "$h/scripts/lib-secrets.sh"

rename_hosts() { sed "s/\.authstar\.svc\.cluster\.local/.${ns}.svc.cluster.local/g"; }
fresh_placeholder() { printf 'dev-%s-%s' "$ns" "$(openssl rand -base64 32 | tr '+/' '-_' | tr -d '=\n')"; }

secrets='{}'
for entry in "${SECRET_ENTRIES[@]}"; do
  svc="${entry%%/*}" name="${entry#*/}"
  var="SECRET_VALUE__${svc}__${name}"; var="${var//-/_}"
  value="${!var}"
  case "$name" in
    dbConnectionString) value="postgres://postgres:postgres@postgresql.postgres.svc.cluster.local:5432/${prefix}${svc}?sslmode=disable" ;;
    rabbitmqUrl) value="$amqp" ;;
    operatorToken) [ "$svc" = tower ] && value="$tower_token" || value="$herald_token" ;;
    *)
      case "$value" in
        "-----BEGIN EC PRIVATE KEY-----"*) value="$(openssl ecparam -name prime256v1 -genkey -noout 2>/dev/null)" ;;
        dev-placeholder-*) value="$(fresh_placeholder)" ;;
        *) if printf %s "$value" | grep -Eq '^[0-9a-f]{32,}$'; then value="$(openssl rand -hex $(( ${#value} / 2 )))"; fi ;;
      esac
      value="$(printf %s "$value" | rename_hosts)" ;;
  esac
  secrets="$(jq --arg s "$svc" --arg k "$name" --arg v "$value" '.[$s][$k] = $v' <<<"$secrets")"
done

bundles='{}'
for f in "$h"/config/authstar/*.yaml; do
  svc="$(basename "$f" .yaml)"
  # tower wants defaultProvider.scope as a list; the committed config gives a string.
  config="$(rename_hosts < "$f" | sed -e "s/^clickhouseDatabase: chronicle$/clickhouseDatabase: ${prefix}chronicle/" \
    -e "s/^\(\s*scope:\) openid email$/\1 [openid, email]/")"
  bundles="$(jq --arg s "$svc" --arg c "$config" --argjson sec "$(jq --arg s "$svc" '.[$s] // {}' <<<"$secrets")" \
    '.[$s] = {secrets: $sec, config: $c}' <<<"$bundles")"
done

jq -n --arg ns "$ns" --argjson b "$bundles" \
  '{ns: $ns, bundles: $b, policies: [["herald","tower"],["herald","keep"],["herald","portcullis"],["tower","portcullis"]]}'
```

### bin/s-0575/kluster-lib/cmd/storens/main.go

```go
// Command storens seeds signet for one story namespace of S-0575's experiment,
// with kluster-lib's own SeedSignetBundle and EnsureSignetPolicy, from a JSON
// spec on standard input:
//
//  {"ns": "s-0901",
//   "bundles": {"tower": {"secrets": {"k": "v"}, "config": "<yaml>"}},
//   "policies": [["herald", "tower"]]}
//
// A policy [source, target] lets spiffe://<trust domain>/ns/<ns>/sa/<source>
// put into <ns>/<target>, as kluster's authstar profile grants in "authstar".
package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "sort"

    "github.com/bytepunx/kluster-lib/addon"
    "github.com/bytepunx/kluster-lib/provider"
    "k8s.io/client-go/tools/clientcmd"
)

type bundle struct {
    Secrets map[string]string `json:"secrets"`
    Config  *string           `json:"config"`
}

type spec struct {
    NS       string            `json:"ns"`
    Bundles  map[string]bundle `json:"bundles"`
    Policies [][2]string       `json:"policies"`
}

func must(err error, what string) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "storens: %s: %v\n", what, err)
        os.Exit(1)
    }
}

func main() {
    kubeconfig := flag.String("kubeconfig", "", "kubeconfig file")
    trustDomain := flag.String("trust-domain", "authstar.cluster.local", "SPIFFE trust domain")
    flag.Parse()

    var s spec
    must(json.NewDecoder(os.Stdin).Decode(&s), "read spec")
    if s.NS == "" || *kubeconfig == "" {
        fmt.Fprintln(os.Stderr, "storens: the spec's ns and -kubeconfig are required")
        os.Exit(2)
    }

    rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    must(err, "load kubeconfig")
    h, err := addon.NewClusterHandle(rc, provider.ClusterConfig{Name: "authstar-integration", TrustDomain: *trustDomain})
    must(err, "cluster handle")
    ctx := context.Background()

    services := make([]string, 0, len(s.Bundles))
    for svc := range s.Bundles {
        services = append(services, svc)
    }
    sort.Strings(services)
    for _, svc := range services {
        b := s.Bundles[svc]
        var config []byte
        if b.Config != nil {
            config = []byte(*b.Config)
        }
        must(addon.SeedSignetBundle(ctx, h, s.NS, svc, b.Secrets, config), "seed "+svc)
        fmt.Fprintf(os.Stderr, "seeded %s/%s: %d secrets, config %t\n", s.NS, svc, len(b.Secrets), config != nil)
    }
    for _, p := range s.Policies {
        id := fmt.Sprintf("spiffe://%s/ns/%s/sa/%s", *trustDomain, s.NS, p[0])
        must(addon.EnsureSignetPolicy(ctx, h, id, s.NS, p[1], []string{"put"}), "policy "+p[0]+"->"+p[1])
        fmt.Fprintf(os.Stderr, "policy %s -> put on %s/%s/*\n", id, s.NS, p[1])
    }
}
```

### The scratch change to kluster-lib

`bin/s-0575/kluster-lib` is a copy of kluster 1.3.1's `kluster-lib`, built with Go 1.26.3 in scratch caches (`GOPATH=/tmp/S-0575/gopath GOTOOLCHAIN=go1.26.3`). Besides the added `cmd/storens` above, this is its whole difference from the original:

```diff
--- kluster/kluster-lib/addon/signet_admin.go
+++ bin/s-0575/kluster-lib/addon/signet_admin.go
@@ -323,6 +323,12 @@
    gz := gzip.NewWriter(&buf)
    tw := tar.NewWriter(gz)
 
+   // S-0575 scratch change: signetd's SyncBundle fails on a bundle with no
+   // secrets folder and extracts regular files only, so every bundle carries
+   // an empty secrets/.keep, which its walk skips as not .yaml.
+   if err := tw.WriteHeader(&tar.Header{Name: "secrets/.keep", Mode: 0o600, Size: 0}); err != nil {
+       return nil, fmt.Errorf("tar header for secrets/.keep: %w", err)
+   }
    for name, content := range files {
        hdr := &tar.Header{Name: name, Mode: 0o600, Size: int64(len(content))}
        if err := tw.WriteHeader(hdr); err != nil {
```

### bin/s-0575/ns-down.sh

```sh
#!/usr/bin/env bash
# S-0575: tear down a story namespace and what it made in the shared services,
# timing it into <logs>/ns-<namespace>-down.times.
#   ns-down.sh <namespace> <logs dir>
# Signet keeps the namespace's bundles and policies. Its admin API has
# DeleteSecret, DeleteConfig, and DeletePolicy, but the signet CLI 0.10.3
# removes only a policy, by ID, and kluster-lib exposes none of the three.
set -euo pipefail
ns="$1" logs="$2"
prefix="${ns//-/}_"
times="$logs/ns-$ns-down.times"
: > "$times"
step() { local name="$1"; shift; local s; s=$(date +%s); "$@"; echo "$name $(( $(date +%s) - s ))s" >> "$times"; }

namespace() { kubectl delete namespace "$ns" --wait=true --timeout=300s > /dev/null; }

rabbitmq() {
  kubectl -n rabbitmq exec rabbitmq-0 -c rabbitmq -- sh -c "
    rabbitmqctl -q list_vhosts | grep -qx '$ns' && rabbitmqctl -q delete_vhost '$ns'
    rabbitmqctl -q list_users | cut -f1 | grep -qx '$ns' && rabbitmqctl -q delete_user '$ns'
    true" > /dev/null
}

databases() {
  for db in tower keep herald; do
    kubectl -n postgres exec postgresql-0 -- env PGPASSWORD=postgres psql -U postgres -qc "DROP DATABASE IF EXISTS ${prefix}${db} WITH (FORCE)" > /dev/null
  done
  kubectl -n clickhouse exec clickhouse-shard0-0 -- clickhouse-client --password clickhouse -q "DROP DATABASE IF EXISTS ${prefix}chronicle"
}

step namespace namespace
step rabbitmq rabbitmq
step databases databases
echo "done $(date -u +%FT%TZ)" >> "$times"
```

### bin/s-0575/promote.sh

```sh
#!/usr/bin/env bash
# S-0575: promote a story's build to the default authstar namespace, as an
# acceptance would: remove the Deployments and CronJobs kluster made there,
# whose specs a kubectl apply would only merge into, then deploy the build
# with ns-up.sh, timing the removal into <logs>/promote.times.
#   promote.sh <image tag> <source dir> <logs dir>
set -euo pipefail
tag="$1" src="$2" logs="$3"
here="$(cd "$(dirname "$0")" && pwd)"
s=$(date +%s)
kubectl -n authstar delete deployment,cronjob --all --wait=true --timeout=300s > /dev/null
echo "remove $(( $(date +%s) - s ))s" > "$logs/promote.times"
bash "$here/ns-up.sh" authstar "$tag" "$src" "$logs"
```

### bin/s-0575/golden.sh

```sh
#!/usr/bin/env bash
# S-0575: run authstar's golden path against one namespace, with that
# namespace's operator tokens, and kill the port-forwards the scenario leaves
# behind when it stops early (its EXIT trap does not reach them).
#   golden.sh <namespace> <tokens file> <log file>
set -uo pipefail
ns="$1" tokens="$2" log="$3"
here="$(cd "$(dirname "$0")" && pwd)"
set -a
# shellcheck source=/dev/null
. "$tokens"
set +a
export AUTHSTAR_NAMESPACE="$ns"
[ "$ns" = authstar ] && export DB_PREFIX="" || export DB_PREFIX="${ns//-/}_"
s=$(date +%s)
bash "$here/scripts/scenario-golden-path.sh" > "$log" 2>&1
rc=$?
echo "rc=$rc seconds=$(( $(date +%s) - s ))" >> "$log"
pkill -f "kubectl port-forward -n ($ns|dex) " || true
exit "$rc"
```

### bin/s-0575/path/nc

```sh
#!/usr/bin/env bash
# nc shim for S-0575: the host has no nc, and authstar's scripts only call
# `nc -z <host> <port>` to wait for a port-forward.
[ "${1:-}" = "-z" ] || { echo "nc shim: only -z is supported" >&2; exit 2; }
exec 3<>"/dev/tcp/$2/$3" || exit 1
exec 3>&-
```

### The changes to the copy of authstar's scripts

`bin/s-0575/scripts` is `integration/scripts` at `e02e430`. Its whole difference:

```diff
--- authstar/integration/scripts/lib.sh
+++ bin/s-0575/scripts/lib.sh
@@ -14,7 +14,7 @@
 # from the environment with no hardcoded default, so this is the value this
 # harness's manifests (a later task) will set for them. lookout/chronicle/
 # steward/web have no live signet wiring yet — see README.md.
-AUTHSTAR_NAMESPACE="authstar"
+AUTHSTAR_NAMESPACE="${AUTHSTAR_NAMESPACE:-authstar}"
 
 # The 8 authstar services, in the order they're discussed in claude.md.
 # scripts/provision.sh iterates this for --config-path layout only; which
--- authstar/integration/scripts/scenario-golden-path.sh
+++ bin/s-0575/scripts/scenario-golden-path.sh
@@ -238,7 +238,7 @@
 # adminEmail) (provisioning.ts) into tower's own `contacts` table -- query
 # it directly instead, same technique step 5's own herald verification uses
 # below.
-contact_id="$(kubectl exec -n postgres statefulset/postgresql -- env PGPASSWORD=postgres psql -U postgres -d tower -tA -c \
+contact_id="$(kubectl exec -n postgres statefulset/postgresql -- env PGPASSWORD=postgres psql -U postgres -d "${DB_PREFIX:-}tower" -tA -c \
   "select id from contacts where tenant = '${TENANT_SLUG}' and email = '${DEX_LOGIN}' limit 1;" 2>/dev/null || echo "")"
 
 if [[ -z "$contact_id" ]]; then
@@ -267,7 +267,7 @@
     outcome_type=""
     outcome_error=""
     for _ in $(seq 1 10); do
-      row="$(kubectl exec -n postgres statefulset/postgresql -- env PGPASSWORD=postgres psql -U postgres -d herald -tA -F'|' -c "$sql" 2>/dev/null || echo "")"
+      row="$(kubectl exec -n postgres statefulset/postgresql -- env PGPASSWORD=postgres psql -U postgres -d "${DB_PREFIX:-}herald" -tA -F'|' -c "$sql" 2>/dev/null || echo "")"
       if [[ -n "$row" ]]; then
         outcome_type="${row%%|*}"
         outcome_error="${row#*|}"
```

### The changes to the cluster's Dex

```sh
kubectl -n dex get secret dex -o jsonpath='{.data.config\.yaml}' | base64 -d > /tmp/S-0575/dex-config.yaml
# 1. under staticClients[authstar].redirectURIs, add:   - https://acme.authstar.app/oauth2/dex/callback
# 2. replace staticPasswords[admin@example.com].hash with a fresh one, checked first:
python3 -I -c "import bcrypt; h=bcrypt.hashpw(b'password', bcrypt.gensalt(10)); assert bcrypt.checkpw(b'password', h); print(h.decode())"
kubectl -n dex create secret generic dex --from-file=config.yaml=/tmp/S-0575/dex-config.yaml --dry-run=client -o yaml | kubectl -n dex apply -f -
kubectl -n dex rollout restart deploy/dex
```
