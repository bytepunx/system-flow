<script lang="ts">
	// /settings/<page>: one Configuration page of the host's settings for this project (S-0105,
	// S-0459), its sections changed through flai once the operator has turned on the settings host
	// action in a shell. Going to another page keeps the panel and what it read, and shows that
	// page's sections.
	import SettingsPanel from '$lib/components/SettingsPanel.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head><title>flaiover — {data.title}</title></svelte:head>

<h1 class="mb-4 text-lg font-semibold text-ink">{data.title}</h1>
<SettingsPanel page={data.slug} />
