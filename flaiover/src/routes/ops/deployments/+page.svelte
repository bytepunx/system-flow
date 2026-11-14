<script lang="ts">
	// /ops/deployments: the namespaces deployed to the project's kluster (S-0581), flai's
	// kluster.deployments through /api/kluster/deployments, grouped in the page's own rail (S-0595,
	// ADR-0263 §5), kept in the URL's hash as the item page keeps its panes:
	// - Overview: main's namespace first, then each story's, in progress and in review alike, as flai
	//   orders them, each through DeploymentRow (S-0663), then the stories in progress or in review
	//   with no deployment, each with its column. It deploys a story that is not deployed, or is
	//   deployed at another commit than its head, runs a story's checks against its namespace, tears a
	//   story's namespace down once the operator confirms on its row, and promotes main, through
	//   flai's kluster.deploy, kluster.checks, kluster.teardown, and kluster.promote. With a block that
	//   names no deploy, which kluster.status says (deployable), Deploy and Promote give way to a link
	//   that opens Setup, never a call flai refuses. It says how many story namespaces stories in
	//   progress hold against kluster.agent_namespaces, and the newest deploy refused for it; the
	//   status badges above the rail carry the count.
	// - Build: for the namespace chosen, the running job's by default, else main's, one panel per build
	//   step of its record, build or build-<component>, each with its status and its log followed live
	//   (LogFeed), then the load.
	// - Logs: the same namespace's every step, each with its log followed live, and the end of the log
	//   of the job that runs or last ran. A row's step badge or log link opens Build or Logs at its
	//   step, that step's log open. Each step's log file and the job's progress log link the log page
	//   (S-0785); the folder of the steps' logs is no file, so it is named, not linked.
	// - Access: how to reach each namespace for a demo, through a port-forward as ADR-0254 decides,
	//   started and stopped through flai's kluster.forward.start and kluster.forward.stop (S-0586),
	//   bound to 0.0.0.0 as the operator chose on TH-0861, linked once it runs, for a story in
	//   progress's namespace as for one in review's; the page names no command.
	// - Setup: the block's namespace, components, build, deploy, teardown, and timeout (DeploySettings),
	//   saved through settings.manifest, so no one edits system-flow.yaml; the cluster and its config
	//   file stay on the Klusters page.
	// flai runs each job in the background and answers at once, so while one runs the page says so and
	// asks again every few seconds until it ends. flai's refusal or error is shown word for word. With
	// the kluster host action off, or no kluster: block, the page says which and links the Klusters
	// page's setup, naming no command and no file to edit (S-0586).
	import type {
		KlusterDeployment,
		KlusterDeployments,
		KlusterForward,
		KlusterJob,
		KlusterStatus,
		KlusterStep,
		KlusterUndeployed
	} from '$lib/server/repo';
	import type { SidebarEntry } from '$lib/pagesidebar.svelte';
	import { api } from '$lib/api';
	import { localTime } from '$lib/localtime';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import PageSidebar from '$lib/components/PageSidebar.svelte';
	import LogFeed from '$lib/components/LogFeed.svelte';
	import LogLink from '$lib/components/LogLink.svelte';
	import LogText from '$lib/components/LogText.svelte';
	import DeploySettings from '$lib/components/DeploySettings.svelte';
	import DeploymentRow from '$lib/components/DeploymentRow.svelte';

	// How often the page asks again while a deploy, a teardown, or a promotion runs; each takes minutes.
	const POLL_MS = 3000;
	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';

	let data = $state<KlusterDeployments | null>(null);
	let error = $state<string | null>(null);
	// The host's flai answered 501: its message, which says what it cannot do.
	let tooOld = $state<string | null>(null);
	// The host action was refused (403) on a write.
	let refused = $state(false);
	// kluster.status, for the setup and whether the block names a deploy; null until read, or from a
	// flai that cannot answer it.
	let kstatus = $state<KlusterStatus | null>(null);
	// flai's refusal or error for the last action, word for word.
	let actionError = $state<string | null>(null);
	let busy = $state(false);
	// flai's refusal or error for a port-forward's start or stop, beside its row, by forwardKey.
	let forwardSaid = $state<Record<string, string>>({});
	// The pane shown, bound to the rail, so a link on a pane can open another.
	let shown = $state<string>();
	// The namespace the Build and Logs panes show, as the reader picked it; the default otherwise.
	let picked = $state<string>();
	// The step whose log Build or Logs opens at, from a row's step badge or its log link.
	let stepAt = $state<string>();

	// Whether the block names a deploy, as kluster.status says; undefined from a flai older than
	// S-0586, or when it could not be read, and then Deploy and Promote are offered as before.
	const deployable = $derived(kstatus?.block ? kstatus.deployable : undefined);
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
	// The stories in progress or in review with no deployment: undeployed from a flai of S-0663 or
	// later, else the stories in review an older flai lists.
	const pending = $derived<KlusterUndeployed[]>(
		data?.undeployed ?? (data?.review ?? []).map((r) => ({ ...r, column: 'review' }))
	);
	const nothing = $derived(!!data && data.deployments.length === 0 && pending.length === 0);
	// The story namespaces of stories in progress against kluster.agent_namespaces, and the newest
	// request refused for it; absent from a flai older than S-0663.
	const holding = $derived(data?.holding ?? []);
	const limitBadge = $derived(
		data?.limit === undefined
			? undefined
			: data.limit === 0
				? 'deploys before review off'
				: `${data.live ?? holding.length} of ${data.limit} story namespaces`
	);
	const refusal = $derived(data?.refusals?.[0]);
	const forwards = $derived(
		(data?.deployments ?? []).reduce(
			(n, d) => n + (d.access.forwards?.filter((f) => f.running).length ?? 0),
			0
		)
	);
	const anyForward = $derived(!!data?.deployments.some((d) => d.access.forwards?.length));

	// The namespace Build and Logs show: the reader's pick, else the running job's, else main's, else
	// the first with a record.
	const records = $derived(data?.deployments ?? []);
	const runningNs = $derived(data?.running?.namespace ?? lockJob?.namespace);
	const chosen = $derived(
		records.find((d) => d.namespace === picked) ??
			records.find((d) => d.namespace === runningNs) ??
			main ??
			records[0]
	);
	const isBuild = (s: KlusterStep) => s.name === 'build' || s.name.startsWith('build-');
	const builds = $derived(chosen?.steps.filter(isBuild) ?? []);
	// the block builds nothing: no build step, or the one build step skipped as having no command
	const buildsNothing = $derived(
		builds.length === 0 ||
			(builds.length === 1 && builds[0].name === 'build' && builds[0].outcome === 'skipped')
	);
	const load = $derived(chosen?.steps.find((s) => s.name === 'load'));

	const entries = $derived<SidebarEntry[]>([
		{
			id: 'overview',
			label: 'Overview',
			mode: 'pane',
			count: data ? data.deployments.length + pending.length : undefined
		},
		{ id: 'build', label: 'Build', mode: 'pane', count: builds.length || undefined },
		{ id: 'logs', label: 'Logs', mode: 'pane' },
		{ id: 'access', label: 'Access', mode: 'pane', count: forwards || undefined },
		{ id: 'setup', label: 'Setup', mode: 'pane' }
	]);

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
	const badgeClass = (o: string | undefined) =>
		o === 'failed' || o === 'interrupted'
			? 'bg-danger-soft text-danger'
			: o === 'passed' || o === 'done'
				? 'bg-good-soft text-good-strong'
				: o === 'running'
					? 'bg-info-soft text-info'
					: 'bg-surface text-muted';
	// a step's name as the page shows it: a component's build by the component
	const stepLabel = (s: KlusterStep) => (s.name.startsWith('build-') ? s.name.slice(6) : s.name);
	// a story's column as the page says it, and its badge's tone, as DeploymentRow shows them
	const columnText = (c: string) =>
		c === 'in-progress' ? 'in progress' : c === 'review' ? 'in review' : c;
	const columnTone = (c: string) =>
		c === 'in-progress'
			? 'ml-1 rounded bg-info-soft px-2 py-0.5 text-xs text-info'
			: 'ml-1 rounded bg-surface px-2 py-0.5 text-xs text-muted';
	const forwardKey = (ns: string, service: string, port: number) => `${ns}/${service}/${port}`;
	const forwardOf = (d: KlusterDeployment, service: string, port: number) =>
		d.access.forwards?.find((f) => f.service === service && f.port === port);
	// A forward listens on 0.0.0.0 of the host, so it answers on whatever address the browser reached
	// the dashboard on: the host's own name, with the forward's port.
	const forwardUrl = (f: KlusterForward) => `http://${window.location.hostname}:${f.local_port}`;
	// how the build steps of the namespace stand, such as "2 passed, 1 running, 1 waiting"
	const tally = (steps: KlusterStep[]) => {
		const n = (o: string) => steps.filter((s) => s.outcome === o).length;
		return [
			[n('passed'), 'passed'],
			[n('running'), 'running'],
			[n('failed'), 'failed'],
			[n('pending'), 'waiting'],
			[n('skipped'), 'skipped']
		]
			.filter(([k]) => (k as number) > 0)
			.map(([k, w]) => `${k} ${w}`)
			.join(', ');
	};

	/** Opens a pane, and for Build and Logs the namespace given, at the step given. */
	function open(pane: string, ns?: string, step?: string) {
		if (ns) picked = ns;
		stepAt = step;
		shown = pane;
	}

	/** Reads the deployments; a failed read keeps what was last shown. */
	async function loadDeployments() {
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

	/** Reads kluster.status, for the setup and whether the block names a deploy; a failed read leaves
	 * it unknown. */
	async function loadStatus() {
		try {
			const r = await api('/api/kluster');
			if (!r.ok) return;
			kstatus = (await r.json()) as KlusterStatus;
		} catch {
			// unknown: Deploy and Promote stay offered, and flai refuses them if it must
		}
	}

	/** Deploys, checks, tears down, or promotes through flai, taking the deployments the host
	 * answers. */
	async function act(action: 'deploy' | 'checks' | 'teardown' | 'promote', story?: string) {
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
				picked = undefined;
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

	/** Starts or stops a port-forward through flai, taking the deployments the host answers. */
	async function forward(ns: string, service: string, port: number, on: boolean) {
		const key = forwardKey(ns, service, port);
		busy = true;
		forwardSaid = Object.fromEntries(Object.entries(forwardSaid).filter(([k]) => k !== key));
		try {
			const r = await api('/api/kluster/forward', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ namespace: ns, service, port, on })
			});
			const body = await r.json().catch(() => ({}));
			if (r.ok) {
				data = body as KlusterDeployments;
				return;
			}
			if (r.status === 501) tooOld = body.error ?? r.statusText;
			else {
				if (r.status === 403) refused = true;
				forwardSaid = { ...forwardSaid, [key]: body.error ?? r.statusText };
			}
		} catch (e) {
			forwardSaid = { ...forwardSaid, [key]: errorText(e) };
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void loadDeployments();
		void loadStatus();
		// asked again only while a job runs: nothing else here changes without the page asking
		const timer = setInterval(() => {
			if (runs) void loadDeployments();
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

{#snippet toSetup(testid: string)}
	<a
		class="{button} border border-line underline"
		href="#setup"
		data-testid={testid}
		onclick={() => open('setup')}>Set up deploys</a
	>
{/snippet}

{#snippet deployButton(story: string)}
	{#if undeployable}
		{@render toSetup('deployment-to-setup')}
	{:else}
		<button
			type="button"
			class="{button} bg-primary text-on-primary"
			{disabled}
			onclick={() => act('deploy', story)}>Deploy</button
		>
	{/if}
{/snippet}

{#snippet promoteButton()}
	{#if undeployable}
		{@render toSetup('deployment-to-setup')}
	{:else}
		<button
			type="button"
			class="{button} bg-primary text-on-primary"
			{disabled}
			onclick={() => act('promote')}>Promote</button
		>
	{/if}
{/snippet}

{#snippet chooser(id: string)}
	{#if records.length > 1}
		<label class="flex flex-wrap items-center gap-2 text-sm">
			<span class="text-muted">Namespace</span>
			<select
				class="rounded border border-line bg-ground px-2 py-1 text-sm"
				value={chosen?.namespace}
				onchange={(e) => (picked = (e.currentTarget as HTMLSelectElement).value)}
				data-testid={id}
			>
				{#each records as d (d.namespace)}
					<option value={d.namespace}
						>{d.default ? 'main' : d.story} · {d.namespace}{d.namespace === runningNs
							? ' · running'
							: ''}</option
					>
				{/each}
			</select>
		</label>
	{/if}
	{#if chosen}
		<p class="mt-2 text-sm" data-testid="deployments-chosen">
			{#if chosen.default}<span class="font-medium">main</span>{:else}{@render storyLink(
					chosen.story,
					chosen.title
				)}{/if}
			in <code class="font-mono">{chosen.namespace}</code>, at
			<code class="font-mono" title={chosen.tag}>{short(chosen.commit)}</code>:
			<span class={outcomeClass(chosen.outcome)}>{chosen.outcome}</span>{chosen.failed_step
				? `, at its ${chosen.failed_step} step`
				: ''}
		</p>
	{/if}
{/snippet}

{#snippet feed(d: KlusterDeployment, s: KlusterStep, label: string)}
	<LogFeed
		namespace={d.namespace}
		step={s.name}
		{label}
		open={s.name === stepAt}
		outcome={s.outcome}
		startedAt={s.started_at ?? ''}
		endedAt={s.ended_at ?? ''}
		exitCode={s.exit_code}
	/>
{/snippet}

{#snippet access(d: KlusterDeployment)}
	<li class="rounded border border-line p-3 text-sm" data-testid="access-row">
		<p data-testid="deployment-who">
			{#if d.default}<span class="font-medium">main</span>{:else}{@render storyLink(
					d.story,
					d.title
				)}{/if}
			<span class="ml-1 text-muted">in <code class="font-mono">{d.namespace}</code></span>
		</p>
		<div class="mt-2" data-testid="deployment-access">
			{#if d.access.services_error}
				<p class="mt-2 text-danger" data-testid="deployment-services-error">
					<LogText text={d.access.services_error} />
				</p>
			{/if}
			{#if d.access.services.length}
				<ul class="mt-2 space-y-2">
					{#each d.access.services as s (s.name)}
						<li data-testid="deployment-service">
							<p><code class="font-mono">{s.host ?? s.name}</code></p>
							<ul class="mt-1 space-y-1">
								{#each s.ports as p (p.port)}
									{@const f = forwardOf(d, s.name, p.port)}
									{@const key = forwardKey(d.namespace, s.name, p.port)}
									<li class="flex flex-wrap items-center gap-2" data-testid="deployment-port">
										<span class="text-muted"
											>port {p.port}{p.protocol && p.protocol !== 'TCP'
												? `/${p.protocol}`
												: ''}</span
										>
										{#if f?.running}
											<!-- eslint-disable svelte/no-navigation-without-resolve -- the forward on the host, not a route of the dashboard -->
											<a
												href={forwardUrl(f)}
												class="underline"
												target="_blank"
												rel="noreferrer"
												data-testid="deployment-forward-link">{forwardUrl(f)}</a
											>
											<!-- eslint-enable svelte/no-navigation-without-resolve -->
											<button
												type="button"
												class="{button} border border-line"
												disabled={busy || !!off}
												onclick={() => forward(d.namespace, s.name, p.port, false)}>Stop</button
											>
										{:else if !p.protocol || p.protocol === 'TCP'}
											<button
												type="button"
												class="{button} border border-line"
												disabled={busy || !!off}
												onclick={() => forward(d.namespace, s.name, p.port, true)}>Start</button
											>
										{/if}
										{#if f && !f.running && f.output}
											<span class="text-xs text-danger" data-testid="deployment-forward-ended"
												>Its port-forward ended: {f.output}</span
											>
										{/if}
										{#if forwardSaid[key]}
											<span
												class="text-xs text-danger"
												role="alert"
												data-testid="deployment-forward-said">{forwardSaid[key]}</span
											>
										{/if}
									</li>
								{/each}
							</ul>
						</li>
					{/each}
				</ul>
			{:else if !d.access.services_error}
				<p class="mt-2 text-muted">The namespace has no services.</p>
			{/if}
		</div>
	</li>
{/snippet}

{#snippet mainRow(d: KlusterDeployment)}
	<li class="rounded border border-line p-3" data-testid="deployment-row">
		<div class="flex flex-wrap items-baseline justify-between gap-2">
			<p class="text-sm" data-testid="deployment-who"><span class="font-medium">main</span></p>
			<div class="flex flex-wrap gap-2">{@render promoteButton()}</div>
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
			<p class="mt-2 text-sm text-danger" data-testid="deployment-error">
				<LogText text={d.error} />
			</p>
		{/if}
		{#if d.job}
			<p class="mt-2 text-sm" role="status" data-testid="deployment-job">
				kluster {d.job.kind} runs against it, since {localTime(d.job.started_at)}.
			</p>
		{/if}
		{#if d.steps.length}
			<ul
				class="mt-3 flex flex-wrap gap-1 text-xs"
				aria-label="main's steps"
				data-testid="deployment-steps"
			>
				{#each d.steps as s, i (`${i}-${s.name}`)}
					<li>
						<a
							class="inline-flex items-center gap-1 rounded px-2 py-0.5 {badgeClass(s.outcome)}"
							href={isBuild(s) ? '#build' : '#logs'}
							data-testid="deployment-step"
							data-step={s.name}
							data-outcome={s.outcome}
							onclick={() => open(isBuild(s) ? 'build' : 'logs', d.namespace, s.name)}
							>{stepLabel(s)} · {s.outcome}</a
						>
					</li>
				{/each}
			</ul>
		{/if}
	</li>
{/snippet}

{#snippet waiting(r: KlusterUndeployed)}
	<li class="rounded border border-line p-3" data-testid="deployment-row" data-story={r.story}>
		<div class="flex flex-wrap items-baseline justify-between gap-2">
			<p class="text-sm" data-testid="deployment-who">
				{@render storyLink(r.story, r.title)}
				<span class={columnTone(r.column)} data-testid="deployment-column"
					>{columnText(r.column)}</span
				>
			</p>
			<div class="flex gap-2">{@render deployButton(r.story)}</div>
		</div>
		<p class="mt-2 text-sm text-muted" data-testid="deployment-outcome">
			{r.column === 'in-progress' ? 'In progress' : 'In review'}, not deployed{#if r.head}; its head
				is <code class="font-mono">{short(r.head)}</code>{/if}.
		</p>
		{#if r.head_error}<p class="mt-1 text-sm text-muted">{r.head_error}</p>{/if}
	</li>
{/snippet}

{#snippet limitNotice()}
	{#if data?.limit !== undefined}
		<div class="mt-4 rounded border border-line p-3 text-sm" data-testid="deployments-limit-line">
			<p>
				{#if data.limit === 0}
					Deploys before review are off: no story in progress may hold a namespace.
				{:else}
					{limitBadge} live for stories in progress{holding.length
						? `: ${holding.join(', ')}`
						: ''}.
				{/if}
				Stories in review do not count.
			</p>
			{#if refusal}
				<p class="mt-2 text-warn" data-testid="deployments-refused">
					Refused for the limit: {refusal.story}'s deploy, asked by {refusal.by || 'a shell'} at
					{localTime(refusal.at)}.
				</p>
				<p class="mt-1 text-xs text-muted" data-testid="deployments-refused-reason">
					{refusal.reason}
				</p>
			{/if}
		</div>
	{/if}
{/snippet}

{#snippet overview()}
	{#if undeployable}
		<div
			class="rounded border border-line p-3 text-sm"
			role="status"
			data-testid="deployments-no-deploy"
		>
			<p>
				Deploys are not set up yet: the cluster can run, but no story is deployed and main is not
				promoted until the project names its deploy command.
			</p>
			<p class="mt-2">{@render toSetup('deployments-to-setup')}</p>
		</div>
	{/if}
	{#if nothing}
		<p class="mt-4 text-sm text-muted" data-testid="deployments-empty">
			Nothing is deployed to the cluster{data?.cluster ? ` ${data.cluster}` : ''}, and no story is
			in progress or in review.
		</p>
	{/if}
	{#if !main}
		<div class="mt-4 flex flex-wrap items-center gap-2 text-sm" data-testid="deployments-no-main">
			<p>
				Main is not deployed{data?.default_namespace
					? ` to the namespace ${data.default_namespace}`
					: ''}.
			</p>
			{@render promoteButton()}
		</div>
	{/if}
	{@render limitNotice()}
	{#if data && !nothing}
		<ul class="mt-4 space-y-3" aria-label="Deployments">
			{#each data.deployments as d (d.namespace)}
				{#if d.default}
					{@render mainRow(d)}
				{:else}
					<DeploymentRow
						deployment={d}
						deployable={!undeployable}
						busy={disabled}
						onDeploy={(id) => act('deploy', id)}
						onChecks={(id) => act('checks', id)}
						onTeardown={(id) => act('teardown', id)}
						onSetup={() => open('setup')}
						onLogs={(ns, step) => open('logs', ns, step)}
						onBuild={(ns, step) => open('build', ns, step)}
					/>
				{/if}
			{/each}
			{#each pending as r (r.story)}
				{@render waiting(r)}
			{/each}
		</ul>
	{/if}
{/snippet}

{#snippet build()}
	<div data-testid="deployments-build">
		<p class="text-sm text-muted">
			Each component the project builds, in the order the deploy builds them, with its status and
			its log as it is written.
		</p>
		{#if !chosen}
			<p class="mt-3 text-sm text-muted" data-testid="deployments-build-empty">
				Nothing has been deployed yet, so nothing has been built.
			</p>
		{:else}
			<div class="mt-3">{@render chooser('deployments-build-namespace')}</div>
			{#if buildsNothing}
				<p class="mt-3 text-sm text-muted" data-testid="deployments-build-none">
					This deploy built nothing: the project names no build command and no components. Name them
					under <a class="underline" href="#setup" onclick={() => open('setup')}>Setup</a>.
				</p>
			{:else}
				<p class="mt-3 text-sm" role="status" data-testid="deployments-build-summary">
					{`${builds.length} ${builds.length === 1 ? 'build' : 'builds'}: ${tally(builds)}.`}
				</p>
				<ul class="mt-2 space-y-2" aria-label="Builds" data-testid="deployments-builds">
					{#each builds as s (`${chosen.namespace}/${s.name}`)}
						<li class="rounded border border-line p-2" data-testid="build-panel" data-step={s.name}>
							{@render feed(chosen, s, stepLabel(s))}
						</li>
					{/each}
				</ul>
				{#if load}
					<div class="mt-3 rounded border border-line p-2" data-testid="build-load">
						{@render feed(chosen, load, 'load the images into the cluster')}
						{#if chosen.images?.length}
							<p class="mt-1 text-xs text-muted" data-testid="build-images">
								{chosen.images?.join(', ')}
							</p>
						{/if}
					</div>
				{/if}
			{/if}
		{/if}
	</div>
{/snippet}

{#snippet logs()}
	<div data-testid="deployments-logs">
		{@render jobLog(data?.running ?? data?.last)}
		{#if !chosen}
			<p class="mt-3 text-sm text-muted">
				No namespace has a deployment, so there are no step logs.
			</p>
		{:else}
			<div class="mt-3">{@render chooser('deployments-logs-namespace')}</div>
			<ul class="mt-2 space-y-2" aria-label="Steps" data-testid="deployments-step-logs">
				{#each chosen.steps as s, i (`${chosen.namespace}/${i}-${s.name}`)}
					<li class="rounded border border-line p-2" data-testid="step-panel" data-step={s.name}>
						{@render feed(chosen, s, s.name)}
						{#if s.log}
							<p class="mt-1 text-xs text-muted" data-testid="step-log">
								Its log: <LogLink path={s.log} />
							</p>
						{/if}
					</li>
				{/each}
			</ul>
		{/if}
	</div>
{/snippet}

{#snippet jobLog(j: KlusterJob | undefined)}
	{#if j && (j.log || j.log_tail)}
		<div class="rounded border border-line p-3 text-sm" data-testid="deployments-job-log">
			<p>{jobLine(j)}{j.outcome ? `, ${j.outcome}` : ', running'}</p>
			{#if j.log}
				<p class="mt-1 text-xs text-muted" data-testid="deployments-job-log-path">
					Its progress: <LogLink path={j.log} />{#if j.logs}; its steps' logs are in the folder
						<code class="font-mono">{j.logs}</code>, each linked with its step{/if}
				</p>
			{/if}
			{#if j.log_tail}
				<pre
					class="mt-2 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs"
					data-testid="deployments-log">{j.log_tail}</pre>
			{/if}
		</div>
	{/if}
{/snippet}

{#snippet accessPane()}
	{#if anyForward}
		<p
			class="rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="deployments-forwards-open"
		>
			A port-forward listens on every address of the host, so anything that can reach the host's
			network can reach it, with no sign-in of its own. Stop it when the demo is over.
		</p>
	{/if}
	<p class="mt-2 text-sm text-muted">
		Start a port-forward to reach a service from your browser. flai runs it on the host, on a port
		of its own.
	</p>
	{#if records.length}
		<ul class="mt-3 space-y-3" aria-label="How to reach each namespace">
			{#each records as d (d.namespace)}
				{@render access(d)}
			{/each}
		</ul>
	{:else}
		<p class="mt-3 text-sm text-muted">No namespace is deployed, so there is nothing to reach.</p>
	{/if}
{/snippet}

{#snippet setupPane()}
	<div data-testid="deployments-setup-pane">
		<p class="text-sm text-muted">
			What the project builds and how it deploys into a namespace. The cluster and its config file
			are on the {@render setup('Klusters page')}.
		</p>
		{#if kstatus?.settings}
			<div class="mt-3">
				<DeploySettings
					settings={kstatus.settings}
					settingsOn={!!kstatus.settings_on}
					settingsEnable={kstatus.settings_enable}
					deployable={!!kstatus.deployable}
					saved={(s) => {
						kstatus = s;
					}}
				/>
			</div>
		{:else}
			<p class="mt-3 text-sm text-muted" data-testid="deployments-setup-unknown">
				This flai does not say the block's settings; set them on the {@render setup(
					'Klusters page'
				)}.
			</p>
		{/if}
	</div>
{/snippet}

<svelte:head><title>Deployments · flaiover</title></svelte:head>

<h1 class="text-2xl font-semibold">Deployments</h1>
<p class="mt-1 text-sm text-muted">
	What runs on the project's kluster: main in the default namespace, and each story in its own. flai
	on the host builds, deploys, tears down, and promotes them.
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
	<ul
		class="mt-4 flex flex-wrap gap-2 text-xs"
		aria-label="Status"
		data-testid="deployments-status"
	>
		<li
			class="rounded px-2 py-0.5 {off ? 'bg-warn-soft text-warn' : 'bg-good-soft text-good-strong'}"
		>
			kluster {off ? 'off' : 'on'}
		</li>
		{#if data.cluster}
			<li class="rounded bg-surface px-2 py-0.5 text-muted">cluster {data.cluster}</li>
		{/if}
		{#if limitBadge}
			<li
				class="rounded px-2 py-0.5 {refusal ? 'bg-warn-soft text-warn' : 'bg-surface text-muted'}"
				data-testid="deployments-limit"
			>
				{limitBadge}
			</li>
		{/if}
		{#if deployable !== undefined}
			<li
				class="rounded px-2 py-0.5 {deployable
					? 'bg-good-soft text-good-strong'
					: 'bg-warn-soft text-warn'}"
				data-testid="deployments-deployable"
			>
				{deployable ? 'deploys set up' : 'deploys not set up'}
			</li>
		{/if}
		<li class="rounded px-2 py-0.5 {runs ? 'bg-info-soft text-info' : 'bg-surface text-muted'}">
			{runs ? 'a job runs' : 'no job runs'}
		</li>
	</ul>
	{#if data.status_error}
		<p
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="deployments-status-error"
		>
			<LogText text={data.status_error} />
		</p>
	{/if}
	{#if off}
		<div
			class="mt-4 rounded border border-warn bg-warn-soft p-3 text-sm text-warn"
			role="status"
			data-testid="deployments-off"
		>
			<p>
				The kluster host action is off, so the dashboard cannot deploy, run checks, tear down, or
				promote.
			</p>
			<p class="mt-1">Turn it on from the {@render setup("Klusters page's setup")}.</p>
		</div>
	{/if}
	{#if actionError}
		<p
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
			data-testid="deployments-action-error"
		>
			<LogText text={actionError} />
		</p>
	{/if}
	{#if data.running}
		<p
			class="mt-4 rounded border border-line p-3 text-sm"
			role="status"
			data-testid="deployments-running"
		>
			{jobLine(data.running)} runs, since {localTime(data.running.started_at)}.
			<a class="ml-1 underline" href="#logs" onclick={() => open('logs', data?.running?.namespace)}
				>Follow it</a
			>
		</p>
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
		<div class="mt-4 rounded border border-line p-3 text-sm" data-testid="deployments-last">
			<p>
				Last: {jobLine(data.last)}
				<span class={outcomeClass(data.last.outcome)}>{data.last.outcome}</span>
				at {data.last.ended_at ? localTime(data.last.ended_at) : '—'}.
			</p>
			{#if data.last.error}
				<p class="mt-1 text-danger" data-testid="deployments-last-error">
					<LogText text={data.last.error} />
				</p>
			{/if}
		</div>
	{/if}

	<div class="mt-6">
		<PageSidebar label="Deployment sections" {entries} fallback="overview" bind:current={shown}>
			{#snippet pane(id: string)}
				{#if id === 'overview'}
					{@render overview()}
				{:else if id === 'build'}
					{@render build()}
				{:else if id === 'logs'}
					{@render logs()}
				{:else if id === 'access'}
					{@render accessPane()}
				{:else if id === 'setup'}
					{@render setupPane()}
				{/if}
			{/snippet}
		</PageSidebar>
	</div>
{/if}
