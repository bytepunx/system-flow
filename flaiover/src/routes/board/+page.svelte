<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { age } from '$lib/age';

	type Card = {
		id: string;
		type: string;
		title: string;
		nature: string;
		parent?: string;
		status: string;
		blocked: boolean;
		age_seconds: number;
	};
	type Board = {
		wip_limits: Record<string, number>;
		order: string[];
		writable: boolean;
		columns: Record<string, Card[]>;
	};

	const states = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'];
	let board = $state<Board | null>(null);
	let all = $state(false);
	let notice = $state<{ kind: 'error' | 'warn' | 'ok'; text: string } | null>(null);
	let dragging = $state<string | null>(null);
	let over = $state<string | null>(null);

	async function load() {
		board = await (await fetch('/api/board')).json();
	}
	onMount(() => {
		load();
		const es = new EventSource('/api/events');
		es.addEventListener('change', () => load());
		return () => es.close();
	});

	const cards = (state: string) =>
		(board?.columns[state] ?? []).filter((c) => all || c.type === 'story');
	const count = (state: string) =>
		(board?.columns[state] ?? []).filter((c) => c.type === 'story').length;

	async function move(id: string, to: string) {
		notice = null;
		let reason: string | undefined;
		const from = board?.columns[
			Object.keys(board.columns).find((s) => board!.columns[s].some((c) => c.id === id)) ?? ''
		]?.find((c) => c.id === id)?.status;
		if (to === 'cancelled' || (from === 'review' && to === 'in-progress')) {
			reason = prompt(`Reason for moving ${id} to ${to}:`) ?? undefined;
			if (!reason) return;
		}
		const r = await fetch(`/api/items/${id}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to, reason })
		});
		const body = await r.json();
		if (!r.ok) notice = { kind: 'error', text: body.error ?? r.statusText };
		else if (body.warnings?.length)
			notice = { kind: 'warn', text: `${id} → ${to}. ${body.warnings.join(' ')}` };
		else notice = { kind: 'ok', text: `${id} → ${to}` };
		await load();
	}
</script>

<svelte:head><title>Board · flaiover</title></svelte:head>

<div class="mb-3 flex flex-wrap items-center gap-4">
	<h1 class="text-2xl font-semibold">Board</h1>
	<label class="flex items-center gap-2 text-sm"
		><input type="checkbox" bind:checked={all} /> epics and tasks too</label
	>
	{#if board && !board.writable}
		<span
			class="rounded bg-amber-100 px-2 py-1 text-xs text-amber-900 dark:bg-amber-900 dark:text-amber-100"
			>read-only: flai is not available to the dashboard</span
		>
	{/if}
	{#if board?.order.length}
		<span class="text-xs text-zinc-500">pull order: {board.order.join(', ')}</span>
	{/if}
</div>
{#if notice}
	<p
		class="mb-3 rounded border p-2 text-sm {notice.kind === 'error'
			? 'border-red-300 bg-red-50 text-red-800'
			: notice.kind === 'warn'
				? 'border-amber-300 bg-amber-50 text-amber-900'
				: 'border-emerald-300 bg-emerald-50 text-emerald-900'}"
	>
		{notice.text}
	</p>
{/if}
{#if board}
	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6">
		{#each states as state (state)}
			<section
				class="min-h-40 rounded border bg-white p-2 dark:bg-zinc-900 {over === state
					? 'border-blue-400'
					: 'border-zinc-200 dark:border-zinc-800'}"
				role="group"
				aria-label={state}
				ondragover={(e) => {
					if (board?.writable && dragging) {
						e.preventDefault();
						over = state;
					}
				}}
				ondragleave={() => (over = null)}
				ondrop={(e) => {
					e.preventDefault();
					const id = dragging;
					dragging = null;
					over = null;
					if (id) move(id, state);
				}}
			>
				<h2 class="mb-2 flex items-baseline justify-between text-sm font-medium">
					<span>{state}</span>
					{#if board.wip_limits[state]}
						<span
							class="text-xs {count(state) > board.wip_limits[state]
								? 'text-red-600'
								: 'text-zinc-500'}">{count(state)}/{board.wip_limits[state]}</span
						>
					{/if}
				</h2>
				{#each cards(state) as c (c.id)}
					<a
						href={resolve('/items/[id]', { id: c.id })}
						draggable={board.writable}
						ondragstart={() => (dragging = c.id)}
						ondragend={() => (dragging = null)}
						class="mb-2 block rounded border border-zinc-200 bg-zinc-50 p-2 text-xs hover:border-zinc-400 dark:border-zinc-700 dark:bg-zinc-800 {dragging ===
						c.id
							? 'opacity-50'
							: ''}"
					>
						<div class="flex items-center justify-between">
							<span class="font-mono font-medium">{c.id}</span>
							<span class="text-zinc-500">{age(c.age_seconds)}</span>
						</div>
						<div class="mt-1 leading-snug">{c.title}</div>
						<div class="mt-1 flex gap-2 text-[10px] text-zinc-500">
							<span>{c.nature}</span>
							{#if c.type !== 'story'}<span>{c.type}</span>{/if}
							{#if c.blocked}<span class="font-semibold text-red-600">BLOCKED</span>{/if}
						</div>
					</a>
				{/each}
			</section>
		{/each}
	</div>
{:else}
	<p class="text-sm text-zinc-500">Loading…</p>
{/if}
