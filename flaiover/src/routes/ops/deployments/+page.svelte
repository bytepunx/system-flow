<script lang="ts">
	// /ops/deployments: the namespaces deployed to the project's kluster (S-0581), flai's
	// kluster.deployments through /api/kluster/deployments. It lists main's namespace first, then each
	// story's, as flai orders them, then the stories in review with no deployment: each with the story
	// and its title linking its item page, the namespace, the commit, when it was deployed, the
	// outcome, each step with its log, and how to reach it for a demo, through a port-forward as
	// ADR-0254 decides. It deploys a story in review that is not deployed or is deployed at an older
	// commit, tears a story's namespace down once the operator confirms, and promotes main, each
	// through flai's kluster.deploy, kluster.teardown, and kluster.promote. flai runs each in the
	// background and answers at once, so while one runs the page says so, shows the end of its log,
	// and asks again every few seconds until it ends. flai's refusal or error is shown word for word.
	// With the kluster host action off, or no kluster: block, the page says which in place of the
	// actions and links the Klusters page's setup, naming no command and no file to edit (S-0586).
	// With a block that names no deploy, which kluster.status says (deployable), Deploy and Promote
	// give way to a sentence that deploys are not set up yet, linking the same setup.
	import type {
		KlusterDeployment,
		KlusterDeployments,
		KlusterJob,
		KlusterReview,
		KlusterStatus
	} from '$lib/server/repo';
	import { api } from '$lib/api';
	import { localTime } from '$lib/localtime';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';

	// How often the page asks again while a deploy, a teardown, or a promotion runs; each takes minutes.
	const POLL_MS = 3000;
	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';

	let data = $state<KlusterDeployments | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// The host action was refused (403) on a write.
	let refused = $state(false);
	// Whether the block names a deploy, as kluster.status says; undefined from a flai older than
	// S-0586, or when it could not be read, and then Deploy and Promote are offered as before.
	let deployable = $state<boolean | undefined>(undefined);
	// flai's refusal or error for the last action, word for word.
	let actionError = $state<string | null>(null);
	let busy = $state(false);
	// The story whose teardown awaits the operator's confirmation.
	let confirming = $state<string | null>(null);

	const off = $derived(refused || (!!data && !data.enabled));
	// no deploy command yet: flai deploys and promotes nothing, so neither is offered
	const undeployable = $derived(deployable === false);
	// a job runs: one this flai serve started, or one that holds kluster's lock against a namespace
	const lockJob = $derived(data?.deployments.find((d) => d.job)?.job);
	const runs = $derived(
		!!data?.running || !!lockJob || !!data?.deployments.some((d) => d.outcome === 'running')
	);
	const main = $derived(data?.deployments.find((d) => d.default));
	const disabled = $derived(busy || runs || !!off);
	const nothing = $derived(!!data && data.deployments.length === 0 && data.review.length === 0);

	const errorText = (e: unknown) => (e instanceof Error ? e.message : String(e));
	const short = (commit: string | undefined) => (commit ? commit.slice(0, 12) : '—');
	const jobLine = (j: KlusterJob) =>
		j.story
			? `kluster ${j.action} of ${j.story} in the namespace ${j.namespace}`
			: j.namespace
				? `kluster ${j.action} of main into the namespace ${j.namespace}`
				: `kluster ${j.action} of the cluster ${j.cluster}`;
	const outcomeClass = (o: string | undefined) =>
		o === 'failed' || o === 'interrupted'
			? 'text-danger'
			: o === 'passed' || o === 'done'
				? 'text-good-strong'
				: 'text-muted';
	// a story's namespace may be deployed again when its branch moved on, or its deploy did not pass
	const redeployable = (d: KlusterDeployment) =>
		d.status === 'review' &&
		((!!d.head && d.head !== d.commit) || d.outcome === 'failed' || d.outcome === 'interrupted');
	const stepId = (ns: string, step: string) => `log-${ns}-${step}`;

	/** Reads the deployments; a failed read keeps what was last shown. */
	async function load() {
		try {
			const r = await api('/api/kluster/deployments');
			const body = await r.json().catch(() => ({}));
			if (r.status === 501) {
				tooOld = body.error ?? r.statusText;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			data = body as KlusterDeployments;
			tooOld = error = null;
			if (data.enabled) refused = false;
		} catch (e) {
			error = `Could not read the deployments: ${errorText(e)}`;
		}
	}

	/** Reads whether the block names a deploy from kluster.status; a failed read leaves it unknown. */
	async function loadStatus() {
		try {
			const r = await api('/api/kluster');
			if (!r.ok) return;
			const body = (await r.json()) as KlusterStatus;
			deployable = body.block ? body.deployable : undefined;
		} catch {
			// unknown: Deploy and Promote stay offered, and flai refuses them if it must
		}
	}

	/** Deploys, tears down, or promotes through flai, taking the deployments the host answers. */
	async function act(action: 'deploy' | 'teardown' | 'promote', story?: string) {
		confirming = null;
		busy = true;
		actionError = null;
		try {
			const r = await api(`/api/kluster/${action}`, {
				method: 'POST',
				...(story
					? { headers: { 'content-type': 'application/json' }, body: JSON.stringify({ story }) }
					: {})
			});
			const body = await r.json().catch(() => ({}));
			if (r.ok) {
				data = body as KlusterDeployments;
				return;
			}
			if (r.status === 501) tooOld = body.error ?? r.statusText;
			else {
				if (r.status === 403) refused = true;
				actionError = body.error ?? r.statusText;
			}
		} catch (e) {
			actionError = errorText(e);
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void load();
		void loadStatus();
		// asked again only while a job runs: nothing else here changes without the page asking
		const timer = setInterval(() => {
			if (runs) void load();
		}, POLL_MS);
		return () => clearInterval(timer);
	});
</script>

{#snippet storyLink(story: string, title: string | undefined)}
	<a class="font-mono font-medium underline" href={resolve('/items/[id]', { id: story })}>{story}</a
	>
	{#if title}<span class="ml-1">{title}</span>{/if}
{/snippet}

{#snippet setup(text: string)}
	<a class="underline" href={resolve('/ops/klusters')} data-testid="deployments-setup">{text}</a>
{/snippet}

{#snippet deployButton(story: string)}
	{#if !undeployable}
		<button
			type="button"
			class="{button} bg-primary text-on-primary"
			{disabled}
			onclick={() => act('deploy', story)}>Deploy</button
		>
	{/if}
{/snippet}

{#snippet promoteButton()}
	{#if !undeployable}
		<button
			type="button"
			class="{button} bg-primary text-on-primary"
			{disabled}
			onclick={() => act('promote')}>Promote</button
		>
	{/if}
{/snippet}

{#snippet job(j: KlusterJob, now: boolean)}
	<div
		class="mt-4 rounded border border-line p-3 text-sm"
		data-testid={now ? 'deployments-running' : 'deployments-last'}
	>
		{#if now}
			<p role="status">{jobLine(j)} runs, since {localTime(j.started_at)}.</p>
		{:else}
			<p>
				Last: {jobLine(j)}
				<span class={outcomeClass(j.outcome)}>{j.outcome}</span>
				at {j.ended_at ? localTime(j.ended_at) : '—'}.
			</p>
			{#if j.error}
				<p class="mt-1 text-danger" data-testid="deployments-last-error">{j.error}</p>
			{/if}
		{/if}
		{#if j.log}
			<p class="mt-1 text-xs text-muted">
				Its progress: <code class="font-mono">{j.log}</code>{#if j.logs}; its steps' logs:
					<code class="font-mono">{j.logs}</code>{/if}
			</p>
		{/if}
		{#if j.log_tail}
			<pre
				class="mt-2 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs"
				data-testid="deployments-log">{j.log_tail}</pre>
		{/if}
	</div>
{/snippet}

{#snippet access(d: KlusterDeployment)}
	<details class="mt-3 text-sm" data-testid="deployment-access">
		<summary class="cursor-pointer">Reach it for a demo</summary>
		<p class="mt-1 text-xs text-muted">
			Run these on the host, in the project's main checkout, where flai keeps the cluster's
			kubeconfig, <code class="font-mono">{d.access.kubeconfig}</code>. Each port-forward serves the
			service on the same port of the host's localhost.
		</p>
		{#if d.access.services_error}
			<p class="mt-2 text-danger" data-testid="deployment-services-error">
				{d.access.services_error}
			</p>
		{/if}
		{#if d.access.services.length}
			<ul class="mt-2 space-y-2">
				{#each d.access.services as s (s.name)}
					<li data-testid="deployment-service">
						<p>
							<code class="font-mono">{s.host ?? s.name}</code>
							{#if s.ports.length}
								<span class="text-muted">
									on {s.ports
										.map(
											(p) =>
												`${p.port}${p.protocol && p.protocol !== 'TCP' ? `/${p.protocol}` : ''}`
										)
										.join(', ')}</span
								>
							{/if}
						</p>
						{#if s.port_forward}
							<pre
								class="mt-1 overflow-auto rounded bg-surface p-2 font-mono text-xs"
								data-testid="deployment-port-forward">{s.port_forward}</pre>
						{/if}
					</li>
				{/each}
			</ul>
		{:else if !d.access.services_error}
			<p class="mt-2 text-muted">The namespace has no services.</p>
		{/if}
		<p class="mt-2 text-xs text-muted">Its services, listed:</p>
		<pre
			class="mt-1 overflow-auto rounded bg-surface p-2 font-mono text-xs"
			data-testid="deployment-list">{d.access.list}</pre>
	</details>
{/snippet}

{#snippet deployed(d: KlusterDeployment)}
	<li class="rounded border border-line p-3" data-testid="deployment-row">
		<div class="flex flex-wrap items-baseline justify-between gap-2">
			<p class="text-sm" data-testid="deployment-who">
				{#if d.default}
					<span class="font-medium">main</span>
				{:else}
					{@render storyLink(d.story, d.title)}
				{/if}
			</p>
			<div class="flex flex-wrap gap-2">
				{#if d.default}
					{@render promoteButton()}
				{:else}
					{#if redeployable(d)}{@render deployButton(d.story)}{/if}
					<button
						type="button"
						class="{button} bg-danger text-on-primary"
						disabled={disabled || confirming === d.story}
						onclick={() => (confirming = d.story)}>Teardown</button
					>
				{/if}
			</div>
		</div>
		<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm">
			<dt class="text-muted">Namespace</dt>
			<dd><code class="font-mono">{d.namespace}</code></dd>
			<dt class="text-muted">Commit</dt>
			<dd>
				<code class="font-mono" title={d.tag} data-testid="deployment-commit"
					>{short(d.commit)}</code
				>
				{#if d.head && d.head !== d.commit}
					<span class="ml-1 text-warn" data-testid="deployment-behind"
						>behind its head, <code class="font-mono">{short(d.head)}</code></span
					>
				{/if}
				{#if d.head_error}<span class="ml-1 text-muted">{d.head_error}</span>{/if}
			</dd>
			<dt class="text-muted">Deployed</dt>
			<dd>{localTime(d.ended_at ?? d.started_at)}</dd>
			<dt class="text-muted">Outcome</dt>
			<dd class={outcomeClass(d.outcome)} data-testid="deployment-outcome">
				{d.outcome}{#if d.failed_step}, at its {d.failed_step} step{/if}
			</dd>
		</dl>
		{#if d.error}
			<p class="mt-2 text-sm text-danger" data-testid="deployment-error">{d.error}</p>
		{/if}
		{#if d.job}
			<p class="mt-2 text-sm" role="status" data-testid="deployment-job">
				kluster {d.job.kind} runs against it, since {localTime(d.job.started_at)}.
			</p>
		{/if}
		{#if confirming === d.story}
			<div
				class="mt-2 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
				role="alertdialog"
				aria-label="Tear down {d.story}'s namespace"
				data-testid="deployment-confirm-teardown"
			>
				<p>
					Tear down {d.story}'s namespace {d.namespace}? The project's teardown removes what its
					deploy made, and flai deletes the namespace.
				</p>
				<div class="mt-2 flex gap-2">
					<button
						type="button"
						class="{button} bg-danger text-on-primary"
						onclick={() => act('teardown', d.story)}>Tear it down</button
					>
					<button
						type="button"
						class="{button} border border-line"
						onclick={() => (confirming = null)}>Keep it</button
					>
				</div>
			</div>
		{/if}
		{#if d.steps.length}
			<div class="mt-3 space-y-1 text-sm" data-testid="deployment-steps">
				{#each d.steps as s (s.name)}
					<details id={stepId(d.namespace, s.name)} data-testid="deployment-step">
						<summary class="cursor-pointer">
							{s.name} <span class={outcomeClass(s.outcome)}>{s.outcome}</span>{#if s.exit_code},
								exit code {s.exit_code}{/if} · its log
						</summary>
						{#if s.log}
							<p class="mt-1 text-xs text-muted"><code class="font-mono">{s.log}</code></p>
						{/if}
						{#if s.log_tail}
							<pre
								class="mt-1 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs"
								data-testid="deployment-step-log">{s.log_tail}</pre>
						{:else}
							<p class="mt-1 text-xs text-muted">The log is empty or gone.</p>
						{/if}
					</details>
				{/each}
			</div>
		{/if}
		{@render access(d)}
	</li>
{/snippet}

{#snippet waiting(r: KlusterReview)}
	<li class="rounded border border-line p-3" data-testid="deployment-row">
		<div class="flex flex-wrap items-baseline justify-between gap-2">
			<p class="text-sm" data-testid="deployment-who">{@render storyLink(r.story, r.title)}</p>
			<div class="flex gap-2">{@render deployButton(r.story)}</div>
		</div>
		<p class="mt-2 text-sm text-muted" data-testid="deployment-outcome">
			In review, not deployed{#if r.head}; its head is <code class="font-mono">{short(r.head)}</code
				>{/if}.
		</p>
		{#if r.head_error}<p class="mt-1 text-sm text-muted">{r.head_error}</p>{/if}
	</li>
{/snippet}

<svelte:head><title>Deployments · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Deployments</h1>
<p class="mt-1 text-sm text-muted">
	What runs on the project's kluster: main in the default namespace, and each story in its own. flai
	on the host deploys, tears down, and promotes them.
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
		data-testid="deployments-too-old"
	>
		<p>This flai cannot show the deployments yet; upgrade flai on the host.</p>
		<p class="mt-1 text-xs">{tooOld}</p>
	</div>
{:else if data === null}
	{#if !error}
		<p class="mt-4 text-sm text-muted">Loading…</p>
	{/if}
{:else if !data.block}
	<div
		class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
		role="status"
		data-testid="deployments-no-block"
	>
		<p>The project's kluster is not set up yet, so nothing is deployed for it.</p>
		<p class="mt-1">Set it up from the {@render setup("Klusters page's setup")}.</p>
	</div>
{:else}
	{#if data.status_error}
		<p
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="deployments-status-error"
		>
			{data.status_error}
		</p>
	{/if}
	{#if off}
		<div
			class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="deployments-off"
		>
			<p>The kluster host action is off, so the dashboard cannot deploy, tear down, or promote.</p>
			<p class="mt-1">Turn it on from the {@render setup("Klusters page's setup")}.</p>
		</div>
	{/if}
	{#if undeployable}
		<div
			class="mt-4 rounded border border-line p-3 text-sm"
			role="status"
			data-testid="deployments-no-deploy"
		>
			<p>
				Deploys are not set up yet: the cluster can run, but no story is deployed and main is not
				promoted until the project names its deploy command.
			</p>
			<p class="mt-1">Name it on the {@render setup("Klusters page's setup")}.</p>
		</div>
	{/if}
	{#if actionError}
		<p
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="deployments-action-error"
		>
			{actionError}
		</p>
	{/if}
	{#if data.running}
		{@render job(data.running, true)}
	{:else if lockJob}
		<p
			class="mt-4 rounded border border-line p-3 text-sm"
			role="status"
			data-testid="deployments-running"
		>
			kluster {lockJob.kind} of {lockJob.story || 'main'} in the namespace {lockJob.namespace} runs, since
			{localTime(lockJob.started_at)}.
		</p>
	{:else if data.last}
		{@render job(data.last, false)}
	{/if}

	{#if nothing}
		<p class="mt-6 text-sm text-muted" data-testid="deployments-empty">
			Nothing is deployed to the cluster{data.cluster ? ` ${data.cluster}` : ''}, and no story is in
			review.
		</p>
	{/if}
	{#if !main}
		<div class="mt-4 flex flex-wrap items-center gap-2 text-sm" data-testid="deployments-no-main">
			<p>
				Main is not deployed{data.default_namespace
					? ` to the namespace ${data.default_namespace}`
					: ''}.
			</p>
			{@render promoteButton()}
		</div>
	{/if}
	{#if !nothing}
		<ul class="mt-6 space-y-3" aria-label="Deployments">
			{#each data.deployments as d (d.namespace)}
				{@render deployed(d)}
			{/each}
			{#each data.review as r (r.story)}
				{@render waiting(r)}
			{/each}
		</ul>
	{/if}
{/if}
