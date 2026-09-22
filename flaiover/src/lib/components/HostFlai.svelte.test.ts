import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostFlai from './HostFlai.svelte';
import { hostFlai } from '$lib/hostflai.svelte';

vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));

describe('HostFlai', () => {
	afterEach(() => {
		hostFlai.status = null;
		document.body.innerHTML = '';
	});
	const badge = () => document.querySelector('[data-host-flai]');

	it('shows nothing when the container was given no agent credential', () => {
		const c = mount(HostFlai, { target: document.body });
		hostFlai.status = { configured: false, connected: false };
		flushSync();
		expect(badge()).toBeNull();
		unmount(c);
	});

	it('says connected when a flai answers everything the dashboard needs', () => {
		const c = mount(HostFlai, { target: document.body });
		hostFlai.status = { configured: true, connected: true, flai: '1.9.0', since: '2026-09-22' };
		flushSync();
		expect(badge()?.getAttribute('data-host-flai')).toBe('connected');
		expect(badge()?.textContent).toContain('connected');
		expect(badge()?.getAttribute('title')).toContain('1.9.0');
		unmount(c);
	});

	it('says not connected, distinctly, when no flai has ever connected', () => {
		const c = mount(HostFlai, { target: document.body });
		hostFlai.status = { configured: true, connected: false };
		flushSync();
		expect(badge()?.getAttribute('data-host-flai')).toBe('gone');
		expect(badge()?.textContent).toContain('not connected');
		unmount(c);
	});

	// Found live (S-0086): the same "not connected" wording used to cover this case too, when a
	// flai genuinely was connected (just too old for what the dashboard asks of it).
	it('tells a connected but outdated flai apart from one never connected', () => {
		const c = mount(HostFlai, { target: document.body });
		hostFlai.status = {
			configured: true,
			connected: true,
			flai: '1.7.0',
			error:
				'flai 1.7.0 on the host is older than this dashboard and lacks 2 of the things it asks for'
		};
		flushSync();
		expect(badge()?.getAttribute('data-host-flai')).toBe('outdated');
		expect(badge()?.textContent).not.toContain('not connected');
		expect(badge()?.getAttribute('title')).toContain('1.7.0');
		unmount(c);
	});
});
