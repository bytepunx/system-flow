<script lang="ts">
	// Record an architecture decision by writing it, in markdown, and nothing else (S-0060).
	// The number, file name, front matter, date, index row, and commit are flai's. Shaped like
	// the new item form (S-0059): the designer's choices and text go out, and what flai
	// answered is shown, with the text kept on a refusal.
	import { api } from '$lib/api';
	import { render } from '$lib/markdown';

	type Finding = { level: string; rule: string; path: string; line: number; message: string };
	type Adr = { id: string; title: string; status: string };

	let { oncreated }: { oncreated: (path: string, id: string) => void } = $props();

	let title = $state('');
	let status = $state<'proposed' | 'accepted'>('proposed');
	let supersedes = $state<string[]>([]);
	let refines = $state<string[]>([]);
	let body = $state('');
	let existing = $state<Adr[]>([]);
	let error = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	let busy = $state(false);

	const html = $derived(render(body, 'design/adrs/new.md'));
	const both = $derived(supersedes.filter((id) => refines.includes(id)));
	const ready = $derived(title.trim() !== '' && body.trim() !== '' && both.length === 0 && !busy);

	$effect(() => {
		void (async () => {
			try {
				const r = await api('/api/adrs/template');
				if (r.ok && body.trim() === '') body = (await r.json()).body;
			} catch {
				// the form still works from an empty body
			}
		})();
		void (async () => {
			try {
				const r = await api('/api/docs/adrs');
				if (r.ok) existing = ((await r.json()) as Adr[]).filter((a) => a.id !== 'ADR-0000');
			} catch {
				existing = [];
			}
		})();
	});

	async function create(e: Event) {
		e.preventDefault();
		if (!ready) return;
		busy = true;
		error = null;
		findings = [];
		try {
			const r = await api('/api/adrs', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ title, status, supersedes, refines, body })
			});
			const answer = await r.json();
			if (!r.ok) {
				error = answer.error ?? r.statusText;
				findings = answer.findings ?? [];
				return;
			}
			oncreated(answer.path, answer.id);
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			busy = false;
		}
	}
</script>

<form class="space-y-4" onsubmit={create} data-testid="new-adr">
	<label class="block text-sm">
		<span class="mb-1 block text-xs text-muted">title: the decision, as a sentence</span>
		<input
			class="w-full rounded border border-line-strong bg-surface px-3 py-2"
			bind:value={title}
			required
			maxlength="200"
			placeholder="The dashboard authenticates every request with a project token"
			data-testid="title"
		/>
	</label>

	<fieldset class="text-sm">
		<legend class="mb-1 text-xs text-muted">status</legend>
		<label class="mr-4"
			><input type="radio" bind:group={status} value="proposed" /> proposed: a draft, still editable,
			accepted later from the ADRs page</label
		>
		<label
			><input type="radio" bind:group={status} value="accepted" /> accepted: decided, and immutable from
			now on</label
		>
	</fieldset>

	{#if existing.length}
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			{#each [{ label: 'supersedes: replaces these decisions, which are marked superseded', key: 'supersedes' }, { label: 'refines: narrows or extends these, which stay in force', key: 'refines' }] as group (group.key)}
				<fieldset class="text-sm" data-testid={group.key}>
					<legend class="mb-1 text-xs text-muted">{group.label}</legend>
					<div class="max-h-32 overflow-auto rounded border border-line bg-surface p-2">
						{#each existing as a (a.id)}
							<label class="block truncate text-xs" title={a.title}>
								{#if group.key === 'supersedes'}
									<input type="checkbox" bind:group={supersedes} value={a.id} />
								{:else}
									<input type="checkbox" bind:group={refines} value={a.id} />
								{/if}
								<span class="font-mono">{a.id}</span>
								{a.title}
							</label>
						{/each}
					</div>
				</fieldset>
			{/each}
		</div>
		{#if both.length}
			<p class="text-xs text-danger" data-testid="both">
				{both.join(', ')} cannot be both superseded and refined; choose one.
			</p>
		{/if}
	{/if}

	<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
		<label class="block text-sm">
			<span class="mb-1 block text-xs text-muted"
				>the decision, in markdown, under your template's sections. The number, heading, front
				matter, and index row are added for you</span
			>
			<textarea
				class="h-[55vh] w-full rounded border border-line-strong bg-surface p-3 font-mono text-sm"
				bind:value={body}
				spellcheck="true"
				data-testid="body"></textarea>
		</label>
		<div>
			<span class="mb-1 block text-xs text-muted">preview</span>
			<div
				class="prose h-[55vh] max-w-none overflow-auto rounded border border-line bg-surface p-3"
				data-testid="preview"
			>
				<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown the designer is typing, rendered client side like the explorer and the editor -->
				{@html html}
			</div>
		</div>
	</div>

	{#if error}
		<div
			class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="refusal"
		>
			<p>{error}</p>
			{#if findings.length}
				<ul class="mt-1 list-disc pl-5">
					{#each findings as f (f.rule + f.path + f.line + f.message)}
						<li><span class="font-mono text-xs">{f.rule}</span> {f.message}</li>
					{/each}
				</ul>
			{/if}
			<p class="mt-1 text-xs">Nothing was created. Your text is still here.</p>
		</div>
	{/if}

	<button
		class="rounded bg-primary px-4 py-2 text-sm text-on-primary disabled:opacity-50"
		disabled={!ready}
	>
		{busy ? 'Recording…' : `Record as ${status}`}
	</button>
</form>
