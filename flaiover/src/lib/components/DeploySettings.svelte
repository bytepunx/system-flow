<script lang="ts">
	// The kluster: block's deploy setup (S-0586, S-0595, ADR-0263), as the Klusters page's step 4 and
	// the Deployments page's Setup pane show it: the default namespace, what builds the images (one
	// build command, or one per component, since a block takes one), the deploy and teardown commands,
	// and the timeout. Each is saved through settings.manifest, which runs flai manifest set on the host
	// and needs the settings host action; flai's refusals are shown beside the field they name. A
	// command is typed one argument a line. Write these with an agent makes a backlog story for an
	// agent to write the project's scripts. After a save the component reads kluster.status again and
	// gives it to saved, so the page shows what flai now says.
	import type { KlusterSetting, KlusterStatus } from '$lib/server/repo';
	import {
		argv,
		componentList,
		componentRows,
		lines,
		sameComponents,
		type BuildComponent,
		type ComponentRow,
		type Refusal
	} from '$lib/settings';
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';

	let {
		settings,
		settingsOn,
		settingsEnable = '',
		deployable,
		notice = true,
		saved
	}: {
		/** kluster.status's settings: the block's keys, each with its kind, meaning, and value. */
		settings: KlusterSetting[];
		/** Whether the settings host action is on, which the block's writes need. */
		settingsOn: boolean;
		/** What the operator runs on the host to turn settings on. */
		settingsEnable?: string;
		/** Whether the block names a deploy command. */
		deployable: boolean;
		/** Whether to say here that settings is off; a page that says it once elsewhere turns it off. */
		notice?: boolean;
		/** Called with kluster.status as flai answers it after a save. */
		saved?: (status: KlusterStatus) => void | Promise<void>;
	} = $props();

	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';
	const field = 'w-full rounded border border-line bg-surface px-2 py-1 font-mono text-sm';
	// The keys saved as they are typed; the build is kluster.build or kluster.components by the mode.
	const PLAIN_KEYS = ['kluster.namespace', 'kluster.deploy', 'kluster.teardown', 'kluster.timeout'];
	const AGENT_TITLE =
		"Write the project's kluster build, deploy, and teardown scripts and name them in the kluster: block";

	type Said = { ok: boolean; text: string };
	type Epic = { id: string; title: string; status: string };
	type Mode = 'one' | 'components';

	let form = $state<Record<string, string>>({});
	let before = $state<Record<string, string>>({});
	let rows = $state<ComponentRow[]>([]);
	let beforeComponents = $state<BuildComponent[]>([]);
	let mode = $state<Mode>('one');
	let refused = $state<Refusal[]>([]);
	let saveSaid = $state<Said | null>(null);
	let busy = $state(false);

	let epics = $state<Epic[]>([]);
	let epic = $state('');
	let made = $state<string | null>(null);
	let storySaid = $state<Said | null>(null);

	const settingOf = (key: string) => settings.find((s) => s.key === key);
	// A flai before ADR-0263 lists no kluster.components: then the build is one command alone.
	const hasComponents = $derived(!!settingOf('kluster.components'));
	const disabled = $derived(!settingsOn || busy);
	const refusalOf = (key: string) => refused.find((r) => r.field === key)?.reason;
	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));
	const typed = (s: KlusterSetting) =>
		s.kind === 'command'
			? lines(Array.isArray(s.value) ? (s.value as string[]) : undefined)
			: typeof s.value === 'string'
				? s.value
				: '';

	// The form follows the settings flai answers: read again whenever they change, as after a save.
	let shown = '';
	$effect(() => {
		const now = JSON.stringify(settings);
		if (now === shown) return;
		shown = now;
		const f: Record<string, string> = {};
		for (const s of settings) if (s.kind !== 'components') f[s.key] = typed(s);
		form = f;
		before = { ...f };
		const comps = settings.find((s) => s.key === 'kluster.components');
		const next = componentRows(comps?.value);
		rows = next;
		beforeComponents = componentList(next);
		mode = next.length ? 'components' : 'one';
	});

	function addRow() {
		rows = [...rows, { name: '', build: '' }];
	}
	function removeRow(i: number) {
		rows = rows.filter((_, j) => j !== i);
	}
	function moveRow(i: number, by: -1 | 1) {
		const j = i + by;
		if (j < 0 || j >= rows.length) return;
		const next = [...rows];
		[next[i], next[j]] = [next[j], next[i]];
		rows = next;
	}

	/** What settings.manifest is asked: only what changed, an emptied key unset. */
	function change(): {
		set: Record<string, string | string[] | BuildComponent[]>;
		unset: string[];
	} {
		const set: Record<string, string | string[] | BuildComponent[]> = {};
		const unset: string[] = [];
		const plain = (key: string) => {
			const s = settingOf(key);
			if (!s || (form[key] ?? '') === (before[key] ?? '')) return;
			if (s.kind === 'command') {
				const args = argv(form[key] ?? '');
				if (args.length) set[key] = args;
				else unset.push(key);
			} else {
				const v = (form[key] ?? '').trim();
				if (v) set[key] = v;
				else unset.push(key);
			}
		};
		plain('kluster.namespace');
		if (mode === 'one') {
			plain('kluster.build');
			// one build command in place of components: flai takes one or the other
			if (beforeComponents.length) unset.push('kluster.components');
		} else {
			const list = componentList(rows);
			if (!sameComponents(list, beforeComponents)) {
				if (list.length) set['kluster.components'] = list;
				else unset.push('kluster.components');
			}
			if ((before['kluster.build'] ?? '').trim()) unset.push('kluster.build');
		}
		for (const key of PLAIN_KEYS.slice(1)) plain(key);
		return { set, unset };
	}

	/** Saves what changed through settings.manifest, then reads kluster.status again for saved. */
	async function save() {
		const { set, unset } = change();
		if (!Object.keys(set).length && !unset.length) {
			saveSaid = { ok: true, text: 'Nothing changed.' };
			return;
		}
		const params: Record<string, unknown> = { kind: 'manifest' };
		if (Object.keys(set).length) params.set = set;
		if (unset.length) params.unset = unset;
		busy = true;
		saveSaid = null;
		try {
			const r = await api('/api/settings', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify(params)
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				refused = Array.isArray(body.refused) ? body.refused : [];
				saveSaid = {
					ok: false,
					text: refused.length
						? 'Not saved: flai refused what is marked.'
						: (body.error ?? r.statusText)
				};
				return;
			}
			refused = [];
			saveSaid = {
				ok: true,
				text: body.commit ? `Saved in ${String(body.commit).slice(0, 7)}.` : 'Saved.'
			};
			const s = await api('/api/kluster');
			if (s.ok) await saved?.((await s.json()) as KlusterStatus);
		} catch (e) {
			saveSaid = { ok: false, text: errorText(e) };
		} finally {
			busy = false;
		}
	}

	/** Reads the open epics a story for the scripts can go under. */
	async function loadEpics() {
		try {
			const r = await api('/api/items?type=epic&archived=false');
			if (!r.ok) return;
			const all = (await r.json()) as Epic[];
			epics = all.filter((e) => e.status !== 'done' && e.status !== 'cancelled');
			if (!epic && epics.length) epic = epics[0].id;
		} catch {
			epics = [];
		}
	}

	/** Makes a backlog story for an agent to write the project's scripts, through item.new. */
	async function writeWithAgent() {
		busy = true;
		storySaid = null;
		const named = componentList(rows)
			.map((c) => c.name)
			.filter(Boolean);
		const goal =
			"Give this project's kluster: block the commands that build its images, deploy a build into a namespace, and tear a namespace down, as design/system/kluster.md describes, so that flai deploys each story in review into its own namespace." +
			(named.length
				? ` Its components are ${named.join(', ')}: each gets a build of its own.`
				: '');
		const criteria = [
			'Scripts under scripts/ build the images, deploy into FLAI_KLUSTER_NAMESPACE, and tear it down',
			named.length
				? 'The kluster: block names a build for each component in components, and the deploy and teardown'
				: 'The kluster: block names them as build, or one build per component in components, and deploy and teardown',
			'A story in review deploys into its namespace and its teardown removes it'
		]
			.map((c) => `- [ ] ${c}`)
			.join('\n');
		try {
			const r = await api('/api/items', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					type: 'story',
					title: AGENT_TITLE,
					nature: 'feature',
					parent: epic,
					touches: ['system-flow.yaml', 'scripts'],
					body: `## Goal\n\n${goal}\n\n## Acceptance criteria\n\n${criteria}\n\n## Tasks\n\n## Notes\n`
				})
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				storySaid = { ok: false, text: `The story was not created: ${body.error ?? r.statusText}` };
				return;
			}
			made = (body as { item: { id: string } }).item.id;
		} catch (e) {
			storySaid = { ok: false, text: errorText(e) };
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void loadEpics();
	});
</script>

{#snippet said(s: Said | null | undefined, id: string)}
	{#if s}
		<p
			class="mt-2 text-sm {s.ok ? 'text-good' : 'text-danger'}"
			role={s.ok ? 'status' : 'alert'}
			data-testid={id}
		>
			{s.text}
		</p>
	{/if}
{/snippet}

{#snippet settingField(key: string)}
	{@const s = settingOf(key)}
	{#if s}
		<label class="mt-3 block text-sm" data-testid="field-{key}">
			<span class="font-mono">{key.replace('kluster.', '')}</span>
			{#if s.kind === 'command'}
				<span class="text-xs text-muted">: one argument a line, the program first</span>
			{/if}
			<span class="block text-xs text-muted">{s.meaning}</span>
			{#if s.kind === 'command'}
				<textarea
					class={field}
					rows="2"
					{disabled}
					placeholder={s.default ?? ''}
					bind:value={form[key]}></textarea>
			{:else}
				<input class={field} {disabled} placeholder={s.default ?? ''} bind:value={form[key]} />
			{/if}
			{#if refusalOf(key)}
				<span class="block text-xs text-danger" data-testid="refused-{key}">{refusalOf(key)}</span>
			{/if}
		</label>
	{/if}
{/snippet}

<div data-testid="deploy-settings">
	{#if !settingsOn && notice && settingsEnable}
		<p
			class="mt-1 rounded border border-warn bg-warn-soft p-2 text-sm text-warn"
			role="status"
			data-testid="deploy-settings-off"
		>
			The settings host action is off, so these can be read but not saved. On the host, in the
			project, run <code class="font-mono">{settingsEnable}</code>, then reload this page.
		</p>
	{/if}
	{#if !deployable}
		<p class="mt-1 text-sm text-warn" data-testid="kluster-no-deploy">
			The cluster can come up and go down, but no story is deployed into it until the deploy command
			is set.
		</p>
	{/if}

	{@render settingField('kluster.namespace')}

	<fieldset class="mt-3 text-sm" data-testid="kluster-build">
		<legend class="font-mono">build</legend>
		{#if hasComponents}
			<div class="mt-1 flex flex-wrap gap-4" data-testid="kluster-build-mode">
				<label class="flex items-center gap-1">
					<input type="radio" value="one" {disabled} bind:group={mode} />
					One build command
				</label>
				<label class="flex items-center gap-1">
					<input type="radio" value="components" {disabled} bind:group={mode} />
					One build per component
				</label>
			</div>
		{/if}
		{#if mode === 'one'}
			{@render settingField('kluster.build')}
		{:else}
			{@const s = settingOf('kluster.components')}
			<p class="mt-2 text-xs text-muted">{s?.meaning}</p>
			<ol class="mt-2 space-y-3" data-testid="kluster-components">
				{#each rows as row, i (i)}
					<li class="rounded border border-line p-2" data-testid="component-row">
						<div class="flex flex-wrap items-end gap-2">
							<label class="grow text-sm">
								<span class="block text-xs text-muted">Name, a DNS label such as api</span>
								<input
									class={field}
									{disabled}
									bind:value={row.name}
									data-testid="component-name"
								/>
							</label>
							<button
								type="button"
								class="{button} border border-line"
								disabled={disabled || i === 0}
								aria-label="Move {row.name || 'this component'} up"
								onclick={() => moveRow(i, -1)}>Up</button
							>
							<button
								type="button"
								class="{button} border border-line"
								disabled={disabled || i === rows.length - 1}
								aria-label="Move {row.name || 'this component'} down"
								onclick={() => moveRow(i, 1)}>Down</button
							>
							<button
								type="button"
								class="{button} border border-line text-danger"
								{disabled}
								onclick={() => removeRow(i)}>Remove</button
							>
						</div>
						<label class="mt-2 block text-sm">
							<span class="block text-xs text-muted"
								>Its build command, one argument a line, the program first</span
							>
							<textarea
								class={field}
								rows="2"
								{disabled}
								bind:value={row.build}
								data-testid="component-build"></textarea>
						</label>
					</li>
				{/each}
			</ol>
			<button type="button" class="{button} mt-2 border border-line" {disabled} onclick={addRow}
				>Add a component</button
			>
			{#if refusalOf('kluster.components')}
				<span class="block text-xs text-danger" data-testid="refused-kluster.components"
					>{refusalOf('kluster.components')}</span
				>
			{/if}
		{/if}
	</fieldset>

	{#each PLAIN_KEYS.slice(1) as key (key)}
		{@render settingField(key)}
	{/each}
	<button type="button" class="{button} mt-3 bg-primary text-on-primary" {disabled} onclick={save}
		>Save</button
	>
	{@render said(saveSaid, 'kluster-save-deploy')}

	<div class="mt-4 rounded border border-line p-3 text-sm" data-testid="kluster-agent">
		<p>
			The build, deploy, and teardown scripts are the project's code. An agent can write them on a
			story, for you to review and accept.
		</p>
		{#if made}
			<p class="mt-2" data-testid="kluster-agent-made">
				Made
				<a class="font-mono underline" href={resolve('/items/[id]', { id: made })}>{made}</a>
				in backlog.
			</p>
		{:else if epics.length}
			<div class="mt-2 flex flex-wrap items-end gap-2">
				<label class="text-sm">
					<span class="block text-xs text-muted">Under the epic</span>
					<select class="rounded border border-line bg-surface px-2 py-1 text-sm" bind:value={epic}>
						{#each epics as e (e.id)}
							<option value={e.id}>{e.id} {e.title}</option>
						{/each}
					</select>
				</label>
				<button
					type="button"
					class="{button} border border-line"
					disabled={busy || !epic}
					onclick={writeWithAgent}>Write these with an agent</button
				>
			</div>
		{:else}
			<p class="mt-2 text-muted">
				A story needs an open epic to go under; make one from the board first.
			</p>
		{/if}
		{@render said(storySaid, 'kluster-agent-said')}
	</div>
</div>
