<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';

	type Adr = {
		path: string;
		id: string;
		title: string;
		status: string;
		date: string;
		supersedes: string[];
		supersededBy: string[];
	};
	let adrs = $state<Adr[]>([]);
	onMount(async () => {
		adrs = await (await fetch('/api/docs/adrs')).json();
	});
	const byId = $derived(new Map(adrs.map((a) => [a.id, a])));
	const link = (id: string) => byId.get(id)?.path;
</script>

<svelte:head><title>ADRs · flaiover</title></svelte:head>
<h1 class="text-2xl font-semibold">Architecture decisions</h1>
<p class="mt-1 text-sm text-zinc-500">
	{adrs.length} records. Superseded decisions point at their successors.
</p>
<table class="mt-4 min-w-full text-sm">
	<thead
		><tr class="text-left text-zinc-500"
			><th class="py-2 pr-4">ADR</th><th class="py-2 pr-4">Title</th><th class="py-2 pr-4"
				>Status</th
			><th class="py-2 pr-4">Date</th><th class="py-2">Chain</th></tr
		></thead
	>
	<tbody>
		{#each adrs as a (a.id)}
			<tr class="border-t border-zinc-200 align-top dark:border-zinc-800">
				<td class="py-2 pr-4 font-mono whitespace-nowrap"
					><a class="underline" href={resolve('/docs/[...path]', { path: a.path })}>{a.id}</a></td
				>
				<td class="py-2 pr-4">{a.title}</td>
				<td class="py-2 pr-4 whitespace-nowrap">{a.status}</td>
				<td class="py-2 pr-4 whitespace-nowrap">{a.date}</td>
				<td class="py-2 text-xs text-zinc-500">
					{#each a.supersedes as s (s)}supersedes {#if link(s)}<a
								class="underline"
								href={resolve('/docs/[...path]', { path: link(s)! })}>{s}</a
							>{:else}{s}{/if}
					{/each}
					{#each a.supersededBy as s (s)}superseded by {#if link(s)}<a
								class="underline"
								href={resolve('/docs/[...path]', { path: link(s)! })}>{s}</a
							>{:else}{s}{/if}
					{/each}
				</td>
			</tr>
		{/each}
	</tbody>
</table>
