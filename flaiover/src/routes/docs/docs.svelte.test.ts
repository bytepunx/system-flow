// S-0223: the documents page shows the analyzer's reports under design/analysis as it shows
// design/experiments, and opens a report with its front matter, the focus and window among it.
// S-0461: the tree heads the page's sidebar, and the open document's headings are its entries.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const page = vi.hoisted(() => ({
	params: { path: '' },
	url: new URL('http://localhost/docs'),
	moved: () => {}
}));
// SvelteKit keeps the page when it opens another document: its params are read as they change.
vi.mock('$app/state', async () => {
	const { createSubscriber } = await import('svelte/reactivity');
	const follows = createSubscriber((update) => {
		page.moved = update;
		return () => {};
	});
	return {
		page: {
			get params() {
				follows();
				return page.params;
			},
			get url() {
				return page.url;
			}
		}
	};
});
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
const enhance = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$lib/markdown', async (actual) => ({
	...(await actual<typeof import('$lib/markdown')>()),
	enhance
}));
vi.mock('$lib/events', () => ({ listen: () => () => {}, follow: () => () => {} }));

import DocsPage from './[...path]/+page.svelte';

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const file = (path: string, title: string) => ({
	name: path.split('/').pop(),
	path,
	kind: 'file',
	title
});
const dir = (path: string, children: unknown[]) => ({
	name: path.split('/').pop(),
	path,
	kind: 'dir',
	children
});
const reportPath = 'design/analysis/2026-10-06-bottlenecks.md';
const tree = [
	dir('design', [
		dir('design/analysis', [
			file(reportPath, 'Bottlenecks in September'),
			file('design/analysis/README.md', 'analysis')
		]),
		dir('design/experiments', [file('design/experiments/S-0001-first.md', 'First experiment')]),
		dir('design/system', [file('design/system/overview.md', 'Overview')])
	]),
	dir('docs', []),
	dir('wip', [])
];
const report = {
	path: reportPath,
	frontMatter: {
		title: 'Bottlenecks in September',
		updated: '2026-10-06T08:00:00Z',
		status: 'active',
		focus: 'bottlenecks',
		from: '2026-09-01',
		to: '2026-09-30'
	},
	body: '\n# Bottlenecks in September\n\nReview waits longest.\n'
};
const overviewPath = 'design/system/overview.md';
const overview = {
	path: overviewPath,
	frontMatter: null,
	body: [
		'# Overview',
		'',
		'The four parts.',
		'',
		'## The `flai` CLI',
		'',
		'```sh',
		'# not a heading',
		'```',
		'',
		'### Commands',
		'',
		'## The dashboard',
		'',
		'### Commands',
		''
	].join('\n')
};
const scrollIntoView = Element.prototype.scrollIntoView;

describe('the documents page (S-0223)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		page.params = { path: '' };
		document.body.innerHTML = '';
		Element.prototype.scrollIntoView = scrollIntoView;
		history.replaceState(null, '', '/');
	});

	const serve = () =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/docs/tree') return answer(tree);
			if (url === `/api/docs/file?path=${encodeURIComponent(reportPath)}`) return answer(report);
			if (url === `/api/docs/file?path=${encodeURIComponent(overviewPath)}`)
				return answer(overview);
			if (url === '/api/board') return answer({ writable: false });
			return answer([]);
		});
	const folder = (name: string) =>
		[...document.querySelectorAll('[data-doc-tree] button')].find((b) =>
			b.textContent!.trim().endsWith(name)
		);
	const links = () =>
		[...document.querySelectorAll('[data-doc-tree] a')].map((a) => [
			a.textContent!.trim(),
			a.getAttribute('href')
		]);

	it('lists analysis under design as it lists experiments, with its reports inside', async () => {
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
		const folders = [...document.querySelectorAll('[data-doc-tree] button')].map((b) =>
			b.textContent!.replace(/[▸▾]/g, '').trim()
		);
		expect(folders).toEqual(['design', 'analysis', 'experiments', 'system', 'docs', 'wip']);
		expect(links()).toEqual([]);

		folder('analysis')!.dispatchEvent(new MouseEvent('click', { bubbles: true }));
		folder('experiments')!.dispatchEvent(new MouseEvent('click', { bubbles: true }));
		flushSync();
		expect(links()).toEqual([
			['Bottlenecks in September', `/docs/${reportPath}`],
			['analysis', '/docs/design/analysis/README.md'],
			['First experiment', '/docs/design/experiments/S-0001-first.md']
		]);
	});

	it('opens a report with its front matter, the focus and the window among it', async () => {
		page.params = { path: reportPath };
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
		// the report's folder starts open, with the report marked as the one shown
		const open = document.querySelector<HTMLAnchorElement>(
			`[data-doc-tree] a[href="/docs/${reportPath}"]`
		)!;
		expect(open.className).toContain('bg-raised');
		expect(document.querySelector('summary')!.textContent).toContain(reportPath);
		const fields = Object.fromEntries(
			[...document.querySelectorAll('dl dt')].map((dt) => [
				dt.textContent!.trim(),
				dt.nextElementSibling!.textContent!.trim()
			])
		);
		expect(fields).toEqual({
			title: 'Bottlenecks in September',
			updated: '2026-10-06T08:00:00Z',
			status: 'active',
			focus: 'bottlenecks',
			from: '2026-09-01',
			to: '2026-09-30'
		});
		expect(document.querySelector('article')!.textContent).toContain('Review waits longest.');
	});
});

describe('the documents page sidebar (S-0461)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		page.params = { path: '' };
		document.body.innerHTML = '';
		Element.prototype.scrollIntoView = scrollIntoView;
		history.replaceState(null, '', '/');
	});

	const serve = () =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/docs/tree') return answer(tree);
			if (url === `/api/docs/file?path=${encodeURIComponent(reportPath)}`) return answer(report);
			if (url === `/api/docs/file?path=${encodeURIComponent(overviewPath)}`)
				return answer(overview);
			if (url === '/api/board') return answer({ writable: true });
			return answer([]);
		});
	const open = async (path: string) => {
		page.params = { path };
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
	};
	// The heading entries, by id and label, each with the ids of the entries listed under it.
	const entries = (list: Element | null): unknown[] =>
		[...(list?.querySelectorAll(':scope > li') ?? [])].map((li) => {
			const a = li.querySelector(':scope > a')!;
			const under = entries(li.querySelector(':scope > ul'));
			return under.length
				? [a.getAttribute('data-sidebar-entry'), a.textContent!.trim(), under]
				: [a.getAttribute('data-sidebar-entry'), a.textContent!.trim()];
		});
	const headings = () =>
		entries(document.querySelector('[data-sidebar-entry="On this page"] + ul'));
	// Record what is scrolled into view.
	const scrolls = () => {
		const into: Element[] = [];
		Element.prototype.scrollIntoView = function (this: Element) {
			into.push(this);
		};
		return into;
	};

	it('puts the tree in the sidebar, above the entries, with the document beside them', async () => {
		await open(overviewPath);
		const nav = document.querySelector('nav[aria-label="Documents"]')!;
		const header = nav.querySelector('[data-sidebar-header]')!;
		expect(header.querySelector('[data-doc-tree]')).not.toBeNull();
		expect(
			header.compareDocumentPosition(nav.querySelector('[data-sidebar-list="top"]')!) &
				Node.DOCUMENT_POSITION_FOLLOWING
		).toBeTruthy();
		// the open document is marked in the tree, and the document is the sidebar's content
		const shown = header.querySelector<HTMLAnchorElement>(`a[href="/docs/${overviewPath}"]`)!;
		expect(shown.className).toContain('bg-raised');
		const content = document.querySelector('[data-sidebar-content]')!;
		expect(content.querySelector('article')!.textContent).toContain('The four parts.');
		expect(content.querySelector('a[href="/edit/design/system/overview.md"]')).not.toBeNull();
	});

	it("lists the open document's headings in order, each under the one it sits in", async () => {
		await open(overviewPath);
		// a fenced comment is no heading, and a repeated heading gets the renderer's unique id
		expect(headings()).toEqual([
			[
				'overview',
				'Overview',
				[
					['the-flai-cli', 'The flai CLI', [['commands', 'Commands']]],
					['the-dashboard', 'The dashboard', [['commands-1', 'Commands']]]
				]
			]
		]);
		// each entry names a heading drawn in the document
		const article = document.querySelector('article')!;
		for (const id of ['overview', 'the-flai-cli', 'commands', 'the-dashboard', 'commands-1'])
			expect(article.querySelector(`[id="${id}"]`)?.tagName).toMatch(/^H[1-6]$/);
		for (const a of document.querySelectorAll('[data-sidebar-list="nested"] a'))
			expect(a.getAttribute('data-mode')).toBe('scroll');
	});

	it('scrolls to the heading chosen, writes its hash, and marks it current', async () => {
		await open(overviewPath);
		const into = scrolls();
		const entry = document.querySelector<HTMLAnchorElement>(
			'[data-sidebar-entry="the-dashboard"]'
		)!;
		entry.click();
		await settle();
		expect(into).toEqual([document.getElementById('the-dashboard')]);
		expect(into[0].tagName).toBe('H2');
		expect(location.hash).toBe('#the-dashboard');
		expect(entry.getAttribute('aria-current')).toBe('true');
		// the second Commands, not the first
		document.querySelector<HTMLAnchorElement>('[data-sidebar-entry="commands-1"]')!.click();
		await settle();
		expect(into[1]).toBe(document.querySelector('article h3:last-of-type'));
		expect(into[1].id).toBe('commands-1');
		expect(location.hash).toBe('#commands-1');
	});

	it('arriving with a heading in the hash lands on it once the document is drawn', async () => {
		history.replaceState(null, '', `/docs/${overviewPath}#the-flai-cli`);
		const into = scrolls();
		await open(overviewPath);
		expect(into.map((e) => e.id)).toContain('the-flai-cli');
		expect(
			document.querySelector('[data-sidebar-entry="the-flai-cli"]')!.getAttribute('aria-current')
		).toBe('true');
	});

	it('follows the document opened: its headings replace the last one', async () => {
		await open(overviewPath);
		expect(headings()).toHaveLength(1);
		page.params = { path: reportPath };
		page.moved();
		await settle();
		expect(headings()).toEqual([['bottlenecks-in-september', 'Bottlenecks in September']]);
		expect(document.querySelector('[data-sidebar-entry="the-flai-cli"]')).toBeNull();
	});

	it('lists no headings and no label before a document is open', async () => {
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-sidebar-header] [data-doc-tree]')).not.toBeNull();
		expect(document.querySelector('[data-sidebar-list="top"]')!.children).toHaveLength(0);
		expect(document.querySelector('[data-sidebar-content]')!.textContent).toContain(
			'Pick a document from the tree.'
		);
	});
});
