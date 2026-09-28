// S-0141: which work item types the board shows, remembered per browser, stories alone by default.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BoardTypes, parseShown } from './boardtypes.svelte';

const KEY = 'flaiover-board-types';

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
	it('shows stories alone when nothing is stored or what is stored is not a choice', () => {
		const stories = { epic: false, story: true, task: false };
		for (const raw of [null, '', 'not json', '[]', '"story"', 'null', '{"story":"yes"}'])
			expect(parseShown(raw)).toEqual(stories);
	});

	it('reads a stored choice, and the default for a type it leaves out', () => {
		expect(parseShown('{"epic":true,"story":false,"task":true}')).toEqual({
			epic: true,
			story: false,
			task: true
		});
		expect(parseShown('{"task":true}')).toEqual({ epic: false, story: true, task: true });
	});
});

describe('BoardTypes', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('starts with stories alone in a fresh browser', () => {
		expect(new BoardTypes().shown).toEqual({ epic: false, story: true, task: false });
	});

	it('writes a toggle at once, and a new page reads it back', () => {
		const first = new BoardTypes();
		first.set('epic', true);
		first.set('story', false);
		expect(first.shown).toEqual({ epic: true, story: false, task: false });
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({
			epic: true,
			story: false,
			task: false
		});
		expect(new BoardTypes().shown).toEqual({ epic: true, story: false, task: false });
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
		expect(types.shown).toEqual({ epic: false, story: true, task: false });
		types.set('task', true);
		expect(types.shown.task).toBe(true);
	});
});
