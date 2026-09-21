// The dashboard's end of the channel to flai on the host (ADR-0029). flai
// opens a WebSocket to /agent and the two speak JSON-RPC 2.0; the dashboard
// never connects to the host. The agent credential is never sent: each side
// proves it holds it over two nonces. Until the second proof arrives nothing
// is served, and afterwards everything the dashboard can ask for is a method
// flai chose to offer.
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
	'publish.run'
];

export type AgentStatus = {
	/** false when the container was given no agent credential (an older flai started it). */
	configured: boolean;
	connected: boolean;
	since?: string;
	flai?: string;
	/** The project the connected flai serves. Not called project: respond() stamps that on every body. */
	serves?: { key: string; name: string };
	/** Methods this dashboard needs and the connected flai does not offer: it is older than the dashboard. */
	missing?: string[];
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
 * Emits 'connected' when a flai has proven itself, 'gone' when it is lost, and 'change' with a
 * repo-relative path when flai says a file of the project changed (S-0073).
 */
export class AgentHub extends EventEmitter {
	private wss = new WebSocketServer({ noServer: true, maxPayload: MAX_MESSAGE });
	private conn: WebSocket | null = null;
	private info: {
		since: string;
		flai: string;
		project: { key: string; name: string };
		missing: string[];
	} | null = null;
	private pending = new Map<number, Pending>();
	private nextId = 1;
	private unproven = 0;
	private pingMs: number;
	private handshakeMs: number;
	private maxUnproven: number;

	constructor(
		private key: string | null,
		opt: AgentOptions = {}
	) {
		super();
		this.pingMs = opt.pingMs ?? 4000;
		this.handshakeMs = opt.handshakeMs ?? 5000;
		this.maxUnproven = opt.maxUnproven ?? 4;
	}

	status(): AgentStatus {
		if (!this.key) return { configured: false, connected: false };
		if (!this.conn || !this.info) return { configured: true, connected: false };
		const { project, ...rest } = this.info;
		return { configured: true, connected: true, ...rest, serves: project };
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
		let hello: { flai?: string; project?: { key?: string; name?: string }; methods?: unknown } = {};
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
			finish();
			ws.off('message', onMessage);
			this.adopt(ws, {
				since: new Date().toISOString(),
				flai: String(hello.flai ?? ''),
				project: { key: String(hello.project?.key ?? ''), name: String(hello.project?.name ?? '') },
				missing: REQUIRED_METHODS.filter(
					(m) => !(Array.isArray(hello.methods) ? hello.methods : []).includes(m)
				)
			});
		};
		ws.on('message', onMessage);
	}

	/** One connection per credential: a newer proven connection replaces the older one. */
	private adopt(ws: WebSocket, info: NonNullable<AgentHub['info']>): void {
		if (this.conn) {
			this.conn.close(4000, 'replaced by a newer connection');
			this.drop(this.conn, 'replaced');
		}
		this.conn = ws;
		this.info = info;
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
		this.wss.close();
	}
}

let hub: AgentHub | null = null;

/** The process-wide hub. The credential comes from FLAIOVER_AGENT_KEY_FILE, read once. */
export function agent(env: Record<string, string | undefined> = process.env): AgentHub {
	if (hub) return hub;
	let key: string | null = null;
	const file = env.FLAIOVER_AGENT_KEY_FILE;
	if (file) {
		try {
			key = readFileSync(file, 'utf8').trim() || null;
		} catch (err) {
			log().warn({ component: 'agent', err: String(err) }, 'agent credential unreadable');
		}
	}
	hub = new AgentHub(key);
	return hub;
}

export function resetAgent(): void {
	hub?.close();
	hub = null;
}

declare global {
	var __flaioverAgentUpgrade:
		((req: IncomingMessage, socket: Duplex, head: Buffer) => void) | undefined;
}

/** Called from the server's init hook: the custom server entry and the dev server find the hub here. */
export function exposeAgentUpgrade(): void {
	globalThis.__flaioverAgentUpgrade = (req, socket, head) =>
		agent().handleUpgrade(req, socket, head);
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
