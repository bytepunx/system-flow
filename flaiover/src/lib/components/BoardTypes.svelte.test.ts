import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import BoardTypes from './BoardTypes.svelte';
import { BoardTypes as Types } from '$lib/boardtypes.svelte';

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('BoardTypes', () => {
	let component: ReturnType<typeof mount> | undefined;
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});

	const boxes = () => [...document.querySelectorAll<HTMLInputElement>('input[type=checkbox]')];

	it('has a checkbox for epics, stories, and tasks, with stories alone ticked', () => {
		component = mount(BoardTypes, { target: document.body, props: { types: new Types() } });
		flushSync();
		expect(
			boxes().map((b) => [b.dataset.type, b.parentElement!.textContent!.trim(), b.checked])
		).toEqual([
			['epic', 'epics', false],
			['story', 'stories', true],
			['task', 'tasks', false]
		]);
	});

	it('toggles a type when its box is clicked and remembers it for the next page', () => {
		const types = new Types();
		component = mount(BoardTypes, { target: document.body, props: { types } });
		flushSync();
		const [epic, story] = boxes();
		epic.click();
		story.click();
		flushSync();
		expect(types.shown).toEqual({ epic: true, story: false, task: false });
		unmount(component);
		component = mount(BoardTypes, { target: document.body, props: { types: new Types() } });
		flushSync();
		expect(boxes().map((b) => b.checked)).toEqual([true, false, false]);
	});
});
