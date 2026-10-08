<script lang="ts">
	// One card on the kanban board: a single link to the item, draggable when
	// the board is writable. The bottom row carries nature, type, and the
	// blocked flag on the left and the parent's ID on the right (S-0048). The background
	// is tinted by nature and the left edge striped by type (S-0055); the text stays. A held story
	// says why it waits after BLOCKED: the hold's code and the stories it waits for (S-0129). A story
	// with tasks counts them by state, with its plan's layers (S-0176). A draft story, one an agent
	// wrote that the operator has not finalized, says DRAFT after BLOCKED (S-0201). A held story
	// whose holding stories' agents were asked about the hold names each of them, and a story
	// started on a share names the shared paths and the holder, each with a link to its conversation
	// on the messages page (S-0338). Those lines sit in a row under the card's link, drawn as part of
	// the card, since a link cannot hold another.
	import { resolve } from '$app/paths';
	import { age } from '$lib/age';
	import { stripeFor, tintFor } from '$lib/cardcolour';
	import {
		askedLine,
		askedTitle,
		elsewhereLine,
		holdLine,
		shareLine,
		shareTitle,
		type Share,
		type StoryActivity
	} from '$lib/activity';
	import { tasksLine, tasksTitle, type TaskSummary } from '$lib/taskplan';
	import AgentDot from './AgentDot.svelte';

	type Card = {
		id: string;
		type: string;
		title: string;
		nature: string;
		parent?: string;
		parent_title?: string;
		status?: string;
		blocked: boolean;
		draft?: boolean;
		age_seconds: number;
		tasks?: TaskSummary;
		shared?: Share[];
	};

	let {
		card,
		draggable = false,
		dragging = false,
		waiting = false,
		activity,
		ondragstart,
		ondragend,
		onkeydown
	}: {
		card: Card;
		draggable?: boolean;
		dragging?: boolean;
		/** Merged to main but not yet published (S-0087); the done column's own state, meaningless
		 * elsewhere. */
		waiting?: boolean;
		/** What the story's agent is doing (S-0104): a dot beside the ID. */
		activity?: StoryActivity;
		ondragstart?: () => void;
		ondragend?: () => void;
		/** Alt+arrow reorders a focused card on the board (S-0057); the card itself decides nothing. */
		onkeydown?: (e: KeyboardEvent) => void;
	} = $props();

	// Epics have no parent. The indicator is not a link: the card is one.
	const parentLabel = $derived(
		card.parent ? (card.parent_title ? `${card.parent} ${card.parent_title}` : card.parent) : ''
	);
	const asked = $derived(activity?.hold?.asked ?? []);
	const shared = $derived(card.type === 'story' ? (card.shared ?? []) : []);
	const talks = $derived(asked.length > 0 || shared.length > 0);
	// Messages.svelte gives each conversation its ID as an anchor.
	const conversation = (id: string) => `${resolve('/messages')}#${id}`;
</script>

<a
	href={card.type === 'story' && card.status === 'review'
		? resolve('/review/[id]', { id: card.id })
		: resolve('/items/[id]', { id: card.id })}
	{draggable}
	{ondragstart}
	{ondragend}
	{onkeydown}
	data-id={card.id}
	data-nature={card.nature}
	data-type={card.type}
	class="{talks
		? 'rounded-t'
		: 'mb-2 rounded'} block border border-line p-2 text-xs hover:border-t-line-strong hover:border-r-line-strong hover:border-b-line-strong {tintFor(
		card.nature
	)} {stripeFor(card.type)} {card.blocked ? 'ring-1 ring-danger' : ''} {dragging
		? 'border-dashed opacity-50'
		: ''}"
>
	<div class="flex items-center justify-between">
		<span class="flex items-center gap-1.5"
			><span class="font-mono font-medium">{card.id}</span><AgentDot {activity} /></span
		>
		<span class="text-muted">{age(card.age_seconds)}</span>
	</div>
	<div class="mt-1 leading-snug break-words">{card.title}</div>
	<!-- The details sit under a thin divider, in one size for the whole row, the title's:
	     the operator chose 12 px from a rendered comparison, and the parent ID inherits it
	     (S-0054). -->
	<div
		class="mt-2 flex flex-wrap items-baseline gap-x-2 gap-y-0.5 border-t border-line pt-2 text-xs text-muted"
		data-testid="details"
	>
		<span>{card.nature}</span>
		{#if card.type !== 'story'}<span>{card.type}</span>{/if}
		{#if card.blocked}<span class="font-semibold text-danger">BLOCKED</span>{/if}
		{#if card.type === 'story' && card.draft}<span
				class="font-semibold text-warn"
				data-testid="draft">DRAFT</span
			>{/if}
		{#if activity?.hold}<span class="font-semibold text-warn" data-testid="held"
				>{holdLine(activity.hold)}</span
			>{/if}
		{#if activity?.elsewhere}<span class="font-semibold text-warn" data-testid="elsewhere"
				>{elsewhereLine(activity.elsewhere)}</span
			>{/if}
		{#if card.type === 'story' && card.tasks}<span
				data-testid="tasks"
				title={tasksTitle(card.tasks)}
				aria-label={tasksTitle(card.tasks)}>{tasksLine(card.tasks)}</span
			>{/if}
		{#if waiting}<span class="font-semibold text-warn" data-testid="waiting-to-publish"
				>waiting to publish</span
			>{/if}
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
{#if talks}
	<div
		class="mb-2 flex flex-col gap-0.5 rounded-b border border-t-0 border-line px-2 pb-2 text-xs text-muted {tintFor(
			card.nature
		)} {stripeFor(card.type)} {dragging ? 'border-dashed opacity-50' : ''}"
		data-testid="talks"
	>
		<!-- eslint-disable svelte/no-navigation-without-resolve -- conversation() resolve()s the page; the rule does not follow the anchor added to it -->
		{#each asked as a (a.by)}<span data-testid="asked" title={askedTitle(a)}
				>{askedLine(a)} in
				<a
					class="font-mono underline"
					href={conversation(a.conversation)}
					data-conversation={a.conversation}>{a.conversation}</a
				></span
			>{/each}
		{#each shared as s, i (i)}<span class="break-all" data-testid="shared" title={shareTitle(s)}
				>{shareLine(s)}{s.conversation ? ' in ' : ''}{#if s.conversation}<a
						class="font-mono underline"
						href={conversation(s.conversation)}
						data-conversation={s.conversation}>{s.conversation}</a
					>{/if}</span
			>{/each}
		<!-- eslint-enable svelte/no-navigation-without-resolve -->
	</div>
{/if}
