import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import AcceptFindings, { findingsOf, unsavedFiles, type Finding } from './AcceptFindings.svelte';

const openEditor = vi.fn();
vi.mock('$lib/editor.svelte', () => ({
	openEditor: (...args: unknown[]) => openEditor(...args)
}));

// S-0378, S-0465: each finding opens the editor on its own checkout, the story's branch or the main
// checkout, and a saved file is marked.
const finding = (over: Partial<Finding> = {}): Finding => ({
	tier: 'markdown',
	checkout: 'story',
	story: 'S-0041',
	path: 'docs/bad.md',
	line: 2,
	column: 3,
	level: 'error',
	rule: 'MD009',
	alias: 'no-trailing-spaces',
	message: 'Trailing spaces [Expected: 0; Actual: 1]',
	...over
});
const onMain = finding({
	checkout: 'main',
	story: undefined,
	path: 'design/notes.md',
	line: 7,
	rule: 'MD041',
	alias: 'first-line-heading',
	message: 'First line in a file should be a top-level heading'
});

describe('AcceptFindings', () => {
	let component: ReturnType<typeof mount> | undefined;
	const render = (props: ComponentProps<typeof AcceptFindings>) => {
		component = mount(AcceptFindings, { target: document.body, props });
		flushSync();
	};
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
		openEditor.mockReset();
	});
	const link = (label: string) =>
		document.querySelector<HTMLButtonElement>(`button[aria-label="${label}"]`)!;
	const words = (e: Element) => (e.textContent ?? '').replace(/\s+/g, ' ').trim();
	const rows = () =>
		[...document.querySelectorAll('[data-testid="accept-finding"]')].map((e) => words(e));

	it('lists each finding with its file and line, rule, message, and the branch it is on', () => {
		render({ id: 'S-0041', status: 'review', findings: [finding(), onMain] });
		const list = document.querySelector('[data-testid="accept-findings"]')!;
		expect(words(list)).toContain('Not accepted. S-0041 is review.');
		expect(rows()[0]).toContain('docs/bad.md:2');
		expect(rows()[0]).toContain("on S-0041's branch");
		expect(rows()[0]).toContain('MD009/no-trailing-spaces');
		expect(rows()[0]).toContain('Trailing spaces [Expected: 0; Actual: 1]');
		expect(rows()[1]).toContain('design/notes.md:7');
		expect(rows()[1]).not.toContain('branch');
		expect(rows()[1]).toContain('MD041/first-line-heading');
	});

	it('says only that the item was not accepted when no status is given', () => {
		render({ id: 'S-0041', findings: [finding()] });
		const first = document.querySelector('[data-testid="accept-findings"] p')!;
		expect(words(first)).toBe('Not accepted.');
	});

	it('opens a finding on the story’s branch with that story', () => {
		render({ id: 'S-0041', status: 'review', findings: [finding(), onMain] });
		link('Edit docs/bad.md at line 2').click();
		expect(openEditor).toHaveBeenLastCalledWith(
			expect.objectContaining({ path: 'docs/bad.md', line: 2, story: 'S-0041' })
		);
	});

	it('opens a finding in the main checkout with no story', () => {
		render({ id: 'S-0041', status: 'review', findings: [finding(), onMain] });
		link('Edit design/notes.md at line 7').click();
		const target = openEditor.mock.calls.at(-1)![0];
		expect(target).toMatchObject({ path: 'design/notes.md', line: 7 });
		expect(target.story).toBeUndefined();
	});

	it('calls onsaved with the file’s key once the editor saves it', () => {
		const onsaved = vi.fn();
		render({ id: 'S-0041', status: 'review', findings: [finding(), onMain], onsaved });
		link('Edit design/notes.md at line 7').click();
		openEditor.mock.calls.at(-1)![0].onsaved({ path: 'design/notes.md', committed: true });
		expect(onsaved).toHaveBeenCalledWith(':design/notes.md');
		link('Edit docs/bad.md at line 2').click();
		openEditor.mock.calls.at(-1)![0].onsaved({ path: 'docs/bad.md', committed: true });
		expect(onsaved).toHaveBeenLastCalledWith('S-0041:docs/bad.md');
	});

	it('marks each finding of a saved file, and only those', () => {
		const again = finding({ line: 9, rule: 'MD012', alias: 'no-multiple-blanks' });
		render({
			id: 'S-0041',
			status: 'review',
			findings: [finding(), again, onMain],
			saved: ['S-0041:docs/bad.md']
		});
		const marked = [...document.querySelectorAll('[data-testid="accept-finding"]')].map(
			(e) => !!e.querySelector('[data-testid="finding-saved"]')
		);
		expect(marked).toEqual([true, true, false]);
	});

	it('shows a file the editor cannot open without a link', () => {
		render({ id: 'S-0041', status: 'review', findings: [finding({ path: 'flaiover/README.md' })] });
		expect(rows()[0]).toContain('flaiover/README.md:2');
		expect(document.querySelector('button')).toBeNull();
	});
});

describe('unsavedFiles', () => {
	it('keys each file the editor can open once, on its checkout, less those saved', () => {
		const list = [finding(), finding({ line: 9 }), onMain, finding({ path: 'flaiover/README.md' })];
		expect(unsavedFiles(list, [])).toEqual(['S-0041:docs/bad.md', ':design/notes.md']);
		expect(unsavedFiles(list, ['S-0041:docs/bad.md'])).toEqual([':design/notes.md']);
	});

	it('keeps a main-checkout file apart from the same path on a story’s branch', () => {
		const onBranch = finding({ path: 'design/notes.md' });
		expect(unsavedFiles([onBranch, onMain], [':design/notes.md'])).toEqual([
			'S-0041:design/notes.md'
		]);
	});
});

describe('findingsOf', () => {
	it('reads the findings of a refusal’s data, and none when it has none', () => {
		expect(findingsOf({ id: 'S-0041', findings: [onMain] })).toEqual([onMain]);
		expect(findingsOf({ id: 'S-0041' })).toEqual([]);
		expect(findingsOf(undefined)).toEqual([]);
	});
});
