<script lang="ts">
	// Which project this tab looks at (S-0080, S-0095), always a control (S-0102): the projects to
	// switch between, and, in a group of their own, the repositories flai serve offers for import
	// (S-0098), where picking one asks whether to import it. Picking switches every screen in place:
	// the layout keys the page on the current project, so the page remounts, asks again, and reopens
	// its event stream for the project picked. A project flai serve serves that is not connected is
	// labelled apart from one that is, and a link goes to the settings page, which says why (S-0122).
	import { resolve } from '$app/paths';
	import { projectState, type Project } from '$lib/project.svelte';

	const projects = $derived(projectState.list.filter((p) => !p.candidate));
	const offered = $derived(projectState.list.filter((p) => p.candidate));
	const down = $derived(projects.filter((p) => !p.connected));

	function label(p: Project): string {
		if (p.connected) return p.name;
		return `${p.name} (${p.served ? 'served, not connected' : 'not connected'})`;
	}

	function onchange(e: Event) {
		projectState.pick((e.target as HTMLSelectElement).value || null);
	}
	// A project that connected since the list was last asked for shows up when the operator reaches
	// for the list, not only on the next periodic look.
	function onfocus() {
		void projectState.refresh();
	}
	// The settings page asks the current project's flai, so from a project that is not connected it
	// goes by way of one that is: the projects flai serve serves are the same from any of them.
	function why() {
		if (projectState.project?.connected) return;
		const via = projects.find((p) => p.connected);
		if (via) projectState.pick(via.key);
	}
</script>

<span class="flex items-center gap-2">
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
				<option value={p.key} title={p.connected ? undefined : p.lastError}>{label(p)}</option>
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
	{#if down.length}
		<a
			class="text-xs text-danger underline"
			href={resolve('/settings')}
			data-testid="project-switcher-why"
			onclick={why}
			title={down.map((p) => `${p.name}: ${p.lastError ?? 'not connected'}`).join('\n')}
			>{down.length} not connected: why?</a
		>
	{/if}
</span>
