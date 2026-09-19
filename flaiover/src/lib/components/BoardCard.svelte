<script lang="ts">
	// One card on the kanban board: a single link to the item, draggable when
	// the board is writable. The bottom row carries nature, type, and the
	// blocked flag on the left and the parent's ID on the right (S-0048).
	import { resolve } from '$app/paths';
	import { age } from '$lib/age';

	type Card = {
		id: string;
		type: string;
		title: string;
		nature: string;
		parent?: string;
		parent_title?: string;
		status?: string;
		blocked: boolean;
		age_seconds: number;
	};

	let {
		card,
		draggable = false,
		dragging = false,
		ondragstart,
		ondragend
	}: {
		card: Card;
		draggable?: boolean;
		dragging?: boolean;
		ondragstart?: () => void;
		ondragend?: () => void;
	} = $props();

	// Epics have no parent. The indicator is not a link: the card is one.
	const parentLabel = $derived(
		card.parent ? (card.parent_title ? `${card.parent} ${card.parent_title}` : card.parent) : ''
	);
</script>

<a
	href={card.type === 'story' && card.status === 'review'
		? resolve('/review/[id]', { id: card.id })
		: resolve('/items/[id]', { id: card.id })}
	{draggable}
	{ondragstart}
	{ondragend}
	class="mb-2 block rounded border border-line bg-ground p-2 text-xs hover:border-line-strong {dragging
		? 'opacity-50'
		: ''}"
>
	<div class="flex items-center justify-between">
		<span class="font-mono font-medium">{card.id}</span>
		<span class="text-muted">{age(card.age_seconds)}</span>
	</div>
	<div class="mt-1 leading-snug break-words">{card.title}</div>
	<div class="mt-1 flex items-baseline gap-2 text-[10px] text-muted">
		<span>{card.nature}</span>
		{#if card.type !== 'story'}<span>{card.type}</span>{/if}
		{#if card.blocked}<span class="font-semibold text-danger">BLOCKED</span>{/if}
		{#if card.parent}
			<span
				class="ml-auto shrink-0 font-mono"
				data-testid="parent"
				title={parentLabel}
				aria-label="parent {parentLabel}">{card.parent}</span
			>
		{/if}
	</div>
</a>
