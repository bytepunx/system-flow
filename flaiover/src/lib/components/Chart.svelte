<script lang="ts">
	// ECharts on core with only the pieces the charts use; the instance
	// follows the option and the container size.
	import { onMount } from 'svelte';
	import type { Theme } from '$lib/viz/palette';

	let {
		option,
		theme,
		height = 360
	}: { option: Record<string, unknown>; theme: Theme; height?: number } = $props();
	let el: HTMLDivElement;
	// Reactive (raw: the instance itself is not proxied) so the redraw effect
	// runs again once the asynchronously created instance exists.
	let chart = $state.raw<import('echarts/core').ECharts | null>(null);

	onMount(() => {
		let disposed = false;
		(async () => {
			const echarts = await import('echarts/core');
			const [
				{ ScatterChart, LineChart, BarChart },
				{ GridComponent, TooltipComponent, LegendComponent, MarkLineComponent },
				{ CanvasRenderer }
			] = await Promise.all([
				import('echarts/charts'),
				import('echarts/components'),
				import('echarts/renderers')
			]);
			echarts.use([
				ScatterChart,
				LineChart,
				BarChart,
				GridComponent,
				TooltipComponent,
				LegendComponent,
				MarkLineComponent,
				CanvasRenderer
			]);
			if (disposed) return;
			chart = echarts.init(el, undefined, { renderer: 'canvas' });
		})();
		const ro = new ResizeObserver(() => chart?.resize());
		ro.observe(el);
		return () => {
			disposed = true;
			ro.disconnect();
			chart?.dispose();
			chart = null;
		};
	});

	// Read the props before touching the instance: an optional chain on a
	// null instance would skip them, and the effect would never track them
	// (S-0045). notMerge replaces the previous chart type entirely.
	$effect(() => {
		const next = option;
		chart?.setOption(next, true);
	});
	$effect(() => {
		void height;
		chart?.resize();
	});
</script>

<div
	bind:this={el}
	style="height: {height}px; background: {theme.surface}"
	class="w-full rounded border border-line"
></div>
