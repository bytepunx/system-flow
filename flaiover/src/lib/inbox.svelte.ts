// The designer's inbox in the browser (S-0042): one fetch shared by the
// navigation badge and the inbox page, refreshed on repository change events,
// with optional desktop notifications for entries that appear while open.
import { api } from '$lib/api';

export type InboxEntry = {
	key: string;
	kind: 'thread' | 'question' | 'review' | 'blocked' | 'overlap';
	title: string;
	detail?: string;
	href: string;
	at?: string;
};
export type Inbox = {
	total: number;
	counts: Record<InboxEntry['kind'], number>;
	entries: InboxEntry[];
	notes: string[];
};

/** Entries whose key was not known before; nothing is new on the first look (known is null). */
export function newEntries(known: Set<string> | null, entries: InboxEntry[]): InboxEntry[] {
	return known === null ? [] : entries.filter((e) => !known.has(e.key));
}

const NOTIFY_KEY = 'flaiover.inbox.notify';

class InboxState {
	data = $state<Inbox | null>(null);
	error = $state<string | null>(null);
	/** The designer's choice in this browser; off by default. */
	notify = $state(false);
	/** "granted", "denied", "default", or "unsupported". */
	permission = $state<string>('default');
	#known: Set<string> | null = null;
	#source: EventSource | null = null;
	#timer: ReturnType<typeof setTimeout> | null = null;

	start(): void {
		if (this.#source || typeof window === 'undefined') return;
		this.permission = 'Notification' in window ? Notification.permission : 'unsupported';
		this.notify = localStorage.getItem(NOTIFY_KEY) === 'on' && this.permission === 'granted';
		void this.refresh();
		this.#source = new EventSource('/api/events');
		this.#source.addEventListener('change', () => {
			// a save touches several files; one refresh after they settle
			if (this.#timer) clearTimeout(this.#timer);
			this.#timer = setTimeout(() => void this.refresh(), 300);
		});
	}

	stop(): void {
		this.#source?.close();
		this.#source = null;
	}

	async refresh(): Promise<void> {
		try {
			const r = await api('/api/inbox');
			if (!r.ok) throw new Error((await r.json().catch(() => ({}))).error ?? r.statusText);
			const next = (await r.json()) as Inbox;
			const fresh = newEntries(this.#known, next.entries);
			this.#known = new Set(next.entries.map((e) => e.key));
			this.data = next;
			this.error = null;
			if (this.notify && this.permission === 'granted')
				for (const e of fresh)
					new Notification('flaiover: ' + e.title, { body: e.detail ?? e.kind, tag: e.key });
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
		}
	}

	/** Turn desktop notifications on or off; asking the browser for permission when needed. */
	async setNotify(on: boolean): Promise<void> {
		if (on && this.permission === 'default' && 'Notification' in window)
			this.permission = await Notification.requestPermission();
		this.notify = on && this.permission === 'granted';
		localStorage.setItem(NOTIFY_KEY, this.notify ? 'on' : 'off');
	}
}

export const inboxState = new InboxState();
