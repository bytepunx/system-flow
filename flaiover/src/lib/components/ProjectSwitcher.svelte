<script lang="ts">
	// Which project this tab looks at (S-0080, S-0095), always a control (S-0102): the projects to
	// switch between, and, in a group of their own, the repositories flai serve offers for import
	// (S-0098), where picking one asks whether to import it. Picking switches every screen in place:
	// the layout keys the page on the current project, so the page remounts, asks again, and reopens
	// its event stream for the project picked. A project flai serve serves that is not connected is
	// labelled apart from one that is, and a link goes to the settings' Projects page, which says why (S-0122, S-0459).
	// Each project says whether its flai is signed, and one whose flai the dashboard refused as
	// unsigned says so, with the reason in the list of why (S-0237). One whose connection the
	// dashboard let through unsigned under FLAIOVER_ALLOW_UNSIGNED names the side that was (S-0239).
	import { resolve } from '$app/paths';
	import { projectState, unsignedLabel, type Project } from '$lib/project.svelte';
	import { refusalText, trust, type FlaiTrust } from '$lib/settings';

	// /api/projects gives each project what the connection list reads of its flai
	type Listed = Project & FlaiTrust;

	const projects = $derived((projectState.list as Listed[]).filter((p) => !p.candidate));
	const offered = $derived(projectState.list.filter((p) => p.candidate));
	const down = $derived(projects.filter((p) => !p.connected));

	function label(p: Listed): string {
		if (p.connected) {
			// a connection let through unsigned names its sides in place of signed or unsigned (S-0239)
			const said = unsignedLabel(p.unsigned) || trust(p);
			return said ? `${p.name} (${said})` : p.name;
		}
		const said = trust(p);
		return `${p.name} (${said || (p.served ? 'served, not connected' : 'not connected')})`;
	}

	/** Why a project is not connected: the dashboard's own refusal first, else what flai serve says. */
	function whyDown(p: Listed): string | undefined {
		return p.refused ? refusalText(p.refused) : p.lastError;
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
				<option value={p.key} title={p.connected ? undefined : whyDown(p)}>{label(p)}</option>
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
			href={resolve('/settings/[page]', { page: 'projects' })}
			data-testid="project-switcher-why"
			onclick={why}
			title={down.map((p) => `${p.name}: ${whyDown(p) ?? 'not connected'}`).join('\n')}
			>{down.length} not connected: why?</a
		>
	{/if}
</span>
