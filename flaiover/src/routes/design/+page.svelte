<script lang="ts">
	// /design: the project's design system (S-0382), flai's design.get through /api/design: the token
	// file's problems, its tokens laid out, its guidance rendered from the document the answer names,
	// and the tokens previewed on sample components in light and dark. A project with no token file
	// is told so, with the command that writes one. Both are read again when a document changes.
	// Edit tokens (S-0383) changes the tokens in place, one TokenEditor per namespace, and the
	// preview, the token sheet, and the contrast follow the edits before anything is saved. Save
	// writes the edits into the token file's raw text, read through /api/design/tokens, and saves it
	// through flai's file.save; flai's refusal, a change made meanwhile, or a flai too old to save is
	// shown with the edits kept. A document change while editing refreshes only the saved tokens, and
	// leaving with unsaved changes asks first.
	import type { DesignContrast, DesignSystem, DesignTheme } from '$lib/server/repo';
	import { api } from '$lib/api';
	import { onMount, tick } from 'svelte';
	import { beforeNavigate } from '$app/navigation';
	import { follow } from '$lib/events';
	import { enhance, render } from '$lib/markdown';
	import { themeState } from '$lib/theme.svelte';
	import {
		changed,
		composeTokenFile,
		editedThemes,
		pruneEdits,
		type Edits
	} from '$lib/design/tokens';
	import { PREVIEW_THEMES, rolePairs } from '$lib/design/preview';
	import TokenSheet from '$lib/components/design/TokenSheet.svelte';
	import PreviewFrame from '$lib/components/design/PreviewFrame.svelte';
	import TokenEditor from '$lib/components/design/TokenEditor.svelte';
	import DismissibleNotice from '$lib/components/DismissibleNotice.svelte';

	/** The token file as GET /api/design/tokens answers it: its raw text and its SHA-256. */
	type Raw = { path: string; content: string; hash: string };
	type Finding = { level: string; rule: string; path: string; line: number; message: string };
	type Conflict = { current: string; hash: string; diff?: string };

	const CANNOT_SAVE = 'This flai cannot save the token file yet; upgrade flai on the host.';
	// The token sheet's order; any other namespace follows in the token file's order.
	const NAMESPACES = ['color', 'font', 'text', 'spacing', 'radius'];
	const HEADINGS: Record<string, string> = {
		color: 'Colours',
		font: 'Font stacks',
		text: 'Type scale',
		spacing: 'Spacing',
		radius: 'Radii',
		shadow: 'Shadows'
	};

	let design = $state<DesignSystem | null>(null);
	let error = $state<string | null>(null);
	let guidance = $state<{ html: string } | { missing: string } | { error: string } | null>(null);
	let content: HTMLElement | undefined = $state();

	let editing = $state(false);
	let edits = $state<Edits>({});
	let raw = $state<Raw | null>(null);
	let rawError = $state<string | null>(null);
	// The host's flai answered 501: it cannot read or save the token file, so nothing is saved.
	let cannotSave = $state(false);
	let message = $state('');
	let saveError = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	let conflict = $state<Conflict | null>(null);
	let notice = $state<string | null>(null);
	let busy = $state(false);
	// Bumped to mount the editors afresh, dropping text they hold that made no edit.
	let generation = $state(0);

	const saved = $derived(design?.themes ?? []);
	// A plain token file's one value serves both themes, so a light edit reaches dark even before,
	// or without, its raw text; a resolver document's dark context is known only from its text.
	const rawText = $derived(raw?.content ?? (design && !design.resolver ? '{}' : null));
	const themes = $derived(editing ? editedThemes(saved, edits, rawText) : saved);
	const pairs = $derived(themes.flatMap((t) => rolePairs(t.name, t.tokens)));
	const contrast = $derived<DesignContrast[]>(
		editing
			? PREVIEW_THEMES.map((theme) => ({ theme, pairs: pairs.filter((p) => p.theme === theme) }))
			: (design?.contrast ?? [])
	);
	const dirty = $derived(changed(saved, edits).length > 0);
	const editable = $derived(!!design?.exists && saved.length > 0);
	const namespaces = $derived(namespacesOf(saved));
	const canSave = $derived(dirty && raw !== null && !cannotSave && !busy);

	const tokens = (name: 'light' | 'dark') => themes.find((t) => t.name === name)?.tokens ?? [];

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));

	/** Each token path's first part, in the token sheet's order, the rest as the file gives them. */
	function namespacesOf(all: readonly DesignTheme[]): string[] {
		const found = [...new Set(all.flatMap((t) => t.tokens.map((k) => k.path.split('.')[0])))];
		const rank = (n: string) => {
			const i = NAMESPACES.indexOf(n);
			return i < 0 ? NAMESPACES.length : i;
		};
		return found.sort((a, b) => rank(a) - rank(b));
	}

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
			guidance = { error: errorText(e) };
		}
	}

	// Reads the design system: while editing, only the saved themes change, and the edits stay.
	async function load() {
		try {
			const r = await api('/api/design');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			design = body as DesignSystem;
			error = null;
			await loadGuidance(design.exists ? design.guidance : '');
		} catch (e) {
			error = errorText(e);
		}
	}

	/** Reads the token file's raw text and hash, which the save writes the edits into. */
	async function loadRaw() {
		rawError = null;
		try {
			const r = await api('/api/design/tokens');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) {
				raw = null;
				cannotSave = true;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			raw = body as Raw;
		} catch (e) {
			raw = null;
			rawError = `Could not read the token file, so it cannot be saved: ${errorText(e)}`;
		}
	}

	function clearOutcome() {
		saveError = notice = null;
		findings = [];
		conflict = null;
	}

	function startEditing() {
		clearOutcome();
		editing = true;
		void loadRaw();
	}

	function stopEditing() {
		if (dirty) return;
		clearOutcome();
		editing = false;
		edits = {};
		raw = null;
		rawError = null;
	}

	function discard() {
		clearOutcome();
		edits = {};
		generation++;
	}

	async function save(overHash?: string) {
		if (!design || !raw || busy || !dirty || cannotSave) return;
		clearOutcome();
		let next: string;
		try {
			next = composeTokenFile(raw.content, pruneEdits(saved, edits));
		} catch (e) {
			saveError = errorText(e);
			return;
		}
		const path = design.path;
		busy = true;
		try {
			const r = await api('/api/design/tokens', {
				method: 'PUT',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					path,
					content: next,
					hash: overHash ?? raw.hash,
					message: message || undefined
				})
			});
			const data = await r.json().catch(() => ({}));
			if (r.status === 409) {
				conflict = { current: data.current, hash: data.hash, diff: data.diff };
				saveError = data.error ?? 'the token file changed after it was read';
			} else if (r.status === 422) {
				findings = data.findings ?? [];
				saveError = data.error ?? 'flai refused this token file';
			} else if (r.status === 501) {
				cannotSave = true;
			} else if (!r.ok) {
				saveError = data.error ?? r.statusText;
			} else {
				// The saved values become the new base: read them, then let the edits go.
				await load();
				await loadRaw();
				edits = {};
				message = '';
				generation++;
				notice = `Saved ${path}.`;
			}
		} catch (e) {
			saveError = errorText(e);
		} finally {
			busy = false;
		}
	}

	/** Takes the conflict's version as the base, keeping each edit that still has a place in it. */
	async function loadCurrent() {
		if (!raw || !conflict) return;
		const current = { ...raw, content: conflict.current, hash: conflict.hash };
		const kept: Edits = {};
		let dropped = 0;
		for (const [key, value] of Object.entries(edits)) {
			try {
				composeTokenFile(current.content, { [key]: value });
				kept[key] = value;
			} catch {
				dropped++;
			}
		}
		raw = current;
		clearOutcome();
		await load();
		edits = pruneEdits(saved, kept);
		generation++;
		const lost = `${dropped} had no place in it and ${dropped === 1 ? 'was' : 'were'} discarded`;
		notice = dropped
			? `Loaded the current version. Your edits are kept where its tokens still take them; ${lost}.`
			: 'Loaded the current version. Your edits are kept on it: save to write them.';
	}

	// Unsaved edits are only in this tab: ask before leaving it, inside the app and out of it.
	beforeNavigate((nav) => {
		if (dirty && !confirm('Leave without saving your changes?')) nav.cancel();
	});

	onMount(() => {
		void load();
		// the token file and the guidance are documents under the design folder
		return follow(['document'], () => void load());
	});
</script>

<svelte:window
	onbeforeunload={(e) => {
		if (dirty) e.preventDefault();
	}}
/>
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
	<div class="mt-1 flex flex-wrap items-baseline gap-3">
		<p class="text-sm text-muted">
			The tokens in <code>{design.path}</code>{#if design.guidance}
				and the guidance in <code>{design.guidance}</code>{/if}, previewed on sample components in
			light and dark.
		</p>
		{#if editable && !editing}
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1 text-sm"
				onclick={startEditing}>Edit tokens</button
			>
		{/if}
	</div>
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
		<div class={['mt-6', editing && 'grid grid-cols-1 gap-6 xl:grid-cols-2']}>
			<section
				class={[editing && 'self-start xl:sticky xl:top-4 xl:max-h-screen xl:overflow-auto']}
				aria-labelledby="design-preview-heading"
			>
				<h2 id="design-preview-heading" class="mb-3 text-lg font-semibold">Preview</h2>
				<PreviewFrame light={tokens('light')} dark={tokens('dark')} />
			</section>
			{#if editing}
				<section class="min-w-0" aria-labelledby="design-edit-heading" data-testid="design-edit">
					<div class="flex flex-wrap items-baseline gap-3">
						<h2 id="design-edit-heading" class="text-lg font-semibold">Edit tokens</h2>
						{#if dirty}<span class="text-xs text-warn">unsaved changes</span>{/if}
					</div>
					<p class="mt-1 text-sm text-muted">
						The preview and the token sheet follow each change; nothing is written until you save.
					</p>
					{#if cannotSave}
						<p
							class="mt-3 rounded border border-warn bg-warn-soft p-2 text-sm text-warn"
							role="status"
							data-testid="design-cannot-save"
						>
							{CANNOT_SAVE}
						</p>
					{/if}
					{#if rawError}
						<p class="mt-3 rounded border border-danger bg-danger-soft p-2 text-sm text-danger">
							{rawError}
						</p>
					{/if}
					{#if notice}
						<DismissibleNotice
							class="mt-3 rounded border border-good bg-good-soft p-2 text-sm text-good"
							testid="design-notice"
							text={notice}
							ondismiss={() => (notice = null)}
						/>
					{/if}
					{#if saveError}
						<div
							class="mt-3 rounded border border-danger bg-danger-soft p-2 text-sm text-danger"
							role="alert"
							data-testid="design-save-error"
						>
							<p>{saveError}</p>
							{#if findings.length}
								<ul class="mt-1 ml-4 list-disc font-mono text-xs">
									{#each findings as f, i (i)}
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
							disabled={!canSave}
							onclick={() => save()}>{busy ? 'Saving…' : 'Save'}</button
						>
						<button
							type="button"
							class="rounded border border-line-strong px-3 py-1 text-sm disabled:opacity-50"
							disabled={!dirty || busy}
							onclick={discard}>Discard changes</button
						>
						<button
							type="button"
							class="rounded border border-line px-3 py-1 text-sm disabled:opacity-50"
							disabled={dirty}
							title={dirty ? 'Save or discard your changes first' : undefined}
							onclick={stopEditing}>Done</button
						>
					</div>
					{#key generation}
						{#each namespaces as namespace (namespace)}
							<section class="mt-5" aria-labelledby="design-edit-{namespace}">
								<h3 id="design-edit-{namespace}" class="mb-1 font-semibold">
									{HEADINGS[namespace] ?? namespace}
								</h3>
								<TokenEditor {namespace} {saved} edited={themes} bind:edits {pairs} />
							</section>
						{/each}
					{/key}
				</section>
			{/if}
		</div>
		<div class="mt-8">
			<TokenSheet {themes} {contrast} />
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
