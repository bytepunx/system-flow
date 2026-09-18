// The chart instance must follow its props after mount: switching chart
// type, filters, theme, and live updates all arrive as a new `option`
// (S-0045, I-0014).
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Chart from './Chart.svelte';
import { theme } from '$lib/viz/palette';

const setOption = vi.fn();
const resize = vi.fn();
const dispose = vi.fn();
vi.mock('echarts/core', () => ({
	use: vi.fn(),
	init: vi.fn(() => ({ setOption, resize, dispose }))
}));
vi.mock('echarts/charts', () => ({ ScatterChart: {}, LineChart: {}, BarChart: {} }));
vi.mock('echarts/components', () => ({
	GridComponent: {},
	TooltipComponent: {},
	LegendComponent: {},
	MarkLineComponent: {}
}));
vi.mock('echarts/renderers', () => ({ CanvasRenderer: {} }));

async function settle() {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
}

describe('Chart', () => {
	beforeEach(() => {
		setOption.mockClear();
		globalThis.ResizeObserver = class {
			observe() {}
			disconnect() {}
			unobserve() {}
		} as unknown as typeof ResizeObserver;
	});
	afterEach(() => vi.clearAllMocks());

	it('draws the first option and every option that follows', async () => {
		const first = { series: [{ name: 'cycle-time' }] };
		const second = { series: [{ name: 'throughput' }] };
		const props = $state({
			option: first as Record<string, unknown>,
			theme: theme(false),
			height: 360
		});
		const target = document.createElement('div');
		document.body.appendChild(target);
		const component = mount(Chart, { target, props });
		await settle();
		expect(setOption).toHaveBeenCalled();
		expect(setOption.mock.calls.at(-1)?.[0]).toEqual(first);

		props.option = second;
		await settle();
		expect(setOption.mock.calls.at(-1)?.[0]).toEqual(second);

		props.option = first;
		props.theme = theme(true);
		await settle();
		expect(setOption.mock.calls.at(-1)?.[0]).toEqual(first);
		// jsdom normalises the hex surface (#272220) to rgb().
		expect(target.querySelector('div')?.getAttribute('style')).toContain('rgb(39, 34, 32)');

		unmount(component);
		expect(dispose).toHaveBeenCalled();
	});
});
