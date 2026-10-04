<script lang="ts">
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { follow } from '$lib/events';
	import Chart from '$lib/components/Chart.svelte';
	import SpendTable from '$lib/components/SpendTable.svelte';
	import {
		bucketsFor,
		build,
		completedIn,
		controls,
		FLOW_KINDS,
		hasSpend,
		human,
		KINDS,
		normalise,
		SPEND_KINDS,
		spendRows,
		titleOf,
		USAGE_KINDS,
		WINDOWS,
		withUsage,
		withCost,
		strategicCost,
		type BucketSize,
		type Kind,
		type Report
	} from '$lib/viz/charts';
	import { count, dollars } from '$lib/usage';
	import { theme } from '$lib/viz/palette';
	import KindChips from '$lib/components/KindChips.svelte';
	import { themeState } from '$lib/theme.svelte';
	import { chartWindow } from '$lib/chartwindow.svelte';

	// the window chosen last in this browser, whichever chart it was chosen on (S-0168)
	const since = $derived(chartWindow.since);
	let type = $state('story');
	let epic = $state('');
	// what spend over time is laid out in (S-0163)
	let bucket = $state<BucketSize>('day');
	let report = $state<Report | null>(null);
	let error = $state<string | null>(null);
	const dark = $derived(themeState.dark);
	let epics = $state<{ id: string; title: string }[]>([]);

	const kind = $derived(
		(KINDS as readonly string[]).includes(page.params.kind ?? '')
			? (page.params.kind as Kind)
			: 'cycle-time'
	);
	const t = $derived(theme(dark));
	const option = $derived(report ? build(kind, report, t, epic || undefined) : null);
	const s = $derived(report?.summary);
	const usageKind = $derived((USAGE_KINDS as readonly string[]).includes(kind));
	const spendKind = $derived(SPEND_KINDS.includes(kind));
	const shown = $derived(controls(kind));
	const buckets = $derived(bucketsFor(since));
	const title = $derived(titleOf(kind, report?.usage?.bucket ?? bucket));
	const spend = $derived(report?.usage?.spend?.[report.type]);
	// the items the window holds, which the tables list as the charts plot them (S-0166)
	// $ / Item also lists the items only strategic agents spent on, as it draws them (ADR-0083)
	const spenders = $derived(
		report ? (kind === 'cost' ? withCost : withUsage)(report, epic || undefined) : []
	);
	const completed = $derived(report ? completedIn(report) : []);
	/** What the items done in the window spent, and what that comes to, in a line. */
	const usageSummary = $derived.by(() => {
		const u = report?.usage;
		if (!report || !u) return '';
		const parts = [
			`${u.items} done`,
			`${count(u.tokens)} tokens`,
			dollars(u.cost) + (u.estimated ? ' (estimated in part)' : ''),
			...u.models.map((m) => `${m.model} ${dollars(m.cost)}`)
		];
		if (spend && spend.items > 0) {
			parts.push(
				`per ${report.type} ${count(spend.tokens_per_item ?? 0)} tokens, ${dollars(spend.cost_per_item ?? 0)}`
			);
			if (spend.tokens_per_minute !== undefined)
				parts.push(`${count(spend.tokens_per_minute)} tokens per agent minute`);
			if (spend.tokens_per_dollar !== undefined)
				parts.push(`${count(spend.tokens_per_dollar)} tokens per dollar`);
		}
		return parts.join(' · ');
	});
	/** A rate per agent minute; a flai older than S-0163 sends it per hour. */
	const perMinute = (m: { tokens_per_minute?: number; tokens_per_hour?: number }) =>
		m.tokens_per_minute ?? (m.tokens_per_hour !== undefined ? m.tokens_per_hour / 60 : undefined);

	// the latest question asked: an answer to an earlier one, arriving after it, is not drawn
	let asked = 0;
	async function load() {
		const mine = ++asked;
		// an hour is laid out over 31 days or less: a longer window goes by the day
		if (!buckets.includes(bucket)) bucket = 'day';
		const r = await api(`/api/stats?since=${since}&type=${type}&bucket=${bucket}`);
		const body = await r.json();
		if (mine !== asked) return;
		if (!r.ok) {
			error = body.error ?? r.statusText;
			report = null;
			return;
		}
		error = null;
		report = normalise(body);
	}
	/** The epics the per-item charts can be narrowed to; none when they cannot be read. */
	async function loadEpics() {
		const r = await api('/api/items?type=epic');
		const body = r.ok ? await r.json() : [];
		epics = Array.isArray(body)
			? body.map((e: { id: string; title: string }) => ({ id: e.id, title: e.title }))
			: [];
	}
	onMount(() => {
		// the statistics are read from the work items (S-0161)
		const stop = follow(['item'], () => void load());
		// the charts do not wait for the epics, and are drawn without them (S-0168)
		void load();
		void loadEpics();
		return stop;
	});
</script>

<svelte:head><title>{title} · flaiover</title></svelte:head>

{#each [{ name: 'flow', kinds: FLOW_KINDS }, { name: 'usage', kinds: USAGE_KINDS }] as group (group.name)}
	<nav
		class="flex flex-wrap items-center gap-2 {group.name === 'usage' ? 'mb-3' : 'mb-1'}"
		aria-label="{group.name} charts"
		data-testid="charts-{group.name}"
	>
		<span class="w-12 text-xs text-muted">{group.name}</span>
		{#each group.kinds as k (k)}
			<a
				href={resolve('/charts/[kind]', { kind: k })}
				class="rounded px-2 py-1 text-sm {k === kind
					? 'bg-raised font-medium '
					: 'text-ink-soft hover:text-ink '}">{titleOf(k, bucket)}</a
			>
		{/each}
	</nav>
{/each}
<div class="mb-4 flex flex-wrap items-center gap-3 text-sm">
	<label
		>window <select
			class="rounded border border-line-strong bg-surface px-2 py-1"
			value={since}
			onchange={(e) => {
				chartWindow.set(e.currentTarget.value);
				void load();
			}}
			data-testid="window"
			>{#each WINDOWS as w (w)}<option value={w}>{w}</option>{/each}</select
		></label
	>
	{#if shown.bucket}
		<label
			>per <select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={bucket}
				onchange={load}
				data-testid="bucket"
				>{#each buckets as b (b)}<option value={b}>{b}</option>{/each}</select
			></label
		>
	{/if}
	{#if shown.type}
		<label
			>type <select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={type}
				onchange={load}
				>{#each ['story', 'task', 'epic'] as ty (ty)}<option value={ty}>{ty}</option>{/each}</select
			></label
		>
	{/if}
	{#if shown.epic}
		<label
			>epic <select class="rounded border border-line-strong bg-surface px-2 py-1" bind:value={epic}
				><option value="">all</option>{#each epics as e (e.id)}<option value={e.id}
						>{e.id} {e.title}</option
					>{/each}</select
			></label
		>
	{/if}
	{#if usageKind && report?.usage && report.usage.items > 0}
		<span class="text-xs text-muted" data-testid="usage-summary">{usageSummary}</span>
	{:else if s}
		<span class="text-xs text-muted"
			>completed {s.completed} · cancelled {s.cancelled} · WIP {s.wip} · throughput {s.throughput_per_week.toFixed(
				1
			)}/week · cycle p50 {human(s.cycle_time.p50_seconds)} p85 {human(
				s.cycle_time.p85_seconds
			)}</span
		>
	{/if}
</div>
{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger">{error}</p>
{:else if option}
	<h1 class="mb-2 text-xl font-semibold">{title}</h1>
	{#if spendKind && report && !hasSpend(report)}
		<p class="mb-2 text-sm text-muted" data-testid="spend-none">
			The flai on the host sends no spend over time, which this chart is drawn from: it is older
			than the dashboard. Upgrade it with <code>flai self-upgrade</code>.
		</p>
	{:else if usageKind && report && !report.items.some((i) => i.usage)}
		<p class="mb-2 text-sm text-muted" data-testid="usage-none">
			No {report.type} here carries usage yet. flai serve records the tokens and cost of the agents it
			starts; <code>flai serve agent usage --all --write</code> fills in stories worked before.
		</p>
	{/if}
	{#if spendKind && report && hasSpend(report)}
		<p class="mb-2 text-xs text-muted" data-testid="spend-note">
			Each {report.usage?.bucket ?? bucket} holds the items that entered done in it, in UTC, with everything
			their agents spent on them.
			{#if kind === 'tokens-spent' || kind === 'cost-spent'}The dashed line is the mean per
				{report.usage?.bucket ?? bucket} so far, from the first with spend.{/if}
		</p>
	{/if}
	{#if kind === 'cost' || kind === 'cost-spent' || kind === 'cost-per-item' || kind === 'cost-per-model'}
		<p class="mb-2 text-xs text-muted">
			* estimated in part: a task's share of its story's session, or a run that ended without its
			totals.
		</p>
	{/if}
	<Chart {option} theme={t} height={380} />
	{#if kind === 'time-in-state' && report}
		<h2 class="mt-6 mb-2 text-base font-medium">Share of lead time per state</h2>
		{#await import('$lib/viz/charts') then m}
			<Chart option={m.stateShare(report, t)} theme={t} height={120} />
		{/await}
	{/if}
	<details class="mt-4 text-xs">
		<summary class="cursor-pointer text-muted">table view</summary>
		<div class="mt-2 overflow-x-auto">
			{#if spendKind && report}
				<SpendTable rows={spendRows(report, kind)} bucket={report.usage?.bucket ?? bucket} />
			{:else if usageKind && report}
				<table class="min-w-full" data-testid="usage-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">item</th><th class="pr-4">completed</th><th class="pr-4">model</th
							><th class="pr-4">tokens</th><th class="pr-4">tokens/agent minute</th><th>cost</th
							></tr
						></thead
					><tbody
						>{#each spenders as i (i.id)}{#each i.usage?.models ?? [] as m (m.model)}<tr
									><td class="pr-4 font-mono">{i.id}</td><td class="pr-4"
										>{i.completed?.slice(0, 10) ?? '-'}</td
									><td class="pr-4">{m.model}</td><td class="pr-4">{count(m.tokens)}</td><td
										class="pr-4">{perMinute(m) !== undefined ? count(perMinute(m)!) : '-'}</td
									><td>{dollars(m.cost)}{i.usage?.estimated ? '*' : ''}</td></tr
								>{/each}{#if kind === 'cost' && strategicCost(i) > 0}<tr
									data-testid="usage-strategic"
									><td class="pr-4 font-mono">{i.id}</td><td class="pr-4"
										>{i.completed?.slice(0, 10) ?? '-'}</td
									><td class="pr-4">strategic</td><td class="pr-4"
										>{count((i.usage?.strategic ?? []).reduce((n, x) => n + x.tokens, 0))}</td
									><td class="pr-4">-</td><td>{dollars(strategicCost(i))}*</td></tr
								>{/if}{/each}</tbody
					>
				</table>
			{:else if kind === 'throughput' && report}
				<table class="min-w-full">
					<thead
						><tr class="text-left text-muted"><th class="pr-4">week</th><th>done</th></tr></thead
					><tbody
						>{#each report.throughput as w (w.week)}<tr
								><td class="pr-4 font-mono">{w.week}</td><td>{w.done}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if report}
				<table class="min-w-full">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">item</th><th class="pr-4">nature</th><th class="pr-4">completed</th
							><th class="pr-4">cycle</th><th class="pr-4">lead</th><th>blocked</th></tr
						></thead
					><tbody
						>{#each completed as i (i.id)}<tr
								><td class="pr-4 font-mono">{i.id}</td><td class="py-0.5 pr-4"
									><KindChips nature={i.nature} /></td
								><td class="pr-4">{i.completed?.slice(0, 10)}</td><td class="pr-4"
									>{i.cycle_time_seconds !== undefined ? human(i.cycle_time_seconds) : '-'}</td
								><td class="pr-4"
									>{i.lead_time_seconds !== undefined ? human(i.lead_time_seconds) : '-'}</td
								><td>{human(i.blocked_seconds)}</td></tr
							>{/each}</tbody
					>
				</table>
			{/if}
		</div>
	</details>
{:else}
	<p class="text-sm text-muted">Loading…</p>
{/if}
