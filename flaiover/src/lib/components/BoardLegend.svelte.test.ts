import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import BoardLegend from './BoardLegend.svelte';
import { natureShownTint, natureTint, typeSwatch } from '$lib/cardcolour';
import { boardNatures, natures } from '$lib/boardnatures.svelte';

const KEY = 'flaiover-board-natures';

/** boardNatures is a module singleton: put every nature back on and forget the choice. */
function reset() {
	for (const n of natures) boardNatures.set(n, true);
	localStorage.clear();
}

function tag(nature: string): HTMLButtonElement {
	return document.querySelector<HTMLButtonElement>(`button[data-nature="${nature}"]`)!;
}

describe('BoardLegend', () => {
	let component: ReturnType<typeof mount> | undefined;
	beforeEach(reset);
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
		reset();
	});

	it('names every nature on its own tint and every type beside its colour', () => {
		component = mount(BoardLegend, { target: document.body });
		flushSync();
		const tags = [...document.querySelectorAll<HTMLElement>('[data-nature]')];
		expect(tags.map((n) => n.textContent!.trim())).toEqual([
			'feature',
			'improvement',
			'remediation',
			'research',
			'experiment'
		]);
		for (const n of tags)
			expect(n.classList.contains(natureShownTint[n.dataset.nature!])).toBe(true);
		const types = [...document.querySelectorAll<HTMLElement>('[data-type]')];
		expect(types.map((t) => t.textContent!.trim())).toEqual(['epic', 'story', 'task']);
		for (const t of types)
			expect(t.querySelector('span')!.className).toContain(typeSwatch[t.dataset.type!]);
	});

	it('starts with every nature shown: pressed, on its brighter tint', () => {
		component = mount(BoardLegend, { target: document.body });
		flushSync();
		for (const n of natures) {
			const b = tag(n);
			expect(b.tagName).toBe('BUTTON');
			expect(b.type).toBe('button');
			expect(b.getAttribute('aria-pressed')).toBe('true');
			expect(b.classList.contains(natureShownTint[n])).toBe(true);
			expect(b.classList.contains(natureTint[n])).toBe(false);
			expect(b.title).toBe(`click to hide ${n}`);
		}
	});

	it('a click hides a nature on its default tint and remembers it; a second shows it again', () => {
		component = mount(BoardLegend, { target: document.body });
		flushSync();
		tag('feature').click();
		flushSync();
		let b = tag('feature');
		expect(b.getAttribute('aria-pressed')).toBe('false');
		expect(b.classList.contains(natureTint.feature)).toBe(true);
		expect(b.classList.contains(natureShownTint.feature)).toBe(false);
		expect(b.title).toBe('click to show feature');
		expect(boardNatures.shown.feature).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({
			feature: false,
			improvement: true,
			remediation: true,
			research: true,
			experiment: true
		});
		expect(tag('improvement').getAttribute('aria-pressed')).toBe('true');

		tag('feature').click();
		flushSync();
		b = tag('feature');
		expect(b.getAttribute('aria-pressed')).toBe('true');
		expect(b.classList.contains(natureShownTint.feature)).toBe(true);
		expect(b.classList.contains(natureTint.feature)).toBe(false);
		expect(b.title).toBe('click to hide feature');
		expect(JSON.parse(localStorage.getItem(KEY)!).feature).toBe(true);
	});
});
