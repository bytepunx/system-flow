// S-0141, S-0303: which work item types the board shows, remembered per browser, every type by
// default.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BoardTypes, parseShown } from './boardtypes.svelte';

const KEY = 'flaiover-board-types';
const all = { epic: true, story: true, task: true };

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('parseShown', () => {
	it('shows every type when nothing is stored or what is stored is not a choice', () => {
		for (const raw of [null, '', 'not json', '[]', '"story"', 'null', '{"story":"yes"}'])
			expect(parseShown(raw)).toEqual(all);
	});

	it('reads a stored choice back as it was, and shows a type it leaves out', () => {
		expect(parseShown('{"epic":false,"story":true,"task":false}')).toEqual({
			epic: false,
			story: true,
			task: false
		});
		expect(parseShown('{"epic":true,"story":false,"task":true}')).toEqual({
			epic: true,
			story: false,
			task: true
		});
		expect(parseShown('{"task":false}')).toEqual({ epic: true, story: true, task: false });
	});
});

describe('BoardTypes', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('starts with every type shown in a fresh browser', () => {
		expect(new BoardTypes().shown).toEqual(all);
	});

	it('keeps a choice stored before the default changed', () => {
		localStorage.setItem(KEY, '{"epic":false,"story":true,"task":false}');
		expect(new BoardTypes().shown).toEqual({ epic: false, story: true, task: false });
	});

	it('writes a change at once, and a new page reads it back', () => {
		const first = new BoardTypes();
		first.set('epic', false);
		first.set('task', false);
		expect(first.shown).toEqual({ epic: false, story: true, task: false });
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({
			epic: false,
			story: true,
			task: false
		});
		expect(new BoardTypes().shown).toEqual({ epic: false, story: true, task: false });
	});

	it('toggles a type hidden and back, writing the whole choice each time', () => {
		const types = new BoardTypes();
		types.toggle('task');
		expect(types.shown).toEqual({ epic: true, story: true, task: false });
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({
			epic: true,
			story: true,
			task: false
		});
		types.toggle('task');
		expect(types.shown).toEqual(all);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(all);
	});

	it('keeps the choice for the page when storage refuses it', () => {
		vi.stubGlobal('localStorage', {
			getItem: () => {
				throw new Error('denied');
			},
			setItem: () => {
				throw new Error('denied');
			}
		});
		const types = new BoardTypes();
		expect(types.shown).toEqual(all);
		types.toggle('task');
		expect(types.shown.task).toBe(false);
	});
});
