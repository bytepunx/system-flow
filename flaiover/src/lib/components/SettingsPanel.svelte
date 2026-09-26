<script lang="ts">
	// The host's settings for this project (S-0105), changed through flai once the operator has
	// turned on the settings host action in a shell. Until then, and for the settings kept for every
	// project until it is on everywhere, each section shows what is set and the command that allows
	// changing it. Commands are written one argument a line: they are run as they stand, never
	// through a shell, so nothing is quoted or split.
	import { api } from '$lib/api';
	import { agentFrom, configText, parseConfig } from '$lib/agent';
	import { projectState } from '$lib/project.svelte';
	import {
		allowed,
		argv,
		health,
		lines,
		type ServedProject,
		type SettingsKind,
		type SettingsView
	} from '$lib/settings';
	import AgentFields from './AgentFields.svelte';

	let view = $state<SettingsView | null>(null);
	let error = $state<string | null>(null);
	let busy = $state<string | null>(null);
	/** What the last change in each section said, by section. */
	let said = $state<Record<string, { ok: boolean; text: string }>>({});
	/** A rotated token's login link or MCP token, shown once. */
	let once = $state<Record<string, string>>({});

	// the forms, filled from what flai said
	let dHarness = $state('');
	let dModel = $state('');
	let dConfig = $state('');
	let aName = $state('');
	let aCommand = $state('');
	let harnessForm = $state<Record<string, { program: string; args: string }>>({});
	let checkForm = $state<Record<string, string>>({});
	let newCheckName = $state('');
	let newCheckCommand = $state('');
	let timeout = $state('');
	let newFolder = $state('');
	/** The root of the project whose Remove is waiting for a yes. */
	let confirming = $state<string | null>(null);

	function fill(v: SettingsView) {
		const h = v.host;
		if (!h) return;
		dHarness = h.default_agent?.harness ?? '';
		dModel = h.default_agent?.model ?? '';
		dConfig = configText(h.default_agent?.config);
		aName = h.agent.name ?? '';
		aCommand = lines(h.agent.command);
		harnessForm = Object.fromEntries(
			Object.entries(h.agent.harnesses).map(([n, x]) => [
				n,
				{ program: x.program, args: lines(x.args) }
			])
		);
		checkForm = Object.fromEntries(h.checks.commands.map((c) => [c.name, lines(c.command)]));
		timeout = String(h.checks.timeout_minutes);
	}

	async function load() {
		try {
			const r = await api('/api/settings');
			const body = await r.json();
			if (!r.ok) {
				error = body.error ?? r.statusText;
				return;
			}
			error = null;
			view = body;
			fill(body);
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
		}
	}
	$effect(() => {
		void load();
	});

	/** One change through flai; the section says what came of it, and the page reads flai again. */
	async function change(section: string, kind: SettingsKind, params: Record<string, unknown>) {
		busy = section;
		try {
			const r = await api('/api/settings', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ kind, ...params })
			});
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				said[section] = { ok: false, text: body.error ?? r.statusText };
				return null;
			}
			said[section] = { ok: true, text: 'saved' };
			await load();
			return body;
		} finally {
			busy = null;
		}
	}

	const can = (kind: SettingsKind) => (view ? allowed(view, kind) : { ok: false, enable: '' });

	function saveDefaultAgent() {
		const parsed = parseConfig(dConfig);
		if ('error' in parsed) {
			said.default = { ok: false, text: `config: ${parsed.error}` };
			return;
		}
		void change('default', 'default_agent', {
			agent: agentFrom(dHarness, dModel, parsed.config) ?? null
		});
	}
	function saveAgent() {
		const params: Record<string, unknown> = {};
		if (aName.trim()) params.name = aName.trim();
		const cmd = argv(aCommand);
		if (cmd.length) params.command = cmd;
		void change('agent', 'agent', params);
	}
	// A project served or removed through flai serve (S-0122). A served one connects within a second
	// or two: the switcher is asked again until it has it, so it appears without a reload.
	async function serveProject(p: ServedProject) {
		if (!(await change('projects', 'serve', { root: p.root, key: p.key || undefined }))) return;
		said.projects = { ok: true, text: `${p.key || p.root} is served; waiting for it to connect…` };
		for (let i = 0; i < 30; i++) {
			await projectState.refresh();
			if (projectState.list.some((x) => x.key === p.key && !x.candidate && x.connected)) {
				said.projects = { ok: true, text: `${p.key} is served, and in the switcher` };
				return;
			}
			await new Promise((r) => setTimeout(r, 500));
		}
		said.projects = {
			ok: true,
			text: `${p.key || p.root} is served but has not connected yet; its state below says why`
		};
		await load();
	}
	async function unserveProject(p: ServedProject) {
		confirming = null;
		if (!(await change('projects', 'unserve', { root: p.root, key: p.key || undefined }))) return;
		said.projects = {
			ok: true,
			text: `${p.key} is no longer served, and none of its files was touched. Serve it again with flai serve project add ${p.root}, on the host.`
		};
		await projectState.refresh();
	}

	async function rotate(section: string, kind: 'mcp_token' | 'dashboard_token') {
		const body = await change(section, kind, {});
		if (!body) return;
		once[section] = kind === 'dashboard_token' ? (body.login_url ?? '') : (body.token ?? '');
	}
</script>

{#snippet gate(kind: SettingsKind)}
	{#if view && !can(kind).ok}
		<p class="text-xs text-muted" data-testid="gate">
			Read-only here. On the host, run <code class="rounded bg-surface px-1"
				>{can(kind).enable}</code
			>
			to change it from the dashboard.
		</p>
	{/if}
{/snippet}

{#snippet result(section: string)}
	{#if said[section]}
		<p
			class="text-xs {said[section].ok ? 'text-good' : 'text-danger'}"
			role={said[section].ok ? 'status' : 'alert'}
			data-testid="said-{section}"
		>
			{said[section].text}
		</p>
	{/if}
{/snippet}

{#if error}
	<div class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</div>
{:else if !view}
	<p class="text-sm text-muted">Asking flai on the host…</p>
{:else if !view.host}
	<p class="text-sm text-muted">flai on the host keeps no settings for this project.</p>
{:else}
	{@const h = view.host}
	<div class="space-y-6 text-sm" data-testid="settings">
		<div class="rounded border border-line-strong bg-raised p-3" data-testid="settings-state">
			{#if view.everywhere}
				<p>Settings can be changed from the dashboard, for this project and every other.</p>
			{:else if view.here}
				<p>
					This project's settings can be changed here. Those kept for every project on the host need
					<code class="rounded bg-surface px-1">{view.enable_everywhere}</code>.
				</p>
			{:else}
				<p>
					Settings are read-only. On the host, in the project, run
					<code class="rounded bg-surface px-1">{view.enable}</code>, or
					<code class="rounded bg-surface px-1">{view.enable_everywhere}</code> for every project.
				</p>
			{/if}
			<p class="mt-1 text-xs text-muted">
				With settings on, whoever holds the dashboard token can change what flai runs on the host,
				as the operator. Only a shell turns it off: <code>flai serve disable settings</code>.
			</p>
		</div>

		<section data-testid="section-actions">
			<h2 class="mb-2 font-medium">Host actions, for this project</h2>
			{@render gate('action')}
			<ul class="space-y-2">
				{#each h.actions as a (a.name)}
					<li class="flex items-start gap-3">
						<input
							type="checkbox"
							class="mt-1"
							checked={a.here || a.everywhere}
							disabled={!can('action').ok || a.name === 'settings' || a.everywhere || busy !== null}
							aria-label="{a.name} on for this project"
							data-testid="action-{a.name}"
							onchange={(e) =>
								change('actions', 'action', { action: a.name, on: e.currentTarget.checked })}
						/>
						<div>
							<span class="font-mono">{a.name}</span>
							{#if a.everywhere}<span class="text-xs text-muted"> on for every project</span>{/if}
							{#if a.name === 'settings'}<span class="text-xs text-muted">
									changed only in a shell</span
								>{/if}
							<p class="text-xs text-muted">{a.means}</p>
						</div>
					</li>
				{/each}
			</ul>
			{@render result('actions')}
		</section>

		<section data-testid="section-default-agent">
			<h2 class="mb-2 font-medium">Default agent, for this project's new stories</h2>
			{@render gate('default_agent')}
			<AgentFields
				bind:harness={dHarness}
				bind:model={dModel}
				bind:config={dConfig}
				isDefault
				note="Stored in system-flow.yaml and committed; empty fields leave it without them"
			/>
			<button
				class="mt-2 rounded border border-line-strong px-3 py-1 disabled:opacity-50"
				disabled={!can('default_agent').ok || busy !== null}
				data-testid="save-default-agent"
				onclick={saveDefaultAgent}>Save default agent</button
			>
			{@render result('default')}
		</section>

		<section data-testid="section-agent">
			<h2 class="mb-2 font-medium">Starting agents, on this host</h2>
			{@render gate('agent')}
			<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
				<label class="block">
					<span class="mb-1 block text-xs text-muted">name: agents work as name-S-0001</span>
					<input
						class="w-full rounded border border-line-strong bg-surface px-2 py-1"
						bind:value={aName}
						placeholder="agent"
						data-testid="agent-name"
					/>
				</label>
			</div>
			<label class="mt-2 block">
				<span class="mb-1 block text-xs text-muted"
					>command for stories that name no harness, one argument a line; {'{story}'}, {'{root}'},
					{'{model}'}, and {'{harness}'} are replaced</span
				>
				<textarea
					class="h-20 w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
					bind:value={aCommand}
					data-testid="agent-command"></textarea>
			</label>
			<div class="mt-2 flex gap-2">
				<button
					class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
					disabled={!can('agent').ok || busy !== null}
					data-testid="save-agent"
					onclick={saveAgent}>Save</button
				>
				<button
					class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
					disabled={!can('agent').ok || busy !== null || h.agent.command.length === 0}
					data-testid="clear-agent-command"
					onclick={() => change('agent', 'agent', { command: null })}>Remove the command</button
				>
			</div>
			{@render result('agent')}

			{#each Object.keys(harnessForm).sort() as name (name)}
				<div class="mt-4 rounded border border-line p-3" data-testid="harness-{name}">
					<h3 class="mb-2 font-mono text-xs">
						harness {name}{h.agent.harnesses[name]?.set ? '' : ' (defaults)'}
					</h3>
					<label class="block">
						<span class="mb-1 block text-xs text-muted">program</span>
						<input
							class="w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
							bind:value={harnessForm[name].program}
							data-testid="harness-program"
						/>
					</label>
					<label class="mt-2 block">
						<span class="mb-1 block text-xs text-muted"
							>arguments: what its agent may do, one a line</span
						>
						<textarea
							class="h-20 w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
							bind:value={harnessForm[name].args}
							data-testid="harness-args"></textarea>
					</label>
					<div class="mt-2 flex gap-2">
						<button
							class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!can('harness').ok || busy !== null}
							data-testid="save-harness"
							onclick={() =>
								change('harness-' + name, 'harness', {
									name,
									program: harnessForm[name].program.trim(),
									args: argv(harnessForm[name].args)
								})}>Save</button
						>
						<button
							class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!can('harness').ok || busy !== null || !h.agent.harnesses[name]?.set}
							data-testid="reset-harness"
							onclick={() => change('harness-' + name, 'harness', { name, reset: true })}
							>Back to the defaults</button
						>
					</div>
					{@render result('harness-' + name)}
				</div>
			{/each}
		</section>

		<section data-testid="section-checks">
			<h2 class="mb-2 font-medium">Checks, run on a story in review, on this host</h2>
			{@render gate('check')}
			{#if h.checks.commands.length === 0}
				<p class="text-xs text-muted">
					None set on the host{h.manifest_checks?.length
						? `: the manifest's are used (${h.manifest_checks.map((c) => c.name).join(', ')})`
						: ''}.
				</p>
			{/if}
			{#each h.checks.commands as c (c.name)}
				<div class="mt-2" data-testid="check-{c.name}">
					<span class="font-mono text-xs">{c.name}</span>
					<textarea
						class="h-16 w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
						bind:value={checkForm[c.name]}></textarea>
					<div class="flex gap-2">
						<button
							class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!can('check').ok || busy !== null}
							onclick={() =>
								change('checks', 'check', { name: c.name, command: argv(checkForm[c.name]) })}
							>Save</button
						>
						<button
							class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!can('check').ok || busy !== null}
							onclick={() => change('checks', 'check', { name: c.name, command: null })}
							>Remove</button
						>
					</div>
				</div>
			{/each}
			<div class="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-3">
				<input
					class="rounded border border-line-strong bg-surface px-2 py-1"
					bind:value={newCheckName}
					placeholder="name"
					data-testid="new-check-name"
				/>
				<textarea
					class="h-16 rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs sm:col-span-2"
					bind:value={newCheckCommand}
					placeholder="scripts/test.sh"
					data-testid="new-check-command"></textarea>
			</div>
			<button
				class="mt-2 rounded border border-line-strong px-3 py-1 disabled:opacity-50"
				disabled={!can('check').ok ||
					busy !== null ||
					!newCheckName.trim() ||
					argv(newCheckCommand).length === 0}
				data-testid="add-check"
				onclick={async () => {
					if (
						await change('checks', 'check', {
							name: newCheckName.trim(),
							command: argv(newCheckCommand)
						})
					) {
						newCheckName = '';
						newCheckCommand = '';
					}
				}}>Add check</button
			>
			<label class="mt-3 flex items-center gap-2">
				<span class="text-xs text-muted">time limit, minutes</span>
				<input
					class="w-20 rounded border border-line-strong bg-surface px-2 py-1"
					bind:value={timeout}
					inputmode="numeric"
					data-testid="checks-timeout"
				/>
				<button
					class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
					disabled={!can('checks_timeout').ok || busy !== null}
					onclick={() => change('checks', 'checks_timeout', { minutes: Number(timeout) })}
					>Save</button
				>
			</label>
			{@render result('checks')}
		</section>

		<section data-testid="section-import">
			<h2 class="mb-2 font-medium">Import folders, on this host</h2>
			{@render gate('import')}
			<p class="text-xs text-muted">
				The board offers the git repositories in these folders for import.
			</p>
			<ul class="mt-2 space-y-1">
				{#each h.import_roots as folder (folder)}
					<li class="flex items-center gap-2">
						<code class="text-xs">{folder}</code>
						<button
							class="rounded border border-line-strong px-2 text-xs disabled:opacity-50"
							disabled={!can('import').ok || busy !== null}
							onclick={() => change('import', 'import', { folder, add: false })}>Remove</button
						>
					</li>
				{/each}
			</ul>
			<div class="mt-2 flex gap-2">
				<input
					class="flex-1 rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
					bind:value={newFolder}
					placeholder="/home/you/git"
					data-testid="new-folder"
				/>
				<button
					class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
					disabled={!can('import').ok || busy !== null || !newFolder.trim()}
					data-testid="add-folder"
					onclick={async () => {
						if (await change('import', 'import', { folder: newFolder.trim(), add: true }))
							newFolder = '';
					}}>Add folder</button
				>
			</div>
			{@render result('import')}
		</section>

		<section data-testid="section-projects">
			<h2 class="mb-2 font-medium">Projects, served by flai on this host</h2>
			{#if h.projects?.error}
				<p class="text-xs text-danger" role="alert">{h.projects.error}</p>
			{:else if h.projects}
				{#if !h.projects.running}
					<p class="text-xs text-muted">
						flai serve is not running, so no project is on the dashboard. On the host, run
						<code class="rounded bg-surface px-1">flai serve start</code>.
					</p>
				{/if}
				<ul class="space-y-2">
					{#each h.projects.served as p (p.root)}
						<li class="rounded border border-line p-2" data-testid="served-{p.key}">
							<div class="flex flex-wrap items-baseline gap-2">
								<span class="font-mono">{p.key}</span>
								<span>{p.name}</span>
								<code class="text-xs text-muted">{p.root}</code>
							</div>
							<p
								class="text-xs {p.state === 'connected' ? 'text-good' : 'text-danger'}"
								data-testid="health"
							>
								{health(p)}
							</p>
							{#if p.from !== 'registry'}
								<p class="text-xs text-muted" data-testid="served-below">
									Served because it is below <code>{p.below}</code>{p.from === 'import'
										? ', named under Import folders; removing that folder stops it'
										: ', the folder flai serve was started in'}.
								</p>
							{:else if confirming === p.root}
								<div class="mt-1 text-xs" data-testid="confirm-remove">
									<p>
										Stop serving {p.name}? None of its files is touched and the dashboard keeps
										running for the other projects. Serve it again with
										<code>flai serve project add {p.root}</code> on the host.
									</p>
									<div class="mt-1 flex gap-2">
										<button
											class="rounded border border-danger px-2 text-danger disabled:opacity-50"
											disabled={busy !== null}
											data-testid="confirm-remove-yes"
											onclick={() => unserveProject(p)}>Remove</button
										>
										<button
											class="rounded border border-line-strong px-2"
											data-testid="confirm-remove-no"
											onclick={() => (confirming = null)}>Keep it</button
										>
									</div>
								</div>
							{:else}
								<button
									class="mt-1 rounded border border-line-strong px-2 text-xs disabled:opacity-50"
									disabled={!p.settings || busy !== null}
									data-testid="remove-project"
									onclick={() => (confirming = p.root)}>Remove</button
								>
							{/if}
							{#if p.from === 'registry' && !p.settings}
								<p class="text-xs text-muted" data-testid="gate">
									To remove it here, run <code class="rounded bg-surface px-1">{view.enable}</code>
									on the host, in {p.root}.
								</p>
							{/if}
						</li>
					{/each}
				</ul>
				{#if h.projects.served.length === 0}
					<p class="text-xs text-muted">No project is served.</p>
				{/if}
				{#if h.projects.unserved.length}
					<h3 class="mt-3 mb-1 text-xs font-medium">Below the import folders, not served</h3>
					<ul class="space-y-2">
						{#each h.projects.unserved as p (p.root)}
							<li class="rounded border border-line p-2" data-testid="unserved-{p.key || p.root}">
								<div class="flex flex-wrap items-baseline gap-2">
									{#if p.key}<span class="font-mono">{p.key}</span>{/if}
									{#if p.name}<span>{p.name}</span>{/if}
									<code class="text-xs text-muted">{p.root}</code>
								</div>
								<p class="text-xs text-danger" data-testid="reason">{p.reason}</p>
								<button
									class="mt-1 rounded border border-line-strong px-2 text-xs disabled:opacity-50"
									disabled={!p.settings || busy !== null}
									data-testid="serve-project"
									onclick={() => serveProject(p)}>Serve</button
								>
								{#if !p.settings}
									<p class="text-xs text-muted" data-testid="gate">
										To serve it here, run <code class="rounded bg-surface px-1">{view.enable}</code>
										on the host, in {p.root}.
									</p>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			{/if}
			{@render result('projects')}
		</section>

		<section data-testid="section-tokens">
			<h2 class="mb-2 font-medium">Tokens</h2>
			<p>
				MCP server for agents over HTTP:
				{#if h.mcp?.running}<code class="text-xs">{h.mcp.url}</code>{:else}not running{/if}
			</p>
			{@render gate('mcp_token')}
			<button
				class="mt-1 rounded border border-line-strong px-3 py-1 disabled:opacity-50"
				disabled={!can('mcp_token').ok || busy !== null}
				data-testid="rotate-mcp"
				onclick={() => rotate('mcp', 'mcp_token')}>Rotate the MCP token</button
			>
			<p class="text-xs text-muted">Every agent connected over HTTP must be given the new one.</p>
			{#if once.mcp}<p class="mt-1 text-xs" data-testid="once-mcp">
					New MCP token, shown once: <code class="break-all">{once.mcp}</code>
				</p>{/if}
			{@render result('mcp')}

			<div class="mt-4">
				{@render gate('dashboard_token')}
				<button
					class="rounded border border-line-strong px-3 py-1 disabled:opacity-50"
					disabled={!can('dashboard_token').ok || busy !== null}
					data-testid="rotate-dashboard"
					onclick={() => rotate('dashboard', 'dashboard_token')}>Rotate the dashboard token</button
				>
				<p class="text-xs text-muted">
					This browser stays logged in; every other session and every agent using the token must log
					in again.
				</p>
				{#if once.dashboard}<p class="mt-1 text-xs" data-testid="once-dashboard">
						New login link, shown once: <code class="break-all">{once.dashboard}</code>
					</p>{/if}
				{@render result('dashboard')}
			</div>
		</section>
	</div>
{/if}
