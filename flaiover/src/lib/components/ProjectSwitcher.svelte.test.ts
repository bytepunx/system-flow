import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ProjectSwitcher from './ProjectSwitcher.svelte';
import { projectState, resetForTests, unsignedLabel, type Project } from '$lib/project.svelte';
import { health, type FlaiTrust } from '$lib/settings';

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
		expect(why?.getAttribute('href')).toBe('/settings/projects');
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

	it('says per project whether its flai is signed, and names a refused one with the reason (S-0237)', () => {
		const list: (Project & FlaiTrust)[] = [
			{ key: 'harbour', name: 'Harbour', connected: true, signed: true },
			{ key: 'pier', name: 'Pier', connected: true, signed: false },
			{
				key: 'quay',
				name: 'Quay',
				connected: false,
				served: true,
				lastError: 'the dashboard refused this flai: unsigned peer: missing stamp',
				refused: { at: '2026-10-12T17:00:00Z', reason: 'missing stamp', flai: '1.43.0' }
			},
			{
				key: 'wharf',
				name: 'Wharf',
				connected: false,
				refused: { at: '2026-10-12T17:00:00Z', reason: 'invalid signature', flai: '' }
			}
		];
		projectState.list = list;
		projectState.pick('harbour', false);
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		const options = [...document.querySelectorAll('option')];
		expect(options.map((o) => o.textContent)).toEqual([
			'Harbour (signed)',
			'Pier (unsigned)',
			'Quay (refused: unsigned flai)',
			'Wharf (refused: unsigned flai)'
		]);
		expect(options[2].title).toBe('refused: unsigned flai: missing stamp (flai 1.43.0)');
		const why = document.querySelector<HTMLAnchorElement>('[data-testid="project-switcher-why"]');
		expect(why?.textContent).toBe('2 not connected: why?');
		expect(why?.title).toBe(
			'Quay: refused: unsigned flai: missing stamp (flai 1.43.0)\n' +
				'Wharf: refused: unsigned flai: invalid signature'
		);
	});

	it('names the side a connection was let through unsigned in place of signed or unsigned; a signed pair keeps signed (S-0239)', () => {
		const list: (Project & FlaiTrust)[] = [
			{ key: 'harbour', name: 'Harbour', connected: true, signed: true },
			{ key: 'pier', name: 'Pier', connected: true, signed: false, unsigned: ['flai'] },
			{ key: 'quay', name: 'Quay', connected: true, signed: true, unsigned: ['dashboard'] },
			{
				key: 'wharf',
				name: 'Wharf',
				connected: true,
				signed: false,
				unsigned: ['flai', 'dashboard']
			},
			{
				key: 'dock',
				name: 'Dock',
				connected: false,
				refused: { at: '2026-10-12T17:00:00Z', reason: 'invalid signature', flai: '9.9.8' }
			}
		];
		projectState.list = list;
		projectState.pick('harbour', false);
		c = mount(ProjectSwitcher, { target: document.body });
		flushSync();
		expect([...document.querySelectorAll('option')].map((o) => o.textContent)).toEqual([
			'Harbour (signed)',
			'Pier (unsigned allowed: flai)',
			'Quay (unsigned allowed: dashboard)',
			'Wharf (unsigned allowed: flai and dashboard)',
			'Dock (refused: unsigned flai)'
		]);
	});

	it('words the sides let through unsigned, and nothing for none (S-0239)', () => {
		expect(unsignedLabel(['flai'])).toBe('unsigned allowed: flai');
		expect(unsignedLabel(['dashboard'])).toBe('unsigned allowed: dashboard');
		expect(unsignedLabel(['flai', 'dashboard'])).toBe('unsigned allowed: flai and dashboard');
		expect(unsignedLabel([])).toBe('');
		expect(unsignedLabel(undefined)).toBe('');
	});

	it('the Settings page says a refused project’s state in the switcher’s words (S-0237)', () => {
		const served = { key: 'quay', name: 'Quay', root: '/r/quay', settings: true } as const;
		expect(
			health({
				...served,
				state: 'not-connected',
				last_error: 'the dashboard refused this flai: unsigned peer: missing stamp'
			})
		).toBe('not connected: refused: unsigned flai: missing stamp');
		expect(health({ ...served, state: 'not-connected', last_error: 'dial: refused' })).toBe(
			'not connected: dial: refused'
		);
	});
});
