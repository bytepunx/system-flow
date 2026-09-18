import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount, type ComponentProps } from 'svelte';
import BoardCard from './BoardCard.svelte';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id)
}));

const base = {
	id: 'S-0048',
	type: 'story',
	title: 'Parent epic indicator on story cards on the board',
	nature: 'improvement',
	blocked: false,
	age_seconds: 3600
};

type Props = ComponentProps<typeof BoardCard>;

function render(card: Props['card'], props: Partial<Props> = {}) {
	const component = mount(BoardCard, { target: document.body, props: { card, ...props } });
	flushSync();
	return component;
}
const parent = () => document.querySelector<HTMLElement>('[data-testid="parent"]');

describe('BoardCard', () => {
	let component: ReturnType<typeof render> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});

	it('shows a story’s parent epic at the end of the bottom row, named by its title', () => {
		component = render({
			...base,
			parent: 'E-0006',
			parent_title: 'flaiover as the designer’s workbench'
		});
		const el = parent()!;
		expect(el.textContent).toBe('E-0006');
		expect(el.getAttribute('title')).toBe('E-0006 flaiover as the designer’s workbench');
		expect(el.getAttribute('aria-label')).toBe(
			'parent E-0006 flaiover as the designer’s workbench'
		);
		// last in the row that holds the nature, pushed to the right edge
		const row = el.parentElement!;
		expect(row.lastElementChild).toBe(el);
		expect(row.firstElementChild!.textContent).toBe('improvement');
		expect(el.className).toContain('ml-auto');
	});

	it('shows nothing, and leaves no empty element, for a story without a parent', () => {
		component = render(base);
		expect(parent()).toBeNull();
		const row = document.querySelector('a')!.lastElementChild!;
		expect([...row.children].map((c) => c.textContent)).toEqual(['improvement']);
	});

	it('shows a task’s parent story the same way, after its type', () => {
		component = render({
			...base,
			id: 'T-0160',
			type: 'task',
			parent: 'S-0048',
			parent_title: 'Parent epic indicator'
		});
		expect(parent()!.textContent).toBe('S-0048');
		expect(parent()!.getAttribute('title')).toBe('S-0048 Parent epic indicator');
		const row = parent()!.parentElement!;
		expect([...row.children].map((c) => c.textContent)).toEqual(['improvement', 'task', 'S-0048']);
	});

	it('shows nothing for an epic', () => {
		component = render({ ...base, id: 'E-0006', type: 'epic' });
		expect(parent()).toBeNull();
	});

	it('falls back to the ID alone when the parent’s title is unknown', () => {
		component = render({ ...base, parent: 'E-0099' });
		expect(parent()!.getAttribute('title')).toBe('E-0099');
		expect(parent()!.getAttribute('aria-label')).toBe('parent E-0099');
	});

	it('keeps the blocked flag on the left of the parent', () => {
		component = render({ ...base, blocked: true, parent: 'E-0006' });
		const row = parent()!.parentElement!;
		expect([...row.children].map((c) => c.textContent)).toEqual([
			'improvement',
			'BLOCKED',
			'E-0006'
		]);
	});

	it('stays one draggable link to the item, with no link inside it', () => {
		const ondragstart = vi.fn();
		component = render({ ...base, parent: 'E-0006' }, { draggable: true, ondragstart });
		const links = document.querySelectorAll('a');
		expect(links).toHaveLength(1);
		expect(links[0].getAttribute('href')).toBe('/items/S-0048');
		expect(links[0].getAttribute('draggable')).toBe('true');
		links[0].dispatchEvent(new Event('dragstart', { bubbles: true }));
		expect(ondragstart).toHaveBeenCalledOnce();
	});
});
