<script lang="ts">
	// The table view of a planning chart (S-0212): each story done in the window with its forecast,
	// its cycle time, and how far its forecast, delivery, and estimate were from what happened.
	import { resolve } from '$app/paths';
	import { human, humanSigned, type ForecastRow } from '$lib/viz/charts';
	import { localDate } from '$lib/localtime';
	import KindChips from './KindChips.svelte';

	let { rows }: { rows: ForecastRow[] } = $props();
	const say = (v: number | undefined, as: (n: number) => string) => (v === undefined ? '-' : as(v));
</script>

<table class="min-w-full" data-testid="forecast-table">
	<thead
		><tr class="text-left text-muted"
			><th class="pr-4">story</th><th class="pr-4">completed</th><th class="pr-4">nature</th><th
				class="pr-4">model</th
			><th class="pr-4">forecast</th><th class="pr-4">actual</th><th class="pr-4">forecast error</th
			><th class="pr-4">delivery error</th><th>estimate error</th></tr
		></thead
	><tbody
		>{#each rows as row (row.id)}<tr
				><td class="pr-4"
					><a class="font-mono underline" href={resolve('/items/[id]', { id: row.id })}>{row.id}</a>
					{row.title}</td
				><td class="pr-4 font-mono whitespace-nowrap">{localDate(row.completed)}</td><td
					class="py-0.5 pr-4"><KindChips nature={row.nature} /></td
				><td class="pr-4">{row.model}</td><td class="pr-4">{say(row.forecast_seconds, human)}</td
				><td class="pr-4">{say(row.cycle_time_seconds, human)}</td><td class="pr-4"
					>{say(row.forecast_error_seconds, humanSigned)}</td
				><td class="pr-4">{say(row.delivery_error_seconds, humanSigned)}</td><td
					>{say(row.estimate_error_seconds, humanSigned)}</td
				></tr
			>{/each}</tbody
	>
</table>
