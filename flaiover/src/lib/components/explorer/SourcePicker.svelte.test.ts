// S-0442: the source picker offers Main and each story in progress or in review, hands up the one
// chosen, and with a story lists the components its branch changes, marked, to pick.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import type { Component } from '$lib/server/repo';
import SourcePicker, { openStories, type ChangedComponent } from './SourcePicker.svelte';

const STORIES = [
	{ id: 'S-0441', title: 'Live tokens', status: 'in-progress' },
	{ id: 'S-0442', title: "A story's worktree", status: 'review' }
];
const component = (path: string): Component => ({
	project: 'flaiover',
	path,
	name: path.slice(path.lastIndexOf('/') + 1, -'.svelte'.length),
	props: []
});
const BUTTON = 'src/lib/components/Button.svelte';
const CHANGED: ChangedComponent[] = [{ project: 'flaiover', component: component(BUTTON) }];

describe('openStories', () => {
	it("lists the board's stories in progress, then in review, and no other card", () => {
		const card = (id: string, type: string, status: string) => ({ id, type, title: id, status });
		const columns = {
			ready: [card('S-0443', 'story', 'ready')],
			'in-progress': [
				card('S-0441', 'story', 'in-progress'),
				card('T-1881', 'task', 'in-progress')
			],
			review: [card('S-0442', 'story', 'review')]
		};
		expect(openStories(columns)).toEqual([
			{ id: 'S-0441', title: 'S-0441', status: 'in-progress' },
			{ id: 'S-0442', title: 'S-0442', status: 'review' }
		]);
		expect(openStories({})).toEqual([]);
	});
});

describe('SourcePicker (S-0442)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	const q = (testid: string) => document.querySelector(`[data-testid="${testid}"]`);
	const select = () => q('source-picker-select') as HTMLSelectElement;
	const options = () => [...select().options].map((o) => ({ value: o.value, text: o.textContent }));

	function show(more: Partial<ComponentProps<typeof SourcePicker>> = {}) {
		const onchoose = vi.fn();
		const onpick = vi.fn();
		c = mount(SourcePicker, {
			target: document.body,
			props: {
				stories: STORIES,
				story: undefined,
				onchoose,
				changed: null,
				picked: null,
				onpick,
				...more
			}
		});
		flushSync();
		return { onchoose, onpick };
	}

	it('offers Main and each story, with Main chosen and nothing else shown without a story', () => {
		show();
		expect(options()).toEqual([
			{ value: '', text: 'Main' },
			{ value: 'S-0441', text: 'S-0441 Live tokens (in-progress)' },
			{ value: 'S-0442', text: "S-0442 A story's worktree (review)" }
		]);
		expect(select().value).toBe('');
		expect(q('source-changed')).toBeNull();
	});

	it('hands up the story chosen, and undefined for Main', () => {
		const { onchoose } = show();
		select().value = 'S-0442';
		select().dispatchEvent(new Event('change', { bubbles: true }));
		expect(onchoose).toHaveBeenLastCalledWith('S-0442');
		select().value = '';
		select().dispatchEvent(new Event('change', { bubbles: true }));
		expect(onchoose).toHaveBeenLastCalledWith(undefined);
	});

	it('shows the story picked, and still offers one the board no longer lists', () => {
		show({ story: 'S-0300' });
		expect(select().value).toBe('S-0300');
		expect(options().at(-1)).toEqual({
			value: 'S-0300',
			text: 'S-0300 not in progress or in review'
		});
		expect(document.body.textContent).toContain('story/S-0300');
	});

	it("lists the components the story's branch changes, marked, and picks one", () => {
		const { onpick } = show({ story: 'S-0442', changed: CHANGED });
		const button = q(`source-changed-${BUTTON}`) as HTMLButtonElement;
		expect(button.textContent).toContain('Button');
		expect(button.querySelector('[data-testid="source-changed-mark"]')!.textContent).toBe(
			'changed'
		);
		expect(button.getAttribute('aria-current')).toBeNull();
		button.click();
		expect(onpick).toHaveBeenCalledWith('flaiover', BUTTON);
	});

	it('marks the changed component framed as current', () => {
		show({ story: 'S-0442', changed: CHANGED, picked: { project: 'flaiover', path: BUTTON } });
		expect(q(`source-changed-${BUTTON}`)!.getAttribute('aria-current')).toBe('true');
	});

	it('says while the branch is read, when it changes no component, and when it cannot be read', () => {
		show({ story: 'S-0442' });
		expect(q('source-changed')!.textContent).toContain('Reading its branch');
		unmount(c!);
		show({ story: 'S-0442', changed: [] });
		expect(q('source-changed-none')).not.toBeNull();
		unmount(c!);
		show({ story: 'S-0442', changed: [], changedError: 'S-0442 has no branch story/S-0442' });
		expect(q('source-changed-error')!.textContent).toContain('has no branch');
	});
});
