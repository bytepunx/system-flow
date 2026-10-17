<!--
	The calls of the framed component's function props (S-0440): the harness gives each function
	prop a logger, which posts an event with the prop's name and its arguments as JSON, or their
	type where they do not serialise. The log lists the last MAX_EVENTS, newest first, each with the
	time the page took it.
-->
<script lang="ts" module>
	/** One call of a logged function prop: when the page took it, the prop, and its arguments. */
	export type Logged = { at: number; name: string; args: unknown[] };

	/** How many calls the log keeps. */
	export const MAX_EVENTS = 50;

	/** events with call first, as the newest, and the oldest past MAX_EVENTS dropped. */
	export function logged(events: readonly Logged[], call: Logged): Logged[] {
		return [call, ...events].slice(0, MAX_EVENTS);
	}
</script>

<script lang="ts">
	let {
		events,
		onclear
	}: {
		/** The calls, newest first. */
		events: readonly Logged[];
		/** Called when the operator clears the log; without it there is no Clear. */
		onclear?: () => void;
	} = $props();

	const shown = $derived(events.slice(0, MAX_EVENTS));
	const pad = (n: number, size = 2) => String(n).padStart(size, '0');
	// The time of day in the browser's zone to the millisecond, as calls come close together.
	const clock = (at: number) => {
		const d = new Date(at);
		return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}.${pad(d.getMilliseconds(), 3)}`;
	};
</script>

<section aria-labelledby="event-log-heading" data-testid="event-log">
	<div class="flex items-center justify-between gap-2">
		<h3 id="event-log-heading" class="text-sm font-semibold">Events</h3>
		{#if onclear && events.length}
			<button
				type="button"
				class="rounded border border-line px-2 py-0.5 text-xs"
				onclick={() => onclear()}>Clear</button
			>
		{/if}
	</div>
	{#if !shown.length}
		<p class="mt-1 text-xs text-muted" data-testid="event-log-none">
			No calls yet. A function prop's calls are listed here, newest first.
		</p>
	{:else}
		<ol class="mt-1 max-h-64 overflow-auto rounded border border-line text-xs">
			{#each shown as e, i (i)}
				<li class="border-b border-line px-2 py-1 last:border-b-0" data-testid="event-log-entry">
					<time class="text-muted" datetime={new Date(e.at).toISOString()}>{clock(e.at)}</time>
					<span class="ml-1 font-mono font-medium">{e.name}</span>
					<code class="block font-mono break-all">{JSON.stringify(e.args)}</code>
				</li>
			{/each}
		</ol>
	{/if}
</section>
