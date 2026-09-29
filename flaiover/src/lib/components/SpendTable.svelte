<script lang="ts">
	// The table view of a chart of spend over time (S-0163): what each bucket spent, and what
	// that comes to per item, per minute of agent work, and per dollar.
	import { bucketLabel, type BucketSize, type SpendRow } from '$lib/viz/charts';
	import { count, dollars } from '$lib/usage';

	let { rows, bucket }: { rows: SpendRow[]; bucket: BucketSize } = $props();
	const say = (v: number | undefined, as: (n: number) => string) => (v === undefined ? '-' : as(v));
	const minutes = (seconds: number) => (seconds / 60).toFixed(1);
</script>

<table class="min-w-full" data-testid="spend-table">
	<thead
		><tr class="text-left text-muted"
			><th class="pr-4">{bucket}</th><th class="pr-4">of</th><th class="pr-4">items done</th><th
				class="pr-4">tokens</th
			><th class="pr-4">cost</th><th class="pr-4">agent minutes</th><th class="pr-4"
				>tokens per item</th
			><th class="pr-4">cost per item</th><th class="pr-4">tokens per agent minute</th><th
				class="pr-4">tokens per dollar</th
			><th class="pr-4">mean tokens per {bucket}</th><th>mean cost per {bucket}</th></tr
		></thead
	><tbody
		>{#each rows as row (row.at + row.of)}<tr
				><td class="pr-4 font-mono whitespace-nowrap">{bucketLabel(row.at, bucket)}</td><td
					class="pr-4">{row.of}</td
				><td class="pr-4">{row.items}</td><td class="pr-4">{count(row.tokens)}</td><td class="pr-4"
					>{dollars(row.cost)}{row.estimated ? '*' : ''}</td
				><td class="pr-4">{minutes(row.seconds)}</td><td class="pr-4"
					>{say(row.tokens_per_item, count)}</td
				><td class="pr-4">{say(row.cost_per_item, dollars)}</td><td class="pr-4"
					>{say(row.tokens_per_minute, count)}</td
				><td class="pr-4">{say(row.tokens_per_dollar, count)}</td><td class="pr-4"
					>{say(row.mean_tokens, count)}</td
				><td>{say(row.mean_cost, dollars)}</td></tr
			>{/each}</tbody
	>
</table>
