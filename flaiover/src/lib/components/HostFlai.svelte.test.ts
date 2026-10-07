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
		hostFlai.status = {
			configured: true,
			connected: true,
			flai: '1.9.0',
			since: '2026-09-22T03:15:42.123Z'
		};
		flushSync();
		expect(badge()?.getAttribute('data-host-flai')).toBe('connected');
		expect(badge()?.textContent).toContain('connected');
		// since when, in the local zone (S-0329): the tests run in New York
		expect(badge()?.getAttribute('title')).toBe(
			'flai 1.9.0 on the host, connected since 2026-09-21 23:15 EDT'
		);
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

	it('says when a connected flai is older than the newest flai release in the project, and how to upgrade it', () => {
		const c = mount(HostFlai, { target: document.body });
		hostFlai.status = {
			configured: true,
			connected: true,
			flai: '1.26.4',
			info: {
				flai_outdated: {
					running: '1.26.4',
					newest: '1.27.0',
					upgrade: 'flai host upgrade (flai self-upgrade where no flai host runs)',
					message: 'this flai is 1.26.4, older than flai 1.27.0'
				}
			}
		};
		flushSync();
		expect(badge()?.getAttribute('data-host-flai')).toBe('behind');
		expect(badge()?.textContent).toContain('1.27.0');
		expect(badge()?.getAttribute('title')).toContain(
			'flai 1.26.4 on the host is older than flai 1.27.0'
		);
		expect(badge()?.getAttribute('title')).toContain('flai host upgrade');
		unmount(c);
	});
});
