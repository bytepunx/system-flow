<script lang="ts">
	import { api } from '$lib/api';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import Chart from '$lib/components/Chart.svelte';
	import { build, human, KINDS, TITLES, type Kind, type Report } from '$lib/viz/charts';
	import { theme } from '$lib/viz/palette';

	let since = $state('30d');
	let type = $state('story');
	let epic = $state('');
	let report = $state<Report | null>(null);
	let error = $state<string | null>(null);
	let dark = $state(false);
	let epics = $state<{ id: string; title: string }[]>([]);

	const kind = $derived(
		(KINDS as readonly string[]).includes(page.params.kind ?? '')
			? (page.params.kind as Kind)
			: 'cycle-time'
	);
	const t = $derived(theme(dark));
	const option = $derived(report ? build(kind, report, t, epic || undefined) : null);
	const s = $derived(report?.summary);

	async function load() {
		error = null;
		const r = await api(`/api/stats?since=${since}&type=${type}`);
		if (!r.ok) {
			error = (await r.json()).error ?? r.statusText;
			report = null;
			return;
		}
		report = await r.json();
	}
	onMount(() => {
		const mq = matchMedia('(prefers-color-scheme: dark)');
		dark = mq.matches;
		const onScheme = (e: MediaQueryListEvent) => (dark = e.matches);
		mq.addEventListener('change', onScheme);
		const es = new EventSource('/api/events');
		es.addEventListener('change', () => load());
		(async () => {
			epics = (await (await api('/api/items?type=epic')).json()).map(
				(e: { id: string; title: string }) => ({ id: e.id, title: e.title })
			);
			await load();
		})();
		return () => {
			es.close();
			mq.removeEventListener('change', onScheme);
		};
	});
</script>

<svelte:head><title>{TITLES[kind]} · flaiover</title></svelte:head>

<div class="mb-3 flex flex-wrap items-center gap-2">
	{#each KINDS as k (k)}
		<a
			href={resolve('/charts/[kind]', { kind: k })}
			class="rounded px-2 py-1 text-sm {k === kind
				? 'bg-zinc-200 font-medium dark:bg-zinc-800'
				: 'text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100'}"
			>{TITLES[k]}</a
		>
	{/each}
</div>
<div class="mb-4 flex flex-wrap items-center gap-3 text-sm">
	<label
		>window <select
			class="rounded border border-zinc-300 bg-white px-2 py-1 dark:border-zinc-700 dark:bg-zinc-900"
			bind:value={since}
			onchange={load}
			>{#each ['7d', '30d', '90d', '365d'] as w (w)}<option value={w}>{w}</option>{/each}</select
		></label
	>
	<label
		>type <select
			class="rounded border border-zinc-300 bg-white px-2 py-1 dark:border-zinc-700 dark:bg-zinc-900"
			bind:value={type}
			onchange={load}
			>{#each ['story', 'task', 'epic'] as ty (ty)}<option value={ty}>{ty}</option>{/each}</select
		></label
	>
	{#if kind !== 'cfd' && kind !== 'throughput' && kind !== 'estimates'}
		<label
			>epic <select
				class="rounded border border-zinc-300 bg-white px-2 py-1 dark:border-zinc-700 dark:bg-zinc-900"
				bind:value={epic}
				><option value="">all</option>{#each epics as e (e.id)}<option value={e.id}
						>{e.id} {e.title}</option
					>{/each}</select
			></label
		>
	{/if}
	{#if s}
		<span class="text-xs text-zinc-500"
			>completed {s.completed} · cancelled {s.cancelled} · WIP {s.wip} · throughput {s.throughput_per_week.toFixed(
				1
			)}/week · cycle p50 {human(s.cycle_time.p50_seconds)} p85 {human(
				s.cycle_time.p85_seconds
			)}</span
		>
	{/if}
</div>
{#if error}
	<p class="rounded border border-red-300 bg-red-50 p-3 text-sm text-red-800">{error}</p>
{:else if option}
	<h1 class="mb-2 text-xl font-semibold">{TITLES[kind]}</h1>
	<Chart
		{option}
		theme={t}
		height={kind === 'aging' ? Math.max(200, 32 * (report?.aging.length ?? 0) + 80) : 380}
	/>
	{#if kind === 'time-in-state' && report}
		<h2 class="mt-6 mb-2 text-base font-medium">Share of lead time per state</h2>
		{#await import('$lib/viz/charts') then m}
			<Chart option={m.stateShare(report, t)} theme={t} height={120} />
		{/await}
	{/if}
	<details class="mt-4 text-xs">
		<summary class="cursor-pointer text-zinc-500">table view</summary>
		<div class="mt-2 overflow-x-auto">
			{#if kind === 'aging' && report}
				<table class="min-w-full">
					<thead
						><tr class="text-left text-zinc-500"
							><th class="pr-4">item</th><th class="pr-4">status</th><th class="pr-4">age</th><th
								>title</th
							></tr
						></thead
					><tbody
						>{#each report.aging as a (a.id)}<tr
								><td class="pr-4 font-mono">{a.id}</td><td class="pr-4">{a.status}</td><td
									class="pr-4">{a.age}{a.over_p85 ? ' (over p85)' : ''}</td
								><td>{a.title}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if kind === 'throughput' && report}
				<table class="min-w-full">
					<thead
						><tr class="text-left text-zinc-500"><th class="pr-4">week</th><th>done</th></tr></thead
					><tbody
						>{#each report.throughput as w (w.week)}<tr
								><td class="pr-4 font-mono">{w.week}</td><td>{w.done}</td></tr
							>{/each}</tbody
					>
				</table>
			{:else if report}
				<table class="min-w-full">
					<thead
						><tr class="text-left text-zinc-500"
							><th class="pr-4">item</th><th class="pr-4">nature</th><th class="pr-4">completed</th
							><th class="pr-4">cycle</th><th class="pr-4">lead</th><th>blocked</th></tr
						></thead
					><tbody
						>{#each report.items.filter((i) => i.completed) as i (i.id)}<tr
								><td class="pr-4 font-mono">{i.id}</td><td class="pr-4">{i.nature}</td><td
									class="pr-4">{i.completed?.slice(0, 10)}</td
								><td class="pr-4"
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
	<p class="text-sm text-zinc-500">Loading…</p>
{/if}
