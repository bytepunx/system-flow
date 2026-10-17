// S-0440: the event log lists the calls of the framed component's function props, newest first,
// each with its time, the prop's name, and its arguments as JSON, and keeps the last 50.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import EventLog, { MAX_EVENTS, logged, type Logged } from './EventLog.svelte';

const AT = new Date(2026, 9, 16, 9, 5, 7, 42).getTime();

describe('EventLog (S-0440)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	const show = (events: Logged[], onclear?: () => void) => {
		c = mount(EventLog, { target: document.body, props: { events, onclear } });
		flushSync();
	};
	const entries = () => [...document.querySelectorAll('[data-testid="event-log-entry"]')];

	it('says there are no calls yet, and offers no clear', () => {
		show([], () => {});
		expect(document.querySelector('[data-testid="event-log-none"]')).not.toBeNull();
		expect(document.querySelector('button')).toBeNull();
	});

	it('lists each call with its time, its name, and its arguments as JSON', () => {
		show([
			{ at: AT + 1000, name: 'onchange', args: [{ value: 'a' }, 2] },
			{ at: AT, name: 'onclick', args: [{ type: 'click' }] }
		]);
		const [newest, oldest] = entries();
		expect(newest.querySelector('time')!.textContent).toBe('09:05:08.042');
		expect(newest.querySelector('time')!.getAttribute('datetime')).toBe(
			new Date(AT + 1000).toISOString()
		);
		expect(newest.textContent).toContain('onchange');
		expect(newest.querySelector('code')!.textContent).toBe('[{"value":"a"},2]');
		expect(oldest.querySelector('code')!.textContent).toBe('[{"type":"click"}]');
	});

	it('keeps the last 50, newest first', () => {
		let events: Logged[] = [];
		for (let i = 0; i < 60; i++) events = logged(events, { at: AT + i, name: `on${i}`, args: [i] });
		expect(MAX_EVENTS).toBe(50);
		expect(events).toHaveLength(50);
		expect(events[0].name).toBe('on59');
		expect(events.at(-1)!.name).toBe('on10');
		// and shows no more than 50 of a longer list
		show([...events, ...events]);
		expect(entries()).toHaveLength(50);
		expect(entries()[0].textContent).toContain('on59');
	});

	it('clears the log through its parent', () => {
		const onclear = vi.fn();
		show([{ at: AT, name: 'onclick', args: [] }], onclear);
		document.querySelector('button')!.click();
		expect(onclear).toHaveBeenCalledOnce();
	});
});
