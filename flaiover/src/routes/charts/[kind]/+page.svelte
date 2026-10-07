<script lang="ts">
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { onMount, untrack } from 'svelte';
	import { follow } from '$lib/events';
	import Chart from '$lib/components/Chart.svelte';
	import SpendTable from '$lib/components/SpendTable.svelte';
	import ForecastTable from '$lib/components/ForecastTable.svelte';
	import WaitTable from '$lib/components/WaitTable.svelte';
	import {
		bucketsFor,
		bucketLabel,
		build,
		COD_COLUMNS,
		COD_KINDS,
		codDayRows,
		codOrderRows,
		codWeekRows,
		completedIn,
		controls,
		driftedIn,
		errorFacets,
		FLOW_KINDS,
		forecastRows,
		hasClaims,
		hasCostOfDelay,
		hasOrder,
		hasSpend,
		hours,
		human,
		isClaimsKind,
		isForecastKind,
		KINDS,
		normalise,
		PLANNING_KINDS,
		plural,
		ORDER_LABEL,
		ORDERS,
		orderSummary,
		SPEND_KINDS,
		spendRows,
		STRATEGIC_AGENTS,
		STRATEGIC_KINDS,
		strategicRatio,
		strategicRows,
		titleOf,
		USAGE_KINDS,
		WINDOWS,
		withUsage,
		withCost,
		spreadFor,
		withoutValueNow,
		strategicCost,
		type BucketSize,
		type ErrorField,
		type ErrorFilter,
		type ErrorSpread,
		type Kind,
		type Report,
		type StrategicUse
	} from '$lib/viz/charts';
	import { count, dollars } from '$lib/usage';
	import { amount } from '$lib/planning';
	import { localDate } from '$lib/localtime';
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
	// the nature and the model the planning charts are narrowed to, all when empty (S-0212)
	let nature = $state('');
	let model = $state('');
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
	const shown = $derived(controls(kind));
	/** The groups of charts the header lists, in order. */
	const groups = [
		{ name: 'flow', kinds: FLOW_KINDS },
		{ name: 'usage', kinds: USAGE_KINDS },
		{ name: 'planning', kinds: [...PLANNING_KINDS, ...COD_KINDS] },
		{ name: 'strategic', kinds: STRATEGIC_KINDS }
	];
	const filter = $derived<ErrorFilter>({
		nature: (shown.nature && nature) || undefined,
		model: (shown.model && model) || undefined
	});
	const option = $derived(report ? build(kind, report, t, epic || undefined, filter) : null);
	const s = $derived(report?.summary);
	const usageKind = $derived((USAGE_KINDS as readonly string[]).includes(kind));
	const spendKind = $derived(SPEND_KINDS.includes(kind));
	const planningKind = $derived((PLANNING_KINDS as readonly string[]).includes(kind));
	// the forecast charts of the planning group, and the claims charts beside them (S-0214)
	const forecastKind = $derived(isForecastKind(kind));
	const claimsKind = $derived(isClaimsKind(kind));
	const strategicKind = $derived((STRATEGIC_KINDS as readonly string[]).includes(kind));
	// the planning charts read stories, whichever type was chosen on another chart (S-0212)
	const asType = $derived(planningKind ? 'story' : type);
	const facets = $derived(
		report && isForecastKind(kind) ? errorFacets(report, kind) : { natures: [], models: [] }
	);
	const forecasts = $derived(report && forecastKind ? forecastRows(report, filter) : []);
	/** Whether any story done in the window has a forecast of its duration or its delivery. */
	const forecastSet = $derived(
		report !== null &&
			forecastRows(report).some(
				(r) => r.forecast_seconds !== undefined || r.delivery_error_seconds !== undefined
			)
	);
	/** flai's spread of an error under the filter; under nature and model, from the rows shown. */
	const spread = (which: 'forecast' | 'delivery', field: ErrorField): ErrorSpread | undefined =>
		report
			? spreadFor(
					report,
					which,
					filter,
					forecasts.flatMap((r) => r[field] ?? [])
				)
			: undefined;
	const p = (v: number | undefined) => (v === undefined ? '-' : human(v));
	/** The stories with a forecast and the p50 and p85 of their absolute errors, in a line. */
	const planningSummary = $derived.by(() => {
		if (!report?.forecasts) return '';
		const f = spread('forecast', 'forecast_error_seconds');
		const d = spread('delivery', 'delivery_error_seconds');
		const of = [filter.nature, filter.model].filter(Boolean).join(', ');
		const scope = of ? `(${of})` : 'in the window';
		const n = f?.count ?? 0;
		return [
			`${n} ${n === 1 ? report.type : plural(report.type)} with a forecast ${scope}`,
			`forecast error p50 ${p(f?.p50_seconds)} p85 ${p(f?.p85_seconds)}`,
			`delivery error p50 ${p(d?.p50_seconds)} p85 ${p(d?.p85_seconds)}`
		].join(' · ');
	});
	const codKind = $derived((COD_KINDS as readonly string[]).includes(kind));
	// what the cost of delay charts state under them (S-0213, ADR-0112)
	const unvalued = $derived(report ? withoutValueNow(report) : undefined);
	const ordered = $derived(report ? orderSummary(report) : undefined);
	/** What the host's flai is missing that this cost of delay chart needs; empty when nothing. */
	const codMissing = $derived.by(() => {
		if (!codKind || !report) return '';
		if (!hasCostOfDelay(report)) return 'sends no cost of delay, which this chart is drawn from';
		if (kind === 'cod-order' && !hasOrder(report))
			return 'sends no projection of the pull order, which this chart is drawn from';
		if (kind !== 'cod-order' && !unvalued) return 'does not count the items without a value';
		return '';
	});
	const buckets = $derived(bucketsFor(since));
	const title = $derived(titleOf(kind, report?.usage?.bucket ?? bucket));
	const spend = $derived(report?.usage?.spend?.[report.type]);
	// the items the window holds, which the tables list as the charts plot them (S-0166)
	// $ / Item also lists the items only strategic agents spent on, as it draws them (ADR-0083)
	const spenders = $derived(
		report ? (kind === 'cost' ? withCost : withUsage)(report, epic || undefined) : []
	);
	const completed = $derived(report ? completedIn(report) : []);
	// the stories touches drift draws, done in the window, as it draws them (S-0214)
	const drifted = $derived(report && kind === 'touches-drift' ? driftedIn(report) : []);
	const share = (v: number | undefined) => (v === undefined ? '-' : `${Math.round(v * 100)}%`);
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
	// flai sends a day for every day of the window: none at all is a flai older than S-0205 (S-0216)
	const strategic = $derived(report && strategicKind ? strategicRows(report) : []);
	const ratio = $derived(report && strategicKind ? strategicRatio(report) : undefined);
	/** What the strategic agents spent and worked over the window, and the items completed in it. */
	const strategicTotal = $derived({
		cost: strategic.reduce((n, d) => n + d.cost, 0),
		seconds: strategic.reduce((n, d) => n + d.seconds, 0),
		completed: strategic.reduce((n, d) => n + d.completed, 0)
	});
	const inHours = (seconds: number | undefined) =>
		seconds === undefined ? '-' : `${hours(seconds)}h`;
	/** A strategic agent's day in the table: its cost, starred when estimated, or its hours. */
	const agentDay = (u: StrategicUse | undefined) =>
		!u
			? '-'
			: kind === 'strategic-cost'
				? dollars(u.cost) + (u.estimated ? '*' : '')
				: inHours(u.seconds);

	// the latest question asked: an answer to an earlier one, arriving after it, is not drawn
	let asked = 0;
	// the type last asked for: a chart that reads another type asks flai again when it opens
	let askedType = 'story';
	$effect(() => {
		if (asType !== askedType) untrack(() => void load());
	});
	async function load() {
		const mine = ++asked;
		askedType = asType;
		// an hour is laid out over 31 days or less: a longer window goes by the day
		if (!buckets.includes(bucket)) bucket = 'day';
		const r = await api(`/api/stats?since=${since}&type=${asType}&bucket=${bucket}`);
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

{#each groups as group (group.name)}
	<nav
		class="flex flex-wrap items-center gap-2 {group.name === 'strategic' ? 'mb-3' : 'mb-1'}"
		aria-label="{group.name} charts"
		data-testid="charts-{group.name}"
	>
		<span class="w-16 text-xs text-muted">{group.name}</span>
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
	{#if shown.nature}
		<label
			>nature <select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={nature}
				data-testid="nature"
				><option value="">all</option>{#each facets.natures as n (n)}<option value={n}>{n}</option
					>{/each}</select
			></label
		>
	{/if}
	{#if shown.model}
		<label
			>model <select
				class="rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={model}
				data-testid="model"
				><option value="">all</option>{#each facets.models as m (m)}<option value={m}>{m}</option
					>{/each}</select
			></label
		>
	{/if}
	{#if forecastKind && planningSummary}
		<span class="text-xs text-muted" data-testid="planning-summary">{planningSummary}</span>
	{:else if usageKind && report?.usage && report.usage.items > 0}
		<span class="text-xs text-muted" data-testid="usage-summary">{usageSummary}</span>
	{:else if strategicKind && report && strategic.length > 0}
		<span class="text-xs text-muted" data-testid="strategic-summary"
			>strategic agents {dollars(strategicTotal.cost)} · {inHours(strategicTotal.seconds)} · {strategicTotal.completed}
			{strategicTotal.completed === 1 ? report.type : plural(report.type)} completed</span
		>
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
	{:else if codMissing}
		<p class="mb-2 text-sm text-muted" data-testid="cod-none">
			The flai on the host {codMissing}: it is older than the dashboard. Upgrade it with
			<code>flai self-upgrade</code>.
		</p>
	{:else if usageKind && report && !report.items.some((i) => i.usage)}
		<p class="mb-2 text-sm text-muted" data-testid="usage-none">
			No {report.type} here carries usage yet. flai serve records the tokens and cost of the agents it
			starts; <code>flai serve agent usage --all --write</code> fills in stories worked before.
		</p>
	{:else if forecastKind && report && !report.forecasts}
		<p class="mb-2 text-sm text-muted" data-testid="forecasts-older">
			The flai on the host sends no forecast errors, which this chart is drawn from: it is older
			than the dashboard. Upgrade it with <code>flai self-upgrade</code>.
		</p>
	{:else if kind === 'agent-waiting' && report && !report.waiting}
		<p class="mb-2 text-sm text-muted" data-testid="waiting-older">
			The flai on the host sends no waiting, which this chart is drawn from: it is older than the
			dashboard. Upgrade it with <code>flai self-upgrade</code>.
		</p>
	{:else if forecastKind && !forecastSet}
		<p class="mb-2 text-sm text-muted" data-testid="forecast-none">
			No story done in the window has a forecast. The planner sets one when it plans a story, or set
			one with <code
				>flai edit &lt;story&gt; --forecast-duration 6h --forecast-delivery &lt;UTC time&gt;</code
			>.
		</p>
	{:else if claimsKind && report && !hasClaims(report)}
		<p class="mb-2 text-sm text-muted" data-testid="claims-older">
			The flai on the host sends no held stories or time held, which this chart is drawn from: it is
			older than the dashboard. Upgrade it with <code>flai self-upgrade</code>.
		</p>
	{:else if kind === 'touches-drift' && report?.claims && report.claims.drift === undefined}
		<p class="mb-2 text-sm text-muted" data-testid="drift-none">
			flai could not read git on the host, which touches drift is drawn from: install git there and
			check that the project is a git repository. <code>flai stats</code> says why.
		</p>
	{:else if strategicKind && strategic.length === 0}
		<p class="mb-2 text-sm text-muted" data-testid="strategic-older">
			The flai on the host sends no strategic use per day, which this chart is drawn from: it is
			older than the dashboard. Upgrade it with <code>flai self-upgrade</code>.
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
	{:else if kind === 'strategic-cost'}
		<p class="mb-2 text-xs text-muted">
			* estimated in part: an activity of that agent that day carries an estimated cost.
		</p>
	{/if}
	<Chart {option} theme={t} height={380} />
	{#if kind === 'strategic-cost' && report && strategic.length > 0}
		<p class="mt-2 text-xs text-muted" data-testid="strategic-note">
			The bars are what planning, orchestration, and analysis cost each day. The line is what the
			agents spent per {report.type} completed that day, the strategic agents left out; the dashed line
			is what the strategic agents add to each {report.type} over the window.
			{#if ratio}<span data-testid="strategic-ratio"
					>Over the window they spent {dollars(ratio.cost)} for {ratio.completed}
					{ratio.completed === 1 ? report.type : plural(report.type)} completed: {dollars(
						ratio.cost_per_item
					)} per {report.type}, {Math.round(ratio.share * 100)}% of the {dollars(
						ratio.agent_cost_per_item
					)} the agents spent per {report.type}.</span
				>{:else if strategicTotal.completed === 0}No {plural(report.type)} were completed in the window,
				so there is no ratio.{:else}The agents spent nothing in the window, so there is no ratio.{/if}
		</p>
	{:else if kind === 'strategic-use' && report && strategic.length > 0}
		<p class="mt-2 text-xs text-muted" data-testid="strategic-note">
			The bars are the hours the planner, the orchestrator, and the analyzer worked each day. The
			lines are the mean cycle time and the mean waiting of the {plural(report.type)} completed that day.
			Waiting or cycle time falling while the strategic agents' hours rise is the return on their time.
		</p>
	{/if}
	{#if kind === 'time-in-state' && report}
		<h2 class="mt-6 mb-2 text-base font-medium">Share of lead time per state</h2>
		{#await import('$lib/viz/charts') then m}
			<Chart option={m.stateShare(report, t)} theme={t} height={120} />
		{/await}
	{/if}
	{#if kind === 'agent-waiting' && report?.waiting}
		<h2 class="mt-6 mb-2 text-base font-medium">Longest waits</h2>
		<div class="overflow-x-auto text-xs">
			<WaitTable rows={report.waiting.longest ?? []} type={report.type} />
		</div>
	{/if}
	{#if kind === 'cod-order' && ordered}
		<p class="mt-2 text-sm" data-testid="cod-saving">
			{#if ordered.saving > 0}Ordering {ORDER_LABEL[ordered.cheaper]} would save {amount(
					ordered.saving
				)} on the pull order.{:else}The pull order is already the cheapest of the three.{/if}
		</p>
		<p class="text-xs text-muted" data-testid="cod-totals">
			Projected from now until each ready story is pulled, in the project's currency: {ORDERS.map(
				(by) => `${ORDER_LABEL[by]} ${amount(ordered.totals[by])}`
			).join(' · ')}.
		</p>
		<p class="text-xs text-muted" data-testid="cod-left-out">
			Ready stories left out, without a value or a forecast duration: {ordered.left_out.length
				? ordered.left_out.join(', ')
				: 'none'}.
		</p>
	{/if}
	{#if codKind && unvalued}
		<p class="mt-2 text-xs text-muted" data-testid="cod-without-value">
			{#if COD_COLUMNS.some((col) => unvalued[col] > 0)}Items without a value now, which add nothing
				to this chart: {COD_COLUMNS.map((col) => `${col} ${unvalued[col]}`).join(
					' · '
				)}.{:else}Items without a value now: none.{/if}
		</p>
	{/if}
	<details class="mt-4 text-xs">
		<summary class="cursor-pointer text-muted">table view</summary>
		<div class="mt-2 overflow-x-auto">
			{#if forecastKind}
				<ForecastTable rows={forecasts} />
			{:else if kind === 'cod-outstanding' && report}
				<table class="min-w-full" data-testid="cod-day-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">date</th>{#each COD_COLUMNS as col (col)}<th class="pr-4">{col}</th
								>{/each}<th>without value</th></tr
						></thead
					><tbody
						>{#each codDayRows(report) as d (d.date)}<tr
								><td class="pr-4 font-mono">{d.date}</td>{#each COD_COLUMNS as col (col)}<td
										class="pr-4">{amount(d.outstanding[col] ?? 0)}</td
									>{/each}<td
									>{d.without_value
										? COD_COLUMNS.reduce((n, col) => n + (d.without_value?.[col] ?? 0), 0)
										: '-'}</td
								></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'cod-incurred' && report}
				<table class="min-w-full" data-testid="cod-week-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">week</th><th class="pr-4">start</th><th class="pr-4">incurred</th
							><th>mean per week</th></tr
						></thead
					><tbody
						>{#each codWeekRows(report) as w (w.week)}<tr
								><td class="pr-4 font-mono">{w.week}</td><td class="pr-4 font-mono">{w.start}</td
								><td class="pr-4">{amount(w.incurred)}</td><td>{amount(w.mean)}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'cod-order' && report}
				<table class="min-w-full" data-testid="cod-order-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">order</th><th class="pr-4">pulled</th><th class="pr-4">story</th><th
								>cumulative incurred</th
							></tr
						></thead
					><tbody
						>{#each codOrderRows(report) as p (`${p.by} ${p.id}`)}<tr
								><td class="pr-4">{ORDER_LABEL[p.by]}</td><td
									class="pr-4 font-mono whitespace-nowrap">{bucketLabel(p.at, 'hour')}</td
								><td class="pr-4 font-mono">{p.id}</td><td>{amount(p.incurred)}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'parallelism' && report}
				<table class="min-w-full" data-testid="parallelism-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">day</th><th class="pr-4">in progress</th><th class="pr-4">held</th
							><th>limit</th></tr
						></thead
					><tbody
						>{#each report.claims?.days ?? [] as d (d.date)}<tr
								><td class="pr-4 font-mono">{d.date}</td><td class="pr-4">{d.in_progress}</td><td
									class="pr-4">{d.held ?? '-'}</td
								><td>{report.claims?.limit ?? '-'}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'hold-time' && report}
				<table class="min-w-full" data-testid="hold-time-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">week</th><th class="pr-4">start</th><th class="pr-4">overlap</th><th
								class="pr-4">after</th
							><th class="pr-4">empty claim</th><th>held</th></tr
						></thead
					><tbody
						>{#each report.claims?.weeks ?? [] as w (w.week)}<tr
								><td class="pr-4 font-mono">{w.week}</td><td class="pr-4 font-mono">{w.start}</td
								><td class="pr-4">{human(w.held_seconds.overlap)}</td><td class="pr-4"
									>{human(w.held_seconds.after)}</td
								><td class="pr-4">{human(w.held_seconds['no-touches'])}</td><td
									>{human(
										w.held_seconds.overlap + w.held_seconds.after + w.held_seconds['no-touches']
									)}</td
								></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'touches-drift' && report}
				<table class="min-w-full" data-testid="drift-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">story</th><th class="pr-4">completed</th><th class="pr-4"
								>outside its touches</th
							><th>touches unchanged</th></tr
						></thead
					><tbody
						>{#each drifted as d (d.id)}<tr
								><td class="pr-4"
									><a class="font-mono underline" href={resolve('/items/[id]', { id: d.id })}
										>{d.id}</a
									>
									{d.title}</td
								><td class="pr-4 font-mono">{localDate(d.completed)}</td><td class="pr-4"
									>{d.outside_count}</td
								><td>{d.unchanged_count}</td></tr
							>{/each}</tbody
					>
				</table>
				<table class="mt-4 min-w-full" data-testid="exact-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">week</th><th class="pr-4">start</th><th class="pr-4">stories</th><th
								class="pr-4">exact</th
							><th>share</th></tr
						></thead
					><tbody
						>{#each report.claims?.weeks ?? [] as w (w.week)}<tr
								><td class="pr-4 font-mono">{w.week}</td><td class="pr-4 font-mono">{w.start}</td
								><td class="pr-4">{w.stories ?? '-'}</td><td class="pr-4">{w.exact ?? '-'}</td><td
									>{share(w.exact_share)}</td
								></tr
							>{/each}</tbody
					>
				</table>
			{:else if spendKind && report}
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
										>{i.completed ? localDate(i.completed) : '-'}</td
									><td class="pr-4">{m.model}</td><td class="pr-4">{count(m.tokens)}</td><td
										class="pr-4">{perMinute(m) !== undefined ? count(perMinute(m)!) : '-'}</td
									><td>{dollars(m.cost)}{i.usage?.estimated ? '*' : ''}</td></tr
								>{/each}{#if kind === 'cost' && strategicCost(i) > 0}<tr
									data-testid="usage-strategic"
									><td class="pr-4 font-mono">{i.id}</td><td class="pr-4"
										>{i.completed ? localDate(i.completed) : '-'}</td
									><td class="pr-4">strategic</td><td class="pr-4"
										>{count((i.usage?.strategic ?? []).reduce((n, x) => n + x.tokens, 0))}</td
									><td class="pr-4">-</td><td>{dollars(strategicCost(i))}*</td></tr
								>{/if}{/each}</tbody
					>
				</table>
			{:else if strategicKind && report}
				<table class="min-w-full" data-testid="strategic-table">
					<thead
						><tr class="text-left text-muted"
							><th class="pr-4">day</th>{#each STRATEGIC_AGENTS as a (a)}<th class="pr-4">{a}</th
								>{/each}<th class="pr-4">total</th><th class="pr-4">completed</th
							>{#if kind === 'strategic-cost'}<th>mean per {report.type}</th>{:else}<th class="pr-4"
									>mean cycle time</th
								><th>mean waiting</th>{/if}</tr
						></thead
					><tbody
						>{#each strategic as d (d.date)}<tr
								><td class="pr-4 font-mono">{d.date}</td>{#each STRATEGIC_AGENTS as a (a)}<td
										class="pr-4">{agentDay(d.agents[a])}</td
									>{/each}<td class="pr-4"
									>{kind === 'strategic-cost' ? dollars(d.cost) : inHours(d.seconds)}</td
								><td class="pr-4">{d.completed}</td>{#if kind === 'strategic-cost'}<td
										>{d.cost_per_item !== undefined ? dollars(d.cost_per_item) : '-'}</td
									>{:else}<td class="pr-4">{inHours(d.cycle_time_seconds)}</td><td
										>{inHours(d.wait_seconds)}</td
									>{/if}</tr
							>{/each}</tbody
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
								><td class="pr-4">{i.completed ? localDate(i.completed) : ''}</td><td class="pr-4"
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
