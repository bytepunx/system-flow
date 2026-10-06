<script lang="ts">
	// One block of the strategic agents' settings in system-flow.yaml (S-0229): the orchestrator's,
	// the planner's, or the analyzer's, as settings.get's strategic block lists them. Each key shows
	// what it does, and a permission what can go wrong while it is on, with an input by its kind and
	// its default while it is unset. Save writes only the keys that changed through settings.manifest,
	// an emptied one as an unset; flai's refusal keeps what was typed and says each reason under its
	// field. While the settings host action is off every input is disabled and the panel says why.
	import { untrack } from 'svelte';
	import { api } from '$lib/api';
	import {
		inBlock,
		strategicChange,
		strategicForm,
		unchanged,
		type Refusal,
		type StrategicBlock,
		type StrategicForm,
		type StrategicSetting,
		type StrategicSettings
	} from '$lib/settings';

	let {
		block,
		strategic,
		onsaved
	}: {
		block: StrategicBlock;
		/** settings.get's host.strategic. */
		strategic: StrategicSettings;
		/** Called after a save that flai took, so that the page reads the settings again. */
		onsaved?: () => void;
	} = $props();

	const TITLE: Record<StrategicBlock, string> = {
		orchestration: "The orchestrator's settings",
		planning: "The planner's settings",
		analysis: "The analyzer's settings"
	};

	const uid = $props.id();
	const rows = $derived(inBlock(strategic, block));
	const shown = $derived(new Set(rows.map((s) => s.key)));

	let before = $state<StrategicForm>(untrack(() => strategicForm(rows)));
	let form = $state<StrategicForm>(untrack(() => strategicForm(rows)));
	let refused = $state<Refusal[]>([]);
	let said = $state<{ ok: boolean; text: string; enable?: string } | null>(null);
	let busy = $state(false);

	// The settings read again replace what the form holds.
	$effect(() => {
		const read = rows;
		untrack(() => {
			before = strategicForm(read);
			form = strategicForm(read);
			refused = [];
		});
	});

	const change = $derived(strategicChange(rows, before, form));
	const dirty = $derived('refused' in change || !unchanged(change));
	const off = $derived(!strategic.editable);
	/** Reasons whose field is not one of the inputs here, said at the top. */
	const elsewhere = $derived(refused.filter((r) => !shown.has(r.field)));

	const name = (key: string) => key.split('.').at(-1) ?? key;
	const group = (key: string) => {
		const parts = key.split('.');
		return parts.length > 2 ? parts[1] : '';
	};
	const id = (key: string) => `${uid}-${key}`;
	const reasons = (key: string) => refused.filter((r) => r.field === key);
	const described = (s: StrategicSetting) =>
		[
			`${id(s.key)}-meaning`,
			s.risk ? `${id(s.key)}-risk` : '',
			reasons(s.key).length ? `${id(s.key)}-refused` : ''
		]
			.filter(Boolean)
			.join(' ');

	async function save() {
		if ('refused' in change) {
			refused = change.refused;
			said = { ok: false, text: 'Not saved: correct what is marked below.' };
			return;
		}
		if (unchanged(change)) return;
		const params: Record<string, unknown> = { kind: 'manifest' };
		if (Object.keys(change.set).length) params.set = change.set;
		if (change.unset.length) params.unset = change.unset;
		busy = true;
		try {
			const r = await api('/api/settings', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify(params)
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				refused = Array.isArray(body.refused) ? body.refused : [];
				said = {
					ok: false,
					text: refused.length
						? 'Not saved: flai refused what is marked.'
						: (body.error ?? r.statusText),
					enable: typeof body.enable === 'string' ? body.enable : undefined
				};
				return;
			}
			refused = [];
			before = $state.snapshot(form);
			said = {
				ok: true,
				text: body.commit ? `saved in ${String(body.commit).slice(0, 7)}` : 'saved'
			};
			onsaved?.();
		} catch (e) {
			said = { ok: false, text: e instanceof Error ? e.message : String(e) };
		} finally {
			busy = false;
		}
	}
</script>

{#snippet unsetNote(s: StrategicSetting)}
	{#if !s.set}
		<span class="text-xs text-muted" data-testid="default-{s.key}">
			{#if s.kind === 'agent'}
				unset: the project's agent
			{:else if s.default}
				unset: default <code>{s.default}</code>
			{:else}
				unset
			{/if}
		</span>
	{/if}
{/snippet}

{#snippet about(s: StrategicSetting)}
	<div class="mt-1 grid grid-cols-1 gap-x-4 gap-y-1 text-xs sm:grid-cols-2">
		<p id="{id(s.key)}-meaning" class="text-muted">{s.meaning}</p>
		{#if s.risk}
			<p id="{id(s.key)}-risk" class="text-warn" data-testid="risk-{s.key}">Risk: {s.risk}</p>
		{/if}
	</div>
	{#if reasons(s.key).length}
		<div id="{id(s.key)}-refused" role="alert" data-testid="refused-{s.key}">
			{#each reasons(s.key) as r, i (i)}
				<p class="text-xs text-danger">{r.reason}</p>
			{/each}
		</div>
	{/if}
{/snippet}

<section class="space-y-3 text-sm" data-testid="strategic-{block}">
	<h2 class="font-medium">{TITLE[block]}</h2>
	{#if off}
		<p
			class="rounded border border-line-strong bg-raised p-3 text-xs"
			data-testid="strategic-readonly"
		>
			Read-only: the settings host action is off for this project, so the dashboard may not change
			system-flow.yaml. On the host, in the project, run
			<code class="rounded bg-surface px-1">{strategic.enable}</code> to change them here.
		</p>
	{/if}
	{#if elsewhere.length}
		<div
			class="rounded border border-danger bg-danger-soft p-2 text-xs text-danger"
			role="alert"
			data-testid="strategic-refused"
		>
			{#each elsewhere as r, i (i)}
				<p><code>{r.field}</code>: {r.reason}</p>
			{/each}
		</div>
	{/if}

	<ul class="space-y-3">
		{#each rows as s, i (s.key)}
			{#if group(s.key) && group(s.key) !== group(rows[i - 1]?.key ?? '')}
				<li><h3 class="mt-2 text-xs font-medium text-muted">{group(s.key)}</h3></li>
			{/if}
			<li class="rounded border border-line p-2" data-testid="row-{s.key}">
				{#if s.kind === 'agent'}
					<fieldset aria-describedby={described(s)}>
						<legend class="flex flex-wrap items-baseline gap-2">
							<span class="font-mono">{name(s.key)}</span>{@render unsetNote(s)}
						</legend>
						<div class="mt-1 grid grid-cols-1 gap-3 sm:grid-cols-3">
							<label class="block">
								<span class="mb-1 block text-xs text-muted">harness</span>
								<input
									class="w-full rounded border border-line-strong bg-surface px-2 py-1 disabled:opacity-50"
									bind:value={form.agents[s.key].harness}
									placeholder="the project's"
									disabled={off || busy}
									aria-invalid={reasons(s.key).length > 0}
									data-testid="input-{s.key}-harness"
								/>
							</label>
							<label class="block">
								<span class="mb-1 block text-xs text-muted">model</span>
								<input
									class="w-full rounded border border-line-strong bg-surface px-2 py-1 disabled:opacity-50"
									bind:value={form.agents[s.key].model}
									placeholder="the project's"
									disabled={off || busy}
									aria-invalid={reasons(s.key).length > 0}
									data-testid="input-{s.key}-model"
								/>
							</label>
							<label class="block">
								<span class="mb-1 block text-xs text-muted">config, one key=value a line</span>
								<textarea
									class="h-16 w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs disabled:opacity-50"
									bind:value={form.agents[s.key].config}
									placeholder="effort=high"
									disabled={off || busy}
									aria-invalid={reasons(s.key).length > 0}
									data-testid="input-{s.key}-config"></textarea>
							</label>
						</div>
					</fieldset>
				{:else}
					<div class="flex flex-wrap items-center gap-3">
						<label class="font-mono" for={id(s.key)}>{name(s.key)}</label>
						{#if s.kind === 'boolean'}
							<input
								id={id(s.key)}
								type="checkbox"
								role="switch"
								checked={form.values[s.key] === 'true'}
								disabled={off || busy}
								aria-describedby={described(s)}
								aria-invalid={reasons(s.key).length > 0}
								data-testid="input-{s.key}"
								onchange={(e) => (form.values[s.key] = String(e.currentTarget.checked))}
							/>
						{:else if s.kind === 'choice'}
							<select
								id={id(s.key)}
								class="rounded border border-line-strong bg-surface px-2 py-1 disabled:opacity-50"
								bind:value={form.values[s.key]}
								disabled={off || busy}
								aria-describedby={described(s)}
								aria-invalid={reasons(s.key).length > 0}
								data-testid="input-{s.key}"
							>
								<option value="">{s.default ? `default (${s.default})` : 'unset'}</option>
								{#each s.values ?? [] as v (v)}
									<option value={v}>{v}</option>
								{/each}
							</select>
						{:else if s.kind === 'number'}
							<input
								id={id(s.key)}
								type="number"
								min="0"
								step={s.whole ? 1 : 'any'}
								inputmode={s.whole ? 'numeric' : 'decimal'}
								class="w-32 rounded border border-line-strong bg-surface px-2 py-1 disabled:opacity-50"
								value={form.values[s.key]}
								placeholder={s.default ?? ''}
								disabled={off || busy}
								aria-describedby={described(s)}
								aria-invalid={reasons(s.key).length > 0}
								data-testid="input-{s.key}"
								oninput={(e) => (form.values[s.key] = e.currentTarget.value)}
							/>
						{:else}
							<input
								id={id(s.key)}
								class="min-w-48 flex-1 rounded border border-line-strong bg-surface px-2 py-1 disabled:opacity-50 {s.kind ===
								'text'
									? ''
									: 'font-mono text-xs'}"
								bind:value={form.values[s.key]}
								placeholder={s.default ??
									(s.kind === 'cron' ? '0 6 * * 1, or daily' : s.kind === 'duration' ? '168h' : '')}
								disabled={off || busy}
								aria-describedby={described(s)}
								aria-invalid={reasons(s.key).length > 0}
								data-testid="input-{s.key}"
							/>
						{/if}
						{@render unsetNote(s)}
					</div>
				{/if}
				{@render about(s)}
			</li>
		{/each}
	</ul>

	<div class="flex flex-wrap items-center gap-3">
		<button
			class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
			disabled={off || busy || !dirty}
			data-testid="strategic-save"
			onclick={save}>Save</button
		>
		<span class="text-xs text-muted">
			Saved to system-flow.yaml and committed; an emptied field is removed, so its default applies.
		</span>
	</div>
	{#if said}
		<p
			class="text-xs {said.ok ? 'text-good' : 'text-danger'}"
			role={said.ok ? 'status' : 'alert'}
			data-testid="strategic-said"
		>
			{said.text}
			{#if said.enable}
				On the host, in the project, run <code class="rounded bg-surface px-1">{said.enable}</code>.
			{/if}
		</p>
	{/if}
</section>
