<script lang="ts">
	import { api } from '$lib/api';
	import { projectState } from '$lib/project.svelte';
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
		refines?: string[];
	};
	let adrs = $state<Adr[]>([]);
	// Recording and accepting are writes through flai (S-0060): offered only when flai answers.
	let writable = $state(false);
	let notice = $state<{ kind: 'error' | 'ok'; text: string } | null>(null);
	let busy = $state<string | null>(null);

	async function load() {
		adrs = await (await api('/api/docs/adrs')).json();
	}
	onMount(() => {
		void load();
		void api('/api/adrs/template')
			.then((r) => (writable = r.ok))
			.catch(() => (writable = false));
		// the list follows the files: an ADR recorded here or from a shell appears without a reload
		const es = new EventSource(projectState.tag('/api/events'));
		es.addEventListener('change', () => void load());
		return () => es.close();
	});

	async function accept(a: Adr) {
		if (
			!confirm(
				`Accept ${a.id}? An accepted ADR is immutable: to change it later you record a new one that supersedes it.`
			)
		)
			return;
		busy = a.id;
		notice = null;
		try {
			const r = await api(`/api/adrs/${a.id}/accept`, { method: 'POST' });
			const body = await r.json();
			notice = r.ok
				? { kind: 'ok', text: `${a.id} accepted` }
				: { kind: 'error', text: body.error ?? r.statusText };
		} finally {
			busy = null;
			await load();
		}
	}
	const byId = $derived(new Map(adrs.map((a) => [a.id, a])));
	const link = (id: string) => byId.get(id)?.path;
</script>

<svelte:head><title>ADRs · flaiover</title></svelte:head>
<div class="flex flex-wrap items-center gap-4">
	<h1 class="text-2xl font-semibold">Architecture decisions</h1>
	{#if writable}
		<a
			class="rounded border border-line-strong bg-surface px-2 py-1 text-sm hover:bg-raised"
			href={resolve('/adrs/new')}
			data-testid="new-adr-link">+ new ADR</a
		>
	{/if}
</div>
{#if notice}
	<p
		class="mt-2 rounded border p-2 text-sm {notice.kind === 'error'
			? 'border-danger bg-danger-soft text-danger'
			: 'border-good bg-good-soft text-good'}"
		role="status"
	>
		{notice.text}
	</p>
{/if}
<p class="mt-1 text-sm text-muted">
	{adrs.length} records. Superseded decisions point at their successors.
</p>
<table class="mt-4 min-w-full text-sm">
	<thead
		><tr class="text-left text-muted"
			><th class="py-2 pr-4">ADR</th><th class="py-2 pr-4">Title</th><th class="py-2 pr-4"
				>Status</th
			><th class="py-2 pr-4">Date</th><th class="py-2">Chain</th></tr
		></thead
	>
	<tbody>
		{#each adrs as a (a.id)}
			<tr class="border-t border-line align-top">
				<td class="py-2 pr-4 font-mono whitespace-nowrap"
					><a class="underline" href={resolve('/docs/[...path]', { path: a.path })}>{a.id}</a></td
				>
				<td class="py-2 pr-4">{a.title}</td>
				<td class="py-2 pr-4 whitespace-nowrap"
					>{a.status}
					{#if writable && a.status === 'proposed'}
						<button
							class="ml-2 rounded border border-line-strong bg-surface px-1.5 text-xs hover:bg-raised disabled:opacity-50"
							disabled={busy === a.id}
							onclick={() => accept(a)}
							data-testid="accept-{a.id}">accept</button
						>
					{/if}</td
				>
				<td class="py-2 pr-4 whitespace-nowrap">{a.date}</td>
				<td class="py-2 text-xs text-muted">
					{#each a.supersedes as s (s)}supersedes {#if link(s)}<a
								class="underline"
								href={resolve('/docs/[...path]', { path: link(s)! })}>{s}</a
							>{:else}{s}{/if}
					{/each}
					{#each a.refines ?? [] as s (s)}refines {#if link(s)}<a
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
