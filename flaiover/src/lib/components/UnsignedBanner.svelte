<script lang="ts">
	// A connection let through unsigned is never silent (S-0239): while the picked project's channel
	// was allowed with a side that has no release stamp, every page names that side. Not dismissible,
	// since the allowance is the operator's on the host and only they can take it away.
	import { projectState, unsignedLabel } from '$lib/project.svelte';

	const sides = $derived(
		projectState.project?.connected ? projectState.project.unsigned : undefined
	);
	const label = $derived(unsignedLabel(sides));
</script>

{#if sides && label}
	<div
		class="border-b border-warn bg-warn-soft px-4 py-3 text-sm"
		role="status"
		data-unsigned-banner
	>
		<div class="mx-auto max-w-6xl">
			<p class="font-medium">{label}</p>
			{#if sides.includes('flai')}
				<p class="mt-1">The flai on the host is a development build, with no release stamp.</p>
			{/if}
			{#if sides.includes('dashboard')}
				<p class="mt-1">This dashboard is a development build, with no release stamp.</p>
			{/if}
			<p class="mt-1">
				The host lets it through on purpose, with <code>dashboard.allow_unsigned</code>; a release
				of flai with the published dashboard image needs no allowance and shows no banner.
			</p>
		</div>
	</div>
{/if}
