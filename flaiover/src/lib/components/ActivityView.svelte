<script lang="ts">
	// Who is working on what (S-0042), from the narratives, and since S-0142 what each story's agent
	// is saying and doing: a card whose story has had an agent flai serve started shows its stream,
	// open while the agent runs, and an agent at work on a story with no narrative yet gets a card of
	// its own.
	import { resolve } from '$app/paths';
	import { age } from '$lib/age';
	import { activityLine, type StoryActivity } from '$lib/activity';
	import AgentDot from './AgentDot.svelte';
	import AgentStream from './AgentStream.svelte';

	type Stream = {
		stream: string;
		title: string;
		agent: string;
		session: string;
		updated: string;
		age_seconds: number;
		status: string;
		blocked: boolean;
		task?: { id: string; title: string };
		last_log?: { at: string; text: string };
		path: string;
	};
	let { streams, agents = {} }: { streams: Stream[]; agents?: Record<string, StoryActivity> } =
		$props();

	const live = (a: StoryActivity | undefined) => a?.state === 'working' || a?.state === 'waiting';
	// An agent that has started, which a held story's stand-in run has not.
	const ran = (a: StoryActivity | undefined): a is StoryActivity => !!a?.run.started;
	const unnarrated = $derived(
		Object.entries(agents)
			.filter(([id, a]) => ran(a) && live(a) && !streams.some((s) => s.stream === id))
			.sort(([a], [b]) => a.localeCompare(b))
	);
</script>

{#if !streams.length && !unnarrated.length}
	<p class="text-sm text-muted">No active streams: no story has an open narrative.</p>
{/if}
<ul class="space-y-3">
	{#each streams as s (s.stream)}
		<li class="rounded border border-line bg-surface p-3 text-sm">
			<div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<a class="font-mono font-medium underline" href={resolve('/items/[id]', { id: s.stream })}
					>{s.stream}</a
				>
				<span class="min-w-0 flex-1">{s.title}</span>
				<span class="rounded border border-line px-2 py-0.5 text-xs">{s.status}</span>
				{#if s.blocked}<span class="text-xs font-semibold text-danger">BLOCKED</span>{/if}
			</div>
			<dl class="mt-2 grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1 text-xs">
				<dt class="text-muted">agent</dt>
				<dd>
					{s.agent || 'unknown'}{#if s.session}<span class="text-muted"
							>&nbsp;· session {s.session}</span
						>{/if}
				</dd>
				<dt class="text-muted">last wrote</dt>
				<dd title={s.updated}>{age(s.age_seconds)} ago</dd>
				<dt class="text-muted">task</dt>
				<dd>
					{#if s.task}<a class="underline" href={resolve('/items/[id]', { id: s.task.id })}
							>{s.task.id}</a
						>
						{s.task.title}{:else}<span class="text-muted">none in progress</span>{/if}
				</dd>
				<dt class="text-muted">last log</dt>
				<dd class="whitespace-pre-wrap">
					{#if s.last_log}<span class="text-muted">{s.last_log.at}</span>
						{s.last_log.text}{:else}<span class="text-muted">nothing logged</span>{/if}
				</dd>
			</dl>
			<p class="mt-2 text-xs">
				<a class="underline" href={resolve('/docs/[...path]', { path: s.path })}>narrative</a>
			</p>
			{#if ran(agents[s.stream])}
				{@const a = agents[s.stream]}
				<p class="mt-2 flex items-center gap-2 text-xs">
					<AgentDot activity={a} /><span class="text-muted">{activityLine(a)}</span>
				</p>
				<AgentStream story={s.stream} started={a.run.started} open={live(a)} />
			{/if}
		</li>
	{/each}
	{#each unnarrated as [id, a] (id)}
		<li class="rounded border border-line bg-surface p-3 text-sm" data-testid="unnarrated">
			<div class="flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<a class="font-mono font-medium underline" href={resolve('/items/[id]', { id })}>{id}</a>
				<span class="min-w-0 flex-1 text-muted">no narrative yet</span>
			</div>
			<p class="mt-2 flex items-center gap-2 text-xs">
				<AgentDot activity={a} /><span class="text-muted">{activityLine(a)} · {a.run.agent}</span>
			</p>
			<AgentStream story={id} started={a.run.started} />
		</li>
	{/each}
</ul>
