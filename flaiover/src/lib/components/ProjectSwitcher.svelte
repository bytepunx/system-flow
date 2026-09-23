<script lang="ts">
	// Which project this tab looks at (S-0080, S-0095), always a control (S-0102): the projects to
	// switch between, and, in a group of their own, the repositories flai serve offers for import
	// (S-0098), where picking one asks whether to import it. Picking switches every screen in place:
	// the layout keys the page on the current project, so the page remounts, asks again, and reopens
	// its event stream for the project picked.
	import { projectState } from '$lib/project.svelte';

	const projects = $derived(projectState.list.filter((p) => !p.candidate));
	const offered = $derived(projectState.list.filter((p) => p.candidate));

	function onchange(e: Event) {
		projectState.pick((e.target as HTMLSelectElement).value || null);
	}
	// A project that connected since the list was last asked for shows up when the operator reaches
	// for the list, not only on the next periodic look.
	function onfocus() {
		void projectState.refresh();
	}
</script>

<label class="flex items-center gap-1 text-sm" data-testid="project-switcher">
	<span class="sr-only">Project</span>
	<select
		class="rounded border border-line-strong bg-surface px-2 py-1 text-sm"
		value={projectState.current ?? ''}
		disabled={projectState.list.length === 0}
		title={projectState.list.length === 0
			? 'No project is connected yet: flai dashboard or flai serve on the host connects them'
			: 'Switch project, or pick a repository to import it'}
		{onchange}
		{onfocus}
	>
		{#if projectState.list.length === 0}
			<option value="">{projectState.loaded ? 'no project connected' : 'projects…'}</option>
		{:else if !projectState.current}
			<option value="" disabled>choose a project…</option>
		{/if}
		{#each projects as p (p.key)}
			<option value={p.key}>{p.name}{p.connected ? '' : ' (not connected)'}</option>
		{/each}
		{#if offered.length}
			<optgroup label="Not imported yet">
				{#each offered as p (p.key)}
					<option value={p.key}>{p.name} (not imported)</option>
				{/each}
			</optgroup>
		{/if}
	</select>
</label>
