# flai-agent-bedrock: the base of every agent image (E-0029), flai and pilot built from this tree
# with git, ca-certificates, and tini. pilot is the entrypoint (ADR-0212): it reads one start request
# on standard input and runs the inner harness a layer above installs. Build from the repository root:
#   docker build -f flai/images/bedrock/Dockerfile -t flai-agent-bedrock:local .
# scripts/agent-image.sh builds it, then the agent images on top of it. Not stamped, not signed.
#
# Alpine, not Debian slim: Claude Code's npm package ships a musl build of its native binary
# (@anthropic-ai/claude-code-linux-*-musl), which runs with libgcc, libstdc++, and the system
# ripgrep, and Alpine's own nodejs package is Node 24, so the claude-code layer needs no other source.

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.24

# ---- flai build ------------------------------------------------------------
FROM golang:${GO_VERSION}-alpine AS build
ARG FLAI_VERSION=dev
ARG FLAI_COMMIT=unknown
ARG FLAI_DATE=
WORKDIR /src/flai
COPY flai/go.mod flai/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY flai/ ./
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -o /out/flai -ldflags "-s -w \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=${FLAI_VERSION} \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=${FLAI_COMMIT} \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Date=${FLAI_DATE}" . && \
    CGO_ENABLED=0 go build -trimpath -o /out/pilot -ldflags "-s -w" ./pilot

# ---- runtime ---------------------------------------------------------------
FROM alpine:${ALPINE_VERSION}
ARG FLAI_VERSION=dev
LABEL dev.system-flow.flai="${FLAI_VERSION}" \
      org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flai-agent-bedrock: flai, pilot, git, and tini, the base of the flai agent images" \
      org.opencontainers.image.licenses="LicenseRef-system-flow-Shield-1.0"
RUN apk add --no-cache git ca-certificates tini
COPY --from=build /out/flai /out/pilot /usr/local/bin/
# The container runs as the operator's uid (--user), which has no entry here: a home any uid can
# write to holds what git and the agent keep there.
RUN addgroup -g 65532 agent && adduser -D -u 65532 -G agent -h /home/agent agent && \
    chmod 1777 /home/agent
ENV HOME=/home/agent
WORKDIR /home/agent
USER 65532:65532
# A start request on standard input, the harness's log on standard output: docker run -i, no command.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/pilot"]
