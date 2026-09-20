<script lang="ts">
	// Which project this tab looks at (S-0080). Shows nothing while at most one project is
	// connected, so a dashboard serving only one looks exactly as it always has; the moment a
	// second one appears, this is how the designer tells the pages apart.
	import { projectState } from '$lib/project.svelte';

	function onchange(e: Event) {
		projectState.pick((e.target as HTMLSelectElement).value || null);
		// Every page loads its own data once, in onMount; a picked project is a fresh look at
		// everything on the page, so a reload is simpler and more robust than wiring each page to
		// react to the choice on its own.
		if (typeof location !== 'undefined') location.reload();
	}
</script>

{#if projectState.needsChoice}
	<label class="flex items-center gap-1 text-sm" data-testid="project-switcher">
		<span class="sr-only">Project</span>
		<select
			class="rounded border border-line-strong bg-surface px-2 py-1 text-sm"
			value={projectState.current ?? ''}
			{onchange}
		>
			<option value="" disabled>choose a project…</option>
			{#each projectState.list as p (p.key)}
				<option value={p.key}>{p.name}{p.connected ? '' : ' (not connected)'}</option>
			{/each}
		</select>
	</label>
{/if}
