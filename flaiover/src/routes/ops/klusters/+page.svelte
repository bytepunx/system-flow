<script lang="ts">
	// /ops/klusters: the project's kluster (S-0580, S-0586), from flai's kluster.status through
	// /api/kluster. Until the cluster is set up, the page leads the operator through it in five
	// steps, each saying whether it is done: the host's tools, turning the kluster host action on
	// (which writes the kluster: block and its config, S-0586), the cluster's config file and its
	// profile, the commands that deploy into it, and bringing it up. Once set up, the steps fold away
	// and the page manages the cluster: its state, Up and Down (Down only once confirmed), the job
	// that runs and the last one, flai check's findings on the block, the nodes and pods kubectl
	// reports, turning the action off, and a link to Deployments. Every write goes through flai;
	// flai's refusals and errors are shown word for word. The block's and the config's writes need
	// the settings host action, which only the operator turns on, on the host (ADR-0039): with it off,
	// the page says so once, in one place, with what to run there. A job's output log, and each log
	// path in flai's refusals and errors, links the log page (S-0785).
	import type {
		KlusterJob,
		KlusterSetting,
		KlusterStatus,
		KlusterWorkloads,
		ProjectFile
	} from '$lib/server/repo';
	import { argv, lines, type Refusal } from '$lib/settings';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import DeploySettings from '$lib/components/DeploySettings.svelte';
	import LogLink from '$lib/components/LogLink.svelte';
	import LogText from '$lib/components/LogText.svelte';

	// How often the page asks again while an up or down runs; either takes minutes.
	const POLL_MS = 3000;
	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';
	const field = 'w-full rounded border border-line bg-surface px-2 py-1 font-mono text-sm';
	type Said = { ok: boolean; text: string };

	let status = $state<KlusterStatus | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// flai's refusal or error for the last Up or Down, word for word.
	let actionError = $state<string | null>(null);
	// A write refused because the kluster host action is off (403).
	let refusedOff = $state(false);
	let busy = $state(false);
	let confirming = $state(false);

	let workloads = $state<KlusterWorkloads | null>(null);
	let workloadsError = $state<string | null>(null);

	// Step 2: turning the kluster host action on or off, and what flai wrote when it was turned on.
	let toggleSaid = $state<Said | null>(null);
	let wrote = $state<{ wrote?: string[]; commit?: string; kept?: boolean } | null>(null);

	// Step 3: the block's config key as typed, and flai's refusals by field. Step 4's keys are
	// DeploySettings', the form the Deployments page's Setup shares (ADR-0263).
	let form = $state<Record<string, string>>({});
	let before = $state<Record<string, string>>({});
	let refused = $state<Refusal[]>([]);
	let saveSaid = $state<Record<string, Said>>({});

	// Step 3: the config file, its profile, and the whole file as typed.
	let config = $state<ProjectFile | null>(null);
	let configError = $state<string | null>(null);
	let configText = $state('');
	let profile = $state('');
	let configSaid = $state<Said | null>(null);

	const cluster = $derived(status?.cluster);
	const running = $derived(status?.running);
	const on = $derived(!!status?.enabled && !refusedOff);
	const settingsOn = $derived(!!status?.settings_on);
	const settings = $derived(status?.settings ?? []);
	const tools = $derived(status?.tools ?? []);
	const file = $derived(status?.config_file);
	const shown = $derived(
		!cluster ? null : cluster.running ? 'Running' : cluster.exists ? 'Stopped' : 'Not created'
	);
	const steps = $derived({
		tools: tools.length > 0 && tools.every((t) => t.found),
		on: on && !!status?.block,
		cluster: !!file?.exists,
		deploy: !!status?.deployable,
		up: !!cluster?.running
	});
	const setUp = $derived(Object.values(steps).every(Boolean));
	const settingOf = (key: string) => settings.find((s) => s.key === key);

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));
	const jobLine = (j: KlusterJob) => `kluster ${j.action} of the cluster ${j.cluster}`;
	const typed = (s: KlusterSetting) =>
		s.kind === 'command'
			? lines(Array.isArray(s.value) ? (s.value as string[]) : undefined)
			: typeof s.value === 'string'
				? s.value
				: '';
	const refusalOf = (key: string) => refused.find((r) => r.field === key)?.reason;

	/** Reads the kluster's status; a failed read keeps what was last shown. */
	async function load() {
		const wasRunning = !!status?.running;
		try {
			const r = await api('/api/kluster');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) {
				tooOld = body.error ?? r.statusText;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			const next = body as KlusterStatus;
			const reform = !status || JSON.stringify(status.settings) !== JSON.stringify(next.settings);
			const reconfig = !status || status.config !== next.config;
			status = next;
			tooOld = error = null;
			if (next.enabled) refusedOff = false;
			if (reform) {
				const f: Record<string, string> = {};
				for (const s of next.settings ?? []) f[s.key] = typed(s);
				form = f;
				before = { ...f };
			}
			if (reconfig) void loadConfig();
			if (!next.running && (wasRunning || !workloads)) void loadWorkloads();
		} catch (e) {
			error = `Could not read the kluster: ${errorText(e)}`;
		}
	}

	/** Reads the cluster's nodes and pods; only once flai has a block to read them from. */
	async function loadWorkloads() {
		if (!status?.block) return;
		try {
			const r = await api('/api/kluster/workloads');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) return;
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			workloads = body as KlusterWorkloads;
			workloadsError = null;
		} catch (e) {
			workloadsError = errorText(e);
		}
	}

	/** Reads the config file the block names, to show and change its profile. */
	async function loadConfig() {
		if (!status?.block || !status.config) {
			config = null;
			return;
		}
		try {
			const r = await api('/api/kluster/config');
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				config = null;
				configError = r.status === 404 ? null : (body.error ?? r.statusText);
				return;
			}
			config = body as ProjectFile;
			configText = config.content;
			profile = status.config_file?.profile || status.config_file?.default_profile || '';
			configError = null;
		} catch (e) {
			configError = errorText(e);
		}
	}

	/** Brings the cluster up or takes it down through flai, taking the status the host answers. */
	async function control(action: 'up' | 'down') {
		confirming = false;
		busy = true;
		actionError = null;
		try {
			const r = await api(`/api/kluster/${action}`, { method: 'POST' });
			const body = await r.json().catch(() => ({}));
			if (r.ok) {
				status = body as KlusterStatus;
				return;
			}
			if (r.status === 501) tooOld = body.error ?? r.statusText;
			else {
				if (r.status === 403) refusedOff = true;
				actionError = body.error ?? r.statusText;
			}
		} catch (e) {
			actionError = errorText(e);
		} finally {
			busy = false;
		}
	}

	/** Turns the kluster host action on or off through settings.action; on writes the block. */
	async function toggle(next: boolean) {
		busy = true;
		toggleSaid = null;
		try {
			const r = await api('/api/settings', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ kind: 'action', action: 'kluster', on: next })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				toggleSaid = { ok: false, text: body.error ?? r.statusText };
				return;
			}
			wrote = next ? (body.kluster ?? null) : null;
			toggleSaid = { ok: true, text: next ? 'The kluster host action is on.' : 'It is off.' };
			await load();
		} catch (e) {
			toggleSaid = { ok: false, text: errorText(e) };
		} finally {
			busy = false;
		}
	}

	/** Saves the keys given that changed through settings.manifest; flai's refusals go by field. */
	async function save(section: string, keys: string[]) {
		const set: Record<string, string | string[]> = {};
		const unset: string[] = [];
		for (const key of keys) {
			const s = settingOf(key);
			if (!s || (form[key] ?? '') === (before[key] ?? '')) continue;
			if (s.kind === 'command') {
				const args = argv(form[key] ?? '');
				if (args.length) set[key] = args;
				else unset.push(key);
			} else {
				const v = (form[key] ?? '').trim();
				if (v) set[key] = v;
				else unset.push(key);
			}
		}
		if (!Object.keys(set).length && !unset.length) {
			saveSaid[section] = { ok: true, text: 'Nothing changed.' };
			return;
		}
		const params: Record<string, unknown> = { kind: 'manifest' };
		if (Object.keys(set).length) params.set = set;
		if (unset.length) params.unset = unset;
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
				saveSaid[section] = {
					ok: false,
					text: refused.length
						? 'Not saved: flai refused what is marked.'
						: (body.error ?? r.statusText)
				};
				return;
			}
			refused = [];
			saveSaid[section] = {
				ok: true,
				text: body.commit ? `Saved in ${String(body.commit).slice(0, 7)}.` : 'Saved.'
			};
			await load();
		} catch (e) {
			saveSaid[section] = { ok: false, text: errorText(e) };
		} finally {
			busy = false;
		}
	}

	/** The config text with its profile: line replaced, or added when it names none. */
	function withProfile(text: string, p: string): string {
		if (/^profile:.*$/m.test(text)) return text.replace(/^profile:.*$/m, `profile: ${p}`);
		return `${text}${text.endsWith('\n') || text === '' ? '' : '\n'}profile: ${p}\n`;
	}

	/** Saves the config file through file.save, with the hash it was read with. */
	async function saveConfig(content: string, message: string) {
		if (!config) return;
		busy = true;
		configSaid = null;
		try {
			const r = await api('/api/kluster/config', {
				method: 'PUT',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ path: config.path, content, hash: config.hash, message })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				configSaid = { ok: false, text: body.error ?? r.statusText };
				return;
			}
			configSaid = { ok: true, text: 'Saved.' };
			// the status first: loadConfig sets the profile picker from it
			await load();
			await loadConfig();
		} catch (e) {
			configSaid = { ok: false, text: errorText(e) };
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void load();
		// asked again only while a job runs: nothing else here changes without the page asking
		const timer = setInterval(() => {
			if (status?.running) void load();
		}, POLL_MS);
		return () => clearInterval(timer);
	});
</script>

{#snippet job(j: KlusterJob, now: boolean)}
	<div
		class="mt-4 rounded border border-line p-3 text-sm"
		data-testid={now ? 'kluster-running' : 'kluster-last'}
	>
		{#if now}
			<p role="status">{jobLine(j)} runs, since {j.started_at}.</p>
		{:else}
			<p>
				Last: {jobLine(j)}
				<span class={j.outcome === 'failed' ? 'text-danger' : 'text-good-strong'}>{j.outcome}</span>
				at {j.ended_at}.
			</p>
			{#if j.error}
				<p class="mt-1 text-danger" data-testid="kluster-last-error"><LogText text={j.error} /></p>
			{/if}
		{/if}
		<p class="mt-1 text-xs text-muted" data-testid="kluster-job-log">
			Its output: <LogLink path={j.log} />
		</p>
		{#if j.log_tail}
			<pre
				class="mt-2 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs"
				data-testid="kluster-log">{j.log_tail}</pre>
		{/if}
	</div>
{/snippet}

{#snippet done(ok: boolean)}
	<span
		class="ml-2 rounded px-2 text-xs {ok ? 'bg-good-soft text-good' : 'bg-warn-soft text-warn'}"
		data-testid="step-state">{ok ? 'Done' : 'To do'}</span
	>
{/snippet}

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
					disabled={!settingsOn || busy}
					placeholder={s.default ?? ''}
					bind:value={form[key]}></textarea>
			{:else}
				<input
					class={field}
					disabled={!settingsOn || busy}
					placeholder={s.default ?? ''}
					bind:value={form[key]}
				/>
			{/if}
			{#if refusalOf(key)}
				<span class="block text-xs text-danger" data-testid="refused-{key}">{refusalOf(key)}</span>
			{/if}
		</label>
	{/if}
{/snippet}

<svelte:head><title>Klusters · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Klusters</h1>
<p class="mt-1 text-sm text-muted">
	The project's local Kubernetes cluster, made by kluster. flai on the host sets it up, brings it
	up, and takes it down, as this page asks.
</p>

{#if error}
	<p class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{/if}
{#if tooOld}
	<div
		class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="kluster-too-old"
	>
		<p>This flai cannot show the kluster yet; upgrade flai on the host.</p>
		<p class="mt-1 text-xs">{tooOld}</p>
	</div>
{:else if status === null}
	{#if !error}
		<p class="mt-4 text-sm text-muted">Loading…</p>
	{/if}
{:else}
	{#if !settingsOn && status.settings_enable}
		<div
			class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="kluster-settings-off"
		>
			<p class="font-medium">One step on the host</p>
			<p class="mt-1">
				The settings host action is off for this project. It is your consent, given on the host, for
				the dashboard to change the project's host actions and its manifest, here the kluster host
				action and the kluster: block, whose commands flai runs on the host as you. Until it is on,
				those steps are shown but cannot be saved. The config file's profile can be changed either
				way.
			</p>
			<p class="mt-1">
				On the host, in the project, run <code class="font-mono">{status.settings_enable}</code>,
				then reload this page.
			</p>
		</div>
	{/if}

	<details class="mt-6 rounded border border-line p-3" open={!setUp} data-testid="kluster-setup">
		<summary class="cursor-pointer text-lg font-semibold">
			Set up{#if setUp}<span class="ml-2 text-sm font-normal text-good">set up</span>{/if}
		</summary>

		<ol class="mt-3 space-y-5">
			<li data-testid="step-tools">
				<h3 class="font-medium">1. The host's tools {@render done(steps.tools)}</h3>
				<ul class="mt-2 space-y-2 text-sm">
					{#each tools as t (t.name)}
						<li data-testid="tool-{t.name}">
							<span class="font-mono">{t.name}</span>
							{#if t.found}
								<span class="text-good">found</span>
								{#if t.path}<span class="font-mono text-xs text-muted">{t.path}</span>{/if}
							{:else}
								<span class="text-danger">missing</span>
								<span class="block text-xs text-muted">{t.need}</span>
								{#if t.install}
									<span class="block text-xs"
										>Install it on the host with <code class="font-mono">{t.install}</code></span
									>
								{/if}
								<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- an external page, the tool's own install docs, not a route of the dashboard -->
								<a class="text-xs underline" href={t.docs} target="_blank" rel="noreferrer"
									>How to install {t.name}</a
								>
							{/if}
						</li>
					{/each}
				</ul>
			</li>

			<li data-testid="step-on">
				<h3 class="font-medium">2. Turn kluster on {@render done(steps.on)}</h3>
				{#if on && status.block}
					<p class="mt-1 text-sm">The kluster host action is on, and the manifest has its block.</p>
				{:else}
					<p class="mt-1 text-sm">
						Turning it on lets flai bring the cluster up and deploy stories into it. For a project
						with no kluster: block it also writes one, naming only the config file, writes that
						config file, and commits both. Nothing is deployed until the deploy command is set in
						step 4.
					</p>
					<button
						type="button"
						class="{button} mt-2 bg-primary text-on-primary"
						disabled={busy || !settingsOn}
						onclick={() => toggle(true)}>Turn kluster on</button
					>
				{/if}
				{#if wrote?.wrote?.length}
					<p class="mt-2 text-sm" data-testid="kluster-wrote">
						flai wrote {wrote.wrote.join(', ')}{wrote.kept
							? ', keeping the config file there'
							: ''}{wrote.commit ? `, committed as ${wrote.commit}` : ', not committed'}.
					</p>
				{/if}
				{@render said(toggleSaid, 'kluster-toggle-said')}
			</li>

			<li data-testid="step-cluster">
				<h3 class="font-medium">3. The cluster {@render done(steps.cluster)}</h3>
				{#if !status.block}
					<p class="mt-1 text-sm text-muted">Turning kluster on in step 2 writes the config.</p>
				{:else}
					{#if file}
						<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm">
							<dt class="text-muted">Config file</dt>
							<dd class="font-mono">{status.config}</dd>
							<dt class="text-muted">Name</dt>
							<dd>{file.name || '—'}</dd>
							<dt class="text-muted">Provider</dt>
							<dd>{file.provider || '—'}</dd>
							<dt class="text-muted">Profile</dt>
							<dd data-testid="kluster-profile">
								{file.profile || `${file.default_profile}, kluster's default`}
							</dd>
							{#if file.addons.length}
								<dt class="text-muted">Addons</dt>
								<dd>{file.addons.join(', ')}</dd>
							{/if}
						</dl>
						{#if file.error}
							<p class="mt-2 text-sm text-danger">{file.error}</p>
						{/if}
						{#if !file.exists}
							<p class="mt-2 text-sm text-warn">The config file is not there yet.</p>
						{/if}
					{/if}
					{#if config}
						<div class="mt-3 flex flex-wrap items-end gap-2">
							<label class="text-sm">
								<span class="block text-xs text-muted"
									>The shared services kluster installs, chosen before the first Up</span
								>
								<select
									class="rounded border border-line bg-surface px-2 py-1 text-sm"
									disabled={busy}
									bind:value={profile}
									data-testid="kluster-profile-choice"
								>
									{#each file?.profiles ?? [] as p (p)}
										<option value={p}>{p}</option>
									{/each}
								</select>
							</label>
							<button
								type="button"
								class="{button} border border-line"
								disabled={busy || !profile || profile === file?.profile}
								onclick={() =>
									saveConfig(withProfile(config!.content, profile), `kluster profile ${profile}`)}
								>Save the profile</button
							>
						</div>
						<details class="mt-3">
							<summary class="cursor-pointer text-sm">Edit the whole config file</summary>
							<textarea
								class="{field} mt-2"
								rows="8"
								disabled={busy}
								bind:value={configText}
								data-testid="kluster-config-text"></textarea>
							<button
								type="button"
								class="{button} mt-2 border border-line"
								disabled={busy || configText === config.content}
								onclick={() => saveConfig(configText, 'kluster config')}>Save the file</button
							>
						</details>
					{/if}
					{#if configError}
						<p class="mt-2 text-sm text-danger" role="alert">{configError}</p>
					{/if}
					{@render said(configSaid, 'kluster-config-said')}
					{@render settingField('kluster.config')}
					<button
						type="button"
						class="{button} mt-2 border border-line"
						disabled={!settingsOn || busy}
						onclick={() => save('config', ['kluster.config'])}>Save the config path</button
					>
					{@render said(saveSaid.config, 'kluster-save-config')}
				{/if}
			</li>

			<li data-testid="step-deploy">
				<h3 class="font-medium">4. Deploying stories {@render done(steps.deploy)}</h3>
				{#if !status.block}
					<p class="mt-1 text-sm text-muted">
						Turn kluster on in step 2 first: the block these are written in comes with it.
					</p>
				{:else}
					<DeploySettings
						settings={status.settings ?? []}
						{settingsOn}
						settingsEnable={status.settings_enable}
						deployable={!!status.deployable}
						notice={false}
						saved={() => load()}
					/>
					<p class="mt-2 text-sm">
						The same setup is on the Deployments page's
						<a
							class="underline"
							href="{resolve('/ops/deployments')}#setup"
							data-testid="kluster-deploy-setup-link">Setup</a
						>, beside each build and its live log.
					</p>
				{/if}
			</li>

			<li data-testid="step-up">
				<h3 class="font-medium">5. Bring it up {@render done(steps.up)}</h3>
				{#if on && status.block}
					<p class="mt-1 text-sm">
						Up makes the cluster with kluster and installs its profile; it takes minutes.
					</p>
					<button
						type="button"
						class="{button} mt-2 bg-primary text-on-primary"
						disabled={busy || !!running || !!cluster?.running}
						onclick={() => control('up')}>Up</button
					>
				{:else}
					<p class="mt-1 text-sm text-muted">Turn kluster on in step 2 first.</p>
				{/if}
			</li>
		</ol>
	</details>

	{#if status.block}
		<section class="mt-6" aria-labelledby="kluster-heading">
			<h2 id="kluster-heading" class="text-lg font-semibold">{cluster?.name ?? 'The cluster'}</h2>
			{#if cluster}
				<dl
					class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm"
					data-testid="kluster-cluster"
				>
					<dt class="text-muted">Name</dt>
					<dd>{cluster.name}</dd>
					<dt class="text-muted">Profile</dt>
					<dd>{cluster.profile || '—'}</dd>
					<dt class="text-muted">Provider</dt>
					<dd>{cluster.provider}</dd>
					<dt class="text-muted">State</dt>
					<dd data-testid="kluster-state">{shown}</dd>
					<dt class="text-muted">Age</dt>
					<dd>{cluster.age || '—'}</dd>
				</dl>
			{/if}
			{#if status.status_error}
				<p
					class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
					role="alert"
					data-testid="kluster-status-error"
				>
					<LogText text={status.status_error} />
				</p>
			{/if}
			{#if status.findings?.length}
				<ul
					class="mt-4 list-disc rounded border border-warn bg-warn-soft p-3 pl-6 text-sm text-warn"
					data-testid="kluster-findings"
				>
					{#each status.findings as f (f)}
						<li>{f}</li>
					{/each}
				</ul>
			{/if}

			{#if !on}
				<p class="mt-4 text-sm text-warn" role="status" data-testid="kluster-off">
					The kluster host action is off, so flai neither brings the cluster up nor deploys into it.
					Turn it on in step 2.
				</p>
			{:else}
				<div class="mt-4 flex flex-wrap items-center gap-2">
					<button
						type="button"
						class="{button} bg-primary text-on-primary"
						disabled={busy || !!running || !!cluster?.running}
						onclick={() => control('up')}>Up</button
					>
					<button
						type="button"
						class="{button} bg-danger text-on-primary"
						disabled={busy || !!running || confirming}
						onclick={() => (confirming = true)}>Down</button
					>
					<button
						type="button"
						class="{button} border border-line"
						disabled={busy || !settingsOn}
						onclick={() => toggle(false)}>Turn kluster off</button
					>
					<a class="text-sm underline" href={resolve('/ops/deployments')}>Deployments</a>
				</div>
				{#if confirming}
					<div
						class="mt-2 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
						role="alertdialog"
						aria-label="Take the cluster down"
						data-testid="kluster-confirm-down"
					>
						<p>
							Take the cluster {cluster?.name ?? ''} down? kluster deletes it, with everything deployed
							in it.
						</p>
						<div class="mt-2 flex gap-2">
							<button
								type="button"
								class="{button} bg-danger text-on-primary"
								onclick={() => control('down')}>Take it down</button
							>
							<button
								type="button"
								class="{button} border border-line"
								onclick={() => (confirming = false)}>Keep it</button
							>
						</div>
					</div>
				{/if}
			{/if}
			{#if actionError}
				<p
					class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
					role="alert"
					data-testid="kluster-action-error"
				>
					<LogText text={actionError} />
				</p>
			{/if}
			{#if running}
				{@render job(running, true)}
			{:else if status.last}
				{@render job(status.last, false)}
			{/if}

			<h3 class="mt-6 font-medium">Nodes and pods</h3>
			{#if workloadsError}
				<p class="mt-2 text-sm text-danger" role="alert">{workloadsError}</p>
			{:else if workloads?.workloads_error}
				<p class="mt-2 text-sm text-muted" data-testid="kluster-workloads-error">
					{workloads.workloads_error}
				</p>
			{:else if workloads}
				<table class="mt-2 text-sm" data-testid="kluster-nodes">
					<thead>
						<tr class="text-left text-muted">
							<th class="pr-4">Node</th>
							<th class="pr-4">Ready</th>
							<th class="pr-4">Roles</th>
							<th>Version</th>
						</tr>
					</thead>
					<tbody>
						{#each workloads.nodes as n (n.name)}
							<tr>
								<td class="pr-4 font-mono">{n.name}</td>
								<td class="pr-4 {n.ready ? 'text-good' : 'text-danger'}"
									>{n.ready ? 'yes' : 'no'}</td
								>
								<td class="pr-4">{n.roles.join(', ') || '—'}</td>
								<td>{n.version || '—'}</td>
							</tr>
						{/each}
					</tbody>
				</table>
				<div class="mt-3 space-y-3" data-testid="kluster-pods">
					{#each workloads.namespaces as ns (ns.namespace)}
						<div>
							<p class="font-mono text-sm font-medium">{ns.namespace}</p>
							<table class="text-sm">
								<tbody>
									{#each ns.pods as p (p.name)}
										<tr>
											<td class="pr-4 font-mono">{p.name}</td>
											<td
												class="pr-4 {p.phase === 'Running' || p.phase === 'Succeeded'
													? 'text-good'
													: 'text-warn'}">{p.phase}</td
											>
											<td class="pr-4">{p.ready}/{p.containers} ready</td>
											<td>{p.restarts} restarts</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/each}
				</div>
			{/if}
		</section>
	{/if}
{/if}
