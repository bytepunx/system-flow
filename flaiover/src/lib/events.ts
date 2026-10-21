// The project's live events (/api/events) for everything on a page that follows them (S-0154):
// one WebSocket per project for the whole tab, opened by the first listener and closed with the
// last, rather than one each (ADR-0182). A story page has the inbox, the item, its threads, and its
// agent panel all listening. A socket that drops is opened again after a backoff while anyone
// listens, and since a new socket does not resume where the dropped one left off, every listener
// is then told to ask again for what it shows.
import { projectState } from '$lib/project.svelte';
import { affects, MANIFEST, type Change, type ChangeKind } from '$lib/changes';

export type Listener = {
	/** A file of the project changed: `path` is repo-relative, `kind` what the server read it as (S-0161). */
	change?: (path: string, kind: ChangeKind) => void;
	/** A story's agent, the planner for an epic or a story, or the orchestrator, started or
	 * ended (flai serve's `agent`), which changes no file. */
	agent?: (story: string) => void;
};

/** The first wait before opening a dropped socket again, doubled after each drop until a `ready`. */
export const RETRY_MS = 1000;
/** The longest wait between attempts to open a dropped socket again. */
export const RETRY_MAX_MS = 30000;

type Stream = {
	url: string;
	listeners: Set<Listener>;
	socket: WebSocket | null;
	retry: ReturnType<typeof setTimeout> | null;
	delay: number;
	/** The socket open now replaces one that dropped, so its `ready` is a reconnection. */
	reconnecting: boolean;
};
const streams = new Map<string, Stream>();

/** The socket URL of a same-origin path, ws: or wss: as the page's own protocol. */
function socketURL(path: string): string {
	return `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}${path}`;
}

function text(v: unknown): string {
	return typeof v === 'string' ? v : '';
}

/** Every listener asks again for all it shows: the manifest changing affects every kind
 * (changes.ts), and an agent with no story is any agent. */
function askAgain(s: Stream) {
	for (const x of [...s.listeners]) {
		x.change?.(MANIFEST, 'project');
		x.agent?.('');
	}
}

function receive(s: Stream, data: unknown) {
	let m: { type?: unknown; path?: unknown; kind?: unknown; story?: unknown };
	try {
		m = JSON.parse(text(data)) ?? {};
	} catch {
		return;
	}
	if (m.type === 'ready') {
		s.delay = RETRY_MS;
		if (s.reconnecting) {
			s.reconnecting = false;
			askAgain(s);
		}
	} else if (m.type === 'change') {
		const path = text(m.path);
		const kind = (text(m.kind) || 'other') as ChangeKind;
		for (const x of [...s.listeners]) x.change?.(path, kind);
	} else if (m.type === 'agent') {
		for (const x of [...s.listeners]) x.agent?.(text(m.story));
	}
}

function connect(s: Stream) {
	const socket = new WebSocket(s.url);
	s.socket = socket;
	socket.onmessage = (e) => receive(s, e.data);
	socket.onclose = () => {
		// the last listener leaving closes it on purpose and forgets it first
		if (s.socket !== socket) return;
		s.socket = null;
		s.reconnecting = true;
		s.retry = setTimeout(() => {
			s.retry = null;
			connect(s);
		}, s.delay);
		s.delay = Math.min(s.delay * 2, RETRY_MAX_MS);
	};
}

function disconnect(s: Stream) {
	if (s.retry) clearTimeout(s.retry);
	s.retry = null;
	const socket = s.socket;
	s.socket = null;
	socket?.close();
}

/** Follows the current project's events until the returned function is called; nothing where
 * there is no WebSocket to follow them with. */
export function listen(l: Listener): () => void {
	if (typeof WebSocket === 'undefined' || typeof location === 'undefined') return () => {};
	const path = projectState.tag('/api/events');
	let s = streams.get(path);
	if (!s) {
		const stream: Stream = {
			url: socketURL(path),
			listeners: new Set(),
			socket: null,
			retry: null,
			delay: RETRY_MS,
			reconnecting: false
		};
		connect(stream);
		streams.set(path, stream);
		s = stream;
	}
	const mine = s;
	mine.listeners.add(l);
	return () => {
		if (!mine.listeners.delete(l) || mine.listeners.size) return;
		disconnect(mine);
		if (streams.get(path) === mine) streams.delete(path);
	};
}

/** How long a page waits for changes that arrive together before it asks again (S-0161). */
export const GATHER_MS = 500;
/** The longest a burst of changes that never settles holds a page back. */
export const GATHER_MAX_MS = 2000;

/**
 * Calls f once calls have settled for ms, with every value they were given (S-0161): a save or an
 * agent's commit touches several files. A burst that never settles still calls f every maxMs.
 * Cancel with stop().
 */
export function debounced<T = unknown>(
	f: (gathered: T[]) => void,
	ms = GATHER_MS,
	maxMs = GATHER_MAX_MS
): { (x?: T): void; stop: () => void } {
	let timer: ReturnType<typeof setTimeout> | null = null;
	let first = 0;
	let gathered: T[] = [];
	const fire = () => {
		timer = null;
		const xs = gathered;
		gathered = [];
		f(xs);
	};
	const call = (x?: T) => {
		if (x !== undefined) gathered.push(x);
		const now = Date.now();
		if (timer) clearTimeout(timer);
		else first = now;
		timer = setTimeout(fire, Math.max(0, Math.min(ms, first + maxMs - now)));
	};
	call.stop = () => {
		if (timer) clearTimeout(timer);
		timer = null;
		gathered = [];
	};
	return call;
}

/**
 * Follows the current project's changes of the kinds a page shows (S-0161): f is called once for
 * the changes that arrive within GATHER_MS of each other, and only when one of them is of kinds, of
 * the manifest, or of a file the server could not tell apart. Stop with the returned function.
 */
export function follow(kinds: ChangeKind[], f: (changes: Change[]) => void): () => void {
	const later = debounced<Change>(f);
	const stop = listen({
		change: (path, kind) => {
			const c = { path, kind };
			if (affects([c], ...kinds)) later(c);
		}
	});
	return () => {
		stop();
		later.stop();
	};
}
