<!--
	The component catalogue (S-0439): each sub-project's components by folder, the directory part of
	each one's path, with one filter on name and path. A sub-project flai does not list components
	for is shown with its kind; a component whose props could not be read is still listed and picked.
-->
<script lang="ts">
	import type { Component, ComponentsProject } from '$lib/server/repo';

	let {
		projects,
		picked,
		onpick
	}: {
		/** The catalogue's sub-projects, as components.list answers them. */
		projects: readonly ComponentsProject[];
		/** The component framed: its sub-project's manifest name and its path. */
		picked: { project: string; path: string } | null;
		/** Called with the sub-project's manifest name and the component's path. */
		onpick: (project: string, path: string) => void;
	} = $props();

	type Folder = { dir: string; components: Component[] };

	let filter = $state('');
	const needle = $derived(filter.trim().toLowerCase());

	const folderOf = (path: string) => {
		const i = path.lastIndexOf('/');
		return i < 0 ? '' : path.slice(0, i);
	};

	function hint(c: Component) {
		return c.error ? `${c.path}: its props could not be read: ${c.error}` : c.path;
	}

	const matches = (c: Component) =>
		!needle || c.name.toLowerCase().includes(needle) || c.path.toLowerCase().includes(needle);

	/** A sub-project's matching components by folder, the folders sorted by path. */
	function foldersOf(p: ComponentsProject): Folder[] {
		const folders: Folder[] = [];
		for (const c of p.components.filter(matches)) {
			const dir = folderOf(c.path);
			const folder = folders.find((f) => f.dir === dir);
			if (folder) folder.components.push(c);
			else folders.push({ dir, components: [c] });
		}
		return folders.sort((a, b) => a.dir.localeCompare(b.dir));
	}

	// Under a filter only the sub-projects with a match are shown; without one, every sub-project.
	const shown = $derived(
		projects
			.map((p) => ({ project: p, folders: foldersOf(p) }))
			.filter((s) => !needle || s.folders.length > 0)
	);
</script>

<div data-testid="component-list">
	<input
		type="search"
		class="w-full rounded border border-line bg-ground px-2 py-1 text-sm"
		placeholder="Filter by name or path"
		aria-label="Filter components by name or path"
		bind:value={filter}
	/>
	{#if needle && !shown.length}
		<p class="mt-3 text-sm text-muted" data-testid="component-list-none">
			No component matches “{filter.trim()}”.
		</p>
	{/if}
	{#each shown as { project, folders } (project.name)}
		<section
			class="mt-4"
			aria-labelledby="components-of-{project.name}"
			data-testid="component-project-{project.name}"
		>
			<h2 id="components-of-{project.name}" class="text-sm font-semibold">
				{project.name} <span class="font-normal text-muted">{project.kind}</span>
			</h2>
			{#if !project.supported}
				<p class="mt-1 text-xs text-muted" data-testid="component-project-unsupported">
					Not Svelte: flai lists the components of SvelteKit and Svelte sub-projects only.
				</p>
			{:else if !project.components.length}
				<p class="mt-1 text-xs text-muted">No components.</p>
			{/if}
			{#each folders as folder (folder.dir)}
				<div class="mt-2" data-testid="component-folder">
					<h3 class="font-mono text-xs text-muted">{folder.dir || '.'}</h3>
					<ul class="mt-1">
						{#each folder.components as c (c.path)}
							{@const current = picked?.project === project.name && picked.path === c.path}
							<li>
								<button
									type="button"
									class={[
										'w-full rounded px-2 py-0.5 text-left text-sm',
										current ? 'bg-primary-soft font-medium' : 'hover:bg-raised'
									]}
									aria-current={current ? 'true' : undefined}
									title={hint(c)}
									data-testid="component-{c.path}"
									onclick={() => onpick(project.name, c.path)}
								>
									{c.name}
									{#if c.error}<span class="text-xs text-warn">props unreadable</span>{/if}
								</button>
							</li>
						{/each}
					</ul>
				</div>
			{/each}
		</section>
	{/each}
</div>
