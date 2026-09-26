// The dashboard's end of the channel to flai on the host (ADR-0029). flai
// opens a WebSocket to /agent and the two speak JSON-RPC 2.0; the dashboard
// never connects to the host. The agent credential is never sent: each side
// proves it holds it over two nonces. Until the second proof arrives nothing
// is served, and afterwards everything the dashboard can ask for is a method
// flai chose to offer.
//
// One dashboard serves every project the host flai serves (S-0080): every
// project's flai proves itself with the same shared credential, so the proof
// is done once, at the registry, before either side knows which project a
// connection is for. Only once hello names a project does a connection join
// (or replace) that project's own AgentHub: "one connection per credential"
// becomes "one connection per project", which is what a dashboard with more
// than one project connected needs.
import { AsyncLocalStorage } from 'node:async_hooks';
import { createHmac, randomBytes, timingSafeEqual } from 'node:crypto';
import { EventEmitter } from 'node:events';
import { readFileSync } from 'node:fs';
import type { IncomingMessage } from 'node:http';
import type { Duplex } from 'node:stream';
import { WebSocketServer, type WebSocket } from 'ws';
import { log } from './log';
import { version } from './metrics';

export const AGENT_PATH = '/agent';
export const PROTOCOL = 1;
/** Where a flai command's output was capped (flai.ts maxBuffer). */
export const MAX_MESSAGE = 16 * 1024 * 1024;

export class AgentError extends Error {
	constructor(
		public status: number,
		message: string,
		public code?: number,
		/** What flai sent with the error: a conflict's current version, a refusal's findings. */
		public data?: Record<string, unknown>
	) {
		super(message);
	}
}

/**
 * Every method this dashboard asks of flai. The image carries no flai of its own (S-0075), so the
 * flai on the host can be older than the dashboard: it says what it offers when it connects, and
 * what is missing here is told to the designer at once, not discovered a page at a time.
 */
export const REQUIRED_METHODS = [
	'project.info',
	'board.get',
	'items.list',
	'item.get',
	'threads.list',
	'docs.tree',
	'doc.get',
	'adrs.list',
	'activity.get',
	'inbox.designer',
	'search.query',
	'item.move',
	'item.move.preview',
	'item.order',
	'item.block',
	'item.unblock',
	'item.new',
	'item.template',
	'item.show',
	'item.edit',
	'accept.preview',
	'accept.run',
	'stream.diff',
	'stream.log',
	'stream.answer',
	'thread.new',
	'thread.reply',
	'thread.resolve',
	'doc.show',
	'doc.save',
	'adr.new',
	'adr.template',
	'adr.accept',
	'stats.get',
	'push.pending',
	'agent.status',
	'push.run',
	'publish.preview',
	'publish.run',
	'dashboard.status',
	'dashboard.check',
	'dashboard.restart',
	'dashboard.upgrade',
	'dashboard.stop',
	'checks.status',
	'checks.tail',
	'checks.run',
	'checks.cancel',
	// S-0106: flai host, reached through flai serve
	'host.status',
	'host.check',
	'host.process',
	'host.upgrade',
	// S-0105: the host's settings, gated by the settings host action
	'settings.get',
	'settings.action',
	'settings.default_agent',
	'settings.agent',
	'settings.harness',
	'settings.check',
	'settings.checks_timeout',
	'settings.import',
	'settings.mcp_token',
	'settings.dashboard_token',
	// S-0116: a new agent for a story whose agent dropped or failed, gated by the agent host action
	'agent.restart',
	// S-0115: a ready story's agent now, gated by the agent host action
	'agent.start'
];

export type AgentStatus = {
	/** false when the dashboard was given no agent credential (an older flai started it). */
	configured: boolean;
	connected: boolean;
	since?: string;
	flai?: string;
	/** The project the connected flai serves. Not called project: respond() stamps that on every body. */
	serves?: { key: string; name: string };
	/** Methods this dashboard needs and the connected flai does not offer: it is older than the dashboard. */
	missing?: string[];
	/** A repository offered for import, not a project yet (S-0098): it answers only import.preview and
	 * import.run. */
	candidate?: boolean;
};

/** One project as the registry knows it, for a project list or switcher. */
export type ConnectedProject = {
	key: string;
	name: string;
	connected: boolean;
	since?: string;
	/** Offered for import, not a project yet (S-0098). */
	candidate?: boolean;
};

type Pending = {
	resolve: (v: unknown) => void;
	reject: (e: Error) => void;
	timer: ReturnType<typeof setTimeout>;
	onProgress?: (value: unknown) => void;
};
type Message = {
	jsonrpc?: string;
	id?: number | string | null;
	method?: string;
	params?: unknown;
	result?: unknown;
	error?: { code: number; message: string; data?: Record<string, unknown> };
};
type ConnInfo = {
	since: string;
	flai: string;
	project: { key: string; name: string };
	missing: string[];
	candidate?: boolean;
};

export type AgentOptions = { pingMs?: number; handshakeMs?: number; maxUnproven?: number };

export function proof(key: string, role: string, first: string, second: string): string {
	return createHmac('sha256', key).update(`${role}|${first}|${second}`).digest('hex');
}

function same(a: string, b: string): boolean {
	const x = Buffer.from(a);
	const y = Buffer.from(b);
	return x.length === y.length && timingSafeEqual(x, y);
}

/**
 * One project's connection: at most one proven flai at a time, and the requests waiting on it.
 * Emits 'connected' when a flai has proven itself, 'gone' when it is lost, and 'change' with a
 * repo-relative path when flai says a file of the project changed (S-0073). An AgentHub never does
 * its own handshake: a registry proves the shared credential and hands it a socket already proven,
 * so "no host flai has ever named this project" (unknown) can be told apart from "flai is not
 * connected right now" (this hub exists, `status().connected` is false).
 */
export class AgentHub extends EventEmitter {
	private conn: WebSocket | null = null;
	private info: ConnInfo | null = null;
	private pending = new Map<number, Pending>();
	private nextId = 1;
	private pingMs: number;
	/** Whether the registry that made this hub was given a shared credential at all: false only when
	 * the dashboard was started with none, which is true for every project alike, not a per-project fact. */
	private configuredFlag: boolean;
	/** Whether the last connection this hub adopted was a candidate (S-0098); kept once it goes. */
	candidate = false;
	/** Whether its flai said the project is no longer served (S-0118): unregistered with flai serve
	 * project remove or flai dashboard stop. A new connection for the key clears it. */
	removed = false;

	constructor(opt: AgentOptions & { configured?: boolean } = {}) {
		super();
		this.pingMs = opt.pingMs ?? 4000;
		this.configuredFlag = opt.configured ?? true;
	}

	status(): AgentStatus {
		if (!this.configuredFlag) return { configured: false, connected: false };
		if (!this.conn || !this.info) return { configured: true, connected: false };
		const { project, ...rest } = this.info;
		return { configured: true, connected: true, ...rest, serves: project };
	}

	/** One connection per project: a newer proven connection replaces the older one. */
	adopt(ws: WebSocket, info: ConnInfo): void {
		if (this.conn) {
			this.conn.close(4000, 'replaced by a newer connection');
			this.drop(this.conn, 'replaced');
		}
		this.conn = ws;
		this.info = info;
		this.candidate = info.candidate === true;
		this.removed = false;
		log().info(
			{ component: 'agent', flai: info.flai, project: info.project.key },
			'host flai connected'
		);
		let alive = true;
		ws.on('pong', () => (alive = true));
		const ping = setInterval(() => {
			if (!alive) return ws.terminate();
			alive = false;
			ws.ping();
		}, this.pingMs);
		ws.on('message', (data: Buffer) => this.settle(data));
		ws.on('close', () => {
			clearInterval(ping);
			this.drop(ws, 'closed');
		});
		this.emit('connected');
	}

	private drop(ws: WebSocket, why: string): void {
		if (this.conn !== ws) return;
		this.conn = null;
		this.info = null;
		for (const [id, p] of this.pending) {
			clearTimeout(p.timer);
			this.pending.delete(id);
			p.reject(new AgentError(502, 'the host flai went away before it answered'));
		}
		log().info({ component: 'agent', why }, 'host flai gone');
		this.emit('gone');
	}

	private settle(data: Buffer): void {
		let m: Message;
		try {
			m = JSON.parse(data.toString());
		} catch {
			return;
		}
		if (m.method === 'change' && m.id === undefined) {
			const path = (m.params as { path?: unknown } | undefined)?.path;
			if (typeof path === 'string' && path) this.emit('change', path);
			return;
		}
		// flai serve no longer serves the project, and is about to close (S-0118)
		if (m.method === 'removed' && m.id === undefined) {
			this.removed = true;
			return;
		}
		// A step of a request still being answered (an acceptance), by the request's ID.
		if (m.method === '$/progress' && m.id === undefined) {
			const { id, value } = (m.params ?? {}) as { id?: unknown; value?: unknown };
			if (typeof id === 'number') this.pending.get(id)?.onProgress?.(value);
			return;
		}
		if (typeof m.id !== 'number') return;
		const p = this.pending.get(m.id);
		if (!p) return;
		clearTimeout(p.timer);
		this.pending.delete(m.id);
		if (m.error) p.reject(new AgentError(502, m.error.message, m.error.code, m.error.data));
		else p.resolve(m.result);
	}

	/** Ask the host flai for a named method. 503 when none is connected, 504 when it does not answer. */
	ask<T>(
		method: string,
		params: Record<string, unknown> = {},
		timeoutMs = 15000,
		onProgress?: (value: unknown) => void
	): Promise<T> {
		const ws = this.conn;
		const info = this.info;
		if (!ws || !info)
			return Promise.reject(
				new AgentError(
					503,
					'no host flai is connected; run flai dashboard in the project, or flai serve start'
				)
			);
		const id = this.nextId++;
		return new Promise<T>((resolve, reject) => {
			const timer = setTimeout(() => {
				this.pending.delete(id);
				ws.send(JSON.stringify({ jsonrpc: '2.0', method: '$/cancel', params: { id } }));
				reject(new AgentError(504, `the host flai did not answer ${method} in time`));
			}, timeoutMs);
			this.pending.set(id, {
				resolve: resolve as (v: unknown) => void,
				reject,
				timer,
				onProgress
			});
			ws.send(
				JSON.stringify({
					jsonrpc: '2.0',
					id,
					method,
					params: { project: info.project.key, ...params }
				})
			);
		});
	}

	close(): void {
		this.conn?.terminate();
	}
}

/** Rejects everything at once: a project no host flai has ever named here (S-0080). Not kept in the
 * registry's map, so it never counts as a known project and is made fresh for every such request. */
class UnknownProjectHub extends AgentHub {
	constructor(private key: string) {
		super();
	}
	override status(): AgentStatus {
		return { configured: false, connected: false };
	}
	override ask<T>(): Promise<T> {
		return Promise.reject(new AgentError(404, `no project "${this.key}" is being served here`));
	}
}

/**
 * Proves the one shared credential every project's flai holds, then routes the proven connection to
 * that project's own AgentHub by the key hello names (S-0080). A dashboard with nothing configured,
 * or asked for a project it has never heard from, answers as AgentHub always has: 503 unconfigured,
 * or (new) 404 unknown.
 *
 * Every hub's 'connected', 'gone', and 'change' are re-emitted here with the project's key first
 * (S-0095), so a listener can follow a project that has not connected yet, or follow "the one
 * project" before there is one, instead of binding to whichever hub it happened to resolve first.
 */
export class AgentRegistry extends EventEmitter {
	private wss = new WebSocketServer({ noServer: true, maxPayload: MAX_MESSAGE });
	private hubs = new Map<string, AgentHub>();
	private unproven = 0;
	private handshakeMs: number;
	private maxUnproven: number;
	private opt: AgentOptions;

	constructor(
		private key: string | null,
		opt: AgentOptions = {}
	) {
		super();
		this.setMaxListeners(0); // one Repo per project listens, however many projects there are
		this.opt = opt;
		this.handshakeMs = opt.handshakeMs ?? 5000;
		this.maxUnproven = opt.maxUnproven ?? 4;
	}

	/** The project's hub, made the first time it is asked for. Every project shares the credential,
	 * so any key is accepted here; a request naming one nothing has ever connected for uses `peek`
	 * instead, to tell that apart from one that is merely not connected right now. */
	hub(key: string): AgentHub {
		let h = this.hubs.get(key);
		if (!h) {
			h = new AgentHub({ ...this.opt, configured: this.key !== null });
			this.hubs.set(key, h);
			h.on('connected', () => this.emit('connected', key));
			h.on('gone', () => this.emit('gone', key));
			h.on('change', (path: string) => this.emit('change', key, path));
		}
		return h;
	}

	/** The project's hub if a flai has ever named it here, else undefined: never creates one. */
	peek(key: string): AgentHub | undefined {
		return this.hubs.get(key);
	}

	/** Every project a flai has named here, connected or not, for a list or a switcher. */
	list(): ConnectedProject[] {
		const out: ConnectedProject[] = [];
		for (const [key, h] of this.hubs) {
			const st = h.status();
			// a candidate that is gone was imported, or its folder is no longer named, and a project
			// flai said was removed is no longer served: neither is a project to show (S-0118)
			if ((h.candidate || h.removed) && !st.connected) continue;
			out.push({
				key,
				name: st.serves?.name ?? key,
				connected: st.connected,
				since: st.since,
				...(h.candidate ? { candidate: true } : {})
			});
		}
		return out.sort((a, b) => a.key.localeCompare(b.key));
	}

	configured(): boolean {
		return this.key !== null;
	}

	private pending: AgentHub | null = null;

	/**
	 * The one project this dashboard serves, when there is exactly one, for code with no project of
	 * its own to name (S-0080): a call with no `?project=` at all keeps working the way it always did
	 * for the overwhelmingly common single-project dashboard, with nothing to configure for it. With
	 * none known yet, or more than one, a placeholder that is never connected stands in: `configured`
	 * still reflects whether the dashboard holds a credential at all, and once a second project is
	 * seen every caller must start naming one, which the client does once it learns there are several.
	 */
	solo(): AgentHub {
		// a repository offered for import is not a project a request with no ?project= could mean
		const known = [...this.hubs.entries()]
			.filter(([, h]) => !h.candidate && !(h.removed && !h.status().connected))
			.map(([k]) => k);
		if (known.length === 1) return this.hubs.get(known[0])!;
		if (!this.pending) this.pending = new AgentHub({ configured: this.key !== null });
		return this.pending;
	}

	/**
	 * The HTTP server's 'upgrade' event for /agent. A browser always sends Origin on a WebSocket
	 * and flai never does, so a page in the designer's browser cannot reach this at all.
	 */
	handleUpgrade(req: IncomingMessage, socket: Duplex, head: Buffer): void {
		const refuse = (status: string) => {
			socket.write(`HTTP/1.1 ${status}\r\nConnection: close\r\n\r\n`);
			socket.destroy();
		};
		const path = (req.url ?? '').split('?')[0];
		if (path !== AGENT_PATH) return refuse('404 Not Found');
		if (!this.key) return refuse('503 Service Unavailable');
		if (req.headers.origin) return refuse('403 Forbidden');
		if (this.unproven >= this.maxUnproven) return refuse('429 Too Many Requests');
		this.wss.handleUpgrade(req, socket, head, (ws) => this.handshake(ws));
	}

	private handshake(ws: WebSocket): void {
		const key = this.key!;
		this.unproven++;
		let mine = '';
		let theirs = '';
		let hello: {
			flai?: string;
			project?: { key?: string; name?: string };
			methods?: unknown;
			kind?: unknown;
		} = {};
		let done = false;
		const finish = () => {
			if (!done) {
				done = true;
				this.unproven--;
				clearTimeout(timer);
			}
		};
		const timer = setTimeout(() => ws.close(4408, 'handshake timed out'), this.handshakeMs);
		ws.on('close', finish);
		ws.on('error', () => ws.terminate());
		const onMessage = (data: Buffer) => {
			let m: Message;
			try {
				m = JSON.parse(data.toString());
			} catch {
				return ws.close(4400, 'not JSON');
			}
			if (!mine) {
				const p = (m.params ?? {}) as { protocol?: number; nonce?: string } & typeof hello;
				if (m.method !== 'hello' || typeof p.nonce !== 'string' || !p.nonce)
					return ws.close(4400, 'hello first');
				if (p.protocol !== PROTOCOL) {
					ws.send(
						JSON.stringify({
							jsonrpc: '2.0',
							id: m.id,
							error: {
								code: -32000,
								message: `this dashboard speaks protocol ${PROTOCOL}, flai sent ${p.protocol}`
							}
						})
					);
					return ws.close(4400, 'protocol');
				}
				theirs = p.nonce;
				hello = p;
				mine = randomBytes(16).toString('hex');
				ws.send(
					JSON.stringify({
						jsonrpc: '2.0',
						id: m.id,
						result: {
							nonce: mine,
							proof: proof(key, 'dashboard', theirs, mine),
							dashboard: version
						}
					})
				);
				return;
			}
			const given = (m.params as { proof?: string } | undefined)?.proof;
			if (
				m.method !== 'hello.prove' ||
				typeof given !== 'string' ||
				!same(given, proof(key, 'flai', mine, theirs))
			) {
				log().warn({ component: 'agent' }, 'a connection to /agent did not prove the credential');
				return ws.close(4401, 'not proven');
			}
			const projectKey = String(hello.project?.key ?? '');
			if (!projectKey) {
				log().warn({ component: 'agent' }, 'a proven connection named no project');
				finish();
				return ws.close(4400, 'project key required');
			}
			finish();
			ws.off('message', onMessage);
			// A repository offered for import answers only the import methods, by design: nothing
			// of a project's is missing from it (S-0098).
			const candidate = hello.kind === 'candidate';
			this.hub(projectKey).adopt(ws, {
				since: new Date().toISOString(),
				flai: String(hello.flai ?? ''),
				project: { key: projectKey, name: String(hello.project?.name ?? '') },
				missing: candidate
					? []
					: REQUIRED_METHODS.filter(
							(m) => !(Array.isArray(hello.methods) ? hello.methods : []).includes(m)
						),
				candidate
			});
		};
		ws.on('message', onMessage);
	}

	close(): void {
		for (const h of this.hubs.values()) h.close();
		this.wss.close();
	}
}

let reg: AgentRegistry | null = null;

/** The process-wide registry. The shared credential comes from FLAIOVER_AGENT_KEY_FILE, read once. */
export function registry(env: Record<string, string | undefined> = process.env): AgentRegistry {
	if (reg) return reg;
	let key: string | null = null;
	const file = env.FLAIOVER_AGENT_KEY_FILE;
	if (file) {
		try {
			key = readFileSync(file, 'utf8').trim() || null;
		} catch (err) {
			log().warn({ component: 'agent', err: String(err) }, 'agent credential unreadable');
		}
	}
	reg = new AgentRegistry(key);
	return reg;
}

/**
 * Which project the current call is for: set for the length of a request by hooks.server.ts, from
 * the URL or a project query parameter (S-0080). Code with no request of its own to draw a project
 * key from (server startup, a test that never sets it) falls back to `defaultProjectKey`, which
 * keeps one dashboard for one project working exactly as it always has.
 */
export const projectContext = new AsyncLocalStorage<{ key: string }>();
export const defaultProjectKey = '__default__';

export function currentProjectKey(): string {
	return projectContext.getStore()?.key ?? defaultProjectKey;
}

/** Run fn with key as the current project for repo()/agent() calls inside it. */
export function withProject<T>(key: string, fn: () => T): T {
	return projectContext.run({ key }, fn);
}

/**
 * The current project's hub. Called with no arguments everywhere in the app, as before S-0080: it
 * resolves the project from the request in progress (or the default, outside one). A project no
 * flai has ever connected for is refused (404) rather than silently starting to track it, unless it
 * is the default, which a single-project dashboard vivifies the way it always has.
 */
export function agent(): AgentHub {
	return agentFor(currentProjectKey());
}

/** The named project's hub, resolved now, whatever request (if any) is in progress. */
export function agentFor(key: string): AgentHub {
	if (key === defaultProjectKey) return registry().solo();
	const h = registry().peek(key);
	return h ?? new UnknownProjectHub(key);
}

/** Whether an event the registry emitted for `from` concerns the project `key` stands for: the
 * default key follows whichever project is the only one known when the event arrives. */
export function concerns(key: string, from: string): boolean {
	if (key !== defaultProjectKey) return key === from;
	const solo = registry().solo();
	return solo === registry().peek(from);
}

/** Every project a flai has named on this dashboard, for a switcher. */
export function connectedProjects(): ConnectedProject[] {
	return registry().list();
}

export function resetAgent(): void {
	reg?.close();
	reg = null;
}

declare global {
	var __flaioverAgentUpgrade:
		((req: IncomingMessage, socket: Duplex, head: Buffer) => void) | undefined;
}

/** Called from the server's init hook: the custom server entry and the dev server find the registry here. */
export function exposeAgentUpgrade(): void {
	globalThis.__flaioverAgentUpgrade = (req, socket, head) =>
		registry().handleUpgrade(req, socket, head);
}

/** Resolves when a flai is connected, or after ms with whether one is. */
export function connectedWithin(hub: AgentHub, ms: number): Promise<boolean> {
	if (hub.status().connected) return Promise.resolve(true);
	return new Promise((resolve) => {
		const done = (ok: boolean) => {
			clearTimeout(timer);
			hub.off('connected', yes);
			resolve(ok);
		};
		const yes = () => done(true);
		const timer = setTimeout(() => done(false), ms);
		hub.on('connected', yes);
	});
}
