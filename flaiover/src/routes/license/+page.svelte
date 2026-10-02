<script lang="ts">
	// /license: the license the dashboard is distributed under (S-0231), in the site menu's Host
	// group. The text is the image's own LICENSE.md, the same one flai license prints.
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { render } from '$lib/markdown';

	let name = $state('License');
	let html = $state('');
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const r = await api('/api/license');
			if (!r.ok) throw new Error((await r.json()).error ?? r.statusText);
			const license = (await r.json()) as { name: string; text: string };
			name = license.name;
			html = render(license.text, 'LICENSE.md');
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	});
</script>

<svelte:head><title>flaiover — license</title></svelte:head>

<h1 class="mb-4 text-lg font-semibold text-ink">License</h1>
<p class="mb-4 text-sm text-muted">
	The terms this dashboard, flai, and the system-flow repository are distributed under. The same
	text is <code>LICENSE.md</code> in the repository and what <code>flai license</code> prints.
</p>
{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">{error}</p>
{:else if html}
	<article class="prose max-w-none rounded border border-line bg-surface p-4" aria-label={name}>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -- the image's own LICENSE.md, rendered client side -->
		{@html html}
	</article>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
