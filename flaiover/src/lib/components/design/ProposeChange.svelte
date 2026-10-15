<!--
	The Design page's request form (S-0384, ADR-0151 § 3): the operator's request for a change to
	the design system, in their words, becomes a story under the epic they choose, made with
	flai's item.new as the board's new-item form makes one, and moved to ready for a story's agent
	to pull. The operator wrote it, so it is not a draft. Its goal is their words; its criteria are
	fixed: the tokens as asked, the stylesheet in step, the UI's tests passing. Its topic is design
	and its touches are the token file and, when the manifest names one, the stylesheet.
-->
<script lang="ts">
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { NATURES } from '$lib/natures';

	type Finding = { level: string; rule: string; path: string; line: number; message: string };

	let {
		epics,
		tokens,
		stylesheet,
		oncreated
	}: {
		/** The open epics the story may go under. */
		epics: readonly { id: string; title: string }[];
		/** The token file's repository path: design.get's path. */
		tokens: string;
		/** The stylesheet the tokens are written into, when the manifest names one. */
		stylesheet?: string;
		/** Called with the new story's ID once it is made, moved or not. */
		oncreated?: (id: string) => void;
	} = $props();

	const CRITERIA = [
		'The tokens are changed as the goal asks.',
		'The generated stylesheet is in step with the tokens: `flai design css --write` is run, and `flai check` reports no `design.stale`.',
		"The UI's tests pass."
	];

	let parent = $state('');
	let nature = $state('feature');
	let title = $state('');
	let request = $state('');
	let busy = $state(false);
	let error = $state<string | null>(null);
	let findings = $state<Finding[]>([]);
	/** The story made, and whether it reached ready. */
	let made = $state<{ id: string; ready: boolean } | null>(null);

	// The epic chosen must be one offered: a closed epic drops out, and a lone epic is chosen.
	$effect(() => {
		if (parent && !epics.some((e) => e.id === parent)) parent = '';
		if (!parent && epics.length === 1) parent = epics[0].id;
	});

	const touches = $derived(stylesheet ? [tokens, stylesheet] : [tokens]);
	const ready = $derived(parent !== '' && title.trim() !== '' && request.trim() !== '' && !busy);

	/** The story's body, in the story template's sections: the goal in the operator's words. */
	function storyBody(goal: string): string {
		const criteria = CRITERIA.map((c) => `- [ ] ${c}`).join('\n');
		return `## Goal\n\n${goal.trim()}\n\n## Acceptance criteria\n\n${criteria}\n\n## Tasks\n\n## Notes\n`;
	}

	/** flai's words for a refusal, or the status when it gave none. */
	async function refusal(r: Response): Promise<{ error?: string; findings?: Finding[] }> {
		const answer = (await r.json().catch(() => ({}))) as { error?: string; findings?: Finding[] };
		return { ...answer, error: answer.error ?? r.statusText };
	}

	async function submit(e: Event) {
		e.preventDefault();
		if (!ready) return;
		busy = true;
		error = null;
		findings = [];
		made = null;
		try {
			const r = await api('/api/items', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					type: 'story',
					title: title.trim(),
					nature,
					parent,
					topics: ['design'],
					touches,
					body: storyBody(request)
				})
			});
			if (!r.ok) {
				// The words stay: a refusal is something to fix, not to retype.
				const refused = await refusal(r);
				error = `The story was not created: ${refused.error}`;
				findings = refused.findings ?? [];
				return;
			}
			const id = ((await r.json()) as { item: { id: string } }).item.id;
			// The story exists from here on, so its words are cleared whatever the move answers.
			title = '';
			request = '';
			const m = await api(`/api/items/${id}/move`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ to: 'ready' })
			});
			made = { id, ready: m.ok };
			if (!m.ok) {
				const why = (await refusal(m)).error;
				error = `${id} was created and is in backlog, but moving it to ready was refused: ${why}`;
			}
			oncreated?.(id);
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			busy = false;
		}
	}
</script>

<form class="space-y-3" onsubmit={submit} data-testid="propose-change">
	<p class="text-sm text-muted" data-testid="propose-about">
		Your request becomes the goal of a story with topic <code>design</code> that changes
		{touches.join(' and ')}, moved to ready for a story's agent to pull.
	</p>
	<div class="flex flex-wrap items-end gap-4">
		<label class="text-sm">
			<span class="mb-1 block text-xs text-muted">epic</span>
			<select
				class="max-w-xs rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={parent}
				data-testid="propose-epic"
			>
				<option value="" disabled>Choose an epic</option>
				{#each epics as e (e.id)}<option value={e.id}>{e.id} {e.title}</option>{/each}
			</select>
		</label>
		<label class="text-sm">
			<span class="mb-1 block text-xs text-muted">nature</span>
			<select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={nature}
				data-testid="propose-nature"
			>
				{#each NATURES as n (n.name)}<option value={n.name}>{n.name}</option>{/each}
			</select>
		</label>
	</div>
	<label class="block text-sm">
		<span class="mb-1 block text-xs text-muted">title</span>
		<input
			class="w-full rounded border border-line-strong bg-surface px-2 py-1"
			bind:value={title}
			maxlength="200"
			placeholder="the change, as a sentence"
			data-testid="propose-title"
		/>
	</label>
	<label class="block text-sm">
		<span class="mb-1 block text-xs text-muted">what to change, and why, in your words</span>
		<textarea
			class="h-32 w-full rounded border border-line-strong bg-surface px-2 py-1 text-sm"
			bind:value={request}
			spellcheck="true"
			data-testid="propose-request"></textarea>
	</label>

	{#if error}
		<div
			class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="propose-error"
		>
			<p>{error}</p>
			{#if findings.length}
				<ul class="mt-1 list-disc pl-5">
					{#each findings as f (f.rule + f.path + f.line + f.message)}
						<li class="font-mono text-xs">{f.path}:{f.line}: {f.level}: {f.rule}: {f.message}</li>
					{/each}
				</ul>
			{/if}
		</div>
	{/if}
	{#if made}
		<p class="text-sm" data-testid="propose-created">
			<a class="underline" href={resolve('/items/[id]', { id: made.id })}>{made.id}</a>
			{made.ready ? 'is ready for a story’s agent to pull.' : 'is in backlog.'}
		</p>
	{/if}

	<button
		class="rounded bg-primary px-4 py-2 text-sm text-on-primary disabled:opacity-50"
		disabled={!ready}
	>
		{busy ? 'Creating…' : 'Request the change'}
	</button>
</form>
