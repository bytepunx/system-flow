// S-0458: which of the nav rail's groups are expanded, the operator's choices kept per browser and
// a page's group expanded on arrival without being stored.
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { NavRailState, parseChosen } from './navrail.svelte';

const KEY = 'flaiover-nav-rail';

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

describe('parseChosen', () => {
	it('chooses nothing when nothing is stored or what is stored is not a choice', () => {
		for (const raw of [null, '', 'not json', '[]', '"flow"', 'null', '{"flow":"yes"}'])
			expect(parseChosen(raw)).toEqual({});
	});

	it('reads the known groups back and leaves out the rest', () => {
		expect(parseChosen('{"flow":false,"design":true,"workflow":true}')).toEqual({
			flow: false,
			design: true
		});
	});

	// S-0484: whether Charts is expanded is kept beside the groups.
	it('reads Charts back when it is a boolean and leaves it out when it is not', () => {
		expect(parseChosen('{"flow":true,"charts":true}')).toEqual({ flow: true, charts: true });
		expect(parseChosen('{"charts":false}')).toEqual({ charts: false });
		expect(parseChosen('{"flow":false,"charts":"yes"}')).toEqual({ flow: false });
	});

	it('reads a value stored by an older dashboard, with groups only', () => {
		expect(parseChosen('{"agents":true,"configuration":false}')).toEqual({
			agents: true,
			configuration: false
		});
	});
});

describe('NavRailState', () => {
	beforeEach(() => vi.stubGlobal('localStorage', fakeStorage()));

	it('starts with every group collapsed in a fresh browser', () => {
		const rail = new NavRailState();
		for (const g of ['flow', 'agents', 'design', 'documentation', 'configuration'] as const)
			expect(rail.isExpanded(g)).toBe(false);
	});

	it('toggles a group and writes the choice, which a new page reads back', () => {
		const rail = new NavRailState();
		rail.toggle('agents');
		expect(rail.isExpanded('agents')).toBe(true);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ agents: true });
		rail.toggle('agents');
		expect(rail.isExpanded('agents')).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ agents: false });
		rail.toggle('design');
		const next = new NavRailState();
		expect(next.isExpanded('agents')).toBe(false);
		expect(next.isExpanded('design')).toBe(true);
	});

	it('expands a group on arrival without storing it, so a collapsed choice survives', () => {
		const rail = new NavRailState();
		rail.toggle('flow');
		rail.toggle('flow');
		rail.expand('flow');
		rail.expand('configuration');
		expect(rail.isExpanded('flow')).toBe(true);
		expect(rail.isExpanded('configuration')).toBe(true);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ flow: false });
		const next = new NavRailState();
		expect(next.isExpanded('flow')).toBe(false);
		expect(next.isExpanded('configuration')).toBe(false);
	});

	it('collapses a group expanded on arrival with one press', () => {
		const rail = new NavRailState();
		rail.expand('flow');
		rail.toggle('flow');
		expect(rail.isExpanded('flow')).toBe(false);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ flow: false });
	});

	it('toggles Charts and writes the choice beside the groups, which a new page reads back', () => {
		const rail = new NavRailState();
		rail.toggle('design');
		rail.toggle('charts');
		expect(rail.isExpanded('charts')).toBe(true);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ design: true, charts: true });
		rail.toggle('charts');
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ design: true, charts: false });
		rail.toggle('charts');
		const next = new NavRailState();
		expect(next.isExpanded('charts')).toBe(true);
		expect(next.isExpanded('design')).toBe(true);
	});

	it('expands Charts on arrival without storing it', () => {
		const rail = new NavRailState();
		rail.toggle('flow');
		rail.expand('flow');
		rail.expand('charts');
		expect(rail.isExpanded('charts')).toBe(true);
		expect(JSON.parse(localStorage.getItem(KEY)!)).toEqual({ flow: true });
		const next = new NavRailState();
		expect(next.isExpanded('charts')).toBe(false);
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
		const rail = new NavRailState();
		expect(rail.isExpanded('flow')).toBe(false);
		rail.toggle('flow');
		expect(rail.isExpanded('flow')).toBe(true);
	});
});
