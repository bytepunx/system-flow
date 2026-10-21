// S-0459: each Configuration page at /settings/<page> shows its title and only its sections, an
// unknown page is a 404, and /settings, with or without an old anchor (S-0458), goes on to the page
// that holds the section, in place of itself in the history and keeping the project.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { isHttpError } from '@sveltejs/kit';
import { SETTINGS_PAGES, settingsPage, type SettingsPage } from '$lib/settingspages';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$lib/project.svelte', () => ({ projectState: { list: [], refresh: async () => {} } }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
const goto = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$app/navigation', () => ({ goto }));
const where = vi.hoisted(() => ({ url: new URL('http://localhost/settings') }));
vi.mock('$app/state', () => ({ page: where }));

import SettingsRedirect from './+page.svelte';
import SettingsPageView from './[page]/+page.svelte';
import { load } from './[page]/+page';

const settle = async () => {
	for (let i = 0; i < 8; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const shown = () =>
	[...document.querySelectorAll('[data-testid^="section-"]')].map((s) =>
		s.getAttribute('data-testid')
	);
const VIEW = {
	here: true,
	everywhere: true,
	enable: 'flai serve enable settings',
	enable_everywhere: 'flai serve enable settings --all-projects',
	host: {
		actions: [{ name: 'settings', means: 'change settings', here: true, everywhere: true }],
		agent: { name: 'builder', command: ['run-agent', '{story}'], harnesses: {} },
		checks: { commands: [], timeout_minutes: 15 },
		manifest_checks: [],
		import_roots: [],
		mcp: { running: false },
		// said, so that the pages that hold them show them
		shared: [],
		planning: { plan: false, replan: 'deterministic', replan_set: false, schedule: '' }
	}
};
const loaded = (slug: string) =>
	(load as (e: { params: { page: string } }) => SettingsPage)({ params: { page: slug } });

describe('the Configuration pages (S-0459)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.head.innerHTML = '';
		api.mockReset();
		goto.mockClear();
	});

	describe.each(SETTINGS_PAGES.map((p) => [p.slug, p] as const))('/settings/%s', (slug, page) => {
		it('loads the page, and shows its title and only its sections', async () => {
			api.mockResolvedValue({ ok: true, status: 200, statusText: '', json: async () => VIEW });
			const data = loaded(slug);
			expect(data).toBe(page);
			c = mount(SettingsPageView, {
				target: document.body,
				props: { data, params: { page: slug } }
			});
			await settle();
			expect(document.querySelector('h1')!.textContent).toBe(page.title);
			expect(document.title).toBe(`flaiover — ${page.title}`);
			expect(shown()).toEqual(page.sections.map((s) => `section-${s}`));
		});
	});

	it('shows the next page’s sections on going to it, without reading the settings again', async () => {
		api.mockResolvedValue({ ok: true, status: 200, statusText: '', json: async () => VIEW });
		const props = $state({ data: settingsPage('features')!, params: { page: 'features' } });
		c = mount(SettingsPageView, { target: document.body, props });
		await settle();
		expect(shown()).toEqual(['section-actions', 'section-checks']);
		props.data = settingsPage('harnesses')!;
		props.params = { page: 'harnesses' };
		await settle();
		expect(document.querySelector('h1')!.textContent).toBe('Harnesses');
		expect(shown()).toEqual(['section-agent']);
		expect(api).toHaveBeenCalledTimes(1);
	});

	it('is a 404 naming the pages for an unknown page', () => {
		let thrown: unknown;
		try {
			loaded('tokens');
		} catch (e) {
			thrown = e;
		}
		expect(isHttpError(thrown, 404)).toBe(true);
		expect((thrown as { body: { message: string } }).body.message).toBe(
			'No settings page "tokens": open one of /settings/features, /settings/agents, ' +
				'/settings/file-system, /settings/projects, /settings/harnesses'
		);
	});

	describe('/settings', () => {
		const arrive = (url: string) => {
			where.url = new URL(url);
			c = mount(SettingsRedirect, { target: document.body });
			flushSync();
		};

		it('goes on to Features in place of itself', () => {
			arrive('http://localhost/settings');
			expect(goto).toHaveBeenCalledExactlyOnceWith('/settings/features', { replaceState: true });
		});

		it.each(SETTINGS_PAGES.map((p) => p.slug))('goes on from the old #%s to its page', (slug) => {
			arrive(`http://localhost/settings#${slug}`);
			expect(goto).toHaveBeenCalledExactlyOnceWith(`/settings/${slug}`, { replaceState: true });
		});

		it('goes on to Features from an anchor no page has', () => {
			arrive('http://localhost/settings#tokens');
			expect(goto).toHaveBeenCalledExactlyOnceWith('/settings/features', { replaceState: true });
		});

		it('keeps the project the URL names', () => {
			arrive('http://localhost/settings?project=a%20b#harnesses');
			expect(goto).toHaveBeenCalledExactlyOnceWith('/settings/harnesses?project=a%20b', {
				replaceState: true
			});
		});
	});
});
