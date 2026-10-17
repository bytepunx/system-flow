<!--
	The picked component's examples (S-0443, ADR-0164): named sets of its props, kept in its
	examples file under design/components. Choosing one hands it up for the page to load its props,
	theme, and width into the panel; Save as example asks for a name and hands it up for the page to
	save the panel as it stands, on main or on the story the page reads. A save that fails shows
	flai's words and its findings, and the name stays to save again.
-->
<script lang="ts" module>
	import type { ComponentExample } from '$lib/server/repo';
	import type { Props, Theme } from '$lib/explorer/protocol';

	/** One finding of flai check's that refused a save. */
	export type Finding = {
		level: string;
		rule: string;
		path: string;
		line: number;
		message: string;
	};
	/** Why a save failed: flai's words, its findings for a refusal, and whether the file changed
	 * between the read and the save. */
	export type SaveFailure = { error: string; findings?: Finding[]; conflict?: boolean };

	/** The example name saves the panel as: its props, the theme when the page shows one, and the
	 * width when one is picked rather than the panel's own (ADR-0164 § 2). */
	export function exampleFrom(
		name: string,
		props: Props,
		themes: readonly Theme[],
		width: number | null
	): ComponentExample {
		return {
			name: name.trim(),
			props,
			...(themes.length === 1 ? { theme: themes[0] } : {}),
			...(width !== null ? { width } : {})
		};
	}
</script>

<script lang="ts">
	let {
		examples,
		path,
		story,
		onchoose,
		onsave
	}: {
		/** The component's examples, in its file's order. */
		examples: readonly ComponentExample[];
		/** Its examples file's repository path; without it there is no Save as example. */
		path?: string;
		/** The story whose branch a save goes on; undefined for main. */
		story?: string;
		/** Called with the example chosen. */
		onchoose: (example: ComponentExample) => void;
		/** Saves the panel as the example name; null once saved, else why not. */
		onsave: (name: string) => Promise<SaveFailure | null>;
	} = $props();

	let chosen = $state<string | null>(null);
	let naming = $state(false);
	let name = $state('');
	let busy = $state(false);
	let failure = $state<SaveFailure | null>(null);
	let saved = $state<string | null>(null);

	const replaces = $derived(examples.some((e) => e.name === name.trim()));

	function choose(example: ComponentExample) {
		chosen = example.name;
		onchoose(example);
	}

	function ask() {
		naming = true;
		name = chosen ?? '';
		failure = saved = null;
	}

	function cancel() {
		naming = false;
		failure = null;
	}

	async function save(e: Event) {
		e.preventDefault();
		const want = name.trim();
		if (!want || busy) return;
		busy = true;
		failure = saved = null;
		try {
			failure = await onsave(want);
			if (!failure) {
				saved = chosen = want;
				naming = false;
			}
		} finally {
			busy = false;
		}
	}
</script>

<section aria-labelledby="example-list-heading" data-testid="example-list">
	<div class="flex items-center justify-between gap-2">
		<h3 id="example-list-heading" class="text-sm font-semibold">Examples</h3>
		{#if path && !naming}
			<button
				type="button"
				class="rounded border border-line px-2 py-0.5 text-xs"
				data-testid="example-save-as"
				onclick={ask}>Save as example</button
			>
		{/if}
	</div>
	{#if !examples.length}
		<p class="mt-1 text-xs text-muted" data-testid="example-list-none">
			No examples yet. Save the panel's props as one to keep them.
		</p>
	{:else}
		<ul class="mt-1 flex flex-wrap gap-2">
			{#each examples as example (example.name)}
				{@const current = chosen === example.name}
				<li>
					<button
						type="button"
						class={[
							'rounded border px-2 py-0.5 text-sm',
							current ? 'border-primary bg-primary-soft font-medium' : 'border-line hover:bg-raised'
						]}
						aria-current={current ? 'true' : undefined}
						data-testid="example-{example.name}"
						onclick={() => choose(example)}>{example.name}</button
					>
				</li>
			{/each}
		</ul>
	{/if}
	{#if naming}
		<form class="mt-2 flex flex-col gap-2" onsubmit={save} data-testid="example-save-form">
			<label class="text-sm">
				<span class="mb-1 block text-xs text-muted">Example name</span>
				<input
					class="w-full rounded border border-line-strong bg-surface px-2 py-1 text-sm"
					bind:value={name}
					data-testid="example-name"
				/>
			</label>
			{#if replaces}
				<p class="text-xs text-muted" data-testid="example-replaces">
					Saving replaces the example of this name.
				</p>
			{/if}
			<div class="flex gap-2">
				<button
					class="rounded bg-primary px-3 py-1 text-sm text-on-primary disabled:opacity-50"
					disabled={!name.trim() || busy}
					data-testid="example-save">{busy ? 'Saving…' : 'Save'}</button
				>
				<button type="button" class="rounded border border-line px-3 py-1 text-sm" onclick={cancel}
					>Cancel</button
				>
			</div>
			<p class="text-xs text-muted">
				Saved to <span class="font-mono">{path}</span>
				{#if story}on <span class="font-mono">story/{story}</span>{:else}on main{/if}, and
				committed.
			</p>
		</form>
	{/if}
	{#if failure}
		<div
			class="mt-2 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="example-save-error"
		>
			<p>{failure.error}</p>
			{#if failure.conflict}
				<p class="mt-1 text-xs">
					The file changed while it was saved; save again to add to it as it is now.
				</p>
			{/if}
			{#if failure.findings?.length}
				<ul class="mt-1 list-disc pl-5">
					{#each failure.findings as f (f.rule + f.path + f.line + f.message)}
						<li class="font-mono text-xs">{f.path}:{f.line}: {f.level}: {f.rule}: {f.message}</li>
					{/each}
				</ul>
			{/if}
		</div>
	{/if}
	{#if saved}
		<p class="mt-1 text-xs text-muted" role="status" data-testid="example-saved">
			Saved the example “{saved}”.
		</p>
	{/if}
</section>
