<script lang="ts">
	// Review of one story (S-0041): what was asked, what the agent says, what
	// is being discussed, what changed, and what accepting will do; then accept
	// as the designer with each step shown, or send back with a reason.
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { criteriaOf, readNdjson, sectionOf } from '$lib/review';
	import DiffView from '$lib/components/DiffView.svelte';
	import Threads from '$lib/components/Threads.svelte';

	type Item = {
		id: string;
		type: string;
		title: string;
		status: string;
		body: string;
		path: string;
	};
	type Version = string | { Major: number; Minor: number; Patch: number };
	type Step = {
		component: { name: string };
		delivered: boolean;
		level: string;
		from: Version;
		to: Version;
		files: string[];
	};
	type Preview = {
		branch?: string;
		blockers?: string[];
		uncommitted?: string[];
		plan?: { level: string; commits: string[]; steps: Step[] | null; skipped?: string } | null;
	};
	type Progress = { step: string; msg: string };

	let { id }: { id: string } = $props();

	let item = $state<Item | null>(null);
	let narrative = $state<string | null>(null);
	let diff = $state<Record<string, unknown> | null>(null);
	let diffError = $state<string | null>(null);
	let preview = $state<Preview | null>(null);
	let previewError = $state<string | null>(null);
	let error = $state<string | null>(null);
	let writable = $state(false);

	let include = $state(false);
	let running = $state(false);
	let progress = $state<Progress[]>([]);
	let warnings = $state<string[]>([]);
	let result = $state<{ tags?: string[]; pushed?: boolean; push_error?: string } | null>(null);
	let failure = $state<string | null>(null);

	let sendingBack = $state(false);
	let reason = $state('');

	const v = (x: Version) => (typeof x === 'string' ? x : `${x.Major}.${x.Minor}.${x.Patch}`);
	const criteria = $derived(item ? criteriaOf(item.body) : []);
	const ticked = $derived(criteria.filter((c) => c.checked).length);
	const inReview = $derived(item?.status === 'review');
	const needsChoice = $derived(!!preview?.uncommitted?.length && !include);
	const canAccept = $derived(
		writable &&
			inReview &&
			!running &&
			!result &&
			!!preview &&
			!preview.blockers?.length &&
			!needsChoice
	);

	async function get<T>(path: string): Promise<T> {
		const r = await api(path);
		const data = await r.json();
		if (!r.ok) throw new Error(data.error ?? r.statusText);
		return data as T;
	}

	async function load(target: string) {
		error = diffError = previewError = failure = null;
		item = null;
		narrative = diff = preview = result = null;
		progress = [];
		warnings = [];
		try {
			item = (await get<{ item: Item }>(`/api/items/${target}`)).item;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
			return;
		}
		api('/api/board')
			.then(async (r) => (writable = r.ok ? (await r.json()).writable : false))
			.catch(() => (writable = false));
		get<{ body: string }>(`/api/docs/file?path=${encodeURIComponent(`wip/agents/${target}.md`)}`)
			.then((d) => (narrative = d.body))
			.catch(() => (narrative = ''));
		get<Record<string, unknown>>(`/api/items/${target}/diff`)
			.then((d) => (diff = d))
			.catch((e) => (diffError = e instanceof Error ? e.message : String(e)));
		if (item.status === 'review')
			get<Preview>(`/api/items/${target}/acceptance`)
				.then((p) => (preview = p))
				.catch((e) => (previewError = e instanceof Error ? e.message : String(e)));
	}
	$effect(() => {
		void load(id);
	});

	async function accept() {
		if (!canAccept) return;
		running = true;
		failure = null;
		progress = [];
		warnings = [];
		try {
			const r = await api(`/api/items/${id}/accept`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ include_uncommitted: include || undefined })
			});
			if (!r.ok) {
				failure = (await r.json().catch(() => ({}))).error ?? r.statusText;
				return;
			}
			await readNdjson(r, (l) => {
				if (l.event === 'progress')
					progress = [...progress, { step: String(l.step), msg: String(l.msg) }];
				else if (l.event === 'warning') warnings = [...warnings, String(l.msg)];
				else if (l.event === 'done') result = (l.result ?? {}) as typeof result;
				else if (l.event === 'error') failure = String(l.error);
			});
			if (!result && !failure) failure = 'the acceptance ended without a result';
		} catch (e) {
			failure = e instanceof Error ? e.message : String(e);
		} finally {
			running = false;
			if (failure) {
				// flai left the story in review; show it as it is now
				const fresh = await get<{ item: Item }>(`/api/items/${id}`).catch(() => null);
				if (fresh) item = fresh.item;
			}
		}
	}

	async function sendBack() {
		if (!reason.trim()) return;
		failure = null;
		const r = await api(`/api/items/${id}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to: 'in-progress', reason: reason.trim() })
		});
		const data = await r.json().catch(() => ({}));
		if (!r.ok) failure = data.error ?? r.statusText;
		else {
			sendingBack = false;
			reason = '';
			await load(id);
		}
	}
</script>

{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{:else if !item}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<div class="mb-4 flex flex-wrap items-baseline gap-3">
		<h1 class="text-xl font-semibold">Review {item.id}</h1>
		<a class="underline" href={resolve('/items/[id]', { id: item.id })}>{item.title}</a>
		<span class="rounded border border-line px-2 py-0.5 text-xs">{item.status}</span>
	</div>

	{#if result}
		<div class="mb-4 rounded border border-good bg-good-soft p-3 text-sm text-good" role="status">
			<p class="font-medium">{item.id} is accepted.</p>
			{#if result.tags?.length}<p>Released: {result.tags.join(', ')}</p>{:else}<p>
					Nothing was released.
				</p>{/if}
			{#if !result.pushed}
				<p>
					Not pushed{result.push_error ? ` (${result.push_error})` : ''}: push the commit{result
						.tags?.length
						? ' and the tags'
						: ''} from a shell.
				</p>
			{/if}
		</div>
	{:else if !inReview}
		<p class="mb-4 rounded border border-line bg-surface p-3 text-sm" role="status">
			{item.id} is {item.status}, not in review, so there is nothing to accept or send back. What
			follows is how it stands.
		</p>
	{/if}

	<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
		<section class="min-w-0 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">
				Acceptance criteria
				<span class="text-xs font-normal text-muted">{ticked} of {criteria.length} checked</span>
			</h2>
			{#if criteria.length}
				<ul class="space-y-1">
					{#each criteria as c (c.text)}
						<li class="flex gap-2">
							<span aria-hidden="true" class={c.checked ? 'text-good' : 'text-danger'}
								>{c.checked ? '☑' : '☐'}</span
							>
							<span class:text-muted={c.checked}
								><span class="sr-only">{c.checked ? 'checked: ' : 'not checked: '}</span
								>{c.text}</span
							>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-muted">The story has no acceptance criteria.</p>
			{/if}
		</section>

		<section class="min-w-0 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">
				Narrative
				<a
					class="text-xs font-normal underline"
					href={resolve('/docs/[...path]', { path: `wip/agents/${item.id}.md` })}>open</a
				>
			</h2>
			{#if narrative === null}
				<p class="text-muted">Loading…</p>
			{:else if !narrative}
				<p class="text-muted">No narrative for this story.</p>
			{:else}
				<h3 class="text-xs font-medium text-muted">Current state</h3>
				<p class="mb-2 whitespace-pre-wrap">{sectionOf(narrative, 'Current state') || '—'}</p>
				<h3 class="text-xs font-medium text-muted">Next steps</h3>
				<p class="whitespace-pre-wrap">{sectionOf(narrative, 'Next steps') || '—'}</p>
			{/if}
		</section>
	</div>

	<section class="mt-4 rounded border border-line bg-surface p-3 text-sm">
		<h2 class="mb-2 font-medium">Changes on the branch</h2>
		{#if diffError}
			<p class="text-muted">{diffError}</p>
		{:else if !diff}
			<p class="text-muted">Loading…</p>
		{:else}
			<DiffView diff={diff as never} />
		{/if}
	</section>

	{#if inReview || result}
		<section class="mt-4 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">What accepting does</h2>
			{#if previewError}
				<p class="rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					{previewError}
				</p>
			{:else if !preview && !result}
				<p class="text-muted">Working it out…</p>
			{:else if preview}
				{#if preview.blockers?.length}
					<div
						class="mb-2 rounded border border-danger bg-danger-soft p-2 text-danger"
						role="alert"
					>
						<p class="font-medium">This cannot be accepted from here yet:</p>
						<ul class="mt-1 ml-4 list-disc">
							{#each preview.blockers as b (b)}<li>{b}</li>{/each}
						</ul>
					</div>
				{/if}
				{#if preview.uncommitted?.length}
					<div class="mb-2 rounded border border-warn bg-warn-soft p-2 text-warn">
						<p class="font-medium">Uncommitted changes outside wip:</p>
						<ul class="mt-1 ml-4 list-disc font-mono text-xs">
							{#each preview.uncommitted as p (p)}<li>{p}</li>{/each}
						</ul>
						<label class="mt-2 flex items-center gap-2">
							<input type="checkbox" bind:checked={include} disabled={running} />
							Include these files in the acceptance commit
						</label>
					</div>
				{/if}
				<ul class="ml-4 list-disc space-y-1">
					{#if preview.branch}
						<li>
							Rebase <code>{preview.branch}</code>, fast-forward it into main, remove its worktree.
						</li>
					{/if}
					<li>Move the story to done, archive it with its tasks and narrative, and commit.</li>
					{#if preview.plan?.skipped}
						<li>No release: {preview.plan.skipped}.</li>
					{:else if preview.plan?.steps?.length}
						<li>
							Tag a {preview.plan.level} release:
							{#each preview.plan.steps as s, i (s.component.name)}{i ? ', ' : ''}<span
									class="font-medium">{s.component.name}</span
								>
								{v(s.from)} → {v(s.to)}{/each}
						</li>
					{/if}
					<li>Nothing is pushed from the dashboard; push from a shell afterwards.</li>
				</ul>
			{/if}

			{#if progress.length || running}
				<ol class="mt-3 space-y-1 border-t border-line pt-2" aria-live="polite">
					{#each progress as p (p.step)}
						<li><span class="text-good" aria-hidden="true">✓</span> {p.msg}</li>
					{/each}
					{#if running}<li class="text-muted">Working…</li>{/if}
				</ol>
			{/if}
			{#each warnings as w (w)}
				<p class="mt-2 rounded border border-warn bg-warn-soft p-2 text-warn">{w}</p>
			{/each}
			{#if failure}
				<div class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					<p class="font-medium">Not accepted. {item.id} is {item.status}. flai said:</p>
					<pre class="mt-1 font-mono text-xs whitespace-pre-wrap">{failure}</pre>
				</div>
			{/if}

			{#if writable && inReview && !result}
				<div class="mt-3 flex flex-wrap gap-2">
					<button
						type="button"
						class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
						disabled={!canAccept}
						onclick={accept}>{running ? 'Accepting…' : 'Accept'}</button
					>
					<button
						type="button"
						class="rounded border border-line-strong px-3 py-1"
						disabled={running}
						onclick={() => (sendingBack = !sendingBack)}>Send back</button
					>
				</div>
				{#if sendingBack}
					<div class="mt-2">
						<label class="block text-xs text-muted" for="send-back-reason"
							>Why it goes back to in-progress (recorded in the story's notes)</label
						>
						<textarea
							id="send-back-reason"
							class="mt-1 h-20 w-full rounded border border-line bg-ground p-2"
							bind:value={reason}></textarea>
						<button
							type="button"
							class="mt-1 rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!reason.trim()}
							onclick={sendBack}>Send back to in-progress</button
						>
					</div>
				{/if}
			{:else if !writable && inReview}
				<p class="mt-3 text-xs text-muted">Read-only: flai is not available to this dashboard.</p>
			{/if}
		</section>
	{/if}

	<section class="mt-4">
		<Threads on={item.id} {writable} />
	</section>
{/if}
