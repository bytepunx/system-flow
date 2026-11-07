// flai-guard hands each OpenCode tool call to `flai guard --harness opencode`
// before it runs, and refuses the call when the guard does, so that an agent
// on OpenCode is held by the same guard as one on Claude Code (S-0407,
// ADR-0130). OpenCode loads every plugin in .opencode/plugins.
//
// A sub-agent's call is told by its session: OpenCode's Task tool runs a
// sub-agent in a child session whose parentID is its caller's, so the plugin
// sends each call with its session's id, parentID, and agent. It also tells
// the guard when a sub-agent's session starts and when its task ends, so that
// the guard keeps the record of running sub-agents (ADR-0092).
//
// The guard answers with its exit status: 2 refuses, with the reason on
// standard error, which the plugin throws so that OpenCode hands it back to
// the model; any other status passes the call. When flai is not on PATH the
// plugin passes every call and says so once, so that a project opened outside
// flai serve still works; in a session flai started (FLAI_STARTED_BY,
// FLAI_STORY, or FLAI_ROLE set) it refuses every call instead, since the
// agent would otherwise run unguarded (ADR-0235).

import { spawn } from "node:child_process";

// COMMAND is the guard's program and the arguments before its own.
const COMMAND = ["flai"];

// REFUSES is the exit status with which flai guard refuses a call.
const REFUSES = 2;

// started says flai started this session, so a missing guard refuses.
const started = Boolean(process.env.FLAI_STARTED_BY || process.env.FLAI_STORY || process.env.FLAI_ROLE);

// run hands input to the guard, as one JSON object on standard input, and
// resolves with its exit status and standard error, or with missing set when
// the guard could not be started.
function run(input, cwd) {
  return new Promise((resolve) => {
    let stderr = "";
    let child;
    try {
      child = spawn(COMMAND[0], [...COMMAND.slice(1), "guard", "--harness", "opencode"], {
        cwd,
        stdio: ["pipe", "ignore", "pipe"],
      });
    } catch (err) {
      resolve({ missing: true, code: -1, stderr: String(err?.message ?? err) });
      return;
    }
    child.on("error", (err) => resolve({ missing: true, code: -1, stderr: String(err?.message ?? err) }));
    child.stderr.on("data", (chunk) => {
      stderr += chunk;
    });
    child.on("close", (code) => resolve({ missing: false, code, stderr }));
    child.stdin.on("error", () => {});
    child.stdin.end(JSON.stringify(input));
  });
}

export const FlaiGuard = async ({ client, directory }) => {
  // sessions are the sessions seen, by id, each with its parentID and agent.
  const sessions = new Map();
  let warned = false;

  // session is the session with an id, from those seen or, the first time,
  // from OpenCode's client; one it cannot read is taken as the agent's own.
  async function session(id) {
    if (!id) return { id: "", parentID: "", agent: "" };
    if (sessions.has(id)) return sessions.get(id);
    let info;
    try {
      const res = await client.session.get({ path: { id }, sessionID: id });
      info = res?.data ?? undefined;
    } catch {
      info = undefined;
    }
    const s = { id, parentID: info?.parentID ?? "", agent: info?.agent ?? "" };
    if (info) sessions.set(id, s);
    return s;
  }

  // missing passes a call when flai is not on PATH, saying so once, or
  // refuses it in a session flai started.
  function missing(why) {
    if (started) {
      throw new Error(`flai guard did not run (${why}): no guarded call passes until flai is installed and on PATH (ADR-0235)`);
    }
    if (!warned) {
      warned = true;
      console.warn(`flai-guard: flai guard did not run (${why}), so every call passes; install flai and put it on PATH to guard this project's agents`);
    }
  }

  // record tells the guard of a sub-agent starting or stopping; it never
  // stops the call it rides on.
  async function record(input) {
    await run(input, directory);
  }

  return {
    "tool.execute.before": async (input, output) => {
      const r = await run(
        {
          event: "tool.execute.before",
          tool: input.tool,
          sessionID: input.sessionID,
          callID: input.callID,
          args: output?.args ?? {},
          session: await session(input.sessionID),
        },
        directory,
      );
      if (r.missing) return missing(r.stderr);
      if (r.code === REFUSES) throw new Error(r.stderr.trim() || "flai guard refused the call");
    },
    "tool.execute.after": async (input, output) => {
      const child = input.tool === "task" ? output?.metadata?.sessionId : undefined;
      if (!child) return;
      const s = await session(child);
      await record({
        event: "tool.execute.after",
        tool: "task",
        sessionID: input.sessionID,
        callID: input.callID,
        session: { ...s, parentID: s.parentID || input.sessionID },
      });
    },
    event: async ({ event }) => {
      if (event?.type !== "session.created") return;
      const info = event.properties?.info;
      if (!info?.id) return;
      const s = { id: info.id, parentID: info.parentID ?? "", agent: info.agent ?? "" };
      sessions.set(info.id, s);
      if (s.parentID) await record({ event: "session.created", sessionID: s.parentID, session: s });
    },
  };
};
