<script module lang="ts">
	import type { KlusterLog } from '$lib/server/repo';

	/** A read of a log from an offset, as kluster.log answers it and a read prop must. */
	export type LogPart = Omit<KlusterLog, 'namespace' | 'step' | 'log'>;
</script>

<script lang="ts">
	// One step of a kluster job, a build, a load, a deploy, or a teardown, in a namespace (S-0595,
	// ADR-0263): its name, a badge of its outcome, how long it took, and, opened, its log as flai's
	// kluster.log reads it through /api/kluster/log. Given read, it shows a log that is not a kluster
	// step's, such as an image job's (S-0757), as read reads it, in a namespace of none. Opened, it
	// reads the log from the start, then from where the last read stopped, about every second while
	// the step runs, and once more when it ends, keeping the newest 256 KB and staying scrolled to the
	// end unless the reader scrolled up. A read that answers from another offset than asked is passed
	// over, but a reset, when a new job cleared the log, starts the text over. Closed or gone, it
	// reads nothing. A read flai refused or could not answer is shown word for word, and tried again
	// while the step runs.
	import { tick, untrack } from 'svelte';
	import { api } from '$lib/api';
	import { localTime } from '$lib/localtime';

	let {
		namespace = '',
		step,
		label = '',
		outcome = '',
		startedAt = '',
		endedAt = '',
		exitCode,
		open = $bindable(false),
		read
	}: {
		/** The namespace whose job ran the step, such as s-0595 or the default namespace; none with read. */
		namespace?: string;
		/**
		 * The step's name: build, build-<component>, load, deploy, or teardown; with read, what the
		 * log is of, such as the image job 20261113T093000Z-001. Another starts the log over.
		 */
		step: string;
		/** What the step is called on the page; the step's name when empty. */
		label?: string;
		/** pending, running, passed, failed, skipped, or empty when the record names none. */
		outcome?: string;
		startedAt?: string;
		endedAt?: string;
		/** The step's exit code, when it ended with one. */
		exitCode?: number;
		/** Whether the log is open; the reader opens and closes it. */
		open?: boolean;
		/** Reads the log from an offset in place of kluster.log; throws a refusal's words. */
		read?: (offset: number) => Promise<LogPart>;
	} = $props();

	/** How often a running step's log is read again. */
	const EVERY = 1000;
	/** How long after a failed read it is tried again, while the step runs. */
	const AFTER_ERROR = 5000;
	/** The most of the log kept on the page; the rest is in the log on the host. */
	const KEEP = 256 * 1024;

	const where = $derived(namespace);
	const which = $derived(step);
	const opened = $derived(open);
	const live = $derived(outcome === 'running');
	// what the log is called to a screen reader: a kluster step's in its namespace, or what read reads
	const named = $derived(
		namespace ? `the log of the ${label || step} step in ${namespace}` : `the log of ${step}`
	);

	let text = $state('');
	let next = 0;
	let earlier = $state(false);
	let exists = $state(true);
	let error = $state<string | null>(null);
	let box = $state<HTMLElement | null>(null);

	// another step's log starts over
	$effect.pre(() => {
		void where;
		void which;
		untrack(() => {
			text = '';
			next = 0;
			earlier = false;
			exists = true;
			error = null;
		});
	});

	$effect(() => {
		void where;
		void which;
		void live;
		if (!opened) return;
		return untrack(follow);
	});

	async function readStep(offset: number): Promise<LogPart> {
		if (read) return read(offset);
		const query = `namespace=${encodeURIComponent(where)}&step=${encodeURIComponent(which)}&offset=${offset}`;
		const r = await api(`/api/kluster/log?${query}`);
		if (!r.ok) {
			const message = (await r.json().catch(() => ({}))).error ?? r.statusText;
			throw new Error(message || `the log could not be read (${r.status})`);
		}
		return r.json();
	}

	async function append(got: LogPart, asked: number) {
		if (got.reset) {
			text = '';
			earlier = false;
		} else if (got.offset !== asked) {
			return;
		}
		const atEnd = !box || box.scrollHeight - box.scrollTop - box.clientHeight < 24;
		let all = text + got.text;
		if (all.length > KEEP) {
			const cut = all.indexOf('\n', all.length - KEEP);
			all = all.slice(cut < 0 ? all.length - KEEP : cut + 1);
			earlier = true;
		}
		text = all;
		next = got.next;
		await tick();
		if (atEnd && box) box.scrollTop = box.scrollHeight;
	}

	function follow() {
		let stopped = false;
		let timer: ReturnType<typeof setTimeout> | undefined;
		const later = (ms: number) => {
			if (!stopped) timer = setTimeout(pass, ms);
		};
		async function pass() {
			const asked = next;
			try {
				const got = await readStep(asked);
				if (stopped) return;
				error = null;
				exists = got.exists;
				await append(got, asked);
				if (stopped) return;
				if (got.next < got.size) later(0);
				else if (live || got.outcome === 'running') later(EVERY);
			} catch (e) {
				if (stopped) return;
				error = e instanceof Error ? e.message : String(e);
				if (live) later(AFTER_ERROR);
			}
		}
		void pass();
		return () => {
			stopped = true;
			clearTimeout(timer);
		};
	}

	const badge = (o: string) =>
		o === 'failed' || o === 'interrupted'
			? 'bg-danger-soft text-danger'
			: o === 'passed' || o === 'done'
				? 'bg-good-soft text-good-strong'
				: o === 'running'
					? 'bg-info-soft text-info'
					: 'bg-surface text-muted';

	function took(from: string, to: string): string {
		const ms = Date.parse(to) - Date.parse(from);
		if (!Number.isFinite(ms) || ms < 0) return '';
		const s = Math.round(ms / 1000);
		const m = Math.floor(s / 60);
		return m ? `${m}m ${String(s % 60).padStart(2, '0')}s` : `${s}s`;
	}
	const duration = $derived(
		startedAt && endedAt
			? took(startedAt, endedAt)
			: startedAt && outcome === 'running'
				? `since ${localTime(startedAt)}`
				: ''
	);
</script>

<details bind:open data-testid="log-feed" data-step={step} data-outcome={outcome}>
	<summary class="flex cursor-pointer flex-wrap items-center gap-2 text-sm select-none">
		<span class="font-medium">{label || step}</span>
		<span class="rounded px-2 py-0.5 text-xs {badge(outcome)}" data-testid="log-feed-status"
			>{outcome || 'not run'}</span
		>
		{#if duration}<span class="text-xs text-muted" data-testid="log-feed-duration">{duration}</span
			>{/if}
		{#if exitCode}<span class="text-xs text-muted">exit code {exitCode}</span>{/if}
		{#if live}<span class="text-xs text-good">· live</span>{/if}
	</summary>
	{#if error}
		<p class="mt-1 text-xs text-danger" role="alert" data-testid="log-feed-error">{error}</p>
	{/if}
	<pre
		bind:this={box}
		class="mt-1 max-h-72 overflow-auto rounded border border-line bg-ground p-2 font-mono text-xs whitespace-pre-wrap"
		aria-label={named}
		data-testid="log-feed-text">{#if earlier}… earlier lines are in the log on the host
		{/if}{text}</pre>
	{#if !text && !error}
		<p class="mt-1 text-xs text-muted" data-testid="log-feed-empty">
			{outcome === 'pending'
				? `${namespace ? 'This step' : 'It'} has not started yet.`
				: exists
					? 'Nothing written yet.'
					: `${namespace ? 'The step' : 'It'} has no log.`}
		</p>
	{/if}
</details>
