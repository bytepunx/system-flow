// S-0080: which project this tab looks at, remembered per browser, with nothing to choose while
// at most one project is connected.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const STORAGE_KEY = 'flaiover-project';

function freshLocation(href: string) {
	const url = new URL(href);
	// jsdom's location cannot be reassigned piecemeal; replace the whole object for each test.
	Object.defineProperty(window, 'location', {
		value: { ...url, href: url.href, search: url.search, pathname: url.pathname },
		writable: true,
		configurable: true
	});
}

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('projectState', () => {
	beforeEach(() => {
		vi.stubGlobal('localStorage', fakeStorage());
		freshLocation('http://dash.example/board');
		vi.stubGlobal(
			'history',
			(() => {
				let current = 'http://dash.example/board';
				return {
					state: null,
					replaceState: vi.fn((_state: unknown, _title: string, url: string | URL) => {
						current = String(url);
						freshLocation(current);
					})
				};
			})()
		);
	});
	afterEach(() => {
		vi.unstubAllGlobals();
		vi.resetModules();
	});

	async function load() {
		return await import('./project.svelte');
	}

	it('starts with nothing chosen when neither the URL nor storage says otherwise', async () => {
		const { projectState } = await load();
		expect(projectState.current).toBeNull();
		expect(projectState.needsChoice).toBe(false); // no list fetched yet
	});

	it('reads the remembered choice, and the URL over it', async () => {
		localStorage.setItem(STORAGE_KEY, 'remembered');
		{
			const { projectState } = await load();
			expect(projectState.current).toBe('remembered');
		}
		vi.resetModules();
		freshLocation('http://dash.example/board?project=from-url');
		{
			const { projectState } = await load();
			expect(projectState.current).toBe('from-url');
		}
	});

	it('tags an API path with the current project, and leaves one already named alone', async () => {
		const { projectState } = await load();
		projectState.pick('harbour');
		expect(projectState.tag('/api/board')).toBe('/api/board?project=harbour');
		expect(projectState.tag('/api/board?x=1')).toBe('/api/board?x=1&project=harbour');
		expect(projectState.tag('/api/board?project=quay')).toBe('/api/board?project=quay');
	});

	it('tags nothing when no project is chosen', async () => {
		const { projectState } = await load();
		expect(projectState.tag('/api/board')).toBe('/api/board');
	});

	it('persists a pick to storage and to the URL, and forgetting clears both', async () => {
		const { projectState } = await load();
		projectState.pick('harbour');
		expect(localStorage.getItem(STORAGE_KEY)).toBe('harbour');
		expect(window.location.search).toContain('project=harbour');
		projectState.pick(null);
		expect(localStorage.getItem(STORAGE_KEY)).toBeNull();
		expect(window.location.search).not.toContain('project=');
	});

	it('refresh() lists what the server knows and auto-picks the one project when there is only one', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({
				ok: true,
				json: async () => ({ projects: [{ key: 'harbour', name: 'Harbour', connected: true }] })
			}))
		);
		const { projectState } = await load();
		await projectState.refresh();
		expect(projectState.list).toEqual([{ key: 'harbour', name: 'Harbour', connected: true }]);
		expect(projectState.current).toBe('harbour');
		expect(projectState.needsChoice).toBe(false);
	});

	it('refresh() asked while a request is in flight shares it, and asks again once it has settled (S-0186)', async () => {
		let answer: (() => void) | undefined;
		const fetchFn = vi.fn(
			() =>
				new Promise((resolve) => {
					answer = () =>
						resolve({
							ok: true,
							json: async () => ({
								projects: [{ key: 'harbour', name: 'Harbour', connected: true, review: 2 }]
							})
						});
				})
		);
		vi.stubGlobal('fetch', fetchFn);
		const { projectState } = await load();
		const first = projectState.refresh();
		const second = projectState.refresh();
		expect(fetchFn).toHaveBeenCalledTimes(1);
		// past the HTTP cache, whose lock would hold a second tab's request behind this one
		expect(fetchFn).toHaveBeenCalledWith('/api/projects', { cache: 'no-store' });
		answer!();
		await Promise.all([first, second]);
		expect(projectState.list[0].review).toBe(2); // the glance is kept with the project
		const third = projectState.refresh();
		expect(fetchFn).toHaveBeenCalledTimes(2);
		answer!();
		await third;
	});

	it('a failed request is not shared with a later refresh()', async () => {
		const fetchFn = vi
			.fn()
			.mockRejectedValueOnce(new Error('offline'))
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ projects: [{ key: 'harbour', name: 'Harbour', connected: true }] })
			});
		vi.stubGlobal('fetch', fetchFn);
		const { projectState } = await load();
		await projectState.refresh();
		expect(projectState.loaded).toBe(false);
		await projectState.refresh();
		expect(fetchFn).toHaveBeenCalledTimes(2);
		expect(projectState.loaded).toBe(true);
	});

	it('does not choose for the designer once there is more than one project', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({
				ok: true,
				json: async () => ({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true },
						{ key: 'quay', name: 'Quay', connected: true }
					]
				})
			}))
		);
		const { projectState } = await load();
		await projectState.refresh();
		expect(projectState.current).toBeNull();
		expect(projectState.needsChoice).toBe(true);
	});

	it('forgets a remembered choice that no longer exists once the real list is known', async () => {
		localStorage.setItem(STORAGE_KEY, 'gone');
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({
				ok: true,
				json: async () => ({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true },
						{ key: 'quay', name: 'Quay', connected: true }
					]
				})
			}))
		);
		const { projectState } = await load();
		expect(projectState.current).toBe('gone');
		await projectState.refresh();
		expect(projectState.current).toBeNull();
	});

	it('keeps a valid remembered choice once the list confirms it', async () => {
		localStorage.setItem(STORAGE_KEY, 'quay');
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => ({
				ok: true,
				json: async () => ({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true },
						{ key: 'quay', name: 'Quay', connected: false }
					]
				})
			}))
		);
		const { projectState } = await load();
		await projectState.refresh();
		expect(projectState.current).toBe('quay');
	});
});
