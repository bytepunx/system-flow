import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostFlaiBanner from './HostFlaiBanner.svelte';
import { hostFlai } from '$lib/hostflai.svelte';

describe('HostFlaiBanner', () => {
	afterEach(() => {
		hostFlai.status = null;
		document.body.innerHTML = '';
	});
	const banner = () => document.querySelector('[data-host-flai-banner]');

	it('says nothing until the status is known, and nothing while a flai is connected', () => {
		const c = mount(HostFlaiBanner, { target: document.body });
		expect(banner()).toBeNull();
		hostFlai.status = { configured: true, connected: true, flai: '1.5.0' };
		flushSync();
		expect(banner()).toBeNull();
		unmount(c);
	});

	it('says what is missing and the commands that bring it back', () => {
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = { configured: true, connected: false };
		flushSync();
		const text = banner()?.textContent ?? '';
		expect(banner()?.getAttribute('role')).toBe('alert');
		expect(text).toContain('No flai on the host is connected');
		expect(text).toContain('flai dashboard');
		expect(text).toContain('flai serve start');
		expect(text).toContain('flai serve status');

		// connected but not answering is not usable either, and the reason is shown
		hostFlai.status = { configured: true, connected: true, error: 'did not answer in time' };
		flushSync();
		expect(banner()?.textContent).toContain('did not answer in time');
		unmount(c);
	});

	it('tells a container started by an older flai apart', () => {
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = { configured: false, connected: false };
		flushSync();
		expect(banner()?.textContent).toContain('started by a flai from before');
		expect(banner()?.textContent).toContain('flai dashboard stop');
		unmount(c);
	});
});
