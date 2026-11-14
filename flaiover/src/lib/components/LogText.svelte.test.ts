// S-0785: LogText renders free text from flai with each log path in it as a LogLink to the log
// page, in the code style, keeping the project the tab views, and the rest as plain text.
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { projectState, resetForTests } from '$lib/project.svelte';
import LogText from './LogText.svelte';

let target: HTMLElement;
let component: ReturnType<typeof mount> | undefined;

function show(text: string) {
	component = mount(LogText, { target, props: { text } });
	flushSync();
}
const links = () => [...target.querySelectorAll<HTMLAnchorElement>('[data-testid="log-link"]')];

beforeEach(() => {
	target = document.createElement('div');
	document.body.appendChild(target);
});
afterEach(() => {
	if (component) unmount(component);
	component = undefined;
	target.remove();
	resetForTests();
});

describe('LogText', () => {
	it('links each log path to the log page and keeps the rest as text', () => {
		const text =
			'planner started (pid 42); its output is in ~/.flai/serve/agents/planner.log on the host, and the build in `.flai-cache/kluster/logs/s-0595/build-api.log`.';
		show(text);
		expect(target.textContent).toBe(text);
		const [planner, build] = links();
		expect(links()).toHaveLength(2);
		expect(planner.textContent).toBe('~/.flai/serve/agents/planner.log');
		expect(planner.querySelector('code')!.className).toContain('font-mono');
		const href = new URL(planner.getAttribute('href')!, 'http://x');
		expect(href.pathname).toBe('/ops/log');
		expect(href.searchParams.get('path')).toBe('~/.flai/serve/agents/planner.log');
		expect(href.searchParams.has('project')).toBe(false);
		expect(new URL(build.href).searchParams.get('path')).toBe(
			'.flai-cache/kluster/logs/s-0595/build-api.log'
		);
	});

	it('keeps the project the tab views in each link', () => {
		projectState.pick('acme', false);
		show('see /tmp/a.log.');
		const href = new URL(links()[0].getAttribute('href')!, 'http://x');
		expect(href.searchParams.get('path')).toBe('/tmp/a.log');
		expect(href.searchParams.get('project')).toBe('acme');
		expect(target.textContent).toBe('see /tmp/a.log.');
	});

	it('renders text without a log path as it is', () => {
		show('the catalog is in /usr/share/catalog and the logs in /tmp/run.logs');
		expect(links()).toEqual([]);
		expect(target.textContent).toBe(
			'the catalog is in /usr/share/catalog and the logs in /tmp/run.logs'
		);
	});
});
