<script lang="ts">
	// The project is read through flai on the host (S-0073). When none is connected the pages have
	// nothing to show, so every page says what is missing and the command that brings it back.
	import { hostFlai } from '$lib/hostflai.svelte';
	import { projectState } from '$lib/project.svelte';

	const status = $derived(hostFlai.status);
	// Several projects and none chosen yet (S-0095): the status was asked for no project in
	// particular, so "not connected" would be false; the switcher and the project list ask for a choice.
	// A repository offered for import has no project to show, connected or not (S-0098): its own
	// page says what there is to say.
	const choosing = $derived(
		(projectState.needsChoice && !projectState.current) || !!projectState.project?.candidate
	);
</script>

{#if status && !hostFlai.usable && !choosing}
	<div
		class="border-b border-warn bg-warn-soft px-4 py-3 text-sm"
		role="alert"
		data-host-flai-banner
	>
		<div class="mx-auto max-w-6xl">
			{#if !status.configured}
				<p class="font-medium">This dashboard was started by a flai from before it needed one.</p>
				<p class="mt-1">
					The project is read through flai on the host, and this container was given no way to meet
					it. On the host, upgrade flai, then run <code>flai dashboard stop</code> and
					<code>flai dashboard</code> in the project.
				</p>
			{:else if status.connected}
				<p class="font-medium">flai on the host is connected, but the project cannot be shown.</p>
				<p class="mt-1">{status.error}</p>
			{:else}
				<p class="font-medium">No flai on the host is connected, so the project cannot be shown.</p>
				<p class="mt-1">
					On the host, run <code>flai dashboard</code> in the project, or
					<code>flai serve start</code>; <code>flai serve status</code> says why it is not connected.
					This page recovers by itself when it is.
				</p>
			{/if}
		</div>
	</div>
{/if}
