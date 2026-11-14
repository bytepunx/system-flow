<script lang="ts">
	// One story namespace's row on the Deployments page's Overview (S-0663): the story and its
	// title linking its item page, its column, the namespace, the commit and whether it is the
	// branch's head, the last deploy and the last checks run, each with its outcome, its log, and
	// who asked, and a badge per step that opens Build for a build step and Logs for the others.
	// Its actions are the page's to run, through flai behind the kluster host action: Deploy, for a
	// story deployed at a commit other than its head or whose deploy failed or was interrupted,
	// giving way to Set up deploys when the block names no deploy; Run checks, once the namespace
	// is deployed at its head with a deploy that passed; and Teardown, which the row asks the
	// operator to confirm before it calls the page. A checks run has no step log for kluster.log
	// to read, so the row shows the end of its log itself. Each log file it names, the checks run's
	// and the failed step's, links the log page, and a log path in its error does too (S-0785); the
	// folder of a deploy's step logs is no file, so Its log opens Logs, which links each step's.
	import type { KlusterDeployment, KlusterStep } from '$lib/server/repo';
	import { localTime } from '$lib/localtime';
	import { resolve } from '$app/paths';
	import LogLink from './LogLink.svelte';
	import LogText from './LogText.svelte';

	let {
		deployment,
		deployable = true,
		busy = false,
		onDeploy,
		onChecks,
		onTeardown,
		onSetup,
		onLogs,
		onBuild
	}: {
		/** The story namespace's deployment, as kluster.deployments answers it. */
		deployment: KlusterDeployment;
		/** Whether the block names a deploy; without one, Deploy gives way to Set up deploys. */
		deployable?: boolean;
		/** Whether a job runs or the kluster host action is off: every action is disabled. */
		busy?: boolean;
		/** Deploys the story again, at its branch's head. */
		onDeploy?: (story: string) => void;
		/** Runs the story's checks against its namespace. */
		onChecks?: (story: string) => void;
		/** Tears the story's namespace down, once the operator has confirmed it here. */
		onTeardown?: (story: string) => void;
		/** Opens the page's Setup, where the block's deploy is named. */
		onSetup?: () => void;
		/** Opens Logs at the namespace, and at the step given. */
		onLogs?: (namespace: string, step?: string) => void;
		/** Opens Build at the namespace, and at the build step given. */
		onBuild?: (namespace: string, step: string) => void;
	} = $props();

	const button = 'rounded px-3 py-1 text-sm disabled:opacity-50';

	// Whether the operator is asked to confirm the namespace's teardown.
	let confirming = $state(false);

	const d = $derived(deployment);
	// the last job's log is its failed step's when it is a file, else its steps' logs' folder
	const failedLog = $derived(d.log?.endsWith('.log') ? d.log : undefined);
	const href = $derived(resolve('/items/[id]', { id: deployment.story }));
	// the story's column; a flai older than S-0663 sends the story's status alone
	const column = $derived(d.column ?? d.status);
	const columnText = $derived(
		column === 'in-progress' ? 'in progress' : column === 'review' ? 'in review' : column
	);
	const tone = $derived(
		column === 'in-progress'
			? 'ml-1 rounded bg-info-soft px-2 py-0.5 text-xs text-info'
			: 'ml-1 rounded bg-surface px-2 py-0.5 text-xs text-muted'
	);
	const open = $derived(column === 'in-progress' || column === 'review');
	const commit = $derived(short(d.commit));
	const behind = $derived(!!d.head && d.head !== d.commit);
	const redeployable = $derived(
		open && (behind || d.outcome === 'failed' || d.outcome === 'interrupted')
	);
	const checkable = $derived(!!d.at_head && d.outcome === 'passed');
	const deployOutcome = $derived(
		d.failed_step ? `${d.outcome}, at its ${d.failed_step} step` : d.outcome
	);
	// who asked for the last job: absent from a shell, and from a flai older than S-0663, which
	// says nothing of who asked
	const askedBy = $derived(d.asked_by || (d.column === undefined ? '—' : 'a shell'));
	const deployWhen = $derived(`${localTime(d.ended_at ?? d.started_at)}, asked by ${askedBy}`);
	const checksWho = $derived(
		`at ${short(d.checks_commit)}, asked by ${d.checks_asked_by || 'a shell or the dashboard'}`
	);

	function short(commit: string | undefined) {
		return commit ? commit.slice(0, 12) : '—';
	}
	const isBuild = (s: KlusterStep) => s.name === 'build' || s.name.startsWith('build-');
	// a step's name as the page shows it: a component's build by the component
	const stepLabel = (s: KlusterStep) => (s.name.startsWith('build-') ? s.name.slice(6) : s.name);
	const outcomeClass = (o: string | undefined) =>
		o === 'failed' || o === 'interrupted' || o === 'timed-out'
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

	function openStep(s: KlusterStep) {
		if (isBuild(s)) onBuild?.(d.namespace, s.name);
		else onLogs?.(d.namespace, s.name);
	}

	function tearDown() {
		confirming = false;
		onTeardown?.(d.story);
	}
</script>

<li class="rounded border border-line p-3" data-testid="deployment-row" data-story={d.story}>
	<div class="flex flex-wrap items-baseline justify-between gap-2">
		<p class="text-sm" data-testid="deployment-who">
			<a class="font-mono font-medium underline" {href}>{d.story}</a>
			{#if d.title}<span class="ml-1">{d.title}</span>{/if}
			{#if column}<span class={tone} data-testid="deployment-column">{columnText}</span>{/if}
		</p>
		<div class="flex flex-wrap gap-2">
			{#if redeployable && !deployable}
				<a
					class="{button} border border-line underline"
					href="#setup"
					data-testid="deployment-to-setup"
					onclick={() => onSetup?.()}>Set up deploys</a
				>
			{:else if redeployable}
				<button
					type="button"
					class="{button} bg-primary text-on-primary"
					disabled={busy}
					onclick={() => onDeploy?.(d.story)}>Deploy</button
				>
			{/if}
			{#if checkable}
				<button
					type="button"
					class="{button} border border-line"
					disabled={busy}
					onclick={() => onChecks?.(d.story)}>Run checks</button
				>
			{/if}
			<button
				type="button"
				class="{button} bg-danger text-on-primary"
				disabled={busy || confirming}
				onclick={() => (confirming = true)}>Teardown</button
			>
		</div>
	</div>
	<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-sm">
		<dt class="text-muted">Namespace</dt>
		<dd><code class="font-mono">{d.namespace}</code></dd>
		<dt class="text-muted">Commit</dt>
		<dd>
			<code class="font-mono" title={d.tag} data-testid="deployment-commit">{commit}</code>
			{#if d.at_head}
				<span class="ml-1 text-good-strong" data-testid="deployment-at-head">at its head</span>
			{:else if behind}
				<span class="ml-1 text-warn" data-testid="deployment-behind"
					>behind its head, <code class="font-mono">{short(d.head)}</code></span
				>
			{/if}
			{#if d.head_error}<span class="ml-1 text-muted">{d.head_error}</span>{/if}
		</dd>
		<dt class="text-muted">Last deploy</dt>
		<dd>
			<p class={outcomeClass(d.outcome)} data-testid="deployment-outcome">{deployOutcome}</p>
			<p class="text-muted" data-testid="deployment-asked-by">{deployWhen}</p>
			<p>
				<a
					class="underline"
					href="#logs"
					data-testid="deployment-logs"
					onclick={() => onLogs?.(d.namespace, d.failed_step)}>Its log</a
				>
			</p>
			{#if failedLog}
				<p data-testid="deployment-failed-log">
					Its {d.failed_step ?? 'failed'} step's log: <LogLink path={failedLog} />
				</p>
			{/if}
		</dd>
		<dt class="text-muted">Last checks</dt>
		<dd data-testid="deployment-checks">
			{#if d.checks_outcome}
				<p class={outcomeClass(d.checks_outcome)} data-testid="deployment-checks-outcome">
					{d.checks_outcome}
				</p>
				<p class="text-muted" data-testid="deployment-checks-asked-by">{checksWho}</p>
				{#if d.checks_log || d.checks_log_tail}
					<details data-testid="deployment-checks-log">
						<summary class="cursor-pointer underline">Its log</summary>
						{#if d.checks_log}
							<p class="mt-1 text-xs text-muted"><LogLink path={d.checks_log} /></p>
						{/if}
						{#if d.checks_log_tail}
							<pre
								class="mt-1 max-h-64 overflow-auto rounded bg-surface p-2 font-mono text-xs whitespace-pre-wrap"
								data-testid="deployment-checks-tail">{d.checks_log_tail}</pre>
						{/if}
					</details>
				{/if}
			{:else}
				<p class="text-muted">none since this deploy</p>
			{/if}
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
	{#if confirming}
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
					disabled={busy}
					onclick={tearDown}>Tear it down</button
				>
				<button
					type="button"
					class="{button} border border-line"
					onclick={() => (confirming = false)}>Keep it</button
				>
			</div>
		</div>
	{/if}
	{#if d.steps.length}
		<ul
			class="mt-3 flex flex-wrap gap-1 text-xs"
			aria-label="{d.story}'s steps"
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
						onclick={() => openStep(s)}>{stepLabel(s)} · {s.outcome}</a
					>
				</li>
			{/each}
		</ul>
	{/if}
</li>
