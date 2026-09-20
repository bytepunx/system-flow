# flaiover: the system-flow dashboard. Build from the repository root:
#   docker build -f flaiover/Dockerfile -t flaiover:local .
# Runs as any UID (flai dashboard passes --user); listens on 3000. The image
# holds no flai, no git, and no ssh (S-0075, ADR-0029): everything it shows
# and everything it changes is asked of flai serve on the host, which dials
# the /agent endpoint. It needs no writable home and nothing of the project.

ARG NODE_VERSION=24

# ---- flaiover build --------------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS build
ARG FLAI_COMMIT=unknown
WORKDIR /src/flaiover
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
RUN corepack enable
COPY flaiover/package.json flaiover/pnpm-lock.yaml flaiover/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY flaiover/ ./
RUN pnpm build && pnpm prune --prod

# ---- runtime ---------------------------------------------------------------
FROM node:${NODE_VERSION}-alpine
ARG FLAI_COMMIT=unknown
ARG FLAIOVER_VERSION=0.0.0
LABEL org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flaiover: the system-flow dashboard; it reaches a project through flai serve on the host" \
      org.opencontainers.image.licenses="UNLICENSED"
RUN apk add --no-cache tini
WORKDIR /app
COPY --from=build /src/flaiover/build ./build
COPY --from=build /src/flaiover/node_modules ./node_modules
COPY --from=build /src/flaiover/package.json ./package.json
COPY --from=build /src/flaiover/server.js ./server.js
ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0 \
    PROJECT_DIR=/project \
    HOME=/tmp \
    LOG_FORMAT=json \
    FLAIOVER_COMMIT=${FLAI_COMMIT} \
    FLAIOVER_VERSION=${FLAIOVER_VERSION}
EXPOSE 3000
# Unprivileged by default; any --user works because nothing here needs a home.
USER 65532:65532
ENTRYPOINT ["/sbin/tini", "--"]
# server.js is adapter-node's handler plus the WebSocket upgrade for /agent (ADR-0029).
CMD ["node", "server.js"]
