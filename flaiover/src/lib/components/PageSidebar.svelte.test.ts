// S-0460: a page's sidebar of its sections. Its entries are links in a named nav, with counts and
// nested entries; a pane entry shows its pane only; the choice is kept in the hash both ways; an
// unknown hash falls back to the page's default; a scroll entry scrolls to its section and is
// current where nothing reports what is in view; the page can choose and read the entry current.
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { createRawSnippet, flushSync, mount, unmount } from 'svelte';
import type { SidebarEntry } from '$lib/pagesidebar.svelte';
import PageSidebar from './PageSidebar.svelte';

const PANES: SidebarEntry[] = [
	{ id: 'story', label: 'Story', mode: 'pane' },
	{ id: 'tasks', label: 'Tasks', mode: 'pane' },
	{ id: 'questions', label: 'Questions and threads', mode: 'pane', count: 2 },
	{ id: 'messages', label: 'Messages', mode: 'pane', count: 0 }
];
const SECTIONS: SidebarEntry[] = [
	{ id: 'goal', label: 'Goal', mode: 'scroll' },
	{ id: 'criteria', label: 'Criteria', mode: 'scroll' }
];

// The page's pane: one element naming the pane it is given, as the page's {#if} would draw.
const pane = createRawSnippet((id: () => string) => ({
	render: () => '<p data-testid="pane"></p>',
	setup: (el: Element) => {
		$effect(() => {
			el.textContent = `the ${id()} pane`;
		});
	}
}));
const sections = createRawSnippet(() => ({
	render: () =>
		'<div><section id="goal">Goal</section><section id="criteria">Criteria</section></div>'
}));

const entry = (id: string) =>
	document.querySelector(`[data-sidebar-entry="${id}"]`) as HTMLAnchorElement;
const shown = () =>
	[...document.querySelectorAll('[data-testid="pane"]')].map((p) => p.textContent);
const currentIds = () =>
	[...document.querySelectorAll('[aria-current]')].map((a) => a.getAttribute('data-sidebar-entry'));
const settle = async () => {
	for (let i = 0; i < 4; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const scrollIntoView = Element.prototype.scrollIntoView;

describe('PageSidebar (S-0460)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => history.replaceState(null, '', '/items/S-0001'));
	// jsdom follows a link's hash, and back and forward, a task later: let them land before the
	// next test sets its own.
	afterEach(async () => {
		await settle();
		if (c) unmount(c);
		c = undefined;
		Element.prototype.scrollIntoView = scrollIntoView;
		history.replaceState(null, '', '/');
		document.body.innerHTML = '';
	});
	const show = (props: Record<string, unknown>) => {
		c = mount(PageSidebar, {
			target: document.body,
			props: { label: 'Story sections', entries: PANES, pane, ...props } as never
		});
		flushSync();
	};

	it('lists its entries as links in a nav the page names, with their counts', () => {
		show({});
		const nav = document.querySelector('nav')!;
		expect(nav.getAttribute('aria-label')).toBe('Story sections');
		const links = [...nav.querySelectorAll('a')];
		expect(links.map((a) => a.getAttribute('href'))).toEqual([
			'#story',
			'#tasks',
			'#questions',
			'#messages'
		]);
		expect(links.map((a) => a.querySelector('span')!.textContent)).toEqual([
			'Story',
			'Tasks',
			'Questions and threads',
			'Messages'
		]);
		expect(entry('questions').querySelector('[data-sidebar-count]')?.textContent).toBe('2');
		expect(entry('messages').querySelector('[data-sidebar-count]')?.textContent).toBe('0');
		expect(entry('story').querySelector('[data-sidebar-count]')).toBeNull();
	});

	it('lists an entry’s children under it', () => {
		show({
			entries: [
				{
					id: 'story',
					label: 'Story',
					mode: 'pane',
					children: [
						{ id: 'goal', label: 'Goal', mode: 'scroll' },
						{ id: 'criteria', label: 'Criteria', mode: 'scroll' }
					]
				},
				{ id: 'history', label: 'History', mode: 'pane' }
			]
		});
		const nested = entry('story').parentElement!.querySelector('[data-sidebar-list="nested"]')!;
		expect([...nested.querySelectorAll('a')].map((a) => a.getAttribute('href'))).toEqual([
			'#goal',
			'#criteria'
		]);
	});

	it('shows the first pane on arrival without a hash, and only that pane', () => {
		show({});
		expect(shown()).toEqual(['the story pane']);
		expect(currentIds()).toEqual(['story']);
		expect(entry('story').getAttribute('aria-current')).toBe('true');
	});

	it('shows the page’s default pane on arrival without a hash, or with an unknown one', () => {
		show({ fallback: 'questions' });
		expect(shown()).toEqual(['the questions pane']);
		unmount(c!);
		history.replaceState(null, '', '#nope');
		show({ fallback: 'questions' });
		expect(shown()).toEqual(['the questions pane']);
		expect(currentIds()).toEqual(['questions']);
	});

	it('shows the pane the hash names on arrival', () => {
		history.replaceState(null, '', '#tasks');
		show({ fallback: 'questions' });
		expect(shown()).toEqual(['the tasks pane']);
		expect(currentIds()).toEqual(['tasks']);
	});

	it('shows the pane match names for a hash that is no entry’s', () => {
		history.replaceState(null, '', '#MS-0002');
		show({ match: (h: string) => (h.startsWith('MS-') ? 'messages' : undefined) });
		expect(shown()).toEqual(['the messages pane']);
	});

	it('shows a pane when its entry is chosen, and puts it in the hash', async () => {
		show({});
		entry('tasks').click();
		flushSync();
		expect(shown()).toEqual(['the tasks pane']);
		expect(currentIds()).toEqual(['tasks']);
		await settle();
		expect(location.hash).toBe('#tasks');
		expect(shown()).toEqual(['the tasks pane']);
	});

	it('follows the hash as it changes, and back and forward', async () => {
		show({});
		history.replaceState(null, '', '#questions');
		window.dispatchEvent(new HashChangeEvent('hashchange'));
		flushSync();
		expect(shown()).toEqual(['the questions pane']);
		entry('tasks').click();
		await settle();
		expect(location.hash).toBe('#tasks');
		history.back();
		await settle();
		expect(location.hash).toBe('#questions');
		expect(shown()).toEqual(['the questions pane']);
		history.forward();
		await settle();
		expect(shown()).toEqual(['the tasks pane']);
	});

	it('leaves a click that opens the link elsewhere to the browser', () => {
		show({});
		const a = entry('tasks');
		a.addEventListener('click', (e) => e.preventDefault());
		a.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true, ctrlKey: true }));
		flushSync();
		expect(shown()).toEqual(['the story pane']);
	});

	it('falls back to the default when the pane shown goes away', () => {
		history.replaceState(null, '', '#messages');
		const props = $state({ label: 'Story sections', entries: PANES, pane, fallback: 'story' });
		c = mount(PageSidebar, { target: document.body, props });
		flushSync();
		expect(shown()).toEqual(['the messages pane']);
		props.entries = PANES.filter((e) => e.id !== 'messages');
		flushSync();
		expect(shown()).toEqual(['the story pane']);
		expect(entry('messages')).toBeNull();
	});

	it('lets the page choose the pane and read the entry current', () => {
		const props = $state<Record<string, unknown>>({
			label: 'Story sections',
			entries: PANES,
			pane,
			current: undefined
		});
		c = mount(PageSidebar, { target: document.body, props: props as never });
		flushSync();
		expect(props.current).toBe('story');
		props.current = 'questions';
		flushSync();
		expect(shown()).toEqual(['the questions pane']);
		entry('tasks').click();
		flushSync();
		expect(props.current).toBe('tasks');
	});

	it('is usable by keyboard: each entry is a link that takes focus and opens on Enter', () => {
		show({});
		for (const a of document.querySelectorAll<HTMLAnchorElement>('nav a')) {
			expect(a.tagName).toBe('A');
			expect(a.hasAttribute('href')).toBe(true);
			expect(a.tabIndex).toBe(0);
			a.focus();
			expect(document.activeElement).toBe(a);
		}
		// Enter on a focused link is a click, which the browser dispatches and jsdom does not.
		entry('questions').focus();
		entry('questions').click();
		flushSync();
		expect(shown()).toEqual(['the questions pane']);
	});

	it('renders the sidebar alone without a pane snippet', () => {
		show({ pane: undefined });
		expect(document.querySelectorAll('nav a')).toHaveLength(4);
		expect(document.querySelector('[data-sidebar-pane]')).toBeNull();
	});

	describe('scroll entries', () => {
		it('makes the first entry current where nothing reports what is in view', () => {
			show({ entries: SECTIONS, pane: undefined, children: sections });
			expect(currentIds()).toEqual(['goal']);
			expect(document.querySelector('#criteria')).not.toBeNull();
		});

		it('scrolls to the section chosen, replaces the hash, and makes it current', async () => {
			const scrolled: string[] = [];
			Element.prototype.scrollIntoView = function (this: Element) {
				scrolled.push(this.id);
			};
			show({ entries: SECTIONS, pane: undefined, children: sections });
			const length = history.length;
			entry('criteria').click();
			await settle();
			expect(scrolled).toContain('criteria');
			expect(location.hash).toBe('#criteria');
			expect(history.length).toBe(length);
			expect(currentIds()).toEqual(['criteria']);
		});

		it('scrolls to the section the hash names on arrival', async () => {
			const scrolled: string[] = [];
			Element.prototype.scrollIntoView = function (this: Element) {
				scrolled.push(this.id);
			};
			history.replaceState(null, '', '#criteria');
			show({ entries: SECTIONS, pane: undefined, children: sections });
			await settle();
			expect(scrolled).toEqual(['criteria']);
			expect(currentIds()).toEqual(['criteria']);
		});

		it('shows the pane a nested section sits in when it is chosen', async () => {
			Element.prototype.scrollIntoView = () => {};
			show({
				entries: [
					{ id: 'story', label: 'Story', mode: 'pane' },
					{
						id: 'questions',
						label: 'Questions',
						mode: 'pane',
						children: [{ id: 'open', label: 'Open', mode: 'scroll' }]
					}
				]
			});
			entry('open').click();
			await settle();
			expect(shown()).toEqual(['the questions pane']);
			expect(currentIds()).toEqual(['open']);
		});
	});

	it('is a row that scrolls in itself below md, and a sticky column beside it from md', () => {
		show({});
		const nav = document.querySelector('[data-sidebar-nav]')!.className;
		expect(nav).toContain('min-w-0');
		expect(nav).toContain('md:sticky');
		const list = document.querySelector('[data-sidebar-list="top"]')!.className;
		expect(list).toContain('overflow-x-auto');
		expect(list).toContain('md:flex-col');
		const box = document.querySelector('[data-page-sidebar]')!.className;
		expect(box).toContain('flex-col');
		expect(box).toContain('md:flex-row');
		expect(document.querySelector('[data-sidebar-content]')!.className).toContain('min-w-0');
	});
});
