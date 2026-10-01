<script lang="ts">
	// A story's task plan on its page (S-0176): each task's state, a waiting task with the tasks it
	// waits for, and the layers of tasks that can run at once, first to last. An open task flai left
	// out of every layer, because its `after` runs in a cycle, is named under them.
	import { resolve } from '$app/paths';
	import { stateLabel, unplaced, waitsFor, type TaskPlan } from '$lib/taskplan';

	let {
		plan,
		titles = {}
	}: {
		plan: TaskPlan;
		/** each task's title by ID, from the story's children, for a link's tooltip */
		titles?: Record<string, string>;
	} = $props();

	const outside = $derived(unplaced(plan));
</script>

{#snippet task(id: string)}<a
		class="font-mono underline"
		href={resolve('/items/[id]', { id })}
		title={titles[id]}>{id}</a
	>{/snippet}
<!-- tasks linked, a comma after each but the last, spaced by the row's gap rather than by text -->
{#snippet links(ids: string[])}<span class="inline-flex flex-wrap gap-x-1"
		>{#each ids as id, i (id)}<span
				>{@render task(id)}{#if i < ids.length - 1},{/if}</span
			>{/each}</span
	>{/snippet}

<section class="rounded border border-line bg-surface p-3" data-testid="task-plan">
	<h2 class="mb-2 font-medium">Task plan</h2>
	<ul class="space-y-1 text-xs">
		{#each plan.tasks as t (t.id)}
			<li data-testid="plan-task" data-state={t.state}>
				{@render task(t.id)}
				<span class={t.state === 'waiting' ? 'text-warn' : 'text-muted'}>{stateLabel(t.state)}</span
				>
				{#if t.state === 'waiting' && waitsFor(t).length}<span class="text-muted"
						>for {@render links(waitsFor(t))}</span
					>{/if}
			</li>
		{/each}
	</ul>
	{#if plan.layers.length}
		<h3 class="mt-3 mb-1 font-medium">Layers</h3>
		<ol class="space-y-1 text-xs" aria-label="layers of tasks that can run at once">
			{#each plan.layers as layer, n (n)}
				<li data-testid="plan-layer">
					<span class="text-muted">layer {n + 1}:</span>
					{@render links(layer)}
				</li>
			{/each}
		</ol>
	{/if}
	{#if outside.length}
		<p class="mt-2 text-xs text-danger" data-testid="plan-cycle">
			in no layer, as what they wait for runs in a cycle:
			{@render links(outside)}
		</p>
	{/if}
</section>
