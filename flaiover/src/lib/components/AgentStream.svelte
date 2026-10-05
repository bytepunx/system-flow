<script lang="ts">
	// What a story's agent is saying and doing (S-0142), as flai on the host reads it from the log it
	// gave the agent: the newest entries, kept scrolled to the end unless the reader scrolled up, and
	// read again every two seconds from where the last read stopped while the agent runs. A new run
	// for the story (another `started`) starts it over from the tail; the same story and run passed
	// again change nothing (S-0178). When flai cannot answer, it says why and tries again now and
	// then, except for a story flai started no agent for. It opens when told to and closes only when
	// the reader closes it, so an agent that ends does not shorten the page under the reader (S-0178).
	// With `plan` it follows the newest planner run for the epic or story named by `story` instead,
	// as flai serve reads it from the planner's log (S-0259).
	import { tick, untrack } from 'svelte';
	import { api } from '$lib/api';
	import type { AgentStreamEntry, AgentStreamRead } from '$lib/activity';

	let {
		story,
		started,
		open = true,
		plan = false
	}: {
		/** The story whose agent it follows, or with `plan` the epic or story the planner plans. */
		story: string;
		started: string;
		/** Opens the stream when it is or becomes true; becoming false leaves it as the reader has it. */
		open?: boolean;
		/** Follows the planner's run for the item rather than a story's agent (S-0259). */
		plan?: boolean;
	} = $props();

	// The page passes these from objects it makes anew on every reload. Derived, they change only
	// when the values do, so a reload does not empty the stream and move the page (S-0178).
	const following = $derived(story);
	const since = $derived(started);
	const opened = $derived(open);
	const planner = $derived(plan);

	/** Entries kept on the page; older ones are in the log on the host. */
	const KEEP = 500;
	const EVERY = 2000;
	const AFTER_ERROR = 10000;

	let entries = $state<AgentStreamEntry[]>([]);
	let running = $state(false);
	let earlier = $state(false);
	let error = $state<string | null>(null);
	let box = $state<HTMLElement | null>(null);
	let shown = $state(false);

	// Opens when told to; closing is the reader's, even when the agent ends (S-0178).
	$effect.pre(() => {
		if (opened) shown = true;
	});

	/** A read flai refused or could not answer, with the dashboard's status for it. */
	type ReadError = Error & { status: number };

	async function read(after?: number): Promise<AgentStreamRead> {
		const query = [planner ? 'plan' : '', after === undefined ? '' : `after=${after}`]
			.filter(Boolean)
			.join('&');
		const r = await api(
			`/api/agent-stream/${encodeURIComponent(following)}${query ? `?${query}` : ''}`
		);
		if (!r.ok) {
			const message = (await r.json().catch(() => ({}))).error ?? r.statusText;
			throw Object.assign(new Error(message), { status: r.status }) as ReadError;
		}
		return r.json();
	}

	async function append(got: AgentStreamRead) {
		const atEnd = !box || box.scrollHeight - box.scrollTop - box.clientHeight < 24;
		if (got.skipped) earlier = true;
		let all = [...entries, ...got.entries];
		if (all.length > KEEP) {
			all = all.slice(-KEEP);
			earlier = true;
		}
		entries = all;
		await tick();
		if (atEnd && box) box.scrollTop = box.scrollHeight;
	}

	$effect(() => {
		void following;
		void since;
		void planner;
		// nothing else read here starts the stream over
		return untrack(follow);
	});

	function follow() {
		let stopped = false;
		let timer: ReturnType<typeof setTimeout> | undefined;
		let next: number | undefined;
		let run = '';
		entries = [];
		earlier = false;
		error = null;
		const later = (ms: number) => {
			if (!stopped) timer = setTimeout(step, ms);
		};
		async function step() {
			try {
				const got = await read(next);
				if (stopped) return;
				if (run && got.started !== run) {
					// another run: its log is another file
					run = '';
					next = undefined;
					entries = [];
					earlier = false;
					later(0);
					return;
				}
				if (!run && got.from > 0) earlier = true;
				run = got.started;
				next = got.next;
				running = got.running;
				error = null;
				await append(got);
				if (got.more) later(0);
				else if (got.running) later(EVERY);
			} catch (e) {
				if (stopped) return;
				error = e instanceof Error ? e.message : String(e);
				if ((e as Partial<ReadError>).status !== 404) later(AFTER_ERROR);
			}
		}
		void step();
		return () => {
			stopped = true;
			clearTimeout(timer);
		};
	}

	function tone(e: AgentStreamEntry): string {
		if (e.error) return 'text-danger';
		switch (e.kind) {
			case 'tool':
				return 'text-accent';
			case 'result':
			case 'output':
			case 'session':
			case 'task':
				return 'text-muted';
			case 'thinking':
				return 'text-muted italic';
			default:
				return '';
		}
	}

	function label(e: AgentStreamEntry): string {
		switch (e.kind) {
			case 'tool':
				return `▸ ${e.tool ?? 'tool'}`;
			case 'result':
				return '↳';
			case 'text':
				return '';
			default:
				return e.kind;
		}
	}
</script>

<details class="mt-2" bind:open={shown} data-testid="agent-stream" data-running={running}>
	<summary class="cursor-pointer text-xs text-muted select-none">
		stream{#if running}<span class="text-good">&nbsp;· live</span>{/if}
	</summary>
	{#if error}
		<p class="mt-1 text-xs text-danger" role="alert">{error}</p>
	{/if}
	<ol
		bind:this={box}
		class="mt-1 max-h-72 space-y-0.5 overflow-y-auto rounded border border-line bg-ground p-2 font-mono text-xs"
		aria-label={plan
			? `what the planner for ${story} said and did`
			: `what ${story}'s agent said and did`}
	>
		{#if earlier}
			<li class="text-muted">
				… earlier entries are in the {plan ? "planner's" : "agent's"} log on the host
			</li>
		{/if}
		{#each entries as e, i (i)}
			<li class="break-words whitespace-pre-wrap {tone(e)}" data-kind={e.kind}>
				{#if label(e)}<span class="font-semibold">{label(e)}</span>&nbsp;{/if}{e.text}
			</li>
		{/each}
		{#if !entries.length && !error}
			<li class="text-muted">nothing written yet</li>
		{/if}
	</ol>
</details>
