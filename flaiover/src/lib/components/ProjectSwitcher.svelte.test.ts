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

	it('is always a control: disabled with nothing connected, and a choice with one project (S-0102)', () => {
		c = mount(ProjectSwitcher, { target: document.body });
		const select = () =>
			document.querySelector<HTMLSelectElement>('[data-testid="project-switcher"] select');
		expect(select()?.disabled).toBe(true);
		projectState.list = [{ key: 'harbour', name: 'Harbour', connected: true }];
		projectState.pick('harbour', false);
		flushSync();
		expect(select()?.disabled).toBe(false);
		expect(select()?.value).toBe('harbour');
		expect([...select()!.options].map((o) => o.textContent)).toEqual(['Harbour']);
	});

	it('offers the repositories to import in a group of their own (S-0102)', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'import-kickr', name: 'kickr', connected: true, candidate: true },
			{ key: 'import-loci', name: 'loci', connected: true, candidate: true }
		];
		projectState.pick('harbour', false);
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const group = document.querySelector('optgroup');
		expect(group?.getAttribute('label')).toBe('Not imported yet');
		expect([...group!.querySelectorAll('option')].map((o) => o.textContent)).toEqual([
			'kickr (not imported)',
			'loci (not imported)'
		]);
		// picking one is picking it: the layout then asks whether to import it
		const select = document.querySelector('select')!;
		select.value = 'import-loci';
		select.dispatchEvent(new Event('change', { bubbles: true }));
		expect(projectState.current).toBe('import-loci');
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

	it('tells a served project that is not connected apart from a connected one, and links to the settings page (S-0122)', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true, served: true },
			{ key: 'quay', name: 'Quay', connected: false, served: true, lastError: 'dial: refused' },
			{ key: 'import-loci', name: 'loci', connected: true, candidate: true }
		];
		projectState.pick('harbour', false);
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const options = [...document.querySelectorAll('option')];
		expect(options.map((o) => o.textContent)).toEqual([
			'Harbour',
			'Quay (served, not connected)',
			'loci (not imported)'
		]);
		expect(options[1].title).toBe('dial: refused');
		const why = document.querySelector<HTMLAnchorElement>('[data-testid="project-switcher-why"]');
		expect(why?.getAttribute('href')).toBe('/settings');
		expect(why?.textContent).toBe('1 not connected: why?');
		expect(why?.title).toBe('Quay: dial: refused');
	});

	it('shows no settings link while every project is connected', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'import-loci', name: 'loci', connected: false, candidate: true }
		];
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		expect(document.querySelector('[data-testid="project-switcher-why"]')).toBeNull();
	});

	it('from a project that is not connected, the settings link goes by way of one that is', () => {
		projectState.list = [
			{ key: 'harbour', name: 'Harbour', connected: true },
			{ key: 'quay', name: 'Quay', connected: false, served: true }
		];
		projectState.pick('quay', false);
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const why = document.querySelector<HTMLAnchorElement>('[data-testid="project-switcher-why"]')!;
		why.addEventListener('click', (e) => e.preventDefault());
		why.click();
		expect(projectState.current).toBe('harbour');
	});
});
