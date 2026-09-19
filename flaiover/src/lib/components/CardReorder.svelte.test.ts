import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import CardReorder from './CardReorder.svelte';

type Props = ComponentProps<typeof CardReorder>;

describe('CardReorder', () => {
	let component: ReturnType<typeof mount> | undefined;
	const render = (props: Props) => {
		component = mount(CardReorder, { target: document.body, props });
		flushSync();
	};
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const buttons = () => [...document.querySelectorAll<HTMLButtonElement>('button')];

	it('offers up and down with names a screen reader can use, and asks for the placement it was given', () => {
		const onplace = vi.fn();
		render({ id: 'S-0059', up: { before: 'S-0061' }, down: { after: 'S-0060' }, onplace });
		expect(buttons().map((b) => b.getAttribute('aria-label'))).toEqual([
			'Move S-0059 up in the pull order',
			'Move S-0059 down in the pull order'
		]);
		expect(buttons().every((b) => b.type === 'button')).toBe(true);
		buttons()[0].click();
		expect(onplace).toHaveBeenLastCalledWith({ before: 'S-0061' }, 'up');
		buttons()[1].click();
		expect(onplace).toHaveBeenLastCalledWith({ after: 'S-0060' }, 'down');
	});

	it('has no up on the first card and no down on the last', () => {
		render({ id: 'S-0061', up: null, down: { after: 'S-0059' }, onplace: vi.fn() });
		expect(buttons().map((b) => b.dataset.step)).toEqual(['down']);
	});

	it('renders nothing for a card with nowhere to go', () => {
		render({ id: 'S-0061', up: null, down: null, onplace: vi.fn() });
		expect(document.querySelector('[data-testid="reorder"]')).toBeNull();
	});

	it('is hidden until the card is hovered or holds keyboard focus', () => {
		render({ id: 'S-0059', up: { top: true }, down: null, onplace: vi.fn() });
		const classes = document.querySelector('[data-testid="reorder"]')!.className.split(/\s+/);
		expect(classes).toEqual(
			expect.arrayContaining(['hidden', 'group-hover:flex', 'group-focus-within:flex'])
		);
	});
});
