<script lang="ts">
	// Create an epic or a story by writing what it is for, in markdown, and nothing else
	// (S-0059). The ID, file name, front matter, parent link, and commit are flai's: this form
	// sends the text and the few choices that are the designer's to make, and shows what flai
	// answered. No front matter is shown, because none of it is the designer's to write.
	import { api } from '$lib/api';
	import { NATURES } from '$lib/natures';
	import { render } from '$lib/markdown';

	type Finding = { level: string; rule: string; path: string; line: number; message: string };
	type Epic = { id: string; title: string; status: string; archived: boolean };

	let {
		oncreated,
		initialType = 'story'
	}: { oncreated: (id: string) => void; initialType?: 'epic' | 'story' } = $props();

	let type = $state<'epic' | 'story'>(initialType);
	let title = $state('');
	let nature = $state('feature');
	let parent = $state('');
	let tags = $state('');
	let touches = $state('');
	let body = $state('');
	let templateBody = $state('');
	let epics = $state<Epic[]>([]);
	let error = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	let busy = $state(false);

	const meaning = $derived(NATURES.find((n) => n.name === nature)?.meaning ?? '');
	const html = $derived(render(body, 'wip/kanban/new.md'));
	const list = (s: string) =>
		s
			.split(/[\s,]+/)
			.map((v) => v.trim())
			.filter(Boolean);
	// A story need not belong to an epic (S-0092): "No epic" (parent === '') is a real choice, not
	// a placeholder to fill in.
	const ready = $derived(title.trim() !== '' && body.trim() !== '' && !busy);

	// The sections come from the project's item template, through flai. Text the designer has
	// already written is never replaced: a change of type swaps the body only while it is still
	// the template's.
	async function loadTemplate(t: 'epic' | 'story') {
		try {
			const r = await api(`/api/items/template?type=${t}`);
			if (!r.ok) return;
			const next = (await r.json()).body as string;
			if (body.trim() === '' || body === templateBody) body = next;
			templateBody = next;
		} catch {
			// the form still works from an empty body
		}
	}
	async function loadEpics() {
		try {
			const r = await api('/api/items?type=epic&archived=false');
			if (!r.ok) return;
			epics = ((await r.json()) as Epic[]).filter(
				(e) => e.status !== 'done' && e.status !== 'cancelled'
			);
			if (!parent && epics.length === 1) parent = epics[0].id;
		} catch {
			epics = [];
		}
	}
	$effect(() => {
		void loadTemplate(type);
	});
	$effect(() => {
		void loadEpics();
	});

	async function create(e: Event) {
		e.preventDefault();
		if (!ready) return;
		busy = true;
		error = null;
		findings = [];
		try {
			const r = await api('/api/items', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					type,
					title,
					nature,
					parent: type === 'story' ? parent : undefined,
					tags: list(tags),
					touches: list(touches),
					body
				})
			});
			const answer = await r.json();
			if (!r.ok) {
				// The text stays: a refusal is something to fix, not to retype.
				error = answer.error ?? r.statusText;
				findings = answer.findings ?? [];
				return;
			}
			oncreated(answer.item.id);
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			busy = false;
		}
	}
</script>

<form class="space-y-4" onsubmit={create} data-testid="new-item">
	<div class="flex flex-wrap items-end gap-4">
		<fieldset class="text-sm">
			<legend class="mb-1 text-xs text-muted">what</legend>
			<label class="mr-3"><input type="radio" bind:group={type} value="story" /> story</label>
			<label><input type="radio" bind:group={type} value="epic" /> epic</label>
		</fieldset>
		{#if type === 'story'}
			<label class="text-sm">
				<span class="mb-1 block text-xs text-muted">parent epic</span>
				<select
					class="max-w-xs rounded border border-line-strong bg-surface px-2 py-1"
					bind:value={parent}
					data-testid="parent"
				>
					<option value="">No epic</option>
					{#each epics as e (e.id)}<option value={e.id}>{e.id} {e.title}</option>{/each}
				</select>
			</label>
		{/if}
		<label class="text-sm">
			<span class="mb-1 block text-xs text-muted">nature</span>
			<select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={nature}
				data-testid="nature"
			>
				{#each NATURES as n (n.name)}<option value={n.name}>{n.name}</option>{/each}
			</select>
		</label>
	</div>
	<p class="text-xs text-muted" data-testid="meaning">{nature}: {meaning}</p>

	<label class="block text-sm">
		<span class="mb-1 block text-xs text-muted">title</span>
		<input
			class="w-full rounded border border-line-strong bg-surface px-3 py-2"
			bind:value={title}
			required
			maxlength="200"
			placeholder="what it delivers, as a sentence"
			data-testid="title"
		/>
	</label>

	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
		<label class="block text-sm">
			<span class="mb-1 block text-xs text-muted"
				>tags, optional: a component name decides which release it cuts</span
			>
			<input
				class="w-full rounded border border-line-strong bg-surface px-3 py-1.5"
				bind:value={tags}
				placeholder="dashboard cli"
				data-testid="tags"
			/>
		</label>
		<label class="block text-sm">
			<span class="mb-1 block text-xs text-muted">touches, optional: paths this work changes</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-3 py-1.5"
				bind:value={touches}
				placeholder="flaiover/src flai/cmd"
				data-testid="touches"
			/>
		</label>
	</div>

	<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
		<label class="block text-sm">
			<span class="mb-1 block text-xs text-muted"
				>what it is for, in markdown. The heading, ID, and front matter are added for you; write
				acceptance criteria as <code>- [ ]</code> lines</span
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
				<p class="mt-1 text-xs">Nothing was created. Your text is still here.</p>
			{/if}
		</div>
	{/if}

	<button
		class="rounded bg-primary px-4 py-2 text-sm text-on-primary disabled:opacity-50"
		disabled={!ready}
	>
		{busy ? 'Creating…' : `Create ${type}`}
	</button>
</form>
