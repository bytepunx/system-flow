<script lang="ts">
	// /design: the project's design system (S-0382), flai's design.get through /api/design: the token
	// file's problems, its tokens laid out, its guidance rendered from the document the answer names,
	// and the tokens previewed on sample components in light and dark. A project with no token file
	// is told so, with the command that writes one. Both are read again when a document changes.
	import type { DesignSystem } from '$lib/server/repo';
	import { api } from '$lib/api';
	import { onMount, tick } from 'svelte';
	import { follow } from '$lib/events';
	import { enhance, render } from '$lib/markdown';
	import { themeState } from '$lib/theme.svelte';
	import TokenSheet from '$lib/components/design/TokenSheet.svelte';
	import PreviewFrame from '$lib/components/design/PreviewFrame.svelte';

	let design = $state<DesignSystem | null>(null);
	let error = $state<string | null>(null);
	let guidance = $state<{ html: string } | { missing: string } | { error: string } | null>(null);
	let content: HTMLElement | undefined = $state();

	const tokens = (name: 'light' | 'dark') =>
		design?.themes.find((t) => t.name === name)?.tokens ?? [];

	async function loadGuidance(path: string) {
		if (!path) {
			guidance = null;
			return;
		}
		try {
			const r = await api(`/api/docs/file?path=${encodeURIComponent(path)}`);
			const body = await r.json().catch(() => ({}));
			if (r.status === 404) guidance = { missing: path };
			else if (!r.ok) guidance = { error: body.error ?? r.statusText };
			else {
				guidance = { html: render(body.body ?? '', path) };
				await tick();
				if (content) await enhance(content, themeState.dark);
			}
		} catch (e) {
			guidance = { error: e instanceof Error ? e.message : String(e) };
		}
	}

	async function load() {
		try {
			const r = await api('/api/design');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			design = body as DesignSystem;
			error = null;
			await loadGuidance(design.exists ? design.guidance : '');
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}

	onMount(() => {
		void load();
		// the token file and the guidance are documents under the design folder
		return follow(['document'], () => void load());
	});
</script>

<svelte:head><title>Design system · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Design system</h1>
<!-- a reload that fails says so above what was last shown rather than in its place -->
{#if error}
	<p class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if design === null}
	{#if !error}<p class="mt-4 text-sm text-muted">Loading…</p>{/if}
{:else if !design.exists}
	<div
		class="mt-4 rounded border border-line bg-surface p-4 text-sm"
		data-testid="no-design-system"
	>
		<p>{design.no_token_file ?? `This project has no design system: ${design.path} is missing.`}</p>
		<p class="mt-2">
			Write one with <code class="rounded bg-raised px-1.5 py-0.5 font-mono">flai design init</code
			>.
		</p>
	</div>
{:else}
	<p class="mt-1 text-sm text-muted">
		The tokens in <code>{design.path}</code>{#if design.guidance}
			and the guidance in <code>{design.guidance}</code>{/if}, previewed on sample components in
		light and dark.
	</p>
	{#if design.problems.length}
		<div
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			data-testid="design-problems"
		>
			<p class="font-medium">The token file has problems:</p>
			<ul class="mt-1 list-disc pl-5">
				{#each design.problems as p, i (i)}
					<li>
						line {p.line}{#if p.path}&nbsp;· <code>{p.path}</code>{/if}: {p.message}
					</li>
				{/each}
			</ul>
		</div>
	{/if}
	{#if design.themes.length}
		<section class="mt-6" aria-labelledby="design-preview-heading">
			<h2 id="design-preview-heading" class="mb-3 text-lg font-semibold">Preview</h2>
			<PreviewFrame light={tokens('light')} dark={tokens('dark')} />
		</section>
		<div class="mt-8">
			<TokenSheet themes={design.themes} contrast={design.contrast} />
		</div>
	{:else}
		<p class="mt-6 text-sm text-muted">
			The token file does not parse, so there are no tokens to show.
		</p>
	{/if}
	<section class="mt-8" aria-labelledby="design-guidance" data-testid="design-guidance">
		<h2 id="design-guidance" class="text-lg font-semibold">Guidance</h2>
		{#if guidance && 'missing' in guidance}
			<p class="mt-2 text-sm text-warn">
				The guidance document <code>{guidance.missing}</code> is missing.
			</p>
		{:else if guidance && 'error' in guidance}
			<p class="mt-2 rounded border border-danger bg-danger-soft p-3 text-sm text-danger">
				{guidance.error}
			</p>
		{:else if guidance}
			<article bind:this={content} class="prose mt-2 max-w-none">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the project's guidance document, rendered client side -->
				{@html guidance.html}
			</article>
		{/if}
	</section>
{/if}
