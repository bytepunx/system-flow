// The project's live events (/api/events) for everything on a page that follows them (S-0154):
// one EventSource per project for the whole tab, opened by the first listener and closed with the
// last, rather than one each. A browser keeps few connections to one origin, and a story page has
// the inbox, the item, its threads, and its agent panel all listening.
import { projectState } from '$lib/project.svelte';

export type Listener = {
	/** A file of the project changed: `path` is repo-relative. */
	change?: (path: string) => void;
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
			for (const x of [...stream.listeners]) x.change?.(field(e, 'path'));
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

/** Calls f once changes have settled for ms: a save touches several files. Cancel with stop(). */
export function debounced(f: () => void, ms = 300): { (): void; stop: () => void } {
	let timer: ReturnType<typeof setTimeout> | null = null;
	const call = () => {
		if (timer) clearTimeout(timer);
		timer = setTimeout(() => {
			timer = null;
			f();
		}, ms);
	};
	call.stop = () => {
		if (timer) clearTimeout(timer);
		timer = null;
	};
	return call;
}
