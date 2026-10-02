// S-0172: the header's two-tier site menu. A group opens on hover, click, or touch; the page shown and
// its group are active; a group carries its pages' indicators.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { inboxState, type Inbox } from '$lib/inbox.svelte';

const where = vi.hoisted(() => ({ url: new URL('http://localhost/board') }));
vi.mock('$app/state', () => ({ page: where }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import SiteMenu from './SiteMenu.svelte';

const BOX: Inbox = {
	total: 3,
	counts: { thread: 1, question: 0, review: 1, blocked: 1, overlap: 0 },
	notes: [],
	entries: []
};

const group = (key: string) =>
	document.querySelector(`[data-menu-group="${key}"]`) as HTMLButtonElement;
const shownPages = () =>
	[...document.querySelectorAll('[data-menu-pages] a')].map((a) => a.textContent);
const pointer = (el: Element, type: string, pointerType = 'mouse') => {
	const e = new Event(type, { bubbles: type !== 'pointerenter' && type !== 'pointerleave' });
	Object.assign(e, { pointerType });
	el.dispatchEvent(e);
	flushSync();
};

describe('the site menu (S-0172)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		inboxState.data = null;
		document.body.innerHTML = '';
	});
	const show = (path: string) => {
		where.url = new URL(`http://localhost${path}`);
		c = mount(SiteMenu, { target: document.body });
		flushSync();
	};

	it('shows every group, and the pages of the group of the page shown', () => {
		show('/board');
		expect([...document.querySelectorAll('[data-menu-group]')].map((b) => b.textContent)).toEqual([
			'Workflow',
			'Status',
			'Host'
		]);
		expect(shownPages()).toEqual(['Overview', 'Board', 'Inbox', 'Threads', 'Activity']);
	});

	it('marks the page shown and its group active', () => {
		show('/settings');
		expect(group('host').dataset.active).toBe('true');
		expect(group('host').getAttribute('aria-expanded')).toBe('true');
		expect(group('workflow').dataset.active).toBe('false');
		const current = document.querySelector('[aria-current="page"]')!;
		expect(current.textContent).toBe('Settings');
		expect(current.className).toContain('font-semibold');
		expect(group('host').className).toContain('font-semibold');
	});

	it('links each page to its route', () => {
		show('/charts/cycle-time');
		const links = [...document.querySelectorAll('[data-menu-pages] a')].map((a) =>
			a.getAttribute('href')
		);
		expect(links).toEqual(['/charts/cycle-time', '/adrs', '/docs/', '/search']);
		pointer(group('host'), 'pointerenter');
		expect(
			[...document.querySelectorAll('[data-menu-pages] a')].map((a) => a.getAttribute('href'))
		).toEqual(['/host', '/settings', '/license']);
	});

	it('opens a group on hover and goes back to the page shown when the pointer leaves', () => {
		show('/board');
		pointer(group('status'), 'pointerenter');
		expect(shownPages()).toEqual(['Charts', 'ADRs', 'Docs', 'Search']);
		expect(group('status').dataset.active).toBe('true');
		expect(group('workflow').dataset.active).toBe('false');
		pointer(document.querySelector('nav')!, 'pointerleave');
		expect(shownPages()).toEqual(['Overview', 'Board', 'Inbox', 'Threads', 'Activity']);
		expect(group('workflow').dataset.active).toBe('true');
	});

	it('opens a group on a click or a touch, and a touch ending does not close it', () => {
		show('/board');
		group('host').click();
		flushSync();
		expect(shownPages()).toEqual(['Updates', 'Settings', 'License']);
		pointer(group('status'), 'pointerenter', 'touch');
		expect(shownPages()).toEqual(['Updates', 'Settings', 'License']);
		group('status').click();
		flushSync();
		pointer(document.querySelector('nav')!, 'pointerleave', 'touch');
		expect(shownPages()).toEqual(['Charts', 'ADRs', 'Docs', 'Search']);
	});

	it('closes an opened group on Escape', () => {
		show('/board');
		group('host').click();
		flushSync();
		document
			.querySelector('nav')!
			.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true }));
		flushSync();
		expect(shownPages()).toEqual(['Overview', 'Board', 'Inbox', 'Threads', 'Activity']);
	});

	it('opens the group of a page only its group claims, with no page current', () => {
		show('/items/S-0172');
		expect(group('workflow').dataset.active).toBe('true');
		expect(document.querySelector('[aria-current="page"]')).toBeNull();
	});

	it('gives the Inbox count to the Workflow group as well, and to no other group', () => {
		show('/settings');
		inboxState.data = BOX;
		flushSync();
		expect(group('workflow').querySelector('[data-testid="inbox-count"]')?.textContent).toBe('3');
		expect(group('status').querySelector('[data-testid="inbox-count"]')).toBeNull();
		expect(group('host').querySelector('[data-testid="inbox-count"]')).toBeNull();
		pointer(group('workflow'), 'pointerenter');
		const inbox = document.querySelector('[data-menu-page="inbox"]')!;
		expect(inbox.querySelector('[data-testid="inbox-count"]')?.textContent).toBe('3');
	});
});
