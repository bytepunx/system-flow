// S-0460: the state behind a page's sidebar. The choice resolves against the entries, falling back
// to the page's default and then the first entry; the pane shown is the choice's own or the one it
// sits in; the hash chooses both ways; a scroll entry is current while its section is in view.
// S-0461: a linked entry is never a pane or watched, and a page of links does not follow the hash.
import { afterEach, describe, expect, it, vi } from 'vitest';
import {
	PageSidebarState,
	flatten,
	hashId,
	hashless,
	kindOf,
	paneOf,
	resolveChoice,
	scrollBehavior,
	type SidebarEntry,
	type SidebarSource
} from './pagesidebar.svelte';

const PANES: SidebarEntry[] = [
	{ id: 'story', label: 'Story', mode: 'pane' },
	{ id: 'tasks', label: 'Tasks', mode: 'pane', count: 3 },
	{
		id: 'questions',
		label: 'Questions and threads',
		mode: 'pane',
		count: 2,
		children: [{ id: 'open-threads', label: 'Open threads', mode: 'scroll' }]
	}
];
const SECTIONS: SidebarEntry[] = [
	{ id: 'goal', label: 'Goal', mode: 'scroll' },
	{ id: 'criteria', label: 'Criteria', mode: 'scroll' },
	{ id: 'notes', label: 'Notes', mode: 'scroll' }
];

// The charts page's entries: a label over links to other pages.
const LINKS: SidebarEntry[] = [
	{
		id: 'flow',
		label: 'Flow',
		children: [
			{ id: 'cfd', label: 'Cumulative flow', href: '/charts/cfd' },
			{ id: 'wip', label: 'WIP', href: '/charts/wip' }
		]
	},
	{ id: 'lead', label: 'Lead time', href: '/charts/lead' }
];

const src = (entries: SidebarEntry[], more: Partial<SidebarSource> = {}): SidebarSource => ({
	entries,
	...more
});

describe('flatten', () => {
	it('lists every entry depth first, each with the pane it sits in', () => {
		expect(flatten(PANES).map((p) => [p.entry.id, p.pane])).toEqual([
			['story', 'story'],
			['tasks', 'tasks'],
			['questions', 'questions'],
			['open-threads', 'questions']
		]);
		expect(flatten(SECTIONS).map((p) => p.pane)).toEqual([undefined, undefined, undefined]);
	});
});

describe('kindOf (S-0461)', () => {
	it('is a link for an href whatever the mode, the mode otherwise, and a label without either', () => {
		expect(kindOf({ id: 'a', label: 'A', href: '/charts/a' })).toBe('link');
		expect(kindOf({ id: 'a', label: 'A', href: '/charts/a', mode: 'pane' })).toBe('link');
		expect(kindOf({ id: 'a', label: 'A', mode: 'pane' })).toBe('pane');
		expect(kindOf({ id: 'a', label: 'A', mode: 'scroll' })).toBe('scroll');
		expect(kindOf({ id: 'a', label: 'A' })).toBe('label');
	});

	it('makes a linked entry no pane for the entries under it', () => {
		const under = [{ id: 'x', label: 'X', href: '/x', mode: 'pane' as const, children: SECTIONS }];
		expect(flatten(under).map((p) => p.pane)).toEqual([undefined, undefined, undefined, undefined]);
	});
});

describe('hashless (S-0461)', () => {
	it('is true for a page of links and labels only, and not while there are no entries', () => {
		expect(hashless(LINKS)).toBe(true);
		expect(hashless([])).toBe(false);
		expect(hashless(PANES)).toBe(false);
		expect(hashless([...LINKS, ...SECTIONS])).toBe(false);
	});
});

describe('hashId', () => {
	it('drops the hash sign and decodes, and is empty for no hash', () => {
		expect(hashId('#tasks')).toBe('tasks');
		expect(hashId('tasks')).toBe('tasks');
		expect(hashId('#a%20b')).toBe('a b');
		expect(hashId('')).toBe('');
		expect(hashId('#')).toBe('');
		expect(hashId('#%E0%A4%A')).toBe('%E0%A4%A');
	});
});

describe('resolveChoice', () => {
	it('keeps a choice that names an entry, nested ones included', () => {
		expect(resolveChoice('tasks', src(PANES))).toBe('tasks');
		expect(resolveChoice('open-threads', src(PANES))).toBe('open-threads');
	});

	it('falls back to the default, then to the first entry, for an unknown or empty choice', () => {
		expect(resolveChoice('nope', src(PANES, { fallback: 'tasks' }))).toBe('tasks');
		expect(resolveChoice('', src(PANES, { fallback: 'tasks' }))).toBe('tasks');
		expect(resolveChoice(undefined, src(PANES))).toBe('story');
		expect(resolveChoice('nope', src(PANES, { fallback: 'gone' }))).toBe('story');
	});

	it('takes the entry match names for a hash that is not an entry', () => {
		const match = (h: string) => (h.startsWith('MS-') ? 'tasks' : undefined);
		expect(resolveChoice('MS-0002', src(PANES, { match, fallback: 'story' }))).toBe('tasks');
		expect(resolveChoice('TH-0001', src(PANES, { match, fallback: 'story' }))).toBe('story');
		expect(resolveChoice('MS-0002', src(PANES, { match: () => 'gone' }))).toBe('story');
	});

	it('never chooses a label, and chooses a link the page names (S-0461)', () => {
		expect(resolveChoice('flow', src(LINKS))).toBe('cfd');
		expect(resolveChoice(undefined, src(LINKS))).toBe('cfd');
		expect(resolveChoice('wip', src(LINKS))).toBe('wip');
		expect(resolveChoice('nope', src(LINKS, { fallback: 'lead' }))).toBe('lead');
		expect(resolveChoice('nope', src(LINKS, { fallback: 'flow' }))).toBe('cfd');
	});

	it('resolves to nothing while there are no entries', () => {
		expect(resolveChoice('tasks', src([]))).toBeUndefined();
	});
});

describe('paneOf', () => {
	it('shows a pane’s own pane, and a nested scroll entry’s pane above it', () => {
		expect(paneOf('tasks', src(PANES))).toBe('tasks');
		expect(paneOf('open-threads', src(PANES))).toBe('questions');
	});

	it('shows the default’s pane, else the first, for a scroll entry outside every pane', () => {
		const mixed = [...SECTIONS, ...PANES];
		expect(paneOf('goal', src(mixed, { fallback: 'tasks' }))).toBe('tasks');
		expect(paneOf('goal', src(mixed))).toBe('story');
	});

	it('shows no pane on a page of scroll entries', () => {
		expect(paneOf('goal', src(SECTIONS))).toBeUndefined();
	});

	it('shows no pane on a page of links (S-0461)', () => {
		expect(paneOf('wip', src(LINKS))).toBeUndefined();
	});
});

describe('scrollBehavior', () => {
	afterEach(() => vi.unstubAllGlobals());

	it('scrolls at once without matchMedia, as under jsdom', () => {
		vi.stubGlobal('matchMedia', undefined);
		expect(scrollBehavior()).toBe('auto');
	});

	it('scrolls at once for reduced motion and smoothly otherwise', () => {
		vi.stubGlobal('matchMedia', (q: string) => ({ matches: q.includes('reduce') }));
		expect(scrollBehavior()).toBe('auto');
		vi.stubGlobal('matchMedia', () => ({ matches: false }));
		expect(scrollBehavior()).toBe('smooth');
	});
});

describe('PageSidebarState', () => {
	afterEach(() => {
		history.replaceState(null, '', location.pathname);
		vi.unstubAllGlobals();
	});

	it('starts on the default without a hash, and on the entry the hash names', () => {
		expect(new PageSidebarState(() => src(PANES, { fallback: 'tasks' })).current).toBe('tasks');
		history.replaceState(null, '', '#questions');
		const s = new PageSidebarState(() => src(PANES));
		expect(s.choice).toBe('questions');
		expect(s.pane).toBe('questions');
		expect(s.current).toBe('questions');
	});

	it('keeps a hash for a pane the page has not drawn yet, and falls back when it goes', () => {
		history.replaceState(null, '', '#tasks');
		let entries = $state<SidebarEntry[]>([]);
		const s = new PageSidebarState(() => src(entries, { fallback: 'story' }));
		expect(s.current).toBeUndefined();
		entries = PANES;
		expect(s.current).toBe('tasks');
		entries = PANES.filter((e) => e.id !== 'tasks');
		expect(s.current).toBe('story');
	});

	it('follows the hash as back and forward change it, until stopped', () => {
		const s = new PageSidebarState(() => src(PANES));
		const stop = s.listen();
		history.replaceState(null, '', '#tasks');
		window.dispatchEvent(new HashChangeEvent('hashchange'));
		expect(s.pane).toBe('tasks');
		history.replaceState(null, '', location.pathname);
		window.dispatchEvent(new PopStateEvent('popstate'));
		expect(s.pane).toBe('story');
		stop();
		history.replaceState(null, '', '#questions');
		window.dispatchEvent(new HashChangeEvent('hashchange'));
		expect(s.pane).toBe('story');
	});

	it('chooses a nested scroll entry, showing the pane it sits in', () => {
		const s = new PageSidebarState(() => src(PANES));
		s.choose('open-threads');
		expect(s.pane).toBe('questions');
		expect(s.current).toBe('open-threads');
		expect(s.watched).toEqual(['open-threads']);
		s.choose('tasks');
		expect(s.watched).toEqual([]);
	});

	it('watches no link, and leaves a page of links to its current on back and forward (S-0461)', () => {
		history.replaceState(null, '', '#wip');
		const s = new PageSidebarState(() => src([...LINKS, ...SECTIONS]));
		expect(s.current).toBe('wip');
		expect(s.watched).toEqual(['goal', 'criteria', 'notes']);
		const links = new PageSidebarState(() => src(LINKS));
		links.choose('lead');
		const stop = links.listen();
		history.replaceState(null, '', location.pathname);
		window.dispatchEvent(new PopStateEvent('popstate'));
		expect(links.current).toBe('lead');
		expect(links.watched).toEqual([]);
		stop();
	});

	it('makes the choice current where there is no IntersectionObserver', () => {
		vi.stubGlobal('IntersectionObserver', undefined);
		const s = new PageSidebarState(() => src(SECTIONS));
		const stop = s.watch(() => document.createElement('section'));
		expect(s.current).toBe('goal');
		s.choose('notes');
		expect(s.current).toBe('notes');
		stop();
	});

	it('makes the first section in view current while it is in view', () => {
		type Fake = { cb: IntersectionObserverCallback; seen: Element[]; gone: boolean };
		const made: Fake[] = [];
		vi.stubGlobal(
			'IntersectionObserver',
			class {
				f: Fake;
				constructor(cb: IntersectionObserverCallback) {
					this.f = { cb, seen: [], gone: false };
					made.push(this.f);
				}
				observe(el: Element) {
					this.f.seen.push(el);
				}
				disconnect() {
					this.f.gone = true;
				}
			}
		);
		const els = Object.fromEntries(
			SECTIONS.map((e) => [e.id, Object.assign(document.createElement('section'), { id: e.id })])
		);
		const s = new PageSidebarState(() => src(SECTIONS));
		const stop = s.watch((id) => els[id] ?? null);
		const io = made[0];
		expect(io.seen).toEqual([els.goal, els.criteria, els.notes]);
		const cross = (id: string, isIntersecting: boolean) =>
			io.cb(
				[{ target: els[id], isIntersecting } as unknown as IntersectionObserverEntry],
				{} as IntersectionObserver
			);
		cross('notes', true);
		expect(s.current).toBe('notes');
		cross('criteria', true);
		expect(s.current).toBe('criteria');
		cross('criteria', false);
		expect(s.current).toBe('notes');
		stop();
		expect(io.gone).toBe(true);
		expect(s.current).toBe('goal');
	});
});
