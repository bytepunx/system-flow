// S-0458: the layout's top bar and nav rail. The top bar holds the flaiover name, the project
// picker, the theme setting, and the status pills, in that order, and no site-menu rows; the rail is
// docked beside the page from md, and below it a drawer the menu button opens, which Escape or
// choosing a page closes, returning focus to the button. The login page has neither. Since S-0534
// every other page names a provider out of budget under the header.
import { afterAll, afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import { createRawSnippet, flushSync, mount, unmount } from 'svelte';
import { inboxState } from '$lib/inbox.svelte';
import { hostFlai } from '$lib/hostflai.svelte';
import { projectState, resetForTests } from '$lib/project.svelte';

const where = vi.hoisted(() => ({ url: new URL('http://localhost/board') }));
vi.mock('$app/state', () => ({ page: where }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import Layout from './+layout.svelte';

// jsdom's <dialog> has no showModal or close. These stand in for the browser's: showModal opens
// the dialog and focuses its first control, Escape inside it closes it, and close fires `close`.
const proto = HTMLDialogElement.prototype as unknown as Record<string, unknown>;
const stubbed = typeof proto.showModal !== 'function';
beforeAll(() => {
	if (!stubbed) return;
	const escape = (e: KeyboardEvent) => {
		if (e.key !== 'Escape') return;
		const d = e.currentTarget as HTMLDialogElement;
		if (d.dispatchEvent(new Event('cancel', { cancelable: true }))) d.close();
	};
	proto.showModal = function (this: HTMLDialogElement) {
		this.setAttribute('open', '');
		this.addEventListener('keydown', escape);
		this.querySelector<HTMLElement>('button, a[href]')?.focus();
	};
	proto.close = function (this: HTMLDialogElement) {
		if (!this.hasAttribute('open')) return;
		this.removeAttribute('open');
		this.removeEventListener('keydown', escape);
		this.dispatchEvent(new Event('close'));
	};
});
afterAll(() => {
	if (!stubbed) return;
	delete proto.showModal;
	delete proto.close;
});

const page = createRawSnippet(() => ({ render: () => '<p data-testid="page">the page</p>' }));
const q = <T extends Element = HTMLElement>(sel: string) => document.querySelector<T>(sel);

describe('the layout’s top bar and nav rail (S-0458)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		localStorage.clear();
		// The layout starts these once signed in; here they stay quiet.
		vi.spyOn(inboxState, 'start').mockImplementation(() => {});
		vi.spyOn(inboxState, 'stop').mockImplementation(() => {});
		vi.spyOn(inboxState, 'restart').mockImplementation(() => {});
		vi.spyOn(hostFlai, 'start').mockImplementation(() => {});
		vi.spyOn(hostFlai, 'stop').mockImplementation(() => {});
		vi.spyOn(hostFlai, 'refresh').mockResolvedValue(undefined);
		vi.spyOn(projectState, 'refresh').mockResolvedValue(undefined);
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({ ok: true, json: async () => ({}) }))
		);
		hostFlai.status = {
			configured: true,
			connected: true,
			flai: '1.2.3',
			since: '2026-10-21T00:00:00Z'
		};
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		vi.restoreAllMocks();
		vi.unstubAllGlobals();
		hostFlai.status = null;
		resetForTests();
		localStorage.clear();
		document.body.innerHTML = '';
	});
	const show = (path: string) => {
		where.url = new URL(`http://localhost${path}`);
		c = mount(Layout, { target: document.body, props: { children: page } });
		flushSync();
	};

	it('holds the menu button, the name, the project picker, the theme, and the pills, in order', () => {
		show('/board');
		const bar = q('[data-top-bar]')!;
		const kinds = [...bar.children].map((el) => {
			if (el.matches('[data-nav-toggle]')) return 'menu';
			if (el.matches('[data-top-name]')) return 'name';
			if (el.querySelector('[data-testid="project-switcher"]')) return 'project';
			if (el.matches('[data-top-theme]')) return 'theme';
			if (el.matches('[data-top-pills]')) return 'pills';
			return el.outerHTML;
		});
		expect(kinds).toEqual(['menu', 'name', 'project', 'theme', 'pills']);
		const name = q<HTMLAnchorElement>('[data-top-name]')!;
		expect(name.textContent).toBe('flaiover');
		expect(name.getAttribute('href')).toBe('/');
		expect(q('[data-top-theme]')!.textContent!.trim()).toBe('theme: system');
		expect(q('[data-top-pills] a')!.getAttribute('href')).toBe('/host');
		expect(q('[data-testid="page"]')).not.toBeNull();
	});

	it('links the name to Overview keeping the project', () => {
		projectState.current = 'harbour';
		show('/board');
		expect(q('[data-top-name]')!.getAttribute('href')).toBe('/?project=harbour');
	});

	it('has no site-menu rows in the header, and the rail docked beside the page', () => {
		show('/board');
		expect(q('header nav')).toBeNull();
		expect(q('[data-menu-group], [data-menu-pages]')).toBeNull();
		const docked = q('[data-nav-docked]')!;
		expect(docked.tagName).toBe('ASIDE');
		expect(docked.className).toContain('md:block');
		expect(docked.className).toContain('hidden');
		expect(docked.querySelector('nav[aria-label="Site"]')).not.toBeNull();
		expect(q('[data-nav-drawer] nav')).toBeNull();
		expect(document.querySelectorAll('nav[aria-label="Site"]')).toHaveLength(1);
	});

	it('opens the drawer from the menu button, with the one rail in it', () => {
		show('/board');
		const toggle = q<HTMLButtonElement>('[data-nav-toggle]')!;
		expect(toggle.getAttribute('aria-label')).toBe('Open navigation');
		expect(toggle.className).toContain('md:hidden');
		expect(toggle.querySelector('svg')!.getAttribute('class')).toMatch(/(^|\s)size-\d+(\s|$)/);
		const drawer = q<HTMLDialogElement>('[data-nav-drawer]')!;
		expect(toggle.getAttribute('aria-controls')).toBe(drawer.id);
		expect(toggle.getAttribute('aria-expanded')).toBe('false');
		expect(drawer.hasAttribute('open')).toBe(false);
		toggle.click();
		flushSync();
		expect(drawer.hasAttribute('open')).toBe(true);
		expect(toggle.getAttribute('aria-expanded')).toBe('true');
		expect(drawer.querySelector('nav[aria-label="Site"]')).not.toBeNull();
		// one rail at a time, so its ids stay unique
		expect(document.querySelectorAll('nav[aria-label="Site"]')).toHaveLength(1);
		expect(drawer.contains(document.activeElement)).toBe(true);
	});

	it('closes the drawer on Escape and returns focus to the menu button', () => {
		show('/board');
		const toggle = q<HTMLButtonElement>('[data-nav-toggle]')!;
		toggle.click();
		flushSync();
		document.activeElement!.dispatchEvent(
			new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })
		);
		flushSync();
		const drawer = q<HTMLDialogElement>('[data-nav-drawer]')!;
		expect(drawer.hasAttribute('open')).toBe(false);
		expect(toggle.getAttribute('aria-expanded')).toBe('false');
		expect(document.activeElement).toBe(toggle);
		expect(drawer.querySelector('nav')).toBeNull();
		expect(q('[data-nav-docked] nav')).not.toBeNull();
	});

	it('closes the drawer when a page is chosen and returns focus to the menu button', () => {
		show('/board');
		const toggle = q<HTMLButtonElement>('[data-nav-toggle]')!;
		toggle.click();
		flushSync();
		const drawer = q<HTMLDialogElement>('[data-nav-drawer]')!;
		const board = drawer.querySelector<HTMLAnchorElement>('[data-rail-page="board"]')!;
		board.addEventListener('click', (e) => e.preventDefault());
		board.click();
		flushSync();
		expect(drawer.hasAttribute('open')).toBe(false);
		expect(document.activeElement).toBe(toggle);
	});

	it('shows no rail, menu button, project picker, or pills on the login page', () => {
		show('/login');
		expect(q('[data-nav-docked]')).toBeNull();
		expect(q('[data-nav-drawer]')).toBeNull();
		expect(q('[data-nav-toggle]')).toBeNull();
		expect(q('nav')).toBeNull();
		expect(q('[data-testid="project-switcher"]')).toBeNull();
		expect(q('[data-top-pills]')).toBeNull();
		expect(q('[data-top-name]')).not.toBeNull();
		expect(q('[data-top-theme]')).not.toBeNull();
		expect(q('[data-testid="page"]')).not.toBeNull();
	});

	// S-0534: a provider out of budget is named under the header on every page while signed in, and
	// the banner goes once flai serve holds none.
	const held = {
		provider: 'claude-code',
		limit: 'usage_limit_reached',
		since: '2026-10-28T05:30:00Z',
		until: '2026-10-28T09:00:00Z'
	};

	it.each(['/board', '/', '/threads'])('names a provider out of budget on %s', (path) => {
		hostFlai.status = { ...hostFlai.status!, budget: [held] };
		show(path);
		const banner = q('[data-budget-banner]')!;
		expect(banner.closest('header')).toBeNull();
		expect(banner.textContent).toContain('claude-code is out of budget (usage_limit_reached)');
		hostFlai.status = { ...hostFlai.status, budget: undefined };
		flushSync();
		expect(q('[data-budget-banner]')).toBeNull();
	});

	it('shows no budget banner on the login page', () => {
		hostFlai.status = { ...hostFlai.status!, budget: [held] };
		show('/login');
		expect(q('[data-budget-banner]')).toBeNull();
	});
});
