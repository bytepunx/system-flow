// The dashboard's live events over a WebSocket (ADR-0182): an upgrade of GET /api/events?project=<key>
// that sends what the project's Repo announces, so that a tab following a project holds none of the
// browser's six HTTP/1.1 connections to the dashboard. Upgrades bypass SvelteKit's handle, so the
// upgrade is authenticated and authorized here, with the functions hooks.server.ts uses for a GET of
// /api/events, before the socket opens.
import type { IncomingMessage } from 'node:http';
import type { Duplex } from 'node:stream';
import { WebSocketServer, type WebSocket } from 'ws';
import { withProject } from './agent';
import { allowedHosts, allowsOrigin } from './allowed-hosts';
import { authenticate, decide } from './auth';
import { log } from './log';
import { repo, type Repo } from './repo';

export const EVENTS_PATH = '/api/events';
/** The close code every events socket gets when the server stops. */
export const CLOSE_GOING_AWAY = 1001;
/** How often a socket is pinged; one that has not answered the ping before is closed. */
export const PING_MS = 25000;

/** What the server sends on an events socket, one JSON text message each. */
export type EventMessage =
	| { type: 'ready' }
	| { type: 'change'; path: string; kind?: string }
	| { type: 'agent'; story: string };

/** The request's headers as a fetch Headers, which authenticate and decide read. */
function headersOf(req: IncomingMessage): Headers {
	const headers = new Headers();
	for (const [name, value] of Object.entries(req.headers)) {
		if (value === undefined) continue;
		headers.set(name, Array.isArray(value) ? value.join(', ') : value);
	}
	return headers;
}

/**
 * Whether a browser's Origin names this request's own host, or one of the allowed hosts a tunnel or a
 * domain reaches the dashboard by (ADR-0193); a request with none is not a browser's.
 */
function sameOrigin(headers: Headers, allowed: readonly string[]): boolean {
	const origin = headers.get('origin');
	if (!origin) return true;
	let url: URL;
	try {
		url = new URL(origin);
	} catch {
		return false;
	}
	return url.host === headers.get('host')?.toLowerCase() || allowsOrigin(url, allowed);
}

/** The upgrade of /api/events and every socket it opened, closed together when the server stops. */
export class EventsSocket {
	private wss = new WebSocketServer({ noServer: true, maxPayload: 1024 });
	/** Each open socket with its teardown: its ping cleared, its Repo listeners off, and it forgotten. */
	private open = new Map<WebSocket, () => void>();
	private stopping = false;

	/** allowed is the hosts an Origin may name beside the request's own, FLAIOVER_ALLOWED_HOSTS's. */
	constructor(
		private pingMs = PING_MS,
		private allowed: readonly string[] = allowedHosts()
	) {}

	/** The HTTP server's 'upgrade' event for /api/events: refused on the raw socket, or opened. */
	upgrade(req: IncomingMessage, socket: Duplex, head: Buffer): void {
		const refuse = (status: string) => {
			log().info({ component: 'events', path: EVENTS_PATH, status }, 'events upgrade refused');
			if (!socket.destroyed) socket.write(`HTTP/1.1 ${status}\r\nConnection: close\r\n\r\n`);
			socket.destroy();
		};
		if (this.stopping) return refuse('503 Service Unavailable');
		const headers = headersOf(req);
		const decision = decide(EVENTS_PATH, 'GET', authenticate(headers), headers);
		if (decision.kind === 'forbidden') return refuse('403 Forbidden');
		if (decision.kind !== 'allow') return refuse('401 Unauthorized');
		// a browser sends the dashboard's cookie on a cross-site WebSocket too, and a GET has no CSRF rule
		if (!sameOrigin(headers, this.allowed)) return refuse('403 Forbidden');
		// the project as hooks.server.ts takes it for a request: the query parameter, else the default
		const project = new URL(req.url ?? '/', 'http://upgrade').searchParams.get('project');
		const r = project ? withProject(project, repo) : repo();
		r.watch().then(
			() => this.wss.handleUpgrade(req, socket, head, (ws) => this.serve(ws, r, project)),
			(err) => {
				log().error(
					{ component: 'events', project, err: String(err) },
					'events upgrade failed: the project could not be watched'
				);
				refuse('500 Internal Server Error');
			}
		);
	}

	/**
	 * Close every events socket with 1001 and refuse new ones, for a server that is stopping. Each is
	 * torn down first, so that no Repo listener outlives the call whenever its close handshake ends: an
	 * upgraded socket is detached from the HTTP server, whose close does not wait for it (I-0205).
	 */
	close(): void {
		this.stopping = true;
		for (const [ws, teardown] of [...this.open]) {
			teardown();
			ws.close(CLOSE_GOING_AWAY, 'server stopping');
		}
	}

	private serve(ws: WebSocket, r: Repo, project: string | null): void {
		if (this.stopping) return void ws.close(CLOSE_GOING_AWAY, 'server stopping');
		const send = (m: EventMessage) => {
			if (ws.readyState === ws.OPEN) ws.send(JSON.stringify(m));
		};
		const change = (path: string, kind?: string) => send({ type: 'change', path, kind });
		const agent = (story: string) => send({ type: 'agent', story });
		let answered = true;
		const ping = setInterval(() => {
			if (!answered) return void ws.terminate();
			answered = false;
			ws.ping();
		}, this.pingMs);
		ws.on('pong', () => (answered = true));
		ws.on('error', (err) =>
			log().debug({ component: 'events', project, err: String(err) }, 'events socket error')
		);
		const teardown = () => {
			clearInterval(ping);
			r.off('change', change);
			r.off('agent', agent);
			this.open.delete(ws);
		};
		this.open.set(ws, teardown);
		ws.on('close', teardown);
		r.on('change', change);
		r.on('agent', agent);
		send({ type: 'ready' });
	}
}
