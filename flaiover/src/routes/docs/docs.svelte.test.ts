// S-0223: the documents page shows the analyzer's reports under design/analysis as it shows
// design/experiments, and opens a report with its front matter, the focus and window among it.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const page = vi.hoisted(() => ({ params: { path: '' }, url: new URL('http://localhost/docs') }));
vi.mock('$app/state', () => ({ page }));
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

describe('the documents page (S-0223)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		page.params.path = '';
		document.body.innerHTML = '';
	});

	const serve = () =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/docs/tree') return answer(tree);
			if (url === `/api/docs/file?path=${encodeURIComponent(reportPath)}`) return answer(report);
			if (url === '/api/board') return answer({ writable: false });
			return answer([]);
		});
	const folder = (name: string) =>
		[...document.querySelectorAll('aside button')].find((b) =>
			b.textContent!.trim().endsWith(name)
		);
	const links = () =>
		[...document.querySelectorAll('aside a')].map((a) => [
			a.textContent!.trim(),
			a.getAttribute('href')
		]);

	it('lists analysis under design as it lists experiments, with its reports inside', async () => {
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
		const folders = [...document.querySelectorAll('aside button')].map((b) =>
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
		page.params.path = reportPath;
		serve();
		c = mount(DocsPage, { target: document.body });
		await settle();
		// the report's folder starts open, with the report marked as the one shown
		const open = document.querySelector<HTMLAnchorElement>(`aside a[href="/docs/${reportPath}"]`)!;
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
