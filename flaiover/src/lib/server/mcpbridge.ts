// MCP over Streamable HTTP, bridged to `flai mcp` (ADR-0024): one process per
// session, JSON-RPC lines on its standard input and output. flaiover carries
// messages and holds none of the rules.
import { spawn, type ChildProcessWithoutNullStreams } from 'node:child_process';
import { randomUUID } from 'node:crypto';
import { join } from 'node:path';
import { flaiBinary } from './flai';
import { log } from './log';
import { RepoError } from './repo';

type Json = Record<string, unknown>;
type Pending = { resolve: (v: Json) => void; reject: (e: Error) => void };

const num = (v: string | undefined, d: number) => (Number(v) > 0 ? Number(v) : d);
export const maxSessions = () => num(process.env.FLAIOVER_MCP_MAX_SESSIONS, 16);
export const idleMs = () => num(process.env.FLAIOVER_MCP_IDLE_MINUTES, 30) * 60_000;

const safeAgent = (s: string) =>
	s
		.replace(/[^A-Za-z0-9._-]+/g, '-')
		.replace(/^[-.]+|[-.]+$/g, '')
		.slice(0, 64);

export class McpSession {
	readonly id = randomUUID();
	private child: ChildProcessWithoutNullStreams;
	private pending = new Map<string, Pending>();
	private buffer = '';
	private timer: ReturnType<typeof setTimeout> | null = null;
	closed = false;

	constructor(
		bin: string,
		projectDir: string,
		readonly agent: string,
		private onClose: (s: McpSession) => void
	) {
		this.child = spawn(bin, ['mcp'], {
			cwd: projectDir,
			env: {
				...process.env,
				FLAI_CONFIG: process.env.FLAI_CONFIG ?? join(projectDir, '.flai-cache', 'config.json'),
				FLAI_CACHE_DIR: process.env.FLAI_CACHE_DIR ?? join(projectDir, '.flai-cache', 'cache'),
				FLAI_AGENT: agent,
				LOG_FORMAT: 'json'
			}
		});
		this.child.stdout.setEncoding('utf8').on('data', (d: string) => this.read(d));
		this.child.stderr.resume(); // flai logs there; stdout carries the protocol only
		this.child.on('error', (e) => this.close(e));
		this.child.on('exit', () => this.close(new Error('flai mcp exited')));
		this.touch();
	}

	private read(chunk: string): void {
		this.buffer += chunk;
		const lines = this.buffer.split('\n');
		this.buffer = lines.pop() ?? '';
		for (const line of lines) {
			if (!line.trim()) continue;
			let msg: Json;
			try {
				msg = JSON.parse(line);
			} catch {
				continue;
			}
			// a response to something we forwarded; anything unprompted has nowhere to go (no GET stream)
			const key = msg.id === undefined || msg.id === null ? null : JSON.stringify(msg.id);
			const waiter = key && !('method' in msg) ? this.pending.get(key) : undefined;
			if (key && waiter) {
				this.pending.delete(key);
				waiter.resolve(msg);
			}
		}
	}

	private touch(): void {
		if (this.timer) clearTimeout(this.timer);
		this.timer = setTimeout(() => this.close(new Error('session idled out')), idleMs());
		this.timer.unref?.();
	}

	/** Forward one JSON-RPC message. A request resolves with its response; anything else with null. */
	send(msg: Json): Promise<Json | null> {
		if (this.closed) return Promise.reject(new Error('session is closed'));
		this.touch();
		const isRequest = typeof msg.method === 'string' && msg.id !== undefined && msg.id !== null;
		const line = JSON.stringify(msg) + '\n';
		if (!isRequest) {
			this.child.stdin.write(line);
			return Promise.resolve(null);
		}
		return new Promise((resolve, reject) => {
			this.pending.set(JSON.stringify(msg.id), {
				resolve: (v) => {
					this.touch();
					resolve(v);
				},
				reject
			});
			this.child.stdin.write(line);
		});
	}

	close(reason?: Error): void {
		if (this.closed) return;
		this.closed = true;
		if (this.timer) clearTimeout(this.timer);
		for (const p of this.pending.values()) p.reject(reason ?? new Error('session closed'));
		this.pending.clear();
		this.child.stdin.end();
		this.child.kill();
		this.onClose(this);
	}
}

const sessions = new Map<string, McpSession>();

export function sessionCount(): number {
	return sessions.size;
}

export function getSession(id: string): McpSession | undefined {
	return sessions.get(id);
}

/** Start a session: its own flai mcp process, named for the agent behind it. */
export async function startSession(projectDir: string, agent: string): Promise<McpSession> {
	const bin = await flaiBinary();
	if (!bin)
		throw new RepoError(503, 'flai is not available to this dashboard, so MCP is not served');
	if (sessions.size >= maxSessions())
		throw new RepoError(
			503,
			`the dashboard already has ${sessions.size} MCP sessions, its limit (FLAIOVER_MCP_MAX_SESSIONS); end one with DELETE /mcp or wait for an idle one to expire`
		);
	const s = new McpSession(bin, projectDir, safeAgent(agent) || 'agent', (closed) => {
		sessions.delete(closed.id);
		log().info(
			{ component: 'mcp', session: closed.id, agent: closed.agent, sessions: sessions.size },
			'mcp session ended'
		);
	});
	sessions.set(s.id, s);
	log().info(
		{ component: 'mcp', session: s.id, agent: s.agent, sessions: sessions.size },
		'mcp session started'
	);
	return s;
}

/** End every session; for shutdown and tests. */
export function closeAllSessions(): void {
	for (const s of [...sessions.values()]) s.close();
}

// A flai mcp process must not outlive the server that started it.
process.once('exit', closeAllSessions);
