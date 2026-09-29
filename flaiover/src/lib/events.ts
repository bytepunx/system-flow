// The project's live events (/api/events) for everything on a page that follows them (S-0154):
// one EventSource per project for the whole tab, opened by the first listener and closed with the
// last, rather than one each. A browser keeps few connections to one origin, and a story page has
// the inbox, the item, its threads, and its agent panel all listening.
import { projectState } from '$lib/project.svelte';
import { affects, type Change, type ChangeKind } from '$lib/changes';

export type Listener = {
	/** A file of the project changed: `path` is repo-relative, `kind` what the server read it as (S-0161). */
	change?: (path: string, kind: ChangeKind) => void;
	/** A story's agent started or ended (flai serve's `agent`), which changes no file. */
	agent?: (story: string) => void;
};

type Stream = { source: EventSource; listeners: Set<Listener> };
const streams = new Map<string, Stream>();

function field(e: Event | undefined, name: string): string {
	try {
		const v = JSON.parse((e as MessageEvent | undefined)?.data ?? '{}')[name];
		return typeof v === 'string' ? v : '';
	} catch {
		return '';
	}
}

/** Follows the current project's events until the returned function is called; nothing where
 * there is no EventSource to follow them with. */
export function listen(l: Listener): () => void {
	if (typeof EventSource === 'undefined') return () => {};
	const url = projectState.tag('/api/events');
	let s = streams.get(url);
	if (!s) {
		const source = new EventSource(url);
		const stream: Stream = { source, listeners: new Set() };
		source.addEventListener('change', (e) => {
			const path = field(e, 'path');
			const kind = (field(e, 'kind') || 'other') as ChangeKind;
			for (const x of [...stream.listeners]) x.change?.(path, kind);
		});
		source.addEventListener('agent', (e) => {
			for (const x of [...stream.listeners]) x.agent?.(field(e, 'story'));
		});
		streams.set(url, stream);
		s = stream;
	}
	const mine = s;
	mine.listeners.add(l);
	return () => {
		if (!mine.listeners.delete(l) || mine.listeners.size) return;
		mine.source.close();
		if (streams.get(url) === mine) streams.delete(url);
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
