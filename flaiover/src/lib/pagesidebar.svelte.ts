// The state behind PageSidebar (S-0460): a page's entries, the one chosen, the pane shown, and
// the entry current. The choice is kept in the URL's hash both ways: a hash the page arrives
// with, or one back and forward bring, chooses its entry, and choosing an entry writes its hash.
// The raw choice is kept as given, so a hash that names a pane the page has not drawn yet chooses
// it once the entries arrive, and a pane that goes away falls back to the page's default. A scroll
// entry is current while its section is in view, which an IntersectionObserver reports; where
// there is none, as under jsdom, the choice is current.

export type SidebarMode = 'pane' | 'scroll';

export type SidebarEntry = {
	/** The entry's id: the URL hash that chooses it, and for a scroll entry its section's element id. */
	id: string;
	label: string;
	/** A count shown beside the label, such as the open threads. */
	count?: number;
	/** pane shows one pane at a time; scroll scrolls to a section on the page or in its pane. */
	mode: SidebarMode;
	/** Entries listed under this one: a scroll entry under a pane is a section of that pane. */
	children?: SidebarEntry[];
};

/** What the page gives the state, read again whenever it changes. */
export type SidebarSource = {
	entries: SidebarEntry[];
	/** The entry chosen when the hash names none; the first entry when this names none either. */
	fallback?: string;
	/** The entry that holds a hash that is not an entry's id, such as a message's anchor. */
	match?: (hash: string) => string | undefined;
};

/** An entry, with the pane it sits in: itself for a pane, its nearest pane above it otherwise. */
export type Placed = { entry: SidebarEntry; pane?: string };

/** Every entry, depth first, each with the pane it sits in. */
export function flatten(entries: SidebarEntry[], pane?: string): Placed[] {
	const out: Placed[] = [];
	for (const entry of entries) {
		const here = entry.mode === 'pane' ? entry.id : pane;
		out.push({ entry, pane: here });
		if (entry.children?.length) out.push(...flatten(entry.children, here));
	}
	return out;
}

/** The id a URL hash names: without its '#', decoded; '' for none. */
export function hashId(hash: string): string {
	const raw = hash.replace(/^#/, '');
	try {
		return decodeURIComponent(raw);
	} catch {
		return raw;
	}
}

/**
 * The entry a choice resolves to: the entry it names, else the entry match says holds it, else
 * the fallback, else the first entry; undefined while there are no entries.
 */
export function resolveChoice(
	id: string | undefined,
	{ entries, fallback, match }: SidebarSource
): string | undefined {
	const all = flatten(entries);
	const known = (x: string | undefined) => !!x && all.some((p) => p.entry.id === x);
	if (known(id)) return id;
	const held = id ? match?.(id) : undefined;
	if (known(held)) return held;
	if (known(fallback)) return fallback;
	return all[0]?.entry.id;
}

/**
 * The pane shown for a resolved choice: the pane it sits in, else the fallback when that sits in
 * a pane, else the first pane; undefined on a page with no panes.
 */
export function paneOf(choice: string | undefined, source: SidebarSource): string | undefined {
	const all = flatten(source.entries);
	const of = (x: string | undefined) => all.find((p) => p.entry.id === x)?.pane;
	return of(choice) ?? of(source.fallback) ?? all.find((p) => p.pane)?.pane;
}

/** How to scroll: at once for a reader who asks for reduced motion, smoothly otherwise. */
export function scrollBehavior(): ScrollBehavior {
	if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return 'auto';
	return window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth';
}

export class PageSidebarState {
	#source: () => SidebarSource;
	/** The choice as given, by a click, the hash, or the page; resolved in choice. */
	chosen = $state<string | undefined>(undefined);
	/** The watched scroll entries whose sections are in view, in entry order. */
	inView = $state<string[]>([]);

	/** source is read again whenever what it reads changes, so give it a closure over props. */
	constructor(source: () => SidebarSource) {
		this.#source = source;
		if (typeof location !== 'undefined') this.chosen = hashId(location.hash);
	}

	get entries(): SidebarEntry[] {
		return this.#source().entries;
	}

	/** The entry chosen, resolved against the entries. */
	choice = $derived.by(() => resolveChoice(this.chosen, this.#source()));
	/** The pane shown. */
	pane = $derived.by(() => paneOf(this.choice, this.#source()));
	/** The scroll entries whose sections are on the page: those outside a pane or in the one shown. */
	watched = $derived.by(() =>
		flatten(this.#source().entries)
			.filter((p) => p.entry.mode === 'scroll' && (!p.pane || p.pane === this.pane))
			.map((p) => p.entry.id)
	);
	/** The entry current: the first section in view, else the choice. */
	current = $derived(this.inView[0] ?? this.choice);

	/** The entry of an id, if there is one. */
	entry(id: string | undefined): SidebarEntry | undefined {
		return flatten(this.entries).find((p) => p.entry.id === id)?.entry;
	}

	/** Choose an entry, or follow a hash; this writes nothing to the URL. */
	choose(id: string): void {
		this.chosen = id;
	}

	/** Follow the URL's hash, as arrival, back, and forward set it. */
	follow(hash: string): void {
		this.chosen = hashId(hash);
	}

	/** Follow the hash as back, forward, and links change it, until the returned stop is called. */
	listen(): () => void {
		if (typeof window === 'undefined') return () => {};
		const on = () => this.follow(location.hash);
		window.addEventListener('hashchange', on);
		window.addEventListener('popstate', on);
		return () => {
			window.removeEventListener('hashchange', on);
			window.removeEventListener('popstate', on);
		};
	}

	/**
	 * Mark the watched sections that are in view, as the reader scrolls, until the returned stop is
	 * called. A section counts as in view while it crosses the top two fifths of the viewport.
	 * Without an IntersectionObserver nothing is marked and the choice stays current.
	 */
	watch(find: (id: string) => Element | null = (id) => document.getElementById(id)): () => void {
		if (typeof IntersectionObserver !== 'function') return () => {};
		const ids = this.watched;
		// Plain records, not reactive ones: only inView, written once per report, is state.
		const of: { el: Element; id: string }[] = [];
		const seen: Record<string, boolean> = {};
		const io = new IntersectionObserver(
			(records) => {
				for (const r of records) {
					const id = of.find((o) => o.el === r.target)?.id;
					if (id) seen[id] = r.isIntersecting;
				}
				this.inView = ids.filter((id) => seen[id]);
			},
			{ rootMargin: '0px 0px -60% 0px' }
		);
		for (const id of ids) {
			const el = find(id);
			if (!el) continue;
			of.push({ el, id });
			io.observe(el);
		}
		return () => {
			io.disconnect();
			this.inView = [];
		};
	}
}
