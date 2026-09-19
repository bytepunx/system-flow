import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import BoardLegend from './BoardLegend.svelte';
import { natureTint, typeSwatch } from '$lib/cardcolour';

describe('BoardLegend', () => {
	let component: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});

	it('names every nature on its own tint and every type beside its colour', () => {
		component = mount(BoardLegend, { target: document.body });
		flushSync();
		const natures = [...document.querySelectorAll<HTMLElement>('[data-nature]')];
		expect(natures.map((n) => n.textContent!.trim())).toEqual([
			'feature',
			'improvement',
			'remediation',
			'research',
			'experiment'
		]);
		for (const n of natures) expect(n.className).toContain(natureTint[n.dataset.nature!]);
		const types = [...document.querySelectorAll<HTMLElement>('[data-type]')];
		expect(types.map((t) => t.textContent!.trim())).toEqual(['epic', 'story', 'task']);
		for (const t of types)
			expect(t.querySelector('span')!.className).toContain(typeSwatch[t.dataset.type!]);
	});
});
