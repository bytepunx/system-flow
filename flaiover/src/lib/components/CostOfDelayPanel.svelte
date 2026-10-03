<script lang="ts">
	// An epic's or a story's cost of delay inputs (S-0204): revenue and penalty per week in the
	// project's currency, and the time lost per cycle. Each is optional, so the panel starts closed and,
	// while closed, says in one line what is set. flai works the value out from them and validates them.
	import { costSummary, type CostFields } from '$lib/costofdelay';

	let {
		fields = $bindable(),
		currency
	}: {
		fields: CostFields;
		/** The ISO 4217 code amounts are in, from the manifest's planning.currency. */
		currency: string;
	} = $props();

	let open = $state(false);
	const summary = $derived(costSummary(fields, currency));
</script>

<details class="rounded border border-line p-3 text-sm" bind:open data-testid="cod-panel">
	<summary class="cursor-pointer select-none">
		<span class="text-xs text-muted">Cost of delay</span>
		{#if !open && summary}<span class="ml-2 text-xs" data-testid="cod-summary">{summary}</span>{/if}
	</summary>
	<div class="mt-2 grid grid-cols-1 gap-3 sm:grid-cols-3">
		<label class="block">
			<span class="mb-1 block text-xs text-muted">revenue per week, {currency}</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={fields.revenue_per_week}
				inputmode="decimal"
				placeholder="1200"
				data-testid="cod-revenue"
			/>
		</label>
		<label class="block">
			<span class="mb-1 block text-xs text-muted">penalty per week, {currency}</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={fields.penalty_per_week}
				inputmode="decimal"
				placeholder="300"
				data-testid="cod-penalty"
			/>
		</label>
		<label class="block">
			<span class="mb-1 block text-xs text-muted">time lost per cycle, a duration</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
				bind:value={fields.time_lost_per_cycle}
				placeholder="6h or 1h30m"
				data-testid="cod-time-lost"
			/>
		</label>
	</div>
	<p class="mt-2 text-xs text-muted">
		Each is optional. The planner works out what a week of waiting costs from them.
	</p>
</details>
