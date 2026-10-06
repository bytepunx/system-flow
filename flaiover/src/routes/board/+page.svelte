<script lang="ts">
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import LaneMenu, { type LaneAction } from '$lib/components/LaneMenu.svelte';
	import CardMenu from '$lib/components/CardMenu.svelte';
	import type { MenuClose } from '$lib/components/Menu.svelte';
	import type { CardEntry } from '$lib/cardmenu';
	import { longPress } from './longpress';
	import LaneMoveDialog from '$lib/components/LaneMoveDialog.svelte';
	import WipLimitDialog from '$lib/components/WipLimitDialog.svelte';
	import { backOf, forwardOf, laneCounts, longCounts, shortCounts } from '$lib/lanes';
	import AcceptConfirm from '$lib/components/AcceptConfirm.svelte';
	import CancelConfirm from '$lib/components/CancelConfirm.svelte';
	import BoardCard from '$lib/components/BoardCard.svelte';
	import BoardLegend from '$lib/components/BoardLegend.svelte';
	import BoardTypes from '$lib/components/BoardTypes.svelte';
	import { boardTypes, type ItemType } from '$lib/boardtypes.svelte';
	import { boardNatures } from '$lib/boardnatures.svelte';
	import HostAgentNotice from '$lib/components/HostAgentNotice.svelte';
	import { anyRunning, storyActivity, type HostAgent, type PlanRun } from '$lib/activity';
	import PublishBanner from '$lib/components/PublishBanner.svelte';
	import { doneLane, type RemoteTags, type Unplanned } from '$lib/publish';
	import DismissibleNotice from '$lib/components/DismissibleNotice.svelte';
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
	import { debounced, follow, listen } from '$lib/events';
	import type { TaskSummary } from '$lib/taskplan';

	type Card = {
		id: string;
		type: string;
		title: string;
		nature: string;
		parent?: string;
		parent_title?: string;
		status: string;
		blocked: boolean;
		draft?: boolean;
		age_seconds: number;
		archived?: boolean;
		tasks?: TaskSummary;
	};
	type Board = {
		wip_limits: Record<string, number>;
		order: string[];
		writable: boolean;
		columns: Record<string, Card[]>;
	};

	const states = ['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'];
	let board = $state<Board | null>(null);
	let notice = $state<{ kind: 'error' | 'warn' | 'ok'; text: string } | null>(null);
	let dragging = $state<string | null>(null);
	let over = $state<string | null>(null);

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
	}

	// Everything accepted and unreleased since each component's last tag (S-0087), for the done
	// column's Publish banner and to mark its own cards as waiting rather than published. Fetched
	// here, not inside the banner, so the two share one answer instead of asking twice (I-0033).
	type PendingItem = { id: string; title: string; level: string };
	type PendingPlan = {
		component: { name: string };
		level: string;
		from: string;
		to: string;
		items: PendingItem[];
	};
	let publishPlans = $state<PendingPlan[]>([]);
	let publishEnabled = $state(false);
	// how the clone's release tags stand against the remote's (S-0174), and what no plan covers
	// (I-0024)
	let publishRemote = $state<RemoteTags | null>(null);
	let publishUnplanned = $state<Unplanned[]>([]);
	async function loadPublish() {
		try {
			const r = await api('/api/publish');
			if (!r.ok) return;
			const body = await r.json();
			publishPlans = body.plans ?? [];
			publishEnabled = body.push_enabled === true;
			publishRemote = body.remote ?? null;
			publishUnplanned = body.unplanned ?? [];
		} catch {
			// keep what we had: a failed question is not news
		}
	}
	// What flai on the host knows of the agents it started (S-0104): a dot on each story's card and
	// the notice above the board, from one answer. An agent ends without changing a file, so while
	// one runs the board asks again now and then.
	let hostAgent = $state<HostAgent | null>(null);
	async function loadAgents() {
		try {
			const r = await api('/api/host-agent');
			if (r.ok) hostAgent = await r.json();
		} catch {
			// keep what we had
		}
	}
	const activity = $derived(storyActivity(hostAgent));
	$effect(() => {
		if (!anyRunning(hostAgent)) return;
		const t = setInterval(() => void loadAgents(), 15000);
		return () => clearInterval(t);
	});

	const waitingIds = $derived(
		new Set([
			...publishPlans.flatMap((p) => p.items.map((it) => it.id)),
			...publishUnplanned.map((u) => u.id)
		])
	);

	onMount(() => {
		load();
		void loadPublish();
		void loadAgents();
		// The board and the Publish banner are read from the work items: a narrative, a thread, or a
		// document changing asks for neither of them, and changes that arrive together ask once
		// (S-0161). The agents are read from the items and the threads, and flai serve says when one
		// starts or ends.
		const agents = debounced(() => void loadAgents());
		const stops = [
			follow(['item'], () => {
				load();
				void loadPublish();
			}),
			follow(['item', 'thread'], () => void loadAgents()),
			listen({ agent: () => agents() })
		];
		return () => {
			for (const stop of stops) stop();
			agents.stop();
		};
	});

	// What a lane holds: while the clone lags its remote's release tags, the done lane leaves out the
	// archived cards it holds only because they look unpublished (S-0174).
	const held = (state: string) => {
		const column = board?.columns[state] ?? [];
		return state === 'done' ? doneLane(column, publishRemote) : column;
	};
	// The types ticked above the board (S-0141) and the natures toggled on in its legend (S-0302):
	// a lane shows a card only when both let it. Neither hides anything from the WIP counts or
	// reordering, which count stories regardless, nor from the lane's counts by type (S-0256), which
	// count every card it holds.
	const cards = (state: string) =>
		held(state).filter(
			(c) => boardTypes.shown[c.type as ItemType] && boardNatures.isShown(c.nature)
		);
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

	// A lane's right-click menu (S-0167): create an item starting there, move its stories a column
	// forward or back, or change its WIP limit. A card's menu offers the same under the lane's name,
	// and a move from there starts with that story.
	let laneMenu = $state<{ lane: string; x: number; y: number } | null>(null);
	let laneMove = $state<{ from: string; to: string; picked: string[] } | null>(null);
	let laneLimit = $state<string | null>(null);
	const laneStories = (lane: string) =>
		(board?.columns[lane] ?? [])
			.filter((c) => c.type === 'story')
			.map((c) => ({ id: c.id, title: c.title }));

	// A card's menu (S-0202): what the item's page offers, then its lane's menu. A right click opens
	// it at the pointer; the keyboard, the card's menu button, and a long press below the card.
	let menuCard = $state<{ id: string; lane: string; x: number; y: number } | null>(null);
	const press = longPress();
	function below(el: Element) {
		const box = el.getBoundingClientRect();
		return { x: box.left, y: box.bottom };
	}
	function openCardMenu(id: string, lane: string, at: { x: number; y: number }) {
		notice = null;
		laneMenu = null;
		menuCard = { id, lane, ...at };
	}

	function openMenu(e: MouseEvent, lane: string) {
		// Shift with the mouse keeps the browser's own menu, for a card's link
		if (!board?.writable || (e.shiftKey && e.button === 2)) return;
		e.preventDefault();
		// from the keyboard there is no pointer: open beside what has focus
		const keyboard = e.clientX === 0 && e.clientY === 0;
		const card = (e.target as HTMLElement | null)?.closest<HTMLElement>('[data-card]');
		if (card) {
			// a long press that already opened it is the same press's native menu
			press.cancel();
			if (press.fired) return;
			openCardMenu(
				card.dataset.card!,
				lane,
				keyboard ? below(card) : { x: e.clientX, y: e.clientY }
			);
			return;
		}
		let { clientX: x, clientY: y } = e;
		if (keyboard) {
			const box = (e.target as HTMLElement).getBoundingClientRect();
			x = box.left + 8;
			y = box.bottom;
		}
		notice = null;
		menuCard = null;
		laneMenu = { lane, x, y };
	}
	function pickLane(action: LaneAction, lane: string, card?: string) {
		if (action === 'create') {
			// eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it
			void goto(resolve('/new') + `?status=${encodeURIComponent(lane)}`);
			return;
		}
		if (action === 'limit') {
			laneLimit = lane;
			return;
		}
		const to = action === 'forward' ? forwardOf(lane) : backOf(lane);
		if (!to) return;
		const picked = card && cardOf(card)?.type === 'story' ? [card] : [];
		laneMove = { from: lane, to, picked };
	}

	// A write from a card's menu, as the item's page makes it: say what came of it, from flai's
	// answer where that says more than the write's name (S-0263), and read the board again, so the
	// card shows it.
	async function write<A>(
		id: string,
		what: string,
		body: object,
		done: string | ((answer: A) => string)
	) {
		notice = null;
		const r = await api(`/api/items/${id}/${what}`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
		});
		const answer = await r.json().catch(() => ({}));
		notice = r.ok
			? { kind: 'ok', text: typeof done === 'string' ? done : done(answer) }
			: { kind: 'error', text: answer.error ?? r.statusText };
		await load();
	}
	async function pickCard(entry: CardEntry) {
		const id = menuCard!.id;
		menuCard = null;
		const c = cardOf(id);
		if (!c) return;
		switch (entry.action) {
			case 'open':
				// the card's own link: a story in review opens on its review
				if (c.type === 'story' && c.status === 'review') void goto(resolve('/review/[id]', { id }));
				else void goto(resolve('/items/[id]', { id }));
				return;
			case 'finalize':
				return write(id, 'finalize', {}, `${id} finalized`);
			case 'block': {
				const reason = prompt('Why is it blocked?');
				if (reason) await write(id, 'block', { reason }, `${id} blocked`);
				return;
			}
			case 'unblock':
				return write(id, 'unblock', {}, `${id} unblocked`);
			case 'agent':
				await write(
					id,
					'agent',
					{ action: entry.agent },
					`${id}: agent ${entry.label === 'Retry' ? 'restarted' : 'started'}`
				);
				return loadAgents();
			case 'plan':
				// as the item's Plan says it (S-0208): the planner's process and where its output goes
				await write(
					id,
					'plan',
					{},
					(run: Partial<PlanRun>) =>
						`planner started for ${id} (pid ${run.pid})${run.log ? `; its output is in ${run.log} on the host` : ''}`
				);
				return loadAgents();
			case 'cancel':
				cancelling = id;
		}
	}
	async function closeCardMenu(how: MenuClose) {
		const id = menuCard?.id;
		menuCard = null;
		// Escape gives focus back to the card it was opened on
		if (how !== 'escape' || !id) return;
		await tick();
		document.querySelector<HTMLElement>(`a[data-id="${id}"]`)?.focus();
	}

	// One move after another, each by flai's rules; one notice says what moved and what was refused.
	async function moveStories(ids: string[], to: string, reason: string) {
		const moved: string[] = [];
		const refused: string[] = [];
		const warnings: string[] = [];
		for (const id of ids) {
			const r = await api(`/api/items/${id}/move`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ to, reason: reason || undefined })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) refused.push(`${id}: ${body.error ?? r.statusText}`);
			else {
				moved.push(id);
				warnings.push(...(body.warnings ?? []));
			}
		}
		const text = [
			moved.length ? `${moved.join(', ')} → ${to}.` : '',
			refused.length ? `Refused: ${refused.join('; ')}.` : '',
			warnings.join(' ')
		]
			.filter(Boolean)
			.join(' ');
		notice = {
			kind: !moved.length ? 'error' : refused.length || warnings.length ? 'warn' : 'ok',
			text
		};
		laneMove = null;
		await load();
	}

	async function setLimit(lane: string, limit: number) {
		const r = await api('/api/board/limit', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ column: lane, limit })
		});
		const body = await r.json().catch(() => ({}));
		if (!r.ok) notice = { kind: 'error', text: body.error ?? r.statusText };
		else
			notice = {
				kind: 'ok',
				text: limit ? `WIP limit for ${lane}: ${limit}` : `${lane} has no WIP limit`
			};
		laneLimit = null;
		await load();
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

{#if laneMenu}
	<LaneMenu
		lane={laneMenu.lane}
		x={laneMenu.x}
		y={laneMenu.y}
		onpick={(action) => {
			const { lane } = laneMenu!;
			laneMenu = null;
			pickLane(action, lane);
		}}
		onclose={() => (laneMenu = null)}
	/>
{/if}

{#if menuCard}
	{@const card = cardOf(menuCard.id)}
	{#if card}
		<CardMenu
			{card}
			lane={menuCard.lane}
			writable={board?.writable === true}
			activity={activity[card.id]}
			agentEnabled={hostAgent?.enabled === true}
			planEnabled={hostAgent?.plan_enabled === true}
			x={menuCard.x}
			y={menuCard.y}
			oncard={pickCard}
			onlane={(action) => {
				const { id, lane } = menuCard!;
				menuCard = null;
				pickLane(action, lane, id);
			}}
			onclose={closeCardMenu}
		/>
	{/if}
{/if}

{#if laneMove}
	<LaneMoveDialog
		from={laneMove.from}
		to={laneMove.to}
		stories={laneStories(laneMove.from)}
		picked={laneMove.picked}
		oncancel={() => (laneMove = null)}
		onconfirm={(ids, reason) => moveStories(ids, laneMove!.to, reason)}
	/>
{/if}

{#if laneLimit && board}
	<WipLimitDialog
		lane={laneLimit}
		limit={board.wip_limits[laneLimit]}
		count={count(laneLimit)}
		oncancel={() => (laneLimit = null)}
		onconfirm={(n) => setLimit(laneLimit!, n)}
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
	<BoardTypes />
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
<HostAgentNotice status={hostAgent} />
{#if notice}
	<DismissibleNotice
		class="mb-3 rounded border p-2 text-sm {notice.kind === 'error'
			? 'border-danger bg-danger-soft text-danger'
			: notice.kind === 'warn'
				? 'border-warn bg-warn-soft text-warn'
				: 'border-good bg-good-soft text-good'}"
		role={notice.kind === 'error' ? 'alert' : 'status'}
		testid="board-notice"
		text={notice.text}
		ondismiss={() => (notice = null)}
	/>
{/if}
{#if board}
	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6">
		{#each states as state (state)}
			{@const counts = laneCounts(held(state))}
			<section
				class="min-h-40 rounded border bg-surface p-2 {over === state
					? 'border-accent'
					: 'border-line '}"
				role="group"
				aria-label={state}
				data-lane={state}
				oncontextmenu={(e) => openMenu(e, state)}
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
				<h2 class="flex items-baseline justify-between text-sm font-medium">
					<span>{state}</span>
					{#if board.wip_limits[state]}
						<span
							class="text-xs {count(state) > board.wip_limits[state]
								? 'text-danger'
								: 'text-muted'}">{count(state)}/{board.wip_limits[state]}</span
						>
					{/if}
				</h2>
				<p
					class="mb-2 text-xs text-muted"
					data-testid="lane-counts"
					title={longCounts(counts)}
					aria-label={longCounts(counts)}
				>
					{shortCounts(counts)}
				</p>
				{#if state === 'done'}
					<PublishBanner
						plans={publishPlans}
						enabled={publishEnabled}
						remote={publishRemote}
						unplanned={publishUnplanned}
						onpublished={() => {
							void loadPublish();
							load();
						}}
					/>
				{/if}
				{#each cards(state) as c (c.id)}
					<!-- The wrapper is the drop target for reordering and holds the controls beside the
					     card's link. A drop on a card of another column falls through to the column. -->
					<!-- On touch a long press opens the card's menu (S-0202); the browser's own callout
					     would cover it. -->
					<div
						class="group relative {board.writable ? '[-webkit-touch-callout:none]' : ''}"
						role="presentation"
						data-card={c.id}
						onpointerdown={(e) => {
							if (!board?.writable) return;
							const card = e.currentTarget;
							press.down(e, () => openCardMenu(c.id, state, below(card)));
						}}
						onpointermove={(e) => press.move(e)}
						onpointerup={() => press.up()}
						onpointercancel={() => press.up()}
						onclickcapture={(e) => press.click(e)}
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
							waiting={state === 'done' && waitingIds.has(c.id)}
							activity={activity[c.id]}
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
						{#if board.writable && !dragging}
							<!-- The card's menu from the keyboard (S-0202): beside the link, never in it, at the
							     top right, clear of the reorder controls. -->
							<button
								type="button"
								class="absolute top-1.5 right-1.5 z-10 rounded border border-line-strong bg-surface px-1.5 text-xs leading-4 text-ink hover:bg-raised {menuCard?.id ===
								c.id
									? 'block'
									: 'hidden group-focus-within:block group-hover:block'}"
								aria-label="Actions for {c.id}"
								aria-haspopup="menu"
								aria-expanded={menuCard?.id === c.id}
								title="Actions"
								data-testid="card-menu-button"
								data-for={c.id}
								onclick={(e) => openCardMenu(c.id, state, below(e.currentTarget.parentElement!))}
								>⋯</button
							>
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
