// S-0485: a chart's explainer, a details element open on the first visit that keeps the user's
// choice for every chart, with what the chart shows, how it is measured, how to read it, and a
// link to its definitions in metrics.md on the Docs page.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { projectState, resetForTests } from '$lib/project.svelte';
import { explain } from '$lib/viz/explain';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

import ChartExplainer from './ChartExplainer.svelte';
import { chartExplainer } from '$lib/chartexplainer.svelte';

const KEY = 'flaiover-chart-explainer';

function fakeStorage() {
	const data = new Map<string, string>();
	return {
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => void data.set(k, v),
		removeItem: (k: string) => void data.delete(k),
		clear: () => data.clear()
	};
}

const details = () => document.querySelector('[data-testid="explainer"]') as HTMLDetailsElement;
const part = (key: string) => document.querySelector(`[data-testid="explainer-${key}"]`);
const link = () => document.querySelector('[data-testid="explainer-metrics"]') as HTMLAnchorElement;

describe('ChartExplainer (S-0485)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		vi.stubGlobal('localStorage', fakeStorage());
		resetForTests();
		chartExplainer.open = true;
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		vi.unstubAllGlobals();
	});

	it('is a details element, About this chart, open on the first visit', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'cycle-time', type: 'story', bucket: 'day' }
		});
		flushSync();
		expect(details().tagName).toBe('DETAILS');
		expect(details().open).toBe(true);
		expect(details().querySelector('summary')!.textContent).toBe('About this chart');
	});

	it('draws what, how, and read as paragraphs, a backticked name as code', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'cycle-time', type: 'story', bucket: 'day' }
		});
		flushSync();
		const e = explain('cycle-time', { type: 'story', bucket: 'day' });
		for (const key of ['what', 'how', 'read'] as const) {
			expect(part(key)!.tagName).toBe('P');
			expect(part(key)!.textContent!.trim()).toBe(e[key].replaceAll('`', ''));
		}
		const codes = [...part('how')!.querySelectorAll('code')].map((x) => x.textContent);
		expect(codes).toContain('summary.cycle_time.p50_seconds');
	});

	it('names the type and the bucket it is given', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'tokens-spent', type: 'task', bucket: 'week' }
		});
		flushSync();
		expect(part('how')!.textContent).toContain('Each week holds');
		expect(part('how')!.textContent).toContain('usage.spend.task.buckets[]');
	});

	it('links the heading of metrics.md that defines the chart, on the Docs page, for the project', () => {
		projectState.list = [{ key: 'harbour', name: 'Harbour', connected: true }];
		projectState.pick('harbour', false);
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'agent-waiting', type: 'story', bucket: 'day' }
		});
		flushSync();
		expect(link().getAttribute('href')).toBe(
			'/docs/design/system/metrics.md?project=harbour#waiting'
		);
	});

	it('encodes an anchor once more, as the Docs page decodes the hash before it finds the heading', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'cost', type: 'story', bucket: 'day' }
		});
		flushSync();
		const hash = link().getAttribute('href')!.split('#')[1];
		expect(decodeURIComponent(hash)).toBe(
			explain('cost', { type: 'story', bucket: 'day' }).metrics
		);
	});

	it('keeps the choice when closed, and opens closed on every chart after', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'cycle-time', type: 'story', bucket: 'day' }
		});
		flushSync();
		details().open = false;
		details().dispatchEvent(new Event('toggle'));
		flushSync();
		expect(chartExplainer.open).toBe(false);
		expect(localStorage.getItem(KEY)).toBe('closed');
		unmount(c);
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'throughput', type: 'story', bucket: 'day' }
		});
		flushSync();
		expect(details().open).toBe(false);
		details().open = true;
		details().dispatchEvent(new Event('toggle'));
		flushSync();
		expect(localStorage.getItem(KEY)).toBe('open');
	});

	it('writes nothing when it opens as it was left', () => {
		c = mount(ChartExplainer, {
			target: document.body,
			props: { kind: 'cycle-time', type: 'story', bucket: 'day' }
		});
		flushSync();
		details().dispatchEvent(new Event('toggle'));
		flushSync();
		expect(localStorage.getItem(KEY)).toBeNull();
	});
});
