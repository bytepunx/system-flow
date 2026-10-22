<script lang="ts">
	// A chart's explainer (S-0485): what the chart shows, how each figure is measured, and how to
	// read it, from viz/explain.ts, with a link to the definitions in metrics.md on the Docs page.
	// Open on the first visit; opening or closing it holds for every chart in this browser.
	import { resolve } from '$app/paths';
	import { explain } from '$lib/viz/explain';
	import type { Kind } from '$lib/viz/catalog';
	import { chartExplainer } from '$lib/chartexplainer.svelte';
	import { projectState } from '$lib/project.svelte';
	import { withProject } from '$lib/sitemenu';

	let {
		kind,
		type,
		bucket
	}: {
		kind: Kind;
		/** The item type the chart reads: story, task, or epic. */
		type: string;
		/** What the chart is laid out in: hour, day, or week. */
		bucket: string;
	} = $props();

	const e = $derived(explain(kind, { type, bucket }));
	const parts = $derived(
		[
			{ key: 'what', text: e.what },
			{ key: 'how', text: e.how },
			{ key: 'read', text: e.read }
		].filter((p) => p.text)
	);
	// the Docs page decodes the hash before it looks the heading up, and a heading's id is already
	// encoded (markdown.ts slug), so the anchor is encoded once more to survive that decoding
	const href = $derived(
		`${withProject(resolve('/docs/[...path]', { path: 'design/system/metrics.md' }), projectState.current)}#${encodeURIComponent(e.metrics)}`
	);
	/** The text cut at its backticks: every second piece names a field, a command, or a file. */
	const pieces = (text: string) => text.split('`');
</script>

<details
	class="text-sm"
	data-testid="explainer"
	open={chartExplainer.open}
	ontoggle={(ev) => {
		if (ev.currentTarget.open !== chartExplainer.open) chartExplainer.set(ev.currentTarget.open);
	}}
>
	<summary class="cursor-pointer text-muted">About this chart</summary>
	<div class="mt-2 flex max-w-3xl flex-col gap-2">
		{#each parts as p (p.key)}
			<p data-testid="explainer-{p.key}">
				{#each pieces(p.text) as s, i (i)}{#if i % 2}<code>{s}</code>{:else}{s}{/if}{/each}
			</p>
		{/each}
		<p class="text-xs text-muted">
			<!-- eslint-disable svelte/no-navigation-without-resolve -- the path is resolve()d; only ?project= and the heading are added -->
			<a class="underline" data-testid="explainer-metrics" {href}
				>The definitions in <code>metrics.md</code></a
			>
			<!-- eslint-enable svelte/no-navigation-without-resolve -->
		</p>
	</div>
</details>
