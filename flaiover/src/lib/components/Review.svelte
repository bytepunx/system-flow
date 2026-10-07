<script lang="ts">
	// Review of one story (S-0041): what was asked, what the agent says, what
	// is being discussed, what changed, and what accepting will do; then accept
	// as the designer with each step shown, or send back with a reason.
	import UnreleasedList from './UnreleasedList.svelte';
	import WorktreeUncommitted, { otherBlockers } from './WorktreeUncommitted.svelte';
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import {
		criteriaOf,
		readNdjson,
		sectionOf,
		verifyView,
		type VerifyReport,
		type VerifyView
	} from '$lib/review';
	import { acceptLabel, issueRows, type Issue, type IssueRow } from '$lib/issues';
	import DiffView from '$lib/components/DiffView.svelte';
	import Threads from '$lib/components/Threads.svelte';
	import { render, enhance } from '$lib/markdown';
	import { themeState } from '$lib/theme.svelte';
	import { tick } from 'svelte';

	type Item = {
		id: string;
		type: string;
		title: string;
		status: string;
		body: string;
		path: string;
	};
	type Version = string | { Major: number; Minor: number; Patch: number };
	type Step = {
		component: { name: string };
		delivered: boolean;
		level: string;
		from: Version;
		to: Version;
		files: string[];
	};
	type Preview = {
		branch?: string;
		blockers?: string[];
		uncommitted?: string[];
		/** Uncommitted paths in the story's worktree: they block acceptance (S-0140). */
		worktree_uncommitted?: string[];
		/** Files the branch changes on a path Claude Code protects: only the operator accepts it (ADR-0106). */
		protected?: string[];
		/** The sentence saying so, present exactly when protected is. */
		operator_only?: string;
		plan?: {
			level: string;
			commits: string[];
			steps: Step[] | null;
			skipped?: string;
			unreleased?: { component: string; files: string[] }[];
		} | null;
	};
	type Progress = { step: string; msg: string };
	type CheckStep = {
		name: string;
		command?: string;
		started?: string;
		ended?: string;
		exit?: number | null;
	};
	type ChecksRun = {
		story?: string;
		running?: boolean;
		current?: string;
		steps?: CheckStep[];
		outcome?: string;
		started?: string;
		ended?: string;
	};

	let { id }: { id: string } = $props();

	let item = $state<Item | null>(null);
	let narrative = $state<string | null>(null);
	let currentStateHtml = $state('');
	let nextStepsHtml = $state('');
	let narrativeContent: HTMLElement | undefined = $state();
	let diff = $state<Record<string, unknown> | null>(null);
	let diffError = $state<string | null>(null);
	let preview = $state<Preview | null>(null);
	let previewError = $state<string | null>(null);
	let error = $state<string | null>(null);
	let writable = $state(false);

	let include = $state(false);
	let running = $state(false);
	let progress = $state<Progress[]>([]);
	let warnings = $state<string[]>([]);
	// by is who accepted, as the done transition records; verified is the commit the orchestrator's
	// verifier passed (ADR-0093). An older flai sends neither.
	let result = $state<{ tags?: string[]; by?: string; verified?: string } | null>(null);
	let failure = $state<string | null>(null);

	let sendingBack = $state(false);
	let reason = $state('');

	// Issues (S-0198): the open issues no open story links, the story's own
	// first and checked. Accepting makes a backlog story for each checked one,
	// once the acceptance has succeeded.
	let issues = $state<IssueRow[]>([]);
	let issuesError = $state<string | null>(null);
	let creating = $state(false);
	let made = $state<{ issue: string; id: string; title: string; note?: string }[]>([]);
	let unmade = $state<{ issue: string; error: string }[]>([]);

	// Checks (S-0082): the operator's named commands, run in the story's
	// worktree, gated on the checks host action. checksEnabled comes from
	// project.info's host_actions, asked afresh each load; checksRun is the
	// last known state, kept live by watchChecks while a run is active. watchToken guards against a stale loop from a previous
	// story still running after id changes.
	let checksEnabled = $state(false);
	let checksRun = $state<ChecksRun | null>(null);
	let checksLines = $state<string[]>([]);
	let checksBusy = $state<'run' | 'cancel' | null>(null);
	let checksError = $state<string | null>(null);
	let tailOffset = 0;
	let watchToken = 0;

	// Verification (S-0270): the last flai verify stored for the story, read
	// through verify.status; null when there is none or no flai is connected,
	// and then the section is left out. Nothing here runs it.
	let verify = $state<VerifyView | null>(null);
	let verifyError = $state<string | null>(null);

	function sleep(ms: number): Promise<void> {
		return new Promise((r) => setTimeout(r, ms));
	}

	const v = (x: Version) => (typeof x === 'string' ? x : `${x.Major}.${x.Minor}.${x.Patch}`);
	const criteria = $derived(item ? criteriaOf(item.body) : []);
	const ticked = $derived(criteria.filter((c) => c.checked).length);
	const inReview = $derived(item?.status === 'review');
	const needsChoice = $derived(!!preview?.uncommitted?.length && !include);
	const chosen = $derived(issues.filter((r) => r.checked));
	const canAccept = $derived(
		writable &&
			inReview &&
			!running &&
			!result &&
			!!preview &&
			!preview.blockers?.length &&
			!preview.worktree_uncommitted?.length &&
			!needsChoice
	);

	async function get<T>(path: string): Promise<T> {
		const r = await api(path);
		const data = await r.json();
		if (!r.ok) throw new Error(data.error ?? r.statusText);
		return data as T;
	}

	async function renderNarrative(target: string): Promise<void> {
		const docPath = `wip/agents/${target}.md`;
		const body = narrative ?? '';
		currentStateHtml = render(sectionOf(body, 'Current state') || '—', docPath);
		nextStepsHtml = render(sectionOf(body, 'Next steps') || '—', docPath);
		await tick();
		if (narrativeContent) await enhance(narrativeContent, themeState.dark);
	}

	async function load(target: string) {
		error = diffError = previewError = failure = null;
		item = null;
		narrative = diff = preview = result = null;
		currentStateHtml = nextStepsHtml = '';
		progress = [];
		warnings = [];
		issues = [];
		issuesError = null;
		made = [];
		unmade = [];
		verify = null;
		verifyError = null;
		try {
			item = (await get<{ item: Item }>(`/api/items/${target}`)).item;
		} catch (e) {
			error = e instanceof Error ? e.message : String(e);
			return;
		}
		api('/api/board')
			.then(async (r) => (writable = r.ok ? (await r.json()).writable : false))
			.catch(() => (writable = false));
		get<{ body: string }>(`/api/docs/file?path=${encodeURIComponent(`wip/agents/${target}.md`)}`)
			.then(async (d) => {
				narrative = d.body;
				await renderNarrative(target);
			})
			.catch(() => (narrative = ''));
		get<Record<string, unknown>>(`/api/items/${target}/diff`)
			.then((d) => (diff = d))
			.catch((e) => (diffError = e instanceof Error ? e.message : String(e)));
		if (item.status === 'review') {
			get<Preview>(`/api/items/${target}/acceptance`)
				.then((p) => (preview = p))
				.catch((e) => (previewError = e instanceof Error ? e.message : String(e)));
			void loadIssues(target);
		}
		void loadChecks(target);
		void loadVerify(target);
	}

	async function loadVerify(target: string) {
		try {
			const { report } = await get<{ report: VerifyReport | null }>(`/api/items/${target}/verify`);
			if (target === id) verify = report ? verifyView(report) : null;
		} catch (e) {
			if (target === id) verifyError = e instanceof Error ? e.message : String(e);
		}
	}

	async function issueList(query: string): Promise<Issue[]> {
		const data = await get<unknown>(`/api/issues${query}`);
		return Array.isArray(data) ? (data as Issue[]) : [];
	}

	// Both lists: the main checkout's, and the story's own, which holds the
	// issues it recorded on its branch until it is accepted.
	async function loadIssues(target: string) {
		try {
			const [main, own] = await Promise.all([
				issueList(''),
				issueList(`?story=${encodeURIComponent(target)}`)
			]);
			if (target === id) issues = issueRows(target, main, own);
		} catch (e) {
			if (target === id) issuesError = e instanceof Error ? e.message : String(e);
		}
	}

	// One at a time, from the main checkout: each makes a commit of its own.
	async function createStories(rows: IssueRow[]) {
		creating = true;
		for (const row of rows) {
			try {
				const r = await api(`/api/issues/${row.id}/story`, {
					method: 'POST',
					headers: { 'content-type': 'application/json' },
					body: '{}'
				});
				const data = await r.json().catch(() => ({}));
				if (!r.ok) unmade = [...unmade, { issue: row.id, error: data.error ?? r.statusText }];
				else
					made = [
						...made,
						{
							issue: row.id,
							id: String(data.id),
							title: String(data.title),
							note: data.commit_error ? `not committed: ${data.commit_error}` : undefined
						}
					];
			} catch (e) {
				unmade = [...unmade, { issue: row.id, error: e instanceof Error ? e.message : String(e) }];
			}
		}
		creating = false;
	}
	$effect(() => {
		void load(id);
	});

	async function refreshChecksStatus(target: string, token: number) {
		try {
			const r = await api(`/api/items/${target}/checks`);
			const data = await r.json();
			if (token !== watchToken) return;
			if (r.ok) {
				checksEnabled = data.checks_enabled === true;
				checksRun = data as ChecksRun;
			}
		} catch {
			// leave the last known state showing rather than clear it
		}
	}

	function resetChecks() {
		// A separate function, not an inline assignment in loadChecks: assigning
		// checksRun = null there and reading it after the await below tripped a
		// TypeScript 6.0 control-flow bug (narrows the read to `never`).
		checksRun = null;
		checksLines = [];
		checksError = null;
		checksBusy = null;
		tailOffset = 0;
	}

	async function loadChecks(target: string) {
		const token = ++watchToken;
		resetChecks();
		await refreshChecksStatus(target, token);
		if (token === watchToken && checksRun?.running) void watchChecks(target, token);
	}

	// Drives the live output while a run is active: checks.status for the
	// summary (which check, pass/fail so far), checks.tail in a loop for
	// the output — each call returns once it has new content or its own
	// wait elapses, then the client calls again with the returned offset
	// (review.ts's readNdjson reads the same NDJSON shape accept already
	// uses). Stops itself once running is false.
	async function watchChecks(target: string, token: number) {
		while (token === watchToken) {
			await refreshChecksStatus(target, token);
			if (token !== watchToken || !checksRun?.running) return;
			let ok = true;
			try {
				const r = await api(`/api/items/${target}/checks/tail?from=${tailOffset}`);
				if (!r.ok) {
					const data = await r.json().catch(() => ({}));
					if (token === watchToken) checksError = data.error ?? r.statusText;
					return;
				}
				await readNdjson(r, (l) => {
					if (token !== watchToken) return;
					if (l.event === 'line') checksLines = [...checksLines, String(l.text)];
					else if (l.event === 'done') tailOffset = Number(l.offset) || tailOffset;
					else if (l.event === 'error') {
						checksError = String(l.error);
						ok = false;
					}
				});
			} catch (e) {
				if (token === watchToken) checksError = e instanceof Error ? e.message : String(e);
				return;
			}
			if (!ok) return;
		}
	}

	async function runChecks() {
		if (checksBusy || !item) return;
		checksBusy = 'run';
		checksError = null;
		checksLines = [];
		tailOffset = 0;
		const target = id;
		// flai checks run blocks until the whole sequence ends, real minutes:
		// race the request against a short wait and fall back to polling and
		// tailing rather than waiting on this response, the same shape
		// HostPanel's own restart/upgrade already uses for a slow action.
		const attempt = api(`/api/items/${target}/checks`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ action: 'run' })
		})
			.then(async (r) => ({ ok: r.ok, body: await r.json().catch(() => ({})) }))
			.catch((e) => ({ ok: false, body: { error: e instanceof Error ? e.message : String(e) } }));
		const outcome = await Promise.race([attempt, sleep(1500).then(() => 'started' as const)]);
		checksBusy = null;
		if (outcome !== 'started' && !outcome.ok) {
			checksError = (outcome.body as { error?: string })?.error ?? 'could not start the checks run';
			return;
		}
		const token = watchToken;
		await refreshChecksStatus(target, token);
		void watchChecks(target, token);
		if (outcome === 'started') {
			// The write had not answered within the short wait: it is still
			// running for real, so once it eventually does, take one more look
			// in case watchChecks's own loop had already stopped on a transient
			// tail error before the run was actually done.
			void attempt.then(() => {
				if (token === watchToken) void refreshChecksStatus(target, token);
			});
		}
	}

	async function cancelChecks() {
		if (checksBusy || !item) return;
		checksBusy = 'cancel';
		checksError = null;
		const target = id;
		try {
			const r = await api(`/api/items/${target}/checks`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action: 'cancel' })
			});
			const data = await r.json().catch(() => ({}));
			if (!r.ok) {
				checksError = data.error ?? r.statusText;
				return;
			}
			await refreshChecksStatus(target, watchToken);
		} finally {
			checksBusy = null;
		}
	}

	async function accept() {
		if (!canAccept) return;
		const toMake = chosen;
		running = true;
		failure = null;
		progress = [];
		warnings = [];
		try {
			const r = await api(`/api/items/${id}/accept`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ include_uncommitted: include || undefined })
			});
			if (!r.ok) {
				failure = (await r.json().catch(() => ({}))).error ?? r.statusText;
				return;
			}
			await readNdjson(r, (l) => {
				if (l.event === 'progress')
					progress = [...progress, { step: String(l.step), msg: String(l.msg) }];
				else if (l.event === 'warning') warnings = [...warnings, String(l.msg)];
				else if (l.event === 'done') result = (l.result ?? {}) as typeof result;
				else if (l.event === 'error') failure = String(l.error);
			});
			if (!result && !failure) failure = 'the acceptance ended without a result';
		} catch (e) {
			failure = e instanceof Error ? e.message : String(e);
		} finally {
			running = false;
			if (failure) {
				// flai left the story in review; show it as it is now
				const fresh = await get<{ item: Item }>(`/api/items/${id}`).catch(() => null);
				if (fresh) item = fresh.item;
			}
		}
		// the stories only once the story is accepted: a failed acceptance makes none
		if (result && !failure && toMake.length) await createStories(toMake);
	}

	async function sendBack() {
		if (!reason.trim()) return;
		failure = null;
		const r = await api(`/api/items/${id}/move`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ to: 'in-progress', reason: reason.trim() })
		});
		const data = await r.json().catch(() => ({}));
		if (!r.ok) failure = data.error ?? r.statusText;
		else {
			sendingBack = false;
			reason = '';
			await load(id);
		}
	}
</script>

{#if error}
	<p class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger" role="alert">
		{error}
	</p>
{:else if !item}
	<p class="text-sm text-muted">Loading…</p>
{:else}
	<div class="mb-4 flex flex-wrap items-baseline gap-3">
		<h1 class="text-xl font-semibold">Review {item.id}</h1>
		<a class="underline" href={resolve('/items/[id]', { id: item.id })}>{item.title}</a>
		<span class="rounded border border-line px-2 py-0.5 text-xs">{item.status}</span>
	</div>

	{#if result}
		<div class="mb-4 rounded border border-good bg-good-soft p-3 text-sm text-good" role="status">
			<p class="font-medium" data-testid="accepted">
				{item.id} is accepted{result.by ? ` by ${result.by}` : ''}{result.verified
					? ` at ${result.verified.slice(0, 12)}`
					: ''}.
			</p>
			{#if result.tags?.length}<p>Released: {result.tags.join(', ')}</p>{:else}<p>
					Nothing was released.
				</p>{/if}
			<p data-testid="accept-unpublished">
				It stays local until it is published: Publish on the board, or
				<code class="rounded bg-surface px-1 text-ink">flai release --pending</code> on the host.
			</p>
			{#if creating}<p>Making stories for the checked issues…</p>{/if}
			{#each made as m (m.issue)}
				<p data-testid="issue-story-made">
					Made <a class="underline" href={resolve('/items/[id]', { id: m.id })}>{m.id}</a>
					{m.title} from {m.issue}, in backlog.{#if m.note}
						{m.note}{/if}
				</p>
			{/each}
		</div>
		{#if unmade.length}
			<div
				class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
				role="alert"
				data-testid="issue-story-failed"
			>
				<p class="font-medium">No story was made for these issues. flai said:</p>
				<ul class="mt-1 ml-4 list-disc">
					{#each unmade as u (u.issue)}<li>
							<span class="font-medium">{u.issue}</span>: {u.error}
						</li>{/each}
				</ul>
			</div>
		{/if}
	{:else if !inReview}
		<p class="mb-4 rounded border border-line bg-surface p-3 text-sm" role="status">
			{item.id} is {item.status}, not in review, so there is nothing to accept or send back. What
			follows is how it stands.
		</p>
	{/if}

	<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
		<section class="min-w-0 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">
				Acceptance criteria
				<span class="text-xs font-normal text-muted">{ticked} of {criteria.length} checked</span>
			</h2>
			{#if criteria.length}
				<ul class="space-y-1">
					{#each criteria as c (c.text)}
						<li class="flex gap-2">
							<span aria-hidden="true" class={c.checked ? 'text-good' : 'text-danger'}
								>{c.checked ? '☑' : '☐'}</span
							>
							<span class:text-muted={c.checked}
								><span class="sr-only">{c.checked ? 'checked: ' : 'not checked: '}</span
								>{c.text}</span
							>
						</li>
					{/each}
				</ul>
			{:else}
				<p class="text-muted">The story has no acceptance criteria.</p>
			{/if}
		</section>

		<section class="min-w-0 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">
				Narrative
				<a
					class="text-xs font-normal underline"
					href={resolve('/docs/[...path]', { path: `wip/agents/${item.id}.md` })}>open</a
				>
			</h2>
			{#if narrative === null}
				<p class="text-muted">Loading…</p>
			{:else if !narrative}
				<p class="text-muted">No narrative for this story.</p>
			{:else}
				<div bind:this={narrativeContent}>
					<h3 class="text-xs font-medium text-muted">Current state</h3>
					<div class="prose mb-2 max-w-none">
						<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the narrative, rendered client side -->
						{@html currentStateHtml}
					</div>
					<h3 class="text-xs font-medium text-muted">Next steps</h3>
					<div class="prose max-w-none">
						<!-- eslint-disable-next-line svelte/no-at-html-tags -- markdown from the narrative, rendered client side -->
						{@html nextStepsHtml}
					</div>
				</div>
			{/if}
		</section>
	</div>

	<section class="mt-4 rounded border border-line bg-surface p-3 text-sm">
		<h2 class="mb-2 font-medium">Changes on the branch</h2>
		{#if diffError}
			<p class="text-muted">{diffError}</p>
		{:else if !diff}
			<p class="text-muted">Loading…</p>
		{:else}
			<DiffView diff={diff as never} />
		{/if}
	</section>

	{#if (inReview || result) && (issues.length || issuesError)}
		<section
			class="mt-4 rounded border border-line bg-surface p-3 text-sm"
			data-testid="issues-section"
		>
			<h2 class="mb-2 font-medium">
				Issues
				<span class="text-xs font-normal text-muted"
					>a backlog story is made for each one checked</span
				>
			</h2>
			{#if issuesError}
				<p class="text-muted">{issuesError}</p>
			{:else}
				<ul class="space-y-1">
					{#each issues as row (row.id)}
						<li>
							<label class="flex items-baseline gap-2">
								<input
									type="checkbox"
									bind:checked={row.checked}
									disabled={running || !!result || !writable}
								/>
								<span>
									<a class="underline" href={resolve('/docs/[...path]', { path: row.doc })}
										>{row.id}</a
									>
									{row.title}
									<span class="text-xs text-muted"
										>{row.class}{row.recorded ? ', recorded by this story' : ''}</span
									>
								</span>
							</label>
						</li>
					{/each}
				</ul>
			{/if}
		</section>
	{/if}

	{#if inReview || result}
		<section class="mt-4 rounded border border-line bg-surface p-3 text-sm">
			<h2 class="mb-2 font-medium">What accepting does</h2>
			{#if previewError}
				<p class="rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					{previewError}
				</p>
			{:else if !preview && !result}
				<p class="text-muted">Working it out…</p>
			{:else if preview}
				{#if otherBlockers(preview).length}
					<div
						class="mb-2 rounded border border-danger bg-danger-soft p-2 text-danger"
						role="alert"
					>
						<p class="font-medium">This cannot be accepted from here yet:</p>
						<ul class="mt-1 ml-4 list-disc">
							{#each otherBlockers(preview) as b (b)}<li>{b}</li>{/each}
						</ul>
					</div>
				{/if}
				{#if preview.worktree_uncommitted?.length}
					<div class="mb-2">
						<WorktreeUncommitted
							{id}
							paths={preview.worktree_uncommitted}
							branch={preview.branch}
							disabled={running || !writable}
						/>
					</div>
				{/if}
				{#if preview.protected?.length}
					<div
						class="mb-2 rounded border border-warn bg-warn-soft p-2 text-warn"
						data-testid="accept-protected"
					>
						<p class="font-medium">{preview.operator_only}</p>
						<ul class="mt-1 ml-4 list-disc font-mono text-xs">
							{#each preview.protected as p (p)}<li>{p}</li>{/each}
						</ul>
					</div>
				{/if}
				{#if preview.uncommitted?.length}
					<div class="mb-2 rounded border border-warn bg-warn-soft p-2 text-warn">
						<p class="font-medium">Uncommitted changes outside wip:</p>
						<ul class="mt-1 ml-4 list-disc font-mono text-xs">
							{#each preview.uncommitted as p (p)}<li>{p}</li>{/each}
						</ul>
						<label class="mt-2 flex items-center gap-2">
							<input type="checkbox" bind:checked={include} disabled={running} />
							Include these files in the acceptance commit
						</label>
					</div>
				{/if}
				<ul class="ml-4 list-disc space-y-1">
					{#if preview.branch}
						<li>
							Rebase <code>{preview.branch}</code>, fast-forward it into main, remove its worktree.
						</li>
					{/if}
					<li>Move the story to done, archive it with its tasks and narrative, and commit.</li>
					{#if preview.plan?.skipped}
						<li>
							No release: {preview.plan.skipped}.
							<UnreleasedList unreleased={preview.plan.unreleased} />
						</li>
					{:else if preview.plan?.steps?.length}
						<li>
							Tag a {preview.plan.level} release:
							{#each preview.plan.steps as s, i (s.component.name)}{i ? ', ' : ''}<span
									class="font-medium">{s.component.name}</span
								>
								{v(s.from)} → {v(s.to)}{/each}
						</li>
					{/if}
					{#if chosen.length}
						<li data-testid="accept-makes-stories">
							Then, once it is accepted, make a backlog story for each checked issue: {chosen
								.map((r) => r.id)
								.join(', ')}.
						</li>
					{/if}
					<li data-testid="accept-stays-local">
						Nothing is pushed: the accepted work stays local until it is published, from Publish on
						the board or with <code>flai release --pending</code> on the host.
					</li>
				</ul>
			{/if}

			{#if progress.length || running}
				<ol class="mt-3 space-y-1 border-t border-line pt-2" aria-live="polite">
					{#each progress as p (p.step)}
						<li><span class="text-good" aria-hidden="true">✓</span> {p.msg}</li>
					{/each}
					{#if running}<li class="text-muted">Working…</li>{/if}
				</ol>
			{/if}
			{#each warnings as w (w)}
				<p class="mt-2 rounded border border-warn bg-warn-soft p-2 text-warn">{w}</p>
			{/each}
			{#if failure}
				<div class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					<p class="font-medium">Not accepted. {item.id} is {item.status}. flai said:</p>
					<pre class="mt-1 font-mono text-xs whitespace-pre-wrap">{failure}</pre>
				</div>
			{/if}

			{#if writable && inReview && !result}
				<div class="mt-3 flex flex-wrap gap-2">
					<button
						type="button"
						class="rounded bg-primary px-3 py-1 text-on-primary disabled:opacity-50"
						disabled={!canAccept}
						onclick={accept}>{acceptLabel(chosen.length, running)}</button
					>
					<button
						type="button"
						class="rounded border border-line-strong px-3 py-1"
						disabled={running}
						onclick={() => (sendingBack = !sendingBack)}>Send back</button
					>
				</div>
				{#if sendingBack}
					<div class="mt-2">
						<label class="block text-xs text-muted" for="send-back-reason"
							>Why it goes back to in-progress (recorded in the story's notes)</label
						>
						<textarea
							id="send-back-reason"
							class="mt-1 h-20 w-full rounded border border-line bg-ground p-2"
							bind:value={reason}></textarea>
						<button
							type="button"
							class="mt-1 rounded border border-line-strong px-3 py-1 disabled:opacity-50"
							disabled={!reason.trim()}
							onclick={sendBack}>Send back to in-progress</button
						>
					</div>
				{/if}
			{:else if !writable && inReview}
				<p class="mt-3 text-xs text-muted">Read-only: flai is not available to this dashboard.</p>
			{/if}
		</section>
	{/if}

	{#if inReview || checksRun?.started}
		<section
			class="mt-4 rounded border border-line bg-surface p-3 text-sm"
			data-testid="checks-section"
		>
			<h2 class="mb-2 font-medium">Checks</h2>
			{#if checksError}
				<p
					class="mb-2 rounded border border-danger bg-danger-soft p-2 text-danger"
					role="alert"
					data-testid="checks-error"
				>
					{checksError}
				</p>
			{/if}
			{#if checksRun?.started}
				<p data-testid="checks-summary">
					{#if checksRun.running}
						Running {checksRun.current}, started {checksRun.started}.
					{:else}
						{checksRun.outcome}, started {checksRun.started}, ended {checksRun.ended}.
					{/if}
				</p>
				{#if checksRun.steps?.length}
					<ul class="mt-1 ml-4 list-disc">
						{#each checksRun.steps as s (s.name)}
							<li>
								{s.name}: {s.exit == null
									? 'did not start'
									: s.exit === 0
										? 'ok'
										: `exit ${s.exit}`}
							</li>
						{/each}
					</ul>
				{/if}
				{#if checksLines.length}
					<pre
						class="mt-2 max-h-64 overflow-y-auto rounded bg-ground p-2 font-mono text-xs whitespace-pre-wrap"
						data-testid="checks-output">{checksLines.join('\n')}</pre>
				{/if}
			{:else if checksEnabled}
				<p class="text-muted">No checks have run yet.</p>
			{/if}

			{#if !checksEnabled}
				<p class="mt-2 text-muted">
					Off: the operator turns this on in a shell on the host with
					<code class="rounded bg-ground px-1 text-ink">flai serve enable checks</code>.
				</p>
			{:else if writable && inReview}
				<p class="mt-3 flex flex-wrap gap-2">
					{#if checksRun?.running}
						<button
							type="button"
							class="rounded border border-warn px-2 py-1 text-warn disabled:opacity-60"
							onclick={cancelChecks}
							disabled={checksBusy !== null}
							data-testid="checks-cancel"
							>{checksBusy === 'cancel' ? 'Cancelling…' : 'Cancel'}</button
						>
					{:else}
						<button
							type="button"
							class="rounded border border-line px-2 py-1 disabled:opacity-60"
							onclick={runChecks}
							disabled={checksBusy !== null}
							data-testid="checks-run">{checksBusy === 'run' ? 'Starting…' : 'Run checks'}</button
						>
					{/if}
				</p>
			{/if}
		</section>
	{/if}

	{#if verify || verifyError}
		<section
			class="mt-4 rounded border border-line bg-surface p-3 text-sm"
			data-testid="verify-section"
		>
			<h2 class="mb-2 font-medium">
				Verification
				<span class="text-xs font-normal text-muted">the last flai verify of the story</span>
			</h2>
			{#if verifyError}
				<p class="text-muted" data-testid="verify-error">{verifyError}</p>
			{:else if verify}
				<p data-testid="verify-summary">
					<span class={verify.passed ? 'font-medium text-good' : 'font-medium text-danger'}
						>{verify.outcome}</span
					>. Verified <code class="rounded bg-ground px-1 text-ink">{verify.commit}</code> against
					<code class="rounded bg-ground px-1 text-ink">{verify.base}</code> at {verify.ranAt}, in
					{verify.duration}.
				</p>
				<ul class="mt-2 space-y-1">
					{#each verify.rows as row, i (i)}
						<li data-testid="verify-step">
							<span class="flex flex-wrap items-baseline gap-2">
								<span
									class={row.state === 'passed'
										? 'text-good'
										: row.state === 'failed'
											? 'font-medium text-danger'
											: 'text-muted'}>{row.label}</span
								>
								<span class="font-medium">{row.name}</span>
								{#if row.tier}<span class="text-xs text-muted">tier</span>{/if}
								{#if row.duration}<span class="text-xs text-muted">{row.duration}</span>{/if}
							</span>
							{#if row.findings.length || row.omitted}
								<ul
									class="mt-1 ml-4 space-y-1 rounded border border-danger bg-danger-soft p-2 font-mono text-xs text-danger"
									aria-label="What {row.name} found"
								>
									{#each row.findings as f, j (j)}
										<li>
											{#if f.where}<span class="font-medium">{f.where}{f.message ? ':' : ''}</span
												>{/if}
											<span class="whitespace-pre-wrap">{f.message}</span>
										</li>
									{/each}
									{#if row.omitted}<li>… {row.omitted} more findings left out</li>{/if}
								</ul>
							{/if}
						</li>
					{/each}
				</ul>
				{#if verify.notes.length}
					<details class="mt-2" data-testid="verify-notes">
						<summary class="cursor-pointer text-muted">
							{verify.notes.length}
							{verify.notes.length === 1 ? 'note' : 'notes'} from flai check outside the story; they do
							not fail it
						</summary>
						<ul class="mt-1 ml-4 list-disc font-mono text-xs">
							{#each verify.notes as n, i (i)}
								<li>
									{#if n.where}<span class="font-medium">{n.where}:</span>{/if}
									<span class="whitespace-pre-wrap">{n.message}</span>
								</li>
							{/each}
						</ul>
					</details>
				{/if}
			{/if}
		</section>
	{/if}

	<section class="mt-4">
		<Threads on={item.id} {writable} />
	</section>
{/if}
