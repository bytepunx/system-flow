<script lang="ts">
	// Editor for one document (ADR-0023, S-0040). flai decides what may be
	// edited and whether a save is accepted; this component shows the text,
	// a preview with the explorer's renderer, and what flai answered.
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { tick } from 'svelte';
	import { themeState } from '$lib/theme.svelte';
	import { render, enhance } from '$lib/markdown';
	import { compose, headingAt, headingsOf, splitRaw } from '$lib/edit';
	import { touching, type Worker } from '$lib/touches';
	import Threads from '$lib/components/Threads.svelte';

	type Doc = {
		path: string;
		content: string;
		hash: string;
		mode: 'full' | 'body' | 'none';
		reason?: string;
	};
	type Finding = { level: string; rule: string; path: string; line: number; message: string };
	type Conflict = { current: string; hash: string; diff?: string };

	let { path, ondirty }: { path: string; ondirty?: (dirty: boolean) => void } = $props();

	let doc = $state<Doc | null>(null);
	let frontMatter = $state<string | null>(null);
	let body = $state('');
	let message = $state('');
	let error = $state<string | null>(null);
	let notice = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	let conflict = $state<Conflict | null>(null);
	let busy = $state(false);
	let workers = $state<Worker[]>([]);
	let acknowledged = $state(false);
	let caretHeading = $state<string | null>(null);
	let threadHeading = $state<string | undefined>(undefined);
	let html = $state('');
	let preview: HTMLElement | undefined = $state();

	const content = $derived(compose(frontMatter, body));
	const dirty = $derived(doc !== null && content !== doc.content);
	const workedOn = $derived(touching(path, workers));
	const blockedByTouches = $derived(workedOn.length > 0 && !acknowledged);

	$effect(() => ondirty?.(dirty));

	function adopt(d: Doc) {
		doc = d;
		const parts = splitRaw(d.content);
		frontMatter = parts.frontMatter;
		body = parts.body;
		conflict = null;
		findings = [];
	}

	async function load(target: string) {
		error = notice = null;
		doc = null;
		try {
			const r = await api(`/api/docs/edit?path=${encodeURIComponent(target)}`);
			const data = await r.json();
			if (!r.ok) throw new Error(data.error ?? r.statusText);
			adopt(data);
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}

	$effect(() => {
		const target = path;
		acknowledged = false;
		void load(target);
		api('/api/items')
			.then(async (r) => (workers = r.ok ? await r.json() : []))
			.catch(() => (workers = []));
	});

	// Preview with the explorer's renderer, a beat after typing stops.
	$effect(() => {
		const text = body;
		const target = path;
		const timer = setTimeout(async () => {
			html = render(text, target);
			await tick();
			if (preview) await enhance(preview, themeState.dark);
		}, 200);
		return () => clearTimeout(timer);
	});

	async function save(overHash?: string) {
		if (!doc || busy) return;
		busy = true;
		error = notice = null;
		findings = [];
		try {
			const r = await api('/api/docs/file', {
				method: 'PUT',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					path: doc.path,
					content,
					hash: overHash ?? doc.hash,
					message: message || undefined
				})
			});
			const data = await r.json().catch(() => ({}));
			if (r.status === 409) {
				conflict = { current: data.current, hash: data.hash, diff: data.diff };
				error = data.error ?? 'the document changed after it was loaded';
			} else if (r.status === 422) {
				findings = data.findings ?? [];
				error = data.error ?? 'flai refused this content';
			} else if (!r.ok) {
				error = data.error ?? r.statusText;
			} else {
				notice = data.unchanged
					? 'Nothing to save: the document is unchanged.'
					: data.committed
						? `Saved and committed as ${data.commit}.`
						: data.commit_error
							? `Saved, but not committed: ${String(data.commit_error).split('\n')[0]}`
							: 'Saved. Not committed: this project leaves dashboard edits uncommitted, or it is not a git repository.';
				message = '';
				// flai may have set the updated date; take what is on disk now.
				const fresh = await api(`/api/docs/edit?path=${encodeURIComponent(doc.path)}`);
				if (fresh.ok) adopt(await fresh.json());
			}
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		} finally {
			busy = false;
		}
	}

	function loadCurrent() {
		if (!doc || !conflict) return;
		adopt({ ...doc, content: conflict.current, hash: conflict.hash });
		error = null;
		notice = 'Loaded the current version. Your edits were discarded.';
	}

	function caret(e: Event) {
		const el = e.currentTarget as HTMLTextAreaElement;
		caretHeading = headingAt(body, el.selectionStart ?? 0);
	}
</script>

{#if error && !doc}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{:else if !doc}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<div class="mb-3 flex flex-wrap items-baseline gap-3">
		<h1 class="text-xl font-semibold">Edit</h1>
		<a class="font-mono text-sm underline" href={resolve('/docs/[...path]', { path: doc.path })}
			>{doc.path}</a
		>
		{#if dirty}<span class="text-xs text-warn">unsaved changes</span>{/if}
	</div>

	{#if workedOn.length}
		<div class="mb-3 rounded border border-warn bg-warn-soft px-3 py-2 text-xs text-warn">
			<p>
				Being worked on by
				{#each workedOn as w, i (w.id)}{i ? ', ' : ' '}<a
						class="font-medium underline"
						href={resolve('/items/[id]', { id: w.id })}>{w.id}</a
					>
					{w.title}{/each}. An edit here may collide with that story's branch.
			</p>
			{#if doc.mode !== 'none'}
				<label class="mt-1 flex items-center gap-2">
					<input type="checkbox" bind:checked={acknowledged} /> I know; let me save anyway
				</label>
			{/if}
		</div>
	{/if}

	{#if doc.mode === 'none'}
		<p class="rounded border border-line bg-surface p-3 text-sm" role="status">
			This document cannot be edited here: {doc.reason}.
		</p>
	{:else}
		{#if frontMatter !== null}
			<section class="mb-3">
				<h2 class="mb-1 text-xs font-medium text-muted">
					Front matter{#if doc.mode === 'body'}
						· read-only: {doc.reason}{/if}
				</h2>
				{#if doc.mode === 'body'}
					<pre
						class="overflow-auto rounded border border-line bg-surface p-2 font-mono text-xs"
						data-testid="front-matter-readonly">{frontMatter}</pre>
				{:else}
					<textarea
						class="h-28 w-full rounded border border-line bg-ground p-2 font-mono text-xs"
						aria-label="Front matter (YAML)"
						spellcheck="false"
						bind:value={frontMatter}></textarea>
				{/if}
			</section>
		{/if}

		<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
			<section class="min-w-0">
				<h2 class="mb-1 text-xs font-medium text-muted">Body (markdown)</h2>
				<textarea
					class="h-[60vh] w-full rounded border border-line bg-ground p-2 font-mono text-sm"
					aria-label="Body (markdown)"
					spellcheck="true"
					bind:value={body}
					onkeyup={caret}
					onclick={caret}></textarea>
			</section>
			<section class="min-w-0">
				<h2 class="mb-1 text-xs font-medium text-muted">Preview</h2>
				<article
					bind:this={preview}
					class="prose h-[60vh] max-w-none overflow-auto rounded border border-line bg-surface p-3"
				>
					<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the editor, rendered client side like the explorer -->
					{@html html}
				</article>
			</section>
		</div>

		{#if notice}
			<p class="mt-3 rounded border border-good bg-good-soft p-2 text-sm text-good" role="status">
				{notice}
			</p>
		{/if}
		{#if error}
			<div
				class="mt-3 rounded border border-danger bg-danger-soft p-2 text-sm text-danger"
				role="alert"
			>
				<p>{error}</p>
				{#if findings.length}
					<ul class="mt-1 ml-4 list-disc font-mono text-xs">
						{#each findings as f (f.rule + f.path + f.line + f.message)}
							<li>{f.path}:{f.line}: {f.level}: {f.rule}: {f.message}</li>
						{/each}
					</ul>
				{/if}
				{#if conflict}
					{#if conflict.diff}
						<p class="mt-2 text-xs">What is there now (−) against what you are saving (+):</p>
						<pre
							class="mt-1 max-h-64 overflow-auto rounded border border-line bg-ground p-2 font-mono text-xs text-ink">{conflict.diff}</pre>
					{/if}
					<div class="mt-2 flex flex-wrap gap-2">
						<button
							type="button"
							class="rounded border border-line-strong px-3 py-1"
							disabled={busy}
							onclick={loadCurrent}>Load the current version</button
						>
						<button
							type="button"
							class="rounded border border-danger px-3 py-1"
							disabled={busy}
							onclick={() => save(conflict!.hash)}>Save mine over it</button
						>
					</div>
				{/if}
			</div>
		{/if}

		<div class="mt-3 flex flex-wrap items-center gap-2">
			<input
				class="min-w-0 flex-1 rounded border border-line bg-ground px-2 py-1 text-sm"
				placeholder="What changed (optional, becomes the commit subject)"
				aria-label="What changed"
				bind:value={message}
			/>
			<button
				type="button"
				class="rounded bg-primary px-3 py-1 text-sm text-on-primary disabled:opacity-50"
				disabled={busy || !dirty || blockedByTouches}
				onclick={() => save()}>{busy ? 'Saving…' : 'Save'}</button
			>
		</div>

		<div class="mt-4 flex flex-wrap items-center gap-2 text-xs">
			<button
				type="button"
				class="rounded border border-line px-2 py-1 disabled:opacity-50"
				disabled={!caretHeading}
				onclick={() => (threadHeading = caretHeading ?? undefined)}
				>Open a thread on {caretHeading
					? `“${caretHeading}”`
					: 'the heading under the cursor'}</button
			>
			<span class="text-muted">Click in the body to pick the heading.</span>
		</div>
	{/if}

	<Threads on={doc.path} headings={headingsOf(body)} compose={threadHeading} />
{/if}
