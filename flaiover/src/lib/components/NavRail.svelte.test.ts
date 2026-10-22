// S-0458: the nav rail. Each group is a button that shows or hides its pages; an expanded group's
// icon is in the accent colour and a collapsed one's muted; the page shown is current and its group
// expanded on arrival; a collapsed group carries its pages' indicators; the operator's choices are
// kept per browser. S-0484: Charts is a button that lists the charts by group, each a link, the
// chart shown current; it expands on arrival at a chart and keeps the operator's choice.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { inboxState, type Inbox } from '$lib/inbox.svelte';
import { projectState, resetForTests } from '$lib/project.svelte';
import { NavRailState } from '$lib/navrail.svelte';
import { CHART_GROUPS, TITLES } from '$lib/viz/catalog';

const where = vi.hoisted(() => ({ url: new URL('http://localhost/board') }));
vi.mock('$app/state', () => ({ page: where }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import NavRail from './NavRail.svelte';

const KEY = 'flaiover-nav-rail';
const BOX: Inbox = {
	total: 3,
	counts: { thread: 1, question: 0, review: 1, blocked: 1, overlap: 0 },
	notes: [],
	entries: []
};

const group = (key: string) =>
	document.querySelector(`[data-rail-group="${key}"]`) as HTMLButtonElement;
const pages = (key: string) => document.querySelector(`[data-rail-pages="${key}"]`) as HTMLElement;
const groupIcon = (key: string) => group(key).querySelector('[data-rail-icon]') as SVGElement;
const charts = () => document.querySelector('[data-rail-page="charts"]') as HTMLButtonElement;
const chart = (kind: string) =>
	document.querySelector(`[data-rail-chart="${kind}"]`) as HTMLAnchorElement;
const press = (el: HTMLElement) => {
	el.click();
	flushSync();
};

describe('the nav rail (S-0458)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let rail: NavRailState;
	beforeEach(() => {
		localStorage.clear();
		rail = new NavRailState();
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		inboxState.data = null;
		resetForTests();
		localStorage.clear();
		document.body.innerHTML = '';
	});
	const show = (path: string, props: { onnavigate?: () => void } = {}) => {
		where.url = new URL(`http://localhost${path}`);
		c = mount(NavRail, { target: document.body, props: { rail, ...props } });
		flushSync();
	};

	it('is the site navigation, with the five groups in order, each a button', () => {
		show('/');
		expect(document.querySelector('nav')!.getAttribute('aria-label')).toBe('Site');
		const groups = [...document.querySelectorAll<HTMLButtonElement>('[data-rail-group]')];
		expect(groups.map((b) => b.textContent!.trim())).toEqual([
			'Flow',
			'Agents',
			'Design',
			'Documentation',
			'Configuration'
		]);
		for (const b of groups) {
			expect(b.tagName).toBe('BUTTON');
			expect(b.getAttribute('aria-expanded')).toBe('false');
			expect(document.getElementById(b.getAttribute('aria-controls')!)).not.toBeNull();
		}
	});

	it('shows or hides a group’s pages when its button is pressed', () => {
		show('/');
		expect(pages('agents').hidden).toBe(true);
		press(group('agents'));
		expect(group('agents').getAttribute('aria-expanded')).toBe('true');
		expect(pages('agents').hidden).toBe(false);
		const links = [...pages('agents').querySelectorAll('a')];
		expect(links.map((a) => a.textContent!.trim())).toEqual([
			'Orchestration',
			'Planning',
			'Analysis',
			'Activity',
			'Threads',
			'Messages'
		]);
		for (const a of links) expect(a.querySelector('svg')).not.toBeNull();
		press(group('agents'));
		expect(group('agents').getAttribute('aria-expanded')).toBe('false');
		expect(pages('agents').hidden).toBe(true);
	});

	it('draws an expanded group’s icon in the accent colour and a collapsed one’s muted', () => {
		show('/');
		expect(groupIcon('design').getAttribute('class')).toContain('text-muted');
		expect(groupIcon('design').getAttribute('class')).not.toContain('text-accent');
		press(group('design'));
		expect(groupIcon('design').getAttribute('class')).toContain('text-accent');
		expect(groupIcon('design').getAttribute('class')).not.toContain('text-muted');
	});

	it('sizes every icon with a size utility only', () => {
		show('/board');
		for (const svg of document.querySelectorAll('nav svg'))
			expect(svg.getAttribute('class')).toMatch(/(^|\s)size-\d+(\s|$)/);
	});

	it('marks the page shown current and expands its group on arrival', () => {
		show('/settings/harnesses');
		expect(group('configuration').getAttribute('aria-expanded')).toBe('true');
		expect(group('flow').getAttribute('aria-expanded')).toBe('false');
		const current = document.querySelectorAll('[aria-current="page"]');
		expect(current).toHaveLength(1);
		expect(current[0].textContent!.trim()).toBe('Harnesses');
		expect(current[0].getAttribute('href')).toBe('/settings/harnesses');
	});

	it('expands the group of a page that belongs to it outside its pages', () => {
		show('/items/S-0001');
		expect(group('flow').getAttribute('aria-expanded')).toBe('true');
		expect(document.querySelector('[aria-current="page"]')).toBeNull();
	});

	it('expands the page’s group on arrival even when the operator collapsed it', () => {
		rail.toggle('flow');
		rail.toggle('flow');
		show('/inbox');
		expect(group('flow').getAttribute('aria-expanded')).toBe('true');
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ flow: false });
	});

	it('shows the Inbox count on Flow while it is collapsed, and on the Inbox link when open', () => {
		inboxState.data = BOX;
		show('/');
		const count = group('flow').querySelector('[data-testid="inbox-count"]');
		expect(count?.textContent).toBe('3');
		press(group('flow'));
		expect(group('flow').querySelector('[data-testid="inbox-count"]')).toBeNull();
		const inbox = pages('flow').querySelector('[data-rail-page="inbox"]')!;
		expect(inbox.querySelector('[data-testid="inbox-count"]')?.textContent).toBe('3');
	});

	it('shows Workflow as coming soon, not a link', () => {
		show('/board');
		const soon = pages('flow').querySelector('[data-rail-page="workflow"]')!;
		expect(soon.tagName).not.toBe('A');
		expect(soon.closest('a')).toBeNull();
		expect(soon.querySelector('a')).toBeNull();
		expect(soon.textContent).toContain('Workflow');
		expect(soon.textContent).toContain('soon');
		expect(soon.querySelector('.sr-only')?.textContent).toBe('(coming soon)');
	});

	it('links each page with the project it was followed from', () => {
		projectState.current = 'harbour';
		show('/settings');
		const hrefs = [...pages('configuration').querySelectorAll('a')].map((a) =>
			a.getAttribute('href')
		);
		expect(hrefs).toEqual([
			'/settings/features?project=harbour',
			'/settings/agents?project=harbour',
			'/settings/file-system?project=harbour',
			'/settings/projects?project=harbour',
			'/settings/harnesses?project=harbour'
		]);
		expect(document.querySelector('[data-rail-license]')!.getAttribute('href')).toBe(
			'/license?project=harbour'
		);
	});

	it('has a License link at its foot, current on the license page', () => {
		show('/license');
		const license = document.querySelector('[data-rail-license]')!;
		expect(license.textContent).toBe('License');
		expect(license.getAttribute('href')).toBe('/license');
		expect(license.getAttribute('aria-current')).toBe('page');
	});

	it('calls onnavigate when a page or the License link is chosen', () => {
		const onnavigate = vi.fn();
		show('/board', { onnavigate });
		const board = pages('flow').querySelector('[data-rail-page="board"]') as HTMLAnchorElement;
		board.addEventListener('click', (e) => e.preventDefault());
		board.click();
		expect(onnavigate).toHaveBeenCalledTimes(1);
		const license = document.querySelector('[data-rail-license]') as HTMLAnchorElement;
		license.addEventListener('click', (e) => e.preventDefault());
		license.click();
		expect(onnavigate).toHaveBeenCalledTimes(2);
		press(group('agents'));
		expect(onnavigate).toHaveBeenCalledTimes(2);
	});

	it('keeps the groups the operator expanded or collapsed in localStorage', () => {
		show('/board');
		press(group('design'));
		press(group('flow'));
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ design: true, flow: false });
		unmount(c!);
		c = undefined;
		document.body.innerHTML = '';
		rail = new NavRailState();
		show('/');
		expect(group('design').getAttribute('aria-expanded')).toBe('true');
		expect(group('flow').getAttribute('aria-expanded')).toBe('false');
	});

	describe('Charts (S-0484)', () => {
		it('is a button that controls its list and links nowhere', () => {
			show('/board');
			expect(charts().tagName).toBe('BUTTON');
			expect(charts().getAttribute('type')).toBe('button');
			expect(charts().hasAttribute('href')).toBe(false);
			expect(charts().closest('a')).toBeNull();
			expect(charts().getAttribute('aria-expanded')).toBe('false');
			expect(charts().getAttribute('aria-controls')).toBe('nav-rail-page-charts');
			expect(charts().textContent!.trim()).toBe('Charts');
		});

		it('draws nothing below it while collapsed', () => {
			show('/board');
			expect(document.getElementById('nav-rail-page-charts')).toBeNull();
			expect(document.querySelector('[data-rail-chart]')).toBeNull();
			expect(document.querySelector('[data-rail-label]')).toBeNull();
		});

		it('lists the four groups’ muted labels and their charts in catalogue order when pressed', () => {
			projectState.current = 'harbour';
			show('/board');
			press(charts());
			expect(charts().getAttribute('aria-expanded')).toBe('true');
			const list = document.getElementById('nav-rail-page-charts')!;
			const labels = [...list.querySelectorAll<HTMLElement>('[data-rail-label]')];
			expect(labels.map((l) => l.textContent!.trim())).toEqual([
				'flow',
				'usage',
				'planning',
				'strategic'
			]);
			for (const l of labels) {
				expect(l.tagName).not.toBe('A');
				expect(l.closest('a')).toBeNull();
				expect(l.getAttribute('class')).toContain('text-muted');
			}
			const links = [...list.querySelectorAll<HTMLAnchorElement>('[data-rail-chart]')];
			const kinds = CHART_GROUPS.flatMap((g) => g.kinds);
			expect(links.map((a) => a.dataset.railChart)).toEqual(kinds);
			for (const a of links) {
				expect(a.tagName).toBe('A');
				expect(a.getAttribute('href')).toBe(`/charts/${a.dataset.railChart}?project=harbour`);
			}
			for (const g of CHART_GROUPS) {
				const under = list.querySelector(`[data-rail-label="${g.key}"]`)!.parentElement!;
				const titles = [...under.querySelectorAll('[data-rail-chart]')].map((a) => a.textContent);
				expect(titles).toEqual(g.kinds.map((k) => TITLES[k]));
			}
			expect(chart('cost-spent').textContent).toBe('$ / Day');
		});

		it('draws its icon in the accent colour while expanded and muted while collapsed', () => {
			show('/board');
			const icon = () => charts().querySelector('[data-rail-icon]') as SVGElement;
			expect(icon().getAttribute('class')).toContain('text-muted');
			expect(icon().getAttribute('class')).not.toContain('text-accent');
			press(charts());
			expect(icon().getAttribute('class')).toContain('text-accent');
			expect(icon().getAttribute('class')).not.toContain('text-muted');
		});

		it('expands Flow and Charts on arrival at a chart, marks it current, and stores nothing', () => {
			show('/charts/throughput');
			expect(group('flow').getAttribute('aria-expanded')).toBe('true');
			expect(charts().getAttribute('aria-expanded')).toBe('true');
			expect(localStorage.getItem(KEY)).toBeNull();
			const current = document.querySelectorAll('[aria-current="page"]');
			expect(current).toHaveLength(1);
			expect(current[0]).toBe(chart('throughput'));
			expect(chart('throughput').getAttribute('class')).toContain('font-semibold');
			expect(chart('throughput').getAttribute('class')).toContain('text-accent');
			expect(chart('cycle-time').getAttribute('class')).not.toContain('font-semibold');
			expect(charts().getAttribute('class')).toContain('font-semibold');
			expect(charts().getAttribute('class')).toContain('text-accent');
		});

		it('marks Charts current with no chart at an unknown kind', () => {
			show('/charts/nope');
			expect(charts().getAttribute('aria-expanded')).toBe('true');
			expect(charts().getAttribute('class')).toContain('font-semibold');
			expect(document.querySelector('[data-rail-chart][aria-current]')).toBeNull();
			expect(document.querySelector('[aria-current="page"]')).toBeNull();
		});

		it('is not current away from the charts', () => {
			show('/board');
			expect(charts().getAttribute('class')).not.toContain('font-semibold');
		});

		it('keeps the operator’s choice in localStorage beside the groups’', () => {
			show('/board');
			press(charts());
			expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ charts: true });
			press(charts());
			expect(charts().getAttribute('aria-expanded')).toBe('false');
			expect(document.querySelector('[data-rail-chart]')).toBeNull();
			expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ charts: false });
			press(charts());
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			rail = new NavRailState();
			show('/board');
			expect(charts().getAttribute('aria-expanded')).toBe('true');
		});

		it('calls onnavigate when a chart is chosen, and not when Charts is pressed', () => {
			const onnavigate = vi.fn();
			show('/charts/cycle-time', { onnavigate });
			press(charts());
			press(charts());
			expect(onnavigate).not.toHaveBeenCalled();
			const a = chart('throughput');
			a.addEventListener('click', (e) => e.preventDefault());
			a.click();
			expect(onnavigate).toHaveBeenCalledTimes(1);
		});
	});
});
