<script lang="ts">
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import AcceptConfirm from '$lib/components/AcceptConfirm.svelte';
	import CancelConfirm from '$lib/components/CancelConfirm.svelte';
	import BoardCard from '$lib/components/BoardCard.svelte';
	import BoardLegend from '$lib/components/BoardLegend.svelte';
	import UnpushedNotice from '$lib/components/UnpushedNotice.svelte';
	import CardReorder from '$lib/components/CardReorder.svelte';
	import {
		canReorderOnto,
		dropPlacement,
		endPlacement,
		keyStep,
		reorderable,
		stepPlacement,
		type Placement
	} from '$lib/reorder';
	import { onMount, tick } from 'svelte';

	type Card = {
		id: string;
		type: string;
		title: string;
		nature: string;
		parent?: string;
		parent_title?: string;
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

	// bumped on every load so the unpushed notice asks again when anything changes
	let loads = $state(0);
	async function load() {
		const r = await api('/api/board');
		const body = await r.json();
		// Without a flai on the host there is no board to show; the banner says why (S-0073).
		if (!r.ok) {
			if (!board) notice = { kind: 'error', text: body.error ?? r.statusText };
			return;
		}
		if (notice?.kind === 'error' && !board) notice = null;
		board = body;
		loads += 1;
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
		else if (to === 'cancelled') cancelling = id;
		else void move(id, to);
	}
	// A cancellation takes everything open under the item with it: show that first (S-0070).
	let cancelling = $state<string | null>(null);

	// Reordering within backlog and ready (S-0057). `marker` is where a drop would put the
	// dragged story: above or below a card, or at the end of its column.
	let marker = $state<{ id: string; half: 'above' | 'below' } | { end: string } | null>(null);
	const stories = (state: string) =>
		(board?.columns[state] ?? []).filter((c) => c.type === 'story').map((c) => c.id);
	const writable = () => board?.writable === true;

	// The column's empty space is below its cards; over a card the card's own handler decides.
	const onCard = (e: DragEvent) => (e.target as HTMLElement | null)?.closest('[data-card]') != null;

	async function place(id: string, p: Placement, refocus?: string) {
		notice = null;
		const r = await api(`/api/items/${id}/order`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(p)
		});
		const body = await r.json();
		if (!r.ok) notice = { kind: 'error', text: body.error ?? r.statusText };
		else notice = { kind: 'ok', text: `${body.status}: ${body.sequence.join(', ')}` };
		await load();
		// A reordered card is moved in the document, which drops keyboard focus: give it back,
		// to the same control when it still exists, else to the card.
		if (refocus) {
			await tick();
			const again =
				document.querySelector<HTMLElement>(refocus) ??
				document.querySelector<HTMLElement>(`a[data-id="${id}"]`);
			again?.focus();
		}
	}
	function step(c: Card, dir: 'up' | 'down', refocus: string) {
		if (!reorderable(c, writable())) return;
		const p = stepPlacement(stories(c.status), c.id, dir);
		if (p) void place(c.id, p, refocus);
	}

	async function move(id: string, to: string, includeUncommitted = false, why?: string) {
		notice = null;
		let reason = why;
		const from = board?.columns[
			Object.keys(board.columns).find((s) => board!.columns[s].some((c) => c.id === id)) ?? ''
		]?.find((c) => c.id === id)?.status;
		if (!reason && (to === 'cancelled' || (from === 'review' && to === 'in-progress'))) {
			reason = prompt(`Reason for moving ${id} to ${to}:`) ?? undefined;
			if (!reason) return;
		}
		const r = await api(`/api/items/${id}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to, reason, include_uncommitted: includeUncommitted || undefined })
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
		else if (body.cancelled?.length)
			notice = {
				kind: 'ok',
				text: `${id} → ${to}, and ${body.cancelled.length} under it: ${body.cancelled.map((c: { id: string }) => c.id).join(', ')}`
			};
		else notice = { kind: 'ok', text: `${id} → ${to}` };
		await load();
	}
</script>

<svelte:head><title>Board · flaiover</title></svelte:head>

{#if accepting}
	<AcceptConfirm
		id={accepting}
		oncancel={() => (accepting = null)}
		onconfirm={async (include) => {
			const id = accepting!;
			await move(id, 'done', include);
			accepting = null;
		}}
	/>
{/if}

{#if cancelling}
	<CancelConfirm
		id={cancelling}
		oncancel={() => (cancelling = null)}
		onconfirm={async (reason) => {
			const id = cancelling!;
			await move(id, 'cancelled', false, reason);
			cancelling = null;
		}}
	/>
{/if}

<div class="mb-3 flex flex-wrap items-center gap-4">
	<h1 class="text-2xl font-semibold">Board</h1>
	<label class="flex items-center gap-2 text-sm"
		><input type="checkbox" bind:checked={all} /> epics and tasks too</label
	>
	{#if board?.writable}
		<!-- Creating is a write through flai (S-0059): no flai, no action. -->
		<a
			class="rounded border border-line-strong bg-surface px-2 py-1 text-sm hover:bg-raised"
			href={resolve('/new')}
			data-testid="new-item-link">+ new</a
		>
	{/if}
	{#if board && !board.writable}
		<span class="rounded bg-warn-soft px-2 py-1 text-xs text-warn"
			>read-only: flai is not available to the dashboard</span
		>
	{/if}
	{#if board?.order.length}
		<span class="text-xs text-muted">pull order: {board.order.join(', ')}</span>
	{/if}
</div>
<div class="mb-3"><BoardLegend /></div>
<UnpushedNotice refresh={loads} />
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
					if (!board?.writable || !dragging) return;
					const from = cardOf(dragging);
					if (from?.status !== state) {
						e.preventDefault();
						over = state;
						marker = null;
					} else if (
						!onCard(e) &&
						reorderable(from, true) &&
						endPlacement(stories(state), from.id)
					) {
						// its own column, below the last card: a drop here sends it to the end
						e.preventDefault();
						marker = { end: state };
					}
					// Otherwise its own column with nowhere to go: not a drop target, so the pointer
					// says no and nothing looks as if it worked.
				}}
				ondragleave={() => {
					over = null;
					marker = null;
				}}
				ondrop={(e) => {
					e.preventDefault();
					const id = dragging;
					dragging = null;
					over = null;
					marker = null;
					const from = id ? cardOf(id) : undefined;
					if (!id || !from) return;
					if (from.status !== state) return request(id, state);
					const p =
						!onCard(e) && reorderable(from, writable()) ? endPlacement(stories(state), id) : null;
					if (p) void place(id, p);
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
					<!-- The wrapper is the drop target for reordering and holds the controls beside the
					     card's link. A drop on a card of another column falls through to the column. -->
					<div
						class="group relative"
						role="presentation"
						data-card={c.id}
						ondragover={(e) => {
							const from = dragging ? cardOf(dragging) : undefined;
							if (!canReorderOnto(from, c, writable())) {
								// not a place to reorder to: no marker left over from the last card
								marker = null;
								return;
							}
							e.preventDefault();
							e.stopPropagation();
							const box = e.currentTarget.getBoundingClientRect();
							const half = e.clientY < box.top + box.height / 2 ? 'above' : 'below';
							marker = dropPlacement(stories(state), from!.id, c.id, half)
								? { id: c.id, half }
								: null;
						}}
						ondrop={(e) => {
							const from = dragging ? cardOf(dragging) : undefined;
							if (!canReorderOnto(from, c, writable())) return;
							e.preventDefault();
							e.stopPropagation();
							const half = marker && 'id' in marker && marker.id === c.id ? marker.half : null;
							const p = half ? dropPlacement(stories(state), from!.id, c.id, half) : null;
							dragging = null;
							over = null;
							marker = null;
							if (p) void place(from!.id, p);
						}}
					>
						{#if marker && 'id' in marker && marker.id === c.id}
							<div
								class="pointer-events-none absolute inset-x-0 h-0.5 rounded bg-accent {marker.half ===
								'above'
									? '-top-[5px]'
									: '-bottom-[5px]'}"
								data-testid="drop-marker"
							></div>
						{/if}
						<BoardCard
							card={c}
							draggable={board.writable}
							dragging={dragging === c.id}
							ondragstart={() => {
								dragging = c.id;
								// the last action's notice must not read as this drag's result
								notice = null;
							}}
							ondragend={() => {
								dragging = null;
								marker = null;
							}}
							onkeydown={(e) => {
								const dir = keyStep(e);
								if (!dir || !reorderable(c, writable())) return;
								e.preventDefault();
								step(c, dir, `a[data-id="${c.id}"]`);
							}}
						/>
						{#if reorderable(c, board.writable) && !dragging}
							<CardReorder
								id={c.id}
								up={stepPlacement(stories(state), c.id, 'up')}
								down={stepPlacement(stories(state), c.id, 'down')}
								onplace={(p, dir) =>
									place(c.id, p, `button[data-step="${dir}"][data-for="${c.id}"]`)}
							/>
						{/if}
					</div>
				{/each}
				{#if marker && 'end' in marker && marker.end === state}
					<div class="pointer-events-none h-0.5 rounded bg-accent" data-testid="drop-marker"></div>
				{/if}
			</section>
		{/each}
	</div>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
