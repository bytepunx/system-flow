<script lang="ts">
	// Edit a story's or an epic's own words where it is read (S-0085): title, nature, tags, what it
	// touches, the stories a story waits for (S-0130), its parent, and the body below its heading. flai on the host makes every change and
	// decides what is allowed; this form only collects it. What is the item's state, its ID, type,
	// status, owner, and dates, is shown and is not a field.
	import { api } from '$lib/api';
	import { render } from '$lib/markdown';
	import { agentFrom, configText, parseConfig, sameAgent, type Agent } from '$lib/agent';
	import AgentFields from './AgentFields.svelte';

	type View = {
		id: string;
		type: string;
		status: string;
		title: string;
		nature: string;
		tags: string[];
		touches: string[];
		/** the stories a story waits for (S-0130); a flai before it does not send it */
		after?: string[];
		parent?: string;
		agent?: Agent;
		default_agent?: Agent;
		body: string;
		path: string;
		hash: string;
		editable: boolean;
		reason?: string;
		natures: string[];
		parents: { id: string; title: string }[];
	};
	type Finding = { path: string; line: number; level: string; rule: string; message: string };

	let {
		id,
		oncancel,
		onsaved
	}: { id: string; oncancel: () => void; onsaved: (changed: string[]) => void } = $props();

	let view = $state<View | null>(null);
	let title = $state('');
	let nature = $state('');
	let tags = $state('');
	let touches = $state('');
	let after = $state('');
	let parent = $state('');
	// the story's own agent (S-0103), which a save replaces
	let harness = $state('');
	let model = $state('');
	let agentConfig = $state('');
	let body = $state('');
	let preview = $state(false);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	let conflict = $state<{ hash: string } | null>(null);

	const list = (s: string) =>
		s
			.split(',')
			.map((v) => v.trim())
			.filter(Boolean);
	const same = (a: string[], b: string[]) => a.length === b.length && a.every((v, i) => v === b[i]);

	function fill(v: View) {
		view = v;
		title = v.title;
		nature = v.nature;
		tags = v.tags.join(', ');
		touches = v.touches.join(', ');
		after = (v.after ?? []).join(', ');
		parent = v.parent ?? '';
		harness = v.agent?.harness ?? '';
		model = v.agent?.model ?? '';
		agentConfig = configText(v.agent?.config);
		body = v.body;
	}

	async function load() {
		error = null;
		findings = [];
		conflict = null;
		const r = await api(`/api/items/${id}/edit`);
		const data = await r.json().catch(() => ({}));
		if (!r.ok) {
			error = data.error ?? r.statusText;
			return;
		}
		fill(data as View);
	}
	$effect(() => {
		void id;
		void load();
	});

	/** Only what differs from what was loaded is sent, so an untouched field is never rewritten. */
	function change(): Record<string, unknown> {
		if (!view) return {};
		const out: Record<string, unknown> = {};
		if (title.trim() !== view.title) out.title = title;
		if (nature !== view.nature) out.nature = nature;
		if (!same(list(tags), view.tags)) out.tags = list(tags);
		if (view.type !== 'epic' && !same(list(touches), view.touches)) out.touches = list(touches);
		if (view.type !== 'epic' && parent && parent !== (view.parent ?? '')) out.parent = parent;
		if (view.type === 'story' && !same(list(after), view.after ?? [])) out.after = list(after);
		if (view.type === 'story') {
			const parsed = parseConfig(agentConfig);
			// a config line that is not key=value is said at save; nothing is sent until it is
			const next = 'error' in parsed ? view.agent : agentFrom(harness, model, parsed.config);
			if (!sameAgent(next, view.agent)) out.agent = next ?? null;
		}
		if (body.trim() !== view.body.trim()) out.body = body;
		return out;
	}
	const configError = $derived.by(() => {
		const parsed = parseConfig(agentConfig);
		return 'error' in parsed ? `agent config: ${parsed.error}` : null;
	});
	const dirty = $derived(Object.keys(change()).length > 0);

	async function save(hash?: string) {
		if (!view || saving) return;
		if (configError) {
			error = configError;
			return;
		}
		saving = true;
		error = null;
		findings = [];
		try {
			const r = await api(`/api/items/${id}/edit`, {
				method: 'PUT',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ ...change(), hash: hash ?? view.hash })
			});
			const data = await r.json().catch(() => ({}));
			if (r.status === 409) {
				conflict = { hash: data.hash };
				return;
			}
			if (r.status === 422) {
				error = data.error ?? 'the check refused this change';
				findings = data.findings ?? [];
				return;
			}
			if (!r.ok) {
				error = data.error ?? r.statusText;
				return;
			}
			conflict = null;
			onsaved(data.changed ?? []);
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		} finally {
			saving = false;
		}
	}
</script>

{#if !view}
	{#if error}
		<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
			{error}
		</p>
	{:else}
		<p class="text-sm text-muted">Loading…</p>
	{/if}
{:else if !view.editable}
	<p class="rounded border border-line-strong bg-raised p-3 text-sm" role="status">
		{view.reason}
		<button type="button" class="ml-2 underline" onclick={oncancel}>back</button>
	</p>
{:else}
	<form
		class="space-y-3"
		data-testid="item-editor"
		onsubmit={(e) => {
			e.preventDefault();
			void save();
		}}
	>
		<p class="text-xs text-muted">
			<span class="font-mono">{view.id}</span> · {view.type} · {view.status} ·
			<span class="font-mono">{view.path}</span>. These are flai's: the status changes by moving the
			card, the rest does not change.
		</p>
		<label class="block text-sm">
			<span class="mb-1 block font-medium">Title</span>
			<input
				class="w-full rounded border border-line-strong px-2 py-1"
				bind:value={title}
				maxlength="200"
				required
				data-testid="edit-title"
			/>
		</label>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Nature</span>
				<select class="w-full rounded border border-line-strong px-2 py-1" bind:value={nature}>
					{#each view.natures as n (n)}<option value={n}>{n}</option>{/each}
				</select>
			</label>
			{#if view.type !== 'epic'}
				<label class="block text-sm">
					<span class="mb-1 block font-medium">Parent</span>
					<select class="w-full rounded border border-line-strong px-2 py-1" bind:value={parent}>
						{#each view.parents as p (p.id)}<option value={p.id}>{p.id} {p.title}</option>{/each}
					</select>
				</label>
			{/if}
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Tags</span>
				<input
					class="w-full rounded border border-line-strong px-2 py-1"
					bind:value={tags}
					placeholder="comma separated"
				/>
			</label>
			{#if view.type !== 'epic'}
				<label class="block text-sm">
					<span class="mb-1 block font-medium">Touches</span>
					<input
						class="w-full rounded border border-line-strong px-2 py-1 font-mono text-xs"
						bind:value={touches}
						placeholder="paths or components, comma separated"
					/>
				</label>
			{/if}
			{#if view.type === 'story'}
				<label class="block text-sm">
					<span class="mb-1 block font-medium">Waits for</span>
					<input
						class="w-full rounded border border-line-strong px-2 py-1 font-mono text-xs"
						bind:value={after}
						placeholder="stories, such as S-0012, comma separated"
						data-testid="edit-after"
					/>
					<span class="mt-1 block text-xs text-muted"
						>It stays in ready, held, until each of these is done.</span
					>
				</label>
			{/if}
		</div>
		{#if view.type === 'story'}
			<AgentFields
				bind:harness
				bind:model
				bind:config={agentConfig}
				defaults={view.default_agent}
				note="Empty fields leave this story without them; the default applies only to new stories"
			/>
		{/if}
		<div>
			<div class="mb-1 flex items-center justify-between text-sm">
				<span class="font-medium">Body</span>
				<button type="button" class="text-xs underline" onclick={() => (preview = !preview)}
					>{preview ? 'write' : 'preview'}</button
				>
			</div>
			{#if preview}
				<article class="prose max-w-none rounded border border-line bg-surface p-3">
					<!-- eslint-disable-next-line svelte/no-at-html-tags -- the designer's own markdown, rendered client side as every document is -->
					{@html render(body, view.path)}
				</article>
			{:else}
				<textarea
					class="h-96 w-full rounded border border-line-strong p-2 font-mono text-sm"
					bind:value={body}
					spellcheck="true"
					data-testid="edit-body"></textarea>
			{/if}
			<p class="mt-1 text-xs text-muted">
				Below the heading: the heading is the ID and the title, and flai writes it.
			</p>
		</div>

		{#if conflict}
			<div
				class="rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
				role="alert"
				data-testid="edit-conflict"
			>
				<p>
					{view.id} changed after you opened it, an agent or someone else edited it. Nothing was saved.
				</p>
				<p class="mt-2 flex flex-wrap gap-2">
					<button type="button" class="rounded border border-warn px-2 py-1" onclick={load}
						>Load the current version (discards your edits)</button
					>
					<button
						type="button"
						class="rounded border border-warn px-2 py-1"
						onclick={() => save(conflict!.hash)}>Save mine over it</button
					>
				</p>
			</div>
		{/if}
		{#if error}
			<div
				class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
				role="alert"
				data-testid="edit-error"
			>
				<p>{error}</p>
				{#if findings.length}
					<ul class="mt-1 list-disc pl-5 text-xs">
						{#each findings as f (f.rule + f.path + f.line + f.message)}
							<li><span class="font-mono">{f.rule}</span>: {f.message}</li>
						{/each}
					</ul>
					<p class="mt-1 text-xs">Nothing was changed; your text is still here.</p>
				{/if}
			</div>
		{/if}
		<div class="flex gap-2">
			<button
				type="submit"
				class="rounded bg-primary px-3 py-1.5 text-sm text-on-primary disabled:opacity-50"
				disabled={!dirty || saving}
				data-testid="edit-save">{saving ? 'Saving…' : 'Save'}</button
			>
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1.5 text-sm"
				onclick={oncancel}>Cancel</button
			>
		</div>
	</form>
{/if}
