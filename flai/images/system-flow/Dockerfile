# flai-agent-sf-claude-code: this repository's layer over the claude-code harness image (S-0510,
# ADR-0203), named by system-flow.yaml's agent.image.layer. It carries the tools the manifest's test
# tiers run, so that flai test and scripts/close-out.sh run inside the agent's container: Go at
# the version agent.image.args names, flai/go.mod's, with gcc and musl-dev for -race, golangci-lint
# v2 built with that Go, ShellCheck, pnpm on the base image's Node 24, and for the smoke tier's install test curl, openssl,
# and GoReleaser (I-0281). Build it after the harness image, and check it:
#   flai image build --layer            (or scripts/agent-image.sh --layer; make agent-images builds all three)
#   scripts/agent-image.sh --check      (runs the smoke tier inside it)
# Both pass BASE, the harness image, and the build arguments system-flow.yaml names under
# agent.image.args (S-0760). The ARG defaults below are the same pins; keep the two in step.
# No login state and no key, as the base has none; it keeps the base's entrypoint, pilot.

ARG BASE=flai-agent-claude-code:local
ARG GO_VERSION=1.26.0
# The Alpine release of the golang image, exact (S-0630): golang publishes 1.26.0 on 3.23 at the newest.
ARG GO_ALPINE_VERSION=3.23
ARG SHELLCHECK_VERSION=v0.10.0

FROM golang:${GO_VERSION}-alpine${GO_ALPINE_VERSION} AS go
ARG GOLANGCI_LINT_VERSION=v2.5.0
# golangci-lint built with the module's own Go, as scripts/install-tools.sh builds it, since a
# release binary built with an older Go refuses to lint a module that targets a newer one.
RUN --mount=type=cache,target=/root/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOBIN=/out go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
ARG GORELEASER_VERSION=v2.18.1
# GoReleaser, for scripts/flai-snapshot.sh, which the smoke tier's install test runs. It builds only
# with the Go its go.mod names, newer than the module's, so the go command fetches that toolchain
# here, checked against Go's checksum database, as scripts/install-tools.sh's go install does on
# the host. Only the binary reaches the image; its Go stays at the module's.
RUN --mount=type=cache,target=/root/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOTOOLCHAIN=auto GOBIN=/out go install "github.com/goreleaser/goreleaser/v2@${GORELEASER_VERSION}"
# flai's modules, fetched once into a cache every uid can read, so that a test run in the
# container needs no download.
COPY flai/go.mod flai/go.sum /src/flai/
RUN cd /src/flai && GOMODCACHE=/opt/go/mod go mod download && chmod -R a+rX /opt/go/mod

FROM koalaman/shellcheck-alpine:${SHELLCHECK_VERSION} AS shellcheck

FROM ${BASE}
ARG GO_VERSION
ARG GOLANGCI_LINT_VERSION=v2.5.0
ARG SHELLCHECK_VERSION
ARG PNPM_VERSION=10.34.5
ARG GORELEASER_VERSION=v2.18.1
LABEL dev.system-flow.go="${GO_VERSION}" \
      dev.system-flow.golangci-lint="${GOLANGCI_LINT_VERSION}" \
      dev.system-flow.goreleaser="${GORELEASER_VERSION}" \
      dev.system-flow.shellcheck="${SHELLCHECK_VERSION}" \
      dev.system-flow.pnpm="${PNPM_VERSION}" \
      org.opencontainers.image.description="flai-agent-sf-claude-code: system-flow's test tools on flai-agent-claude-code"
USER root
# gcc and musl-dev for go test -race, which needs cgo; curl and openssl for install.sh, which the
# smoke tier's install test runs and which refuses without them (I-0281); procps-ng's ps for
# scripts/release-server.sh, whose ps -o stat= -p BusyBox's ps refuses, so that it read its server
# as exited before it listened; the rest is the base's.
RUN apk add --no-cache build-base curl openssl procps-ng && \
    npm install -g "pnpm@${PNPM_VERSION}" && \
    npm cache clean --force && rm -rf /root/.npm /tmp/node-compile-cache
COPY --from=go /usr/local/go /usr/local/go
COPY --from=go /out/golangci-lint /usr/local/bin/golangci-lint
COPY --from=go /out/goreleaser /usr/local/bin/goreleaser
COPY --from=go /opt/go/mod /opt/go/mod
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck
# The toolchain the image holds, never one go downloads; the module cache read-only, filled above;
# the build cache in the home any uid can write. go build stamps VCS information as on the host: for
# a story's worktree Go runs git in the main checkout, which git accepts through the safe.directory
# a story's container is given (ADR-0253, I-0255), and scripts/agent-vcs-test.sh checks it.
ENV PATH=/usr/local/go/bin:$PATH \
    GOTOOLCHAIN=local \
    GOMODCACHE=/opt/go/mod
USER 65532:65532
