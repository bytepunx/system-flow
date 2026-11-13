# flai-agent-bedrock: the base of every agent image (E-0029): flai, pilot, git, ca-certificates,
# and tini. pilot is the entrypoint (ADR-0212): it reads one start request on standard input and
# runs the inner harness a layer above installs. Two targets give flai and pilot (ADR-0295):
#
#   tree, the default: built from this tree. Build from the repository root:
#     docker build -f flai/images/bedrock/Dockerfile -t flai-agent-bedrock:local .
#   scripts/agent-image.sh builds it, then the agent images on top of it. Not stamped, not signed.
#
#   prebuilt: copies a release's flai and pilot for Linux from the build context, a folder that
#   holds linux-<arch>/flai and linux-<arch>/pilot for each architecture built, the arch as
#   buildx's TARGETARCH names it (amd64, arm64). release-flai.yml builds and publishes the
#   release's image from it, and a host build from a release's archives uses it too:
#     docker buildx build -f flai/images/bedrock/Dockerfile --target prebuilt \
#       --build-arg FLAI_VERSION=1.2.3 --build-arg FLAI_COMMIT=abc1234 <folder>
#
# Alpine, not Debian slim: Claude Code's npm package ships a musl build of its native binary
# (@anthropic-ai/claude-code-linux-*-musl), which runs with libgcc, libstdc++, and the system
# ripgrep, and Alpine's own nodejs package is Node 24, so the claude-code layer needs no other source.

# Exact versions (S-0630, I-0297), each base named by its tag and its index digest, since a signed
# flai release builds from this file (ADR-0295); dependabot's docker entry moves both together.
# GO_VERSION is flai/go.mod's; golang publishes 1.26.0 on Alpine 3.23 at the newest, which only
# builds the binaries, so the runtime's Alpine is its own.
ARG GO_VERSION=1.26.0
ARG GO_ALPINE_VERSION=3.23
ARG ALPINE_VERSION=3.24.4

# ---- flai build from the tree ----------------------------------------------
FROM golang:${GO_VERSION}-alpine${GO_ALPINE_VERSION}@sha256:d4c4845f5d60c6a974c6000ce58ae079328d03ab7f721a0734277e69905473e5 AS build
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

# ---- runtime, shared by both targets ----------------------------------------
# dev.system-flow.flai is the version the image's flai reports, and dev.system-flow.commit the
# commit it was built from: the labels that tie the image, and every layer over it, to its flai
# (ADR-0295). Both targets take them from the build arguments; a release passes its own.
FROM alpine:${ALPINE_VERSION}@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS runtime
ARG FLAI_VERSION=dev
ARG FLAI_COMMIT=unknown
LABEL dev.system-flow.flai="${FLAI_VERSION}" \
      dev.system-flow.commit="${FLAI_COMMIT}" \
      org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flai-agent-bedrock: flai, pilot, git, and tini, the base of the flai agent images" \
      org.opencontainers.image.licenses="LicenseRef-system-flow-Shield-1.0"
RUN apk add --no-cache git ca-certificates tini
# The container runs as the operator's uid (--user), which has no entry here: a home any uid can
# write to holds what git and the agent keep there.
RUN addgroup -g 65532 agent && adduser -D -u 65532 -G agent -h /home/agent agent && \
    chmod 1777 /home/agent
ENV HOME=/home/agent
WORKDIR /home/agent
# A start request on standard input, the harness's log on standard output: docker run -i, no command.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/pilot"]

# ---- prebuilt: a release's flai and pilot from the build context -----------
FROM runtime AS prebuilt
ARG TARGETARCH
COPY --chmod=0755 linux-${TARGETARCH}/flai linux-${TARGETARCH}/pilot /usr/local/bin/
USER 65532:65532

# ---- tree, the default: flai and pilot built above -------------------------
FROM runtime AS tree
COPY --from=build /out/flai /out/pilot /usr/local/bin/
USER 65532:65532
