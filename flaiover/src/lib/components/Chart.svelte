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
	let chart: import('echarts/core').ECharts | null = null;

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
			chart.setOption(option);
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

	$effect(() => {
		chart?.setOption(option, true);
	});
</script>

<div
	bind:this={el}
	style="height: {height}px; background: {theme.surface}"
	class="w-full rounded border border-zinc-200 dark:border-zinc-800"
></div>
