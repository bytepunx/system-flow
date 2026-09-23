import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import HostFlaiBanner from './HostFlaiBanner.svelte';
import { hostFlai } from '$lib/hostflai.svelte';
import { projectState, resetForTests } from '$lib/project.svelte';

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
		unmount(c);
	});

	// Found live (S-0095): with two projects connected and none chosen yet, the status is asked for
	// no project in particular and says "not connected", which was false. Until a choice is made the
	// banner says nothing; once one is, it says what is true of that project.
	it('says nothing while several projects are known and none is chosen', () => {
		projectState.list = [
			{ key: 'alpha', name: 'Alpha', connected: true },
			{ key: 'beta', name: 'Beta', connected: true }
		];
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = { configured: true, connected: false };
		flushSync();
		expect(banner()).toBeNull();
		projectState.pick('beta', false);
		flushSync();
		expect(banner()?.textContent).toContain('No flai on the host is connected');
		unmount(c);
		resetForTests();
	});

	// Found live (S-0086): a flai still running an old in-memory version after the binary on disk
	// was upgraded reported as "not connected" to the operator, when /api/agent already carried the
	// true story (connected: true, and why it cannot be used). A connected flai never gets the
	// not-connected wording, whatever the reason it cannot be used.
	it('tells a connected flai that lacks what the dashboard needs apart from one never connected', () => {
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = {
			configured: true,
			connected: true,
			flai: '1.7.0',
			error:
				'flai 1.7.0 on the host is older than this dashboard and lacks 2 of the things it asks for (checks.run, checks.cancel); upgrade flai on the host, then run flai serve stop and flai dashboard'
		};
		flushSync();
		const text = banner()?.textContent ?? '';
		expect(text).toContain('flai on the host is connected');
		expect(text).not.toContain('No flai on the host is connected');
		expect(text).toContain('1.7.0');
		expect(text).toContain('checks.run');
		expect(text).toContain('flai serve stop');
		unmount(c);
	});

	// Connected but not answering (a timeout, say) is the same "connected but unusable" shape,
	// whatever the reason in status.error.
	it('shows a connected flai that failed to answer the same way, not as not-connected', () => {
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = { configured: true, connected: true, error: 'did not answer in time' };
		flushSync();
		const text = banner()?.textContent ?? '';
		expect(text).not.toContain('No flai on the host is connected');
		expect(text).toContain('did not answer in time');
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

	// Found live (S-0098): a repository offered for import answers no project.info and, once
	// imported, is no longer connected as a candidate; neither is news the banner should give,
	// since the import prompt says what there is to say.
	it('says nothing for a repository offered for import, connected or not', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'import-widget', name: 'widget', connected: true, candidate: true }
		];
		projectState.pick('import-widget', false);
		const c = mount(HostFlaiBanner, { target: document.body });
		hostFlai.status = { configured: true, connected: false };
		flushSync();
		expect(banner()).toBeNull();
		hostFlai.status = {
			configured: true,
			connected: true,
			error: 'flai offers no method project.info'
		};
		flushSync();
		expect(banner()).toBeNull();
		unmount(c);
		resetForTests();
	});
});
