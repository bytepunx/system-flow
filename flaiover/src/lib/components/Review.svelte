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
		patchLines,
		readNdjson,
		sectionOf,
		verifyView,
		type PatchKind,
		type VerifyReport,
		type VerifyView
	} from '$lib/review';
	import { acceptLabel, issueRows, type Issue, type IssueRow } from '$lib/issues';
	import DiffView from '$lib/components/DiffView.svelte';
	import Threads from '$lib/components/Threads.svelte';
	import ForgeStatus, { type PullRequestStatus } from '$lib/components/ForgeStatus.svelte';
	import { render, enhance } from '$lib/markdown';
	import { localTime } from '$lib/localtime';
	import { themeState } from '$lib/theme.svelte';
	import { editable, editLabel, fileRef, linkSegments, type LinkSegment } from '$lib/filelinks';
	import { openEditor } from '$lib/editor.svelte';
	import { changedComponents, componentFile, diffPaths } from '$lib/explorer/changed';
	import type { Component, ComponentsCatalogue } from '$lib/server/repo';
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
		/** Paths the branch undoes of main's change: each blocks until restored or kept (ADR-0155). */
		undoes_main?: Undone[];
		/** What the pull request holds acceptance on, under pull-request acceptance, each also in blockers (ADR-0169). */
		forge_blockers?: ForgeBlocker[];
		/** The story's pull request's URL, when flai serve kept one. */
		pull_request?: string;
		plan?: {
			level: string;
			commits: string[];
			steps: Step[] | null;
			skipped?: string;
			unreleased?: { component: string; files: string[] }[];
		} | null;
	};
	/** A path the branch undoes of main's change, as flai's preview and refusal give it (ADR-0155). */
	type Undone = { path: string; lines: string[]; old_base: string };
	/** A forge blocker as flai's preview gives it: name is a check's or a reviewer's, url its page (ADR-0169). */
	type ForgeBlocker = { code: string; message: string; name?: string; url?: string };
	/** A file of the story's diff, as /diff gives it and DiffView shows it. */
	type DiffFile = {
		path: string;
		status: string;
		binary: boolean;
		truncated: boolean;
		patch: string;
	};
	type Progress = { step: string; msg: string };
	/** A finding of the lint the fix left, as flai's refusal of the acceptance gives it (ADR-0140). */
	type Finding = {
		tier: string;
		/** main, or story: the story's branch holds the file, and story names it. */
		checkout: string;
		story?: string;
		path: string;
		line?: number;
		column?: number;
		level?: string;
		rule?: string;
		alias?: string;
		message: string;
	};
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
	// A refusal for lint the fix left (S-0378): each finding, and the files saved in the editor since,
	// by fileKey. Accept is on again once each file the editor can open has been saved.
	let findings = $state<Finding[]>([]);
	let saved = $state<string[]>([]);

	let sendingBack = $state(false);
	let reason = $state('');

	// Paths the branch undoes of main's change (S-0388, ADR-0155): a card each, from the preview or
	// a refusal. Each is restored, and drops out when the preview runs again; or kept with a reason,
	// by path, which Accept passes as keep_removals; or the story goes back naming them. keeping
	// holds the reason typed so far on a card whose Keep is open.
	let undone = $state<Undone[]>([]);
	let kept = $state<Record<string, string>>({});
	let keeping = $state<Record<string, string>>({});
	let restoring = $state<string | null>(null);
	let undoneErrors = $state<Record<string, string>>({});
	let sendingUndone = $state(false);
	let undoneError = $state<string | null>(null);

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

	// The pull request (S-0394): the story's, its required checks and its reviews, as flai serve
	// last read them from the forge; null when it has none or it cannot be read, and then the
	// section is left out.
	let pullRequest = $state<PullRequestStatus | null>(null);

	// Components (S-0442): the catalogue as the story's worktree has it, so a component the branch
	// adds counts; null while unread, or when it cannot be read, and then nothing is offered.
	let catalogue = $state.raw<ComponentsCatalogue | null>(null);

	function sleep(ms: number): Promise<void> {
		return new Promise((r) => setTimeout(r, ms));
	}

	const v = (x: Version) => (typeof x === 'string' ? x : `${x.Major}.${x.Minor}.${x.Patch}`);
	const criteria = $derived(item ? criteriaOf(item.body) : []);
	const ticked = $derived(criteria.filter((c) => c.checked).length);
	const inReview = $derived(item?.status === 'review');
	const needsChoice = $derived(!!preview?.uncommitted?.length && !include);
	const chosen = $derived(issues.filter((r) => r.checked));
	// A main-checkout finding is saved without a story; a story's on its branch.
	const storyOf = (f: Finding) => (f.checkout === 'story' ? f.story : undefined);
	const fileKey = (f: Finding) => `${storyOf(f) ?? ''}:${f.path}`;
	// The editor opens a markdown file under the layout's folders that flai lets it edit.
	const opens = (f: Finding) => !!fileRef(f.path) && editable(f.path);
	// Another refusal's text, each file it names opening the editor on the story's branch.
	const refusal = $derived(failure ? linkSegments(failure, { story: id }) : []);
	const unsaved = $derived(
		[...new Set(findings.filter(opens).map(fileKey))].filter((k) => !saved.includes(k))
	);
	// an undone path neither kept nor yet restored holds Accept; its blocker is left to the card
	const unresolved = $derived(undone.filter((u) => !kept[u.path]));
	// The forge blockers are listed by kind, so the general list leaves out their messages.
	const forgeBlockers = $derived(preview?.forge_blockers ?? []);
	const blockers = $derived.by(() => {
		if (!preview) return [];
		const forge = new Set(forgeBlockers.map((b) => b.message));
		return otherBlockers(preview).filter((b) => !forge.has(b));
	});
	const canAccept = $derived(
		writable &&
			inReview &&
			!running &&
			!result &&
			!!preview &&
			!blockers.length &&
			!forgeBlockers.length &&
			!preview.worktree_uncommitted?.length &&
			!unresolved.length &&
			!needsChoice &&
			!unsaved.length
	);
	const undoneActive = $derived(
		writable && inReview && !result && !running && !restoring && !sendingUndone
	);

	// a line of the branch's change, as DiffView colours it
	const patchLine: Record<PatchKind, string> = {
		add: 'bg-good-soft text-good',
		del: 'bg-danger-soft text-danger',
		hunk: 'text-muted',
		context: '',
		note: 'text-muted'
	};
	// a forge blocker's kind in plain words, and what its url opens (ADR-0169)
	const forgeKind: Record<string, string> = {
		pull_request_missing: 'No pull request',
		forge_stale: 'Not read at the head',
		check_failing: 'Check failing',
		check_pending: 'Check pending',
		changes_requested: 'Changes requested',
		approvals_missing: 'Approvals missing'
	};
	const forgeLink = (code: string) =>
		code.startsWith('check_') ? 'log' : code === 'changes_requested' ? 'review' : 'pull request';
	// The catalogue's components the branch changes, from the diff the page already holds.
	const changedList = $derived.by((): Component[] => {
		if (!catalogue || !diff) return [];
		const changed = changedComponents(catalogue, diffPaths(diff));
		return catalogue.projects.flatMap((p) =>
			p.components.filter((c) => changed.has(componentFile(c)))
		);
	});
	const componentsHref = $derived(
		`${resolve('/design/components')}?story=${encodeURIComponent(id)}`
	);
	const branchFile = (path: string) =>
		((diff?.files ?? []) as DiffFile[]).find((f) => f.path === path);

	function findingsOf(data: unknown): Finding[] {
		const list = (data as { findings?: unknown } | undefined)?.findings;
		return Array.isArray(list) ? (list as Finding[]) : [];
	}

	function undoneOf(data: unknown): Undone[] {
		const list = (data as { undoes_main?: unknown } | undefined)?.undoes_main;
		return Array.isArray(list) ? (list as Undone[]) : [];
	}

	function openFinding(f: Finding) {
		const key = fileKey(f);
		openEditor({
			path: f.path,
			line: f.line,
			story: storyOf(f),
			onsaved: () => {
				if (!saved.includes(key)) saved = [...saved, key];
			}
		});
	}

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
		findings = [];
		saved = [];
		issues = [];
		issuesError = null;
		made = [];
		unmade = [];
		verify = null;
		verifyError = null;
		pullRequest = null;
		catalogue = null;
		undone = [];
		kept = {};
		keeping = {};
		undoneErrors = {};
		undoneError = null;
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
		void loadDiff(target);
		if (item.status === 'in-progress' || item.status === 'review') void loadCatalogue(target);
		if (item.status === 'review') {
			void loadPreview(target);
			void loadIssues(target);
		}
		void loadChecks(target);
		void loadVerify(target);
		if (item.type === 'story') void loadPullRequest(target);
	}

	async function loadDiff(target: string) {
		try {
			const d = await get<Record<string, unknown>>(`/api/items/${target}/diff`);
			if (target === id) diff = d;
		} catch (e) {
			if (target === id) diffError = e instanceof Error ? e.message : String(e);
		}
	}

	// The worktree's catalogue, for Open in Components. A refusal (an older flai, no Svelte
	// sub-project, no worktree) leaves catalogue null: the page offers nothing, and says nothing.
	async function loadCatalogue(target: string) {
		try {
			const r = await api(`/api/components?story=${encodeURIComponent(target)}`);
			const body = await r.json().catch(() => null);
			const read = (body as { catalogue?: ComponentsCatalogue } | null)?.catalogue;
			if (target === id && r.ok && Array.isArray(read?.projects)) catalogue = read;
		} catch {
			// the same: a failed read offers nothing rather than an error beside the diff
		}
	}

	// The preview's undone paths are the cards: a path restored since drops out (S-0388).
	async function loadPreview(target: string) {
		try {
			const p = await get<Preview>(`/api/items/${target}/acceptance`);
			if (target !== id) return;
			preview = p;
			undone = p.undoes_main ?? [];
		} catch (e) {
			if (target === id) previewError = e instanceof Error ? e.message : String(e);
		}
	}

	// flai's stream.restore commits main's version of the path onto the branch; the preview and the
	// diff are read again, so its card clears and the branch's changes are as they now are.
	async function restore(path: string) {
		restoring = path;
		delete undoneErrors[path];
		try {
			const r = await api(`/api/items/${id}/restore`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ path })
			});
			const data = await r.json().catch(() => ({}));
			if (!r.ok) {
				undoneErrors[path] = data.error ?? r.statusText;
				return;
			}
			await Promise.all([loadPreview(id), loadDiff(id)]);
		} catch (e) {
			undoneErrors[path] = e instanceof Error ? e.message : String(e);
		} finally {
			restoring = null;
		}
	}

	function keep(path: string) {
		const why = keeping[path]?.trim();
		if (!why) return;
		kept[path] = why;
		delete keeping[path];
	}

	// Two calls in order: a thread on the story naming the paths, for its agent; then the move back,
	// its reason naming them too. The page is loaded again as a send-back does.
	async function sendBackUndone() {
		const paths = undone.map((u) => u.path).join(', ');
		const branch = preview?.branch ?? `story/${id}`;
		sendingUndone = true;
		undoneError = null;
		try {
			const t = await api('/api/threads', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					on: id,
					title: `${branch} undoes main's change to ${paths}`,
					text: `${branch} removes lines main added since the branch left it, in ${paths}. Restore main's version with flai stream restore ${id} --path <path>, or redo the change, then move ${id} to review again (ADR-0155).`
				})
			});
			if (!t.ok) {
				const data = await t.json().catch(() => ({}));
				undoneError = `No thread was opened, and ${id} stays in review: ${data.error ?? t.statusText}`;
				return;
			}
			const m = await api(`/api/items/${id}/move`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					to: 'in-progress',
					reason: `${branch} undoes main's change to ${paths}: restore main's version or redo the change (ADR-0155)`
				})
			});
			if (!m.ok) {
				const data = await m.json().catch(() => ({}));
				undoneError = `The thread is open, but ${id} was not moved back: ${data.error ?? m.statusText}`;
				return;
			}
			await load(id);
		} catch (e) {
			undoneError = e instanceof Error ? e.message : String(e);
		} finally {
			sendingUndone = false;
		}
	}

	async function loadVerify(target: string) {
		try {
			const { report } = await get<{ report: VerifyReport | null }>(`/api/items/${target}/verify`);
			if (target === id) verify = report ? verifyView(report) : null;
		} catch (e) {
			if (target === id) verifyError = e instanceof Error ? e.message : String(e);
		}
	}

	async function loadPullRequest(target: string) {
		try {
			const data = await get<{ pull_request?: PullRequestStatus | null }>(
				`/api/items/${target}/forge`
			);
			if (target === id) pullRequest = data.pull_request ?? null;
		} catch {
			// no pull request read, or no flai connected: the page leaves the section out
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
		findings = [];
		saved = [];
		progress = [];
		warnings = [];
		const keeps = undone
			.filter((u) => kept[u.path])
			.map((u) => ({ path: u.path, reason: kept[u.path] }));
		try {
			const r = await api(`/api/items/${id}/accept`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					include_uncommitted: include || undefined,
					keep_removals: keeps.length ? keeps : undefined
				})
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
				else if (l.event === 'error') {
					failure = String(l.error);
					findings = findingsOf(l.data);
					// a path only the acceptance's own rebase found joins the cards (S-0388)
					const refused = undoneOf(l.data).filter((u) => !undone.some((x) => x.path === u.path));
					if (refused.length) undone = [...undone, ...refused];
				}
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

<!-- A refusal's text, each file it names a control that opens the editor at its line (S-0378). -->
{#snippet linked(segments: LinkSegment[])}
	{#each segments as s, i (i)}
		{#if s.edit}
			<button
				type="button"
				class="cursor-pointer underline"
				aria-label={s.edit.label}
				onclick={s.edit.open}>{s.text}</button
			>
		{:else}
			{s.text}
		{/if}
	{/each}
{/snippet}

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
		{#if changedList.length}
			<p class="mb-2" data-testid="open-in-components">
				It changes the components
				<span class="font-medium">{changedList.map((c) => c.name).join(', ')}</span>.
				<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; only ?story= is added -->
				<a class="underline" href={componentsHref}>Open in Components</a>
			</p>
		{/if}
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
				<!-- Each path the branch undoes of main's change, to restore or keep (S-0388, ADR-0155). -->
				{#each undone as u (u.path)}
					{@const file = branchFile(u.path)}
					<div
						class="mb-2 rounded border p-2 {kept[u.path]
							? 'border-warn bg-warn-soft'
							: 'border-danger bg-danger-soft'}"
						role="group"
						aria-label="The branch undoes main's change to {u.path}"
						data-testid="undoes-main"
					>
						<p class="font-medium {kept[u.path] ? 'text-warn' : 'text-danger'}">
							<code>{preview.branch ?? `story/${item.id}`}</code> undoes main's change to
							<span class="font-mono">{u.path}</span>: it removes {u.lines.length === 1
								? '1 line'
								: `${u.lines.length} lines`} main added since the branch left it at
							<code>{u.old_base.slice(0, 12)}</code>.
						</p>
						<div class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-2">
							<div class="min-w-0">
								<h3 class="text-xs text-muted">Main's lines the branch removes</h3>
								<pre
									class="mt-1 overflow-x-auto rounded bg-ground p-1 font-mono text-xs leading-snug text-danger"
									data-testid="undoes-main-lines">{u.lines.join('\n')}</pre>
							</div>
							<div class="min-w-0" data-testid="undoes-main-branch">
								<h3 class="text-xs text-muted">The branch's change to it</h3>
								{#if diffError}
									<p class="mt-1 text-muted">{diffError}</p>
								{:else if !diff}
									<p class="mt-1 text-muted">Loading…</p>
								{:else if !file}
									<p class="mt-1">The branch's version of this path is not in its diff.</p>
								{:else if !file.patch}
									<p class="mt-1 text-muted">
										{file.binary
											? 'Binary file, no hunks.'
											: file.truncated
												? 'Left out for size.'
												: 'No textual changes.'}
									</p>
								{:else}
									<div
										class="mt-1 overflow-x-auto rounded bg-ground p-1 font-mono text-xs leading-snug"
									>
										<div class="w-max min-w-full">
											{#each patchLines(file.patch) as l, n (n)}
												<div class="flex {patchLine[l.kind]}" data-line={l.kind}>
													<span class="w-5 shrink-0 text-center select-none">{l.sign}</span><span
														class="pr-2 whitespace-pre">{l.text || ' '}</span
													>
												</div>
											{/each}
										</div>
									</div>
								{/if}
							</div>
						</div>
						{#if kept[u.path]}
							<p class="mt-2" data-testid="undoes-main-kept">
								Kept: {kept[u.path]}
								{#if undoneActive}<button
										type="button"
										class="ml-2 rounded border border-line-strong bg-surface px-2 py-0.5 text-ink"
										onclick={() => delete kept[u.path]}>Undo</button
									>{/if}
							</p>
						{:else if keeping[u.path] !== undefined}
							<div class="mt-2">
								<label class="block text-xs text-muted" for="keep-{u.path}"
									>Why the removal is kept (recorded in the story's notes)</label
								>
								<input
									id="keep-{u.path}"
									type="text"
									class="mt-1 w-full rounded border border-line bg-ground p-1"
									bind:value={keeping[u.path]}
								/>
								<button
									type="button"
									class="mt-1 rounded border border-line-strong bg-surface px-2 py-0.5 text-ink disabled:opacity-50"
									disabled={!keeping[u.path].trim() || !undoneActive}
									onclick={() => keep(u.path)}>Keep the removal</button
								>
								<button
									type="button"
									class="mt-1 rounded border border-line bg-surface px-2 py-0.5 text-ink"
									onclick={() => delete keeping[u.path]}>Cancel</button
								>
							</div>
						{:else if writable && inReview && !result}
							<p class="mt-2 flex flex-wrap gap-2">
								<button
									type="button"
									class="rounded border border-line-strong bg-surface px-2 py-0.5 text-ink disabled:opacity-50"
									disabled={!undoneActive}
									onclick={() => restore(u.path)}
									>{restoring === u.path ? 'Restoring…' : "Restore main's version"}</button
								>
								<button
									type="button"
									class="rounded border border-line-strong bg-surface px-2 py-0.5 text-ink disabled:opacity-50"
									disabled={!undoneActive}
									onclick={() => (keeping[u.path] = '')}>Keep</button
								>
								<button
									type="button"
									class="rounded border border-line-strong bg-surface px-2 py-0.5 text-ink disabled:opacity-50"
									disabled={!undoneActive}
									onclick={sendBackUndone}
									>{sendingUndone ? 'Sending back…' : 'Send back, naming the undone paths'}</button
								>
							</p>
						{/if}
						{#if undoneErrors[u.path]}
							<p class="mt-2 text-danger" role="alert">
								{u.path} was not restored. flai said: {undoneErrors[u.path]}
							</p>
						{/if}
					</div>
				{/each}
				{#if undone.length && writable && inReview && !result}
					<div class="mb-2">
						{#if unresolved.length}
							<p class="text-xs text-muted">
								Accept is on once each of these is restored or kept; or send {item.id} back for its agent
								to put right.
							</p>
						{/if}
						{#if undoneError}<p class="mt-1 text-danger" role="alert">{undoneError}</p>{/if}
					</div>
				{/if}
				<!-- What the pull request holds acceptance on, by kind, under pull-request acceptance (ADR-0169). -->
				{#if forgeBlockers.length}
					<div
						class="mb-2 rounded border border-danger bg-danger-soft p-2 text-danger"
						role="alert"
						data-testid="forge-blockers"
					>
						<p class="font-medium">
							{#if preview.pull_request}The <a
									class="underline"
									href={preview.pull_request}
									target="_blank"
									rel="external noopener noreferrer">pull request</a
								>{:else}The pull request{/if} holds this acceptance:
						</p>
						<ul class="mt-1 ml-4 list-disc space-y-1">
							{#each forgeBlockers as b (b.message)}
								<li data-code={b.code}>
									<span class="flex flex-wrap items-baseline gap-2">
										<span class="font-medium">{forgeKind[b.code] ?? b.code}</span>
										{#if b.name}<span class="font-mono">{b.name}</span>{/if}
										{#if b.url}<a
												class="text-xs underline"
												href={b.url}
												target="_blank"
												rel="external noopener noreferrer"
												aria-label="{b.name ?? forgeKind[b.code] ?? b.code} {forgeLink(b.code)}"
												>{forgeLink(b.code)}</a
											>{/if}
									</span>
									<p class="text-xs">{b.message}</p>
								</li>
							{/each}
						</ul>
					</div>
				{/if}
				{#if blockers.length}
					<div
						class="mb-2 rounded border border-danger bg-danger-soft p-2 text-danger"
						role="alert"
					>
						<p class="font-medium">This cannot be accepted from here yet:</p>
						<ul class="mt-1 ml-4 list-disc">
							{#each blockers as b (b)}<li>{b}</li>{/each}
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
					{#if preview.pull_request && !forgeBlockers.length}
						<li data-testid="accept-pull-request">
							Its <a
								class="underline"
								href={preview.pull_request}
								target="_blank"
								rel="external noopener noreferrer">pull request</a
							> passes its required checks and reviews.
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
			{#if failure && findings.length}
				<div
					class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger"
					role="alert"
					data-testid="accept-findings"
				>
					<p class="font-medium">Not accepted. {item.id} is {item.status}.</p>
					<p>The fix left lint in this markdown: correct each file at its line and save it.</p>
					<p>Accept is on again once each file is saved.</p>
					<ul class="mt-1 space-y-1 font-mono text-xs">
						{#each findings as f, i (i)}
							{@const where = f.line ? `${f.path}:${f.line}` : f.path}
							<li data-testid="accept-finding">
								{#if opens(f)}
									<button
										type="button"
										class="cursor-pointer underline"
										aria-label={editLabel({ path: f.path, line: f.line })}
										onclick={() => openFinding(f)}>{where}</button
									>
								{:else}
									<span class="font-medium">{where}</span>
								{/if}
								{#if storyOf(f)}<span class="text-muted">on {storyOf(f)}'s branch</span>{/if}
								{#if f.rule}<span>{f.alias ? `${f.rule}/${f.alias}` : f.rule}</span>{/if}
								<span class="whitespace-pre-wrap">{f.message}</span>
								{#if saved.includes(fileKey(f))}
									<span class="text-good" data-testid="finding-saved">saved</span>
								{/if}
							</li>
						{/each}
					</ul>
				</div>
			{:else if failure}
				<div class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
					<p class="font-medium">Not accepted. {item.id} is {item.status}. flai said:</p>
					<pre class="mt-1 font-mono text-xs whitespace-pre-wrap">{@render linked(refusal)}</pre>
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
						Running {checksRun.current}, started {localTime(checksRun.started)}.
					{:else}
						{checksRun.outcome}, started {localTime(checksRun.started)}, ended
						{localTime(checksRun.ended ?? '')}.
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

	<ForgeStatus status={pullRequest} class="mt-4" />

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
					<code class="rounded bg-ground px-1 text-ink">{verify.base}</code> at
					{localTime(verify.ranAt)}, in {verify.duration}.
				</p>
				<ul class="mt-2 space-y-1">
					{#each verify.rows as row, i (i)}
						<li data-testid="verify-step">
							<span class="flex flex-wrap items-baseline gap-2">
								<span
									class={row.state === 'passed'
										? 'text-good'
										: row.state === 'reused'
											? 'text-good opacity-70'
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
