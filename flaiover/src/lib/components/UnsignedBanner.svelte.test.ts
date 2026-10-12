import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import UnsignedBanner from './UnsignedBanner.svelte';
import { projectState, resetForTests, type Project } from '$lib/project.svelte';

const FLAI_WORDS = 'The flai on the host is a development build, with no release stamp.';
const DASHBOARD_WORDS = 'This dashboard is a development build, with no release stamp.';

describe('UnsignedBanner', () => {
	afterEach(() => {
		resetForTests();
		document.body.innerHTML = '';
	});
	const banner = () => document.querySelector('[data-unsigned-banner]');
	const show = (project: Project) => {
		projectState.list = [project];
		projectState.pick(project.key, false);
		return mount(UnsignedBanner, { target: document.body });
	};

	it('names the flai on the host when it is the unsigned side', () => {
		const c = show({ key: 'alpha', name: 'Alpha', connected: true, unsigned: ['flai'] });
		const text = banner()?.textContent ?? '';
		expect(banner()?.getAttribute('role')).toBe('status');
		expect(text).toContain('unsigned allowed: flai');
		expect(text).toContain(FLAI_WORDS);
		expect(text).not.toContain(DASHBOARD_WORDS);
		expect(text).toContain('dashboard.allow_unsigned');
		unmount(c);
	});

	it('names this dashboard when it is the unsigned side', () => {
		const c = show({ key: 'alpha', name: 'Alpha', connected: true, unsigned: ['dashboard'] });
		const text = banner()?.textContent ?? '';
		expect(text).toContain('unsigned allowed: dashboard');
		expect(text).toContain(DASHBOARD_WORDS);
		expect(text).not.toContain(FLAI_WORDS);
		unmount(c);
	});

	it('names both sides when both are unsigned', () => {
		const c = show({
			key: 'alpha',
			name: 'Alpha',
			connected: true,
			unsigned: ['flai', 'dashboard']
		});
		const text = banner()?.textContent ?? '';
		expect(text).toContain('unsigned allowed: flai and dashboard');
		expect(text).toContain(FLAI_WORDS);
		expect(text).toContain(DASHBOARD_WORDS);
		unmount(c);
	});

	// Allowed is never silent, so there is nothing to dismiss it with.
	it('offers no way to dismiss it', () => {
		const c = show({ key: 'alpha', name: 'Alpha', connected: true, unsigned: ['flai'] });
		expect(banner()?.querySelector('button')).toBeNull();
		unmount(c);
	});

	it('says nothing for a signed pair', () => {
		const c = show({ key: 'alpha', name: 'Alpha', connected: true });
		expect(banner()).toBeNull();
		unmount(c);
	});

	it('says nothing for a project with no connection', () => {
		const c = show({ key: 'alpha', name: 'Alpha', connected: false, unsigned: ['flai'] });
		expect(banner()).toBeNull();
		unmount(c);
	});

	it('says nothing with no project picked, and follows the pick', () => {
		projectState.list = [
			{ key: 'alpha', name: 'Alpha', connected: true },
			{ key: 'beta', name: 'Beta', connected: true, unsigned: ['dashboard'] }
		];
		const c = mount(UnsignedBanner, { target: document.body });
		expect(banner()).toBeNull();
		projectState.pick('beta', false);
		flushSync();
		expect(banner()?.textContent).toContain(DASHBOARD_WORDS);
		projectState.pick('alpha', false);
		flushSync();
		expect(banner()).toBeNull();
		unmount(c);
	});
});
