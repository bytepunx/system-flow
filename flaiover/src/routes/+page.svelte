<script lang="ts">
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	type Manifest = { name: string; description?: string; template?: { version?: string } };
	type Item = { id: string; type: string; status: string; title: string; archived: boolean };

	let manifest = $state<Manifest | null>(null);
	let items = $state<Item[]>([]);
	let error = $state<string | null>(null);

	const statuses = ['backlog', 'ready', 'in-progress', 'review', 'done'];
	const active = $derived(items.filter((i) => !i.archived));
	const count = (type: string, status: string) =>
		active.filter((i) => i.type === type && i.status === status).length;

	onMount(async () => {
		try {
			const [m, i] = await Promise.all([api('/api/manifest'), api('/api/items')]);
			if (!m.ok) throw new Error((await m.json()).error ?? m.statusText);
			manifest = await m.json();
			items = await i.json();
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	});
</script>

{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">
		Cannot read the project: {error}
	</p>
{:else if manifest}
	<h1 class="text-2xl font-semibold">{manifest.name}</h1>
	{#if manifest.description}<p class="mt-1 text-ink-soft">
			{manifest.description}
		</p>{/if}
	<p class="mt-1 text-xs text-muted">
		template {manifest.template?.version ?? '?'} · {active.length} active items · {items.length -
			active.length} archived
	</p>

	<div class="mt-6 overflow-x-auto">
		<table class="min-w-full text-sm">
			<thead>
				<tr class="text-left text-muted">
					<th class="py-2 pr-4"></th>
					{#each statuses as s (s)}<th class="py-2 pr-4 font-medium">{s}</th>{/each}
				</tr>
			</thead>
			<tbody>
				{#each ['epic', 'story', 'task'] as t (t)}
					<tr class="border-t border-line">
						<td class="py-2 pr-4 font-medium">{t}s</td>
						{#each statuses as s (s)}<td class="py-2 pr-4 tabular-nums">{count(t, s)}</td>{/each}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
