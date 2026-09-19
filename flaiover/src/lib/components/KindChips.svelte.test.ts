import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import KindChips from './KindChips.svelte';

describe('KindChips', () => {
	let component: ReturnType<typeof mount> | undefined;
	const render = (props: { type?: string; nature?: string }) => {
		component = mount(KindChips, { target: document.body, props });
		flushSync();
	};
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});

	it('writes the type beside its stripe colour and the nature on its tint', () => {
		render({ type: 'story', nature: 'remediation' });
		const type = document.querySelector<HTMLElement>('[data-type]')!;
		expect(type.textContent).toBe('story');
		expect(type.querySelector('span')!.className).toContain('bg-type-story');
		const nature = document.querySelector<HTMLElement>('[data-nature]')!;
		expect(nature.textContent).toBe('remediation');
		expect(nature.className).toContain('bg-nature-remediation');
		expect(document.body.textContent).toBe('story·remediation');
	});

	it('shows the nature alone, for a table cell', () => {
		render({ nature: 'research' });
		expect(document.querySelector('[data-type]')).toBeNull();
		expect(document.body.textContent).toBe('research');
	});

	it('keeps the plain words for values the schema does not know', () => {
		render({ type: 'theme', nature: 'chore' });
		expect(document.body.textContent).toBe('theme·chore');
		expect(document.querySelector('[data-type] span')).toBeNull();
		expect(document.querySelector<HTMLElement>('[data-nature]')!.className).toBe('');
	});
});
