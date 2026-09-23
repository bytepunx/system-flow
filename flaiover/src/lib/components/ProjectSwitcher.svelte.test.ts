import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ProjectSwitcher from './ProjectSwitcher.svelte';
import { projectState, resetForTests } from '$lib/project.svelte';

describe('ProjectSwitcher (S-0080)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		resetForTests();
		vi.stubGlobal('location', { ...window.location, reload: vi.fn() });
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.unstubAllGlobals();
	});

	it('shows nothing before any project is known, and names the only one once there is one (S-0095)', () => {
		c = mount(ProjectSwitcher, { target: document.body });
		expect(document.querySelector('[data-testid="project-switcher"]')).toBeNull();
		expect(document.querySelector('[data-testid="project-name"]')).toBeNull();
		projectState.list = [{ key: 'harbour', name: 'Harbour', connected: true }];
		flushSync();
		expect(document.querySelector('[data-testid="project-switcher"]')).toBeNull();
		expect(document.querySelector('[data-testid="project-name"]')?.textContent).toBe('Harbour');
	});

	it('lists every known project once there is more than one, current or not connected', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'quay', name: 'Quay', connected: false }
		];
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const options = [...document.querySelectorAll('option')].map((o) => o.textContent);
		expect(options).toContain('Harbour');
		expect(options).toContain('Quay (not connected)');
	});

	it('picking an option switches to that project in place, without a reload (S-0095)', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'quay', name: 'Quay', connected: true }
		];
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const select = document.querySelector('select')!;
		select.value = 'quay';
		select.dispatchEvent(new Event('change', { bubbles: true }));
		expect(projectState.current).toBe('quay');
		expect(location.reload).not.toHaveBeenCalled();
	});
});
