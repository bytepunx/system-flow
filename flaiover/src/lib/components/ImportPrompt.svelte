<script lang="ts">
	// A repository the host offers for import, picked in the switcher (S-0098). It is not a
	// system-flow project yet, so there is no board or inbox to show: this asks whether to import it.
	// Declining is not remembered, so picking it again asks again; importing runs flai import --commit
	// on the host and, once it has been applied, switches to the project it became.
	import { onDestroy, onMount } from 'svelte';
	import { api } from '$lib/api';
	import { projectState, type Project } from '$lib/project.svelte';
	import { readNdjson } from '$lib/review';

	let { project }: { project: Project } = $props();

	type TestCommand = { name: string; dir: string; command: string[] };
	type TestResult = TestCommand & {
		ok: boolean;
		exit_code: number;
		seconds: number;
		output: string;
	};
	type Preview = {
		analysis?: { root?: string };
		plan?: { projects?: { name: string; kind: string; path: string }[] };
		tests?: TestCommand[];
		tests_from?: string;
	};
	type Result = {
		key?: string;
		commit?: { tests?: TestResult[]; committed?: boolean; commit?: string; reason?: string };
	};

	let preview = $state<Preview | null>(null);
	let previewError = $state<string | null>(null);
	let declined = $state(false);
	let busy = $state(false);
	let steps = $state<string[]>([]);
	let result = $state<Result | null>(null);
	let committed = $state(false);
	let error = $state<string | null>(null);
	let opening = $state(false);

	onMount(async () => {
		try {
			const r = await api('/api/import');
			const body = await r.json();
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			preview = body as Preview;
		} catch (e) {
			previewError = e instanceof Error ? e.message : String(e);
		}
	});

	// what became of the import stays on the screen after flai serve stops offering the repository
	onDestroy(() => {
		if (projectState.hold?.key === project.key) projectState.hold = null;
	});

	async function importIt() {
		projectState.hold = project;
		busy = true;
		error = null;
		steps = [];
		try {
			const r = await api('/api/import', { method: 'POST' });
			await readNdjson(r, (line) => {
				if (line.event === 'progress') steps = [...steps, String(line.msg)];
				else if (line.event === 'done') {
					result = line.result as Result;
					committed = line.committed === true;
				} else if (line.event === 'error') error = String(line.error);
			});
			if (!result && !error) error = 'the import ended without an answer; look on the host';
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		} finally {
			busy = false;
		}
		if (result && committed) await open();
	}

	// flai serve serves the imported repository as a project within a second or two: wait for its
	// connection, then show it.
	async function open() {
		const key = result?.key;
		if (!key) return;
		opening = true;
		for (let i = 0; i < 30; i++) {
			await projectState.refresh();
			if (projectState.list.some((p) => p.key === key && !p.candidate && p.connected)) {
				projectState.pick(key);
				return;
			}
			await new Promise((r) => setTimeout(r, 500));
		}
		opening = false;
		error = `the project ${key} has not connected yet; pick it in the switcher once it does`;
	}

	const failed = $derived((result?.commit?.tests ?? []).filter((t) => !t.ok));
</script>

<section class="max-w-2xl" data-testid="import-prompt">
	<h1 class="text-2xl font-semibold">{project.name}</h1>
	<p class="mt-2 text-ink-soft">
		{result
			? 'Imported into system-flow'
			: 'This repository is not a system-flow project yet'}{preview?.analysis?.root
			? ` (${preview.analysis.root})`
			: ''}.
	</p>

	{#if result}
		{#if committed}
			<p class="mt-4" data-testid="import-done">
				Imported and committed as <code>{result.commit?.commit}</code>.
				{opening ? 'Opening it…' : ''}
			</p>
		{:else}
			<div
				class="mt-4 rounded border border-warn bg-warn-soft p-3"
				data-testid="import-not-committed"
			>
				<p class="font-medium">Imported, but not committed: a test failed.</p>
				<p class="mt-1 text-sm">
					The standard's files are in the repository, uncommitted. Fix what failed, then commit them
					there. The repository is a project from now on either way.
				</p>
			</div>
			{#each failed as t (t.name)}
				<h2 class="mt-4 font-medium">{t.name} failed (exit {t.exit_code})</h2>
				<pre class="mt-1 max-h-72 overflow-auto rounded bg-surface p-2 text-xs">{t.output}</pre>
			{/each}
			<button
				type="button"
				class="mt-4 rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={opening}
				onclick={open}>{opening ? 'Opening…' : 'Open the project'}</button
			>
		{/if}
	{:else if declined}
		<p class="mt-4" data-testid="import-declined">
			Not imported. There is nothing to show for it until it is: select a different project in the
			switcher above.
		</p>
		<button
			type="button"
			class="mt-3 rounded border border-line-strong px-3 py-1"
			onclick={() => (declined = false)}>Import it after all</button
		>
	{:else}
		<p class="mt-4">Import it into system-flow? flai on the host will:</p>
		<ul class="mt-2 list-disc pl-6 text-sm">
			<li>add the standard's folders and files to it, overwriting nothing that is there;</li>
			{#if previewError}
				<li>run the tests it finds (it could not say which: {previewError});</li>
			{:else if !preview}
				<li>run the tests it finds…</li>
			{:else if (preview.tests ?? []).length === 0}
				<li>run its tests: none were found, so none will run;</li>
			{:else}
				<li>
					run {preview.tests_from === 'host' ? 'the checks this host names' : 'its tests'}:
					{#each preview.tests ?? [] as t, i (t.name)}{i > 0 ? ', ' : ''}<code>{t.name}</code
						>{/each};
				</li>
			{/if}
			<li>commit the import if they pass, and leave it uncommitted if one fails.</li>
		</ul>
		<p class="mt-2 text-sm text-muted">
			The repository must have no uncommitted changes; the commit holds the import and nothing else.
		</p>
		{#if steps.length}
			<ul class="mt-4 text-sm" data-testid="import-steps">
				{#each steps as s, i (i)}<li>{s}</li>{/each}
			</ul>
		{/if}
		<div class="mt-4 flex gap-2">
			<button
				type="button"
				class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={busy}
				onclick={importIt}>{busy ? 'Importing…' : 'Import'}</button
			>
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1"
				disabled={busy}
				onclick={() => (declined = true)}>Not now</button
			>
		</div>
	{/if}
	{#if error}
		<p class="mt-4 text-danger" role="alert">{error}</p>
	{/if}
</section>
