# flai-agent-claude-code: Claude Code on flai-agent-bedrock (E-0029). Build after the bedrock image:
#   docker build -f flai/images/claude-code/Dockerfile -t flai-agent-claude-code:local .
# scripts/agent-image.sh builds both in order. No login state and no key: the agent's credentials
# reach the container when it runs, never the image.

ARG BASE=flai-agent-bedrock:local
FROM ${BASE}
ARG CLAUDE_CODE_VERSION=2.1.312
LABEL dev.system-flow.claude-code="${CLAUDE_CODE_VERSION}" \
      org.opencontainers.image.description="flai-agent-claude-code: Claude Code on flai-agent-bedrock"
USER root
# Installed as Anthropic documents it, from npm; on musl it needs libgcc, libstdc++, and ripgrep.
# Its Bash tool needs bash or zsh: with busybox sh alone it answers "No suitable shell found" (S-0505).
RUN apk add --no-cache nodejs npm libgcc libstdc++ ripgrep bash && \
    npm install -g "@anthropic-ai/claude-code@${CLAUDE_CODE_VERSION}" && \
    npm cache clean --force && rm -rf /root/.npm /tmp/node-compile-cache
# The system ripgrep, as Claude Code on musl needs; the image is the update, not the agent.
ENV USE_BUILTIN_RIPGREP=0 \
    DISABLE_AUTOUPDATER=1 \
    SHELL=/bin/bash
USER 65532:65532
