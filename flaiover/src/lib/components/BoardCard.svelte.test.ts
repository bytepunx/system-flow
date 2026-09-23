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
const details = () => document.querySelector<HTMLElement>('[data-testid="details"]');
// a font size utility, not a colour: text-[11px], text-xs, text-sm
const sizeClass = /(^|\s)text-(\[\d+px\]|xs|sm|base)(\s|$)/;

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

	it('puts the details under a divider, in one size for the whole row', () => {
		component = render({ ...base, id: 'T-0160', type: 'task', blocked: true, parent: 'S-0048' });
		const row = details()!;
		expect(row).toBe(document.querySelector('a')!.lastElementChild);
		expect(row.className).toContain('border-t');
		// 12 px, the title's size, chosen by the operator
		expect(row.className).toMatch(/(^|\s)text-xs(\s|$)/);
		// nothing in the row sets a size of its own, the parent ID included
		expect(row.children).toHaveLength(4);
		for (const child of row.children) expect(child.className).not.toMatch(sizeClass);
	});

	it('lets a crowded row wrap instead of overflowing the card', () => {
		component = render({ ...base, blocked: true, parent: 'E-0006' });
		expect(details()!.className).toContain('flex-wrap');
		expect(parent()!.className).toContain('shrink-0');
	});

	it.each([
		['feature', 'bg-nature-feature'],
		['improvement', 'bg-nature-improvement'],
		['remediation', 'bg-nature-remediation'],
		['research', 'bg-nature-research'],
		['experiment', 'bg-nature-experiment']
	])('tints a %s card with its nature token, and still says the nature', (nature, tint) => {
		component = render({ ...base, nature });
		const card = document.querySelector('a')!;
		expect(card.className.split(/\s+/)).toContain(tint);
		expect(card.className).not.toContain('bg-ground');
		expect(card.dataset.nature).toBe(nature);
		expect(details()!.firstElementChild!.textContent).toBe(nature);
	});

	it.each([
		['epic', 'E-0006', 'border-l-type-epic'],
		['story', 'S-0048', 'border-l-type-story'],
		['task', 'T-0160', 'border-l-type-task']
	])('stripes the left edge of a %s with its type token', (type, id, stripe) => {
		component = render({ ...base, id, type });
		const classes = document.querySelector('a')!.className.split(/\s+/);
		expect(classes).toContain('border-l-4');
		expect(classes).toContain(stripe);
		// hovering strengthens the other three edges and leaves the stripe alone
		expect(classes).not.toContain('hover:border-line-strong');
		expect(classes).not.toContain('hover:border-l-line-strong');
	});

	it('keeps the plain card for a nature or a type it does not know', () => {
		component = render({ ...base, nature: 'chore', type: 'theme' });
		const classes = document.querySelector('a')!.className.split(/\s+/);
		expect(classes).toContain('bg-ground');
		expect(classes).not.toContain('border-l-4');
		expect(classes).toContain('hover:border-l-line-strong');
		expect(classes.filter((c) => /nature-|type-/.test(c))).toEqual([]);
	});

	it('rings a blocked card on its tint, and keeps the written flag', () => {
		component = render({ ...base, nature: 'research', blocked: true });
		const card = document.querySelector('a')!;
		expect(card.className).toContain('ring-danger');
		expect(card.className).toContain('bg-nature-research');
		expect(card.textContent).toContain('BLOCKED');
	});

	it('does not ring a card that is not blocked', () => {
		component = render(base);
		expect(document.querySelector('a')!.className).not.toContain('ring-danger');
	});

	it('shows a "waiting to publish" flag when told the card is merged but unreleased', () => {
		component = render(base, { waiting: true });
		expect(document.querySelector('[data-testid="waiting-to-publish"]')?.textContent).toBe(
			'waiting to publish'
		);
	});

	it('shows nothing of the sort by default', () => {
		component = render(base);
		expect(document.querySelector('[data-testid="waiting-to-publish"]')).toBeNull();
	});

	it('fades and dashes the card being dragged, on its tint', () => {
		component = render({ ...base, nature: 'feature' }, { dragging: true });
		const classes = document.querySelector('a')!.className.split(/\s+/);
		expect(classes).toEqual(expect.arrayContaining(['opacity-50', 'border-dashed']));
		expect(classes).toContain('bg-nature-feature');
	});

	it('hands key presses to the board and carries its ID for the board to find it again', () => {
		const onkeydown = vi.fn();
		component = render(base, { onkeydown });
		const card = document.querySelector('a')!;
		expect(card.dataset.id).toBe('S-0048');
		card.dispatchEvent(
			new KeyboardEvent('keydown', { key: 'ArrowUp', altKey: true, bubbles: true })
		);
		expect(onkeydown).toHaveBeenCalledOnce();
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

	// S-0104: what the story's agent is doing, as a dot beside the ID
	it('shows a green, yellow, or red dot for the story’s agent, and none once it finished', () => {
		const run = {
			story: 'S-0048',
			harness: 'claude-code',
			model: 'claude-haiku-4-5',
			command: 'claude',
			agent: 'agent-S-0048',
			started: '2026-09-23T18:00:00Z'
		};
		const dot = () => document.querySelector<HTMLElement>('[data-testid="agent-dot"]');
		component = render(base);
		expect(dot()).toBeNull();
		for (const [state, colour] of [
			['working', 'bg-good'],
			['waiting', 'bg-warn'],
			['failed', 'bg-danger']
		] as const) {
			unmount(component);
			component = render(base, { activity: { state, why: 'why so', run } });
			expect(dot()!.dataset.state).toBe(state);
			expect(dot()!.className).toContain(colour);
			expect(dot()!.getAttribute('aria-label')).toContain(
				`agent ${state} (claude-code, claude-haiku-4-5)`
			);
		}
		unmount(component);
		component = render(base, { activity: { state: 'worked', run } });
		expect(dot()).toBeNull();
	});
});
