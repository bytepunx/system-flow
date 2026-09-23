<script lang="ts">
	// Which project this tab looks at (S-0080, S-0095). Always names the active project once one is
	// known, so the operator can tell at a glance which repository every screen is showing; with more
	// than one it is a choice, and picking another switches every screen in place: the layout keys
	// the page on the current project, so the page remounts, asks again, and reopens its event stream
	// for the project picked.
	import { projectState } from '$lib/project.svelte';

	function onchange(e: Event) {
		projectState.pick((e.target as HTMLSelectElement).value || null);
	}
	// A project that connected since the list was last asked for shows up when the operator reaches
	// for the list, not only on the next periodic look.
	function onfocus() {
		void projectState.refresh();
	}
</script>

{#if projectState.needsChoice}
	<label class="flex items-center gap-1 text-sm" data-testid="project-switcher">
		<span class="sr-only">Project</span>
		<select
			class="rounded border border-line-strong bg-surface px-2 py-1 text-sm"
			value={projectState.current ?? ''}
			{onchange}
			{onfocus}
		>
			<option value="" disabled>choose a project…</option>
			{#each projectState.list as p (p.key)}
				<option value={p.key}
					>{p.name}{p.candidate ? ' (not imported)' : p.connected ? '' : ' (not connected)'}</option
				>
			{/each}
		</select>
	</label>
{:else if projectState.list.length === 1}
	<span class="text-sm font-medium text-ink" data-testid="project-name" title="The project shown"
		>{projectState.list[0].name}</span
	>
{/if}
