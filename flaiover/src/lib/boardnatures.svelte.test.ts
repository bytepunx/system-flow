// S-0302: which natures the board shows, remembered per browser, every nature by default.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BoardNatures, parseShownNatures } from './boardnatures.svelte';

const KEY = 'flaiover-board-natures';
const all = {
	feature: true,
	improvement: true,
	remediation: true,
	research: true,
	experiment: true
};

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('parseShownNatures', () => {
	it('shows every nature when nothing is stored or what is stored is not a choice', () => {
		for (const raw of [null, '', 'not json', '[]', '"feature"', 'null', '{"feature":"no"}'])
			expect(parseShownNatures(raw)).toEqual(all);
	});

	it('reads a stored choice, and the default for a nature it leaves out or garbles', () => {
		expect(
			parseShownNatures(
				'{"feature":false,"improvement":true,"remediation":false,"research":true,"experiment":false}'
			)
		).toEqual({
			feature: false,
			improvement: true,
			remediation: false,
			research: true,
			experiment: false
		});
		expect(parseShownNatures('{"research":false,"experiment":0,"spike":false}')).toEqual({
			...all,
			research: false
		});
	});
});

describe('BoardNatures', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('starts with every nature shown in a fresh browser', () => {
		expect(new BoardNatures().shown).toEqual(all);
	});

	it('writes a toggle at once, and a new page reads it back', () => {
		const first = new BoardNatures();
		first.toggle('research');
		first.set('feature', false);
		const expected = { ...all, research: false, feature: false };
		expect(first.shown).toEqual(expected);
		expect(first.isShown('research')).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual(expected);
		expect(new BoardNatures().shown).toEqual(expected);
		first.toggle('research');
		expect(first.isShown('research')).toBe(true);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ ...all, feature: false });
	});

	it('always shows a nature the schema does not know', () => {
		const n = new BoardNatures();
		for (const nature of [...Object.keys(all)] as (keyof typeof all)[]) n.set(nature, false);
		expect(n.isShown('spike')).toBe(true);
		expect(n.isShown('')).toBe(true);
		expect(n.isShown('feature')).toBe(false);
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
		const n = new BoardNatures();
		expect(n.shown).toEqual(all);
		n.toggle('experiment');
		expect(n.shown.experiment).toBe(false);
		expect(n.isShown('experiment')).toBe(false);
	});
});
