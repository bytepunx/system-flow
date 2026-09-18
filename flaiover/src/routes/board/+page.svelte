<script lang="ts">
	import { api } from '$lib/api';
	import AcceptConfirm from '$lib/components/AcceptConfirm.svelte';
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
		board = await (await api('/api/board')).json();
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

	// A story dropped on done is an acceptance: confirm with the plan first (S-0046).
	let accepting = $state<string | null>(null);
	function cardOf(id: string): Card | undefined {
		for (const cards of Object.values(board?.columns ?? {})) {
			const c = cards.find((x) => x.id === id);
			if (c) return c;
		}
	}
	function request(id: string, to: string) {
		const card = cardOf(id);
		if (to === 'done' && card?.type === 'story' && card.status === 'review') accepting = id;
		else void move(id, to);
	}

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
		const r = await api(`/api/items/${id}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to, reason })
		});
		const body = await r.json();
		if (!r.ok) notice = { kind: 'error', text: body.error ?? r.statusText };
		else if (body.warnings?.length)
			notice = { kind: 'warn', text: `${id} → ${to}. ${body.warnings.join(' ')}` };
		else if (body.push_error)
			notice = {
				kind: 'warn',
				text: `${id} accepted locally${body.tags?.length ? ` (${body.tags.join(', ')})` : ''} but not pushed: ${body.push_error}. Push the commit and tags from a shell.`
			};
		else if (body.tags?.length)
			notice = { kind: 'ok', text: `${id} accepted: released ${body.tags.join(', ')}` };
		else notice = { kind: 'ok', text: `${id} → ${to}` };
		await load();
	}
</script>

<svelte:head><title>Board · flaiover</title></svelte:head>

{#if accepting}
	<AcceptConfirm
		id={accepting}
		oncancel={() => (accepting = null)}
		onconfirm={async () => {
			const id = accepting!;
			await move(id, 'done');
			accepting = null;
		}}
	/>
{/if}

<div class="mb-3 flex flex-wrap items-center gap-4">
	<h1 class="text-2xl font-semibold">Board</h1>
	<label class="flex items-center gap-2 text-sm"
		><input type="checkbox" bind:checked={all} /> epics and tasks too</label
	>
	{#if board && !board.writable}
		<span class="rounded bg-warn-soft px-2 py-1 text-xs text-warn"
			>read-only: flai is not available to the dashboard</span
		>
	{/if}
	{#if board?.order.length}
		<span class="text-xs text-muted">pull order: {board.order.join(', ')}</span>
	{/if}
</div>
{#if notice}
	<p
		class="mb-3 rounded border p-2 text-sm {notice.kind === 'error'
			? 'border-danger bg-danger-soft text-danger'
			: notice.kind === 'warn'
				? 'border-warn bg-warn-soft text-warn'
				: 'border-good bg-good-soft text-good'}"
	>
		{notice.text}
	</p>
{/if}
{#if board}
	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6">
		{#each states as state (state)}
			<section
				class="min-h-40 rounded border bg-surface p-2 {over === state
					? 'border-accent'
					: 'border-line '}"
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
					if (id) request(id, state);
				}}
			>
				<h2 class="mb-2 flex items-baseline justify-between text-sm font-medium">
					<span>{state}</span>
					{#if board.wip_limits[state]}
						<span
							class="text-xs {count(state) > board.wip_limits[state]
								? 'text-danger'
								: 'text-muted'}">{count(state)}/{board.wip_limits[state]}</span
						>
					{/if}
				</h2>
				{#each cards(state) as c (c.id)}
					<a
						href={resolve('/items/[id]', { id: c.id })}
						draggable={board.writable}
						ondragstart={() => (dragging = c.id)}
						ondragend={() => (dragging = null)}
						class="mb-2 block rounded border border-line bg-ground p-2 text-xs hover:border-line-strong {dragging ===
						c.id
							? 'opacity-50'
							: ''}"
					>
						<div class="flex items-center justify-between">
							<span class="font-mono font-medium">{c.id}</span>
							<span class="text-muted">{age(c.age_seconds)}</span>
						</div>
						<div class="mt-1 leading-snug">{c.title}</div>
						<div class="mt-1 flex gap-2 text-[10px] text-muted">
							<span>{c.nature}</span>
							{#if c.type !== 'story'}<span>{c.type}</span>{/if}
							{#if c.blocked}<span class="font-semibold text-danger">BLOCKED</span>{/if}
						</div>
					</a>
				{/each}
			</section>
		{/each}
	</div>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
